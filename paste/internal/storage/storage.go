package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	Save(url, content string, expiresIn time.Duration) error
	Get(url string) (string, error)
	Exists(url string) bool
	Delete(url string) error
	List() []string
}

// For maps

// type MemoryStorage struct {
// 	pastes map[string]string
// }

// func NewMemoryStorage() *MemoryStorage {
// 	return &MemoryStorage{
// 		pastes: make(map[string]string),
// 	}
// }

// func (m *MemoryStorage) Save(url, content string) error {
// 	if _, exists := m.pastes[url]; exists {
// 		return fmt.Errorf("URL already exists")
// 	}
// 	m.pastes[url] = content
// 	return nil
// }

// func (m *MemoryStorage) Get(url string) (string, error) {
// 	val, exists := m.pastes[url]
// 	if !exists {
// 		return "", fmt.Errorf("URL not found")
// 	}
// 	return val, nil
// }

// func (m *MemoryStorage) Exists(url string) bool {
// 	_, exists := m.pastes[url]
// 	return exists
// }

// func (m *MemoryStorage) Delete(url string) error {
// 	if _, exists := m.pastes[url]; exists {
// 		delete(m.pastes, url)
// 		return nil
// 	}
// 	return fmt.Errorf("URL not found")
// }

// func (m *MemoryStorage) List() []string {
// 	var urls []string
// 	for key := range m.pastes {
// 		urls = append(urls, key)
// 	}
// 	return urls
// }

// SQL
type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS pastes(
			url TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			expires_at DATETIME
		)`)
	if err != nil {
		return nil, fmt.Errorf("Error creating a paste: %v", err)
	}
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Save(url, content string, expiresIn time.Duration) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM pastes WHERE url = ?", url).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking URL: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("URL already exists")
	}

	var expiresAt interface{}
	if expiresIn > 0 {
		expiresAt = time.Now().Add(expiresIn).Format(time.RFC3339)
	} else {
		expiresAt = nil
	}
	_, err = s.db.Exec(
		"INSERT INTO pastes (url, content, created_at, expires_at) VALUES (?, ?, ?, ?)",
		url, content, time.Now().Format(time.RFC3339), expiresAt)
	if err != nil {
		return fmt.Errorf("error saving paste: %v", err)
	}
	return nil
}

func (s *SQLiteStorage) Get(url string) (string, error) {
	var content string
	var expiresAt sql.NullString

	err := s.db.QueryRow("SELECT content, expires_at FROM pastes WHERE url = ?", url).Scan(&content, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("URL cannot be found")
		}
		return "", fmt.Errorf("error querying the database: %v", err)
	}

	if expiresAt.Valid {
		expiry, err := time.Parse(time.RFC3339, expiresAt.String)
		if err != nil {
			return "", fmt.Errorf("error parsing the expiry: %v", err)
		}

		if time.Now().After(expiry) {
			s.Delete(url)
			return "", fmt.Errorf("paste has expired")
		}
	}
	return content, nil
}

func (s *SQLiteStorage) Exists(url string) bool {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM pastes WHERE url = ?", url).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (s *SQLiteStorage) Delete(url string) error {
	res, err := s.db.Exec("DELETE FROM pastes WHERE url = ?", url)
	if err != nil {
		return fmt.Errorf("error deleting paste: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("URL not found")
	}
	return nil
}

func (s *SQLiteStorage) List() []string {
	var urls []string
	rows, err := s.db.Query("SELECT * FROM pastes")
	if err != nil {
		return urls
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			continue
		}
		urls = append(urls, url)
	}
	return urls
}
