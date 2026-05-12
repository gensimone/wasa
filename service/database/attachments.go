package database

import "database/sql"

type MediaType string

const (
	Image MediaType = "image"
	Video MediaType = "video"
	Audio MediaType = "audio"
	File  MediaType = "file"
)

// Returns true if the provided string is a valid MediaType, otherwise false.
func IsValidMediaType(s string) bool {
	switch MediaType(s) {
	case Image, Video, Audio, File:
		return true
	default:
		return false
	}
}

// Adds a new attachment with the specified url and media type and returns
// the attachment id of the created record.
func (db *appdbimpl) AddAttachment(url string, mediaType MediaType) (*int64, error) {
	if res, err := db.c.Exec(
		`INSERT INTO attachments (url, media_type) VALUES (?, ?)`,
		url, mediaType,
	); err != nil {
		return nil, err
	} else if attachmentId, err := res.LastInsertId(); err != nil {
		return nil, err
	} else {
		return &attachmentId, nil
	}
}

// Deletes the attachment identified by the specified attachment id.
func (db *appdbimpl) DeleteAttachment(attachmentId int64) (sql.Result, error) {
	return db.c.Exec(
		`DELETE FROM attachments WHERE attachment_id = ?`,
		attachmentId,
	)
}

// Returns the attachment associated with the specified attachment id.
func (db *appdbimpl) GetAttachment(attachmentId int64) (*Attachment, error) {
	var attachment Attachment

	if err := db.c.QueryRow(
		`SELECT attachment_id, url, media_type FROM attachments WHERE attachment_id = ?`,
		attachmentId,
	).Scan(
		&attachment.AttachmentId,
		&attachment.Url,
		&attachment.MediaType,
	); err != nil {
		return nil, err
	}

	return &attachment, nil
}
