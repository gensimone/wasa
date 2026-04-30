package api

import (
	"encoding/json"
	"net/http"
)

func sendResponse(w http.ResponseWriter, v any, s int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		}
	}
}
