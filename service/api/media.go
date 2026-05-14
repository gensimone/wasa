package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

// Remove the provided file from the filesystem.
func (rt *_router) removeMediaFile(w http.ResponseWriter, path string) error {
	filename := filepath.Base(path)
	fullpath := filepath.Join(rt.rootMedia, filename)
	// FIXME: Check that fullpath is actually a file.

	err := os.Remove(fullpath)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error deleting file %s: %w", fullpath, err)
	}

	return err
}

// Upload the provided file from the filesystem.
func (rt *_router) uploadMediaFile(w http.ResponseWriter, r *http.Request, key string) (*string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		rt.sendResponse(w, "Invalid multipart form", http.StatusBadRequest)
		return nil, err
	}

	file, _, err := r.FormFile(key)
	if err != nil {
		rt.sendResponse(w, fmt.Sprintf("Missing %s field", key), http.StatusBadRequest)
		return nil, err
	}

	defer file.Close()

	filename, err := uuid.NewV4()
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error calling uuid.NewV4: %w", err)
		return nil, err
	}

	fullpath := filepath.Join(rt.rootMedia, filename.String())

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

	url := filepath.Join(rt.media, filename.String())
	return &url, nil
}
