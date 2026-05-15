package database

type MediaType string

const (
	Image MediaType = "image"
	Video MediaType = "video"
	Audio MediaType = "audio"
	File  MediaType = "file"
)

type InvalidMediaTypeError struct {
	Type MediaType
}

func (e *InvalidMediaTypeError) Error() string {
	return "Invalid media type: " + string(e.Type)
}

func ValidateMediaType(mediaType MediaType) error {
	switch mediaType {
	case Image, Video, Audio, File:
		return &InvalidMediaTypeError{mediaType}
	default:
		return nil
	}
}
