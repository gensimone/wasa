package database

import (
	"os"
	"database/sql"
)

func (db *appdbimpl) SetGroupName(id int64, name string) (sql.Result, error) {
	return db.c.Exec(`UPDATE groups SET name = ? WHERE id = ?`, name, id)
}

func (db *appdbimpl) SetGroupPhoto(id int64, photo os.File) (sql.Result, error) {
	return db.c.Exec(`UPDATE groups SET photo = ? WHERE id = ?`, photo, id)
}

func (db *appdbimpl) DeleteGroup(conversation int64) (sql.Result, error) {
	return db.c.Exec(`DELETE FROM groups WHERE conversation = ?`, conversation)
}

func (db *appdbimpl) IsFounder(conversation int64, user int64) (bool, error) {
	var founder int64
	err := db.c.QueryRow(
		`SELECT founder FROM groups WHERE conversation = ?`, conversation,
	).Scan(&founder)

	return err != nil && founder == user, err
}

func (db *appdbimpl) GetGroupById(conversation int64) (*Group, error) {
	var g Group

	if err := db.c.QueryRow(
		`SELECT conversation, founder, name, photo, timestamp FROM groups WHERE conversation = ?`, conversation,
	).Scan(&g.Conversation, &g.Founder, &g.Name, &g.Photo, &g.Timestamp); err != nil {
		return nil, err
	}

	return &g, nil
}

func (db *appdbimpl) IsGroup(conversation int64) (bool, error) {
	rows, err := db.c.Query(`SELECT conversation FROM groups WHERE conversation = ?`, conversation)
	if err != nil { return false, err }
	defer rows.Close()
	return rows.Next(), rows.Err()
}

