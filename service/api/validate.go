package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

var nameRegex = regexp.MustCompile(`^.+$`)

func (rt *_router) checkGroup(w http.ResponseWriter, ps httprouter.Params) (*database.Group, error) {
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
