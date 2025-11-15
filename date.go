package schema

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/nyxstack/i18n"
)

// Default error messages for date validation
var (
	dateRequiredError = i18n.S("value is required")
	dateTypeError     = i18n.S("value must be a date string")
	dateFormatError   = i18n.S("value must be a valid date format")
	dateEnumError     = i18n.S("value must be one of the allowed dates")
)

func dateConstError(value string) i18n.TranslatedFunc {
	return i18n.F("value must be exactly: %s", value)
}

func dateRangeError(min, max string) i18n.TranslatedFunc {
	return i18n.F("value must be between %s and %s", min, max)
}

// DateFormat represents supported date/time formats
type DateFormat string

const (
	// Standard JSON Schema formats
	FormatDate     DateFormat = "date"      // YYYY-MM-DD
	FormatDateTime DateFormat = "date-time" // RFC3339: 2006-01-02T15:04:05Z07:00
	FormatTime     DateFormat = "time"      // HH:MM:SS or HH:MM:SS.sss

	// Additional common formats
	FormatDateOnly DateFormat = "date-only" // YYYY-MM-DD (same as date)
	FormatTimeOnly DateFormat = "time-only" // HH:MM:SS (same as time)
	FormatISO8601  DateFormat = "iso8601"   // ISO 8601 format
	FormatRFC3339  DateFormat = "rfc3339"   // RFC 3339 format
	FormatUnix     DateFormat = "unix"      // Unix timestamp (as string)
)

// DateSchema represents a JSON Schema for date/time values
type DateSchema struct {
	Schema
	// Date-specific validation
	format   DateFormat // Date format to validate against
	minDate  *time.Time // Minimum date/time
	maxDate  *time.Time // Maximum date/time
	nullable bool       // Allow null values

	// Error messages for validation failures (support i18n)
	requiredError     ErrorMessage
	enumError         ErrorMessage
	constError        ErrorMessage
	formatError       ErrorMessage
	rangeError        ErrorMessage
	typeMismatchError ErrorMessage
}

