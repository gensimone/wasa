package api

import (
	"net/http"

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
	if err != nil || messageIds == nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessageIds: %v", err)
		return
	}

	rt.sendResponse(w, struct {
		MessageIds []int64 `json:"messageIds"`
	}{MessageIds: messageIds}, http.StatusOK)
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
