package schema

import (
	"encoding/json"

	"github.com/nyxstack/i18n"
)

// Default error messages for boolean validation
var (
	boolRequiredError = i18n.S("value is required")
	boolTypeError     = i18n.S("value must be a boolean")
	boolEnumError     = i18n.S("value must be one of the allowed values")
)

func boolConstError(value bool) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %v", value)
}

// BoolSchema represents a JSON Schema for boolean values
type BoolSchema struct {
	Schema
	// Bool-specific validation (private fields)
	nullable bool

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	typeMismatchError ErrorMessage
}

// Bool creates a new bool schema with optional type error message
func Bool(errorMessage ...interface{}) *BoolSchema {
	schema := &BoolSchema{
		Schema: Schema{
			schemaType: "boolean",
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
func (s *BoolSchema) Title(title string) *BoolSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *BoolSchema) Description(description string) *BoolSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *BoolSchema) Default(value interface{}) *BoolSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *BoolSchema) Example(example bool) *BoolSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values with optional custom error message
func (s *BoolSchema) Enum(values []bool, errorMessage ...interface{}) *BoolSchema {
	s.Schema.enum = make([]interface{}, len(values))
	for i, v := range values {
		s.Schema.enum[i] = v
	}
	if len(errorMessage) > 0 {
		s.enumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Const sets a constant value with optional custom error message
func (s *BoolSchema) Const(value bool, errorMessage ...interface{}) *BoolSchema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *BoolSchema) Optional() *BoolSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *BoolSchema) Required(errorMessage ...interface{}) *BoolSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *BoolSchema) Nullable() *BoolSchema {
	s.nullable = true
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *BoolSchema) TypeError(message string) *BoolSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Convenience methods

// True creates a boolean schema that only accepts true
func (s *BoolSchema) True() *BoolSchema {
	return s.Const(true)
}

// False creates a boolean schema that only accepts false
func (s *BoolSchema) False() *BoolSchema {
	return s.Const(false)
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *BoolSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *BoolSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *BoolSchema) IsNullable() bool {
	return s.nullable
}

// GetDefault returns the default value as a bool
func (s *BoolSchema) GetDefaultBool() *bool {
	if s.GetDefault() != nil {
		if b, ok := s.GetDefault().(bool); ok {
			return &b
		}
	}
	return nil
}

// Validation

// Parse validates and parses a boolean value, returning the final parsed value
func (s *BoolSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := boolRequiredError(ctx.Locale)
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

	// Type check
	boolValue, ok := value.(bool)
	if !ok {
		message := boolTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Now validate the bool value against all constraints
	finalValue := boolValue // This is our parsed value

	// Check enum
	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == boolValue {
				valid = true
				break
			}
		}
		if !valid {
			message := boolEnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(boolValue, message, "enum"))
		}
	}

	// Check const
	if s.Schema.constVal != nil {
		if constBool, ok := s.Schema.constVal.(bool); ok && constBool != boolValue {
			message := boolConstError(constBool)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(boolValue, message, "const"))
		}
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *BoolSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("boolean")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"boolean", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize BoolSchema for JSON schema generation
func (s *BoolSchema) MarshalJSON() ([]byte, error) {
	type jsonBoolSchema struct {
		Schema
		Nullable bool `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonBoolSchema{
		Schema:   s.Schema,
		Nullable: s.nullable,
	})
}

// Interface implementations for BoolSchema

// SetTitle implements SetTitle interface
func (s *BoolSchema) SetTitle(title string) {
	s.Title(title)
}

// SetDescription implements SetDescription interface
func (s *BoolSchema) SetDescription(description string) {
	s.Description(description)
}

// SetRequired implements SetRequired interface
func (s *BoolSchema) SetRequired() {
	s.Required()
}

// SetOptional implements SetOptional interface
func (s *BoolSchema) SetOptional() {
	s.Optional()
}

// SetNullable implements SetNullable interface
func (s *BoolSchema) SetNullable() {
	s.Nullable()
}

// SetDefault implements SetDefault interface
func (s *BoolSchema) SetDefault(value interface{}) {
	s.Default(value)
}

// SetExample implements SetExample interface
func (s *BoolSchema) SetExample(example interface{}) {
	if val, ok := example.(bool); ok {
		s.Example(val)
	}
}
