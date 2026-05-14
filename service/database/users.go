package database

import (
	"database/sql"
	"errors"
)

// Sets the name of the user identified by the specified user id.
func (db *appdbimpl) SetMyUserName(userId int64, name string) (sql.Result, error) {
	return db.c.Exec(
		`UPDATE users SET name = ? WHERE user_id = ?`,
		name, userId,
	)
}

// Sets the photo of the user identified by the specified user id.
func (db *appdbimpl) SetMyPhotoUrl(userId int64, photoUrl string) (sql.Result, error) {
	return db.c.Exec(
		`UPDATE users SET photo_url = ? WHERE user_id = ?`,
		photoUrl, userId,
	)
}

// Creates a new user with the specified name and returns its user id.
func (db *appdbimpl) CreateUser(name string, photoUrl string) (*int64, error) {
	if res, err := db.c.Exec(
		`INSERT INTO users (name, photo_url) VALUES (?, ?)`,
		name, photoUrl,
	); err != nil {
		return nil, err
	} else if userId, err := res.LastInsertId(); err != nil {
		return nil, err
	} else {
		return &userId, nil
	}
}

// Returns all the user ids found in the users table.
func (db *appdbimpl) GetUserIds() ([]int64, error) {
	rows, err := db.c.Query(`SELECT user_id FROM users`)
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

// Returns all the users found in the users table.
func (db *appdbimpl) GetUsers() ([]User, error) {
	userIds, err := db.GetUserIds()
	if err != nil {
		return nil, err
	}

	var users []User
	for _, userId := range userIds {
		user, err := db.GetUserById(userId)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

// Returns the user object associated with the specified user id.
func (db *appdbimpl) GetUserById(userId int64) (*User, error) {
	var user User

	if err := db.c.QueryRow(
		`SELECT user_id, name, photo_url FROM users WHERE user_id = ?`,
		userId,
	).Scan(
		&user.UserId,
		&user.Name,
		&user.PhotoUrl,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *appdbimpl) IsUserById(userId int64) (bool, error) {
	var dummy int64
	err := db.c.QueryRow(
		`SELECT 1 FROM users WHERE user_id = ? LIMIT 1`,
		userId,
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

// Returns the user object associated with the specified username.
func (db *appdbimpl) GetUserByName(name string) (*User, error) {
	var user User

	if err := db.c.QueryRow(
		`SELECT user_id, name, photo_url FROM users WHERE name = ?`,
		name,
	).Scan(
		&user.UserId,
		&user.Name,
		&user.PhotoUrl,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

// Deletes the user specified by the user id.
func (db *appdbimpl) DeleteUser(userId int64) (sql.Result, error) {
	return db.c.Exec(
		`DELETE FROM users WHERE user_id = ?`,
		userId,
	)
}
