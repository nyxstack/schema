package schema

import (
	"testing"
	"time"
)

// Test Any Schema
func TestAnySchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Any()

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"string", "hello", true},
		{"number", 42, true},
		{"float", 3.14, true},
		{"boolean", true, true},
		{"nil", nil, true},
		{"array", []string{"test"}, true},
		{"object", map[string]string{"key": "value"}, true},
		{"empty string", "", true},
		{"zero", 0, true},
		{"false", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Any.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
			}
		})
	}
}

// Test Tuple Schema
func TestTupleSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Tuple(
		String().MinLength(2),
		Int().Min(0),
		Bool(),
	)

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid tuple", []interface{}{"hello", 42, true}, true},
		{"valid with exact constraints", []interface{}{"hi", 0, false}, true},
		{"string too short", []interface{}{"a", 42, true}, false},
		{"negative number", []interface{}{"hello", -1, true}, false},
		{"wrong type in position", []interface{}{"hello", "42", true}, false},
		{"too few items", []interface{}{"hello", 42}, false},
		{"too many items", []interface{}{"hello", 42, true, "extra"}, false},
		{"not an array", "not an array", false},
		{"nil", nil, false},
		{"empty array", []interface{}{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Tuple.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test AllOf Schema
func TestAllOfSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := AllOf(
		String(),
		String().MinLength(3),
		String().MaxLength(10),
	)

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid all constraints", "hello", true},
		{"min length exact", "abc", true},
		{"max length exact", "1234567890", true},
		{"too short", "hi", false},
		{"too long", "this is too long", false},
		{"not a string", 123, false},
		{"nil", nil, false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("AllOf.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test AnyOf Schema
func TestAnyOfSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := AnyOf(
		String().MinLength(5),
		Int().Min(100),
	)

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid string", "hello world", true},
		{"valid number", 150, true},
		{"string min length exact", "hello", true},
		{"number min exact", 100, true},
		{"string too short", "hi", false},
		{"number too small", 50, false},
		{"matches neither", true, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("AnyOf.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Date Schema
func TestDateSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *DateSchema
		value    interface{}
		expected bool
	}{
		{"valid date", Date(), "2024-12-25", true},
		{"valid datetime", DateTime(), "2024-12-25T15:30:00Z", true},
		{"valid time", Time(), "15:30:00", true},
		{"invalid date format", Date(), "25/12/2024", false},
		{"invalid date value", Date(), "2024-02-30", false},
		{"datetime for date schema", Date(), "2024-12-25T15:30:00Z", false},
		{"date for datetime schema", DateTime(), "2024-12-25", false},
		{"not a string", Date(), 20241225, false},
		{"nil", Date(), nil, false},
		{"empty string", Date(), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("DateSchema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

func TestDateSchema_Range(t *testing.T) {
	ctx := DefaultValidationContext()
	minDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	schema := Date().MinDate(minDate).MaxDate(maxDate)

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid in range", "2024-06-15", true},
		{"min date exact", "2024-01-01", true},
		{"max date exact", "2024-12-31", true},
		{"before min", "2023-12-31", false},
		{"after max", "2025-01-01", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("DateSchema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test UUID Schema
func TestUUIDSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *UUIDSchema
		value    interface{}
		expected bool
	}{
		{"valid uuid v4", UUID(), "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid uuid any version", UUID(), "6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"uuid v4 specific", UUID().Version(UUIDVersion4), "550e8400-e29b-41d4-a716-446655440000", true},
		{"uuid v1 for v4 schema", UUID().Version(UUIDVersion4), "6ba7b810-9dad-11d1-80b4-00c04fd430c8", false},
		{"compact format", UUID().Format(UUIDFormatCompact), "550e8400e29b41d4a716446655440000", true},
		{"hyphenated for compact", UUID().Format(UUIDFormatCompact), "550e8400-e29b-41d4-a716-446655440000", false},
		{"invalid format", UUID(), "not-a-uuid", false},
		{"not a string", UUID(), 123, false},
		{"nil", UUID(), nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("UUIDSchema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Not Schema
func TestNotSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Not(Int().Max(-1)) // Not a negative number

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"positive number", 42, true},
		{"zero", 0, true},
		{"negative number", -5, false},
		{"not a number", "hello", true},
		{"boolean", true, true},
		{"nil", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Not.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Conditional Schema
func TestConditionalSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Conditional(String()).
		Then(String().MinLength(5)).
		Else(Int())

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"string matches then", "hello world", true},
		{"string fails then", "hi", false},
		{"number matches else", 42, true},
		{"boolean fails else", true, false},
		{"nil fails both", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Conditional.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Ref Schema
func TestRefSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	registry := NewSchemaRegistry()

	// Define schemas
	registry.Define("PersonName", String().MinLength(2).MaxLength(50))
	registry.Define("PersonAge", Int().Min(0).Max(150))

	nameRef := Ref("#/PersonName", registry)
	ageRef := Ref("#/PersonAge", registry)

	tests := []struct {
		name     string
		schema   *RefSchema
		value    interface{}
		expected bool
	}{
		{"valid name ref", nameRef, "John Doe", true},
		{"invalid name ref", nameRef, "A", false},
		{"valid age ref", ageRef, 25, true},
		{"invalid age ref", ageRef, 200, false},
		{"name ref wrong type", nameRef, 123, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Ref.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

func TestRefSchema_NotFound(t *testing.T) {
	ctx := DefaultValidationContext()
	registry := NewSchemaRegistry()

	// Reference to non-existent schema
	ref := Ref("#/NonExistent", registry)

	result := ref.Parse("test", ctx)
	if result.Valid {
		t.Error("Expected invalid result for non-existent reference")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error for non-existent reference")
	}
}

func TestRefSchema_CircularReference(t *testing.T) {
	ctx := DefaultValidationContext()
	registry := NewSchemaRegistry()

	// This would create a circular reference in a real scenario
	// For testing, we'll just verify the detection works
	ref := Ref("#/Circular", registry)
	registry.Define("Circular", ref) // Self-referencing

	result := ref.Parse("test", ctx)
	if result.Valid {
		t.Error("Expected invalid result for circular reference")
	}
}

// Test Binary Schema
func TestBinarySchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *BinarySchema
		value    interface{}
		expected bool
	}{
		{"valid base64", Base64(), "SGVsbG8gV29ybGQ=", true},      // "Hello World"
		{"valid base64url", Base64URL(), "SGVsbG8gV29ybGQ", true}, // No padding
		{"valid hex", Hex(), "48656c6c6f20576f726c64", true},      // "Hello World"
		{"invalid base64", Base64(), "invalid-base64!", false},
		{"invalid hex", Hex(), "invalid-hex-data", false},
		{"not a string", Base64(), 123, false},
		{"nil", Base64(), nil, false},
		{"empty string", Base64(), "", true}, // Empty is valid if not required
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Binary.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

func TestBinarySchema_Size(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Base64().MinSize(5).MaxSize(100)

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid size", "SGVsbG8gV29ybGQ=", true}, // "Hello World" = 11 bytes
		{"too small", "SGk=", false},             // "Hi" = 2 bytes
		{"empty string", "", true},               // Empty is allowed (not required)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Binary.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test JSON Schema Generation for Advanced Types
func TestAdvancedSchemas_JSON(t *testing.T) {
	tests := []struct {
		name     string
		schema   interface{ JSON() map[string]interface{} }
		expected map[string]interface{}
	}{
		{
			name:     "any schema",
			schema:   Any(),
			expected: map[string]interface{}{
				// Any schema should not have a type field - it accepts everything
			},
		},
		{
			name:   "not schema",
			schema: Not(String()),
			expected: map[string]interface{}{
				"not": map[string]interface{}{"type": "string"},
			},
		},
		{
			name:   "uuid schema",
			schema: UUID(),
			expected: map[string]interface{}{
				"type":   "string",
				"format": "uuid",
			},
		},
		{
			name:   "binary schema",
			schema: Base64(),
			expected: map[string]interface{}{
				"type":            "string",
				"contentEncoding": "base64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.JSON()

			for key, expectedValue := range tt.expected {
				actualValue, exists := result[key]
				if !exists {
					t.Errorf("JSON() missing field %s", key)
					continue
				}

				// Handle different value types
				switch expectedValue := expectedValue.(type) {
				case []interface{}:
					actualSlice, ok := actualValue.([]interface{})
					if !ok {
						t.Errorf("JSON()[%s] type mismatch", key)
						continue
					}
					if len(expectedValue) != len(actualSlice) {
						t.Errorf("JSON()[%s] length mismatch", key)
						continue
					}
					for i, expectedItem := range expectedValue {
						if actualSlice[i] != expectedItem {
							t.Errorf("JSON()[%s][%d] = %v, want %v", key, i, actualSlice[i], expectedItem)
						}
					}
				case map[string]interface{}:
					actualMap, ok := actualValue.(map[string]interface{})
					if !ok {
						t.Errorf("JSON()[%s] type mismatch", key)
						continue
					}
					for nestedKey, nestedExpected := range expectedValue {
						if actualMap[nestedKey] != nestedExpected {
							t.Errorf("JSON()[%s][%s] = %v, want %v", key, nestedKey, actualMap[nestedKey], nestedExpected)
						}
					}
				default:
					if actualValue != expectedValue {
						t.Errorf("JSON()[%s] = %v, want %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}

// Test Complex Combinations
func TestAdvancedSchemas_Complex(t *testing.T) {
	ctx := DefaultValidationContext()

	// Complex schema: AnyOf(AllOf(string constraints), number)
	complexSchema := AnyOf(
		AllOf(
			String(),
			String().MinLength(5),
			String().MaxLength(20),
		),
		Int().Min(100),
	)

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid string all constraints", "hello world", true},
		{"valid number", 150, true},
		{"string too short", "hi", false},
		{"string too long", "this string is way too long for our constraints", false},
		{"number too small", 50, false},
		{"invalid type", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := complexSchema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("ComplexSchema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}
