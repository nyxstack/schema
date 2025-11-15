package schema

import (
	"encoding/json"
	"math"

	"github.com/nyxstack/i18n"
)

var (
	floatRequiredError = i18n.S("value is required")
	floatTypeError     = i18n.S("value must be a 32-bit float")
	floatEnumError     = i18n.S("value must be one of the allowed values")
)

func floatMinimumError(min float32) i18n.TranslatedFunc {
	return i18n.F("value must be at least %g", min)
}

func floatMaximumError(max float32) i18n.TranslatedFunc {
	return i18n.F("value must be at most %g", max)
}

func floatMultipleOfError(multiple float32) i18n.TranslatedFunc {
	return i18n.F("value must be a multiple of %g", multiple)
}

func floatConstError(value float32) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %g", value)
}

type FloatSchema struct {
	Schema
	minimum    *float32
	maximum    *float32
	multipleOf *float32
	nullable   bool

	requiredError     ErrorMessage
	minimumError      ErrorMessage
	maximumError      ErrorMessage
	multipleOfError   ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	typeMismatchError ErrorMessage
}

func Float(errorMessage ...interface{}) *FloatSchema {
	schema := &FloatSchema{
		Schema: Schema{
			schemaType: "number",
			required:   true,
		},
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

func (s *FloatSchema) Title(title string) *FloatSchema { s.Schema.title = title; return s }
func (s *FloatSchema) Description(description string) *FloatSchema {
	s.Schema.description = description
	return s
}
func (s *FloatSchema) Default(value interface{}) *FloatSchema {
	s.Schema.defaultValue = value
	return s
}
func (s *FloatSchema) Example(example float32) *FloatSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}
func (s *FloatSchema) Optional() *FloatSchema { s.Schema.required = false; return s }
func (s *FloatSchema) Nullable() *FloatSchema { s.nullable = true; return s }

func (s *FloatSchema) Enum(values []float32, errorMessage ...interface{}) *FloatSchema {
	s.Schema.enum = make([]interface{}, len(values))
	for i, v := range values {
		s.Schema.enum[i] = v
	}
	if len(errorMessage) > 0 {
		s.enumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *FloatSchema) Const(value float32, errorMessage ...interface{}) *FloatSchema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *FloatSchema) Min(min float32, errorMessage ...interface{}) *FloatSchema {
	s.minimum = &min
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *FloatSchema) Max(max float32, errorMessage ...interface{}) *FloatSchema {
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *FloatSchema) Range(min, max float32, errorMessage ...interface{}) *FloatSchema {
	s.minimum = &min
	s.maximum = &max
	if len(errorMessage) > 0 {
		s.minimumError = toErrorMessage(errorMessage[0])
		s.maximumError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *FloatSchema) MultipleOf(multiple float32, errorMessage ...interface{}) *FloatSchema {
	s.multipleOf = &multiple
	if len(errorMessage) > 0 {
		s.multipleOfError = toErrorMessage(errorMessage[0])
	}
	return s
}

func (s *FloatSchema) IsRequired() bool        { return s.Schema.required }
func (s *FloatSchema) IsOptional() bool        { return !s.Schema.required }
func (s *FloatSchema) IsNullable() bool        { return s.nullable }
func (s *FloatSchema) GetMinimum() *float32    { return s.minimum }
func (s *FloatSchema) GetMaximum() *float32    { return s.maximum }
func (s *FloatSchema) GetMultipleOf() *float32 { return s.multipleOf }

func (s *FloatSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	if value == nil {
		if s.nullable {
			return ParseResult{Valid: true, Value: nil, Errors: nil}
		}
		if s.Schema.required {
			if defaultVal := s.GetDefault(); defaultVal != nil {
				return s.Parse(defaultVal, ctx)
			}
			message := floatRequiredError(ctx.Locale)
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

	var floatValue float32
	var typeValid bool

	switch v := value.(type) {
	case float32:
		floatValue = v
		typeValid = true
	case float64:
		if v >= -math.MaxFloat32 && v <= math.MaxFloat32 {
			floatValue = float32(v)
			typeValid = true
		}
	case int:
		floatValue = float32(v)
		typeValid = true
	case int8:
		floatValue = float32(v)
		typeValid = true
	case int16:
		floatValue = float32(v)
		typeValid = true
	case int32:
		floatValue = float32(v)
		typeValid = true
	case int64:
		floatValue = float32(v)
		typeValid = true
	}

	if !typeValid {
		message := floatTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{Valid: false, Value: nil, Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")}}
	}

	finalValue := floatValue

	if s.minimum != nil && floatValue < *s.minimum {
		message := floatMinimumError(*s.minimum)(ctx.Locale)
		if !isEmptyErrorMessage(s.minimumError) {
			message = resolveErrorMessage(s.minimumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(floatValue, message, "minimum"))
	}

	if s.maximum != nil && floatValue > *s.maximum {
		message := floatMaximumError(*s.maximum)(ctx.Locale)
		if !isEmptyErrorMessage(s.maximumError) {
			message = resolveErrorMessage(s.maximumError, ctx)
		}
		errors = append(errors, NewPrimitiveError(floatValue, message, "maximum"))
	}

	if s.multipleOf != nil {
		quotient := floatValue / *s.multipleOf
		if quotient != float32(int(quotient+0.5)) {
			message := floatMultipleOfError(*s.multipleOf)(ctx.Locale)
			if !isEmptyErrorMessage(s.multipleOfError) {
				message = resolveErrorMessage(s.multipleOfError, ctx)
			}
			errors = append(errors, NewPrimitiveError(floatValue, message, "multiple_of"))
		}
	}

	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == floatValue {
				valid = true
				break
			}
		}
		if !valid {
			message := floatEnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(floatValue, message, "enum"))
		}
	}

	if s.Schema.constVal != nil {
		if constFloat, ok := s.Schema.constVal.(float32); ok && constFloat != floatValue {
			message := floatConstError(constFloat)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(floatValue, message, "const"))
		}
	}

	return ParseResult{Valid: len(errors) == 0, Value: finalValue, Errors: errors}
}

func (s *FloatSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("number")
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

	schema["format"] = "float"

	if s.nullable {
		schema["type"] = []string{"number", "null"}
	}

	return schema
}

func (s *FloatSchema) MarshalJSON() ([]byte, error) {
	type jsonFloatSchema struct {
		Schema
		Minimum    *float32 `json:"minimum,omitempty"`
		Maximum    *float32 `json:"maximum,omitempty"`
		MultipleOf *float32 `json:"multipleOf,omitempty"`
		Format     string   `json:"format"`
		Nullable   bool     `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonFloatSchema{
		Schema:     s.Schema,
		Minimum:    s.minimum,
		Maximum:    s.maximum,
		MultipleOf: s.multipleOf,
		Format:     "float",
		Nullable:   s.nullable,
	})
}
