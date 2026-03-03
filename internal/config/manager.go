package config

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
)

func ReadJSONFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]any), nil
		}
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func WriteJSONFile(path string, data map[string]any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, content, 0o600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpFile, path); err != nil {
		if rmErr := os.Remove(tmpFile); rmErr != nil {
			return fmt.Errorf("rename failed: %w; cleanup failed: %v", err, rmErr)
		}
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func UpdateJSONFields(path string, updates map[string]any) error {
	existing, err := ReadJSONFile(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		if err := CreateBackup(path); err != nil {
			return err
		}
	}

	merged := DeepMerge(existing, updates)
	return WriteJSONFile(path, merged)
}

func DeepMerge(base, overlay map[string]any) map[string]any {
	result := make(map[string]any)

	maps.Copy(result, base)

	for k, v := range overlay {
		if baseMap, ok := base[k].(map[string]any); ok {
			if overlayMap, ok := v.(map[string]any); ok {
				result[k] = DeepMerge(baseMap, overlayMap)
				continue
			}
		}
		result[k] = v
	}

	return result
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
