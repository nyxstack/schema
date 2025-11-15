package schema

import (
	"math"

	"github.com/nyxstack/i18n"
)

var (
	int32RequiredError = i18n.S("value is required")
	int32TypeError     = i18n.S("value must be a 32-bit integer")
	int32EnumError     = i18n.S("value must be one of the allowed values")
)

func int32MinimumError(min int32) i18n.TranslatedFunc {
	return i18n.F("value must be at least %d", min)
}

func int32MaximumError(max int32) i18n.TranslatedFunc {
	return i18n.F("value must be at most %d", max)
}

func int32MultipleOfError(multiple int32) i18n.TranslatedFunc {
	return i18n.F("value must be a multiple of %d", multiple)
}

func int32ConstError(value int32) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %d", value)
}

type Int32Schema struct {
	Schema
	minimum    *int32
	maximum    *int32
	multipleOf *int32
	nullable   bool

	requiredError     ErrorMessage
	minimumError      ErrorMessage
	maximumError      ErrorMessage
	multipleOfError   ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	typeMismatchError ErrorMessage
}

func Int32(errorMessage ...interface{}) *Int32Schema {
	schema := &Int32Schema{
		Schema: Schema{
			schemaType: "integer",
			required:   true,
		},
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

func (s *Int32Schema) Title(title string) *Int32Schema {
	s.Schema.title = title
	return s
}

func (s *Int32Schema) Description(description string) *Int32Schema {
	s.Schema.description = description
	return s
}

func (s *Int32Schema) Default(value interface{}) *Int32Schema {
	s.Schema.defaultValue = value
	return s
}

func (s *Int32Schema) Example(example int32) *Int32Schema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

func (s *Int32Schema) Enum(values []int32, errorMessage ...interface{}) *Int32Schema {
	s.Schema.enum = make([]interface{}, len(values))
	for i, v := range values {
		s.Schema.enum[i] = v
	}
	if len(errorMessage) > 0 {
		s.enumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int32Schema) Const(value int32, errorMessage ...interface{}) *Int32Schema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int32Schema) Optional() *Int32Schema {
	s.Schema.required = false
	return s
}

func (s *Int32Schema) Required(errorMessage ...interface{}) *Int32Schema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int32Schema) Nullable() *Int32Schema {
	s.nullable = true
	return s
}

func (s *Int32Schema) Min(min int32, errorMessage ...interface{}) *Int32Schema {
	s.minimum = &min
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int32Schema) Max(max int32, errorMessage ...interface{}) *Int32Schema {
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int32Schema) Range(min, max int32, errorMessage ...interface{}) *Int32Schema {
	s.minimum = &min
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int32Schema) MultipleOf(multiple int32, errorMessage ...interface{}) *Int32Schema {
	s.multipleOf = &multiple
	if len(errorMessage) > 0 {
		s.multipleOfError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int32Schema) IsRequired() bool      { return s.Schema.required }
func (s *Int32Schema) IsOptional() bool      { return !s.Schema.required }
func (s *Int32Schema) IsNullable() bool      { return s.nullable }
func (s *Int32Schema) GetMinimum() *int32    { return s.minimum }
func (s *Int32Schema) GetMaximum() *int32    { return s.maximum }
func (s *Int32Schema) GetMultipleOf() *int32 { return s.multipleOf }

func (s *Int32Schema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	if value == nil {
		if s.nullable {
			return ParseResult{Valid: true, Value: nil, Errors: nil}
		}
		if s.Schema.required {
			if defaultVal := s.GetDefault(); defaultVal != nil {
				return s.Parse(defaultVal, ctx)
			}
			message := int32RequiredError(ctx.Locale)
			if !isEmptyErrorMessage(s.requiredError) {
				message = resolveErrorMessage(s.requiredError, ctx)
			}
			return ParseResult{Valid: false, Value: nil, Errors: []ValidationError{NewPrimitiveError(value, message, "required")}}
		}
		if defaultVal := s.GetDefault(); defaultVal != nil {
			return s.Parse(defaultVal, ctx)
		}
		return ParseResult{Valid: true, Value: nil, Errors: nil}
	}

	var int32Value int32
	var typeValid bool

	switch v := value.(type) {
	case int32:
		int32Value = v
		typeValid = true
	case int8:
		int32Value = int32(v)
		typeValid = true
	case int16:
		int32Value = int32(v)
		typeValid = true
	case int:
		if v >= math.MinInt32 && v <= math.MaxInt32 {
			int32Value = int32(v)
			typeValid = true
		}
	case int64:
		if v >= math.MinInt32 && v <= math.MaxInt32 {
			int32Value = int32(v)
			typeValid = true
		}
	case float32:
		if v == float32(int(v)) && v >= math.MinInt32 && v <= math.MaxInt32 {
			int32Value = int32(v)
			typeValid = true
		}
	case float64:
		if v == float64(int(v)) && v >= math.MinInt32 && v <= math.MaxInt32 {
			int32Value = int32(v)
			typeValid = true
		}
	}

	if !typeValid {
		message := int32TypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{Valid: false, Value: nil, Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")}}
	}

	finalValue := int32Value

	if s.minimum != nil && int32Value < *s.minimum {
		message := int32MinimumError(*s.minimum)(ctx.Locale)
		if !isEmptyErrorMessage(s.minimumError) {
			message = resolveErrorMessage(s.minimumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int32Value, message, "minimum"))
	}

	if s.maximum != nil && int32Value > *s.maximum {
		message := int32MaximumError(*s.maximum)(ctx.Locale)
		if !isEmptyErrorMessage(s.maximumError) {
			message = resolveErrorMessage(s.maximumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int32Value, message, "maximum"))
	}

	if s.multipleOf != nil && int32Value%*s.multipleOf != 0 {
		message := int32MultipleOfError(*s.multipleOf)(ctx.Locale)
		if !isEmptyErrorMessage(s.multipleOfError) {
			message = resolveErrorMessage(s.multipleOfError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int32Value, message, "multiple_of"))
	}

	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == int32Value {
				valid = true
				break
			}
		}
		if !valid {
			message := int32EnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int32Value, message, "enum"))
		}
	}

	if s.Schema.constVal != nil {
		if constInt32, ok := s.Schema.constVal.(int32); ok && constInt32 != int32Value {
			message := int32ConstError(constInt32)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int32Value, message, "const"))
		}
	}

	return ParseResult{Valid: len(errors) == 0, Value: finalValue, Errors: errors}
}

func (s *Int32Schema) JSON() map[string]interface{} {
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

	schema["format"] = "int32"

	if s.nullable {
		schema["type"] = []string{"integer", "null"}
	}

	return schema
}
