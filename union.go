package schema

import (
	"encoding/json"
	"fmt"

	"github.com/nyxstack/i18n"
)

// Default error messages for union validation
var (
	unionRequiredError      = i18n.S("value is required")
	unionNoMatchError       = i18n.S("value does not match any of the allowed schemas")
	unionMultipleMatchError = i18n.S("value matches multiple schemas, only one is allowed")
)

// UnionSchema represents a JSON Schema oneOf for union types
type UnionSchema struct {
	Schema
	schemas   []Parseable // The schemas to validate against
	nullable  bool        // Allow null values
	allowNone bool        // Allow values that match none of the schemas

	// Error messages for validation failures (support i18n)
	requiredError      ErrorMessage
	noMatchError       ErrorMessage
	multipleMatchError ErrorMessage
	typeMismatchError  ErrorMessage
}

// Union creates a new union schema with the provided schemas
func Union(schemas ...Parseable) *UnionSchema {
	schema := &UnionSchema{
		Schema: Schema{
			schemaType: "oneOf",
			required:   true, // Default to required
		},
		schemas: schemas,
	}
	return schema
}

// OneOf is an alias for Union for JSON Schema compatibility
func OneOf(schemas ...Parseable) *UnionSchema {
	return Union(schemas...)
}

// Core fluent API methods

// Title sets the title of the schema
func (s *UnionSchema) Title(title string) *UnionSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *UnionSchema) Description(description string) *UnionSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *UnionSchema) Default(value interface{}) *UnionSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *UnionSchema) Example(example interface{}) *UnionSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Schema manipulation

// Add appends additional schemas to the union
func (s *UnionSchema) Add(schemas ...Parseable) *UnionSchema {
	s.schemas = append(s.schemas, schemas...)
	return s
}

// Schemas returns all schemas in the union
func (s *UnionSchema) Schemas() []Parseable {
	return s.schemas
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *UnionSchema) Optional() *UnionSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *UnionSchema) Required(errorMessage ...interface{}) *UnionSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *UnionSchema) Nullable() *UnionSchema {
	s.nullable = true
	return s
}

// AllowNone allows values that don't match any schema (makes union more permissive)
func (s *UnionSchema) AllowNone() *UnionSchema {
	s.allowNone = true
	return s
}

// Error customization

// NoMatchError sets a custom error message when no schemas match
func (s *UnionSchema) NoMatchError(message string) *UnionSchema {
	s.noMatchError = toErrorMessage(message)
	return s
}

// MultipleMatchError sets a custom error message when multiple schemas match
func (s *UnionSchema) MultipleMatchError(message string) *UnionSchema {
	s.multipleMatchError = toErrorMessage(message)
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *UnionSchema) TypeError(message string) *UnionSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *UnionSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *UnionSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *UnionSchema) IsNullable() bool {
	return s.nullable
}

// GetSchemaCount returns the number of schemas in the union
func (s *UnionSchema) GetSchemaCount() int {
	return len(s.schemas)
}

// Validation

// Parse validates and parses a union value, returning the final parsed value
func (s *UnionSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	// Handle nil values
	if value == nil {
		if s.nullable {
			// For nullable schemas, nil is a valid value
			return ParseResult{Valid: true, Value: nil, Errors: nil}
		}
		if s.Schema.required {
			// Check if we have a default value to use instead
			if defaultVal := s.GetDefault(); defaultVal != nil {
				// Use default value and re-parse it
				return s.Parse(defaultVal, ctx)
			}
			// No default, required field is missing
			message := unionRequiredError(ctx.Locale)
			if !isEmptyErrorMessage(s.requiredError) {
				message = resolveErrorMessage(s.requiredError, ctx)
			}
			return ParseResult{
				Valid:  false,
				Value:  nil,
				Errors: []ValidationError{NewPrimitiveError(value, message, "required")},
			}
		}
		// Optional field, use default if available
		if defaultVal := s.GetDefault(); defaultVal != nil {
			return s.Parse(defaultVal, ctx)
		}
		// Optional field with no default
		return ParseResult{Valid: true, Value: nil, Errors: nil}
	}

	// Validate against each schema in the union
	var validResults []ParseResult
	var allErrors []ValidationError

	for i, schema := range s.schemas {
		result := schema.Parse(value, ctx)
		if result.Valid {
			validResults = append(validResults, result)
		} else {
			// Collect errors from failed schemas for debugging
			for _, err := range result.Errors {
				// Add context about which schema failed
				contextualErr := ValidationError{
					Path:    append([]string{fmt.Sprintf("schema_%d", i)}, err.Path...),
					Value:   err.Value,
					Message: err.Message,
					Code:    err.Code,
				}
				allErrors = append(allErrors, contextualErr)
			}
		}
	}

	// Check validation results
	if len(validResults) == 0 {
		// No schemas matched
		if s.allowNone {
			// Allow values that don't match any schema
			return ParseResult{Valid: true, Value: value, Errors: nil}
		}
		message := unionNoMatchError(ctx.Locale)
		if !isEmptyErrorMessage(s.noMatchError) {
			message = resolveErrorMessage(s.noMatchError, ctx)
		}
		// Return the original value with no match error, plus all schema errors for context
		errors = append(errors, NewPrimitiveError(value, message, "no_match"))
		// Also include all the individual schema errors for debugging
		errors = append(errors, allErrors...)
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: errors,
		}
	}

	if len(validResults) > 1 {
		// Multiple schemas matched - this violates oneOf semantics
		message := unionMultipleMatchError(ctx.Locale)
		if !isEmptyErrorMessage(s.multipleMatchError) {
			message = resolveErrorMessage(s.multipleMatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "multiple_match")},
		}
	}

	// Exactly one schema matched - this is what we want
	return validResults[0]
}

// JSON generates JSON Schema representation
func (s *UnionSchema) JSON() map[string]interface{} {
	schema := make(map[string]interface{})

	// Generate oneOf array with all schemas
	oneOfSchemas := make([]interface{}, len(s.schemas))
	for i, subSchema := range s.schemas {
		if jsonSchema, ok := subSchema.(interface{ JSON() map[string]interface{} }); ok {
			oneOfSchemas[i] = jsonSchema.JSON()
		} else {
			// Fallback for schemas that don't implement JSON method
			oneOfSchemas[i] = map[string]interface{}{"type": "unknown"}
		}
	}
	schema["oneOf"] = oneOfSchemas

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())

	// Add nullable if true
	if s.nullable {
		// Add null to the oneOf array
		oneOfSchemas = append(oneOfSchemas, map[string]interface{}{"type": "null"})
		schema["oneOf"] = oneOfSchemas
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize UnionSchema for JSON schema generation
func (s *UnionSchema) MarshalJSON() ([]byte, error) {
	type jsonUnionSchema struct {
		Schema
		Schemas   []Parseable `json:"schemas"`
		Nullable  bool        `json:"nullable,omitempty"`
		AllowNone bool        `json:"allowNone,omitempty"`
	}

	return json.Marshal(jsonUnionSchema{
		Schema:    s.Schema,
		Schemas:   s.schemas,
		Nullable:  s.nullable,
		AllowNone: s.allowNone,
	})
}
