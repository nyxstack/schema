package schema

import (
	"encoding/json"
	"math"

	"github.com/nyxstack/i18n"
)

// Default error messages for int16 validation
var (
	int16RequiredError = i18n.S("value is required")
	int16TypeError     = i18n.S("value must be a 16-bit integer")
	int16EnumError     = i18n.S("value must be one of the allowed values")
	int16RangeError    = i18n.S("value must be between -32768 and 32767")
)

// Default error message functions that take parameters
func int16MinimumError(min int16) i18n.TranslatedFunc {
	return i18n.F("value must be at least %d", min)
}

func int16MaximumError(max int16) i18n.TranslatedFunc {
	return i18n.F("value must be at most %d", max)
}

func int16MultipleOfError(multiple int16) i18n.TranslatedFunc {
	return i18n.F("value must be a multiple of %d", multiple)
}

func int16ConstError(value int16) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %d", value)
}

// Int16Schema represents a JSON Schema for int16 values
type Int16Schema struct {
	Schema
	// Int16-specific validation (private fields)
	minimum    *int16
	maximum    *int16
	multipleOf *int16
	nullable   bool

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	minimumError      ErrorMessage
	maximumError      ErrorMessage
	multipleOfError   ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	typeMismatchError ErrorMessage
	rangeError        ErrorMessage
}

// Int16 creates a new int16 schema with optional type error message
func Int16(errorMessage ...interface{}) *Int16Schema {
	schema := &Int16Schema{
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
func (s *Int16Schema) Title(title string) *Int16Schema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *Int16Schema) Description(description string) *Int16Schema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *Int16Schema) Default(value interface{}) *Int16Schema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *Int16Schema) Example(example int16) *Int16Schema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values with optional custom error message
func (s *Int16Schema) Enum(values []int16, errorMessage ...interface{}) *Int16Schema {
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
func (s *Int16Schema) Const(value int16, errorMessage ...interface{}) *Int16Schema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *Int16Schema) Optional() *Int16Schema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *Int16Schema) Required(errorMessage ...interface{}) *Int16Schema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *Int16Schema) Nullable() *Int16Schema {
	s.nullable = true
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *Int16Schema) TypeError(message string) *Int16Schema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Int16-specific fluent API methods

