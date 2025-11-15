package schema

import (
	"encoding/json"

	"github.com/nyxstack/i18n"
)

// Default error messages for number validation
var (
	numberRequiredError = i18n.S("value is required")
	numberTypeError     = i18n.S("value must be a number")
	numberEnumError     = i18n.S("value must be one of the allowed values")
)

// Default error message functions that take parameters
func numberMinimumError(min float64) i18n.TranslatedFunc {
	return i18n.F("value must be at least %g", min)
}

func numberMaximumError(max float64) i18n.TranslatedFunc {
	return i18n.F("value must be at most %g", max)
}

func numberMultipleOfError(multiple float64) i18n.TranslatedFunc {
	return i18n.F("value must be a multiple of %g", multiple)
}

func numberConstError(value float64) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %g", value)
}

// NumberSchema represents a JSON Schema for float64 values
type NumberSchema struct {
	Schema
	// Number-specific validation (private fields)
	minimum    *float64
	maximum    *float64
	multipleOf *float64
	nullable   bool

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	minimumError      ErrorMessage
	maximumError      ErrorMessage
	multipleOfError   ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	typeMismatchError ErrorMessage
}

// Number creates a new number schema with optional type error message
func Number(errorMessage ...interface{}) *NumberSchema {
	schema := &NumberSchema{
		Schema: Schema{
			schemaType: "number",
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
func (s *NumberSchema) Title(title string) *NumberSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *NumberSchema) Description(description string) *NumberSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *NumberSchema) Default(value interface{}) *NumberSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *NumberSchema) Example(example float64) *NumberSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values with optional custom error message
func (s *NumberSchema) Enum(values []float64, errorMessage ...interface{}) *NumberSchema {
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
func (s *NumberSchema) Const(value float64, errorMessage ...interface{}) *NumberSchema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *NumberSchema) Optional() *NumberSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *NumberSchema) Required(errorMessage ...interface{}) *NumberSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *NumberSchema) Nullable() *NumberSchema {
	s.nullable = true
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *NumberSchema) TypeError(message string) *NumberSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Number-specific fluent API methods

