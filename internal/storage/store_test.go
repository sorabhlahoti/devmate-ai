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

func TestListUnsyncedAndMarkSynced(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("expected store to open, got error: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	savedNote, err := store.SaveNote(ctx, "Go channel blocking", "Send blocks when no receiver is ready")
	if err != nil {
		t.Fatalf("expected note to save, got error: %v", err)
	}

	unsyncedNotes, err := store.ListUnsyncedNotes(ctx, 10)
	if err != nil {
		t.Fatalf("expected unsynced notes to list, got error: %v", err)
	}

	if len(unsyncedNotes) != 1 {
		t.Fatalf("expected 1 unsynced note, got %d", len(unsyncedNotes))
	}

	err = store.MarkNoteSynced(ctx, savedNote.ID, 101)
	if err != nil {
		t.Fatalf("expected note to be marked synced, got error: %v", err)
	}

	unsyncedNotes, err = store.ListUnsyncedNotes(ctx, 10)
	if err != nil {
		t.Fatalf("expected unsynced notes to list, got error: %v", err)
	}

	if len(unsyncedNotes) != 0 {
		t.Fatalf("expected 0 unsynced notes, got %d", len(unsyncedNotes))
	}

	notes, err := store.ListNotes(ctx, 10)
	if err != nil {
		t.Fatalf("expected notes to list, got error: %v", err)
	}

	if notes[0].RemoteID == nil {
		t.Fatal("expected remote id to be set")
	}

	if *notes[0].RemoteID != 101 {
		t.Fatalf("expected remote id 101, got %d", *notes[0].RemoteID)
	}

	if notes[0].SyncedAt == nil {
		t.Fatal("expected synced at to be set")
	}
}