// Min sets the minimum value constraint with optional custom error message
func (s *Int16Schema) Min(min int16, errorMessage ...interface{}) *Int16Schema {
	s.minimum = &min
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Max sets the maximum value constraint with optional custom error message
func (s *Int16Schema) Max(max int16, errorMessage ...interface{}) *Int16Schema {
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Range sets both minimum and maximum values with optional custom error message
func (s *Int16Schema) Range(min, max int16, errorMessage ...interface{}) *Int16Schema {
	s.minimum = &min
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MultipleOf sets the multiple constraint with optional custom error message
func (s *Int16Schema) MultipleOf(multiple int16, errorMessage ...interface{}) *Int16Schema {
	s.multipleOf = &multiple
	if len(errorMessage) > 0 {
		s.multipleOfError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *Int16Schema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *Int16Schema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *Int16Schema) IsNullable() bool {
	return s.nullable
}

// GetMinimum returns the minimum value constraint
func (s *Int16Schema) GetMinimum() *int16 {
	return s.minimum
}

// GetMaximum returns the maximum value constraint
func (s *Int16Schema) GetMaximum() *int16 {
	return s.maximum
}

// GetMultipleOf returns the multiple constraint
func (s *Int16Schema) GetMultipleOf() *int16 {
	return s.multipleOf
}

// GetDefault returns the default value as an int16
func (s *Int16Schema) GetDefaultInt16() *int16 {
	if s.GetDefault() != nil {
		if i, ok := s.GetDefault().(int16); ok {
			return &i
		}
	}
	return nil
}

// Validation

// Parse validates and parses an int16 value, returning the final parsed value
func (s *Int16Schema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	// Handle nil values
	if value == nil {
		if s.nullable {
			return ParseResult{Valid: true, Value: nil, Errors: nil}
		}
		if s.Schema.required {
			if defaultVal := s.GetDefault(); defaultVal != nil {
				return s.Parse(defaultVal, ctx)
			}
			message := int16RequiredError(ctx.Locale)
			if !isEmptyErrorMessage(s.requiredError) {
				message = resolveErrorMessage(s.requiredError, ctx)
			}
			return ParseResult{
				Valid:  false,
				Value:  nil,
				Errors: []ValidationError{NewPrimitiveError(value, message, "required")},
			}
		}
		if defaultVal := s.GetDefault(); defaultVal != nil {
			return s.Parse(defaultVal, ctx)
		}
		return ParseResult{Valid: true, Value: nil, Errors: nil}
	}

	// Type coercion and validation
	var int16Value int16
	var typeValid bool

	switch v := value.(type) {
	case int16:
		int16Value = v
		typeValid = true
	case int8:
		int16Value = int16(v)
		typeValid = true
	case int:
		if v >= math.MinInt16 && v <= math.MaxInt16 {
			int16Value = int16(v)
			typeValid = true
		}
	case int32:
		if v >= math.MinInt16 && v <= math.MaxInt16 {
			int16Value = int16(v)
			typeValid = true
		}
	case int64:
		if v >= math.MinInt16 && v <= math.MaxInt16 {
			int16Value = int16(v)
			typeValid = true
		}
	case float32:
		if v == float32(int(v)) && v >= math.MinInt16 && v <= math.MaxInt16 {
			int16Value = int16(v)
			typeValid = true
		}
	case float64:
		if v == float64(int(v)) && v >= math.MinInt16 && v <= math.MaxInt16 {
			int16Value = int16(v)
			typeValid = true
		}
	}

	if !typeValid {
		message := int16TypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	finalValue := int16Value

	// Validation constraints
	if s.minimum != nil && int16Value < *s.minimum {
		message := int16MinimumError(*s.minimum)(ctx.Locale)
		if !isEmptyErrorMessage(s.minimumError) {
			message = resolveErrorMessage(s.minimumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int16Value, message, "minimum"))
	}

	if s.maximum != nil && int16Value > *s.maximum {
		message := int16MaximumError(*s.maximum)(ctx.Locale)
		if !isEmptyErrorMessage(s.maximumError) {
			message = resolveErrorMessage(s.maximumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int16Value, message, "maximum"))
	}

	if s.multipleOf != nil && int16Value%*s.multipleOf != 0 {
		message := int16MultipleOfError(*s.multipleOf)(ctx.Locale)
		if !isEmptyErrorMessage(s.multipleOfError) {
			message = resolveErrorMessage(s.multipleOfError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int16Value, message, "multiple_of"))
	}

	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == int16Value {
				valid = true
				break
			}
		}
		if !valid {
			message := int16EnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int16Value, message, "enum"))
		}
	}

	if s.Schema.constVal != nil {
		if constInt16, ok := s.Schema.constVal.(int16); ok && constInt16 != int16Value {
			message := int16ConstError(constInt16)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int16Value, message, "const"))
		}
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *Int16Schema) JSON() map[string]interface{} {
	schema := baseJSONSchema("integer")

	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	if s.minimum != nil {
		schema["minimum"] = int(*s.minimum)
	}
	if s.maximum != nil {
		schema["maximum"] = int(*s.maximum)
	}
	if s.multipleOf != nil {
		schema["multipleOf"] = int(*s.multipleOf)
	}

	schema["format"] = "int16"

	if s.nullable {
		schema["type"] = []string{"integer", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler
func (s *Int16Schema) MarshalJSON() ([]byte, error) {
	type jsonInt16Schema struct {
		Schema
		Minimum    *int16 `json:"minimum,omitempty"`
		Maximum    *int16 `json:"maximum,omitempty"`
		MultipleOf *int16 `json:"multipleOf,omitempty"`
		Format     string `json:"format"`
		Nullable   bool   `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonInt16Schema{
		Schema:     s.Schema,
		Minimum:    s.minimum,
		Maximum:    s.maximum,
		MultipleOf: s.multipleOf,
		Format:     "int16",
		Nullable:   s.nullable,
	})
}
