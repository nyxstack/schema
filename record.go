package schema

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nyxstack/i18n"
)

// Default error messages for record validation
var (
	recordRequiredError = i18n.S("value is required")
	recordTypeError     = i18n.S("value must be an object")
	recordKeyError      = i18n.S("record key is invalid")
	recordValueError    = i18n.S("record value is invalid")
)

func recordMinPropsError(min int) i18n.TranslatedFunc {
	return i18n.F("record must contain at least %d properties", min)
}

func recordMaxPropsError(max int) i18n.TranslatedFunc {
	return i18n.F("record must contain at most %d properties", max)
}

// RecordSchema represents a JSON Schema for key-value record/map validation
// This is similar to additionalProperties in JSON Schema
type RecordSchema struct {
	Schema
	// Record-specific validation
	keySchema   Parseable // Schema for validating keys (usually string schema)
	valueSchema Parseable // Schema for validating values
	minProps    *int      // Minimum number of properties
	maxProps    *int      // Maximum number of properties
	nullable    bool      // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	minPropsError     ErrorMessage
	maxPropsError     ErrorMessage
	keyError          ErrorMessage
	valueError        ErrorMessage
	typeMismatchError ErrorMessage
}

// Record creates a new record schema with key and value schemas
func Record(keySchema, valueSchema Parseable, errorMessage ...interface{}) *RecordSchema {
	schema := &RecordSchema{
		Schema: Schema{
			schemaType: "object", // Records are objects in JSON Schema
			required:   true,     // Default to required
		},
		keySchema:   keySchema,
		valueSchema: valueSchema,
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *RecordSchema) Title(title string) *RecordSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *RecordSchema) Description(description string) *RecordSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *RecordSchema) Default(value interface{}) *RecordSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *RecordSchema) Example(example map[string]interface{}) *RecordSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Record-specific validation

// Keys sets the schema for record keys
func (s *RecordSchema) Keys(keySchema Parseable) *RecordSchema {
	s.keySchema = keySchema
	return s
}

// Values sets the schema for record values
func (s *RecordSchema) Values(valueSchema Parseable) *RecordSchema {
	s.valueSchema = valueSchema
	return s
}

