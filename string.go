package schema

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/nyxstack/i18n"
)

// StringFormat represents the format constraint for string schemas
type StringFormat string

// Common string formats as defined in JSON Schema specification
const (
	StringFormatEmail    StringFormat = "email"
	StringFormatURI      StringFormat = "uri"
	StringFormatURL      StringFormat = "url"
	StringFormatDateTime StringFormat = "date-time"
	StringFormatDate     StringFormat = "date"
	StringFormatTime     StringFormat = "time"
	StringFormatUUID     StringFormat = "uuid"
	StringFormatHostname StringFormat = "hostname"
	StringFormatIPv4     StringFormat = "ipv4"
	StringFormatIPv6     StringFormat = "ipv6"
	StringFormatPassword StringFormat = "password"
	StringFormatBinary   StringFormat = "binary"
	StringFormatByte     StringFormat = "byte"
)

// Default error messages for string validation
var (
	stringRequiredError = i18n.S("value is required")
	stringTypeError     = i18n.S("value must be a string")
	stringPatternError  = i18n.S("value format is invalid")
	stringEnumError     = i18n.S("value must be one of the allowed values")
)

// Default error message functions that take parameters
func stringMinLengthError(min int) i18n.TranslatedFunc {
	return i18n.F("value must be at least %d characters long", min)
}

func stringMaxLengthError(max int) i18n.TranslatedFunc {
	return i18n.F("value must be at most %d characters long", max)
}

func stringFormatError(format string) i18n.TranslatedFunc {
	return i18n.F("value must be a valid %s", format)
}

func stringConstError(value string) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %v", value)
}

// StringSchema represents a JSON Schema for string values
type StringSchema struct {
	Schema
	// String-specific validation (private fields)
	minLength *int
	maxLength *int
	pattern   *string
	format    *StringFormat
	nullable  bool

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	minLengthError    ErrorMessage
	maxLengthError    ErrorMessage
	patternError      ErrorMessage
	formatError       ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	typeMismatchError ErrorMessage
}

// String creates a new string schema with optional type error message
func String(errorMessage ...interface{}) *StringSchema {
	schema := &StringSchema{
		Schema: Schema{
			schemaType: "string",
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
func (s *StringSchema) Title(title string) *StringSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *StringSchema) Description(description string) *StringSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *StringSchema) Default(value interface{}) *StringSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *StringSchema) Example(example string) *StringSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values with optional custom error message
func (s *StringSchema) Enum(values []string, errorMessage ...interface{}) *StringSchema {
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
func (s *StringSchema) Const(value string, errorMessage ...interface{}) *StringSchema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *StringSchema) Optional() *StringSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *StringSchema) Required(errorMessage ...interface{}) *StringSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *StringSchema) Nullable() *StringSchema {
	s.nullable = true
	return s
}

// TypeError sets a custom error message for type mismatch validation
func (s *StringSchema) TypeError(message string) *StringSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// String-specific fluent API methods

