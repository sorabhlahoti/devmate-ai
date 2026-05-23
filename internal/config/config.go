package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	APIURL string `json:"api_url"`
}

func DefaultConfig() Config {
	return Config{
		APIURL: "http://localhost:8080",
	}
}

func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to find user config directory: %w", err)
	}

	return filepath.Join(configDir, "devmate-ai", "config.json"), nil
}

func Load(path string) (Config, error) {
	if strings.TrimSpace(path) == "" {
		defaultPath, err := DefaultPath()
		if err != nil {
			return Config{}, err
		}

		path = defaultPath
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}

		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config

	if err := json.Unmarshal(content, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	if strings.TrimSpace(cfg.APIURL) == "" {
		cfg.APIURL = DefaultConfig().APIURL
	}

	return cfg, nil
}

func Save(path string, cfg Config) error {
	if strings.TrimSpace(path) == "" {
		defaultPath, err := DefaultPath()
		if err != nil {
			return err
		}

		path = defaultPath
	}

	configDir := filepath.Dir(path)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	if err := os.WriteFile(path, content, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func SetValue(path string, key string, value string) error {
	cleanKey := strings.TrimSpace(key)
	cleanValue := strings.TrimSpace(value)

	if cleanKey == "" {
		return fmt.Errorf("config key cannot be empty")
	}

	if cleanValue == "" {
		return fmt.Errorf("config value cannot be empty")
	}

	cfg, err := Load(path)
	if err != nil {
		return err
	}

	switch cleanKey {
	case "api_url":
		cfg.APIURL = cleanValue
	default:
		return fmt.Errorf("unsupported config key: %s", cleanKey)
	}

	return Save(path, cfg)
}

func GetValue(path string, key string) (string, error) {
	cleanKey := strings.TrimSpace(key)

	if cleanKey == "" {
		return "", fmt.Errorf("config key cannot be empty")
	}

	cfg, err := Load(path)
	if err != nil {
		return "", err
	}

	switch cleanKey {
	case "api_url":
		return cfg.APIURL, nil
	default:
		return "", fmt.Errorf("unsupported config key: %s", cleanKey)
	}
}