// MinProperties sets the minimum number of properties with optional custom error message
func (s *RecordSchema) MinProperties(min int, errorMessage ...interface{}) *RecordSchema {
	s.minProps = &min
	if len(errorMessage) > 0 {
		s.minPropsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MaxProperties sets the maximum number of properties with optional custom error message
func (s *RecordSchema) MaxProperties(max int, errorMessage ...interface{}) *RecordSchema {
	s.maxProps = &max
	if len(errorMessage) > 0 {
		s.maxPropsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Size sets both min and max properties to the same value
func (s *RecordSchema) Size(size int) *RecordSchema {
	s.minProps = &size
	s.maxProps = &size
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *RecordSchema) Optional() *RecordSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *RecordSchema) Required(errorMessage ...interface{}) *RecordSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *RecordSchema) Nullable() *RecordSchema {
	s.nullable = true
	return s
}

// Error customization

// TypeError sets a custom error message for type mismatch validation
func (s *RecordSchema) TypeError(message string) *RecordSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// KeyError sets a custom error message for key validation failures
func (s *RecordSchema) KeyError(message string) *RecordSchema {
	s.keyError = toErrorMessage(message)
	return s
}

// ValueError sets a custom error message for value validation failures
func (s *RecordSchema) ValueError(message string) *RecordSchema {
	s.valueError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *RecordSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *RecordSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *RecordSchema) IsNullable() bool {
	return s.nullable
}

// GetKeySchema returns the schema for record keys
func (s *RecordSchema) GetKeySchema() Parseable {
	return s.keySchema
}

// GetValueSchema returns the schema for record values
func (s *RecordSchema) GetValueSchema() Parseable {
	return s.valueSchema
}

// GetMinProperties returns the minimum number of properties
func (s *RecordSchema) GetMinProperties() *int {
	return s.minProps
}

// GetMaxProperties returns the maximum number of properties
func (s *RecordSchema) GetMaxProperties() *int {
	return s.maxProps
}

// Validation

// Parse validates and parses a record value, returning the final parsed value
func (s *RecordSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := recordRequiredError(ctx.Locale)
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

	// Type check - accept map or struct
	var recordMap map[string]interface{}
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Map:
		// Convert map to map[string]interface{}
		recordMap = make(map[string]interface{})
		for _, key := range v.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			recordMap[keyStr] = v.MapIndex(key).Interface()
		}
	case reflect.Struct:
		// Convert struct to map[string]interface{}
		recordMap = make(map[string]interface{})
		structType := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := structType.Field(i)
			if field.IsExported() { // Only exported fields
				recordMap[field.Name] = v.Field(i).Interface()
			}
		}
	default:
		message := recordTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Now validate the record against all constraints
	finalValue := make(map[string]interface{}, len(recordMap)) // This will be our parsed record

	// Validate size constraints
	size := len(recordMap)
	if s.minProps != nil && size < *s.minProps {
		message := recordMinPropsError(*s.minProps)(ctx.Locale)
		if !isEmptyErrorMessage(s.minPropsError) {
			message = resolveErrorMessage(s.minPropsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(recordMap, message, "min_properties"))
	}

	if s.maxProps != nil && size > *s.maxProps {
		message := recordMaxPropsError(*s.maxProps)(ctx.Locale)
		if !isEmptyErrorMessage(s.maxPropsError) {
			message = resolveErrorMessage(s.maxPropsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(recordMap, message, "max_properties"))
	}

	// Validate each key-value pair
	for key, val := range recordMap {
		var finalKey string = key
		var finalVal interface{} = val

		// Validate key using key schema
		if s.keySchema != nil {
			keyResult := s.keySchema.Parse(key, ctx)
			if !keyResult.Valid {
				// Key validation failed
				message := recordKeyError(ctx.Locale)
				if !isEmptyErrorMessage(s.keyError) {
					message = resolveErrorMessage(s.keyError, ctx)
				}
				errors = append(errors, NewFieldError([]string{key}, key, message, "key_invalid"))
				// Also add the specific key validation errors
				for _, keyErr := range keyResult.Errors {
					errors = append(errors, NewFieldError([]string{key + "_key"}, keyErr.Value, keyErr.Message, keyErr.Code))
				}
				continue // Skip this key-value pair
			} else {
				// Use the parsed key
				if parsedKey, ok := keyResult.Value.(string); ok {
					finalKey = parsedKey
				}
			}
		}

		// Validate value using value schema
		if s.valueSchema != nil {
			valueResult := s.valueSchema.Parse(val, ctx)
			if !valueResult.Valid {
				// Value validation failed
				message := recordValueError(ctx.Locale)
				if !isEmptyErrorMessage(s.valueError) {
					message = resolveErrorMessage(s.valueError, ctx)
				}
				errors = append(errors, NewFieldError([]string{key}, val, message, "value_invalid"))
				// Also add the specific value validation errors
				for _, valErr := range valueResult.Errors {
					// Prefix the path with the key
					errors = append(errors, NewFieldError(append([]string{key}, valErr.Path...), valErr.Value, valErr.Message, valErr.Code))
				}
			} else {
				// Use the parsed value
				finalVal = valueResult.Value
			}
		}

		// Store the final key-value pair
		finalValue[finalKey] = finalVal
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *RecordSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("object")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// For records, we use additionalProperties to represent value schema
	if s.valueSchema != nil {
		if jsonSchema, ok := s.valueSchema.(interface{ JSON() map[string]interface{} }); ok {
			schema["additionalProperties"] = jsonSchema.JSON()
		} else {
			schema["additionalProperties"] = true
		}
	} else {
		schema["additionalProperties"] = true
	}

	// Add property count constraints
	if s.minProps != nil {
		schema["minProperties"] = *s.minProps
	}

	if s.maxProps != nil {
		schema["maxProperties"] = *s.maxProps
	}

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"object", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize RecordSchema for JSON schema generation
func (s *RecordSchema) MarshalJSON() ([]byte, error) {
	type jsonRecordSchema struct {
		Schema
		KeySchema   Parseable `json:"keySchema,omitempty"`
		ValueSchema Parseable `json:"valueSchema,omitempty"`
		MinProps    *int      `json:"minProps,omitempty"`
		MaxProps    *int      `json:"maxProps,omitempty"`
		Nullable    bool      `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonRecordSchema{
		Schema:      s.Schema,
		KeySchema:   s.keySchema,
		ValueSchema: s.valueSchema,
		MinProps:    s.minProps,
		MaxProps:    s.maxProps,
		Nullable:    s.nullable,
	})
}
