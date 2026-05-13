package database

import (
	"database/sql"
	"errors"
)

// Returns the user ids associated with the specified conversation id.
func (db *appdbimpl) GetMembers(conversationId int64) ([]int64, error) {
	rows, err := db.c.Query(
		`SELECT user_id FROM user_conversations WHERE conversation_id = ?`,
		conversationId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userIds []int64
	for rows.Next() {
		var userId int64
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}
		userIds = append(userIds, userId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userIds, nil
}

// Returns true if the specified user id is part of the specified conversation id.
func (db *appdbimpl) IsMember(userId int64, conversationId int64) (bool, error) {
	var dummy int64
	err := db.c.QueryRow(
		`SELECT 1 FROM user_conversations WHERE conversation_id = ? AND user_id = ? LIMIT 1`,
		conversationId, userId,
	).Scan(&dummy)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, err
	}
}

// Adds a record with the specified conversation and user id in the user_conversations table.
func (db *appdbimpl) AddConversation(conversationId int64, userId int64) (sql.Result, error) {
	return db.c.Exec(
		`INSERT INTO user_conversations (conversation_id, user_id) VALUES (?, ?)`,
		conversationId, userId,
	)
}

// Deletes the records associated with the specified conversation id from the conversations table.
func (db *appdbimpl) DeleteConversation(conversationId int64) (sql.Result, error) {
	return db.c.Exec(
		`DELETE FROM conversations WHERE conversation_id = ?`,
		conversationId,
	)
}

// Deletes the record associated with the specified conversation and user id from the user_conversations table.
func (db *appdbimpl) DeleteUserConversation(conversationId int64, userId int64) (sql.Result, error) {
	return db.c.Exec(
		`DELETE FROM user_conversations WHERE conversation_id = ? AND user_id = ?`,
		conversationId, userId,
	)
}

// Creates a new conversation between the specified senderId and receiverId.
// A new message record is created inside the messages table with the specified values.
func (db *appdbimpl) CreateConversation(senderId int64, receiverId int64) (*int64, error) {
	tx, err := db.GetTransaction()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.Exec(`INSERT INTO conversations DEFAULT VALUES`)
	if err != nil {
		return nil, err
	}

	conversationId, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO user_conversations (conversation_id, user_id) VALUES (?, ?)`,
		conversationId, senderId,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO user_conversations (conversation_id, user_id) VALUES (?, ?)`,
		conversationId, receiverId,
	)
	if err != nil {
		return nil, err
	}

	return &conversationId, nil
}

// Returns the conversation id between userIdA and userIdB if exists, otherwise nil.
func (db *appdbimpl) GetConversation(userIdA int64, userIdB int64) (*int64, error) {
	conversationIdsA, err := db.GetConversations(userIdA)
	if err != nil || len(conversationIdsA) == 0 {
		return nil, err
	}

	conversationIdsB, err := db.GetConversations(userIdB)
	if err != nil || len(conversationIdsB) == 0 {
		return nil, err
	}

	set := make(map[int64]struct{}, len(conversationIdsA))
	for _, x := range conversationIdsA {
		set[x] = struct{}{}
	}

	for _, conversationId := range conversationIdsB {
		_, ok := set[conversationId]
		if !ok {
			continue
		}

		isGroup, err := db.IsGroup(conversationId)
		if err != nil {
			return nil, err
		}

		if !isGroup {
			return &conversationId, nil
		}
	}

	return nil, nil
}

// Returns the conversation ids associated to the specified user id.
func (db *appdbimpl) GetConversations(userId int64) ([]int64, error) {
	rows, err := db.c.Query(
		`SELECT conversation_id FROM user_conversations WHERE user_id = ?`,
		userId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var conversationIds []int64
	for rows.Next() {
		var conversationId int64
		if err := rows.Scan(&conversationId); err != nil {
			return nil, err
		}
		conversationIds = append(conversationIds, conversationId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversationIds, nil
}
