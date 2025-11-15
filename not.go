package schema

import (
	"github.com/nyxstack/i18n"
)

// Default error messages for not validation
var (
	notShouldNotMatchError = i18n.S("value should not match the specified schema")
)

// NotErrors defines error message functions
var NotErrors = struct {
	ShouldNotMatch i18n.TranslatedFunc
}{
	ShouldNotMatch: notShouldNotMatchError,
}

// NotSchema represents a "not" validation schema that rejects values matching the given schema
type NotSchema struct {
	schema   Parseable
	notError ErrorMessage
}

// Not creates a new Not schema that rejects values matching the given schema
func Not(schema Parseable) *NotSchema {
	return &NotSchema{
		schema: schema,
	}
}

// NotError sets a custom error message for when the value matches (and should not)
func (s *NotSchema) NotError(err ErrorMessage) *NotSchema {
	s.notError = err
	return s
}

// Parse validates that a value does NOT match the specified schema
func (s *NotSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	// Try to parse with the inner schema
	result := s.schema.Parse(value, ctx)

	// If the inner schema validation succeeded, this should fail
	if result.Valid {
		message := NotErrors.ShouldNotMatch(ctx.Locale)
		if !isEmptyErrorMessage(s.notError) {
			message = resolveErrorMessage(s.notError, ctx)
		}

		return ParseResult{
			Valid:  false,
			Value:  value,
			Errors: []ValidationError{NewPrimitiveError(value, message, "not_match")},
		}
	}

	// If the inner schema validation failed, this succeeds
	return ParseResult{
		Valid:  true,
		Value:  value,
		Errors: nil,
	}
}

// JSON generates JSON Schema for Not validation
func (s *NotSchema) JSON() map[string]interface{} {
	if jsonSchema, ok := s.schema.(interface{ JSON() map[string]interface{} }); ok {
		return map[string]interface{}{
			"not": jsonSchema.JSON(),
		}
	}

	// Fallback if schema doesn't support JSON generation
	return map[string]interface{}{
		"not": map[string]interface{}{"type": "unknown"},
	}
}
