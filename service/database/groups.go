package database

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/gensimone/WASA-project/service/globaltime"
)

func (db *appdbimpl) SetGroupName(id int64, name string) (sql.Result, error) {
	return db.c.Exec(`UPDATE groups SET name = ? WHERE id = ?`, name, id)
}

func (db *appdbimpl) SetGroupPhoto(id int64, photo os.File) (sql.Result, error) {
	return db.c.Exec(`UPDATE groups SET photo = ? WHERE id = ?`, photo, id)
}

func (db *appdbimpl) IsFounder(conv int64, user int64) (bool, error) {
	var founder int64
	err := db.c.QueryRow(
		`SELECT founder FROM groups WHERE conversation = ?`, conv,
	).Scan(&founder)

	return err == nil && founder == user, err
}

func (db *appdbimpl) GroupExists(founder int64, name string) (bool, error) {
	var conv int64
	if err := db.c.QueryRow(
		`SELECT conversation FROM groups WHERE founder = ? AND name = ?`, founder, name,
	).Scan(&conv); errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (db *appdbimpl) CreateGroup(founder int64, name string, photo *os.File) (*Group, error) {
	tx, err := db.GetTransaction()
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(`INSERT INTO conversations DEFAULT VALUES`)
	if err != nil {
		return nil, err
	}

	conv, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res, err = tx.Exec(
		`INSERT INTO userConversations (conversation, user) VALUES (?, ?)`, conv, founder,
	)
	if err != nil {
		return nil, err
	}

	timestamp := globaltime.Now().Format(time.DateTime)
	res, err = tx.Exec(
		`INSERT INTO groups (conversation, founder, name, photo, timestamp) VALUES (?, ?, ?, ?, ?)`,
		conv, founder, name, photo, timestamp,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Group{
		Conversation: conv,
		Founder: founder,
		Name: name,
		Photo: photo,
		Timestamp: timestamp,
	}, nil
}

func (db *appdbimpl) GetGroupById(conv int64) (*Group, error) {
	var g Group
	if err := db.c.QueryRow(
		`SELECT conversation, founder, name, photo, timestamp FROM groups WHERE conversation = ?`, conv,
	).Scan(&g.Conversation, &g.Founder, &g.Name, &g.Photo, &g.Timestamp); err != nil {
		return nil, err
	}

	return &g, nil
}

func (db *appdbimpl) IsGroup(conv int64) (bool, error) {
	rows, err := db.c.Query(`SELECT conversation FROM groups WHERE conversation = ?`, conv)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), rows.Err()
}
