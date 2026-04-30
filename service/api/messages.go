package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, sender int64,
) {
	content, err := validateMessage(w, r)
	if err != nil {
		return
	}

	rt.insertMessage(w, ps, sender, content.Text, content.Photo, false)
}

func (rt *_router) forwardMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, sender int64,
) {
	msgId, err := validateMsgId(w, r)
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msgId)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if yes, err := rt.db.IsMember(sender, m.Conversation); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		rt.insertMessage(w, ps, sender, m.Text, m.Photo, true)
	}
}

func (rt *_router) insertMessage(
	w http.ResponseWriter, ps httprouter.Params, sender int64, text string, photo *os.File, isForwarded bool,
) {
	receiver, err := validateParameterInt64(w, ps, "userId")
	if err != nil {
		return
	}

	conv, err := rt.db.GetConversation(sender, receiver)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetConversation: %w", err)
		return
	}

	var m *database.Message
	if conv == nil {
		m, err = rt.db.CreateConversation(sender, receiver, text, photo, isForwarded)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("CreateConversation: %w", err)
			return
		}
	} else {
		m, err = rt.db.InsertMessage(sender, *conv, text, photo, isForwarded, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("InsertMessage: %w", err)
			return
		}
	}

	sendResponse(w, m, http.StatusCreated)
}

func (rt *_router) commentMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, sender int64,
) {
	msg, err := validateParameterInt64(w, ps, "messageId")
	if err != nil {
		return
	}

	content, err := validateMessage(w, r)
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msg)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if yes, err := rt.db.IsMember(sender, m.Conversation); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		m, err = rt.db.InsertMessage(sender, m.Conversation, content.Text, content.Photo, false, &msg)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("InsertMessage: %w", err)
		} else {
			sendResponse(w, m, http.StatusOK)
		}
	}
}

func (rt *_router) deleteMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user int64,
) {
	msg, err := validateParameterInt64(w, ps, "messageId")
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msg)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if m.Sender != user {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if _, err = rt.db.DeleteMessage(msg); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteMessage: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) getMessage(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user int64,
) {
	msg, err := validateParameterInt64(w, ps, "messageId")
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msg)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if yes, err := rt.db.IsMember(user, m.Conversation); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		if m.Sender != user {
			s, err := rt.db.GetStatusOf(m.Id, user)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("GetStatusOf: %w", err)
				return
			} else if s.Info == "not received" {
				rt.db.UpdateStatus(msg, user, "received")
			}
		}

		sendResponse(w, m, http.StatusOK)
	}
}

func (rt *_router) getStatus(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user int64,
) {
	msg, err := validateParameterInt64(w, ps, "messageId")
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msg)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if yes, err := rt.db.IsMember(user, m.Conversation); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if user != m.Sender {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else {
		status, err := rt.db.GetStatus(msg)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetStatus: %w", err)
		} else {
			sendResponse(w, struct {
				Status []database.Status `json:"status"`
			}{Status: status}, http.StatusOK)
		}
	}
}

func (rt *_router) addReaction(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, sender int64,
) {
	msg, err := validateParameterInt64(w, ps, "messageId")
	if err != nil {
		return
	}

	emoji, err := validateEmoji(w, r)
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msg)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if yes, err := rt.db.IsMember(sender, m.Conversation); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if r, err := rt.db.GetReaction(msg, sender); errors.Is(err, sql.ErrNoRows) {
		if _, err = rt.db.AddReaction(emoji, msg, sender); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("AddReaction: %w", err)
		} else if r, err = rt.db.GetReaction(msg, sender); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetReaction: %w", err)
		} else {
			sendResponse(w, *r, http.StatusCreated)
		}
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetReaction: %w", err)
	} else if r.Emoji == emoji {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else if _, err = rt.db.UpdateReaction(emoji, msg, sender); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("UpdateReaction: %w", err)
	} else if r, err = rt.db.GetReaction(msg, sender); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetReaction: %w", err)
	} else {
		sendResponse(w, *r, http.StatusCreated)
	}
}

func (rt *_router) deleteReaction(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, sender int64,
) {
	msg, err := validateParameterInt64(w, ps, "messageId")
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msg)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if yes, err := rt.db.IsMember(sender, m.Conversation); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if _, err := rt.db.GetReaction(msg, sender); errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteReaction: %w", err)
	} else if _, err := rt.db.DeleteReaction(msg, sender); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteReaction: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) getReactions(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user int64,
) {
	msg, err := validateParameterInt64(w, ps, "messageId")
	if err != nil {
		return
	}

	m, err := rt.db.GetMessage(msg)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
	} else if yes, err := rt.db.IsMember(user, m.Conversation); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if reactions, err := rt.db.GetReactions(msg); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetReactions: %w", err)
	} else {
		sendResponse(w, struct {
			Reactions []database.Reaction
		}{Reactions: reactions}, http.StatusOK)
	}
}
