package database

import "database/sql"

type Info string

const (
	Read        Info = "read"
	NotReceived Info = "not received"
	Received    Info = "received"
)

// Updates the status of the message associated with the specified message and user id.
func (db *appdbimpl) UpdateStatus(messageId int64, userId int64, info Info) (sql.Result, error) {
	return db.c.Exec(
		`UPDATE status SET info = ? WHERE message_id = ? AND user_id = ?`,
		info, messageId, userId,
	)
}

// Returns the status of the message associated with the specified message and user id.
func (db *appdbimpl) GetStatusOf(messageId int64, userId int64) (*Status, error) {
	var status Status
	if err := db.c.QueryRow(
		`SELECT message_id, user_id, info FROM status WHERE message_id = ? AND user_id = ?`,
		messageId, userId,
	).Scan(
		&status.MessageId,
		&status.UserId,
		&status.Info,
	); err != nil {
		return nil, err
	}

	return &status, nil
}

// Returns the various status associated with the specified message id.
func (db *appdbimpl) GetStatus(messageId int64) ([]Status, error) {
	rows, err := db.c.Query(
		`SELECT message_id, user_id, info FROM status WHERE message_id = ?`,
		messageId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var multipleStatus []Status
	for rows.Next() {
		var singleStatus Status
		if err := rows.Scan(
			&singleStatus.MessageId,
			&singleStatus.UserId,
			&singleStatus.Info,
		); err != nil {
			return nil, err
		}

		multipleStatus = append(multipleStatus, singleStatus)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return multipleStatus, nil
}
