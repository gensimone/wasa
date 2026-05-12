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
	messageId, err := rt.checkMessageId(w, r)
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
		rt.baseLogger.Errorf("GetMessage: %w", err)
		return
	}

	conversationId, err := rt._getConversation(w, r, ps, user)
	if err != nil {
		return
	}

	fmessage, err := rt.db.InsertMessage(
		user.UserId,
		*conversationId,
		message.Text,
		message.AttachmentId,
		true,
		nil,
	)

	rt.sendResponse(w, fmessage, http.StatusCreated)
}

// operationId: commentMessage
func (rt *_router) commentMessage(
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

	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		rt._insertMessage(w, r, user.UserId, message.ConversationId, &message.MessageId)
	default:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
	}
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

	if message.SenderId != user.UserId {
		rt.sendResponse(w, "Not athorized to delete message", http.StatusUnauthorized)
		return
	}

	_, err = rt.db.DeleteMessage(messageId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteMessage: %w", err)
		return
	}

	rt.sendResponse(w, nil, http.StatusNoContent)
}

// operationId: getMessage
func (rt *_router) getMessage(
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

	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %w", err)
		return
	case isMember && message.SenderId != user.UserId:
		rt._updateStatus(w, messageId, user.UserId)
	default:
		rt.sendResponse(w, "Not authorized to read message", http.StatusUnauthorized)
		return
	}

	rt.sendResponse(w, message, http.StatusOK)
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

	_, textOk := r.MultipartForm.Value["text"]
	files, fileOk := r.MultipartForm.File["file"]
	_, mediaTypeOk := r.MultipartForm.Value["mediaType"]

	if !(textOk && fileOk && mediaTypeOk) {
		errMsg := "Missing one or more required fields in multipart form"
		rt.sendResponse(w, errMsg, http.StatusBadRequest)
		return nil, errors.New(errMsg)
	}

	mediaType := r.FormValue("mediaType")

	switch {
		case len(files) == 0 && mediaType != "":
			errMsg := "Found media type but no file was provided"
			rt.sendResponse(w, errMsg, http.StatusBadRequest)
			return nil, errors.New(errMsg)
		case len(files) > 0 && mediaType == "":
			errMsg := "Found file but no media type was provided"
			rt.sendResponse(w, errMsg, http.StatusBadRequest)
			return nil, errors.New(errMsg)
		case len(files) > 0 && !database.IsValidMediaType(mediaType):
			errMsg := fmt.Sprintf("Invalid media type: %s", mediaType)
			rt.sendResponse(w, errMsg, http.StatusBadRequest)
			return nil, errors.New(errMsg)
	}

	url, err := rt.uploadFile(w, r, "file")
	if err != nil {
		return nil, err
	}

	attachmentId, err := rt.db.AddAttachment(*url, database.MediaType(mediaType))
	if err != nil {
		rt.baseLogger.Error("AddAttachment: %w", err)

		if err = rt.removeFile(w, *url); err != nil { // Issue Internal Server Error
			return nil, err
		}

		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	text := r.FormValue("text")
	message, err := rt.db.InsertMessage(
		senderId,
		conversationId,
		text,
		attachmentId,
		false,
		commentTo,
	)

	if err != nil {
		rt.baseLogger.Error("InsertMessage: %w", err)

		if err = rt.removeFile(w, *url); err != nil { // Issue Internal Server Error
			return nil, err
		}

		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	return message, nil
}

// Used internally by:
// - sendMessage
// - forwardMessage
// - commentMessage
func (rt *_router) _getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) (*int64, error) {
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
		rt.baseLogger.Errorf("IsUserById: %w", err)
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
		rt.baseLogger.Errorf("GetConversation: %w", err)
		return nil, err
	}

	if conversationId == nil {
		conversationId, err = rt.db.CreateConversation(user.UserId, userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("CreateConversation: %w", err)
			return nil, err
		}
	}

	return conversationId, err
}

