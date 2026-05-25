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
func (rt *_router) getConversation(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := rt.authConversationAccessParam(w, ps, user)
	if err != nil {
		return
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

func (rt *_router) getConversationByUserId(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	if userId == user.UserId {
		rt.sendResponse(
			w,
			"Parameter userId must be different to the provided authentication id",
			http.StatusBadRequest,
		)
		return
	}

	_, err = rt.db.GetUserById(userId)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("User id %d not found", userId), http.StatusNotFound)
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserById: %v", err)
	default:
		conversationId, err := rt.db.GetConversation(user.UserId, userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetConversation: %v", err)
			return
		}

		type MessageIds struct {
			MessageIds []int64 `json:"messageIds"`
		}

		if conversationId == nil {
			rt.sendResponse(w, MessageIds{MessageIds: []int64{}}, http.StatusOK)
			return
		}

		messageIds, err := rt.db.GetMessageIds(*conversationId)
		if err != nil {
			rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetMessageIds: %v", err)
			return
		}

		rt.sendResponse(w, MessageIds{MessageIds: messageIds}, http.StatusOK)
	}
}

// operationId: getLastMessage
func (rt *_router) getLastMessage(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := rt.authConversationAccessParam(w, ps, user)
	if err != nil {
		return
	}

	messageIds, err := rt.db.GetMessageIds(*conversationId)
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessageIds: %v", err)
		return
	}

	if len(messageIds) == 0 {
		rt.sendResponse(w, nil, http.StatusOK)
		return
	}

	message, err := rt.db.GetMessage(messageIds[len(messageIds)-1])
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %v", err)
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

// operationId: getMemberIds
func (rt *_router) getMemberIds(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := rt.authConversationAccessParam(w, ps, user)
	if err != nil {
		return
	}

	userIds, err := rt.db.GetMembers(*conversationId)
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMembers: %v", err)
		return
	}

	rt.sendResponse(w, struct {
		UserIds []int64 `json:"userIds"`
	}{UserIds: userIds}, http.StatusOK)
}

// operationId: sendMessageToConversation
func (rt *_router) sendMessageToConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := rt.authConversationAccessParam(w, ps, user)
	if err != nil {
		return
	}

	message, err := rt._insertMessage(w, r, user.UserId, *conversationId, nil)
	if err != nil {
		return
	}

	rt.sendResponse(w, message, http.StatusCreated)
}

// operationId: forwardMessageToConversation
func (rt *_router) forwardMessageToConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := rt.authConversationAccessParam(w, ps, user)
	if err != nil {
		return
	}

	messageId, err := rt.getMessageIdFromReq(w, r)
	if err != nil {
		return
	}

	message, err := rt.db.GetMessage(*messageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %v", err)
	default:
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
}
