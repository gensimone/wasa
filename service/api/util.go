package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (rt *_router) sendResponse(w http.ResponseWriter, content any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if content == nil {
		return
	}

	// For semplicity, string are treated as errors.
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
		rt.baseLogger.Errorf("json.NewEncoder: %v", err)
	}
}

func (rt *_router) checkRowsAffected(w http.ResponseWriter, expected int64, result sql.Result) error {
	ra, err := result.RowsAffected()

	switch {
	case err != nil:
		rt.baseLogger.Errorf("RowsAffected: %v", err)
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return err

	case ra != expected:
		errMsg := fmt.Sprintf("Invalid number of affected rows: %d. Expected: %d", ra, expected)
		rt.baseLogger.Error(errMsg)
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return errors.New(errMsg)

	default:
		return nil
	}
}
