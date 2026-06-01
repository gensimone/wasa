package database

import (
	"database/sql"
	"errors"
)

type AppDatabase interface {
	// Users.
	SetMyUserName(int64, string) (sql.Result, error)
	SetMyPhotoUrl(int64, *string) (sql.Result, error)
	CreateUser(string, *string) (*int64, error)
	GetUserIds() ([]int64, error)
	GetUsers() ([]User, error)
	GetUserById(int64) (*User, error)
	IsUserById(int64) (bool, error)
	GetUserByName(string) (*User, error)

	// Groups.
	SetGroupName(int64, string) (sql.Result, error)
	SetGroupPhotoUrl(int64, *string) (sql.Result, error)
	IsFounder(int64, int64) (bool, error)
	GroupExists(int64, string) (bool, error)
	CreateGroup(int64, string, *string) (*Group, error)
	GetGroupById(int64) (*Group, error)
	IsGroup(int64) (bool, error)

	// Conversations.
	GetMembers(int64) ([]int64, error)
	GetOtherMembers(int64, int64) ([]int64, error)
	IsMember(int64, int64) (bool, error)
	AddConversation(int64, int64) (sql.Result, error)
	DeleteConversation(int64) (sql.Result, error)
	DeleteUserConversation(int64, int64) (sql.Result, error)
	CreateConversation(int64, int64) (*int64, error)
	GetConversation(int64, int64) (*int64, error)
	GetConversations(int64) ([]int64, error)

	// Messages.
	InsertMessage(string, int64, int64, bool, *int64, string, MediaType) (*Message, error)
	DeleteMessage(int64) (sql.Result, error)
	GetMessage(int64) (*Message, error)
	GetMessageIds(int64) ([]int64, error)

	// Status.
	SetReceiptStatus(int64, int64, Status) (sql.Result, error)
	GetReceipt(int64, int64) (*Receipt, error)
	GetReceipts(int64) ([]Receipt, error)

	// Reactions.
	GetReactions(int64) ([]Reaction, error)
	GetReaction(int64, int64) (*Reaction, error)
	AddReaction(Emoji, int64, int64) (sql.Result, error)
	DeleteReaction(int64, int64) (sql.Result, error)
	UpdateReaction(Emoji, int64, int64) (sql.Result, error)
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
        user_id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        photo_url TEXT
      );`,

			`CREATE TABLE groups (
        conversation_id INTEGER PRIMARY KEY,
        founder_id INTEGER NOT NULL,
				name TEXT NOT NULL,
        photo_url TEXT,
				created_at TIMESTAMP,
				FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id) ON DELETE CASCADE,
				FOREIGN KEY (founder_id) REFERENCES users(user_id) ON DELETE CASCADE
      );`,

			`CREATE TABLE conversations (
				conversation_id INTEGER PRIMARY KEY
			);`,

			`CREATE TABLE user_conversations (
				conversation_id INTEGER NOT NULL,
			  user_id INTEGER NOT NULL,
				PRIMARY KEY (conversation_id, user_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
				FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id) ON DELETE CASCADE
      );`,

			`CREATE TABLE messages (
        message_id INTEGER PRIMARY KEY,
        text TEXT,
        sender_id INTEGER NOT NULL,
				conversation_id INTEGER NOT NULL,
				created_at TIMESTAMP,
				is_forwarded BOOL DEFAULT 0,
				comment_to INTEGER,
				attachment_url TEXT,
				media_type TEXT,
				FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id) ON DELETE CASCADE,
				FOREIGN KEY (sender_id) REFERENCES users(user_id),
				FOREIGN KEY (comment_to) REFERENCES messages(message_id)
      );`,

			`CREATE TABLE receipts (
				message_id INTEGER NOT NULL,
				user_id INTEGER NOT NULL,
        status TEXT NOT NULL CHECK (status  IN ('sent', 'received', 'read')),
				sent_at TIMESTAMP,
				received_at TIMESTAMP,
				read_at TIMESTAMP,
				PRIMARY KEY (message_id, user_id),
        FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE CASCADE
				FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
      );`,

			`CREATE TABLE reactions (
				emoji TEXT NOT NULL,
				message_id INTEGER NOT NULL,
				sender_id INTEGER NOT NULL,
				PRIMARY KEY (message_id, sender_id),
        FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE CASCADE,
				FOREIGN KEY (sender_id) REFERENCES users(user_id) ON DELETE CASCADE
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
