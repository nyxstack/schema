package schema

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nyxstack/i18n"
)

// Default error messages for object validation
var (
	objectRequiredError        = i18n.S("value is required")
	objectTypeError            = i18n.S("value must be an object")
	objectAdditionalPropsError = i18n.S("additional property is not allowed")
)

func objectMinPropsError(min int) i18n.TranslatedFunc {
	return i18n.F("object must have at least %d properties", min)
}

func objectMaxPropsError(max int) i18n.TranslatedFunc {
	return i18n.F("object must have at most %d properties", max)
}

func objectPropertyError(prop string) i18n.TranslatedFunc {
	return i18n.F("property %s is invalid", prop)
}

func objectRequiredPropError(prop string) i18n.TranslatedFunc {
	return i18n.F("property %s is required", prop)
}

// Shape represents a map of property names to their schemas for object construction
type Shape map[string]interface{}

// AsObject converts a Shape to an ObjectSchema
func (s Shape) AsObject() *ObjectSchema {
	return Object(s)
}

// ObjectProperty represents a single property in an object schema
type ObjectProperty struct {
	Schema   Parseable // The schema validator for this property
	Required bool      // Whether this property is required
	Name     string    // The property name
}

// ObjectSchema represents a JSON Schema for object values with structured properties
type ObjectSchema struct {
	Schema
	// Object-specific validation
	properties      map[string]ObjectProperty // Property schemas
	requiredProps   []string                  // List of required property names
	additionalProps bool                      // Allow additional properties
	minProps        *int                      // Minimum number of properties
	maxProps        *int                      // Maximum number of properties
	nullable        bool                      // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError        ErrorMessage
	minPropsError        ErrorMessage
	maxPropsError        ErrorMessage
	additionalPropsError ErrorMessage
	propertyError        ErrorMessage
	typeMismatchError    ErrorMessage
}

// Object creates a new object schema with optional Shape and error message
func Object(shapeAndError ...interface{}) *ObjectSchema {
	schema := &ObjectSchema{
		Schema: Schema{
			schemaType: "object",
			required:   true, // Default to required
		},
		properties:      make(map[string]ObjectProperty),
		requiredProps:   []string{},
		additionalProps: false, // Strict by default
	}

	// Process optional parameters
	for _, param := range shapeAndError {
		switch p := param.(type) {
		case Shape:
			// Add properties from the shape (determine required from schema)
			for name, schemaVal := range p {
				schema.Property(name, schemaVal)
			}
		case string:
			// Set custom type mismatch error message
			schema.typeMismatchError = toErrorMessage(p)
		}
	}

	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *ObjectSchema) Title(title string) *ObjectSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *ObjectSchema) Description(description string) *ObjectSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *ObjectSchema) Default(value interface{}) *ObjectSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *ObjectSchema) Example(example map[string]interface{}) *ObjectSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Property definition methods

// Property adds a property to the object schema (infers required/optional from schema)
func (s *ObjectSchema) Property(name string, schema interface{}) *ObjectSchema {
	// Convert to Parseable interface
	var parseable Parseable
	if p, ok := schema.(Parseable); ok {
		parseable = p
	} else {
		// Try to wrap in a simple interface - this is a fallback
		return s // Skip if not Parseable
	}

	// Check if the schema has IsRequired() method to determine if it's required
	isRequired := true
	if requiredChecker, ok := schema.(interface{ IsRequired() bool }); ok {
		isRequired = requiredChecker.IsRequired()
	}

	s.properties[name] = ObjectProperty{
		Schema:   parseable,
		Required: isRequired,
		Name:     name,
	}

	// Add to required list if required and not already there
	if isRequired {
		for _, req := range s.requiredProps {
			if req == name {
				return s
			}
		}
		s.requiredProps = append(s.requiredProps, name)
	}
	return s
}

// OptionalProperty explicitly adds an optional property
func (s *ObjectSchema) OptionalProperty(name string, schema interface{}) *ObjectSchema {
	var parseable Parseable
	if p, ok := schema.(Parseable); ok {
		parseable = p
	} else {
		return s // Skip if not Parseable
	}

	s.properties[name] = ObjectProperty{
		Schema:   parseable,
		Required: false,
		Name:     name,
	}
	return s
}

// RequiredProperty explicitly adds a required property
func (s *ObjectSchema) RequiredProperty(name string, schema interface{}) *ObjectSchema {
	var parseable Parseable
	if p, ok := schema.(Parseable); ok {
		parseable = p
	} else {
		return s // Skip if not Parseable
	}

	s.properties[name] = ObjectProperty{
		Schema:   parseable,
		Required: true,
		Name:     name,
	}

	// Add to required list if not already there
	for _, req := range s.requiredProps {
		if req == name {
			return s
		}
	}
	s.requiredProps = append(s.requiredProps, name)
	return s
}

