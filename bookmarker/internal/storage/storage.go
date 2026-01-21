package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Storage interface {
	AddBookmark(url, title, description string, tags []string) error
	GetBookmark(id int) (*Bookmark, error)
	ListBookmarks() ([]Bookmark, error)
	ListByTag(tag string) ([]Bookmark, error)
	Search(query string) ([]Bookmark, error)
	DeleteBookmark(id int) error
}

type Bookmark struct {
	ID          int
	URL         string
	Title       string
	Description string
	CreatedAt   time.Time
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}
	queries := []string{
		`CREATE TABLE IF NOT EXISTS bookmarks(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			title TEXT,
			description TEXT,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tags(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS bookmark_tags(
			bookmark_id INTEGER,
			tag_id INTEGER,
			FOREIGN KEY(tag_id) REFERENCES tags(id),
			PRIMARY KEY (bookmark_id, tag_id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("Error creating table: %v", err)
		}
	}
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) AddBookmark(url, title, description string, tags []string) error {
	result, err := s.db.Exec(
		`INSERT INTO bookmarks (url, title, description, created_at) VALUES (?, ?, ?, ?)`, url, title, description, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("error inserting bookmark: %v", err)
	}

	bookmarkId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting bookmark ID: %v", err)
	}

	for _, tagName := range tags {
		var tagID int64

		err := s.db.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			result, err := s.db.Exec("INSERT INTO tags(name) VALUES (?)", tagName)
			if err != nil {
				return fmt.Errorf("error creating tag: %v", err)
			}
			tagID, _ = result.LastInsertId()
		} else if err != nil {
			return fmt.Errorf("error checking tag: %v", err)
		}
		_, err = s.db.Exec("INSERT INTO bookmark_tags(bookmark_id, tag_id) VALUES (?, ?)", bookmarkId, tagID)
		if err != nil {
			return fmt.Errorf("error linking bookmark to tag: %v", err)
		}
	}
	return nil
}

func (s *SQLiteStorage) GetBookmark(id int) (*Bookmark, error) {
	var bookmark Bookmark

	err := s.db.QueryRow("SELECT * FROM bookmarks WHERE id = ?", id).Scan(&bookmark)
	if err == sql.ErrNoRows {
		return &Bookmark{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("error getting bookmark: %v", err)
	}

	return &Bookmark{
		ID:          bookmark.ID,
		URL:         bookmark.URL,
		Title:       bookmark.Title,
		Description: bookmark.Description,
		CreatedAt:   bookmark.CreatedAt,
	}, nil
}
