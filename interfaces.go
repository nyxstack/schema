package schema

// Common interfaces for schema setter methods
// These interfaces allow for type-safe casting and method calls without knowing the concrete type

// SetTitle interface for schemas that support setting a title
type SetTitle interface {
	SetTitle(title string)
}

// SetDescription interface for schemas that support setting a description
type SetDescription interface {
	SetDescription(description string)
}

// SetRequired interface for schemas that support setting required status
type SetRequired interface {
	SetRequired()
}

// SetOptional interface for schemas that support setting optional status
type SetOptional interface {
	SetOptional()
}

// SetMinimum interface for schemas that support setting minimum values (integers/numbers)
type SetMinimum interface {
	SetMinimum(min int)
}

// SetMaximum interface for schemas that support setting maximum values (integers/numbers)
type SetMaximum interface {
	SetMaximum(max int)
}

// SetMinimumFloat interface for schemas that support setting minimum float values (numbers)
type SetMinimumFloat interface {
	SetMinimumFloat(min float64)
}

// SetMaximumFloat interface for schemas that support setting maximum float values (numbers)
type SetMaximumFloat interface {
	SetMaximumFloat(max float64)
}

// SetMinLength interface for schemas that support setting minimum length (strings/arrays)
type SetMinLength interface {
	SetMinLength(minLen int)
}

// SetMaxLength interface for schemas that support setting maximum length (strings/arrays)
type SetMaxLength interface {
	SetMaxLength(maxLen int)
}

// SetPattern interface for schemas that support setting regex patterns (strings)
type SetPattern interface {
	SetPattern(pattern string)
}

// SetMinItems interface for schemas that support setting minimum items (arrays)
type SetMinItems interface {
	SetMinItems(min int)
}

// SetMaxItems interface for schemas that support setting maximum items (arrays)
type SetMaxItems interface {
	SetMaxItems(max int)
}

// SetNullable interface for schemas that support setting nullable status
type SetNullable interface {
	SetNullable()
}

// SetDefault interface for schemas that support setting default values
type SetDefault interface {
	SetDefault(value interface{})
}

// SetExample interface for schemas that support setting example values
type SetExample interface {
	SetExample(example interface{})
}
