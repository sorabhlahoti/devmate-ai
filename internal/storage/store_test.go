package storage

import (
    "context"
    "path/filepath"
    "testing"
)

func TestSaveListAndSearchNotes(t *testing.T) {
    dbPath := filepath.Join(t.TempDir(), "test.db")

    store, err := NewSQLiteStore(dbPath)
    if err != nil {
        t.Fatalf("expected store to open, got error: %v", err)
    }
    defer store.Close()

    ctx := context.Background()

    savedNote, err := store.SaveNote(ctx, "Docker permission fix", "Check mounted volume permission")
    if err != nil {
        t.Fatalf("expected note to save, got error: %v", err)
    }

    if savedNote.ID == 0 {
        t.Fatal("expected saved note to have an id")
    }

    notes, err := store.ListNotes(ctx, 10)
    if err != nil {
        t.Fatalf("expected notes to list, got error: %v", err)
    }

    if len(notes) != 1 {
        t.Fatalf("expected 1 note, got %d", len(notes))
    }

    results, err := store.SearchNotes(ctx, "docker", 10)
    if err != nil {
        t.Fatalf("expected notes to search, got error: %v", err)
    }

    if len(results) != 1 {
        t.Fatalf("expected 1 search result, got %d", len(results))
    }
}

func TestSaveNoteValidation(t *testing.T) {
    dbPath := filepath.Join(t.TempDir(), "test.db")

    store, err := NewSQLiteStore(dbPath)
    if err != nil {
        t.Fatalf("expected store to open, got error: %v", err)
    }
    defer store.Close()

    ctx := context.Background()

    _, err = store.SaveNote(ctx, "", "body")
    if err == nil {
        t.Fatal("expected error for empty title")
    }

    _, err = store.SaveNote(ctx, "title", "")
    if err == nil {
        t.Fatal("expected error for empty body")
    }
}
