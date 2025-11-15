package schema

import (
	"encoding/json"
	"fmt"

	"github.com/nyxstack/i18n"
)

// Default error messages for anyof validation
var (
	anyofRequiredError = i18n.S("value is required")
	anyofNoMatchError  = i18n.S("value must match at least one of the provided schemas")
)

// AnyOfSchema represents a JSON Schema anyOf (value must match AT LEAST one schema)
type AnyOfSchema struct {
	Schema
	schemas  []Parseable // The schemas where AT LEAST one must match
	nullable bool        // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	noMatchError      ErrorMessage
	typeMismatchError ErrorMessage
}

// AnyOf creates a new anyof schema with the provided schemas (at least one must match)
func AnyOf(schemas ...Parseable) *AnyOfSchema {
	schema := &AnyOfSchema{
		Schema: Schema{
			schemaType: "anyOf",
			required:   true, // Default to required
		},
		schemas: schemas,
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *AnyOfSchema) Title(title string) *AnyOfSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *AnyOfSchema) Description(description string) *AnyOfSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *AnyOfSchema) Default(value interface{}) *AnyOfSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *AnyOfSchema) Example(example interface{}) *AnyOfSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Schema manipulation

// Add appends additional schemas to the anyof
func (s *AnyOfSchema) Add(schemas ...Parseable) *AnyOfSchema {
	s.schemas = append(s.schemas, schemas...)
	return s
}

// Schemas returns all schemas in the anyof
func (s *AnyOfSchema) Schemas() []Parseable {
	return s.schemas
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *AnyOfSchema) Optional() *AnyOfSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *AnyOfSchema) Required(errorMessage ...interface{}) *AnyOfSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *AnyOfSchema) Nullable() *AnyOfSchema {
	s.nullable = true
	return s
}

// Error customization

// NoMatchError sets a custom error message when no schemas match
func (s *AnyOfSchema) NoMatchError(message string) *AnyOfSchema {
	s.noMatchError = toErrorMessage(message)
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *AnyOfSchema) TypeError(message string) *AnyOfSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *AnyOfSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *AnyOfSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *AnyOfSchema) IsNullable() bool {
	return s.nullable
}

// GetSchemaCount returns the number of schemas in the anyof
func (s *AnyOfSchema) GetSchemaCount() int {
	return len(s.schemas)
}

// Validation

// Parse validates and parses an anyof value, returning the final parsed value
func (s *AnyOfSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := anyofRequiredError(ctx.Locale)
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

	// Validate against each schema in the anyof
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
					Path:    append([]string{fmt.Sprintf("anyOf[%d]", i)}, err.Path...),
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
		message := anyofNoMatchError(ctx.Locale)
		if !isEmptyErrorMessage(s.noMatchError) {
			message = resolveErrorMessage(s.noMatchError, ctx)
		}
		// Return the original value with no match error, plus all schema errors for context
		errors = append(errors, NewPrimitiveError(value, message, "anyof_no_match"))
		// Also include all the individual schema errors for debugging
		errors = append(errors, allErrors...)
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: errors,
		}
	}

	// At least one schema matched - this is what we want for anyOf
	// Use the first successful result's value
	// (You could implement different strategies here, like using the "best" match)
	return validResults[0]
}

// JSON generates JSON Schema representation
func (s *AnyOfSchema) JSON() map[string]interface{} {
	schema := make(map[string]interface{})

	// Generate anyOf array with all schemas
	anyOfSchemas := make([]interface{}, len(s.schemas))
	for i, subSchema := range s.schemas {
		if jsonSchema, ok := subSchema.(interface{ JSON() map[string]interface{} }); ok {
			anyOfSchemas[i] = jsonSchema.JSON()
		} else {
			// Fallback for schemas that don't implement JSON method
			anyOfSchemas[i] = map[string]interface{}{"type": "unknown"}
		}
	}
	schema["anyOf"] = anyOfSchemas

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())

	// Add nullable if true
	if s.nullable {
		// Add null to the anyOf array
		anyOfSchemas = append(anyOfSchemas, map[string]interface{}{"type": "null"})
		schema["anyOf"] = anyOfSchemas
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize AnyOfSchema for JSON schema generation
func (s *AnyOfSchema) MarshalJSON() ([]byte, error) {
	type jsonAnyOfSchema struct {
		Schema
		Schemas  []Parseable `json:"schemas"`
		Nullable bool        `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonAnyOfSchema{
		Schema:   s.Schema,
		Schemas:  s.schemas,
		Nullable: s.nullable,
	})
}
