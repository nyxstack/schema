package schema

import (
	"encoding/json"
	"fmt"

	"github.com/nyxstack/i18n"
)

// Default error messages for transform validation
var (
	transformRequiredError = i18n.S("field is required")
)

func transformFailedError(err error) i18n.TranslatedFunc {
	return i18n.F("transformation failed: %v", err)
}

// TransformFunc represents a function that transforms one type to another
type TransformFunc func(input interface{}) (interface{}, error)

// TransformSchema represents a schema that validates input, transforms it, then validates output
type TransformSchema struct {
	Schema
	inputSchema   Parseable     // Schema to validate input
	outputSchema  Parseable     // Schema to validate output
	transformFunc TransformFunc // Function to transform input to output
	nullable      bool          // Whether the schema allows null values

	// Error messages
	requiredError  ErrorMessage `json:"-"`
	transformError ErrorMessage `json:"-"`
}

// Transform creates a new transform schema
func Transform(
	inputSchema Parseable,
	outputSchema Parseable,
	transformFunc TransformFunc,
	errorMessage ...interface{},
) *TransformSchema {
	schema := &TransformSchema{
		Schema: Schema{
			schemaType: "transform",
			required:   false,
		},
		inputSchema:   inputSchema,
		outputSchema:  outputSchema,
		transformFunc: transformFunc,
	}

	// Set custom error message if provided
	if len(errorMessage) > 0 {
		schema.transformError = parseErrorMessageToErrorMessage(errorMessage...)
	}

	return schema
}

// Title sets the title of the transform schema
func (s *TransformSchema) Title(title string) *TransformSchema {
	s.Schema.title = title
	return s
}

// Description sets the description of the transform schema
func (s *TransformSchema) Description(description string) *TransformSchema {
	s.Schema.description = description
	return s
}

// Required marks the schema as required with optional custom error message
func (s *TransformSchema) Required(errorMessage ...interface{}) *TransformSchema {
	s.Schema.required = true
	if len(errorMessage) > 0 {
		s.requiredError = parseErrorMessageToErrorMessage(errorMessage...)
	}
	return s
}

// Optional marks the schema as optional
func (s *TransformSchema) Optional() *TransformSchema {
	s.Schema.required = false
	return s
}

// Nullable marks the schema as nullable
func (s *TransformSchema) Nullable() *TransformSchema {
	s.nullable = true
	return s
}

// Default sets a default value for the schema
func (s *TransformSchema) Default(value interface{}) *TransformSchema {
	s.Schema.defaultValue = value
	return s
}

// WithTransformError sets a custom error message for transformation failures
func (s *TransformSchema) WithTransformError(errorMessage ...interface{}) *TransformSchema {
	s.transformError = parseErrorMessageToErrorMessage(errorMessage...)
	return s
}

// Parse validates input, transforms it, then validates output
func (s *TransformSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	// Handle nil values
	if value == nil {
		if s.nullable {
			return ParseResult{Valid: true, Value: nil, Errors: nil}
		}
		if s.Schema.required {
			// Check if we have a default value
			if defaultVal := s.GetDefault(); defaultVal != nil {
				return s.Parse(defaultVal, ctx)
			}
			// Required field is missing
			message := transformRequiredError(ctx.Locale)
			if s.requiredError != nil {
				message = s.requiredError.Resolve(ctx)
			}
			return ParseResult{
				Valid:  false,
				Value:  nil,
				Errors: []ValidationError{NewPrimitiveError(value, message, "required")},
			}
		}
		// Use default value if available for optional fields
		if defaultVal := s.GetDefault(); defaultVal != nil {
			return s.Parse(defaultVal, ctx)
		}
		// Optional and no default, return nil
		return ParseResult{Valid: true, Value: nil, Errors: nil}
	}

	// Step 1: Validate and parse input against input schema
	inputResult := s.inputSchema.Parse(value, ctx)
	if !inputResult.Valid {
		// Prefix input validation errors
		var prefixedErrors []ValidationError
		for _, err := range inputResult.Errors {
			prefixedErrors = append(prefixedErrors, ValidationError{
				Path:    err.Path,
				Value:   err.Value,
				Message: "input validation: " + err.Message,
				Code:    "input_" + err.Code,
			})
		}
		return ParseResult{
			Valid:  false,
			Value:  value,
			Errors: prefixedErrors,
		}
	}

	// Step 2: Transform the validated input value
	transformed, transformErr := s.transformFunc(inputResult.Value)
	if transformErr != nil {
		message := transformFailedError(transformErr)(ctx.Locale)
		if s.transformError != nil {
			message = s.transformError.Resolve(ctx)
		}

		return ParseResult{
			Valid:  false,
			Value:  value,
			Errors: []ValidationError{NewPrimitiveError(value, message, "transform")},
		}
	}

	// Step 3: Validate and parse transformed output against output schema
	outputResult := s.outputSchema.Parse(transformed, ctx)
	if !outputResult.Valid {
		// Prefix output validation errors
		var prefixedErrors []ValidationError
		for _, err := range outputResult.Errors {
			prefixedErrors = append(prefixedErrors, ValidationError{
				Path:    err.Path,
				Value:   err.Value,
				Message: "output validation: " + err.Message,
				Code:    "output_" + err.Code,
			})
		}
		return ParseResult{
			Valid:  false,
			Value:  transformed,
			Errors: prefixedErrors,
		}
	}

	// Success: return the final transformed and validated value
	return ParseResult{
		Valid:  true,
		Value:  outputResult.Value,
		Errors: nil,
	}
}

// JSON returns the JSON representation of the transform schema
func (s *TransformSchema) JSON() map[string]interface{} {
	result := make(map[string]interface{})

	// Set basic properties
	result["type"] = "transform"

	// Add input schema
	if inputJSON, ok := s.inputSchema.(interface{ JSON() map[string]interface{} }); ok {
		result["inputSchema"] = inputJSON.JSON()
	}

	// Add output schema
	if outputJSON, ok := s.outputSchema.(interface{ JSON() map[string]interface{} }); ok {
		result["outputSchema"] = outputJSON.JSON()
	}

	// Add metadata
	if s.title != "" {
		result["title"] = s.title
	}
	if s.description != "" {
		result["description"] = s.description
	}

	// Add schema flags
	if s.nullable {
		result["nullable"] = true
	}
	if defaultVal := s.GetDefault(); defaultVal != nil {
		result["default"] = defaultVal
	}

	return result
}

// MarshalJSON implements json.Marshaler
func (s *TransformSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.JSON())
}

// Helper function to convert error message parameters to ErrorMessage
func parseErrorMessageToErrorMessage(errorMessage ...interface{}) ErrorMessage {
	if len(errorMessage) == 0 {
		return nil
	}

	switch msg := errorMessage[0].(type) {
	case string:
		return Msg(msg)
	case i18n.TranslatedFunc:
		return I18nMessage(msg)
	case ErrorMessage:
		return msg
	default:
		return Msg(fmt.Sprintf("%v", msg))
	}
}
