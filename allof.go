package schema

import (
	"encoding/json"
	"fmt"

	"github.com/nyxstack/i18n"
)

// Default error messages for allof validation
var (
	allofRequiredError    = i18n.S("value is required")
	allofNotAllMatchError = i18n.S("value must match all provided schemas")
)

func allofSchemaError(index int) i18n.TranslatedFunc {
	return i18n.F("value failed to match schema %d", index)
}

// AllOfSchema represents a JSON Schema allOf for composition (value must match ALL schemas)
type AllOfSchema struct {
	Schema
	schemas  []Parseable // The schemas that ALL must match
	nullable bool        // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	notAllMatchError  ErrorMessage
	typeMismatchError ErrorMessage
}

// AllOf creates a new allof schema with the provided schemas (all must match)
func AllOf(schemas ...Parseable) *AllOfSchema {
	schema := &AllOfSchema{
		Schema: Schema{
			schemaType: "allOf",
			required:   true, // Default to required
		},
		schemas: schemas,
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *AllOfSchema) Title(title string) *AllOfSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *AllOfSchema) Description(description string) *AllOfSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *AllOfSchema) Default(value interface{}) *AllOfSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *AllOfSchema) Example(example interface{}) *AllOfSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Schema manipulation

// Add appends additional schemas to the allof (all must match)
func (s *AllOfSchema) Add(schemas ...Parseable) *AllOfSchema {
	s.schemas = append(s.schemas, schemas...)
	return s
}

// Schemas returns all schemas in the allof
func (s *AllOfSchema) Schemas() []Parseable {
	return s.schemas
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *AllOfSchema) Optional() *AllOfSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *AllOfSchema) Required(errorMessage ...interface{}) *AllOfSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *AllOfSchema) Nullable() *AllOfSchema {
	s.nullable = true
	return s
}

// Error customization

// NotAllMatchError sets a custom error message when not all schemas match
func (s *AllOfSchema) NotAllMatchError(message string) *AllOfSchema {
	s.notAllMatchError = toErrorMessage(message)
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *AllOfSchema) TypeError(message string) *AllOfSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *AllOfSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *AllOfSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *AllOfSchema) IsNullable() bool {
	return s.nullable
}

// GetSchemaCount returns the number of schemas in the allof
func (s *AllOfSchema) GetSchemaCount() int {
	return len(s.schemas)
}

// Validation

// Parse validates and parses an allof value, returning the final parsed value
func (s *AllOfSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := allofRequiredError(ctx.Locale)
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

	// Validate against ALL schemas in the allof
	var finalValue interface{} = value
	var allErrors []ValidationError

	for i, schema := range s.schemas {
		result := schema.Parse(value, ctx)
		if !result.Valid {
			// This schema failed - collect errors
			message := allofSchemaError(i)(ctx.Locale)
			errors = append(errors, NewPrimitiveError(value, message, "allof_schema_failed"))

			// Add context about which schema failed
			for _, err := range result.Errors {
				contextualErr := ValidationError{
					Path:    append([]string{fmt.Sprintf("allOf[%d]", i)}, err.Path...),
					Value:   err.Value,
					Message: err.Message,
					Code:    err.Code,
				}
				allErrors = append(allErrors, contextualErr)
			}
		} else {
			// This schema passed - use its parsed value
			// For allOf, we typically want the most "parsed" version of the value
			// If multiple schemas transform the value, the last successful one wins
			finalValue = result.Value
		}
	}

	// Check if ALL schemas passed
	if len(errors) > 0 {
		// Not all schemas matched
		message := allofNotAllMatchError(ctx.Locale)
		if !isEmptyErrorMessage(s.notAllMatchError) {
			message = resolveErrorMessage(s.notAllMatchError, ctx)
		}

		// Return the main error plus all schema-specific errors
		mainError := NewPrimitiveError(value, message, "allof_not_all_match")
		allErrorsList := append([]ValidationError{mainError}, errors...)
		allErrorsList = append(allErrorsList, allErrors...)

		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: allErrorsList,
		}
	}

	// All schemas matched
	return ParseResult{
		Valid:  true,
		Value:  finalValue,
		Errors: nil,
	}
}

// JSON generates JSON Schema representation
func (s *AllOfSchema) JSON() map[string]interface{} {
	schema := make(map[string]interface{})

	// Generate allOf array with all schemas
	allOfSchemas := make([]interface{}, len(s.schemas))
	for i, subSchema := range s.schemas {
		if jsonSchema, ok := subSchema.(interface{ JSON() map[string]interface{} }); ok {
			allOfSchemas[i] = jsonSchema.JSON()
		} else {
			// Fallback for schemas that don't implement JSON method
			allOfSchemas[i] = map[string]interface{}{"type": "unknown"}
		}
	}
	schema["allOf"] = allOfSchemas

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())

	// Add nullable if true
	if s.nullable {
		// For allOf with nullable, we add a oneOf wrapper
		schema = map[string]interface{}{
			"oneOf": []interface{}{
				schema,
				map[string]interface{}{"type": "null"},
			},
		}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize AllOfSchema for JSON schema generation
func (s *AllOfSchema) MarshalJSON() ([]byte, error) {
	type jsonAllOfSchema struct {
		Schema
		Schemas  []Parseable `json:"schemas"`
		Nullable bool        `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonAllOfSchema{
		Schema:   s.Schema,
		Schemas:  s.schemas,
		Nullable: s.nullable,
	})
}
