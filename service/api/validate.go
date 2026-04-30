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

var validateNameRegex = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)

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

func (rt *_router) validateGroup(w http.ResponseWriter, id int64) (*database.Group, error) {
	g, err := rt.db.GetGroupById(id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, fmt.Sprintf("Group not found: %d", id), http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetGroupById: %w", err)
	}
	return g, err
}

func (rt *_router) validateFounder(w http.ResponseWriter, groupId int64, userId int64) (bool, error) {
	yes, err := rt.db.IsFounder(groupId, userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsFounder: %w", err)
	} else if !yes {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	return yes, err
}

func (rt *_router) validateUser(w http.ResponseWriter, id int64) (*database.User, error) {
	u, err := rt.db.GetUserById(id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, fmt.Sprintf("User not found: %d", id), http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserById: %w", err)
	}
	return u, err
}

func validateParameterInt64(w http.ResponseWriter, ps httprouter.Params, p string) (int64, error) {
	id, err := strconv.ParseInt(ps.ByName(p), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid `%s` type. Expected type: int64", p), http.StatusBadRequest)
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
		http.Error(w, "Invalid json: expected `photo:<binary>`", http.StatusBadRequest)
	}
	return s.Photo, err
}

func validateName(w http.ResponseWriter, r *http.Request) (string, error) {
	var s name
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Invalid json: expected `name:<string>`", http.StatusBadRequest)
	} else if !validateNameRegex.MatchString(s.Name) {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		err = fmt.Errorf("Invalid format")
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
