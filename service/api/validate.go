package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type userId struct {
	UserId int64 `json:"userId"`
}
type msgId struct {
	MsgId int64 `json:"messageId"`
}
type name struct {
	Name string `json:"name"`
}
type emoji struct {
	Emoji string `json:"emoji"`
}
type photo struct {
	Photo os.File `json:"photo"`
}
type Content struct {
	Text  string   `json:"text"`
	Photo *os.File `json:"photo"`
}
type GroupInfo struct {
	Name  string   `json:"name"`
	Photo *os.File `json:"photo"`
}

// FIXME: Allow spaces.
var validateNameRegex = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)

func (rt *_router) validateUserId(w http.ResponseWriter, r *http.Request) (int64, error) {
	var s userId
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if _, err = rt.db.GetUserById(s.UserId); errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserById: %w", err)
	}

	return s.UserId, err
}

func validateGroupInfo(w http.ResponseWriter, r *http.Request) (GroupInfo, error) {
	var s GroupInfo
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if !validateNameRegex.MatchString(s.Name) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		err = fmt.Errorf("Bad Request")
	}

	return s, err
}

func validateEmoji(w http.ResponseWriter, r *http.Request) (string, error) {
	var s emoji
	err := json.NewDecoder(r.Body).Decode(&s)
	// TODO: check emoji validity
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return s.Emoji, err
}

func validateMsgId(w http.ResponseWriter, r *http.Request) (int64, error) {
	var s msgId
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return s.MsgId, err
}

func (rt *_router) validateGroup(w http.ResponseWriter, group int64) (*database.Group, error) {
	g, err := rt.db.GetGroupById(group)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetGroupById: %w", err)
	}
	return g, err
}

func (rt *_router) validateFounder(w http.ResponseWriter, group int64, user int64) (bool, error) {
	yes, err := rt.db.IsFounder(group, user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsFounder: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	return yes, err
}

func (rt *_router) validateUser(w http.ResponseWriter, user int64) (*database.User, error) {
	u, err := rt.db.GetUserById(user)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserById: %w", err)
	}
	return u, err
}

func validateParameterInt64(w http.ResponseWriter, ps httprouter.Params, param string) (int64, error) {
	id, err := strconv.ParseInt(ps.ByName(param), 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	return id, err
}

func validateMessage(w http.ResponseWriter, r *http.Request) (Content, error) {
	var s Content
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return s, err
}

func validatePhoto(w http.ResponseWriter, r *http.Request) (os.File, error) {
	var s photo
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	return s.Photo, err
}

func validateName(w http.ResponseWriter, r *http.Request) (string, error) {
	var s name
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else if !validateNameRegex.MatchString(s.Name) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		err = fmt.Errorf("Bad Request")
	}
	return s.Name, err
}

func (rt *_router) validateAuthorization(
	fn func(http.ResponseWriter, *http.Request, httprouter.Params, int64),
) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		id, err := strconv.ParseInt(auth, 10, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if _, err := rt.db.GetUserById(id); errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetUserById: %w", err)
		} else {
			fn(w, r, ps, id)
		}
	}
}
