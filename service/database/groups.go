package database

import (
	"database/sql"

	"github.com/gensimone/WASA-project/service/globaltime"
)

// Sets the group identified by the group id with the specified group name.
func (db *appdbimpl) SetGroupName(groupId int64, name string) (sql.Result, error) {
	return db.c.Exec(
		`UPDATE groups SET name = ?
		WHERE conversation_id = ?`,
		name,
		groupId,
	)
}

// Sets the photo (URL) of the group identified by the specified group id.
func (db *appdbimpl) SetGroupPhotoUrl(groupId int64, photoUrl string) (sql.Result, error) {
	return db.c.Exec(
		`UPDATE groups
		SET photo_url = ?
		WHERE conversation_id = ?`,
		photoUrl,
		groupId,
	)
}

// Returns true if the user identified by the user id is the founder of the group
// identified by the group id.
func (db *appdbimpl) IsFounder(conversationId, userId int64) (bool, error) {
	var founderId int64
	err := db.c.QueryRow(
		`SELECT founder_id
		FROM groups
		WHERE conversation_id = ?`,
		conversationId,
	).Scan(&founderId)

	return err == nil && founderId == userId, err
}

// Returns true if a group with the specified founder id and name exists,
// otherwise returns false.
func (db *appdbimpl) GroupExists(founderId int64, name string) (bool, error) {
	var exists bool

	err := db.c.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM groups
			WHERE founder_id = ? AND name = ?
		)`,
		founderId,
		name,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// Creates a new record in the groups table with the specified parameters and returns
// a Group struct that describes the created record.
func (db *appdbimpl) CreateGroup(founderId int64, name, photoUrl string) (*Group, error) {
	tx, err := db.GetTransaction()
	if err != nil {
		return nil, err
	}

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
		`INSERT INTO user_conversations (
			conversation_id,
			user_id
		) VALUES (?, ?)`,
		conversationId,
		founderId,
	)
	if err != nil {
		return nil, err
	}

	createdAt := globaltime.Now()
	_, err = tx.Exec(
		`INSERT INTO groups (
			conversation_id,
			founder_id,
			name,
			photo_url,
			created_at
		) VALUES (?, ?, ?, ?, ?)`,
		conversationId,
		founderId,
		name,
		photoUrl,
		createdAt,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Group{
		ConversationId: conversationId,
		FounderId:      founderId,
		Name:           name,
		PhotoUrl:       photoUrl,
		CreatedAt:      createdAt,
	}, nil
}

// Returns the group object associated with the specified group id.
func (db *appdbimpl) GetGroupById(conversationId int64) (*Group, error) {
	var group Group
	if err := db.c.QueryRow(
		`SELECT
		conversation_id,
		founder_id,
		name,
		photo_url,
		created_at
		FROM groups
		WHERE conversation_id = ?`,
		conversationId,
	).Scan(
		&group.ConversationId,
		&group.FounderId,
		&group.Name,
		&group.PhotoUrl,
		&group.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &group, nil
}

// Returns true if the specified conversation id identify a group in the groups table,
// otherwise returns false.
func (db *appdbimpl) IsGroup(conversationId int64) (bool, error) {
	rows, err := db.c.Query(
		`SELECT
		conversation_id
		FROM groups
		WHERE conversation_id = ?`,
		conversationId,
	)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), rows.Err()
}
