package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
	convs, err := rt.db.GetConversations(user)
	if err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else {
		sendResponse(w, struct {
			Convs []int64 `json:"conversations"`
		}{Convs: convs}, http.StatusOK)
	}
}

func (rt *_router) getConversation(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
	conv, err := validateParameterInt64(w, ps, "conversationId")
	if err != nil {
		return
	}

	if yes, err := rt.db.IsMember(user, conv); err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if convs, err := rt.db.GetMessages(conv); err != nil || convs == nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessages: %w", err)
	} else {
		sendResponse(w, struct {
			Messages []int64 `json:"messages"`
		}{Messages: convs}, http.StatusOK)
	}
}

func (rt *_router) getMembers(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
	conv, err := validateParameterInt64(w, ps, "conversationId")
	if err != nil {
		return
	}

	if yes, err := rt.db.IsMember(user, conv); err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if members, err := rt.db.GetMembers(conv); err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else {
		sendResponse(w, struct {
			Members []int64 `json:"members"`
		}{Members: members}, http.StatusOK)
	}
}
