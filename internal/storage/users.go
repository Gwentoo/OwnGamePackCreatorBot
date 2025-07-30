package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) SaveUser(userID int64) error {
	_, err := s.db.Exec(
		`INSERT INTO users (user_id) 
        VALUES ($1) 
        ON CONFLICT (user_id) DO NOTHING`,
		userID,
	)
	return err
}

func (s *Storage) CheckUserExists(userID int64) (bool, error) {
	var exists bool
	err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`,
		userID,
	).Scan(&exists)

	return exists, err
}

func (s *Storage) GetUserPackCount(userID int64) (int, error) {
	var value int
	query := `SELECT pack_count FROM users WHERE user_id = $1`
	err := s.db.QueryRow(query, userID).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, errors.New("user not found")
		}
		return -1, fmt.Errorf("failed to get field: %w", err)
	}
	return value, nil
}
