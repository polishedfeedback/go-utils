package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	Save(shortCode, longUrl string) error
	Get(shortCode string) (string, error)
	Exists(shortCode string) bool
	IncrementClicks(shortCode string) error
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS urls(
			short_code TEXT PRIMARY KEY,
			long_url TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			clicks INTEGER DEFAULT 0
		)`)
	if err != nil {
		return nil, fmt.Errorf("Error creating the table: %v", err)
	}
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Save(shortCode, longUrl string) error {

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = ?", shortCode).Scan(&count)
	if err != nil {
		return fmt.Errorf("Error checking the url: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("Short Code already exists: %v", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO urls (short_code, long_url, created_at) VALUES (?, ?, ?)", shortCode, longUrl, time.Now().Format(time.RFC3339))

	if err != nil {
		return fmt.Errorf("Error saving URL: %v", err)
	}
	return nil
}

func (s *SQLiteStorage) Get(shortCode string) (string, error) {
	var longUrl string
	err := s.db.QueryRow("SELECT long_url FROM urls WHERE short_code = ?", shortCode).Scan(&longUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("No rows found: %v", err)
		}
		return "", fmt.Errorf("Error getting the URL: %v", err)
	}
	return longUrl, nil
}

func (s *SQLiteStorage) Exists(shortCode string) bool {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = ?", shortCode).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (s *SQLiteStorage) IncrementClicks(shortCode string) error {
	_, err := s.db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?", shortCode)
	if err != nil {
		return fmt.Errorf("error incrementing clicks: %v", err)
	}
	return nil
}
