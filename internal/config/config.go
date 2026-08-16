package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const DefaultConfigFileName = ".envsync.json"
const DefaultServerURL = "http://localhost:8080"

type ProjectConfig struct {
	ProjectID   string    `json:"project_id"`
	ProjectName string    `json:"project_name"`
	ServerURL   string    `json:"server_url"`
	DefaultEnv  string    `json:"default_env"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetConfigPath returns the path to .envsync.json
func GetConfigPath(customPath string) string {
	if customPath != "" {
		return customPath
	}
	return DefaultConfigFileName
}

// ConfigExists checks if .envsync.json exists
func ConfigExists(customPath string) bool {
	path := GetConfigPath(customPath)
	_, err := os.Stat(path)
	return err == nil
}

// LoadProjectConfig reads and parses .envsync.json
func LoadProjectConfig(customPath string) (*ProjectConfig, error) {
	path := GetConfigPath(customPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("project not initialized. Run 'envsync init' first")
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg ProjectConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file format %s: %w", path, err)
	}

	return &cfg, nil
}

// SaveProjectConfig writes ProjectConfig to .envsync.json
func SaveProjectConfig(cfg *ProjectConfig, customPath string) error {
	path := GetConfigPath(customPath)

	cfg.UpdatedAt = time.Now()
	if cfg.CreatedAt.IsZero() {
		cfg.CreatedAt = cfg.UpdatedAt
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	return nil
}
