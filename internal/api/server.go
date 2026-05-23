package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "github.com/sorabhlahoti/devmate-ai/internal/storage"

    httpSwagger "github.com/swaggo/http-swagger/v2"

    _ "github.com/sorabhlahoti/devmate-ai/internal/swaggerdocs"
)

type Server struct {
    store *storage.Store
    mux   *http.ServeMux
}

type CreateNoteRequest struct {
    Title string `json:"title" example:"Docker permission fix"`
    Body  string `json:"body" example:"Check mounted volume permission"`
}

type NoteResponse struct {
    ID        int64  `json:"id" example:"1"`
    Title     string `json:"title" example:"Docker permission fix"`
    Body      string `json:"body" example:"Check mounted volume permission"`
    CreatedAt string `json:"created_at" example:"2026-05-23T10:30:00Z"`
}

type NotesResponse struct {
    Notes []NoteResponse `json:"notes"`
}

type ErrorResponse struct {
    Error string `json:"error" example:"title cannot be empty"`
}

func NewServer(store *storage.Store) *Server {
    server := &Server{
        store: store,
        mux:   http.NewServeMux(),
    }

    server.registerRoutes()

    return server
}

func (s *Server) Handler() http.Handler {
    return s.mux
}

func (s *Server) registerRoutes() {
    s.mux.HandleFunc("/health", s.handleHealth)
    s.mux.HandleFunc("/api/v1/notes", s.handleNotes)
    s.mux.HandleFunc("/api/v1/notes/search", s.handleSearchNotes)
    s.mux.HandleFunc("/swagger/", httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"),
    ))
}

// handleHealth godoc
// @Summary Health check
// @Description Check whether DevMate API is running.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 405 {object} ErrorResponse
// @Router /health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        respondError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{
        "status":  "ok",
        "service": "devmate-api",
    })
}

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        s.handleCreateNote(w, r)
    case http.MethodGet:
        s.handleListNotes(w, r)
    default:
        respondError(w, http.StatusMethodNotAllowed, "method not allowed")
    }
}

// handleCreateNote godoc
// @Summary Create note
// @Description Create a developer note in the backend database.
// @Tags notes
// @Accept json
// @Produce json
// @Param request body CreateNoteRequest true "Create note request"
// @Success 201 {object} NoteResponse
// @Failure 400 {object} ErrorResponse
// @Failure 405 {object} ErrorResponse
// @Router /api/v1/notes [post]
func (s *Server) handleCreateNote(w http.ResponseWriter, r *http.Request) {
    var req CreateNoteRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid json body")
        return
    }

    note, err := s.store.SaveNote(r.Context(), req.Title, req.Body)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusCreated, toNoteResponse(note))
}

// handleListNotes godoc
// @Summary List notes
// @Description List saved developer notes from the backend database.
// @Tags notes
// @Produce json
// @Param limit query int false "Maximum number of notes to return"
// @Success 200 {object} NotesResponse
// @Failure 405 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notes [get]
func (s *Server) handleListNotes(w http.ResponseWriter, r *http.Request) {
    limit := parseLimit(r, 20)

    notes, err := s.store.ListNotes(r.Context(), limit)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, NotesResponse{
        Notes: toNoteResponses(notes),
    })
}

// handleSearchNotes godoc
// @Summary Search notes
// @Description Search developer notes by title or body.
// @Tags notes
// @Produce json
// @Param q query string true "Search text"
// @Param limit query int false "Maximum number of notes to return"
// @Success 200 {object} NotesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 405 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notes/search [get]
func (s *Server) handleSearchNotes(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        respondError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    query := strings.TrimSpace(r.URL.Query().Get("q"))
    if query == "" {
        respondError(w, http.StatusBadRequest, "query parameter q is required")
        return
    }

    limit := parseLimit(r, 20)

    notes, err := s.store.SearchNotes(r.Context(), query, limit)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, NotesResponse{
        Notes: toNoteResponses(notes),
    })
}

func parseLimit(r *http.Request, defaultLimit int) int {
    rawLimit := strings.TrimSpace(r.URL.Query().Get("limit"))
    if rawLimit == "" {
        return defaultLimit
    }

    limit, err := strconv.Atoi(rawLimit)
    if err != nil || limit <= 0 {
        return defaultLimit
    }

    if limit > 100 {
        return 100
    }

    return limit
}

func toNoteResponse(note storage.Note) NoteResponse {
    return NoteResponse{
        ID:        note.ID,
        Title:     note.Title,
        Body:      note.Body,
        CreatedAt: note.CreatedAt.Format("2006-01-02T15:04:05Z"),
    }
}

func toNoteResponses(notes []storage.Note) []NoteResponse {
    responses := make([]NoteResponse, 0, len(notes))

    for _, note := range notes {
        responses = append(responses, toNoteResponse(note))
    }

    return responses
}

func respondJSON(w http.ResponseWriter, statusCode int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        fmt.Println("failed to write json response:", err)
    }
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
    respondJSON(w, statusCode, ErrorResponse{
        Error: message,
    })
}
