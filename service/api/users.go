package api

import (
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, user int64) {
	name, err := validateName(w, r)
	if err != nil {
		return
	}

	if _, err = rt.db.GetUserByName(name); errors.Is(err, sql.ErrNoRows) {
		if _, err := rt.db.SetMyUserName(user, name); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("SetMyUserName: %w", err)
		} else {
			sendResponse(w, nil, http.StatusNoContent)
		}
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserByName: %w", err)
	} else {
		sendResponse(w, "User name already taken", http.StatusConflict)
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, user int64) {
	photo, err := validatePhoto(w, r)
	if err != nil {
		return
	}

	if _, err := rt.db.SetMyPhoto(user, photo); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetMyPhoto: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) deleteUser(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, user int64) {
	if _, err := rt.db.DeleteUser(user); err != nil {
		sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteUser: %w", err)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) getUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, _ int64) {
	if users, err := rt.db.GetUserIds(); err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserIds: %w", err)
	} else {
		sendResponse(w, struct {
			Users []int64 `json:"users"`
		}{Users: users}, http.StatusOK)
	}
}

func (rt *_router) getUserById(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, _ int64) {
	user, err := validateParameterInt64(w, ps, "userId")
	if err != nil {
		return
	}

	if u, err := rt.validateUser(w, user); err == nil {
		sendResponse(w, u, http.StatusOK)
	}
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name, err := validateName(w, r)
	if err != nil {
		return
	}

	if id, err := rt.db.DoLogin(name); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DoLogin: %w", err)
	} else {
		sendResponse(w, struct{ Id int64 }{Id: *id}, http.StatusCreated)
	}
}
