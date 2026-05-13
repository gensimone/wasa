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
	messageId, err := rt.checkMessageId(w, r)
	if err != nil {
		return
	}

	message, err := rt.db.GetMessage(*messageId)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %w", err)
		return
	case !isMember:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// reaction, err := rt.db.GetReaction(*messageId, user.UserId)
	// TODO: finish implement this
}

// TODO: implement this
// operationId: deleteReaction
func (rt *_router) deleteReaction(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
}

// TODO: implement this
// operationId: getReactions
func (rt *_router) getReactions(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
}
