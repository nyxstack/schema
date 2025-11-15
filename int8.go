package schema

import (
	"encoding/json"
	"math"

	"github.com/nyxstack/i18n"
)

// Default error messages for int8 validation
var (
	int8RequiredError = i18n.S("value is required")
	int8TypeError     = i18n.S("value must be an 8-bit integer")
	int8EnumError     = i18n.S("value must be one of the allowed values")
	int8RangeError    = i18n.S("value must be between -128 and 127")
)

// Default error message functions that take parameters
func int8MinimumError(min int8) i18n.TranslatedFunc {
	return i18n.F("value must be at least %d", min)
}

func int8MaximumError(max int8) i18n.TranslatedFunc {
	return i18n.F("value must be at most %d", max)
}

func int8MultipleOfError(multiple int8) i18n.TranslatedFunc {
	return i18n.F("value must be a multiple of %d", multiple)
}

func int8ConstError(value int8) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %d", value)
}

// Int8Schema represents a JSON Schema for int8 values
type Int8Schema struct {
	Schema
	// Int8-specific validation (private fields)
	minimum    *int8
	maximum    *int8
	multipleOf *int8
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

// Int8 creates a new int8 schema with optional type error message
func Int8(errorMessage ...interface{}) *Int8Schema {
	schema := &Int8Schema{
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
func (s *Int8Schema) Title(title string) *Int8Schema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *Int8Schema) Description(description string) *Int8Schema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *Int8Schema) Default(value interface{}) *Int8Schema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *Int8Schema) Example(example int8) *Int8Schema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values with optional custom error message
func (s *Int8Schema) Enum(values []int8, errorMessage ...interface{}) *Int8Schema {
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
func (s *Int8Schema) Const(value int8, errorMessage ...interface{}) *Int8Schema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *Int8Schema) Optional() *Int8Schema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *Int8Schema) Required(errorMessage ...interface{}) *Int8Schema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *Int8Schema) Nullable() *Int8Schema {
	s.nullable = true
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *Int8Schema) TypeError(message string) *Int8Schema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// Int8-specific fluent API methods