// Object constraint methods

// MinProperties sets the minimum number of properties with optional custom error message
func (s *ObjectSchema) MinProperties(min int, errorMessage ...interface{}) *ObjectSchema {
	s.minProps = &min
	if len(errorMessage) > 0 {
		s.minPropsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MaxProperties sets the maximum number of properties with optional custom error message
func (s *ObjectSchema) MaxProperties(max int, errorMessage ...interface{}) *ObjectSchema {
	s.maxProps = &max
	if len(errorMessage) > 0 {
		s.maxPropsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// PropertyRange sets both min and max property constraints
func (s *ObjectSchema) PropertyRange(min, max int, errorMessage ...interface{}) *ObjectSchema {
	s.minProps = &min
	s.maxProps = &max
	if len(errorMessage) > 0 {
		s.minPropsError = toErrorMessage(errorMessage[0])
		s.maxPropsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Strict disallows additional properties (default behavior)
func (s *ObjectSchema) Strict() *ObjectSchema {
	s.additionalProps = false
	return s
}

// Passthrough allows additional properties
func (s *ObjectSchema) Passthrough() *ObjectSchema {
	s.additionalProps = true
	return s
}

// AdditionalProperties sets whether additional properties are allowed with optional custom error message
func (s *ObjectSchema) AdditionalProperties(allowed bool, errorMessage ...interface{}) *ObjectSchema {
	s.additionalProps = allowed
	if !allowed && len(errorMessage) > 0 {
		s.additionalPropsError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *ObjectSchema) Optional() *ObjectSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *ObjectSchema) Required(errorMessage ...interface{}) *ObjectSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *ObjectSchema) Nullable() *ObjectSchema {
	s.nullable = true
	return s
}

// Error customization

// TypeError sets a custom error message for type mismatch validation
func (s *ObjectSchema) TypeError(message string) *ObjectSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// PropertyError sets a custom error prefix for property validation errors
func (s *ObjectSchema) PropertyError(message string) *ObjectSchema {
	s.propertyError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *ObjectSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *ObjectSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *ObjectSchema) IsNullable() bool {
	return s.nullable
}

// GetProperties returns the object properties
func (s *ObjectSchema) GetProperties() map[string]ObjectProperty {
	return s.properties
}

// GetRequiredProperties returns the list of required property names
func (s *ObjectSchema) GetRequiredProperties() []string {
	return s.requiredProps
}

// AllowsAdditionalProperties returns whether additional properties are allowed
func (s *ObjectSchema) AllowsAdditionalProperties() bool {
	return s.additionalProps
}

// GetMinProperties returns the minimum number of properties
func (s *ObjectSchema) GetMinProperties() *int {
	return s.minProps
}

// GetMaxProperties returns the maximum number of properties
func (s *ObjectSchema) GetMaxProperties() *int {
	return s.maxProps
}

// Helper methods for converting input to map[string]interface{}

// convertToMap converts various input types to map[string]interface{}
func convertToMap(value interface{}) (map[string]interface{}, bool) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		// Convert map to map[string]interface{}
		result := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			result[keyStr] = v.MapIndex(key).Interface()
		}
		return result, true

	case reflect.Struct:
		// Convert struct to map[string]interface{}
		result := make(map[string]interface{})
		structType := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := structType.Field(i)
			if !field.IsExported() {
				continue // Skip unexported fields
			}

			// Use json tag if available, otherwise use field name
			fieldName := field.Name
			if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
				// Handle "fieldname,omitempty" format
				commaIdx := len(tag)
				for idx := 0; idx < len(tag); idx++ {
					if tag[idx] == ',' {
						commaIdx = idx
						break
					}
				}
				fieldName = tag[:commaIdx]
			}

			result[fieldName] = v.Field(i).Interface()
		}
		return result, true

	default:
		return nil, false
	}
}

// Validation

