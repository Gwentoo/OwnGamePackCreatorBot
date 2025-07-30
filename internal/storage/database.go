package storage

import (
	"OwnGamePack/config"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.DBHost, cfg.DB.DBPort, cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}
	fmt.Println(cfg.DB)
	if err1 := db.Ping(); err1 != nil {
		return nil, fmt.Errorf("failed to ping storage: %w", err1)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
