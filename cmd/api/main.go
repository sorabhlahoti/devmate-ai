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
