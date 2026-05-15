package database

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
