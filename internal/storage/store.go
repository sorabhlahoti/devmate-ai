package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Note struct {
	ID        int64
	Title     string
	Body      string
	CreatedAt time.Time
}

type Store struct {
	db *sql.DB
}

func DefaultDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to find user config directory: %w", err)
	}

	appDir := filepath.Join(configDir, "devmate-ai")

	return filepath.Join(appDir, "devmate.db"), nil
}

func NewSQLiteStore(dbPath string) (*Store, error) {
	if strings.TrimSpace(dbPath) == "" {
		defaultPath, err := DefaultDBPath()
		if err != nil {
			return nil, err
		}

		dbPath = defaultPath
	}

	dbDir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect sqlite database: %w", err)
	}

	store := &Store{db: db}

	if err := store.Init(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Init(ctx context.Context) error {
	query := `
CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TEXT NOT NULL
);`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to initialize notes table: %w", err)
	}

	return nil
}

func (s *Store) SaveNote(ctx context.Context, title string, body string) (Note, error) {
	cleanTitle := strings.TrimSpace(title)
	cleanBody := strings.TrimSpace(body)

	if cleanTitle == "" {
		return Note{}, fmt.Errorf("title cannot be empty")
	}

	if cleanBody == "" {
		return Note{}, fmt.Errorf("body cannot be empty")
	}

	createdAt := time.Now().UTC().Truncate(time.Second)

	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO notes (title, body, created_at) VALUES (?, ?, ?)`,
		cleanTitle,
		cleanBody,
		createdAt.Format(time.RFC3339),
	)
	if err != nil {
		return Note{}, fmt.Errorf("failed to save note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Note{}, fmt.Errorf("failed to get saved note id: %w", err)
	}

	return Note{
		ID:        id,
		Title:     cleanTitle,
		Body:      cleanBody,
		CreatedAt: createdAt,
	}, nil
}

func (s *Store) ListNotes(ctx context.Context, limit int) ([]Note, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, title, body, created_at FROM notes ORDER BY id DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}
	defer rows.Close()

	return scanNotes(rows)
}

func (s *Store) SearchNotes(ctx context.Context, searchText string, limit int) ([]Note, error) {
	cleanSearchText := strings.TrimSpace(searchText)

	if cleanSearchText == "" {
		return nil, fmt.Errorf("search text cannot be empty")
	}

	if limit <= 0 {
		limit = 20
	}

	likeValue := "%" + strings.ToLower(cleanSearchText) + "%"

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, title, body, created_at
         FROM notes
         WHERE lower(title) LIKE ? OR lower(body) LIKE ?
         ORDER BY id DESC
         LIMIT ?`,
		likeValue,
		likeValue,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}
	defer rows.Close()

	return scanNotes(rows)
}

func scanNotes(rows *sql.Rows) ([]Note, error) {
	var notes []Note

	for rows.Next() {
		var note Note
		var createdAtText string

		if err := rows.Scan(&note.ID, &note.Title, &note.Body, &createdAtText); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtText)
		if err != nil {
			return nil, fmt.Errorf("failed to parse note created_at: %w", err)
		}

		note.CreatedAt = createdAt
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading notes: %w", err)
	}

	return notes, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
