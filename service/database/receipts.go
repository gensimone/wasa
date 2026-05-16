package database

import (
	"database/sql"
	"time"

	"github.com/gensimone/WASA-project/service/globaltime"
)

// Status the status of the message associated with the specified message and user id.
func (db *appdbimpl) SetReceiptStatus(messageId, userId int64, status Status) (sql.Result, error) {
	receipt, err := db.GetReceipt(messageId, userId)
	if err != nil {
		return nil, err
	}

	switch {
	case receipt.Status == Read:
	case receipt.Status == status:
	case receipt.Status == Received && status != Read:
		return nil, nil
	}

	timestamp := globaltime.Now().Format(time.DateTime)

	if status == Received {
		return db.c.Exec(
			`UPDATE receipts
			SET
			status = ?,
			received_at = ?
			WHERE message_id = ? AND user_id = ?`,
			Received,
			timestamp,
			messageId,
			userId,
		)
	}

	return db.c.Exec(
		`UPDATE receipts
		SET
		status = ?,
		read_at = ?
		WHERE message_id = ? AND user_id = ?`,
		Read,
		timestamp,
		messageId,
		userId,
	)
}

// Returns the status of the message associated with the specified message and user id.
func (db *appdbimpl) GetReceipt(messageId, userId int64) (*Receipt, error) {
	var receipt Receipt
	if err := db.c.QueryRow(
		`SELECT
		message_id,
		user_id,
		status,
		sent_at,
		received_at,
		read_at
		FROM receipts
		WHERE message_id = ? AND user_id = ?`,
		messageId,
		userId,
	).Scan(
		&receipt.MessageId,
		&receipt.UserId,
		&receipt.Status,
		&receipt.SentAt,
		&receipt.ReceivedAt,
		&receipt.ReadAt,
	); err != nil {
		return nil, err
	}

	return &receipt, nil
}

// Returns the various status associated with the specified message id.
func (db *appdbimpl) GetReceipts(messageId int64) ([]Receipt, error) {
	rows, err := db.c.Query(
		`SELECT
		message_id,
		user_id,
		status,
		sent_at,
		received_at,
		read_at
		FROM receipts
		WHERE message_id = ?`,
		messageId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var receipts []Receipt
	for rows.Next() {
		var receipt Receipt
		if err := rows.Scan(
			&receipt.MessageId,
			&receipt.UserId,
			&receipt.Status,
			&receipt.SentAt,
			&receipt.ReceivedAt,
			&receipt.ReadAt,
		); err != nil {
			return nil, err
		}

		receipts = append(receipts, receipt)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return receipts, nil
}
