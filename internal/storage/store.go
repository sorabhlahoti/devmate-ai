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
    RemoteID  *int64
    SyncedAt  *time.Time
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

    if err := s.ensureColumn(ctx, "notes", "remote_id", "INTEGER"); err != nil {
        return err
    }

    if err := s.ensureColumn(ctx, "notes", "synced_at", "TEXT"); err != nil {
        return err
    }

    return nil
}

func (s *Store) ensureColumn(ctx context.Context, table string, column string, columnType string) error {
    exists, err := s.columnExists(ctx, table, column)
    if err != nil {
        return err
    }

    if exists {
        return nil
    }

    query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, columnType)

    if _, err := s.db.ExecContext(ctx, query); err != nil {
        return fmt.Errorf("failed to add column %s: %w", column, err)
    }

    return nil
}

func (s *Store) columnExists(ctx context.Context, table string, column string) (bool, error) {
    rows, err := s.db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", table))
    if err != nil {
        return false, fmt.Errorf("failed to read table info: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var cid int
        var name string
        var dataType string
        var notNull int
        var defaultValue sql.NullString
        var primaryKey int

        if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
            return false, fmt.Errorf("failed to scan table info: %w", err)
        }

        if name == column {
            return true, nil
        }
    }

    if err := rows.Err(); err != nil {
        return false, fmt.Errorf("failed while reading table info: %w", err)
    }

    return false, nil
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
        `SELECT id, title, body, created_at, remote_id, synced_at
         FROM notes
         ORDER BY id DESC
         LIMIT ?`,
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
        `SELECT id, title, body, created_at, remote_id, synced_at
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

func (s *Store) ListUnsyncedNotes(ctx context.Context, limit int) ([]Note, error) {
    if limit <= 0 {
        limit = 20
    }

    rows, err := s.db.QueryContext(
        ctx,
        `SELECT id, title, body, created_at, remote_id, synced_at
         FROM notes
         WHERE remote_id IS NULL
         ORDER BY id ASC
         LIMIT ?`,
        limit,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to list unsynced notes: %w", err)
    }
    defer rows.Close()

    return scanNotes(rows)
}

func (s *Store) MarkNoteSynced(ctx context.Context, localID int64, remoteID int64) error {
    if localID <= 0 {
        return fmt.Errorf("local id must be positive")
    }

    if remoteID <= 0 {
        return fmt.Errorf("remote id must be positive")
    }

    syncedAt := time.Now().UTC().Truncate(time.Second)

    result, err := s.db.ExecContext(
        ctx,
        `UPDATE notes SET remote_id = ?, synced_at = ? WHERE id = ?`,
        remoteID,
        syncedAt.Format(time.RFC3339),
        localID,
    )
    if err != nil {
        return fmt.Errorf("failed to mark note as synced: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check synced note update: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("note not found: %d", localID)
    }

    return nil
}

func scanNotes(rows *sql.Rows) ([]Note, error) {
    var notes []Note

    for rows.Next() {
        var note Note
        var createdAtText string
        var remoteID sql.NullInt64
        var syncedAtText sql.NullString

        if err := rows.Scan(&note.ID, &note.Title, &note.Body, &createdAtText, &remoteID, &syncedAtText); err != nil {
            return nil, fmt.Errorf("failed to scan note: %w", err)
        }

        createdAt, err := time.Parse(time.RFC3339, createdAtText)
        if err != nil {
            return nil, fmt.Errorf("failed to parse note created_at: %w", err)
        }

        note.CreatedAt = createdAt

        if remoteID.Valid {
            value := remoteID.Int64
            note.RemoteID = &value
        }

        if syncedAtText.Valid && strings.TrimSpace(syncedAtText.String) != "" {
            syncedAt, err := time.Parse(time.RFC3339, syncedAtText.String)
            if err != nil {
                return nil, fmt.Errorf("failed to parse note synced_at: %w", err)
            }

            note.SyncedAt = &syncedAt
        }

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
