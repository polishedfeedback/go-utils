package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	GetPassword(name string) (Cred, error)
	AddPassword(name, url, email, password string) error
	ListPasswords() ([]Cred, error)
	DeletePassword(name string) error
	Exists(name string) bool
}

type Cred struct {
	ID        int
	Name      string
	URL       string
	Email     string
	Password  string
	CreatedAt time.Time
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS passwords(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME
		)`)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) GetPassword(name string) (*Cred, error) {
	var cred Cred
	err := s.db.QueryRow("SELECT id, name, url, email, password, created_at FROM passwords WHERE name = ?", name).Scan(&cred.ID, &cred.Name, &cred.URL, &cred.Email, &cred.Password, &cred.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("password not found")
	} else if err != nil {
		return nil, fmt.Errorf("error querying password: %v", err)
	}
	return &cred, nil
}

func (s *SQLiteStorage) Exists(name string) bool {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM passwords WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (s *SQLiteStorage) AddPassword(name, url, email, password string) error {
	if s.Exists(name) {
		return fmt.Errorf("name already exists")
	}
	_, err := s.db.Exec("INSERT INTO passwords (name, url, email, password, created_at) VALUES(?, ?, ?, ?, ?)", name, url, email, password, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("error adding password: %v", err)
	}
	return nil
}

func (s *SQLiteStorage) ListPasswords() ([]Cred, error) {
	var creds []Cred
	rows, err := s.db.Query("SELECT id, name, url, email, password, created_at FROM passwords")
	if err != nil {
		return nil, fmt.Errorf("error listing passwords: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cred Cred
		err := rows.Scan(&cred.ID, &cred.Name, &cred.URL, &cred.Email, &cred.Password, &cred.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error getting rows from list of passwords: %v", err)
		}
		creds = append(creds, cred)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting passwords: %v", err)
	}

	return creds, nil
}

func (s *SQLiteStorage) DeletePassword(name string) error {
	if !s.Exists(name) {
		return fmt.Errorf("password doesn't exist")
	}
	_, err := s.db.Exec("DELETE FROM passwords WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("error deleting password: %v", err)
	}
	return nil
}
