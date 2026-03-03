package constants

import (
	"testing"
)

func TestConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		got     string
		wantNot string
	}{
		{"Version is set", Version, ""},
		{"APIBaseURL is set", APIBaseURL, ""},
		{"APIKeyURL is set", APIKeyURL, ""},
		{"DefaultModel is set", DefaultModel, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got == tt.wantNot {
				t.Errorf("constant should not be empty")
			}
		})
	}
}

func TestAPITimeout(t *testing.T) {
	t.Parallel()

	if APITimeout <= 0 {
		t.Errorf("APITimeout should be positive, got %d", APITimeout)
	}
}

func TestDefaultModelFormat(t *testing.T) {
	t.Parallel()

	if DefaultModel != "nordlys/hypernova" {
		t.Errorf("DefaultModel = %q, want %q", DefaultModel, "nordlys/hypernova")
	}
}

func TestAPIBaseURLFormat(t *testing.T) {
	t.Parallel()

	if APIBaseURL != "https://api.nordlyslabs.com" {
		t.Errorf("APIBaseURL = %q, want %q", APIBaseURL, "https://api.nordlyslabs.com")
	}
}
