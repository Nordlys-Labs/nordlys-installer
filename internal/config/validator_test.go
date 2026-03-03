package config

import (
	"testing"
)

func TestValidateAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "sk-1234567890abcdefghij",
			wantErr: false,
		},
		{
			name:    "valid with underscores",
			apiKey:  "key_with_underscores_1234567890",
			wantErr: false,
		},
		{
			name:    "valid with dashes",
			apiKey:  "key-with-dashes-1234567890",
			wantErr: false,
		},
		{
			name:    "valid with dots",
			apiKey:  "key.with.dots.1234567890",
			wantErr: false,
		},
		{
			name:    "empty key",
			apiKey:  "",
			wantErr: true,
		},
		{
			name:    "too short",
			apiKey:  "short",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			apiKey:  "key-with-invalid-chars!@#$%",
			wantErr: true,
		},
		{
			name:    "with spaces",
			apiKey:  "key with spaces 1234567890",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAPIKey(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		model   string
		wantErr bool
	}{
		{
			name:    "valid model",
			model:   "nordlys/hypernova",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			model:   "author123/model456",
			wantErr: false,
		},
		{
			name:    "valid with dashes",
			model:   "author-name/model-name",
			wantErr: false,
		},
		{
			name:    "valid with underscores",
			model:   "author_name/model_name",
			wantErr: false,
		},
		{
			name:    "valid with dots",
			model:   "author.name/model.name",
			wantErr: false,
		},
		{
			name:    "empty model",
			model:   "",
			wantErr: false,
		},
		{
			name:    "missing slash",
			model:   "authormodel",
			wantErr: true,
		},
		{
			name:    "no author",
			model:   "/model",
			wantErr: true,
		},
		{
			name:    "no model",
			model:   "author/",
			wantErr: true,
		},
		{
			name:    "too many slashes",
			model:   "author/model/extra",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			model:   "author@name/model#name",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateModel(tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAPIConnection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "invalid API key will fail",
			apiKey:  "invalid-test-key-1234567890",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAPIConnection(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAPIConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
