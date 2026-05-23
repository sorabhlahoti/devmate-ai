package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/sorabhlahoti/devmate-ai/internal/api"
    "github.com/sorabhlahoti/devmate-ai/internal/storage"
)

// @title DevMate AI API
// @version 0.1.0
// @description DevMate AI backend API for syncing developer notes from the CLI.
// @termsOfService http://swagger.io/terms/

// @contact.name DevMate AI Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
    port := getEnv("PORT", "8080")
    dbPath := strings.TrimSpace(os.Getenv("DEVMATE_DB_PATH"))

    store, err := storage.NewSQLiteStore(dbPath)
    if err != nil {
        log.Fatal("failed to create storage:", err)
    }
    defer store.Close()

    server := api.NewServer(store)

    addr := ":" + port

    fmt.Println("DevMate API is running on", addr)
    fmt.Println("Swagger UI is available at /swagger/index.html")

    if err := http.ListenAndServe(addr, server.Handler()); err != nil {
        log.Fatal("server failed:", err)
    }
}

func getEnv(key string, fallback string) string {
    value := strings.TrimSpace(os.Getenv(key))
    if value == "" {
        return fallback
    }

    return value
}
