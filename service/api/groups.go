package api

import (
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

	if yes, err := rt.db.IsMember(group, user); err != nil {
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

	if _, err = rt.db.DeleteGroup(group); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteGroup: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) leaveGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
	group, err := validateParameterInt64(w, ps, "groupId")
	if err != nil {
		return
	}

	if _, err := rt.validateGroup(w, group); err != nil {
		return
	}

	if yes, err := rt.db.IsMember(user, group); err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	} else if !yes {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else {
		// TODO: Check rows affected?
		if _, err := rt.db.DeleteUserConversation(group, user); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("DeleteUserConversation: %w", err)
		}
	}
}

func (rt *_router) addToGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
}

func (rt *_router) createGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
}

func (rt *_router) forwardMessageToGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
}

func (rt *_router) sendMessageToGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
}

func (rt *_router) removeUser(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user int64) {
}
