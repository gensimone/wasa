package api

import (
	"encoding/json"
	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.GET("/session", rt.doLogin)
	return rt.router
}

type Name struct {
	name string
}

type User struct {
	id int64
	name string
	photo string
}

func (rt *_router) doLogin(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// Extract the provided user name.
    var name_struct Name
    decoder := json.NewDecoder(request.Body)
    err := decoder.Decode(&name_struct)
    if err != nil {
		http.Error(writer, `{"message":"Bad Request"}`, http.StatusBadRequest)
    }

	// Get the user id
	id, err := rt.db.DoLogin(name_struct.name)
	if err != nil {
		http.Error(writer, `{"message":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Write the response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(User{id: id, name: name_struct.name, photo: ""})
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
