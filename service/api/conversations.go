package api

import (
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
		rt.baseLogger.Errorf("GetConversations: %w", err)
		return
	}

	type data struct {
		ConversationId int64 `json:"conversationId"`
		IsGroup 	   bool  `json:"isGroup"`
	}

	var conversations []data
	for _, conversationId := range conversationIds {
		isGroup, err := rt.db.IsGroup(conversationId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("IsGroup: %w", err)
			return
		}

		conversations = append(conversations, data{
			ConversationId: conversationId,
			IsGroup: isGroup,
		})
	}

	rt.sendResponse(w, struct {
		Conversations []data `json:"conversations"`
	}{Conversations: conversations}, http.StatusOK)
}

// operationId: getConversation
func (rt *_router) getConversation(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter conversationId must be an int64", http.StatusBadRequest)
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, conversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		messageIds, err := rt.db.GetMessageIds(conversationId)
		if err != nil || messageIds == nil {
			rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetMessageIds: %w", err)
			return
		}

		rt.sendResponse(w, struct {
			MessageIds []int64 `json:"messageIds"`
		}{MessageIds: messageIds}, http.StatusOK)
	default:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// operationId: getMembers
func (rt *_router) getMembers(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter conversationId must be an int64", http.StatusBadRequest)
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, conversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		userIds, err := rt.db.GetMembers(conversationId)
		if err != nil {
			rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetMembers: %w", err)
			return
		}

		rt.sendResponse(w, struct {
			UserIds []int64 `json:"userIds"`
		}{UserIds: userIds}, http.StatusOK)
	default:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
	}
}
