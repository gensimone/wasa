package api

import (
	"net/http"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

// operationId: getAttachment
func (rt *_router) getAttachment(
	w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authorizeMessageAccess(w, ps, user)
	if err != nil {
		return
	}

	attachment, err := rt.db.GetAttachment(message.MessageId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetAttachment: %w", err)
		return
	}

	rt.sendResponse(w, attachment, http.StatusOK)
}
