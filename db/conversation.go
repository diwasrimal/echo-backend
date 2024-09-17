package db

import (
	"context"
	"time"

	"github.com/diwasrimal/echo-backend/models"
)

func GetConversationsOf(userId uint64) ([]models.Conversation, error) {
	var conversations []models.Conversation
	rows, err := pool.Query(
		context.Background(),
		"SELECT * FROM conversations WHERE "+
			"user1_id = $1 OR user2_id = $1 "+
			"ORDER BY timestamp DESC",
		userId,
	)
	if err != nil {
		return conversations, err
	}
	defer rows.Close()
	for rows.Next() {
		var conv models.Conversation
		if err := rows.Scan(&conv.UserId1, &conv.UserId2, &conv.Timestamp); err != nil {
			return conversations, err
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil // TODO: maybe add limit
}

// Updates an exsiting conversation's timestamp between two users
// or creates a new one
func UpdateOrCreateConversation(senderId, receiverId uint64, timestamp time.Time) error {
	_, err := pool.Exec(
		context.Background(),
		`INSERT INTO conversations(user1_id, user2_id, timestamp) VALUES
			($1, $2, $3) ON CONFLICT(LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id)) 
			DO UPDATE 
			SET timestamp = excluded.timestamp`,
		senderId,
		receiverId,
		timestamp,
	)
	return err
}
