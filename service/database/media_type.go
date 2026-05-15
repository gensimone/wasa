package database

type MediaType string

const (
	Image MediaType = "image"
	Video MediaType = "video"
	Audio MediaType = "audio"
	File  MediaType = "file"
)

var validMediaTypes = map[MediaType]struct{}{
	Image: {},
	Video: {},
	Audio: {},
	File:  {},
}

func IsValidMediaType(mediaType MediaType) bool {
	_, ok := validMediaTypes[mediaType]
	return ok
}
