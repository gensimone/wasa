package database

import "database/sql"

// Returns the reactions associated with the specified message id.
func (db *appdbimpl) GetReactions(messageId int64) ([]Reaction, error) {
	rows, err := db.c.Query(
		`SELECT
		emoji,
		message_id,
		sender_id
		FROM reactions
		WHERE message_id = ?`,
		messageId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var reaction Reaction
		if err := rows.Scan(
			&reaction.Emoji,
			&reaction.MessageId,
			&reaction.SenderId,
		); err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
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
		`SELECT
		emoji,
		message_id,
		sender_id
		FROM reactions
		WHERE message_id = ? AND sender_id = ?`,
		messageId,
		senderId,
	).Scan(
		&reaction.Emoji,
		&reaction.MessageId,
		&reaction.SenderId,
	); err != nil {
		return nil, err
	}

	return &reaction, nil
}

// Adds a reaction to the message specified by the message id with the specified parameters.
func (db *appdbimpl) AddReaction(emoji Emoji, messageId, senderId int64) (sql.Result, error) {
	return db.c.Exec(
		`INSERT INTO reactions (
			emoji,
			message_id,
			sender_id
		) VALUES (?, ?, ?)`,
		emoji,
		messageId,
		senderId,
	)
}

// Deletes the reaction associated with the specified message and user id.
func (db *appdbimpl) DeleteReaction(messageId, senderId int64) (sql.Result, error) {
	return db.c.Exec(
		`DELETE FROM reactions
		WHERE message_id = ? AND sender_id = ?`,
		messageId,
		senderId,
	)
}

// Updates the reaction specified by the message and user id with the specified emoji code.
func (db *appdbimpl) UpdateReaction(emoji Emoji, messageId, senderId int64) (sql.Result, error) {
	return db.c.Exec(
		`UPDATE reactions
		SET emoji = ?
		WHERE message_id = ? AND sender_id = ?`,
		emoji,
		messageId,
		senderId,
	)
}