// Min sets the minimum value constraint with optional custom error message
func (s *Int8Schema) Min(min int8, errorMessage ...interface{}) *Int8Schema {
	s.minimum = &min
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Max sets the maximum value constraint with optional custom error message
func (s *Int8Schema) Max(max int8, errorMessage ...interface{}) *Int8Schema {
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Range sets both minimum and maximum values with optional custom error message
func (s *Int8Schema) Range(min, max int8, errorMessage ...interface{}) *Int8Schema {
	s.minimum = &min
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MultipleOf sets the multiple constraint with optional custom error message
func (s *Int8Schema) MultipleOf(multiple int8, errorMessage ...interface{}) *Int8Schema {
	s.multipleOf = &multiple
	if len(errorMessage) > 0 {
		s.multipleOfError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *Int8Schema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *Int8Schema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *Int8Schema) IsNullable() bool {
	return s.nullable
}

// GetMinimum returns the minimum value constraint
func (s *Int8Schema) GetMinimum() *int8 {
	return s.minimum
}

// GetMaximum returns the maximum value constraint
func (s *Int8Schema) GetMaximum() *int8 {
	return s.maximum
}

// GetMultipleOf returns the multiple constraint
func (s *Int8Schema) GetMultipleOf() *int8 {
	return s.multipleOf
}

// GetDefault returns the default value as an int8
func (s *Int8Schema) GetDefaultInt8() *int8 {
	if s.GetDefault() != nil {
		if i, ok := s.GetDefault().(int8); ok {
			return &i
		}
	}
	return nil
}

// Validation

// Parse validates and parses an int8 value, returning the final parsed value
func (s *Int8Schema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := int8RequiredError(ctx.Locale)
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
	var int8Value int8
	var typeValid bool

	switch v := value.(type) {
	case int8:
		int8Value = v
		typeValid = true
	case int:
		if v >= math.MinInt8 && v <= math.MaxInt8 {
			int8Value = int8(v)
			typeValid = true
		} else {
			typeValid = false
		}
	case int16:
		if v >= math.MinInt8 && v <= math.MaxInt8 {
			int8Value = int8(v)
			typeValid = true
		} else {
			typeValid = false
		}
	case int32:
		if v >= math.MinInt8 && v <= math.MaxInt8 {
			int8Value = int8(v)
			typeValid = true
		} else {
			typeValid = false
		}
	case int64:
		if v >= math.MinInt8 && v <= math.MaxInt8 {
			int8Value = int8(v)
			typeValid = true
		} else {
			typeValid = false
		}
	case float32:
		// Allow conversion from whole number floats within range
		if v == float32(int(v)) && v >= math.MinInt8 && v <= math.MaxInt8 {
			int8Value = int8(v)
			typeValid = true
		} else {
			typeValid = false
		}
	case float64:
		// Allow conversion from whole number floats within range
		if v == float64(int(v)) && v >= math.MinInt8 && v <= math.MaxInt8 {
			int8Value = int8(v)
			typeValid = true
		} else {
			typeValid = false
		}
	default:
		typeValid = false
	}

	if !typeValid {
		message := int8TypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		} else if !isEmptyErrorMessage(s.rangeError) {
			message = int8RangeError(ctx.Locale)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Now validate the int8 value against all constraints
	finalValue := int8Value // This is our parsed value

	// Check minimum
	if s.minimum != nil && int8Value < *s.minimum {
		message := int8MinimumError(*s.minimum)(ctx.Locale)
		if !isEmptyErrorMessage(s.minimumError) {
			message = resolveErrorMessage(s.minimumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int8Value, message, "minimum"))
	}

	// Check maximum
	if s.maximum != nil && int8Value > *s.maximum {
		message := int8MaximumError(*s.maximum)(ctx.Locale)
		if !isEmptyErrorMessage(s.maximumError) {
			message = resolveErrorMessage(s.maximumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int8Value, message, "maximum"))
	}

	// Check multipleOf
	if s.multipleOf != nil && int8Value%*s.multipleOf != 0 {
		message := int8MultipleOfError(*s.multipleOf)(ctx.Locale)
		if !isEmptyErrorMessage(s.multipleOfError) {
			message = resolveErrorMessage(s.multipleOfError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int8Value, message, "multiple_of"))
	}

	// Check enum
	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == int8Value {
				valid = true
				break
			}
		}
		if !valid {
			message := int8EnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int8Value, message, "enum"))
		}
	}

	// Check const
	if s.Schema.constVal != nil {
		if constInt8, ok := s.Schema.constVal.(int8); ok && constInt8 != int8Value {
			message := int8ConstError(constInt8)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int8Value, message, "const"))
		}
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *Int8Schema) JSON() map[string]interface{} {
	schema := baseJSONSchema("integer")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add int8-specific fields (converted to regular int for JSON)
	if s.minimum != nil {
		schema["minimum"] = int(*s.minimum)
	}
	if s.maximum != nil {
		schema["maximum"] = int(*s.maximum)
	}
	if s.multipleOf != nil {
		schema["multipleOf"] = int(*s.multipleOf)
	}

	// Add format to indicate this is an int8
	schema["format"] = "int8"

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"integer", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize Int8Schema for JSON schema generation
func (s *Int8Schema) MarshalJSON() ([]byte, error) {
	type jsonInt8Schema struct {
		Schema
		Minimum    *int8  `json:"minimum,omitempty"`
		Maximum    *int8  `json:"maximum,omitempty"`
		MultipleOf *int8  `json:"multipleOf,omitempty"`
		Format     string `json:"format"`
		Nullable   bool   `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonInt8Schema{
		Schema:     s.Schema,
		Minimum:    s.minimum,
		Maximum:    s.maximum,
		MultipleOf: s.multipleOf,
		Format:     "int8",
		Nullable:   s.nullable,
	})
}
