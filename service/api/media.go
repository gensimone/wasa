package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

// Remove the provided file from the filesystem.
func (rt *_router) removeMediaFile(path string) error {
	filename := filepath.Base(path)
	fullpath := filepath.Join(rt.rootMedia, filename)
	err := os.Remove(fullpath)
	if err != nil {
		rt.baseLogger.Errorf("Error deleting file %s: %v", fullpath, err)
	}

	return err
}

// Upload the provided file from the filesystem.
func (rt *_router) uploadMediaFile(w http.ResponseWriter, r *http.Request, key string) (string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		rt.sendResponse(w, "Invalid multipart form", http.StatusBadRequest)
		return "", err
	}

	file, _, err := r.FormFile(key)
	if err != nil {
		rt.sendResponse(w, "No such file", http.StatusBadRequest)
		return "", err
	}

	defer file.Close()

	filename, err := uuid.NewV4()
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error calling uuid.NewV4: %v", err)
		return "", err
	}

	fullpath := filepath.Join(rt.rootMedia, filename.String())

	dst, err := os.Create(fullpath)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error creating file %s: %v", fullpath, err)
		return "", err
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		rt.sendResponse(w, "Internal Server Error", http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error copying file %s: %v", fullpath, err)
		return "", err
	}

	url := filepath.Join(rt.media, filename.String())
	return url, nil
}
