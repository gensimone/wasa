package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

func (rt *_router) removeFile(w http.ResponseWriter, path string) error {
	err := os.Remove(path)
	if err == nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error deleting %s: %w", path, err)
	}

	return err
}

// Uploads the file inside the multipart/form-data of the provided request
// inside the directory rt.uploads using uuid.NewV4 as filename.
func (rt *_router) uploadFile(w http.ResponseWriter, r *http.Request, key string) (*string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		rt.sendResponse(w, "Invalid multipart form", http.StatusBadRequest)
		return nil, err
	}

	file, _, err := r.FormFile(key)
	if err != nil {
		rt.sendResponse(w, fmt.Sprint("Missing '%s'", key), http.StatusBadRequest)
		return nil, err
	}

	defer file.Close()

	filename, err := uuid.NewV4()
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error calling uuid.NewV4: %w", err)
		return nil, err
	}

	fullpath := filepath.Join(rt.uploads, filename.String())

	dst, err := os.Create(fullpath)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error creating file %s: %w", fullpath, err)
		return nil, err
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error copying file %s: %w", fullpath, err)
		return nil, err
	}

	return &fullpath, nil
}
