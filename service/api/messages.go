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

// operationId: sendMessage
func (rt *_router) sendMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	conversationId, err := rt._getConversation(w, r, ps, user)
	if err != nil {
		return
	}

	message, err := rt._insertMessage(w, r, user.UserId, *conversationId, nil)
	if err != nil {
		return
	}

	rt.sendResponse(w, message, http.StatusCreated)
}

// operationId: forwardMessage
func (rt *_router) forwardMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	messageId, err := rt.getMessageIdFromReq(w, r)
	if err != nil {
		return
	}

	message, err := rt.db.GetMessage(*messageId)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %v", err)
		return
	}

	conversationId, err := rt._getConversation(w, r, ps, user)
	if err != nil {
		return
	}

	fmessage, err := rt.db.InsertMessage(
		message.Text,
		user.UserId,
		*conversationId,
		true,
		nil,
		message.AttachmentUrl,
		message.MediaType,
	)

	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("InsertMessage: %v", err)
		return
	}

	rt.sendResponse(w, fmessage, http.StatusCreated)
}

// operationId: commentMessage
func (rt *_router) commentMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
	if err != nil {
		return
	}

	comment, err := rt._insertMessage(w, r, user.UserId, message.ConversationId, &message.MessageId)
	if err != nil {
		return
	}

	rt.sendResponse(w, comment, http.StatusCreated)
}

// operationId: uncommentMessage
func (rt *_router) uncommentMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	rt.deleteMessage(w, r, ps, user)
}

// operationId: deleteMessage
func (rt *_router) deleteMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
	if err != nil {
		return
	}

	// NOTE: We must be the sender of the message in order to delete the message.
	if message.SenderId != user.UserId {
		rt.sendResponse(w, "Not athorized to delete message", http.StatusUnauthorized)
		return
	}

	_, err = rt.db.DeleteMessage(message.MessageId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteMessage: %v", err)
		return
	}

	rt.sendResponse(w, nil, http.StatusNoContent)
}

// operationId: getMessage
func (rt *_router) getMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	message, err := rt.authMessageAccessParam(w, ps, user)
	if err != nil {
		return
	}

	// NOTE: We must be the receiver of the message in order to update its receipt.
	if message.SenderId == user.UserId {
		rt.sendResponse(w, message, http.StatusOK)
		return
	}

	receipt, err := rt.db.GetReceipt(message.MessageId, user.UserId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, message, http.StatusOK)
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetReceipt: %v", err)
	case receipt.Status == database.Received || receipt.Status == database.Read:
		rt.sendResponse(w, message, http.StatusOK)
	case receipt.Status == database.Sent:
		_, err = rt.db.SetReceiptStatus(message.MessageId, user.UserId, database.Received)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("SetReceiptStatus: %v", err)
			return
		}

		rt.sendResponse(w, message, http.StatusOK)
	}
}

// Used internally by:
// - sendMessage
// - sendMessageToGroup
// - commentMessage
func (rt *_router) _insertMessage(
	w http.ResponseWriter,
	r *http.Request,
	senderId int64,
	conversationId int64,
	commentTo *int64,
) (*database.Message, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		rt.sendResponse(w, "Invalid multipart form", http.StatusBadRequest)
		return nil, err
	}

	missingText := false
	text := r.FormValue("text")
	if text == "" {
		missingText = true
	}

	missingFile := len(r.MultipartForm.File["file"]) == 0

	if missingText && missingFile {
		errMsg := "Empty message content"
		rt.sendResponse(w, errMsg, http.StatusBadRequest)
		return nil, errors.New(errMsg)
	}

	attachmentUrl := ""
	var mediaType database.MediaType = ""

	if !missingFile {
		mediaTypeStr := r.FormValue("mediaType")
		if mediaTypeStr == "" {
			errMsg := "The file was provided but mediaType was not"
			rt.sendResponse(w, errMsg, http.StatusBadRequest)
			return nil, errors.New(errMsg)
		}

		mediaType = database.MediaType(mediaTypeStr)

		if !database.IsValidMediaType(mediaType) {
			errMsg := "Invalid media type"
			rt.sendResponse(w, errMsg, http.StatusBadRequest)
			return nil, errors.New(errMsg)
		}

		attachmentUrl, err = rt.uploadMediaFile(w, r, "file")
		if err != nil {
			return nil, err
		}
	}

	message, err := rt.db.InsertMessage(
		text,
		senderId,
		conversationId,
		false,
		commentTo,
		attachmentUrl,
		mediaType,
	)

	if err != nil {
		rt.baseLogger.Errorf("InsertMessage: %v", err)
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)

		if attachmentUrl != "" {
			_ = rt.removeMediaFile(attachmentUrl)
		}

		return nil, err
	}

	return message, nil
}

// Used internally by:
// - sendMessage
// - forwardMessage
// - commentMessage
func (rt *_router) _getConversation(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) (*int64, error) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return nil, err
	}

	if userId == user.UserId {
		errMsg := "Sender and receiver are the same"
		rt.sendResponse(w, errMsg, http.StatusBadRequest)
		return nil, errors.New(errMsg)
	}

	isValid, err := rt.db.IsUserById(userId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsUserById: %v", err)
		return nil, err
	}

	if !isValid {
		errMsg := fmt.Sprintf("User %d not found", userId)
		rt.sendResponse(w, errMsg, http.StatusNotFound)
		return nil, errors.New(errMsg)
	}

	conversationId, err := rt.db.GetConversation(user.UserId, userId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetConversation: %v", err)
		return nil, err
	}

	if conversationId == nil {
		conversationId, err = rt.db.CreateConversation(user.UserId, userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("CreateConversation: %v", err)
			return nil, err
		}
	}

	return conversationId, err
}
