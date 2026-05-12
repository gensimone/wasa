package api

import (
	"encoding/json"
	"net/http"
)

func (rt *_router) sendResponse(w http.ResponseWriter, content any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if content == nil {
		return
	}

	if s, ok := content.(string); ok {
		content = struct {
			Error string `json:"error"`
		}{
			Error: s,
		}
	}

	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		rt.sendResponse(w, "Internal Sever Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("json.NewEncoder: %w", err)
	}
}
