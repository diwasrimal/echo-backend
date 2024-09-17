package db

import (
	"context"

	"github.com/diwasrimal/echo-backend/models"
	"github.com/jackc/pgx/v5"
)

func CreateUserSession(userId uint64, sessionId string) error {
	_, err := pool.Exec(
		context.Background(),
		"INSERT INTO user_sessions(user_id, session_id) "+
			"VALUES($1, $2) "+
			"ON CONFLICT(user_id) DO UPDATE "+
			"SET session_id = excluded.session_id",
		userId,
		sessionId,
	)
	return err
}

func DeleteUserSession(sessionId string) error {
	_, err := pool.Exec(
		context.Background(),
		"DELETE FROM user_sessions WHERE session_id = $1",
		sessionId,
	)
	return err
}

func GetSession(sessionId string) (*models.Session, error) {
	var session models.Session
	if err := pool.QueryRow(
		context.Background(),
		"SELECT * FROM user_sessions WHERE session_id = $1",
		sessionId,
	).Scan(&session.UserId, &session.SessionId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}
