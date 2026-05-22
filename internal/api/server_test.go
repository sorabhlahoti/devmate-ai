package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/sorabhlahoti/devmate-ai/internal/storage"
)

func TestHealthEndpoint(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestCreateListAndSearchNotes(t *testing.T) {
	server := newTestServer(t)

	createBody := []byte(`{"title":"Docker permission fix","body":"Check mounted volume permission"}`)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	server.Handler().ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", createRec.Code, createRec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/notes", nil)
	listRec := httptest.NewRecorder()

	server.Handler().ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", listRec.Code)
	}

	var listResp NotesResponse
	if err := json.NewDecoder(listRec.Body).Decode(&listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}

	if len(listResp.Notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(listResp.Notes))
	}

	searchReq := httptest.NewRequest(http.MethodGet, "/api/v1/notes/search?q=docker", nil)
	searchRec := httptest.NewRecorder()

	server.Handler().ServeHTTP(searchRec, searchReq)

	if searchRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", searchRec.Code)
	}

	var searchResp NotesResponse
	if err := json.NewDecoder(searchRec.Body).Decode(&searchResp); err != nil {
		t.Fatalf("failed to decode search response: %v", err)
	}

	if len(searchResp.Notes) != 1 {
		t.Fatalf("expected 1 search result, got %d", len(searchResp.Notes))
	}
}

func TestCreateNoteInvalidJSON(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewReader([]byte(`bad-json`)))
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func newTestServer(t *testing.T) *Server {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")

	store, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}

	t.Cleanup(func() {
		store.Close()
	})

	return NewServer(store)
}
