/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// DoLogin(string) (int64)
}

type appdbimpl struct {
	c *sql.DB
}

func (db *appdbimpl) DoLogin(name string) (int64, error) {
	var id int64
	err := db.c.QueryRow(`SELECT id FROM user WHERE username = ?`, name).Scan(&id)
	switch err {
	case nil:
		return id, nil
	case sql.ErrNoRows:
		result, err := db.c.Exec(`INSERT INTO user (name, photo) VALUES (?, ?)`, name, nil)
		if err == nil {
			return -1, err
		}
		id, err = result.LastInsertId()
		if err == nil {
			return -1, err
		}
		return id, nil
	default:
		return -1, err
	}
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		tables := []string{
			`CREATE TABLE IF NOT EXISTS user (
                id INTEGER PRIMARY KEY,
                name TEXT,
                photo TEXT
            );`,
			// `CREATE TABLE conversation (
			//              id,
			// 	userId
			//          );`,
			// `CREATE TABLE group (
			// 	conversationId,
			// 	userId,
			// 	name,
			// 	photo,
			// 	timestamp
			//          );`,
			// `CREATE TABLE message (
			//              id,
			//              content,
			// 	userId,
			//              conversationId,
			// 	timestamp,
			// 	isForwarded,
			// 	messageId
			//          );`,
			// `CREATE TABLE reaction (
			// 	emoji,
			// 	messageId,
			// 	userId
			//          );`,
		}
		for _, stmt := range tables {
			if _, err := db.Exec(stmt); err != nil {
				return nil, fmt.Errorf("error creating database structure: %w", err)
			}
		}
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error checking database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}
