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
	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, group.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
	case isMember:
		rt.sendResponse(w, group, http.StatusOK)
	default:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// operationId: setGroupName
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	name, err := rt.getNameFromReq(w, r)
	if err != nil {
		return
	}

	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	isFounder, err := rt.authUserAsFounder(w, group.ConversationId, user.UserId)
	if !isFounder || err != nil {
		return
	}

	_, err = rt.db.SetGroupName(group.ConversationId, *name)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetGroupName: %v", err)
		return
	}

	rt.sendResponse(w, struct {
		Name string `json:"name"`
	}{Name: *name}, http.StatusOK)
}

// operationId: setGroupPhoto
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	isFounder, err := rt.authUserAsFounder(w, group.ConversationId, user.UserId)
	if !isFounder || err != nil {
		return
	}

	photoUrl, err := rt.uploadMediaFile(w, r, "photo")
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("uploadFile: %v", err)
		return
	}

	_, err = rt.db.SetGroupPhotoUrl(group.ConversationId, &photoUrl)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetGroupPhotoUrl: %v", err)
		return
	}

	// NOTE: Here we must delete the old photo.
	//       Or maybe we could remove old photos (or even attachments) somewhere else (e.g. a script)
	rt.sendResponse(w, struct {
		PhotoUrl string `json:"photoUrl"`
	}{PhotoUrl: photoUrl}, http.StatusOK)
}

// operationId: getMemberIds
func (rt *_router) getMemberIds(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	groupId, err := rt.authConversationAccessParam(w, ps, user, "groupId")
	if err != nil {
		return
	}

	userIds, err := rt.db.GetMembers(*groupId)
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetMembers: %v", err)
		return
	}

	rt.sendResponse(w, struct {
		UserIds []int64 `json:"userIds"`
	}{UserIds: userIds}, http.StatusOK)
}

// operationId: deleteGroupPhoto
func (rt *_router) deleteGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	isFounder, err := rt.authUserAsFounder(w, group.ConversationId, user.UserId)
	if !isFounder || err != nil {
		return
	}

	_, err = rt.db.SetGroupPhotoUrl(group.ConversationId, nil)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("SetGroupPhotoUrl: %v", err)
		return
	}

	rt.sendResponse(w, nil, http.StatusNoContent)
}

// operationId: deleteGroup
func (rt *_router) deleteGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	isFounder, err := rt.authUserAsFounder(w, group.ConversationId, user.UserId)
	if !isFounder || err != nil {
		return
	}

	_, err = rt.db.DeleteConversation(group.ConversationId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("DeleteConversation: %v", err)
		return
	}

	rt.sendResponse(w, nil, http.StatusNoContent)
}

// operationId: leaveGroup
func (rt *_router) leaveGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, group.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
	case isMember:
		if group.FounderId == user.UserId {
			_, err := rt.db.DeleteConversation(group.ConversationId)
			if err != nil {
				rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("DeleteConversation: %v", err)
				return
			}
		} else {
			_, err := rt.db.DeleteUserConversation(group.ConversationId, user.UserId)
			if err != nil {
				rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
				rt.baseLogger.Errorf("DeleteUserConversation: %v", err)
				return
			}
		}
		rt.sendResponse(w, nil, http.StatusNoContent)
	default:
		rt.sendResponse(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// operationId: addToGroup
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	userId, err := rt.getUserIdFromReq(w, r)
	if err != nil {
		return
	}

	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	err = rt.authConversationAccess(w, user, group.ConversationId)
	if err != nil {
		return
	}

	userToAdd, err := rt.db.GetUserById(*userId)
	if errors.Is(err, sql.ErrNoRows) {
		rt.sendResponse(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetUserById: %v", err)
		return
	}

	isMember, err := rt.db.IsMember(*userId, group.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
	case isMember:
		rt.sendResponse(w, "User already in the group", http.StatusBadRequest)
	default:
		_, err = rt.db.AddConversation(group.ConversationId, *userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("AddConversation: %v", err)
			return
		}

		rt.sendResponse(w, userToAdd, http.StatusCreated)
	}
}

// operationId: removeUser
func (rt *_router) removeUser(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	group, err := rt.IsValidGroup(w, ps)
	if err != nil {
		return
	}

	isFounder, err := rt.authUserAsFounder(w, group.ConversationId, user.UserId)
	if !isFounder || err != nil {
		return
	}

	if userId == user.UserId {
		rt.sendResponse(w, "Cannot remove founder. Use leaveGroup instead", http.StatusBadRequest)
		return
	}

	isMember, err := rt.db.IsMember(userId, group.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
	case isMember:
		_, err = rt.db.DeleteUserConversation(group.ConversationId, userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("DeleteUserConversation: %v", err)
			return
		}

		rt.sendResponse(w, nil, http.StatusNoContent)
	default:
		rt.sendResponse(w, fmt.Sprintf("User %d is not a member of group %d", userId, group.ConversationId), http.StatusBadRequest)
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
		rt.baseLogger.Errorf("GroupExists: %v", err)
	case exists:
		rt.sendResponse(w, "Group name already used", http.StatusConflict)
	default:
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			rt.sendResponse(w, "Invalid multipart form", http.StatusBadRequest)
			return
		}

		var photoUrl string
		if len(r.MultipartForm.File["photo"]) > 0 {
			photoUrl, err = rt.uploadMediaFile(w, r, "photo")
			if err != nil {
				return
			}
		}

		group, err := rt.db.CreateGroup(user.UserId, name, &photoUrl)
		if err != nil {
			rt.baseLogger.Errorf("CreateGroup: %v", err)
			_ = rt.removeMediaFile(photoUrl)
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rt.sendResponse(w, group, http.StatusCreated)
	}
}

func (rt *_router) IsValidGroup(w http.ResponseWriter, ps httprouter.Params) (*database.Group, error) {
	groupId, err := strconv.ParseInt(ps.ByName("groupId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter groupId must be an int64", http.StatusBadRequest)
		return nil, err
	}

	group, err := rt.db.GetGroupById(groupId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		rt.sendResponse(w, fmt.Sprintf("Group id %d not found", groupId), http.StatusNotFound)
		return nil, err
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("GetGroupById: %v", err)
		return nil, err
	default:
		return group, nil
	}
}
