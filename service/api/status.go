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

// operationId: updateStatus
func (rt *_router) updateStatus(
	w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User,
) {
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter messageId must be an int64", http.StatusBadRequest)
		return
	}

	message, err := rt.db.GetMessage(messageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetMessage: %w", err)
		return
	}

	// NOTE: We must be the receiver of the message in order to update its status.
	if message.SenderId == user.UserId {
		rt.sendResponse(w, "Bad Request", http.StatusBadRequest)
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		_, err = rt.db.UpdateStatus(messageId, user.UserId, database.Read)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("UpdateStatus: %w", err)
			return
		}
		rt.sendResponse(w, nil, http.StatusNoContent)
	default:
		rt.sendResponse(w, "Not authorized to read message", http.StatusUnauthorized)
	}
}

// operationId: getStatus
func (rt *_router) getStatus(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter messageId must be an int64", http.StatusBadRequest)
		return
	}

	message, err := rt.db.GetMessage(messageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("GetMessage: %w", err)
		return
	}

	// NOTE:
	// We can only access the status of our messages.
	if message.SenderId != user.UserId {
		rt.sendResponse(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// NOTE:
	// Even though we are the authors of the message, we must be inside the message conversation
	// to be able to read the message status. For example, if we have been kicked out of a group
	// we should not be able to read the status of that group's messages
	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		status, err := rt.db.GetStatus(messageId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetStatus: %v", err)
			return
		}

		rt.sendResponse(w, struct {
			Status []database.Status `json:"status"`
		}{Status: status}, http.StatusOK)
	default:
		rt.sendResponse(w, "Not authorized to get message status", http.StatusUnauthorized)
	}
}

// Used internally by:
// - api/getMessage
func (rt *_router) _updateStatus(w http.ResponseWriter, messageId int64, userId int64) error {
	status, err := rt.db.GetStatusOf(messageId, userId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetStatusOf: %w", err)
		return err
	}

	var info database.Info

	switch status.Info {
	case database.Read:
		return nil
	case database.NotReceived:
		info = database.Received
	case database.Received:
		info = database.Read
	}

	_, err = rt.db.UpdateStatus(messageId, userId, info)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("UpdateStatus: %w", err)
		return err
	}

	return nil
}
