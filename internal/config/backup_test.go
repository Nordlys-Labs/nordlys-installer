package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateBackup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setup   func(string) error
		name    string
		wantErr bool
	}{
		{
			name: "backup existing file",
			setup: func(path string) error {
				return os.WriteFile(path, []byte("original content"), 0o600)
			},
			wantErr: false,
		},
		{
			name: "backup non-existent file",
			setup: func(path string) error {
				return nil
			},
			wantErr: false,
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

			err := CreateBackup(testFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && FileExists(testFile) {
				entries, err := os.ReadDir(tmpDir)
				if err != nil {
					t.Fatalf("ReadDir() error = %v", err)
				}

				foundBackup := false
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".bak") {
						foundBackup = true
						break
					}
				}

				if !foundBackup {
					t.Error("CreateBackup() should create .bak file")
				}
			}
		})
	}
}

func TestGetLatestBackup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setup      func(string) error
		name       string
		wantErr    bool
		wantLatest bool
	}{
		{
			name: "single backup",
			setup: func(dir string) error {
				testFile := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(testFile, []byte("content"), 0o600); err != nil {
					return err
				}
				return CreateBackup(testFile)
			},
			wantErr:    false,
			wantLatest: true,
		},
		{
			name: "multiple backups",
			setup: func(dir string) error {
				testFile := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(testFile, []byte("content1"), 0o600); err != nil {
					return err
				}
				if err := CreateBackup(testFile); err != nil {
					return err
				}
				time.Sleep(10 * time.Millisecond)
				if err := os.WriteFile(testFile, []byte("content2"), 0o600); err != nil {
					return err
				}
				return CreateBackup(testFile)
			},
			wantErr:    false,
			wantLatest: true,
		},
		{
			name: "no backups",
			setup: func(dir string) error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.txt")

			if err := tt.setup(tmpDir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got, err := GetLatestBackup(testFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == "" {
					t.Error("GetLatestBackup() returned empty path")
				}
				if !strings.HasSuffix(got, ".bak") {
					t.Errorf("GetLatestBackup() = %q, should end with .bak", got)
				}
				if !FileExists(got) {
					t.Errorf("GetLatestBackup() returned non-existent file: %q", got)
				}
			}
		})
	}
}

func TestRestoreFromBackup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setup   func(string) (string, error)
		name    string
		wantErr bool
	}{
		{
			name: "restore from backup",
			setup: func(dir string) (string, error) {
				testFile := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(testFile, []byte("original"), 0o600); err != nil {
					return "", err
				}
				if err := CreateBackup(testFile); err != nil {
					return "", err
				}
				if err := os.WriteFile(testFile, []byte("modified"), 0o600); err != nil {
					return "", err
				}
				return testFile, nil
			},
			wantErr: false,
		},
		{
			name: "no backup to restore",
			setup: func(dir string) (string, error) {
				testFile := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(testFile, []byte("content"), 0o600); err != nil {
					return "", err
				}
				return testFile, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			testFile, err := tt.setup(tmpDir)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err = RestoreFromBackup(testFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestoreFromBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				content, err := os.ReadFile(testFile)
				if err != nil {
					t.Fatalf("ReadFile() after restore error = %v", err)
				}
				if string(content) != "original" {
					t.Errorf("RestoreFromBackup() content = %q, want %q", string(content), "original")
				}
			}
		})
	}
}
