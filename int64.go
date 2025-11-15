package schema

import (
	"github.com/nyxstack/i18n"
)

var (
	int64RequiredError = i18n.S("value is required")
	int64TypeError     = i18n.S("value must be a 64-bit integer")
	int64EnumError     = i18n.S("value must be one of the allowed values")
)

func int64MinimumError(min int64) i18n.TranslatedFunc {
	return i18n.F("value must be at least %d", min)
}

func int64MaximumError(max int64) i18n.TranslatedFunc {
	return i18n.F("value must be at most %d", max)
}

func int64MultipleOfError(multiple int64) i18n.TranslatedFunc {
	return i18n.F("value must be a multiple of %d", multiple)
}

func int64ConstError(value int64) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %d", value)
}

type Int64Schema struct {
	Schema
	minimum    *int64
	maximum    *int64
	multipleOf *int64
	nullable   bool

	requiredError     ErrorMessage
	minimumError      ErrorMessage
	maximumError      ErrorMessage
	multipleOfError   ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	typeMismatchError ErrorMessage
}

func Int64(errorMessage ...interface{}) *Int64Schema {
	schema := &Int64Schema{
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

func (s *Int64Schema) Title(title string) *Int64Schema { s.Schema.title = title; return s }
func (s *Int64Schema) Description(description string) *Int64Schema {
	s.Schema.description = description
	return s
}
func (s *Int64Schema) Default(value interface{}) *Int64Schema {
	s.Schema.defaultValue = value
	return s
}
func (s *Int64Schema) Example(example int64) *Int64Schema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}
func (s *Int64Schema) Optional() *Int64Schema { s.Schema.required = false; return s }
func (s *Int64Schema) Nullable() *Int64Schema { s.nullable = true; return s }

func (s *Int64Schema) Enum(values []int64, errorMessage ...interface{}) *Int64Schema {
	s.Schema.enum = make([]interface{}, len(values))
	for i, v := range values {
		s.Schema.enum[i] = v
	}
	if len(errorMessage) > 0 {
		s.enumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int64Schema) Const(value int64, errorMessage ...interface{}) *Int64Schema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int64Schema) Min(min int64, errorMessage ...interface{}) *Int64Schema {
	s.minimum = &min
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int64Schema) Max(max int64, errorMessage ...interface{}) *Int64Schema {
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int64Schema) Range(min, max int64, errorMessage ...interface{}) *Int64Schema {
	s.minimum = &min
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int64Schema) MultipleOf(multiple int64, errorMessage ...interface{}) *Int64Schema {
	s.multipleOf = &multiple
	if len(errorMessage) > 0 {
		s.multipleOfError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *Int64Schema) IsRequired() bool      { return s.Schema.required }
func (s *Int64Schema) IsOptional() bool      { return !s.Schema.required }
func (s *Int64Schema) IsNullable() bool      { return s.nullable }
func (s *Int64Schema) GetMinimum() *int64    { return s.minimum }
func (s *Int64Schema) GetMaximum() *int64    { return s.maximum }
func (s *Int64Schema) GetMultipleOf() *int64 { return s.multipleOf }

func (s *Int64Schema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	if value == nil {
		if s.nullable {
			return ParseResult{Valid: true, Value: nil, Errors: nil}
		}
		if s.Schema.required {
			if defaultVal := s.GetDefault(); defaultVal != nil {
				return s.Parse(defaultVal, ctx)
			}
			message := int64RequiredError(ctx.Locale)
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

	var int64Value int64
	var typeValid bool

	switch v := value.(type) {
	case int64:
		int64Value = v
		typeValid = true
	case int:
		int64Value = int64(v)
		typeValid = true
	case int8:
		int64Value = int64(v)
		typeValid = true
	case int16:
		int64Value = int64(v)
		typeValid = true
	case int32:
		int64Value = int64(v)
		typeValid = true
	case float32:
		if v == float32(int64(v)) {
			int64Value = int64(v)
			typeValid = true
		}
	case float64:
		if v == float64(int64(v)) {
			int64Value = int64(v)
			typeValid = true
		}
	}

	if !typeValid {
		message := int64TypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{Valid: false, Value: nil, Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")}}
	}

	finalValue := int64Value

	if s.minimum != nil && int64Value < *s.minimum {
		message := int64MinimumError(*s.minimum)(ctx.Locale)
		if !isEmptyErrorMessage(s.minimumError) {
			message = resolveErrorMessage(s.minimumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int64Value, message, "minimum"))
	}

	if s.maximum != nil && int64Value > *s.maximum {
		message := int64MaximumError(*s.maximum)(ctx.Locale)
		if !isEmptyErrorMessage(s.maximumError) {
			message = resolveErrorMessage(s.maximumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int64Value, message, "maximum"))
	}

	if s.multipleOf != nil && int64Value%*s.multipleOf != 0 {
		message := int64MultipleOfError(*s.multipleOf)(ctx.Locale)
		if !isEmptyErrorMessage(s.multipleOfError) {
			message = resolveErrorMessage(s.multipleOfError, ctx)
		}
		errors = append(errors, NewPrimitiveError(int64Value, message, "multiple_of"))
	}

	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == int64Value {
				valid = true
				break
			}
		}
		if !valid {
			message := int64EnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int64Value, message, "enum"))
		}
	}

	if s.Schema.constVal != nil {
		if constInt64, ok := s.Schema.constVal.(int64); ok && constInt64 != int64Value {
			message := int64ConstError(constInt64)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(int64Value, message, "const"))
		}
	}

	return ParseResult{Valid: len(errors) == 0, Value: finalValue, Errors: errors}
}

func (s *Int64Schema) JSON() map[string]interface{} {
	schema := baseJSONSchema("integer")
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	if s.minimum != nil {
		schema["minimum"] = *s.minimum
	}
	if s.maximum != nil {
		schema["maximum"] = *s.maximum
	}
	if s.multipleOf != nil {
		schema["multipleOf"] = *s.multipleOf
	}

	schema["format"] = "int64"

	if s.nullable {
		schema["type"] = []string{"integer", "null"}
	}

	return schema
}