// MinLength sets the minimum length constraint with optional custom error message
func (s *StringSchema) MinLength(min int, errorMessage ...interface{}) *StringSchema {
	s.minLength = &min
	if len(errorMessage) > 0 {
		s.minLengthError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MaxLength sets the maximum length constraint with optional custom error message
func (s *StringSchema) MaxLength(max int, errorMessage ...interface{}) *StringSchema {
	s.maxLength = &max
	if len(errorMessage) > 0 {
		s.maxLengthError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Length sets both min and max length to the same value with optional custom error message
func (s *StringSchema) Length(length int, errorMessage ...interface{}) *StringSchema {
	s.minLength = &length
	s.maxLength = &length
	if len(errorMessage) > 0 {
		s.minLengthError = toErrorMessage(errorMessage[0])
		s.maxLengthError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Pattern sets a regex pattern constraint with optional custom error message
func (s *StringSchema) Pattern(pattern string, errorMessage ...interface{}) *StringSchema {
	s.pattern = &pattern
	if len(errorMessage) > 0 {
		s.patternError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Format sets the string format with optional custom error message
func (s *StringSchema) Format(format StringFormat, errorMessage ...interface{}) *StringSchema {
	s.format = &format
	if len(errorMessage) > 0 {
		s.formatError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *StringSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *StringSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *StringSchema) IsNullable() bool {
	return s.nullable
}

// GetMinLength returns the minimum length constraint
func (s *StringSchema) GetMinLength() *int {
	return s.minLength
}

// GetMaxLength returns the maximum length constraint
func (s *StringSchema) GetMaxLength() *int {
	return s.maxLength
}

// GetPattern returns the pattern constraint
func (s *StringSchema) GetPattern() *string {
	return s.pattern
}

// GetFormat returns the format constraint
func (s *StringSchema) GetFormat() *StringFormat {
	return s.format
}

// GetDefault returns the default value as a string
func (s *StringSchema) GetDefaultString() *string {
	if s.GetDefault() != nil {
		if str, ok := s.GetDefault().(string); ok {
			return &str
		}
	}
	return nil
}

// Convenience methods for common formats

// Email sets the format to email
func (s *StringSchema) Email() *StringSchema {
	return s.Format(StringFormatEmail)
}

// URI sets the format to URI
func (s *StringSchema) URI() *StringSchema {
	return s.Format(StringFormatURI)
}

// URL sets the format to URL
func (s *StringSchema) URL() *StringSchema {
	return s.Format(StringFormatURL)
}

// DateTime sets the format to date-time
func (s *StringSchema) DateTime() *StringSchema {
	return s.Format(StringFormatDateTime)
}

// Date sets the format to date
func (s *StringSchema) Date() *StringSchema {
	return s.Format(StringFormatDate)
}

// Time sets the format to time
func (s *StringSchema) Time() *StringSchema {
	return s.Format(StringFormatTime)
}

// UUID sets the format to UUID
func (s *StringSchema) UUID() *StringSchema {
	return s.Format(StringFormatUUID)
}

// Password sets the format to password
func (s *StringSchema) Password() *StringSchema {
	return s.Format(StringFormatPassword)
}

// Validation

// Validate validates a string value against this schema with context
// Parse validates and parses a string value, returning the final parsed value
func (s *StringSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := stringRequiredError(ctx.Locale)
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

	// Type check first
	strValue, ok := value.(string)
	if !ok {
		message := stringTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Check required (empty string case)
	if s.Schema.required && strValue == "" {
		// Check if we have a default value for empty strings
		if defaultVal := s.GetDefault(); defaultVal != nil {
			return s.Parse(defaultVal, ctx)
		}

		message := stringRequiredError(ctx.Locale)
		if !isEmptyErrorMessage(s.requiredError) {
			message = resolveErrorMessage(s.requiredError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(strValue, message, "required")},
		}
	}

	// If value is empty and not required, it's valid - return empty string or default
	if strValue == "" && !s.Schema.required {
		if defaultVal := s.GetDefault(); defaultVal != nil {
			// Return default instead of empty string
			return s.Parse(defaultVal, ctx)
		}
		return ParseResult{Valid: true, Value: "", Errors: nil}
	}

	// Now validate the string value against all constraints
	finalValue := strValue // This is our parsed value

	// Check minimum length
	if s.minLength != nil && len(strValue) < *s.minLength {
		message := stringMinLengthError(*s.minLength)(ctx.Locale)
		if !isEmptyErrorMessage(s.minLengthError) {
			message = resolveErrorMessage(s.minLengthError, ctx)
		}
		errors = append(errors, NewPrimitiveError(strValue, message, "min_length"))
	}

	// Check maximum length
	if s.maxLength != nil && len(strValue) > *s.maxLength {
		message := stringMaxLengthError(*s.maxLength)(ctx.Locale)
		if !isEmptyErrorMessage(s.maxLengthError) {
			message = resolveErrorMessage(s.maxLengthError, ctx)
		}
		errors = append(errors, NewPrimitiveError(strValue, message, "max_length"))
	}

	// Check pattern
	if s.pattern != nil {
		matched, err := regexp.MatchString(*s.pattern, strValue)
		if err != nil || !matched {
			message := stringPatternError(ctx.Locale)
			if !isEmptyErrorMessage(s.patternError) {
				message = resolveErrorMessage(s.patternError, ctx)
			}
			errors = append(errors, NewPrimitiveError(strValue, message, "pattern"))
		}
	}

	// Check format
	if s.format != nil {
		if !s.validateFormat(strValue, *s.format) {
			message := stringFormatError(string(*s.format))(ctx.Locale)
			if !isEmptyErrorMessage(s.formatError) {
				message = resolveErrorMessage(s.formatError, ctx)
			}
			errors = append(errors, NewPrimitiveError(strValue, message, "format"))
		}
	}

	// Check enum
	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == strValue {
				valid = true
				break
			}
		}
		if !valid {
			message := stringEnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(strValue, message, "enum"))
		}
	}

	// Check const
	if s.Schema.constVal != nil && s.Schema.constVal != strValue {
		message := stringConstError(fmt.Sprintf("%v", s.Schema.constVal))(ctx.Locale)
		if !isEmptyErrorMessage(s.constError) {
			message = resolveErrorMessage(s.constError, ctx)
		}
		errors = append(errors, NewPrimitiveError(strValue, message, "const"))
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  finalValue,
		Errors: errors,
	}
}

// MarshalJSON implements json.Marshaler to properly serialize StringSchema for JSON schema generation
func (s *StringSchema) MarshalJSON() ([]byte, error) {
	type jsonStringSchema struct {
		Schema
		MinLength *int          `json:"minLength,omitempty"`
		MaxLength *int          `json:"maxLength,omitempty"`
		Pattern   *string       `json:"pattern,omitempty"`
		Format    *StringFormat `json:"format,omitempty"`
		Nullable  bool          `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonStringSchema{
		Schema:    s.Schema,
		MinLength: s.minLength,
		MaxLength: s.maxLength,
		Pattern:   s.pattern,
		Format:    s.format,
		Nullable:  s.nullable,
	})
}

// validateFormat validates a string against a specific format
func (s *StringSchema) validateFormat(value string, format StringFormat) bool {
	switch format {
	case StringFormatEmail:
		// Simple email validation regex
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, value)
		return matched
	case StringFormatURI, StringFormatURL:
		// Basic URL validation - starts with http/https or is a valid URI
		urlRegex := `^https?://[^\s/$.?#].[^\s]*$|^[a-zA-Z][a-zA-Z0-9+.-]*:[^\s]*$`
		matched, _ := regexp.MatchString(urlRegex, value)
		return matched
	case StringFormatUUID:
		// UUID v4 format validation
		uuidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`
		matched, _ := regexp.MatchString(uuidRegex, value)
		return matched
	case StringFormatDateTime:
		// ISO 8601 date-time format (basic validation)
		dateTimeRegex := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?([+-]\d{2}:\d{2}|Z)$`
		matched, _ := regexp.MatchString(dateTimeRegex, value)
		return matched
	case StringFormatDate:
		// ISO 8601 date format
		dateRegex := `^\d{4}-\d{2}-\d{2}$`
		matched, _ := regexp.MatchString(dateRegex, value)
		return matched
	case StringFormatTime:
		// ISO 8601 time format
		timeRegex := `^\d{2}:\d{2}:\d{2}(\.\d{3})?$`
		matched, _ := regexp.MatchString(timeRegex, value)
		return matched
	case StringFormatIPv4:
		// IPv4 format validation
		ipv4Regex := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
		matched, _ := regexp.MatchString(ipv4Regex, value)
		return matched
	case StringFormatIPv6:
		// IPv6 format validation (simplified)
		ipv6Regex := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$`
		matched, _ := regexp.MatchString(ipv6Regex, value)
		return matched
	case StringFormatHostname:
		// Basic hostname validation
		hostnameRegex := `^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
		matched, _ := regexp.MatchString(hostnameRegex, value)
		return matched
	default:
		// For custom formats or unsupported formats, assume valid
		return true
	}
}

// JSON generates JSON Schema representation
func (s *StringSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("string")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add string-specific fields
	addOptionalField(schema, "minLength", s.minLength)
	addOptionalField(schema, "maxLength", s.maxLength)
	addOptionalField(schema, "pattern", s.pattern)
	if s.format != nil {
		schema["format"] = string(*s.format)
	}

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"string", "null"}
	}

	return schema
}

// Interface implementations for StringSchema

// SetTitle implements SetTitle interface
func (s *StringSchema) SetTitle(title string) {
	s.Title(title)
}

// SetDescription implements SetDescription interface
func (s *StringSchema) SetDescription(description string) {
	s.Description(description)
}

// SetRequired implements SetRequired interface
func (s *StringSchema) SetRequired() {
	s.Required()
}

// SetOptional implements SetOptional interface
func (s *StringSchema) SetOptional() {
	s.Optional()
}

// SetMinLength implements SetMinLength interface
func (s *StringSchema) SetMinLength(minLen int) {
	s.MinLength(minLen)
}

// SetMaxLength implements SetMaxLength interface
func (s *StringSchema) SetMaxLength(maxLen int) {
	s.MaxLength(maxLen)
}

// SetPattern implements SetPattern interface
func (s *StringSchema) SetPattern(pattern string) {
	s.Pattern(pattern)
}

// SetNullable implements SetNullable interface
func (s *StringSchema) SetNullable() {
	s.Nullable()
}

// SetDefault implements SetDefault interface
func (s *StringSchema) SetDefault(value interface{}) {
	s.Default(value)
}

// SetExample implements SetExample interface
func (s *StringSchema) SetExample(example interface{}) {
	if str, ok := example.(string); ok {
		s.Example(str)
	}
}
