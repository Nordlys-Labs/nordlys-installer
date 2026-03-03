package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// ReadTOMLFile reads and parses a TOML file into a map
func ReadTOMLFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]any), nil
		}
		return nil, err
	}

	var result map[string]any
	if err := toml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}
	return result, nil
}

// WriteTOMLFile writes a map to a TOML file
func WriteTOMLFile(path string, data map[string]any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	content, err := toml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal TOML: %w", err)
	}

	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, content, 0o644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpFile, path); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// UpdateTOMLFields reads a TOML file, merges updates, and writes back
func UpdateTOMLFields(path string, updates map[string]any) error {
	existing, err := ReadTOMLFile(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		if err := CreateBackup(path); err != nil {
			return err
		}
	}

	merged := DeepMerge(existing, updates)
	return WriteTOMLFile(path, merged)
}
