package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gensimone/WASA-project/service/database"
)

func (rt *_router) getUserIdFromReq(w http.ResponseWriter, r *http.Request) (*int64, error) {
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

func (rt *_router) getEmojiFromReq(w http.ResponseWriter, r *http.Request) (*database.Emoji, error) {
	body := struct {
		Emoji string `json:"emoji"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		rt.sendResponse(
			w,
			"Invalid format. Field 'emoji' of type string is required",
			http.StatusBadRequest,
		)
		return nil, err
	}

	emoji := database.Emoji(body.Emoji)

	return &emoji, nil
}

func (rt *_router) getMessageIdFromReq(w http.ResponseWriter, r *http.Request) (*int64, error) {
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

func (rt *_router) getNameFromReq(w http.ResponseWriter, r *http.Request) (*string, error) {
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
