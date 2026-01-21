package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	AddMovie(title string, year int, actors []string) error
	GetMovie(id int) (*Movie, error)
	ListMovies() ([]Movie, error)
}

type Movie struct {
	ID        int
	Title     string
	Year      int
	CreatedAt time.Time
	Actors    []string
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	queries := []string{
		`CREATE TABLE IF NOT EXISTS movies(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			year INTEGER,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS actors(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS movie_actors(
			movie_id INTEGER,
			actor_id INTEGER,
			FOREIGN KEY (movie_id) REFERENCES movies(id),
			FOREIGN KEY (actor_id) REFERENCES actors(id),
			PRIMARY KEY (movie_id, actor_id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %v", err)
		}
	}
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) AddMovie(title string, year int, actors []string) error {
	result, err := s.db.Exec("INSERT INTO movies (title, year, created_at) VALUES (?, ?, ?)", title, year, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("error adding movie: %v", err)
	}

	movieId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting movie ID: %v", err)
	}

	for _, actorName := range actors {
		var actorID int64

		err = s.db.QueryRow("SELECT id FROM actors WHERE name = ?", actorName).Scan(&actorID)
		if err == sql.ErrNoRows {
			result, err := s.db.Exec("INSERT INTO actors(name) VALUES (?)", actorName)
			if err != nil {
				return fmt.Errorf("error adding actor: %v", err)
			}
			actorID, _ = result.LastInsertId()
		} else if err != nil {
			return fmt.Errorf("error adding actor: %v", err)
		}
		_, err := s.db.Exec("INSERT INTO movie_actors(movie_id, actor_id) VALUES (?, ?)", movieId, actorID)
		if err != nil {
			return fmt.Errorf("error linking movies and actors: %v", err)
		}
	}
	return nil
}

func (s *SQLiteStorage) GetMovie(id int) (*Movie, error) {
	var movie Movie
	err := s.db.QueryRow("SELECT id, title, year, created_at FROM movies WHERE id = ?", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("movie not found: %v", err)
	} else if err != nil {
		return nil, fmt.Errorf("error getting movie: %v", err)
	}

	rows, err := s.db.Query(
		`SELECT actors.name
		 FROM actors
		 JOIN movie_actors ON actors.id = movie_actors.actor_id
		 WHERE movie_actors.movie_id = ?
		`, id)
	if err != nil {
		return nil, fmt.Errorf("error getting actors: %v", err)
	}
	defer rows.Close()

	var actors []string
	for rows.Next() {
		var actorName string
		rows.Scan(&actorName)
		actors = append(actors, actorName)
	}
	movie.Actors = actors
	return &movie, nil
}

func (s *SQLiteStorage) ListMovies() ([]Movie, error) {
	var movies []Movie
	rows, err := s.db.Query("SELECT id, title, year, created_at FROM movies")
	if err != nil {
		return nil, fmt.Errorf("error getting movies: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.CreatedAt); err != nil {
			return nil, fmt.Errorf("error getting movies: %v", err)
		}

		var actors []string
		actorRows, err := s.db.Query(
			`SELECT actors.name
			FROM actors
			JOIN movie_actors ON actors.id = movie_actors.actor_id 
			WHERE movie_actors.movie_id = ?`, movie.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting actors: %v", err)
		}
		defer actorRows.Close()
		for actorRows.Next() {
			var actor string
			if err := actorRows.Scan(&actor); err != nil {
				return nil, fmt.Errorf("error scanning actor: %v", err)
			}
			actors = append(actors, actor)
		}
		movie.Actors = actors
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting movies: %v", err)
	}

	return movies, nil
}
