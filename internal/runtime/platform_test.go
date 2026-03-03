package runtime

import (
	"runtime"
	"testing"
)

func TestGetPlatform(t *testing.T) {
	t.Parallel()

	got := GetPlatform()
	if got == "" {
		t.Error("GetPlatform() returned empty string")
	}

	if got != runtime.GOOS {
		t.Errorf("GetPlatform() = %q, want %q", got, runtime.GOOS)
	}
}

func TestGetArch(t *testing.T) {
	t.Parallel()

	got := GetArch()
	if got == "" {
		t.Error("GetArch() returned empty string")
	}

	if got != runtime.GOARCH {
		t.Errorf("GetArch() = %q, want %q", got, runtime.GOARCH)
	}
}

func TestIsPosix(t *testing.T) {
	t.Parallel()

	got := IsPosix()
	expected := runtime.GOOS == "linux" || runtime.GOOS == "darwin"

	if got != expected {
		t.Errorf("IsPosix() = %v, want %v for OS %q", got, expected, runtime.GOOS)
	}
}

func TestIsWindows(t *testing.T) {
	t.Parallel()

	got := IsWindows()
	expected := runtime.GOOS == "windows"

	if got != expected {
		t.Errorf("IsWindows() = %v, want %v for OS %q", got, expected, runtime.GOOS)
	}
}
