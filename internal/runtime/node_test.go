package runtime

import (
	"os/exec"
	"testing"
)

func TestIsNodeInstalled(t *testing.T) {
	t.Parallel()

	got := IsNodeInstalled()
	_, err := exec.LookPath("node")
	expected := err == nil

	if got != expected {
		t.Errorf("IsNodeInstalled() = %v, want %v", got, expected)
	}
}

func TestGetNodeVersion(t *testing.T) {
	t.Parallel()

	if !IsNodeInstalled() {
		t.Skip("Node.js not installed, skipping version test")
	}

	version, err := GetNodeVersion()
	if err != nil {
		t.Fatalf("GetNodeVersion() error = %v", err)
	}

	if version <= 0 {
		t.Errorf("GetNodeVersion() = %d, want positive number", version)
	}
}

func TestEnsureNodeJS(t *testing.T) {
	t.Parallel()

	if !IsNodeInstalled() {
		err := EnsureNodeJS()
		if err == nil {
			t.Error("EnsureNodeJS() should error when Node.js not installed")
		}
		return
	}

	version, err := GetNodeVersion()
	if err != nil {
		t.Fatalf("GetNodeVersion() error = %v", err)
	}

	err = EnsureNodeJS()
	if version >= MinNodeVersion {
		if err != nil {
			t.Errorf("EnsureNodeJS() error = %v, want nil for Node.js v%d", err, version)
		}
	} else {
		if err == nil {
			t.Errorf("EnsureNodeJS() should error for Node.js v%d (< v%d)", version, MinNodeVersion)
		}
	}
}
