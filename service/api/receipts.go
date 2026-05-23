package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

// operationId: setMessageStatusAsRead
func (rt *_router) setMessageStatusAsRead(
	w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
	if err != nil {
		return
	}

	// NOTE: We must be the receiver of the message in order to update its receipt.
	if message.SenderId == user.UserId {
		rt.sendResponse(w, "Unauthorized to update message status", http.StatusUnauthorized)
		return
	}

	receipt, err := rt.db.GetReceipt(message.MessageId, user.UserId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// No receipt for this user. It happens when a user read a message in a group
		// that was sent before the user joined the group.
		rt.sendResponse(w, nil, http.StatusNoContent)
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetReceipt: %v", err)
	case receipt.Status == database.Sent:
		rt.sendResponse(w, "You must access the message first", http.StatusUnauthorized)
	case receipt.Status == database.Read:
		rt.sendResponse(w, nil, http.StatusNoContent)
	default:
		_, err = rt.db.SetReceiptStatus(message.MessageId, user.UserId, database.Read)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("SetReceiptStatus: %v", err)
			return
		}

		rt.sendResponse(w, nil, http.StatusNoContent)
	}
}

// operationId: getReceipts
func (rt *_router) getReceipts(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
	if err != nil {
		return
	}

	// NOTE: We can only access the receipts of our messages.
	if message.SenderId != user.UserId {
		rt.sendResponse(w, "Unauthorized to read message status", http.StatusUnauthorized)
		return
	}

	receipts, err := rt.db.GetReceipts(message.MessageId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetReceipts: %v", err)
		return
	}

	rt.sendResponse(w, struct {
		Receipts []database.Receipt `json:"receipts"`
	}{Receipts: receipts}, http.StatusOK)
}
