package database

import (
	"database/sql"
	"errors"
	"os"
)

type AppDatabase interface {
	// Messages.
	InsertMessage(int64, int64, string, *os.File, bool, *int64) (*Message, error)
	DeleteMessage(int64) (sql.Result, error)
	GetMessage(int64) (*Message, error)
	GetStatus(int64) ([]Status, error)
	GetStatusOf(int64, int64) (*Status, error)
	UpdateStatus(int64, int64, string) (sql.Result, error)
	AddReaction(string, int64, int64) (sql.Result, error)
	UpdateReaction(string, int64, int64) (sql.Result, error)
	DeleteReaction(int64, int64) (sql.Result, error)
	GetReactions(int64) ([]Reaction, error)
	GetReaction(int64, int64) (*Reaction, error)

	// Conversations.
	GetMembers(int64) ([]int64, error)
	IsMember(int64, int64) (bool, error)
	DeleteConversation(int64) (sql.Result, error)
	DeleteUserConversation(int64, int64) (sql.Result, error)
	CreateConversation(int64, int64, string, *os.File, bool) (*Message, error)
	GetConversation(int64, int64) (*int64, error)
	GetConversations(int64) ([]int64, error)
	GetMessages(int64) ([]int64, error)

	// Groups.
	SetGroupName(int64, string) (sql.Result, error)
	SetGroupPhoto(int64, os.File) (sql.Result, error)
	DeleteGroup(int64) (sql.Result, error)
	GetGroupById(int64) (*Group, error)
	IsFounder(int64, int64) (bool, error)
	IsGroup(int64) (bool, error)

	// Users.
	SetMyPhoto(int64, os.File) (sql.Result, error)
	SetMyUserName(int64, string) (sql.Result, error)
	GetUserIds() ([]int64, error)
	GetUsers() ([]User, error)
	GetUserById(int64) (*User, error)
	GetUserByName(string) (*User, error)
	DeleteUser(int64) (sql.Result, error)
	DoLogin(string) (*int64, error)
}

type appdbimpl struct {
	c *sql.DB
}

func (db *appdbimpl) GetTransaction() (*sql.Tx, error) {
	return db.c.Begin()
}

func New(db *sql.DB) (AppDatabase, error) {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		stmts := []string{
			`CREATE TABLE users (
                id INTEGER PRIMARY KEY,
                name TEXT,
                photo BLOB
            );`,
			`CREATE TABLE groups (
                conversation INTEGER PRIMARY KEY,
                founder INTEGER,
				name TEXT,
                photo BLOB,
				timestamp TIMESTAMP,
				FOREIGN KEY (conversation) REFERENCES conversations(id) ON DELETE CASCADE,
				FOREIGN KEY (founder) REFERENCES users(id) ON DELETE CASCADE
            );`,
			`CREATE TABLE messages (
                id INTEGER PRIMARY KEY,
                text INTEGER,
				photo BLOB,
                sender INTEGER,
				conversation INTEGER,
				timestamp TIMESTAMP,
				isForwarded BOOL DEFAULT 0,
				commentTo INTEGER,
				FOREIGN KEY (conversation) REFERENCES conversations(id),
				FOREIGN KEY (sender) REFERENCES users(id),
				FOREIGN KEY (commentTo) REFERENCES messages(id)
            );`,
			`CREATE TABLE status (
				message INTEGER,
				user INTEGER,
                info TEXT NOT NULL CHECK (info IN ('read', 'not received', 'received')),
				PRIMARY KEY (message, user),
                FOREIGN KEY (message) REFERENCES messages(id) ON DELETE CASCADE
				FOREIGN KEY (user) REFERENCES users(id)
            );`,
			`CREATE TABLE conversations (
				id INTEGER PRIMARY KEY
			);`,
			`CREATE TABLE userConversations (
				conversation INTEGER,
			    user INTEGER,
				PRIMARY KEY (conversation, user),
                FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE,
				FOREIGN KEY (conversation) REFERENCES conversations(id) ON DELETE CASCADE
            );`,
			`CREATE TABLE reactions (
				emoji TEXT,
				message INTEGER,
				sender INTEGER,
				PRIMARY KEY (message, sender),
                FOREIGN KEY (message) REFERENCES messages(id) ON DELETE CASCADE,
				FOREIGN KEY (sender) REFERENCES users(id) ON DELETE CASCADE
            );`,
		}
		for _, stmt := range stmts {
			if _, err := db.Exec(stmt); err != nil {
				return nil, err
			}
		}
	} else if err != nil {
		return nil, err
	}
	return &appdbimpl{c: db}, nil
}
