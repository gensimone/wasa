package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	g, err := rt.validateGroup(w, group)
	if err != nil {
		return
	}

	if yes, err := rt.db.IsMember(user, group); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		sendResponse(w, g, http.StatusOK)
	}
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user int64) {
	name, err := validateName(w, r)
	if err != nil {
		return
	}

	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	if _, err = rt.validateGroup(w, group); err != nil {
		return
	}

	if yes, err := rt.validateFounder(w, group, user); !yes || err != nil {
		return
	}

	if _, err = rt.db.SetGroupName(group, name); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetGroupName: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user int64) {
	photo, err := validatePhoto(w, r)
	if err != nil {
		return
	}

	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	if _, err := rt.validateGroup(w, group); err != nil {
		return
	}

	if yes, err := rt.validateFounder(w, group, user); !yes || err != nil {
		return
	}

	if _, err = rt.db.SetGroupPhoto(group, photo); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetGroupPhoto: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) deleteGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	if _, err := rt.validateGroup(w, group); err != nil {
		return
	}

	if yes, err := rt.validateFounder(w, group, user); !yes || err != nil {
		return
	}

	if _, err = rt.db.DeleteConversation(group); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteConversation: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) leaveGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	g, err := rt.validateGroup(w, group)
	if err != nil {
		return
	}

	if yes, err := rt.db.IsMember(user, group); err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else if g.Founder == user {
		if _, err := rt.db.DeleteConversation(group); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("DeleteConversation: %w", err)
		} else {
			sendResponse(w, nil, http.StatusNoContent)
		}
	} else if _, err := rt.db.DeleteUserConversation(group, user); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteUserConversation: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, founder int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	user, err := rt.validateUserId(w, r)
	if err != nil {
		return
	}

	g, err := rt.validateGroup(w, group)
	if err != nil {
		return
	}

	if g.Founder != founder {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if user == founder {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else if yes, err := rt.db.IsMember(user, group); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if yes {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else if _, err = rt.db.AddConversation(group, user); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("AddConversation: %w", err)
	} else {
		sendResponse(w, struct {
			User int64 `json:"userId"`
		}{User: user}, http.StatusCreated)
	}
}

func (rt *_router) removeUser(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, founder int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	user, err := validateParameterInt64(w, ps, "userId")
	if err != nil {
		return
	}

	g, err := rt.validateGroup(w, group)
	if err != nil {
		return
	}

	if g.Founder != founder {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if _, err = rt.db.DeleteUserConversation(group, user); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteUserConversation: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, founder int64) {
	i, err := validateGroupInfo(w, r)
	if err != nil {
		return
	}

	if yes, err := rt.db.GroupExists(founder, i.Name); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GroupExists: %w", err)
	} else if yes {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else if g, err := rt.db.CreateGroup(founder, i.Name, i.Photo); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("CreateGroup: %w", err)
	} else {
		sendResponse(w, g, http.StatusCreated)
	}
}

func (rt *_router) forwardMessageToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sender int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

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
	} else if m, err := rt.db.InsertMessage(sender, group, m.Text, m.Photo, true, nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("InsertMessage: %w", err)
	} else {
		sendResponse(w, m, http.StatusCreated)
	}
}

func (rt *_router) sendMessageToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sender int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	content, err := validateMessage(w, r)
	if err != nil {
		return
	}

	m, err := rt.db.InsertMessage(sender, group, content.Text, content.Photo, false, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("InsertMessage: %w", err)
	} else {
		sendResponse(w, m, http.StatusCreated)
	}
}