// Parse validates and parses an object value, returning the final parsed value
func (s *ObjectSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := objectRequiredError(ctx.Locale)
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

	// Type check and convert to map
	objectMap, ok := convertToMap(value)
	if !ok {
		message := objectTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Now validate the object against all constraints
	finalValue := make(map[string]interface{}, len(objectMap)) // This will be our parsed object

	// Validate property count constraints
	propCount := len(objectMap)
	if s.minProps != nil && propCount < *s.minProps {
		message := objectMinPropsError(*s.minProps)(ctx.Locale)
		if !isEmptyErrorMessage(s.minPropsError) {
			message = resolveErrorMessage(s.minPropsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(objectMap, message, "min_properties"))
	}

	if s.maxProps != nil && propCount > *s.maxProps {
		message := objectMaxPropsError(*s.maxProps)(ctx.Locale)
		if !isEmptyErrorMessage(s.maxPropsError) {
			message = resolveErrorMessage(s.maxPropsError, ctx)
		}
		errors = append(errors, NewPrimitiveError(objectMap, message, "max_properties"))
	}

	// Check required properties
	for _, requiredProp := range s.requiredProps {
		if _, exists := objectMap[requiredProp]; !exists {
			message := objectRequiredPropError(requiredProp)(ctx.Locale)
			errors = append(errors, NewFieldError([]string{requiredProp}, "<missing>", message, "required"))
		}
	}

	// Validate each property
	for propName, propValue := range objectMap {
		// Check if property is defined in schema
		propSchema, isDefined := s.properties[propName]
		if !isDefined {
			if !s.additionalProps {
				message := objectAdditionalPropsError(ctx.Locale)
				if !isEmptyErrorMessage(s.additionalPropsError) {
					message = resolveErrorMessage(s.additionalPropsError, ctx)
				}
				errors = append(errors, NewFieldError([]string{propName}, propValue, message, "additional_property"))
			} else {
				// Additional property allowed, use as-is
				finalValue[propName] = propValue
			}
			continue
		}

		// Validate the property value using its schema
		propResult := propSchema.Schema.Parse(propValue, ctx)
		if !propResult.Valid {
			// Property validation failed
			message := objectPropertyError(propName)(ctx.Locale)
			if !isEmptyErrorMessage(s.propertyError) {
				message = resolveErrorMessage(s.propertyError, ctx)
			}
			// Add the main property error
			errors = append(errors, NewFieldError([]string{propName}, propValue, message, "property_invalid"))
			// Also add the specific validation errors for this property
			for _, propErr := range propResult.Errors {
				// Prefix the path with property name
				errors = append(errors, NewFieldError(append([]string{propName}, propErr.Path...), propErr.Value, propErr.Message, propErr.Code))
			}
		} else {
			// Use the parsed value from property validation
			finalValue[propName] = propResult.Value
		}
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *ObjectSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("object")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add object-specific fields
	if len(s.properties) > 0 {
		properties := make(map[string]interface{})
		for name, prop := range s.properties {
			if jsonSchema, ok := prop.Schema.(interface{ JSON() map[string]interface{} }); ok {
				properties[name] = jsonSchema.JSON()
			}
		}
		schema["properties"] = properties
	}

	if len(s.requiredProps) > 0 {
		schema["required"] = s.requiredProps
	}

	schema["additionalProperties"] = s.additionalProps

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

// MarshalJSON implements json.Marshaler to properly serialize ObjectSchema for JSON schema generation
func (s *ObjectSchema) MarshalJSON() ([]byte, error) {
	type jsonObjectSchema struct {
		Schema
		Properties      map[string]ObjectProperty `json:"properties"`
		RequiredProps   []string                  `json:"required,omitempty"`
		AdditionalProps bool                      `json:"additionalProperties"`
		MinProps        *int                      `json:"minProperties,omitempty"`
		MaxProps        *int                      `json:"maxProperties,omitempty"`
		Nullable        bool                      `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonObjectSchema{
		Schema:          s.Schema,
		Properties:      s.properties,
		RequiredProps:   s.requiredProps,
		AdditionalProps: s.additionalProps,
		MinProps:        s.minProps,
		MaxProps:        s.maxProps,
		Nullable:        s.nullable,
	})
}

// Interface implementations for ObjectSchema

// SetTitle implements SetTitle interface
func (s *ObjectSchema) SetTitle(title string) {
	s.Title(title)
}

// SetDescription implements SetDescription interface
func (s *ObjectSchema) SetDescription(description string) {
	s.Description(description)
}

// SetRequired implements SetRequired interface
func (s *ObjectSchema) SetRequired() {
	s.Required()
}

// SetOptional implements SetOptional interface
func (s *ObjectSchema) SetOptional() {
	s.Optional()
}

// SetNullable implements SetNullable interface
func (s *ObjectSchema) SetNullable() {
	s.Nullable()
}

// SetDefault implements SetDefault interface
func (s *ObjectSchema) SetDefault(value interface{}) {
	s.Default(value)
}

// SetExample implements SetExample interface
func (s *ObjectSchema) SetExample(example interface{}) {
	if val, ok := example.(map[string]interface{}); ok {
		s.Example(val)
	}
}
