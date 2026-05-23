package apiclient

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
}

type CreateNoteRequest struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}

type NoteResponse struct {
    ID        int64  `json:"id"`
    Title     string `json:"title"`
    Body      string `json:"body"`
    CreatedAt string `json:"created_at"`
}

func New(baseURL string) (*Client, error) {
    cleanBaseURL := strings.TrimRight(strings.TrimSpace(baseURL), "/")

    if cleanBaseURL == "" {
        return nil, fmt.Errorf("api url cannot be empty")
    }

    parsedURL, err := url.Parse(cleanBaseURL)
    if err != nil {
        return nil, fmt.Errorf("invalid api url: %w", err)
    }

    if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
        return nil, fmt.Errorf("api url must start with http:// or https://")
    }

    if parsedURL.Host == "" {
        return nil, fmt.Errorf("api url host cannot be empty")
    }

    return &Client{
        baseURL: cleanBaseURL,
        httpClient: &http.Client{
            Timeout: 15 * time.Second,
        },
    }, nil
}

func (c *Client) CreateNote(ctx context.Context, title string, body string) (NoteResponse, error) {
    requestBody := CreateNoteRequest{
        Title: title,
        Body:  body,
    }

    jsonBody, err := json.Marshal(requestBody)
    if err != nil {
        return NoteResponse{}, fmt.Errorf("failed to encode note request: %w", err)
    }

    endpoint := c.baseURL + "/api/v1/notes"

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
    if err != nil {
        return NoteResponse{}, fmt.Errorf("failed to create note request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return NoteResponse{}, fmt.Errorf("failed to call api: %w", err)
    }
    defer resp.Body.Close()

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return NoteResponse{}, fmt.Errorf("failed to read api response: %w", err)
    }

    if resp.StatusCode != http.StatusCreated {
        return NoteResponse{}, fmt.Errorf("api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
    }

    var note NoteResponse

    if err := json.Unmarshal(responseBody, &note); err != nil {
        return NoteResponse{}, fmt.Errorf("failed to decode note response: %w", err)
    }

    return note, nil
}
