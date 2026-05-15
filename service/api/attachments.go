package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

// operationId: getAttachment
func (rt *_router) getAttachment(
	w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccess(w, ps, user)
	if err != nil {
		return
	}

	attachment, err := rt.db.GetAttachment(message.MessageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, "Message has no attachment", http.StatusBadRequest)
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetAttachment: %w", err)
	default:
		rt.sendResponse(w, attachment, http.StatusOK)
	}
}
