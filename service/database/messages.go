package database

import (
	"database/sql"
	"os"
	"time"

	"github.com/gensimone/WASA-project/service/globaltime"
)

func (db *appdbimpl) InsertMessage(
	sender int64, conv int64, text string, photo *os.File, isForwarded bool, commentTo *int64,
) (*Message, error) {
	tx, err := db.GetTransaction()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	timestamp := globaltime.Now().Format(time.DateTime)
	res, err := tx.Exec(
		`INSERT INTO messages (text, photo, sender, conversation, timestamp, isForwarded, commentTo) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		text, photo, sender, conv, timestamp, isForwarded, &commentTo,
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
	members, err := db.GetMembers(conv)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	for _, m := range members {
		if m == sender {
			continue
		}
		if _, err = tx.Exec(
			`INSERT INTO status (message, user, info) VALUES (?, ?, ?)`,
			msg, m, "not received",
		); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
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
		CommentTo:    commentTo,
	}, nil
}

func (db *appdbimpl) GetMessages(conv int64) ([]int64, error) {
	rows, err := db.c.Query(`SELECT id FROM messages WHERE conversation = ?`, conv)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (db *appdbimpl) GetMessage(id int64) (*Message, error) {
	var m Message

	if err := db.c.QueryRow(
		`SELECT id, text, photo, sender, conversation, timestamp, isForwarded, commentTo FROM messages WHERE id = ?`, id,
	).Scan(&m.Id, &m.Text, &m.Photo, &m.Sender, &m.Conversation, &m.Timestamp, &m.IsForwarded, &m.CommentTo); err != nil {
		return nil, err
	}

	return &m, nil
}

func (db *appdbimpl) DeleteMessage(id int64) (sql.Result, error) {
	return db.c.Exec(`DELETE FROM messages WHERE id = ?`, id)
}

func (db *appdbimpl) UpdateStatus(msg int64, user int64, info string) (sql.Result, error) {
	return db.c.Exec(`UPDATE status SET info = ? WHERE message = ? AND user = ?`, info, msg, user)
}

func (db *appdbimpl) GetStatusOf(msg int64, user int64) (*Status, error) {
	var s Status
	if err := db.c.QueryRow(
		`SELECT message, user, info FROM status WHERE message = ? AND user = ?`, msg, user,
	).Scan(&s.Message, &s.User, &s.Info); err != nil {
		return nil, err
	}

	return &s, nil
}

func (db *appdbimpl) GetStatus(msg int64) ([]Status, error) {
	rows, err := db.c.Query(`SELECT message, user, info FROM status WHERE message = ?`, msg)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var status []Status
	for rows.Next() {
		var s Status
		if err := rows.Scan(&s.Message, &s.User, &s.Info); err != nil {
			return nil, err
		}
		status = append(status, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return status, nil
}

func (db *appdbimpl) GetReactions(msg int64) ([]Reaction, error) {
	rows, err := db.c.Query(`SELECT emoji, message, sender FROM reactions WHERE message = ?`, msg)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var r Reaction
		if err := rows.Scan(&r.Emoji, &r.Message, &r.Sender); err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

func (db *appdbimpl) GetReaction(msg int64, sender int64) (*Reaction, error) {
	var r Reaction
	if err := db.c.QueryRow(
		`SELECT emoji, message, sender FROM reactions WHERE message = ? AND sender = ?`,
		msg, sender,
	).Scan(&r.Emoji, &r.Message, &r.Sender); err != nil {
		return nil, err
	}

	return &r, nil
}

func (db *appdbimpl) AddReaction(emoji string, msg int64, sender int64) (sql.Result, error) {
	return db.c.Exec(
		`INSERT INTO reactions (emoji, message, sender) VALUES (?, ?, ?)`,
		emoji, msg, sender,
	)
}

func (db *appdbimpl) DeleteReaction(msg int64, sender int64) (sql.Result, error) {
	return db.c.Exec(`DELETE FROM reactions WHERE message = ? AND sender = ?`, msg, sender)
}

func (db *appdbimpl) UpdateReaction(emoji string, msg int64, sender int64) (sql.Result, error) {
	return db.c.Exec(`UPDATE reactions SET emoji = ? WHERE message = ? AND sender = ?`, emoji, msg, sender)
}