// Date creates a new date schema with default date format (YYYY-MM-DD)
func Date(errorMessage ...interface{}) *DateSchema {
	schema := &DateSchema{
		Schema: Schema{
			schemaType: "string",
			required:   true, // Default to required
		},
		format: FormatDate,
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

// DateTime creates a new datetime schema with RFC3339 format
func DateTime(errorMessage ...interface{}) *DateSchema {
	schema := &DateSchema{
		Schema: Schema{
			schemaType: "string",
			required:   true, // Default to required
		},
		format: FormatDateTime,
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

// Time creates a new time schema with time format (HH:MM:SS)
func Time(errorMessage ...interface{}) *DateSchema {
	schema := &DateSchema{
		Schema: Schema{
			schemaType: "string",
			required:   true, // Default to required
		},
		format: FormatTime,
	}
	if len(errorMessage) > 0 {
		schema.typeMismatchError = toErrorMessage(errorMessage[0])
	}
	return schema
}

// Core fluent API methods

// Title sets the title of the schema
func (s *DateSchema) Title(title string) *DateSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the schema
func (s *DateSchema) Description(description string) *DateSchema {
	s.Schema.description = description
	return s
}

// Default sets the default value
func (s *DateSchema) Default(value interface{}) *DateSchema {
	s.Schema.defaultValue = value
	return s
}

// Example adds an example value
func (s *DateSchema) Example(example string) *DateSchema {
	s.Schema.examples = append(s.Schema.examples, example)
	return s
}

// Enum sets the allowed enum values with optional custom error message
func (s *DateSchema) Enum(values []string, errorMessage ...interface{}) *DateSchema {
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
func (s *DateSchema) Const(value string, errorMessage ...interface{}) *DateSchema {
	s.Schema.constVal = value
	if len(errorMessage) > 0 {
		s.constError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Date-specific validation

// Format sets the date format to validate against
func (s *DateSchema) Format(format DateFormat) *DateSchema {
	s.format = format
	return s
}

// MinDate sets the minimum date/time constraint
func (s *DateSchema) MinDate(min time.Time, errorMessage ...interface{}) *DateSchema {
	s.minDate = &min
	if len(errorMessage) > 0 {
		s.rangeError = toErrorMessage(errorMessage[0])
	}
	return s
}

// MaxDate sets the maximum date/time constraint
func (s *DateSchema) MaxDate(max time.Time, errorMessage ...interface{}) *DateSchema {
	s.maxDate = &max
	if len(errorMessage) > 0 {
		s.rangeError = toErrorMessage(errorMessage[0])
	}
	return s
}

// DateRange sets both min and max date constraints
func (s *DateSchema) DateRange(min, max time.Time) *DateSchema {
	s.minDate = &min
	s.maxDate = &max
	return s
}

// Required/Optional/Nullable control

// Optional marks the schema as optional
func (s *DateSchema) Optional() *DateSchema {
	s.Schema.required = false
	return s
}

// Required marks the schema as required (default behavior) with optional custom error message
func (s *DateSchema) Required(errorMessage ...interface{}) *DateSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = toErrorMessage(errorMessage[0])
	}
	return s
}

// Nullable marks the schema as nullable (allows nil values)
func (s *DateSchema) Nullable() *DateSchema {
	s.nullable = true
	return s
}

// Error customization

// TypeError sets a custom error message for type mismatch validation
func (s *DateSchema) TypeError(message string) *DateSchema {
	s.typeMismatchError = toErrorMessage(message)
	return s
}

// FormatError sets a custom error message for format validation
func (s *DateSchema) FormatError(message string) *DateSchema {
	s.formatError = toErrorMessage(message)
	return s
}

// Getters for accessing private fields

// IsRequired returns whether the schema is marked as required
func (s *DateSchema) IsRequired() bool {
	return s.Schema.required
}

// IsOptional returns whether the schema is marked as optional
func (s *DateSchema) IsOptional() bool {
	return !s.Schema.required
}

// IsNullable returns whether the schema allows nil values
func (s *DateSchema) IsNullable() bool {
	return s.nullable
}

// GetFormat returns the date format
func (s *DateSchema) GetFormat() DateFormat {
	return s.format
}

// GetMinDate returns the minimum date constraint
func (s *DateSchema) GetMinDate() *time.Time {
	return s.minDate
}

// GetMaxDate returns the maximum date constraint
func (s *DateSchema) GetMaxDate() *time.Time {
	return s.maxDate
}

// Validation helpers

// validateDateFormat validates a date string against the specified format
func (s *DateSchema) validateDateFormat(dateStr string) (*time.Time, error) {
	var layout string
	var pattern *regexp.Regexp

	switch s.format {
	case FormatDate, FormatDateOnly:
		layout = "2006-01-02"
		pattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	case FormatDateTime, FormatRFC3339, FormatISO8601:
		layout = time.RFC3339
		// More flexible pattern for RFC3339
		pattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)

	case FormatTime, FormatTimeOnly:
		layout = "15:04:05"
		pattern = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}`)

	case FormatUnix:
		// Unix timestamp validation (numbers only)
		pattern = regexp.MustCompile(`^\d+$`)
		// For unix timestamp, we don't parse as time.Time here
		if pattern.MatchString(dateStr) {
			return nil, nil // Valid unix timestamp format
		}
		return nil, &time.ParseError{Layout: "unix", Value: dateStr, Message: dateFormatError("en")}

	default:
		// Default to RFC3339
		layout = time.RFC3339
		pattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)
	}

	// First check pattern
	if pattern != nil && !pattern.MatchString(dateStr) {
		return nil, &time.ParseError{Layout: layout, Value: dateStr, Message: dateFormatError("en")}
	}

	// Then parse the actual date
	if s.format != FormatUnix {
		parsed, err := time.Parse(layout, dateStr)
		if err != nil {
			return nil, err
		}
		return &parsed, nil
	}

	return nil, nil
}

// Validation

// Parse validates and parses a date value, returning the final parsed value
func (s *DateSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
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
			message := dateRequiredError(ctx.Locale)
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

	// Type check
	dateString, ok := value.(string)
	if !ok {
		message := dateTypeError(ctx.Locale)
		if !isEmptyErrorMessage(s.typeMismatchError) {
			message = resolveErrorMessage(s.typeMismatchError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  nil,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_type")},
		}
	}

	// Validate format
	parsedTime, err := s.validateDateFormat(dateString)
	if err != nil {
		message := dateFormatError(ctx.Locale)
		if !isEmptyErrorMessage(s.formatError) {
			message = resolveErrorMessage(s.formatError, ctx)
		}
		errors = append(errors, NewPrimitiveError(dateString, message, "format"))
	}

	// Check enum
	if len(s.Schema.enum) > 0 {
		valid := false
		for _, enumValue := range s.Schema.enum {
			if enumValue == dateString {
				valid = true
				break
			}
		}
		if !valid {
			message := dateEnumError(ctx.Locale)
			if !isEmptyErrorMessage(s.enumError) {
				message = resolveErrorMessage(s.enumError, ctx)
			}
			errors = append(errors, NewPrimitiveError(dateString, message, "enum"))
		}
	}

	// Check const
	if s.Schema.constVal != nil {
		if constStr, ok := s.Schema.constVal.(string); ok && constStr != dateString {
			message := dateConstError(constStr)(ctx.Locale)
			if !isEmptyErrorMessage(s.constError) {
				message = resolveErrorMessage(s.constError, ctx)
			}
			errors = append(errors, NewPrimitiveError(dateString, message, "const"))
		}
	}

	// Check date range constraints (only if we successfully parsed the date)
	if parsedTime != nil {
		if s.minDate != nil && parsedTime.Before(*s.minDate) {
			minStr := s.minDate.Format("2006-01-02")
			maxStr := ""
			if s.maxDate != nil {
				maxStr = s.maxDate.Format("2006-01-02")
			} else {
				maxStr = "∞"
			}
			message := dateRangeError(minStr, maxStr)(ctx.Locale)
			if !isEmptyErrorMessage(s.rangeError) {
				message = resolveErrorMessage(s.rangeError, ctx)
			}
			errors = append(errors, NewPrimitiveError(dateString, message, "min_date"))
		}

		if s.maxDate != nil && parsedTime.After(*s.maxDate) {
			minStr := ""
			if s.minDate != nil {
				minStr = s.minDate.Format("2006-01-02")
			} else {
				minStr = "-∞"
			}
			maxStr := s.maxDate.Format("2006-01-02")
			message := dateRangeError(minStr, maxStr)(ctx.Locale)
			if !isEmptyErrorMessage(s.rangeError) {
				message = resolveErrorMessage(s.rangeError, ctx)
			}
			errors = append(errors, NewPrimitiveError(dateString, message, "max_date"))
		}
	}

	return ParseResult{
		Valid:  len(errors) == 0,
		Value:  dateString, // Return the original string value
		Errors: errors,
	}
}

// JSON generates JSON Schema representation
func (s *DateSchema) JSON() map[string]interface{} {
	schema := baseJSONSchema("string")

	// Add base schema fields
	addTitle(schema, s.GetTitle())
	addDescription(schema, s.GetDescription())
	addOptionalField(schema, "default", s.GetDefault())
	addOptionalArray(schema, "examples", s.GetExamples())
	addOptionalArray(schema, "enum", s.GetEnum())
	addOptionalField(schema, "const", s.GetConst())

	// Add format
	schema["format"] = string(s.format)

	// Add nullable if true
	if s.nullable {
		schema["type"] = []string{"string", "null"}
	}

	return schema
}

// MarshalJSON implements json.Marshaler to properly serialize DateSchema for JSON schema generation
func (s *DateSchema) MarshalJSON() ([]byte, error) {
	type jsonDateSchema struct {
		Schema
		Format   DateFormat `json:"format"`
		MinDate  *time.Time `json:"minDate,omitempty"`
		MaxDate  *time.Time `json:"maxDate,omitempty"`
		Nullable bool       `json:"nullable,omitempty"`
	}

	return json.Marshal(jsonDateSchema{
		Schema:   s.Schema,
		Format:   s.format,
		MinDate:  s.minDate,
		MaxDate:  s.maxDate,
		Nullable: s.nullable,
	})
}
