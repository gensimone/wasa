package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

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
			rt.sendResponse(w, "User not found", http.StatusNotFound)
		case err != nil:
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetUserById: %v", userId, err)
		default:
			fn(w, r, ps, *user)
		}
	}
}

func (rt *_router) authConversationAccessParam(
	w http.ResponseWriter, ps httprouter.Params, user database.User,
) (*int64, error) {
	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter conversationId must be an int64", http.StatusBadRequest)
		return nil, err
	}

	return rt.authConversationAccess(w, user, conversationId)
}

func (rt *_router) authConversationAccess(
	w http.ResponseWriter, user database.User, conversationId int64,
) (*int64, error) {
	isMember, err := rt.db.IsMember(user.UserId, conversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
	case !isMember:
		errMsg := "Unauthorized to access conversation."
		rt.sendResponse(w, errMsg, http.StatusUnauthorized)
		err = errors.New(errMsg)
	}

	return &conversationId, err
}

func (rt *_router) authMessageAccessParam(
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
		rt.sendResponse(w, "Message not found", http.StatusNotFound)
		return nil, err
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %v", err)
		return nil, err
	}

	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %v", err)
		return message, err
	case !isMember:
		errMsg := "Unauthorized to access message"
		rt.sendResponse(w, errMsg, http.StatusUnauthorized)
		return message, errors.New(errMsg)
	default:
		return message, nil
	}
}

func (rt *_router) authUserAsFounder(w http.ResponseWriter, groupId int64, userId int64) (bool, error) {
	isFounder, err := rt.db.IsFounder(groupId, userId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsFounder: %v", err)
		return false, err
	case isFounder:
		return true, nil
	default:
		rt.sendResponse(w, "Unauthorized to perform previliged actions on group.", http.StatusUnauthorized)
		return false, nil
	}
}
