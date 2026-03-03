package config

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

var (
	apiKeyPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{20,}$`)
	modelPattern  = regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)
)

func ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	if !apiKeyPattern.MatchString(apiKey) {
		return fmt.Errorf("API key format appears invalid")
	}
	return nil
}

func ValidateModel(model string) error {
	if model == "" {
		return nil
	}
	if !modelPattern.MatchString(model) {
		return fmt.Errorf("model format invalid, use: author/model_id (e.g., nordlys/hypernova)")
	}
	return nil
}

func ValidateAPIConnection(apiKey string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", constants.APIBaseURL+"/v1/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid API key")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}
