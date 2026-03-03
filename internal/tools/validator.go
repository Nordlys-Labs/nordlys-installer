package tools

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

// Schema URL constants for JSON schema validation
const (
	OpenCodeSchemaURL  = "https://opencode.ai/config.json"
	GeminiCLISchemaURL = "https://raw.githubusercontent.com/google-gemini/gemini-cli/main/schemas/settings.schema.json"
	CodexSchemaURL     = "https://developers.openai.com/codex/config-schema.json"
)

// SchemaCache stores compiled JSON schemas with a TTL
type SchemaCache struct {
	schemas map[string]*gojsonschema.Schema
	mu      sync.RWMutex
	ttl     time.Duration
}

// globalSchemaCache is a singleton cache for compiled schemas
var globalSchemaCache = &SchemaCache{
	schemas: make(map[string]*gojsonschema.Schema),
	ttl:     24 * time.Hour,
}

// Get retrieves a compiled schema from cache or compiles it
func (c *SchemaCache) Get(url string) (*gojsonschema.Schema, error) {
	c.mu.RLock()
	if schema, ok := c.schemas[url]; ok {
		c.mu.RUnlock()
		return schema, nil
	}
	c.mu.RUnlock()

	// Fetch and compile schema
	loader := gojsonschema.NewReferenceLoader(url)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.schemas[url] = schema
	c.mu.Unlock()

	return schema, nil
}

// ValidateConfig validates a config against a schema URL
// Returns nil on success, logs warning and returns nil on failure (warn and continue)
func ValidateConfig(schemaURL string, config interface{}) error {
	if schemaURL == "" {
		return nil // No schema to validate against
	}

	schema, err := globalSchemaCache.Get(schemaURL)
	if err != nil {
		log.Printf("Warning: failed to load schema %s: %v", schemaURL, err)
		return nil // Warn and continue
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Printf("Warning: failed to marshal config for validation: %v", err)
		return nil
	}

	documentLoader := gojsonschema.NewBytesLoader(configJSON)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		log.Printf("Warning: schema validation error: %v", err)
		return nil
	}

	if !result.Valid() {
		log.Println("Warning: config validation issues:")
		for _, e := range result.Errors() {
			log.Printf("  - %s", e)
		}
	}

	return nil
}

// structToMap converts a struct to map[string]any for merging with existing configs
func structToMap(v interface{}) map[string]any {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}
