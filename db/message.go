package db

import (
	"context"
	"time"

	"github.com/diwasrimal/echo-backend/models"
)

func GetMessagesAmong(userId1, userId2 uint64) ([]models.Message, error) {
	var messages []models.Message
	rows, err := pool.Query(
		context.Background(),
		"SELECT * FROM messages WHERE "+
			"(sender_id = $1 AND receiver_id = $2) OR "+
			"(sender_id = $2 AND receiver_id = $1)"+
			"ORDER BY timestamp DESC",
		userId1,
		userId2,
	)
	if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.Id, &msg.SenderId, &msg.ReceiverId, &msg.Text, &msg.Timestamp); err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}
	return messages, nil // TODO: maybe add limit
}

func RecordMessage(senderId, receiverId uint64, text string, timestamp time.Time) (uint64, error) {
	var msgId uint64
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO messages(sender_id, receiver_id, text, timestamp)
			VALUES ($1, $2, $3, $4)
			RETURNING id`,
		senderId,
		receiverId,
		text,
		timestamp,
	).Scan(&msgId)
	return msgId, err
}
