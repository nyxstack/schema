package schema

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nyxstack/i18n"
)

// Default error messages for tuple validation
var (
	tupleRequiredError = i18n.S("value is required")
	tupleTypeError     = i18n.S("value must be an array")
	tupleUniqueError   = i18n.S("tuple items must be unique")
)

func tupleLengthError(expected int) i18n.TranslatedFunc {
	return i18n.F("tuple must have exactly %d items", expected)
}

func tupleMinLengthError(min int) i18n.TranslatedFunc {
	return i18n.F("tuple must have at least %d items", min)
}

func tupleItemError(index int) i18n.TranslatedFunc {
	return i18n.F("tuple item at index %d is invalid", index)
}

// TupleSchema represents a JSON Schema for fixed-length arrays with position-specific types
type TupleSchema struct {
	Schema
	// Tuple-specific validation
	itemSchemas     []Parseable // Schemas for each position (order matters)
	additionalItems bool        // Allow additional items beyond defined positions
	uniqueItems     bool        // Items must be unique
	nullable        bool        // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	lengthError       ErrorMessage
	uniqueItemsError  ErrorMessage
	itemError         ErrorMessage
	typeMismatchError ErrorMessage
}

// Tuple creates a new tuple schema with position-specific item schemas
func Tuple(itemSchemas ...Parseable) *TupleSchema {
	schema := &TupleSchema{
		Schema: Schema{
			schemaType: "array",
			required:   true, // Default to required
		},
		itemSchemas:     itemSchemas,
		additionalItems: false, // Strict by default - exact length required
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *TupleSchema) Title(title string) *TupleSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *TupleSchema) Description(description string) *TupleSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *TupleSchema) Default(value interface{}) *TupleSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *TupleSchema) Example(example []interface{}) *TupleSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Tuple-specific validation

// AllowAdditionalItems allows extra items beyond the defined positions
func (s *TupleSchema) AllowAdditionalItems() *TupleSchema {
	s.additionalItems = true
	return s
}

// Strict requires exact length matching (default behavior)
func (s *TupleSchema) Strict() *TupleSchema {
	s.additionalItems = false
	return s
}

// UniqueItems requires all items to be unique with optional custom error message
func (s *TupleSchema) UniqueItems(errorMessage ...interface{}) *TupleSchema {
	s.uniqueItems = true
	if len(errorMessage) > 0 {
		s.uniqueItemsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *TupleSchema) Optional() *TupleSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *TupleSchema) Required(errorMessage ...interface{}) *TupleSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *TupleSchema) Nullable() *TupleSchema {
	s.nullable = true
	return s
}

// Error customization

// TypeError sets a custom error message for type mismatch validation
func (s *TupleSchema) TypeError(message string) *TupleSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// LengthError sets a custom error message for length validation
func (s *TupleSchema) LengthError(message string) *TupleSchema {
	s.lengthError = toErrorMessage(message)
	return s
}

// ItemError sets a custom error message for item validation failures
func (s *TupleSchema) ItemError(message string) *TupleSchema {
	s.itemError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *TupleSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *TupleSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *TupleSchema) IsNullable() bool {
	return s.nullable
}

// GetItemSchemas returns the schemas for each tuple position
func (s *TupleSchema) GetItemSchemas() []Parseable {
	return s.itemSchemas
}

// GetExpectedLength returns the expected tuple length
func (s *TupleSchema) GetExpectedLength() int {
	return len(s.itemSchemas)
}

// AllowsAdditionalItems returns whether additional items are allowed
func (s *TupleSchema) AllowsAdditionalItems() bool {
	return s.additionalItems
}

// IsUniqueItems returns whether items must be unique
func (s *TupleSchema) IsUniqueItems() bool {
	return s.uniqueItems
}

// Validation helpers

// isUnique checks if all items in a slice are unique
func isTupleUnique(slice []interface{}) bool {
	seen := make(map[interface{}]bool)
	for _, item := range slice {
		key := getTupleComparableKey(item)
		if seen[key] {
			return false
		}
		seen[key] = true
	}
	return true
}

// getTupleComparableKey converts an interface{} to a comparable key
func getTupleComparableKey(item interface{}) interface{} {
	if item == nil {
		return nil
	}

	v := reflect.ValueOf(item)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Func:
		// These types aren't directly comparable, use their string representation
		return v.String()
	default:
		return item
	}
}

// Validation

