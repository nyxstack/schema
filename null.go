package schema

import (
	"encoding/json"

	"github.com/nyxstack/i18n"
)

// Default error messages for null validation
var (
	nullRequiredError = i18n.S("value is required")
	nullTypeError     = i18n.S("value must be null")
)

// NullSchema represents a JSON Schema for null values
type NullSchema struct {
	Schema
	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	typeMismatchError ErrorMessage
}

// Null creates a new null schema with optional type error message
func Null(errorMessage ...interface{}) *NullSchema {
	schema := &NullSchema{
		Schema: Schema{
			schemaType: "null",
			required:   true, // Default to required
		},
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *NullSchema) Title(title string) *NullSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *NullSchema) Description(description string) *NullSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value (always nil for null schemas)
func (s *NullSchema) Default(value interface{}) *NullSchema {
	if value == nil {
		s.Schema.defaultValue = nil
	}
	// Ignore non-nil defaults for null schemas
	return s
}

// Example adds an example value (always nil for null schemas)
func (s *NullSchema) Example(example interface{}) *NullSchema {
	if example == nil {
		s.Schema.examples = append(s.Schema.examples, nil)
	}
	// Ignore non-nil examples for null schemas
	return s
}

// Required/Optional control

// Optional marks the schema as optional
func (s *NullSchema) Optional() *NullSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *NullSchema) Required(errorMessage ...interface{}) *NullSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *NullSchema) TypeError(message string) *NullSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *NullSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *NullSchema) IsOptional() bool {
	return !s.Schema.required
}

// Validation

// Parse validates and parses a null value, returning the final parsed value
func (s *NullSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	// Handle nil values
	if value == nil {
		// Null values are always valid for null schemas
		return ParseResult{Valid: true, Value: nil, Errors: nil}
	}

	// Handle missing required field
	if s.Schema.required {
		// Check if we have a default value (should be nil)
		if defaultVal := s.GetDefault(); defaultVal == nil {
			// Use default value and re-parse it
			return s.Parse(defaultVal, ctx)
		}
		// Required null field but got non-nil value
		message := nullRequiredError(ctx.Locale)
		if !isEmptyErrorMessage(s.requiredError) {
			message = resolveErrorMessage(s.requiredError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "required")},
		}
	}

	// Non-nil value for null schema is type error
	message := nullTypeError(ctx.Locale)
	if !isEmptyErrorMessage(s.typeMismatchError) {
		message = resolveErrorMessage(s.typeMismatchError, ctx)
	}
	return ParseResult{
		Valid:  false,
		Value:  nil,
		Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
	}
}

// JSON generates JSON Schema representation
func (s *NullSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("null")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	// Default and examples should always be null for null schemas
	if s.GetDefault() == nil {
		schema["default"] = nil
	}
	if len(s.GetExamples()) > 0 {
		schema["examples"] = []interface{}{nil}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize NullSchema for JSON schema generation
func (s *NullSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Schema)
}
