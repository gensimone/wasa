package database

import (
	"database/sql"
	"time"

	"github.com/gensimone/WASA-project/service/globaltime"
)

// Insert a message in the messages table with the specified parameters.
func (db *appdbimpl) InsertMessage(
	senderId int64,
	conversationId int64,
	text string,
	attachmentId *int64,
	isForwarded bool,
	commentTo *int64,
) (*Message, error) {
	tx, err := db.GetTransaction()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	createdAt := globaltime.Now().Format(time.DateTime)
	res, err := tx.Exec(
		`INSERT INTO messages (text, attachment_id, sender_id, conversation_id, created_at, is_forwarded, comment_to) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		text, &attachmentId, senderId, conversationId, createdAt, isForwarded, &commentTo,
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

		if _, err = tx.Exec(
			`INSERT INTO status (message_id, user_id, info) VALUES (?, ?, ?)`,
			messageId, memberId, "not received",
		); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Message{
		MessageId:      messageId,
		Text:           text,
		AttachmentId:   attachmentId,
		SenderId:       senderId,
		ConversationId: conversationId,
		CreatedAt:      createdAt,
		IsForwarded:    isForwarded,
		CommentTo:      commentTo,
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
		`SELECT message_id, text, attachment_id, sender_id, conversation_id, created_at, is_forwarded, comment_to FROM messages WHERE message_id = ?`,
		messageId,
	).Scan(
		&message.MessageId,
		&message.Text,
		&message.AttachmentId,
		&message.SenderId,
		&message.ConversationId,
		&message.CreatedAt,
		&message.IsForwarded,
		&message.CommentTo,
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
