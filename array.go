package schema

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nyxstack/i18n"
)

// Default error messages for array validation
var (
	arrayRequiredError = i18n.S("value is required")
	arrayTypeError     = i18n.S("value must be an array")
	arrayUniqueError   = i18n.S("array must contain unique items")
)

func arrayMinItemsError(min int) i18n.TranslatedFunc {
	return i18n.F("array must contain at least %d items", min)
}

func arrayMaxItemsError(max int) i18n.TranslatedFunc {
	return i18n.F("array must contain at most %d items", max)
}

func arrayItemError(index int) i18n.TranslatedFunc {
	return i18n.F("array item at index %d is invalid", index)
}

// ArraySchema represents a JSON Schema for array values
type ArraySchema struct {
	Schema
	// Array-specific validation
	itemSchema  Parseable // Schema for validating items
	minItems    *int      // Minimum number of items
	maxItems    *int      // Maximum number of items
	uniqueItems bool      // Items must be unique
	nullable    bool      // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	minItemsError     ErrorMessage
	maxItemsError     ErrorMessage
	uniqueItemsError  ErrorMessage
	itemError         ErrorMessage
	typeMismatchError ErrorMessage
}

// Array creates a new array schema with an item schema
func Array(itemSchema Parseable, errorMessage ...interface{}) *ArraySchema {
	schema := &ArraySchema{
		Schema: Schema{
			schemaType: "array",
			required:   true, // Default to required
		},
		itemSchema: itemSchema,
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *ArraySchema) Title(title string) *ArraySchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *ArraySchema) Description(description string) *ArraySchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *ArraySchema) Default(value interface{}) *ArraySchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *ArraySchema) Example(example []interface{}) *ArraySchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Array-specific validation

// Items sets the schema for array items
func (s *ArraySchema) Items(itemSchema Parseable) *ArraySchema {
	s.itemSchema = itemSchema
	return s
}

