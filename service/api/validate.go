package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

var nameRegex = regexp.MustCompile(`^.*$`)

func (rt *_router) authorize(
	fn func(http.ResponseWriter, *http.Request, httprouter.Params, database.User),
) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := strconv.ParseInt(auth, 10, 64)
		if err != nil {
			rt.sendResponse(w, "User id must be an int64", http.StatusBadRequest)
			return
		}

		user, err := rt.db.GetUserById(userId)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			rt.sendResponse(w, "Not Found", http.StatusNotFound)
		case err != nil:
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("GetUserById (user_id: %d): %w", userId, err)
		default:
			fn(w, r, ps, *user)
		}
	}
}

// To process the operation IDs for addReaction, removeReaction, and getReactions,
// we first need to verify that the authenticated user has the privileges to
// perform these operations. This function accomplishes exactly that.
//
// Actions performed by this function:
// - Check message validity.
// - Check user privileges for the message.
func (rt *_router) authorizeMessageAccess(
	w http.ResponseWriter, ps httprouter.Params, user database.User,
) (*database.Message, error) {
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter messageId must be an int64", http.StatusBadRequest)
		return nil, err
	}

	message, err := rt.db.GetMessage(messageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return nil, err
	case err != nil:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
		return nil, err
	}

	isMember, err := rt.db.IsMember(user.UserId, message.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusNotFound)
		rt.baseLogger.Errorf("IsMember: %w", err)
		return nil, err
	case !isMember:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
		return nil, err
	}

	return message, nil
}

func (rt *_router) checkGroup(w http.ResponseWriter, groupId int64) (*database.Group, error) {
	group, err := rt.db.GetGroupById(groupId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Group id %d not found", groupId), http.StatusNotFound)
		return nil, err
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetGroupById: %w", err)
		return nil, err
	default:
		return group, nil
	}
}

func (rt *_router) checkFounder(w http.ResponseWriter, groupId int64, userId int64) (bool, error) {
	isFounder, err := rt.db.IsFounder(groupId, userId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsFounder: %w", err)
		return false, err
	case isFounder:
		return true, err
	default:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
		return false, nil
	}
}

func (rt *_router) checkUserId(w http.ResponseWriter, r *http.Request) (*int64, error) {
	body := struct {
		UserId int64 `json:"userId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		rt.sendResponse(
			w,
			"Invalid format. Field 'userId' of type int64 is required",
			http.StatusBadRequest,
		)
		return nil, err
	}

	return &body.UserId, err
}

func (rt *_router) checkEmojiCode(w http.ResponseWriter, r *http.Request) (*database.EmojiCode, error) {
	body := struct {
		EmojiCode string `json:"emojiCode"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		rt.sendResponse(
			w,
			"Invalid format. Field 'emojiCode' of type string is required",
			http.StatusBadRequest,
		)
		return nil, err
	}

	emojiCode := database.EmojiCode(body.EmojiCode)

	err = database.IsValidEmojiCode(emojiCode)
	if err != nil {
		rt.sendResponse(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &emojiCode, nil
}

func (rt *_router) checkMessageId(w http.ResponseWriter, r *http.Request) (*int64, error) {
	body := struct {
		MessageId int64 `json:"messageId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		rt.sendResponse(
			w,
			"Invalid format. Field 'messageId' of type int64 is required",
			http.StatusBadRequest,
		)
	}

	return &body.MessageId, err
}

func (rt *_router) checkName(w http.ResponseWriter, r *http.Request) (*string, error) {
	body := struct {
		Name string `json:"name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		rt.sendResponse(
			w,
			"Invalid format. Field 'name' of type string is required",
			http.StatusBadRequest,
		)
		return nil, err
	}

	if !nameRegex.MatchString(body.Name) {
		errMsg := "Invalid name format"
		rt.sendResponse(w, errMsg, http.StatusBadRequest)
		return nil, errors.New(errMsg)
	}

	return &body.Name, nil
}
