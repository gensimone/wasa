package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

// operationId: addReaction
func (rt *_router) addReaction(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
	if err != nil {
		return
	}

	emoji, err := rt.getEmojiFromReq(w, r)
	if err != nil {
		return
	}

	isValidEmoji := database.IsValidEmoji(*emoji)
	if !isValidEmoji {
		rt.sendResponse(w, "Invalid emoji", http.StatusBadRequest)
		return
	}

	reaction, err := rt.db.GetReaction(message.MessageId, user.UserId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = rt.db.AddReaction(*emoji, message.MessageId, user.UserId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
			rt.baseLogger.Errorf("AddReaction: %v", err)
			return
		}
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetReaction: %v", err)
		return
	default:
		if reaction.Emoji == *emoji {
			rt.sendResponse(w, reaction, http.StatusCreated)
			return
		}

		_, err = rt.db.UpdateReaction(*emoji, message.MessageId, user.UserId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
			rt.baseLogger.Errorf("UpdateReaction: %v", err)
			return
		}
	}

	reaction, err = rt.db.GetReaction(message.MessageId, user.UserId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetReaction: %v", err)
	default:
		rt.sendResponse(w, reaction, http.StatusCreated)
	}
}

// operationId: deleteReaction
func (rt *_router) deleteReaction(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
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
		rt.baseLogger.Errorf("GetReaction: %v", err)
	default:
		_, err = rt.db.DeleteReaction(message.MessageId, user.UserId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
			rt.baseLogger.Errorf("DeleteReaction: %v", err)
			return
		}

		rt.sendResponse(w, nil, http.StatusNoContent)
	}
}

// operationId: getReactions
func (rt *_router) getReactions(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
	if err != nil {
		return
	}

	reactions, err := rt.db.GetReactions(message.MessageId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetReactions: %v", err)
		return
	}

	rt.sendResponse(w, struct {
		Reactions []database.Reaction `json:"reactions"`
	}{Reactions: reactions}, http.StatusOK)
}
