package config

import (
	"path/filepath"
	"testing"
)

func TestLoadReturnsDefaultWhenFileDoesNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.APIURL != "http://localhost:8080" {
		t.Fatalf("expected default api url, got %s", cfg.APIURL)
	}
}

func TestSetAndGetAPIURL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	err := SetValue(path, "api_url", "http://example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	value, err := GetValue(path, "api_url")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if value != "http://example.com" {
		t.Fatalf("expected http://example.com, got %s", value)
	}
}

func TestUnsupportedKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	err := SetValue(path, "unknown", "value")
	if err == nil {
		t.Fatal("expected error for unsupported key")
	}
}
