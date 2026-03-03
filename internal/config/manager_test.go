package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadJSONFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setup     func(string) error
		name      string
		wantErr   bool
		wantEmpty bool
	}{
		{
			name: "existing valid JSON",
			setup: func(path string) error {
				return os.WriteFile(path, []byte(`{"key": "value"}`), 0o600)
			},
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name: "non-existent file",
			setup: func(path string) error {
				return nil
			},
			wantErr:   false,
			wantEmpty: true,
		},
		{
			name: "invalid JSON",
			setup: func(path string) error {
				return os.WriteFile(path, []byte(`{invalid}`), 0o600)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.json")

			if err := tt.setup(testFile); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got, err := ReadJSONFile(testFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantEmpty && len(got) != 0 {
					t.Errorf("ReadJSONFile() should return empty map for non-existent file")
				}
				if !tt.wantEmpty && len(got) == 0 {
					t.Errorf("ReadJSONFile() should return non-empty map")
				}
			}
		})
	}
}

func TestWriteJSONFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		data    map[string]any
		name    string
		wantErr bool
	}{
		{
			name:    "simple map",
			data:    map[string]any{"key": "value"},
			wantErr: false,
		},
		{
			name:    "nested map",
			data:    map[string]any{"outer": map[string]any{"inner": "value"}},
			wantErr: false,
		},
		{
			name:    "empty map",
			data:    map[string]any{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.json")

			err := WriteJSONFile(testFile, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !FileExists(testFile) {
					t.Error("WriteJSONFile() should create file")
				}

				got, err := ReadJSONFile(testFile)
				if err != nil {
					t.Fatalf("ReadJSONFile() after write error = %v", err)
				}

				if len(got) != len(tt.data) {
					t.Errorf("written data length = %d, want %d", len(got), len(tt.data))
				}
			}
		})
	}
}

func TestDeepMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     map[string]any
		overlay  map[string]any
		wantKeys []string
	}{
		{
			name:     "merge simple maps",
			base:     map[string]any{"a": 1, "b": 2},
			overlay:  map[string]any{"c": 3},
			wantKeys: []string{"a", "b", "c"},
		},
		{
			name:     "overlay overwrites",
			base:     map[string]any{"a": 1},
			overlay:  map[string]any{"a": 2},
			wantKeys: []string{"a"},
		},
		{
			name:     "merge nested maps",
			base:     map[string]any{"outer": map[string]any{"a": 1}},
			overlay:  map[string]any{"outer": map[string]any{"b": 2}},
			wantKeys: []string{"outer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DeepMerge(tt.base, tt.overlay)

			for _, key := range tt.wantKeys {
				if _, exists := got[key]; !exists {
					t.Errorf("DeepMerge() missing key %q", key)
				}
			}
		})
	}
}

func TestUpdateJSONFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		existing map[string]any
		updates  map[string]any
		name     string
		wantErr  bool
	}{
		{
			name:     "update existing file",
			existing: map[string]any{"keep": "this", "change": "old"},
			updates:  map[string]any{"change": "new", "add": "field"},
			wantErr:  false,
		},
		{
			name:     "create new file",
			existing: nil,
			updates:  map[string]any{"new": "field"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.json")

			if tt.existing != nil {
				if err := WriteJSONFile(testFile, tt.existing); err != nil {
					t.Fatalf("setup WriteJSONFile() error = %v", err)
				}
			}

			err := UpdateJSONFields(testFile, tt.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateJSONFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := ReadJSONFile(testFile)
				if err != nil {
					t.Fatalf("ReadJSONFile() after update error = %v", err)
				}

				for key, value := range tt.updates {
					if got[key] != value {
						t.Errorf("UpdateJSONFields() key %q = %v, want %v", key, got[key], value)
					}
				}

				if tt.existing != nil {
					for key := range tt.existing {
						if _, exists := got[key]; !exists && tt.updates[key] == nil {
							t.Errorf("UpdateJSONFields() lost existing key %q", key)
						}
					}
				}
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setup    func(string) error
		name     string
		wantTrue bool
	}{
		{
			name: "file exists",
			setup: func(path string) error {
				return os.WriteFile(path, []byte("test"), 0o600)
			},
			wantTrue: true,
		},
		{
			name: "file does not exist",
			setup: func(path string) error {
				return nil
			},
			wantTrue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.txt")

			if err := tt.setup(testFile); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := FileExists(testFile)
			if got != tt.wantTrue {
				t.Errorf("FileExists() = %v, want %v", got, tt.wantTrue)
			}
		})
	}
}
