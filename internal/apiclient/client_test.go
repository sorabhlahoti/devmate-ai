package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNote(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/notes" {
			t.Fatalf("expected path /api/v1/notes, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_ = json.NewEncoder(w).Encode(NoteResponse{
			ID:        99,
			Title:     "Docker fix",
			Body:      "Check permission",
			CreatedAt: "2026-05-23T00:00:00Z",
		})
	}))
	defer testServer.Close()

	client, err := New(testServer.URL)
	if err != nil {
		t.Fatalf("expected client to be created, got error: %v", err)
	}

	note, err := client.CreateNote(context.Background(), "Docker fix", "Check permission")
	if err != nil {
		t.Fatalf("expected note to be created, got error: %v", err)
	}

	if note.ID != 99 {
		t.Fatalf("expected remote id 99, got %d", note.ID)
	}
}

func TestNewRejectsInvalidURL(t *testing.T) {
	_, err := New("localhost:8080")
	if err == nil {
		t.Fatal("expected error for URL without scheme")
	}
}
