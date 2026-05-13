package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

// To process the operation IDs for addReaction, removeReaction, and getReactions,
// we first need to verify that the authenticated user has the privileges to
// perform these operations. This function accomplishes exactly that.
//
// Actions performed by this function:
// - Check message validity.
// - Check user privileges for the message.
func (rt *_router) _authorizeReactionOperation(
	w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User,
) (*database.Message, error) {
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter messageId must be an int64", http.StatusBadRequest)
		return nil, err
	}

	message, err := rt.db.GetMessage(messageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return nil, err
	case err != nil:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
		return nil, err
	}

	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %w", err)
		return nil, err
	case !isMember:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
		return nil, err
	}

	return message, nil
}

// operationId: addReaction
func (rt *_router) addReaction(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt._authorizeReactionOperation(w, r, ps, user)
	if err != nil {
		return
	}

	emojiCode, err := rt.checkEmojiCode(w, r)
	if err != nil {
		return
	}

	reaction, err := rt.db.GetReaction(message.MessageId, user.UserId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = rt.db.AddReaction(*emojiCode, message.MessageId, user.UserId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
			rt.baseLogger.Errorf("AddReaction: %w", err)
			return
		}
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetReaction: %w", err)
		return
	default:
		if reaction.EmojiCode == *emojiCode {
			rt.sendResponse(w, "The same emoji code was provied", http.StatusBadRequest)
			return
		}
		_, err = rt.db.UpdateReaction(*emojiCode, message.MessageId, user.UserId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
			rt.baseLogger.Errorf("UpdateReaction: %w", err)
			return
		}
	}

	reaction, err = rt.db.GetReaction(message.MessageId, user.UserId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetReaction: %w", err)
	default:
		rt.sendResponse(w, reaction, http.StatusCreated)
	}
}

// operationId: deleteReaction
func (rt *_router) deleteReaction(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt._authorizeReactionOperation(w, r, ps, user)
	if err != nil {
		return
	}

	_, err = rt.db.GetReaction(message.MessageId, user.UserId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(
			w,
			fmt.Sprintf("Reaction not found for user %d and message %d", user.UserId, message.MessageId),
			http.StatusBadRequest,
		)
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetReaction: %w", err)
	default:
		_, err = rt.db.DeleteReaction(message.MessageId, user.UserId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
			rt.baseLogger.Errorf("DeleteReaction: %w", err)
			return
		}

		rt.sendResponse(w, nil, http.StatusNoContent)
	}
}

// operationId: getReactions
func (rt *_router) getReactions(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt._authorizeReactionOperation(w, r, ps, user)
	if err != nil {
		return
	}

	reactions, err := rt.db.GetReactions(message.MessageId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetReactions: %w", err)
		return
	}

	rt.sendResponse(w, struct {
		Reactions []database.Reaction `json:"reactions"`
	}{Reactions: reactions}, http.StatusOK)
}
