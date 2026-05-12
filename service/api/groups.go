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

// operationId: getGroup
func (rt *_router) getGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	group, err := rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, groupId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		rt.sendResponse(w, group, http.StatusOK)
	default:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// operationId: setGroupName
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	name, err := rt.checkName(w, r)
	if err != nil {
		return
	}

	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	_, err = rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	if isFounder, err := rt.checkFounder(w, groupId, user.UserId); !isFounder || err != nil {
		return
	}

	_, err = rt.db.SetGroupName(groupId, *name)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetGroupName: %w", err)
		return
	}

	rt.sendResponse(w, name, http.StatusOK)
}

// operationId: setGroupPhoto
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	_, err = rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	if isFounder, err := rt.checkFounder(w, groupId, user.UserId); !isFounder || err != nil {
		return
	}

	photoUrl, err := rt.uploadFile(w, r, "photo")
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("uploadFile: %w", user.UserId, err)
		return
	}

	if _, err := rt.db.SetGroupPhotoUrl(groupId, *photoUrl); err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetGroupPhotoUrl: %w", user.UserId, err)
		return
	}

	rt.sendResponse(w, photoUrl, http.StatusOK)
}

// operationId: deleteGroup
func (rt *_router) deleteGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	_, err = rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	if isFounder, err := rt.checkFounder(w, groupId, user.UserId); !isFounder || err != nil {
		return
	}

	_, err = rt.db.DeleteConversation(groupId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteConversation: %w", err)
		return
	}

	rt.sendResponse(w, nil, http.StatusNoContent)
}

// operationId: leaveGroup
func (rt *_router) leaveGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	group, err := rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, groupId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		if group.FounderId == user.UserId {
			_, err := rt.db.DeleteConversation(groupId)
			if err != nil {
				rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("DeleteConversation: %w", err)
				return
			}

			rt.sendResponse(w, nil, http.StatusNoContent)
			return
		}

		_, err := rt.db.DeleteUserConversation(groupId, user.UserId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("DeleteUserConversation: %w", err)
			return
		}

		rt.sendResponse(w, nil, http.StatusNoContent)

	default:
		rt.sendResponse(w, "Bad Request", http.StatusBadRequest)
	}
}

// operationId: addToGroup
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	_, err = rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	userId, err := rt.checkUserId(w, r)
	if err != nil {
		return
	}

	if isFounder, err := rt.checkFounder(w, groupId, user.UserId); !isFounder || err != nil {
		return
	}

	if *userId == user.UserId {
		rt.sendResponse(w, "Cannot add founder to his own group", http.StatusBadRequest)
		return
	}

	isMember, err := rt.db.IsMember(*userId, groupId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		rt.sendResponse(w, fmt.Sprintf("User %d already in group %d", userId, groupId), http.StatusBadRequest)
	default:
		_, err = rt.db.AddConversation(groupId, *userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("AddConversation: %w", err)
			return
		}

		rt.sendResponse(w, struct {
			UserId int64 `json:"userId"`
		}{UserId: *userId}, http.StatusCreated)
	}
}

// operationId: removeUser
func (rt *_router) removeUser(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	_, err = rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	if isFounder, err := rt.checkFounder(w, groupId, user.UserId); !isFounder || err != nil {
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, groupId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		_, err = rt.db.DeleteUserConversation(groupId, userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("DeleteUserConversation: %w", err)
			return
		}

		rt.sendResponse(w, nil, http.StatusNoContent)
	default:
		rt.sendResponse(w, fmt.Sprintf("User %d is not a member of group %d", userId, groupId), http.StatusBadRequest)
	}
}

// operationId: createGroup
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	name := r.FormValue("name")
	if name == "" {
		rt.sendResponse(w, "Missing group name", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.GroupExists(user.UserId, name)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GroupExists: %w", err)
	case exists:
		rt.sendResponse(w, "Group name already used", http.StatusBadRequest)
	default:
		photoUrl, err := rt.uploadFile(w, r, "photo")
		if err != nil {
			return
		}

		group, err := rt.db.CreateGroup(user.UserId, name, *photoUrl)
		if err != nil {
			rt.baseLogger.Errorf("CreateGroup: %w", err)

			if err = rt.removeFile(w, *photoUrl); err != nil { // Issue Internal Server Error
				return
			}

			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rt.sendResponse(w, group, http.StatusCreated)
	}
}

// operationId: forwardMessageToGroup
func (rt *_router) forwardMessageToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	_, err = rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	messageId, err := rt.checkMessageId(w, r)
	if err != nil {
		return
	}

	message, err := rt.db.GetMessage(*messageId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Message %d not found", messageId), http.StatusNotFound)
		return
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMessage: %w", err)
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, groupId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
	case isMember:
		fmessage, err := rt.db.InsertMessage(
			user.UserId,
			groupId,
			message.Text,
			message.AttachmentId,
			true,
			nil,
		)

		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("InsertMessage: %w", err)
			return
		}

		rt.sendResponse(w, fmessage, http.StatusCreated)
	default:
		rt.sendResponse(
			w,
			fmt.Sprintf("User %d is not a member of group %d", user.UserId, groupId),
			http.StatusUnauthorized,
		)
	}
}

// operationId: sendMessageToGroup
func (rt *_router) sendMessageToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return
	}

	_, err = rt.checkGroup(w, groupId)
	if err != nil {
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, groupId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %w", err)
		return
	}

	if !isMember {
		rt.sendResponse(
			w,
			fmt.Sprintf("User %d is not a member of group %d", user.UserId, groupId),
			http.StatusUnauthorized,
		)
		return
	}

	message, err := rt._insertMessage(w, r, user.UserId, groupId, nil)
	if err != nil {
		return
	}

	rt.sendResponse(w, message, http.StatusCreated)
}
