package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, id int64) {
	name, err := validateName(w, r)
	if err != nil { return }

	if _, err = rt.db.GetUserByName(name); errors.Is(err, sql.ErrNoRows) {
		if err := rt.db.SetMyUserName(id, name); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			sendResponse(w, nil, http.StatusNoContent)
		}
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		sendResponse(w, "The username provided is already in use", http.StatusConflict)
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, id int64) {
	photo, err := validatePhoto(w, r)
	if err != nil { return }

	if err := rt.db.SetMyPhoto(id, photo); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) deleteUser( w http.ResponseWriter, _ *http.Request, _ httprouter.Params, id int64) {
	if err := rt.db.DeleteUser(id); err != nil {
		sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func (rt *_router) getUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, _ int64) {
	if ids, err := rt.db.GetUserIds(); err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else {
		sendResponse(w, struct{Users []int64 `json:"users"`}{Users: ids}, http.StatusOK)
	}
}

func (rt *_router) getUserById(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, _ int64) {
	id, err := validateParameterInt64(w, ps, "userId")
	if err != nil { return }

	if u, err := rt.validateUser(w, id); err == nil {
		sendResponse(w, u, http.StatusOK)
	}
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name, err := validateName(w, r)
	if err != nil { return }

	if id, err := rt.db.DoLogin(name); err != nil
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		sendResponse(w, struct{Id int64}{Id: id}, http.StatusCreated)
	}
}
