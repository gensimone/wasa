package database

import (
	"os"
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserIds() ([]int64, error) {
	rows, err := db.c.Query(`SELECT id FROM users`)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
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

	return ids, nil
}

func (db *appdbimpl) GetUsers() ([]User, error) {
	ids, err := db.GetUserIds()
	if err != nil {
		return nil, err
	}

	var users []User
	for _, id := range ids {
		user, err := db.GetUserById(id)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (db *appdbimpl) GetUserById(id int64) (*User, error) {
	var u User

	if err := db.c.QueryRow(
		`SELECT id, name, photo FROM users WHERE id = ?`, id,
	).Scan(&u.Id, &u.Name, &u.Photo); err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *appdbimpl) GetUserByName(name string) (*User, error) {
	var u User

	if err := db.c.QueryRow(
		`SELECT id, name, photo FROM users WHERE name = ?`, name,
	).Scan(&u.Id, &u.Name, &u.Photo); err != nil {
		return nil, err
	}

	return &u, nil
}


func (db *appdbimpl) SetMyUserName(id int64, name string) error {
	_, err := db.c.Exec(`UPDATE users SET name = ? WHERE id = ?`, name, id)
	return err
}

func (db *appdbimpl) SetMyPhoto(id int64, photo os.File) error {
	_, err := db.c.Exec(`UPDATE users SET photo = ? WHERE id = ?`, photo, id)
	return err
}

func (db *appdbimpl) DeleteUser(id int64) error {
	_, err := db.c.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (db *appdbimpl) DoLogin(name string) (int64, error) {
	var id int64

	err := db.c.QueryRow(`SELECT id FROM users WHERE name = ?`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		res, err := db.c.Exec(`INSERT INTO users (name, photo) VALUES (?, ?)`, name, nil)
		id, err = res.LastInsertId()
		if err != nil {
			return -1, err
		}
		return id, nil
	} else if err != nil {
		return -1, err
	}

	return id, nil
}
