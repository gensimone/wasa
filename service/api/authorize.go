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

// Authorize the request provided by the client.
// A request is authorized when the 'Authorization' header field
// contains a valid User ID. If the extract User ID is valid, the
// wrapper function is invoked.
func (rt *_router) authRequest(
	fn func(http.ResponseWriter, *http.Request, httprouter.Params, database.User),
) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := strconv.ParseInt(auth, 10, 64)
		if err != nil {
			rt.sendResponse(w, "User id must be an int64", http.StatusBadRequest)
			return
		}

		user, err := rt.db.GetUserById(userId)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			rt.sendResponse(w, fmt.Sprintf("User %d not found", userId), http.StatusNotFound)
		case err != nil:
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetUserById: %w", userId, err)
		default:
			fn(w, r, ps, *user)
		}
	}
}

// To process the operation IDs for addReaction, removeReaction, and getReactions,
// we first need to verify that the authenticated user has the privileges to
// perform these operations. This function accomplishes exactly that.
//
// Actions performed by this function:
// - Check message validity.
// - Check user privileges for the message.
func (rt *_router) authMessageAccess(
	w http.ResponseWriter, ps httprouter.Params, user database.User,
) (*database.Message, error) {
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter messageId must be an int64", http.StatusBadRequest)
		return nil, err
	}

	message, err := rt.db.GetMessage(messageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return nil, err
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
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
		errMsg := fmt.Sprintf("Unauthorized to access message %d", messageId)
		rt.sendResponse(w, errMsg, http.StatusUnauthorized)
		return nil, errors.New(errMsg)
	}

	return message, nil
}

func (rt *_router) authUserAsFounder(w http.ResponseWriter, groupId int64, userId int64) (bool, error) {
	isFounder, err := rt.db.IsFounder(groupId, userId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsFounder: %w", err)
		return false, err
	case isFounder:
		return true, nil
	default:
		rt.sendResponse(w, "Unauthorized to perform previliged actions on group.", http.StatusUnauthorized)
		return false, nil
	}
}
