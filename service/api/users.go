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

// operationId: getUsers
func (rt *_router) getUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, _ database.User) {
	userIds, err := rt.db.GetUserIds()
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserIds: %w", err)
		return
	}

	rt.sendResponse(w, struct {
		UserIds []int64 `json:"users"`
	}{UserIds: userIds}, http.StatusOK)
}

// operationId: getUserById
func (rt *_router) getUserById(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, _ database.User) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	user, err := rt.db.GetUserById(userId)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("User id %d not found", userId), http.StatusNotFound)
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserById: %w", userId, err)
	default:
		rt.sendResponse(w, user, http.StatusOK)
	}
}

// operationId: deleteUser
func (rt *_router) deleteUser(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	if userId != user.UserId {
		rt.sendResponse(
			w,
			"Parameter userId must be equal to the provided authentication id",
			http.StatusBadRequest,
		)
		return
	}
	_, err = rt.db.DeleteUser(user.UserId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteUser: %w", user.UserId, err)
		return
	}

	rt.sendResponse(w, nil, http.StatusNoContent)
}

// operationId: setMyUserName
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	if userId != user.UserId {
		rt.sendResponse(
			w,
			"Parameter userId must be equal to the provided authentication id",
			http.StatusBadRequest,
		)
		return
	}

	name, err := rt.checkName(w, r)
	if err != nil {
		return
	}

	_, err = rt.db.GetUserByName(*name)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		if _, err := rt.db.SetMyUserName(user.UserId, *name); err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("SetMyUserName: %w", user.UserId, err)
			return
		}

		rt.sendResponse(w, name, http.StatusOK)

	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserByName: %w", user.UserId, err)

	default:
		rt.sendResponse(w, "User name already taken", http.StatusConflict)
	}
}

// operationId: setMyPhoto
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	if userId != user.UserId {
		rt.sendResponse(
			w,
			"Parameter userId must be equal to the provided authentication id",
			http.StatusBadRequest,
		)
		return
	}

	photoUrl, err := rt.uploadFile(w, r, "photo")
	if err != nil {
		return
	}

	if _, err := rt.db.SetMyPhotoUrl(user.UserId, *photoUrl); err != nil {
		rt.baseLogger.Errorf("SetMyPhotoUrl: %w", user.UserId, err)

		if err = rt.removeFile(w, *photoUrl); err != nil { // Issue Internal Server Error
			return
		}

		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rt.sendResponse(w, photoUrl, http.StatusOK)
}

// operationId: doLogin
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name, err := rt.checkName(w, r)
	if err != nil {
		return
	}

	userId, err := rt.db.DoLogin(*name)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DoLogin: %w", err)
		return
	}

	rt.sendResponse(w, struct {
		UserId int64 `json:"userId"`
	}{UserId: *userId}, http.StatusCreated)
}
