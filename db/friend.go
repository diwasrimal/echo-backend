package db

import (
	"context"

	"github.com/diwasrimal/echo-backend/models"
)

func RecordFriendRequest(from, to uint64) error {
	_, err := pool.Exec(
		context.Background(),
		`INSERT INTO friend_requests(requestor_id, receiver_id)
			VALUES($1, $2)
			ON CONFLICT DO NOTHING`,
		from,
		to,
	)
	return err
}

func DeleteFriendRequest(userId1, userId2 uint64) error {
	_, err := pool.Exec(
		context.Background(),
		`DELETE FROM friend_requests WHERE
			requestor_id = $1 AND receiver_id = $2 OR
			requestor_id = $2 AND receiver_id = $1`,
		userId1,
		userId2,
	)
	return err
}

func RecordFriendship(userId1, userId2 uint64) error {
	_, err := pool.Exec(
		context.Background(),
		"INSERT INTO friends(user1_id, user2_id) VALUES($1, $2)",
		userId1,
		userId2,
	)
	return err
}

func DeleteFriendship(userId1, userId2 uint64) error {
	_, err := pool.Exec(
		context.Background(),
		`DELETE FROM friends WHERE
			user1_id = $1 AND user2_id = $2 OR
			user1_id = $2 AND user2_id = $1`,
		userId1,
		userId2,
	)
	return err
}

// Returns status of friendship for two users from first user's point of view.
// Can give 4 statuses, "friends", "req-sent", "req-received", "unknown".
// Ex. "req-sent" means first user has sent a request to second.
func GetFriendshipStatus(userId, otherUserId uint64) (string, error) {
	var status string
	if err := pool.QueryRow(
		context.Background(),
		`SELECT CASE
			WHEN EXISTS (
				SELECT 1 FROM friends WHERE
				(user1_id = $1 AND user2_id = $2) OR
				(user2_id = $1 AND user1_id = $2) ) THEN 'friends'
			WHEN EXISTS (
				SELECT 1 FROM friend_requests WHERE requestor_id = $1 AND receiver_id = $2
			) THEN 'req-sent'
			WHEN EXISTS (
				SELECT 1 FROM friend_requests WHERE receiver_id = $1 AND requestor_id = $2
			) THEN 'req-received'
			ELSE 'unknown'
		END AS status`,
		userId,
		otherUserId,
	).Scan(&status); err != nil {
		return "", err
	}
	return status, nil
}

// Returns list of users that are friends to user with given id
func GetFriends(userId uint64) ([]models.User, error) {
	var friends []models.User
	rows, err := pool.Query(
		context.Background(),
		`SELECT * FROM users WHERE id IN (
			SELECT CASE WHEN user1_id = $1 THEN user2_id ELSE user1_id END
			FROM friends WHERE
			user1_id = $1 OR user2_id = $1
		)`,
		userId,
	)
	if err != nil {
		return friends, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
			return friends, err
		}
		friends = append(friends, user)
	}
	return friends, nil // TODO: maybe add limit
}

// Returns users that sent friend request to user with provided id
func GetFriendRequestSenders(toUserId uint64) ([]models.User, error) {
	var requestors []models.User
	rows, err := pool.Query(
		context.Background(),
		`SELECT * FROM users WHERE id IN (
			SELECT requestor_id FROM friend_requests WHERE
			receiver_id = $1
		)`,
		toUserId,
	)
	if err != nil {
		return requestors, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
			return requestors, err
		}
		requestors = append(requestors, user)
	}
	return requestors, nil // TODO: maybe add limit
}

// Returns users that user with provided id has sent requests to
func GetFriendRequestReceivers(fromUserId uint64) ([]models.User, error) {
	var receivers []models.User
	rows, err := pool.Query(
		context.Background(),
		`SELECT * FROM users WHERE id IN (
			SELECT receiver_id FROM friend_requests WHERE
			requestor_id = $1
		)`,
		fromUserId,
	)
	if err != nil {
		return receivers, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
			return receivers, err
		}
		receivers = append(receivers, user)
	}
	return receivers, nil // TODO: maybe add limit
}

func GetSentFriendRequests(requestorId uint64) ([]models.FriendRequest, error) {
	var reqs []models.FriendRequest
	rows, err := pool.Query(
		context.Background(),
		`SELECT * FROM friend_requests WHERE requestor_id = $1`,
		requestorId,
	)
	if err != nil {
		return reqs, err
	}
	defer rows.Close()
	for rows.Next() {
		var req models.FriendRequest
		if err := rows.Scan(&req.RequestorId, &req.ReceiverId); err != nil {
			return reqs, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil // TODO: maybe add limit
}

func GetReceivedFriendRequests(receiverId uint64) ([]models.FriendRequest, error) {
	var reqs []models.FriendRequest
	rows, err := pool.Query(
		context.Background(),
		`SELECT * FROM friend_requests WHERE receiver_id = $1`,
		receiverId,
	)
	if err != nil {
		return reqs, err
	}
	defer rows.Close()
	for rows.Next() {
		var req models.FriendRequest
		if err := rows.Scan(&req.RequestorId, &req.ReceiverId); err != nil {
			return reqs, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil // TODO: maybe add limit
}
