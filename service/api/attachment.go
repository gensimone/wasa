package api

import (
	"net/http"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
)

// TODO: implement this
// operationId: getAttachment
func (rt *_router) getAttachment(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params, user database.User,
) {
}