// Parse validates and parses a tuple value, returning the final parsed value
func (s *TupleSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := tupleRequiredError(ctx.Locale)
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

	// Type check - convert to slice
	var tupleValue []interface{}
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		message := tupleTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Convert to []interface{}
	tupleValue = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		tupleValue[i] = v.Index(i).Interface()
	}

	// Validate length constraints
	actualLength := len(tupleValue)
	expectedLength := len(s.itemSchemas)

	if !s.additionalItems && actualLength != expectedLength {
		message := tupleLengthError(expectedLength)(ctx.Locale)
		if !isEmptyErrorMessage(s.lengthError) {
			message = resolveErrorMessage(s.lengthError, ctx)
		}
		errors = append(errors, NewPrimitiveError(tupleValue, message, "tuple_length"))
	}

	if s.additionalItems && actualLength < expectedLength {
		message := tupleMinLengthError(expectedLength)(ctx.Locale)
		if !isEmptyErrorMessage(s.lengthError) {
			message = resolveErrorMessage(s.lengthError, ctx)
		}
		errors = append(errors, NewPrimitiveError(tupleValue, message, "min_length"))
	}

	// Prepare final value array
	finalValue := make([]interface{}, len(tupleValue))

	// Validate each item at its position using the corresponding schema
	for i, item := range tupleValue {
		if i < len(s.itemSchemas) {
			// Validate using position-specific schema
			itemResult := s.itemSchemas[i].Parse(item, ctx)
			if !itemResult.Valid {
				// Create error for this item
				message := tupleItemError(i)(ctx.Locale)
				if !isEmptyErrorMessage(s.itemError) {
					message = resolveErrorMessage(s.itemError, ctx)
				}
				// Add the main item error
				errors = append(errors, NewFieldError([]string{fmt.Sprintf("[%d]", i)}, item, message, "item_invalid"))
				// Also add the specific validation errors for this item
				for _, itemErr := range itemResult.Errors {
					// Prefix the path with tuple index
					errors = append(errors, NewFieldError(append([]string{fmt.Sprintf("[%d]", i)}, itemErr.Path...), itemErr.Value, itemErr.Message, itemErr.Code))
				}
			} else {
				// Use the parsed value from item validation
				finalValue[i] = itemResult.Value
			}
		} else if s.additionalItems {
			// Additional items beyond defined positions - accept as-is
			finalValue[i] = item
		}
	}

	// Check uniqueness constraint
	if s.uniqueItems && !isTupleUnique(tupleValue) {
		message := tupleUniqueError(ctx.Locale)
		if !isEmptyErrorMessage(s.uniqueItemsError) {
			message = resolveErrorMessage(s.uniqueItemsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(tupleValue, message, "unique_items"))
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *TupleSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("array")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add tuple-specific fields using "items" as array of schemas
	if len(s.itemSchemas) > 0 {
		items := make([]interface{}, len(s.itemSchemas))
		for i, itemSchema := range s.itemSchemas {
			if jsonSchema, ok := itemSchema.(interface{ JSON() map[string]interface{} }); ok {
				items[i] = jsonSchema.JSON()
			} else {
				items[i] = map[string]interface{}{"type": "unknown"}
			}
		}
		schema["items"] = items
	}

	// Add additionalItems
	schema["additionalItems"] = s.additionalItems

	// Add uniqueItems if true
	if s.uniqueItems {
		schema["uniqueItems"] = true
	}

	// Set exact length constraints for strict tuples
	if !s.additionalItems && len(s.itemSchemas) > 0 {
		schema["minItems"] = len(s.itemSchemas)
		schema["maxItems"] = len(s.itemSchemas)
	} else if len(s.itemSchemas) > 0 {
		schema["minItems"] = len(s.itemSchemas)
	}

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"array", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize TupleSchema for JSON schema generation
func (s *TupleSchema) MarshalJSON() ([]byte, error) {
	type jsonTupleSchema struct {
		Schema
		ItemSchemas     []Parseable `json:"itemSchemas"`
		AdditionalItems bool        `json:"additionalItems"`
		UniqueItems     bool        `json:"uniqueItems,omitempty"`
		Nullable        bool        `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonTupleSchema{
		Schema:          s.Schema,
		ItemSchemas:     s.itemSchemas,
		AdditionalItems: s.additionalItems,
		UniqueItems:     s.uniqueItems,
		Nullable:        s.nullable,
	})
}
