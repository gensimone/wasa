package database

import "database/sql"

// Returns the reactions associated with the specified message id.
func (db *appdbimpl) GetReactions(messageId int64) ([]Reaction, error) {
	rows, err := db.c.Query(
		`SELECT emoji_code, message_id, sender_id FROM reactions WHERE message_id = ?`,
		messageId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var r Reaction
		if err := rows.Scan(&r.EmojiCode, &r.MessageId, &r.SenderId); err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

// Returns the reaction associated with the specified message and user id.
func (db *appdbimpl) GetReaction(messageId int64, senderId int64) (*Reaction, error) {
	var reaction Reaction
	if err := db.c.QueryRow(
		`SELECT emoji_code, message_id, sender_id FROM reactions WHERE message_id = ? AND sender_id = ?`,
		messageId, senderId,
	).Scan(
		&reaction.EmojiCode,
		&reaction.MessageId,
		&reaction.SenderId,
	); err != nil {
		return nil, err
	}

	return &reaction, nil
}

// Adds a reaction to the message specified by the message id with the specified parameters.
func (db *appdbimpl) AddReaction(emojiCode EmojiCode, messageId int64, senderId int64) (sql.Result, error) {
	// FIXME: We should check the validity of emojiCode here instead of service/api/*
	return db.c.Exec(
		`INSERT INTO reactions (emoji_code, message_id, sender_id) VALUES (?, ?, ?)`,
		emojiCode, messageId, senderId,
	)
}

// Deletes the reaction associated with the specified message and user id.
func (db *appdbimpl) DeleteReaction(messageId int64, senderId int64) (sql.Result, error) {
	return db.c.Exec(
		`DELETE FROM reactions WHERE message_id = ? AND sender_id = ?`,
		messageId, senderId,
	)
}

// Updates the reaction specified by the message and user id with the specified emoji code.
func (db *appdbimpl) UpdateReaction(emojiCode EmojiCode, messageId int64, senderId int64) (sql.Result, error) {
	// FIXME: We should check the validity of emojiCode here instead of service/api/*
	return db.c.Exec(
		`UPDATE reactions SET emojiCode = ? WHERE message_id = ? AND sender_id = ?`,
		emojiCode, messageId, senderId,
	)
}
