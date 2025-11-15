package schema

import (
	"github.com/nyxstack/i18n"
)

// Default error messages for conditional validation
var (
	conditionalThenFailedError = i18n.S("value matches the 'if' condition but fails the 'then' validation")
	conditionalElseFailedError = i18n.S("value does not match the 'if' condition but fails the 'else' validation")
)

// ConditionalErrors defines error message functions
var ConditionalErrors = struct {
	ThenFailed i18n.TranslatedFunc
	ElseFailed i18n.TranslatedFunc
}{
	ThenFailed: conditionalThenFailedError,
	ElseFailed: conditionalElseFailedError,
}

// ConditionalSchema represents an if-then-else validation schema
type ConditionalSchema struct {
	ifSchema   Parseable
	thenSchema Parseable
	elseSchema Parseable
	thenError  ErrorMessage
	elseError  ErrorMessage
}

// Conditional creates a new Conditional schema with if condition
func Conditional(ifSchema Parseable) *ConditionalSchema {
	return &ConditionalSchema{
		ifSchema: ifSchema,
	}
}

// Then sets the schema that must be valid if the 'if' condition matches
func (s *ConditionalSchema) Then(thenSchema Parseable) *ConditionalSchema {
	s.thenSchema = thenSchema
	return s
}

// Else sets the schema that must be valid if the 'if' condition does not match
func (s *ConditionalSchema) Else(elseSchema Parseable) *ConditionalSchema {
	s.elseSchema = elseSchema
	return s
}

// ThenError sets a custom error message for when the 'then' validation fails
func (s *ConditionalSchema) ThenError(err ErrorMessage) *ConditionalSchema {
	s.thenError = err
	return s
}

// ElseError sets a custom error message for when the 'else' validation fails
func (s *ConditionalSchema) ElseError(err ErrorMessage) *ConditionalSchema {
	s.elseError = err
	return s
}

// Parse validates using if-then-else logic
func (s *ConditionalSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	// First, test the 'if' condition
	ifResult := s.ifSchema.Parse(value, ctx)

	if ifResult.Valid {
		// If condition matched, apply 'then' schema
		if s.thenSchema != nil {
			thenResult := s.thenSchema.Parse(value, ctx)
			if !thenResult.Valid {
				// 'Then' schema failed
				message := ConditionalErrors.ThenFailed(ctx.Locale)
				if !isEmptyErrorMessage(s.thenError) {
					message = resolveErrorMessage(s.thenError, ctx)
				}

				// Combine the original errors with our conditional error
				errors := []ValidationError{NewPrimitiveError(value, message, "then_failed")}
				errors = append(errors, thenResult.Errors...)

				return ParseResult{
					Valid:  false,
					Value:  value,
					Errors: errors,
				}
			}

			// 'Then' schema passed, use its transformed value
			return thenResult
		}

		// No 'then' schema specified, just return the value
		return ParseResult{
			Valid:  true,
			Value:  value,
			Errors: nil,
		}
	} else {
		// If condition did not match, apply 'else' schema if present
		if s.elseSchema != nil {
			elseResult := s.elseSchema.Parse(value, ctx)
			if !elseResult.Valid {
				// 'Else' schema failed
				message := ConditionalErrors.ElseFailed(ctx.Locale)
				if !isEmptyErrorMessage(s.elseError) {
					message = resolveErrorMessage(s.elseError, ctx)
				}

				// Combine the original errors with our conditional error
				errors := []ValidationError{NewPrimitiveError(value, message, "else_failed")}
				errors = append(errors, elseResult.Errors...)

				return ParseResult{
					Valid:  false,
					Value:  value,
					Errors: errors,
				}
			}

			// 'Else' schema passed, use its transformed value
			return elseResult
		}

		// No 'else' schema specified, just return the value
		return ParseResult{
			Valid:  true,
			Value:  value,
			Errors: nil,
		}
	}
}

// JSON generates JSON Schema for Conditional validation
func (s *ConditionalSchema) JSON() map[string]interface{} {
	schema := map[string]interface{}{}

	// Add 'if' schema
	if ifSchema, ok := s.ifSchema.(interface{ JSON() map[string]interface{} }); ok {
		schema["if"] = ifSchema.JSON()
	} else {
		schema["if"] = map[string]interface{}{"type": "unknown"}
	}

	// Add 'then' schema if present
	if s.thenSchema != nil {
		if thenSchema, ok := s.thenSchema.(interface{ JSON() map[string]interface{} }); ok {
			schema["then"] = thenSchema.JSON()
		} else {
			schema["then"] = map[string]interface{}{"type": "unknown"}
		}
	}

	// Add 'else' schema if present
	if s.elseSchema != nil {
		if elseSchema, ok := s.elseSchema.(interface{ JSON() map[string]interface{} }); ok {
			schema["else"] = elseSchema.JSON()
		} else {
			schema["else"] = map[string]interface{}{"type": "unknown"}
		}
	}

	return schema
}
