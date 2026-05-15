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

	if message.AttachmentId == nil {
		rt.sendResponse(w, "Message has no attachment", http.StatusBadRequest)
		return
	}

	attachment, err := rt.db.GetAttachment(*message.AttachmentId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetAttachment: %w", err)

		// FIXME: This should not be here. I need to fix the database to prevent this.
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.Errorf(
				"Inconsistent state detected: message %d has attachment with id %d "+
					"but no attachment with that id was found",
				message.MessageId,
				*message.AttachmentId,
			)
		}

		return
	}

	rt.sendResponse(w, attachment, http.StatusOK)
}
