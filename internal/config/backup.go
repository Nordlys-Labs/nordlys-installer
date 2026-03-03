package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func CreateBackup(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.bak", path, timestamp)

	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy to backup: %w", err)
	}

	return nil
}

func GetLatestBackup(path string) (string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, base+".") && strings.HasSuffix(name, ".bak") {
			backups = append(backups, filepath.Join(dir, name))
		}
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("no backups found")
	}

	sort.Strings(backups)
	return backups[len(backups)-1], nil
}

func RestoreFromBackup(path string) error {
	backupPath, err := GetLatestBackup(path)
	if err != nil {
		return err
	}

	src, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
