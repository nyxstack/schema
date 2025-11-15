package schema

import (
	"encoding/json"
)

// JSONSchemaGenerator interface for types that can generate JSON Schema
type JSONSchemaGenerator interface {
	JSON() map[string]interface{}
}

// JSONSchema converts any schema to JSONSchema Schema format
func JSONSchema(s JSONSchemaGenerator) map[string]interface{} {
	return s.JSON()
}

// JSON converts any schema to JSON Schema bytes
func JSON(s JSONSchemaGenerator) ([]byte, error) {
	schema := s.JSON()
	return json.MarshalIndent(schema, "", "  ")
}

// Helper functions for common JSON Schema patterns

// baseJSONSchema creates a basic JSON Schema with type
func baseJSONSchema(schemaType string) map[string]interface{} {
	return map[string]interface{}{
		"type": schemaType,
	}
}

// addOptionalField adds a field to JSON Schema if value is not nil
func addOptionalField(schema map[string]interface{}, key string, value interface{}) {
	if value != nil {
		// Handle pointer types
		switch v := value.(type) {
		case *string:
			if v != nil {
				schema[key] = *v
			}
		case *int:
			if v != nil {
				schema[key] = *v
			}
		case *int64:
			if v != nil {
				schema[key] = *v
			}
		case *float64:
			if v != nil {
				schema[key] = *v
			}
		case *bool:
			if v != nil {
				schema[key] = *v
			}
		default:
			schema[key] = value
		}
	}
}

// addOptionalArray adds an array field to JSON Schema if slice is not empty
func addOptionalArray(schema map[string]interface{}, key string, value interface{}) {
	switch v := value.(type) {
	case []string:
		if len(v) > 0 {
			schema[key] = v
		}
	case []interface{}:
		if len(v) > 0 {
			schema[key] = v
		}
	}
}

// addTitle adds title if not empty
func addTitle(schema map[string]interface{}, title string) {
	if title != "" {
		schema["title"] = title
	}
}

// addDescription adds description if not empty
func addDescription(schema map[string]interface{}, description string) {
	if description != "" {
		schema["description"] = description
	}
}
