package schema

import (
	"context"
	"fmt"
)

// ValidationContext contains locale and other context information for validation
type ValidationContext struct {
	Locale string
	Ctx    context.Context
}

// DefaultValidationContext returns a context with English locale
func DefaultValidationContext() *ValidationContext {
	return &ValidationContext{
		Locale: "en",
		Ctx:    context.Background(),
	}
}

// NewValidationContext creates a validation context with specified locale
func NewValidationContext(locale string) *ValidationContext {
	return &ValidationContext{
		Locale: locale,
		Ctx:    context.Background(),
	}
}

// WithContext sets the Go context
func (vc *ValidationContext) WithContext(ctx context.Context) *ValidationContext {
	vc.Ctx = ctx
	return vc
}

// Parseable interface that all schemas should implement
type Parseable interface {
	Parse(value interface{}, ctx *ValidationContext) ParseResult
}

// ValidationError represents a validation error with details
type ValidationError struct {
	Path    []string `json:"path"`    // Path to the field (empty for primitive values)
	Value   string   `json:"value"`   // String representation of the invalid value
	Message string   `json:"message"` // Human-readable error message
	Code    string   `json:"code"`    // Machine-readable error code
}

// NewPrimitiveError creates a validation error for primitive value validation
func NewPrimitiveError(value interface{}, message, code string) ValidationError {
	return ValidationError{
		Path:    []string{}, // Empty path for primitive values
		Value:   fmt.Sprintf("%v", value),
		Message: message,
		Code:    code,
	}
}

// NewFieldError creates a validation error for object field validation
func NewFieldError(path []string, value interface{}, message, code string) ValidationError {
	return ValidationError{
		Path:    path,
		Value:   fmt.Sprintf("%v", value),
		Message: message,
		Code:    code,
	}
}

// ParseResult contains parsing and validation results with the final parsed value
type ParseResult struct {
	Valid  bool              `json:"valid"`
	Value  interface{}       `json:"value"` // The final parsed/transformed value
	Errors []ValidationError `json:"errors"`
}

// ValidationResult contains validation results (deprecated, use ParseResult)
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}
