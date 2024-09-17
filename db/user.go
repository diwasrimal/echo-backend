package db

import (
	"context"
	"log"

	"github.com/diwasrimal/echo-backend/models"
	"github.com/jackc/pgx/v5"
)

func CreateUser(fullname, username, passwordHash string) error {
	_, err := pool.Exec(
		context.Background(),
		"INSERT INTO users(fullname, username, password_hash) VALUES($1, $2, $3)",
		fullname,
		username,
		passwordHash,
	)
	return err
}

func UpdateUser(userId uint64, newUser models.User) error {
	_, err := pool.Exec(
		context.Background(),
		"UPDATE users SET "+
			"fullname = $1, "+
			"bio = $2 "+
			"WHERE id = $3",
		newUser.Fullname,
		newUser.Bio,
		userId,
	)
	return err
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE username = $1",
		username,
	).Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserById(id uint64) (*models.User, error) {
	var user models.User
	if err := pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE id = $1",
		id,
	).Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserBySessionId(sessionId string) (*models.User, error) {
	var user models.User
	if err := pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE id = ( "+
			"SELECT user_id FROM user_sessions WHERE session_id = $1 "+
			")",
		sessionId,
	).Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func IsUsernameTaken(username string) (bool, error) {
	rows, err := pool.Query(
		context.Background(),
		"SELECT username FROM users WHERE username = $1",
		username,
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func GetRecentChatPartners(userId uint64) ([]models.User, error) {
	var partners []models.User
	rows, err := pool.Query(
		context.Background(),
		`SELECT u.* FROM users u JOIN (
			SELECT
				CASE WHEN user1_id = $1 THEN user2_id ELSE user1_id END
			AS id, timestamp
			FROM conversations WHERE
			user1_id = $1 OR user2_id = $1
		) as subq ON u.id = subq.id ORDER by subq.timestamp DESC`,
		userId,
	)
	if err != nil {
		return partners, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
			return partners, err
		}
		partners = append(partners, user)
	}
	return partners, nil // TODO: maybe add limit
}

func SearchUser(searchType, searchQuery string) ([]models.User, error) {
	var matches []models.User
	var rows pgx.Rows
	var err error

	// Fuzzy search using likeness and levenshtein distance.
	if searchType == "normal" {
		maxLevDist := int(0.5 * float64(len(searchQuery)))
		log.Println("maxLevList:", maxLevDist)
		rows, err = pool.Query(
			context.Background(),
			` SELECT * FROM users WHERE
				fullname ILIKE '%' || $1 || '%' OR
				levenshtein(fullname, $1) <= $2
				ORDER BY levenshtein(fullname, $1) ASC;
			`,
			searchQuery,
			maxLevDist,
		)
	} else if searchType == "by-username" {
		rows, err = pool.Query(
			context.Background(),
			"SELECT * FROM users WHERE username ILIKE '%' || $1 || '%'",
			searchQuery,
		)
	}
	if err != nil {
		return matches, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Fullname, &user.Username, &user.PasswordHash, &user.Bio); err != nil {
			return matches, err
		}
		matches = append(matches, user)
	}
	return matches, nil // TODO: maybe add limit
}
