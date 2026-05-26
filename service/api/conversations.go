package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

// operationId: getMyConversations
func (rt *_router) getMyConversations(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationIds, err := rt.db.GetConversations(user.UserId)
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetConversations: %v", err)
		return
	}

	type data struct {
		Id      int64 `json:"id"`
		IsGroup bool  `json:"isGroup"`
	}

	var conversations []data
	for _, conversationId := range conversationIds {
		isGroup, err := rt.db.IsGroup(conversationId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("IsGroup: %v", err)
			return
		}

		var entry data
		if !isGroup {
			otherMembers, err := rt.db.GetOtherMembers(conversationId, user.UserId)
			switch {
			case err != nil:
				rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("GetOtherMembers: %v", err)
				return
			case len(otherMembers) > 1:
				rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("Invalid state: non-group conversations must have exactly two members. Found: %d", len(otherMembers))
				return
			default:
				entry = data{
					Id:      otherMembers[0],
					IsGroup: false,
				}
			}
		} else {
			entry = data{
				Id:      conversationId,
				IsGroup: true,
			}
		}

		conversations = append(conversations, entry)
	}

	rt.sendResponse(w, struct {
		Conversations []data `json:"conversations"`
	}{Conversations: conversations}, http.StatusOK)
}

// operationId: getConversation
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	direct, err := rt.validateDirect(w, r)
	if err != nil {
		return
	}

	var conversationId *int64
	if direct {
		conversationId, _, err = rt._getConversation(w, r, ps, user)
	} else {
		conversationId, err = rt.authConversationAccessParam(w, ps, user, "id")
	}

	if err != nil {
		return
	}

	if conversationId == nil {
		rt.sendResponse(w, struct {
			MessageIds []int64 `json:"messageIds"`
		}{MessageIds: []int64{}}, http.StatusOK)
	}

	messageIds, err := rt.db.GetMessageIds(*conversationId)
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessageIds: %v", err)
		return
	}

	rt.sendResponse(w, struct {
		MessageIds []int64 `json:"messageIds"`
	}{MessageIds: messageIds}, http.StatusOK)
}

// operationId: getLastMessage
// func (rt *_router) getLastMessage(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
// 	conversationId, err := rt.authConversationAccessParam(w, ps, user)
// 	if err != nil {
// 		return
// 	}
//
// 	messageIds, err := rt.db.GetMessageIds(*conversationId)
// 	if err != nil {
// 		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
// 		rt.baseLogger.Errorf("GetMessageIds: %v", err)
// 		return
// 	}
//
// 	if len(messageIds) == 0 {
// 		rt.sendResponse(w, nil, http.StatusOK)
// 		return
// 	}
//
// 	message, err := rt.db.GetMessage(messageIds[len(messageIds)-1])
// 	if err != nil {
// 		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
// 		rt.baseLogger.Errorf("GetMessage: %v", err)
// 		return
// 	}
//
// 	// NOTE: We must be the receiver of the message in order to update its receipt.
// 	if message.SenderId == user.UserId {
// 		rt.sendResponse(w, message, http.StatusOK)
// 		return
// 	}
//
// 	receipt, err := rt.db.GetReceipt(message.MessageId, user.UserId)
// 	switch {
// 	case errors.Is(err, sql.ErrNoRows):
// 		rt.sendResponse(w, message, http.StatusOK)
// 	case err != nil:
// 		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
// 		rt.baseLogger.Errorf("GetReceipt: %v", err)
// 	case receipt.Status == database.Received || receipt.Status == database.Read:
// 		rt.sendResponse(w, message, http.StatusOK)
// 	case receipt.Status == database.Sent:
// 		_, err = rt.db.SetReceiptStatus(message.MessageId, user.UserId, database.Received)
// 		if err != nil {
// 			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
// 			rt.baseLogger.Errorf("SetReceiptStatus: %v", err)
// 			return
// 		}
// 		rt.sendResponse(w, message, http.StatusOK)
// 	}
// }

// operationId: sendMessage
func (rt *_router) sendMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
	direct, err := rt.validateDirect(w, r)
	if err != nil {
		return
	}

	var conversationId *int64
	if direct {
		conversationId2, userId, err := rt._getConversation(w, r, ps, user)
		if err != nil {
			return
		}

		if conversationId2 == nil {
			conversationId2, err = rt.db.CreateConversation(user.UserId, *userId)
			if err != nil {
				rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("CreateConversation: %v", err)
				return
			}
		}
		conversationId = conversationId2
	} else {
		conversationId, err = rt.authConversationAccessParam(w, ps, user, "id")
		if err != nil {
			return
		}
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
	direct, err := rt.validateDirect(w, r)
	if err != nil {
		return
	}

	var conversationId *int64
	if direct {
		conversationId2, userId, err := rt._getConversation(w, r, ps, user)
		if err != nil {
			return
		}

		if conversationId2 == nil {
			conversationId2, err = rt.db.CreateConversation(user.UserId, *userId)
			if err != nil {
				rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("CreateConversation: %v", err)
				return
			}
		}

		conversationId = conversationId2
	} else {
		conversationId, err = rt.authConversationAccessParam(w, ps, user, "id")
		if err != nil {
			return
		}
	}

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
	text := strings.TrimSpace(r.FormValue("text"))
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

func (rt *_router) _getConversation(
	w http.ResponseWriter,
	_ *http.Request,
	ps httprouter.Params,
	user database.User,
) (*int64, *int64, error) {
	userId, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter id must be an int64", http.StatusBadRequest)
		return nil, nil, err
	}

	if userId == user.UserId {
		errMsg := "Sender and receiver are the same"
		rt.sendResponse(w, errMsg, http.StatusBadRequest)
		return nil, &userId, errors.New(errMsg)
	}

	isValid, err := rt.db.IsUserById(userId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsUserById: %v", err)
		return nil, &userId, err
	}

	if !isValid {
		errMsg := fmt.Sprintf("User %d not found", userId)
		rt.sendResponse(w, errMsg, http.StatusNotFound)
		return nil, &userId, errors.New(errMsg)
	}

	conversationId, err := rt.db.GetConversation(user.UserId, userId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetConversation: %v", err)
		return nil, &userId, err
	}

	return conversationId, &userId, err
}

// Helper function that extracts the "direct" query from the request.
func (rt *_router) validateDirect(w http.ResponseWriter, r *http.Request) (bool, error) {
	raw := r.URL.Query().Get("direct")
	if raw == "" {
		errMsg := "Missing required query: direct"
		rt.sendResponse(w, errMsg, http.StatusBadRequest)
		return false, errors.New(errMsg)
	}

	b, err := strconv.ParseBool(raw)
	if err != nil {
		errMsg := "Invalid value for query: direct"
		rt.sendResponse(w, errMsg, http.StatusBadRequest)
		return false, errors.New(errMsg)
	}

	return b, nil
}
