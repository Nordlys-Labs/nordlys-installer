package runtime

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	MinNodeVersion = 18
)

func GetNodeVersion() (int, error) {
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	version := strings.TrimSpace(string(output))
	version = strings.TrimPrefix(version, "v")

	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version format")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	return major, nil
}

func IsNodeInstalled() bool {
	_, err := exec.LookPath("node")
	return err == nil
}

func EnsureNodeJS() error {
	if !IsNodeInstalled() {
		return fmt.Errorf("node.js is not installed, install v%d+ from https://nodejs.org/", MinNodeVersion)
	}

	version, err := GetNodeVersion()
	if err != nil {
		return fmt.Errorf("failed to get Node.js version: %w", err)
	}

	if version < MinNodeVersion {
		return fmt.Errorf("node.js v%d found, v%d+ required, upgrade from https://nodejs.org/", version, MinNodeVersion)
	}

	return nil
}
