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
	group, err := rt.checkGroup(w, ps)
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

	group, err := rt.checkGroup(w, ps)
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
	group, err := rt.checkGroup(w, ps)
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

	if _, err := rt.db.SetGroupPhotoUrl(group.ConversationId, photoUrl); err != nil {
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

// operationId: deleteGroup
func (rt *_router) deleteGroup(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	group, err := rt.checkGroup(w, ps)
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
	group, err := rt.checkGroup(w, ps)
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

	group, err := rt.checkGroup(w, ps)
	if err != nil {
		return
	}

	isFounder, err := rt.authUserAsFounder(w, group.ConversationId, user.UserId)
	if !isFounder || err != nil {
		return
	}

	if *userId == user.UserId {
		rt.sendResponse(w, "Cannot add founder to his own group", http.StatusBadRequest)
		return
	}

	isMember, err := rt.db.IsMember(*userId, group.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
	case isMember:
		rt.sendResponse(w, fmt.Sprintf("User %d already in group %d", userId, group.ConversationId), http.StatusBadRequest)
	default:
		_, err = rt.db.AddConversation(group.ConversationId, *userId)
		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("AddConversation: %v", err)
			return
		}

		rt.sendResponse(w, struct {
			UserId int64 `json:"userId"`
		}{UserId: *userId}, http.StatusCreated)
	}
}

// operationId: removeUser
func (rt *_router) removeUser(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, user database.User) {
	userId, err := strconv.ParseInt(ps.ByName("userId"), 10, 64)
	if err != nil {
		rt.sendResponse(w, "Parameter userId must be an int64", http.StatusBadRequest)
		return
	}

	group, err := rt.checkGroup(w, ps)
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

		// NOTE: If the client does not provide the group photo we use the default one.
		var photoUrl string
		if len(r.MultipartForm.File["photo"]) == 0 {
			photoUrl = rt.defaultGroupPhoto
		} else {
			photoUrl, err = rt.uploadMediaFile(w, r, "photo")
			if err != nil {
				return
			}
		}

		group, err := rt.db.CreateGroup(user.UserId, name, photoUrl)
		if err != nil {
			rt.baseLogger.Errorf("CreateGroup: %v", err)
			_ = rt.removeMediaFile(photoUrl)
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rt.sendResponse(w, group, http.StatusCreated)
	}
}

// operationId: forwardMessageToGroup
func (rt *_router) forwardMessageToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	group, err := rt.checkGroup(w, ps)
	if err != nil {
		return
	}

	messageId, err := rt.getMessageIdFromReq(w, r)
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
		rt.baseLogger.Errorf("GetMessage: %v", err)
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, group.ConversationId)
	switch {
	case err != nil:
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
	case isMember:
		fmessage, err := rt.db.InsertMessage(
			message.Text,
			user.UserId,
			group.ConversationId,
			true,
			nil,
			message.AttachmentUrl,
			message.MediaType,
		)

		if err != nil {
			rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
			rt.baseLogger.Errorf("InsertMessage: %v", err)
			return
		}

		rt.sendResponse(w, fmessage, http.StatusCreated)
	default:
		rt.sendResponse(
			w,
			fmt.Sprintf(
				"User %d is not a member of group %d",
				user.UserId,
				group.ConversationId,
			),
			http.StatusUnauthorized,
		)
	}
}

// operationId: sendMessageToGroup
func (rt *_router) sendMessageToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User) {
	group, err := rt.checkGroup(w, ps)
	if err != nil {
		return
	}

	isMember, err := rt.db.IsMember(user.UserId, group.ConversationId)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("IsMember: %v", err)
		return
	}

	if !isMember {
		rt.sendResponse(
			w,
			fmt.Sprintf(
				"User %d is not a member of group %d",
				user.UserId,
				group.ConversationId,
			),
			http.StatusUnauthorized,
		)
		return
	}

	message, err := rt._insertMessage(w, r, user.UserId, group.ConversationId, nil)
	if err != nil {
		return
	}

	rt.sendResponse(w, message, http.StatusCreated)
}
