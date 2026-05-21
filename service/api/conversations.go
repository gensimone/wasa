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
		ConversationId int64 `json:"conversationId"`
		IsGroup        bool  `json:"isGroup"`
	}

	var conversations []data
	for _, conversationId := range conversationIds {
		isGroup, err := rt.db.IsGroup(conversationId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("IsGroup: %v", err)
			return
		}

		conversations = append(conversations, data{
			ConversationId: conversationId,
			IsGroup:        isGroup,
		})
	}

	rt.sendResponse(w, struct {
		Conversations []data `json:"conversations"`
	}{Conversations: conversations}, http.StatusOK)
}

// operationId: getConversation
func (rt *_router) getConversation(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := rt.authConversationAccess(w, ps, user)
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
	conversationId, err := rt.authConversationAccess(w, ps, user)
	if err != nil {
		return
	}

	messageIds, err := rt.db.GetMessageIds(*conversationId)
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessageIds: %v", err)
		return
	}

	var message *database.Message
	if len(messageIds) == 0 {
		message = nil
	} else {
		message, err = rt.db.GetMessage(messageIds[len(messageIds)-1])
		if err != nil {
			rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetMessage: %v", err)
			return
		}
	}

	rt.sendResponse(w, message, http.StatusOK)
}

// operationId: getMembers
func (rt *_router) getMembers(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := rt.authConversationAccess(w, ps, user)
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