// MinItems sets the minimum number of items with optional custom error message
func (s *ArraySchema) MinItems(min int, errorMessage ...interface{}) *ArraySchema {
	s.minItems = &min
	if len(errorMessage) > 0 {
		s.minItemsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MaxItems sets the maximum number of items with optional custom error message
func (s *ArraySchema) MaxItems(max int, errorMessage ...interface{}) *ArraySchema {
	s.maxItems = &max
	if len(errorMessage) > 0 {
		s.maxItemsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Length sets both min and max items to the same value
func (s *ArraySchema) Length(length int) *ArraySchema {
	s.minItems = &length
	s.maxItems = &length
	return s
}

// UniqueItems requires all items to be unique with optional custom error message
func (s *ArraySchema) UniqueItems(errorMessage ...interface{}) *ArraySchema {
	s.uniqueItems = true
	if len(errorMessage) > 0 {
		s.uniqueItemsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *ArraySchema) Optional() *ArraySchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *ArraySchema) Required(errorMessage ...interface{}) *ArraySchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *ArraySchema) Nullable() *ArraySchema {
	s.nullable = true
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *ArraySchema) TypeError(message string) *ArraySchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// ItemError sets a custom error message for item validation failures
func (s *ArraySchema) ItemError(message string) *ArraySchema {
	s.itemError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *ArraySchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *ArraySchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *ArraySchema) IsNullable() bool {
	return s.nullable
}

// GetItemSchema returns the schema for array items
func (s *ArraySchema) GetItemSchema() Parseable {
	return s.itemSchema
}

// GetMinItems returns the minimum number of items
func (s *ArraySchema) GetMinItems() *int {
	return s.minItems
}

// GetMaxItems returns the maximum number of items
func (s *ArraySchema) GetMaxItems() *int {
	return s.maxItems
}

// IsUniqueItems returns whether items must be unique
func (s *ArraySchema) IsUniqueItems() bool {
	return s.uniqueItems
}

// Validation helpers

// isUnique checks if all items in a slice are unique
func isUnique(slice []interface{}) bool {
	seen := make(map[interface{}]bool)
	for _, item := range slice {
		key := getComparableKey(item)
		if seen[key] {
			return false
		}
		seen[key] = true
	}
	return true
}

// getComparableKey converts an interface{} to a comparable key
func getComparableKey(item interface{}) interface{} {
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

// Parse validates and parses an array value, returning the final parsed value
func (s *ArraySchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := arrayRequiredError(ctx.Locale)
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
	var arrayValue []interface{}
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		message := arrayTypeError(ctx.Locale)
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
	arrayValue = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		arrayValue[i] = v.Index(i).Interface()
	}

	// Now validate the array against all constraints
	finalValue := make([]interface{}, len(arrayValue)) // This will be our parsed array

	// Validate length constraints
	length := len(arrayValue)
	if s.minItems != nil && length < *s.minItems {
		message := arrayMinItemsError(*s.minItems)(ctx.Locale)
		if !isEmptyErrorMessage(s.minItemsError) {
			message = resolveErrorMessage(s.minItemsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(arrayValue, message, "min_items"))
	}

	if s.maxItems != nil && length > *s.maxItems {
		message := arrayMaxItemsError(*s.maxItems)(ctx.Locale)
		if !isEmptyErrorMessage(s.maxItemsError) {
			message = resolveErrorMessage(s.maxItemsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(arrayValue, message, "max_items"))
	}

	// Validate each item using the item schema
	for i, item := range arrayValue {
		if s.itemSchema != nil {
			itemResult := s.itemSchema.Parse(item, ctx)
			if !itemResult.Valid {
				// Create error for this item
				message := arrayItemError(i)(ctx.Locale)
				if !isEmptyErrorMessage(s.itemError) {
					message = resolveErrorMessage(s.itemError, ctx)
				}
				// Add the main item error
				errors = append(errors, NewFieldError([]string{fmt.Sprintf("[%d]", i)}, item, message, "item_invalid"))
				// Also add the specific validation errors for this item
				for _, itemErr := range itemResult.Errors {
					// Prefix the path with array index
					errors = append(errors, NewFieldError(append([]string{fmt.Sprintf("[%d]", i)}, itemErr.Path...), itemErr.Value, itemErr.Message, itemErr.Code))
				}
			} else {
				// Use the parsed value from item validation
				finalValue[i] = itemResult.Value
			}
		} else {
			// No item schema, use original value
			finalValue[i] = item
		}
	}

	// Check uniqueness constraint
	if s.uniqueItems && !isUnique(arrayValue) {
		message := arrayUniqueError(ctx.Locale)
		if !isEmptyErrorMessage(s.uniqueItemsError) {
			message = resolveErrorMessage(s.uniqueItemsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(arrayValue, message, "unique_items"))
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *ArraySchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("array")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add array-specific fields
	if s.itemSchema != nil {
		if jsonSchema, ok := s.itemSchema.(interface{ JSON() map[string]interface{} }); ok {
			schema["items"] = jsonSchema.JSON()
		}
	}

	if s.minItems != nil {
		schema["minItems"] = *s.minItems
	}

	if s.maxItems != nil {
		schema["maxItems"] = *s.maxItems
	}

	if s.uniqueItems {
		schema["uniqueItems"] = true
	}

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"array", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize ArraySchema for JSON schema generation
func (s *ArraySchema) MarshalJSON() ([]byte, error) {
	type jsonArraySchema struct {
		Schema
		ItemSchema  Parseable `json:"itemSchema,omitempty"`
		MinItems    *int      `json:"minItems,omitempty"`
		MaxItems    *int      `json:"maxItems,omitempty"`
		UniqueItems bool      `json:"uniqueItems,omitempty"`
		Nullable    bool      `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonArraySchema{
		Schema:      s.Schema,
		ItemSchema:  s.itemSchema,
		MinItems:    s.minItems,
		MaxItems:    s.maxItems,
		UniqueItems: s.uniqueItems,
		Nullable:    s.nullable,
	})
}

// Interface implementations for ArraySchema

// SetTitle implements SetTitle interface
func (s *ArraySchema) SetTitle(title string) {
	s.Title(title)
}

// SetDescription implements SetDescription interface
func (s *ArraySchema) SetDescription(description string) {
	s.Description(description)
}

// SetRequired implements SetRequired interface
func (s *ArraySchema) SetRequired() {
	s.Required()
}

// SetOptional implements SetOptional interface
func (s *ArraySchema) SetOptional() {
	s.Optional()
}

// SetMinItems implements SetMinItems interface
func (s *ArraySchema) SetMinItems(min int) {
	s.MinItems(min)
}

// SetMaxItems implements SetMaxItems interface
func (s *ArraySchema) SetMaxItems(max int) {
	s.MaxItems(max)
}

// SetNullable implements SetNullable interface
func (s *ArraySchema) SetNullable() {
	s.Nullable()
}

// SetDefault implements SetDefault interface
func (s *ArraySchema) SetDefault(value interface{}) {
	s.Default(value)
}

// SetExample implements SetExample interface
func (s *ArraySchema) SetExample(example interface{}) {
	if val, ok := example.([]interface{}); ok {
		s.Example(val)
	}
}
