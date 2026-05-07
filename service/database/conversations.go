package database

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/gensimone/WASA-project/service/globaltime"
)

func (db *appdbimpl) GetConversation(userA int64, userB int64) (*int64, error) {
	convsA, err := db.GetConversations(userA)
	if err != nil || len(convsA) == 0 {
		return nil, err
	}

	convsB, err := db.GetConversations(userB)
	if err != nil || len(convsB) == 0 {
		return nil, err
	}

	set := make(map[int64]struct{}, len(convsA))
	for _, x := range convsA {
		set[x] = struct{}{}
	}

	for _, c := range convsB {
		if _, ok := set[c]; ok {
			if yes, err := db.IsGroup(c); err != nil {
				return nil, err
			} else if !yes {
				return &c, nil
			}
		}
	}

	return nil, nil
}

func (db *appdbimpl) CreateConversation(
	sender int64, receiver int64, text string, photo *os.File, isForwarded bool,
) (*Message, error) {
	tx, err := db.GetTransaction()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO conversations DEFAULT VALUES`)
	if err != nil {
		return nil, err
	}

	conv, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO userConversations (conversation, user) VALUES (?, ?)`, conv, sender,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO userConversations (conversation, user) VALUES (?, ?)`, conv, receiver,
	)
	if err != nil {
		return nil, err
	}

	timestamp := globaltime.Now().Format(time.DateTime)
	res, err = tx.Exec(
		`INSERT INTO messages (text, photo, sender, conversation, timestamp, isForwarded, commentTo) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		text, photo, sender, conv, timestamp, isForwarded, nil,
	)
	if err != nil {
		return nil, err
	}

	msg, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Status
	if _, err = tx.Exec(
		`INSERT INTO status (message, user, info) VALUES (?, ?, ?)`,
		msg, receiver, "not received",
	); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Message{
		Id:           msg,
		Text:         text,
		Photo:        photo,
		Sender:       sender,
		Conversation: conv,
		Timestamp:    timestamp,
		IsForwarded:  isForwarded,
		CommentTo:    nil,
	}, nil
}

func (db *appdbimpl) AddConversation(conv int64, user int64) (sql.Result, error) {
	return db.c.Exec(`INSERT INTO userConversations (conversation, user) VALUES (?, ?)`, conv, user)
}

func (db *appdbimpl) DeleteConversation(conv int64) (sql.Result, error) {
	return db.c.Exec(`DELETE FROM conversations WHERE id = ?`, conv)
}

func (db *appdbimpl) DeleteUserConversation(conv int64, user int64) (sql.Result, error) {
	return db.c.Exec(`DELETE FROM userConversations WHERE conversation = ? AND user = ?`, conv, user)
}

func (db *appdbimpl) IsMember(user int64, conv int64) (bool, error) {
	var x int64
	if err := db.c.QueryRow(
		`SELECT conversation FROM userConversations WHERE conversation = ? AND user = ?`, conv, user,
	).Scan(&x); errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, err
	}
}

func (db *appdbimpl) GetMembers(conv int64) ([]int64, error) {
	rows, err := db.c.Query(`SELECT user FROM userConversations WHERE conversation = ?`, conv)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []int64
	for rows.Next() {
		var user int64
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *appdbimpl) GetConversations(user int64) ([]int64, error) {
	rows, err := db.c.Query(`SELECT conversation FROM userConversations WHERE user = ?`, user)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return ids, nil
}
