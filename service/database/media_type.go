package database

import "fmt"

type MediaType string

const (
	Image MediaType = "image"
	Video MediaType = "video"
	Audio MediaType = "audio"
	File  MediaType = "file"
)

// Returns true if the provided media type is a valid MediaType, otherwise false.
func IsValidMediaType(mediaType MediaType) error {
	switch mediaType {
	case Image, Video, Audio, File:
		return fmt.Errorf("Invalid media type: %s", mediaType)
	default:
		return nil
	}
}