// Min sets the minimum value constraint with optional custom error message
func (s *NumberSchema) Min(min float64, errorMessage ...interface{}) *NumberSchema {
	s.minimum = &min
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Max sets the maximum value constraint with optional custom error message
func (s *NumberSchema) Max(max float64, errorMessage ...interface{}) *NumberSchema {
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Range sets both minimum and maximum values with optional custom error message
func (s *NumberSchema) Range(min, max float64, errorMessage ...interface{}) *NumberSchema {
	s.minimum = &min
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MultipleOf sets the multiple constraint with optional custom error message
func (s *NumberSchema) MultipleOf(multiple float64, errorMessage ...interface{}) *NumberSchema {
	s.multipleOf = &multiple
	if len(errorMessage) > 0 {
		s.multipleOfError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *NumberSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *NumberSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *NumberSchema) IsNullable() bool {
	return s.nullable
}

// GetMinimum returns the minimum value constraint
func (s *NumberSchema) GetMinimum() *float64 {
	return s.minimum
}

// GetMaximum returns the maximum value constraint
func (s *NumberSchema) GetMaximum() *float64 {
	return s.maximum
}

// GetMultipleOf returns the multiple constraint
func (s *NumberSchema) GetMultipleOf() *float64 {
	return s.multipleOf
}

// GetDefault returns the default value as a float64
func (s *NumberSchema) GetDefaultNumber() *float64 {
	if s.GetDefault() != nil {
		if f, ok := s.GetDefault().(float64); ok {
			return &f
		}
	}
	return nil
}

// Validation

// Parse validates and parses a number value, returning the final parsed value
func (s *NumberSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := numberRequiredError(ctx.Locale)
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

	// Type coercion and validation
	var numValue float64
	var typeValid bool

	switch v := value.(type) {
	case float64:
		numValue = v
		typeValid = true
	case float32:
		numValue = float64(v)
		typeValid = true
	case int:
		numValue = float64(v)
		typeValid = true
	case int8:
		numValue = float64(v)
		typeValid = true
	case int16:
		numValue = float64(v)
		typeValid = true
	case int32:
		numValue = float64(v)
		typeValid = true
	case int64:
		numValue = float64(v)
		typeValid = true
	default:
		typeValid = false
	}

	if !typeValid {
		message := numberTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Now validate the number value against all constraints
	finalValue := numValue // This is our parsed value

	// Check minimum
	if s.minimum != nil && numValue < *s.minimum {
		message := numberMinimumError(*s.minimum)(ctx.Locale)
		if !isEmptyErrorMessage(s.minimumError) {
			message = resolveErrorMessage(s.minimumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(numValue, message, "minimum"))
	}

	// Check maximum
	if s.maximum != nil && numValue > *s.maximum {
		message := numberMaximumError(*s.maximum)(ctx.Locale)
		if !isEmptyErrorMessage(s.maximumError) {
			message = resolveErrorMessage(s.maximumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(numValue, message, "maximum"))
	}

	// Check multipleOf (for numbers, we need to handle floating point precision)
	if s.multipleOf != nil {
		quotient := numValue / *s.multipleOf
		if quotient != float64(int64(quotient+0.5)) { // Check if it's close to an integer
			message := numberMultipleOfError(*s.multipleOf)(ctx.Locale)
			if !isEmptyErrorMessage(s.multipleOfError) {
				message = resolveErrorMessage(s.multipleOfError, ctx)
			}
			errors = append(errors, NewPrimitiveError(numValue, message, "multiple_of"))
		}
	}

	// Check enum
	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == numValue {
				valid = true
				break
			}
		}
		if !valid {
			message := numberEnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(numValue, message, "enum"))
		}
	}

	// Check const
	if s.Schema.constVal != nil {
		if constFloat, ok := s.Schema.constVal.(float64); ok && constFloat != numValue {
			message := numberConstError(constFloat)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(numValue, message, "const"))
		}
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *NumberSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("number")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add number-specific fields
	addOptionalField(schema, "minimum", s.minimum)
	addOptionalField(schema, "maximum", s.maximum)
	addOptionalField(schema, "multipleOf", s.multipleOf)

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"number", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize NumberSchema for JSON schema generation
func (s *NumberSchema) MarshalJSON() ([]byte, error) {
	type jsonNumberSchema struct {
		Schema
		Minimum    *float64 `json:"minimum,omitempty"`
		Maximum    *float64 `json:"maximum,omitempty"`
		MultipleOf *float64 `json:"multipleOf,omitempty"`
		Nullable   bool     `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonNumberSchema{
		Schema:     s.Schema,
		Minimum:    s.minimum,
		Maximum:    s.maximum,
		MultipleOf: s.multipleOf,
		Nullable:   s.nullable,
	})
}

// Interface implementations for NumberSchema

// SetTitle implements SetTitle interface
func (s *NumberSchema) SetTitle(title string) {
	s.Title(title)
}

// SetDescription implements SetDescription interface
func (s *NumberSchema) SetDescription(description string) {
	s.Description(description)
}

// SetRequired implements SetRequired interface
func (s *NumberSchema) SetRequired() {
	s.Required()
}

// SetOptional implements SetOptional interface
func (s *NumberSchema) SetOptional() {
	s.Optional()
}

// SetMinimumFloat implements SetMinimumFloat interface
func (s *NumberSchema) SetMinimumFloat(min float64) {
	s.Min(min)
}

// SetMaximumFloat implements SetMaximumFloat interface
func (s *NumberSchema) SetMaximumFloat(max float64) {
	s.Max(max)
}

// SetNullable implements SetNullable interface
func (s *NumberSchema) SetNullable() {
	s.Nullable()
}

// SetDefault implements SetDefault interface
func (s *NumberSchema) SetDefault(value interface{}) {
	s.Default(value)
}

// SetExample implements SetExample interface
func (s *NumberSchema) SetExample(example interface{}) {
	if val, ok := example.(float64); ok {
		s.Example(val)
	}
}
