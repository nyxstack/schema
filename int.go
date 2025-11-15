package schema

import (
	"encoding/json"

	"github.com/nyxstack/i18n"
)

// Default error messages for int validation
var (
	intRequiredError = i18n.S("value is required")
	intTypeError     = i18n.S("value must be an integer")
	intEnumError     = i18n.S("value must be one of the allowed values")
)

// Default error message functions that take parameters
func intMinimumError(min int) i18n.TranslatedFunc {
	return i18n.F("value must be at least %d", min)
}

func intMaximumError(max int) i18n.TranslatedFunc {
	return i18n.F("value must be at most %d", max)
}

func intMultipleOfError(multiple int) i18n.TranslatedFunc {
	return i18n.F("value must be a multiple of %d", multiple)
}

func intConstError(value int) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %d", value)
}

// IntSchema represents a JSON Schema for integer values
type IntSchema struct {
	Schema
	// Int-specific validation (private fields)
	minimum    *int
	maximum    *int
	multipleOf *int
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

// Int creates a new int schema with optional type error message
func Int(errorMessage ...interface{}) *IntSchema {
	schema := &IntSchema{
		Schema: Schema{
			schemaType: "integer",
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
func (s *IntSchema) Title(title string) *IntSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *IntSchema) Description(description string) *IntSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *IntSchema) Default(value interface{}) *IntSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *IntSchema) Example(example int) *IntSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values with optional custom error message
func (s *IntSchema) Enum(values []int, errorMessage ...interface{}) *IntSchema {
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
func (s *IntSchema) Const(value int, errorMessage ...interface{}) *IntSchema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *IntSchema) Optional() *IntSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *IntSchema) Required(errorMessage ...interface{}) *IntSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *IntSchema) Nullable() *IntSchema {
	s.nullable = true
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *IntSchema) TypeError(message string) *IntSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Int-specific fluent API methods

