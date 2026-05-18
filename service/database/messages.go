package database

import (
	"database/sql"

	"github.com/gensimone/WASA-project/service/globaltime"
)

// Insert a message in the messages table with the specified parameters.
func (db *appdbimpl) InsertMessage(
	text string,
	senderId int64,
	conversationId int64,
	isForwarded bool,
	commentTo *int64,
	attachmentUrl string,
	mediaType MediaType,
) (*Message, error) {
	tx, err := db.GetTransaction()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	createdAt := globaltime.Now()
	res, err := tx.Exec(
		`INSERT INTO messages (
			text,
			sender_id,
			conversation_id,
			created_at,
			is_forwarded,
			comment_to,
			attachment_url,
			media_type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		text,
		senderId,
		conversationId,
		createdAt,
		isForwarded,
		&commentTo,
		attachmentUrl,
		mediaType,
	)

	if err != nil {
		return nil, err
	}

	messageId, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	memberIds, err := db.GetMembers(conversationId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for _, memberId := range memberIds {
		if memberId == senderId {
			continue
		}

		_, err = tx.Exec(
			`INSERT INTO receipts (
				message_id,
				user_id,
				status,
				sent_at
			) VALUES (?, ?, ?, ?)`,
			messageId,
			memberId,
			Sent,
			createdAt, // send_at = created_at
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Message{
		MessageId:      messageId,
		Text:           text,
		SenderId:       senderId,
		ConversationId: conversationId,
		CreatedAt:      createdAt,
		IsForwarded:    isForwarded,
		CommentTo:      commentTo,
		AttachmentUrl:  attachmentUrl,
		MediaType:      mediaType,
	}, nil
}

// Deletes the message record from the messages table associated with the specified message id.
func (db *appdbimpl) DeleteMessage(messageId int64) (sql.Result, error) {
	return db.c.Exec(
		`DELETE FROM messages WHERE message_id = ?`,
		messageId,
	)
}

// Returns the message object associated with the specified message id.
func (db *appdbimpl) GetMessage(messageId int64) (*Message, error) {
	var message Message

	if err := db.c.QueryRow(
		`SELECT
		message_id,
		text,
		sender_id,
		conversation_id,
		created_at,
		is_forwarded,
		comment_to,
		attachment_url,
		media_type
		FROM messages
		WHERE message_id = ?`,
		messageId,
	).Scan(
		&message.MessageId,
		&message.Text,
		&message.SenderId,
		&message.ConversationId,
		&message.CreatedAt,
		&message.IsForwarded,
		&message.CommentTo,
		&message.AttachmentUrl,
		&message.MediaType,
	); err != nil {
		return nil, err
	}

	return &message, nil
}

// Returns the message ids associated with the specified conversation id.
func (db *appdbimpl) GetMessageIds(conversationId int64) ([]int64, error) {
	rows, err := db.c.Query(
		`SELECT message_id FROM messages WHERE conversation_id = ?`,
		conversationId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messageIds []int64
	for rows.Next() {
		var messageId int64
		if err := rows.Scan(&messageId); err != nil {
			return nil, err
		}
		messageIds = append(messageIds, messageId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messageIds, nil
}
