package schema

import (
	"encoding/json"

	"github.com/nyxstack/i18n"
)

// Default error messages for any validation
var (
	anyRequiredError = i18n.S("value is required")
	anyEnumError     = i18n.S("value must be one of the allowed values")
	anyConstError    = i18n.S("value must be exactly the specified constant")
)

// AnySchema represents a JSON Schema that accepts any value
type AnySchema struct {
	Schema
	nullable bool // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError ErrorMessage
}

// Any creates a new any schema that accepts any value
func Any(errorMessage ...interface{}) *AnySchema {
	schema := &AnySchema{
		Schema: Schema{
			schemaType: "",    // No specific type for "any"
			required:   false, // Any schema should accept nil by default
		},
		nullable: true, // Any schema is nullable by default
	}
	if len(errorMessage) > 0 {
		schema.requiredError = toErrorMessage(errorMessage[0])
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *AnySchema) Title(title string) *AnySchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *AnySchema) Description(description string) *AnySchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *AnySchema) Default(value interface{}) *AnySchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *AnySchema) Example(example interface{}) *AnySchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values (any types allowed)
func (s *AnySchema) Enum(values []interface{}) *AnySchema {
	s.Schema.enum = values
	return s
}

// Const sets a constant value
func (s *AnySchema) Const(value interface{}) *AnySchema {
	s.Schema.constVal = value
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *AnySchema) Optional() *AnySchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *AnySchema) Required(errorMessage ...interface{}) *AnySchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *AnySchema) Nullable() *AnySchema {
	s.nullable = true
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *AnySchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *AnySchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *AnySchema) IsNullable() bool {
	return s.nullable
}

// Validation

// Parse validates and parses any value, returning the final parsed value
func (s *AnySchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	// Handle nil values
	if value == nil {
		if s.nullable || !s.Schema.required {
			// For nullable or optional schemas, nil is valid
			if defaultVal := s.GetDefault(); defaultVal != nil {
				// Use default value if available
				return ParseResult{Valid: true, Value: defaultVal, Errors: nil}
			}
			return ParseResult{Valid: true, Value: nil, Errors: nil}
		}
		if s.Schema.required {
			// Check if we have a default value to use instead
			if defaultVal := s.GetDefault(); defaultVal != nil {
				// Use default value and re-parse it
				return s.Parse(defaultVal, ctx)
			}
			// No default, required field is missing
			message := anyRequiredError(ctx.Locale)
			if !isEmptyErrorMessage(s.requiredError) {
				message = resolveErrorMessage(s.requiredError, ctx)
			}
			return ParseResult{
				Valid:  false,
				Value:  nil,
				Errors: []ValidationError{NewPrimitiveError(value, message, "required")},
			}
		}
	}

	// For any schema, we accept all non-nil values as-is
	finalValue := value

	// Check enum constraint if present
	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == value {
				valid = true
				break
			}
		}
		if !valid {
			message := anyEnumError(ctx.Locale)
			errors = append(errors, NewPrimitiveError(value, message, "enum"))
		}
	}

	// Check const constraint if present
	if s.Schema.constVal != nil && s.Schema.constVal != value {
		message := anyConstError(ctx.Locale)
		errors = append(errors, NewPrimitiveError(value, message, "const"))
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *AnySchema) JSON() map[string]interface{} {
	schema := make(map[string]interface{})

	// Any schema doesn't specify a type - it accepts everything
	// This is represented by omitting the "type" field entirely

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Any schema can be represented as an empty object {} in JSON Schema
	// which means "accepts anything"

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize AnySchema for JSON schema generation
func (s *AnySchema) MarshalJSON() ([]byte, error) {
	type jsonAnySchema struct {
		Schema
		Nullable bool `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonAnySchema{
		Schema:   s.Schema,
		Nullable: s.nullable,
	})
}
