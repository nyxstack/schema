package schema

import (
	"math"
	"testing"
)

// Test Int Schema
func TestIntSchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Int()

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid int", 42, true},
		{"zero", 0, true},
		{"negative int", -42, true},
		{"max int", math.MaxInt, true},
		{"min int", math.MinInt, true},
		{"float", 3.14, false},
		{"string", "42", false},
		{"boolean", true, false},
		{"nil", nil, false},
		{"array", []int{1, 2, 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Int.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
			}
		})
	}
}

func TestIntSchema_MinMax(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *IntSchema
		value    int
		expected bool
	}{
		{"min valid", Int().Min(10), 15, true},
		{"min invalid", Int().Min(10), 5, false},
		{"min exact", Int().Min(10), 10, true},
		{"max valid", Int().Max(100), 50, true},
		{"max invalid", Int().Max(100), 150, false},
		{"max exact", Int().Max(100), 100, true},
		{"range valid", Int().Min(10).Max(100), 50, true},
		{"range too small", Int().Min(10).Max(100), 5, false},
		{"range too large", Int().Min(10).Max(100), 150, false},
		{"negative range", Int().Min(-100).Max(-10), -50, true},
		{"zero in range", Int().Min(-10).Max(10), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("IntSchema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

func TestIntSchema_Enum(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Int().Enum([]int{1, 2, 3, 5, 8, 13})

	tests := []struct {
		name     string
		value    int
		expected bool
	}{
		{"enum valid 1", 1, true},
		{"enum valid 13", 13, true},
		{"enum invalid", 4, false},
		{"enum negative", -1, false},
		{"enum zero", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("IntSchema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
			}
		})
	}
}

// Test Int8 Schema
func TestInt8Schema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Int8()

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid int8", int8(42), true},
		{"valid int within range", 100, true},
		{"max int8", int8(127), true},
		{"min int8", int8(-128), true},
		{"int too large", 200, false},
		{"int too small", -200, false},
		{"float", 3.14, false},
		{"string", "42", false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Int8.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Int16 Schema
func TestInt16Schema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Int16()

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid int16", int16(1000), true},
		{"valid int within range", 5000, true},
		{"max int16", int16(32767), true},
		{"min int16", int16(-32768), true},
		{"int too large", 40000, false},
		{"int too small", -40000, false},
		{"float", 3.14, false},
		{"string", "1000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Int16.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
			}
		})
	}
}

// Test Int32 Schema
func TestInt32Schema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Int32()

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid int32", int32(1000000), true},
		{"valid int within range", 1000000, true},
		{"max int32", int32(2147483647), true},
		{"min int32", int32(-2147483648), true},
		{"int64 too large", int64(3000000000), false},
		{"int64 too small", int64(-3000000000), false},
		{"float", 3.14, false},
		{"string", "1000000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Int32.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
			}
		})
	}
}

// Test Int64 Schema
func TestInt64Schema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Int64()

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid int64", int64(9223372036854775807), true},
		{"valid int", 1000000, true},
		{"max int64", int64(9223372036854775807), true},
		{"min int64", int64(-9223372036854775808), true},
		{"zero", int64(0), true},
		{"float", 3.14, false},
		{"string", "1000000", false},
		{"boolean", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Int64.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
			}
		})
	}
}

// Test Combined Integer Constraints
func TestIntegerSchemas_Combined(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   Parseable
		value    interface{}
		expected bool
	}{
		{"int with min/max/enum", Int().Min(0).Max(100).Enum([]int{10, 20, 30}), 20, true},
		{"int with constraints invalid", Int().Min(0).Max(100).Enum([]int{10, 20, 30}), 15, false},
		{"int8 with min/max", Int8().Min(0).Max(100), 50, true},
		{"int8 with constraints invalid", Int8().Min(0).Max(100), 150, false},
		{"int16 required", Int16().Required(), int16(1000), true},
		{"int32 const", Int32().Const(12345), int32(12345), true},
		{"int32 const invalid", Int32().Const(12345), int32(54321), false},
		{"int64 default", Int64().Default(999), int64(123), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Schema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test JSON Schema Generation
func TestIntegerSchemas_JSON(t *testing.T) {
	tests := []struct {
		name     string
		schema   Parseable
		expected map[string]interface{}
	}{
		{
			name:   "basic int",
			schema: Int(),
			expected: map[string]interface{}{
				"type": "integer",
			},
		},
		{
			name:   "int with constraints",
			schema: Int().Min(0).Max(100),
			expected: map[string]interface{}{
				"type":    "integer",
				"minimum": 0,
				"maximum": 100,
			},
		},
		{
			name:   "int8 with enum",
			schema: Int8().Enum([]int8{1, 2, 3}),
			expected: map[string]interface{}{
				"type": "integer",
				"enum": []interface{}{int8(1), int8(2), int8(3)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use type assertion to access JSON method
			var result map[string]interface{}
			switch s := tt.schema.(type) {
			case *IntSchema:
				result = s.JSON()
			case *Int8Schema:
				result = s.JSON()
			case *Int16Schema:
				result = s.JSON()
			case *Int32Schema:
				result = s.JSON()
			case *Int64Schema:
				result = s.JSON()
			default:
				t.Fatalf("Unknown schema type: %T", tt.schema)
			}

			for key, expectedValue := range tt.expected {
				actualValue, exists := result[key]
				if !exists {
					t.Errorf("JSON() missing field %s", key)
					continue
				}

				// Handle slice comparison specially for enum
				if key == "enum" {
					expectedSlice, ok1 := expectedValue.([]interface{})
					actualSlice, ok2 := actualValue.([]interface{})
					if !ok1 || !ok2 {
						t.Errorf("JSON()[%s] type mismatch", key)
						continue
					}
					if len(expectedSlice) != len(actualSlice) {
						t.Errorf("JSON()[%s] length mismatch: got %v, want %v", key, actualSlice, expectedSlice)
						continue
					}
					for i, expectedItem := range expectedSlice {
						if actualSlice[i] != expectedItem {
							t.Errorf("JSON()[%s][%d] = %v, want %v", key, i, actualSlice[i], expectedItem)
						}
					}
				} else {
					if actualValue != expectedValue {
						t.Errorf("JSON()[%s] = %v, want %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}

// Test Edge Cases
func TestIntegerSchemas_EdgeCases(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   Parseable
		value    interface{}
		expected bool
	}{
		{"int8 overflow", Int8(), 128, false},
		{"int8 underflow", Int8(), -129, false},
		{"int16 overflow", Int16(), 32768, false},
		{"int16 underflow", Int16(), -32769, false},
		{"int32 overflow", Int32(), int64(2147483648), false},
		{"int32 underflow", Int32(), int64(-2147483649), false},
		{"negative zero", Int(), -0, true},
		{"int min constraint", Int().Min(math.MinInt), math.MinInt, true},
		{"int max constraint", Int().Max(math.MaxInt), math.MaxInt, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Schema.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if !result.Valid && len(result.Errors) > 0 {
					t.Logf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}