// Min sets the minimum value constraint with optional custom error message
func (s *IntSchema) Min(min int, errorMessage ...interface{}) *IntSchema {
	s.minimum = &min
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Max sets the maximum value constraint with optional custom error message
func (s *IntSchema) Max(max int, errorMessage ...interface{}) *IntSchema {
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Range sets both minimum and maximum values with optional custom error message
func (s *IntSchema) Range(min, max int, errorMessage ...interface{}) *IntSchema {
	s.minimum = &min
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MultipleOf sets the multiple constraint with optional custom error message
func (s *IntSchema) MultipleOf(multiple int, errorMessage ...interface{}) *IntSchema {
	s.multipleOf = &multiple
	if len(errorMessage) > 0 {
		s.multipleOfError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *IntSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *IntSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *IntSchema) IsNullable() bool {
	return s.nullable
}

// GetMinimum returns the minimum value constraint
func (s *IntSchema) GetMinimum() *int {
	return s.minimum
}

// GetMaximum returns the maximum value constraint
func (s *IntSchema) GetMaximum() *int {
	return s.maximum
}

// GetMultipleOf returns the multiple constraint
func (s *IntSchema) GetMultipleOf() *int {
	return s.multipleOf
}

// GetDefault returns the default value as an int
func (s *IntSchema) GetDefaultInt() *int {
	if s.GetDefault() != nil {
		if i, ok := s.GetDefault().(int); ok {
			return &i
		}
	}
	return nil
}

// Validation

// Parse validates and parses an integer value, returning the final parsed value
func (s *IntSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := intRequiredError(ctx.Locale)
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
	var intValue int
	var typeValid bool

	switch v := value.(type) {
	case int:
		intValue = v
		typeValid = true
	case int8:
		intValue = int(v)
		typeValid = true
	case int16:
		intValue = int(v)
		typeValid = true
	case int32:
		intValue = int(v)
		typeValid = true
	case int64:
		// Check for overflow when converting int64 to int
		if v >= int64(^uint(0)>>1) || v <= int64(-1-int(^uint(0)>>1)) {
			// Value outside int range
			typeValid = false
		} else {
			intValue = int(v)
			typeValid = true
		}
	case float32:
		// Allow conversion from whole number floats
		if v == float32(int(v)) {
			intValue = int(v)
			typeValid = true
		} else {
			typeValid = false
		}
	case float64:
		// Allow conversion from whole number floats
		if v == float64(int(v)) {
			intValue = int(v)
			typeValid = true
		} else {
			typeValid = false
		}
	default:
		typeValid = false
	}

	if !typeValid {
		message := intTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Now validate the int value against all constraints
	finalValue := intValue // This is our parsed value

	// Check minimum
	if s.minimum != nil && intValue < *s.minimum {
		message := intMinimumError(*s.minimum)(ctx.Locale)
		if !isEmptyErrorMessage(s.minimumError) {
			message = resolveErrorMessage(s.minimumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(intValue, message, "minimum"))
	}

	// Check maximum
	if s.maximum != nil && intValue > *s.maximum {
		message := intMaximumError(*s.maximum)(ctx.Locale)
		if !isEmptyErrorMessage(s.maximumError) {
			message = resolveErrorMessage(s.maximumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(intValue, message, "maximum"))
	}

	// Check multipleOf
	if s.multipleOf != nil && intValue%*s.multipleOf != 0 {
		message := intMultipleOfError(*s.multipleOf)(ctx.Locale)
		if !isEmptyErrorMessage(s.multipleOfError) {
			message = resolveErrorMessage(s.multipleOfError, ctx)
		}
		errors = append(errors, NewPrimitiveError(intValue, message, "multiple_of"))
	}

	// Check enum
	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == intValue {
				valid = true
				break
			}
		}
		if !valid {
			message := intEnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(intValue, message, "enum"))
		}
	}

	// Check const
	if s.Schema.constVal != nil {
		if constInt, ok := s.Schema.constVal.(int); ok && constInt != intValue {
			message := intConstError(constInt)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(intValue, message, "const"))
		}
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *IntSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("integer")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add int-specific fields
	addOptionalField(schema, "minimum", s.minimum)
	addOptionalField(schema, "maximum", s.maximum)
	addOptionalField(schema, "multipleOf", s.multipleOf)

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"integer", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize IntSchema for JSON schema generation
func (s *IntSchema) MarshalJSON() ([]byte, error) {
	type jsonIntSchema struct {
		Schema
		Minimum    *int `json:"minimum,omitempty"`
		Maximum    *int `json:"maximum,omitempty"`
		MultipleOf *int `json:"multipleOf,omitempty"`
		Nullable   bool `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonIntSchema{
		Schema:     s.Schema,
		Minimum:    s.minimum,
		Maximum:    s.maximum,
		MultipleOf: s.multipleOf,
		Nullable:   s.nullable,
	})
}

// Interface implementations for IntSchema

// SetTitle implements SetTitle interface
func (s *IntSchema) SetTitle(title string) {
	s.Title(title)
}

// SetDescription implements SetDescription interface
func (s *IntSchema) SetDescription(description string) {
	s.Description(description)
}

// SetRequired implements SetRequired interface
func (s *IntSchema) SetRequired() {
	s.Required()
}

// SetOptional implements SetOptional interface
func (s *IntSchema) SetOptional() {
	s.Optional()
}

// SetMinimum implements SetMinimum interface
func (s *IntSchema) SetMinimum(min int) {
	s.Min(min)
}

// SetMaximum implements SetMaximum interface
func (s *IntSchema) SetMaximum(max int) {
	s.Max(max)
}

// SetNullable implements SetNullable interface
func (s *IntSchema) SetNullable() {
	s.Nullable()
}

// SetDefault implements SetDefault interface
func (s *IntSchema) SetDefault(value interface{}) {
	s.Default(value)
}

// SetExample implements SetExample interface
func (s *IntSchema) SetExample(example interface{}) {
	if val, ok := example.(int); ok {
		s.Example(val)
	}
}
