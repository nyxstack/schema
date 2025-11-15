package schema

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Test Array Schema Basic Validation
func TestArraySchema_Basic(t *testing.T) {
	ctx := DefaultValidationContext()
	schema := Array(String())

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"valid string array", []string{"hello", "world"}, true},
		{"valid interface array with strings", []interface{}{"hello", "world"}, true},
		{"empty array", []string{}, true},
		{"single item", []string{"hello"}, true},
		{"invalid item type", []interface{}{"hello", 123}, false}, // 123 should fail String validation
		{"not an array", "hello", false},
		{"number", 123, false},
		{"boolean", true, false},
		{"nil", nil, false},
		{"object", map[string]interface{}{"key": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Array Schema with Different Item Types
func TestArraySchema_ItemTypes(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *ArraySchema
		value    interface{}
		expected bool
	}{
		{"int array valid", Array(Int()), []int{1, 2, 3}, true},
		{"int array invalid", Array(Int()), []string{"1", "2", "3"}, false},
		{"bool array valid", Array(Bool()), []bool{true, false, true}, true},
		{"bool array invalid", Array(Bool()), []int{1, 0, 1}, false},
		{"nested array valid", Array(Array(String())), [][]string{{"a", "b"}, {"c", "d"}}, true},
		{"nested array invalid", Array(Array(String())), [][]interface{}{{"a", 123}, {"c", "d"}}, false},
		{"object array valid", Array(Object().Property("name", String())), []map[string]interface{}{{"name": "John"}, {"name": "Jane"}}, true},
		{"object array invalid", Array(Object().Property("name", String())), []map[string]interface{}{{"name": "John"}, {"age": 30}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Array Length Constraints
func TestArraySchema_Length(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *ArraySchema
		value    interface{}
		expected bool
	}{
		// MinItems tests
		{"min items valid", Array(String()).MinItems(2), []string{"a", "b", "c"}, true},
		{"min items exact", Array(String()).MinItems(2), []string{"a", "b"}, true},
		{"min items invalid", Array(String()).MinItems(2), []string{"a"}, false},
		{"min items empty", Array(String()).MinItems(1), []string{}, false},

		// MaxItems tests
		{"max items valid", Array(String()).MaxItems(3), []string{"a", "b"}, true},
		{"max items exact", Array(String()).MaxItems(3), []string{"a", "b", "c"}, true},
		{"max items invalid", Array(String()).MaxItems(3), []string{"a", "b", "c", "d"}, false},
		{"max items empty", Array(String()).MaxItems(2), []string{}, true},

		// Range tests
		{"range valid", Array(String()).MinItems(2).MaxItems(4), []string{"a", "b", "c"}, true},
		{"range min exact", Array(String()).MinItems(2).MaxItems(4), []string{"a", "b"}, true},
		{"range max exact", Array(String()).MinItems(2).MaxItems(4), []string{"a", "b", "c", "d"}, true},
		{"range too few", Array(String()).MinItems(2).MaxItems(4), []string{"a"}, false},
		{"range too many", Array(String()).MinItems(2).MaxItems(4), []string{"a", "b", "c", "d", "e"}, false},

		// Length tests (exact)
		{"exact length valid", Array(String()).Length(3), []string{"a", "b", "c"}, true},
		{"exact length invalid short", Array(String()).Length(3), []string{"a", "b"}, false},
		{"exact length invalid long", Array(String()).Length(3), []string{"a", "b", "c", "d"}, false},
		{"exact length zero", Array(String()).Length(0), []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Array Unique Items
func TestArraySchema_UniqueItems(t *testing.T) {
	ctx := DefaultValidationContext()
	stringSchema := Array(String()).UniqueItems()
	anySchema := Array(Any()).UniqueItems()

	tests := []struct {
		name     string
		schema   *ArraySchema
		value    interface{}
		expected bool
	}{
		{"unique strings", stringSchema, []string{"a", "b", "c"}, true},
		{"duplicate strings", stringSchema, []string{"a", "b", "a"}, false},
		{"empty array", stringSchema, []string{}, true},
		{"single item", stringSchema, []string{"a"}, true},
		{"case sensitive", stringSchema, []string{"Hello", "hello"}, true},
		{"numbers unique", anySchema, []interface{}{1, 2, 3}, true},
		{"numbers duplicate", anySchema, []interface{}{1, 2, 1}, false},
		{"mixed types unique", anySchema, []interface{}{"hello", 42, true}, true},
		{"mixed types duplicate", anySchema, []interface{}{"hello", 42, "hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Array Required/Optional/Nullable
func TestArraySchema_RequiredOptionalNullable(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *ArraySchema
		value    interface{}
		expected bool
	}{
		// Required tests (default)
		{"required with array", Array(String()), []string{"a"}, true},
		{"required with nil", Array(String()), nil, false},

		// Optional tests
		{"optional with array", Array(String()).Optional(), []string{"a"}, true},
		{"optional with nil", Array(String()).Optional(), nil, true},

		// Nullable tests
		{"nullable with array", Array(String()).Nullable(), []string{"a"}, true},
		{"nullable with nil", Array(String()).Nullable(), nil, true},

		// Optional + Nullable
		{"optional nullable with array", Array(String()).Optional().Nullable(), []string{"a"}, true},
		{"optional nullable with nil", Array(String()).Optional().Nullable(), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Array Default Values
func TestArraySchema_DefaultValues(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name          string
		schema        *ArraySchema
		value         interface{}
		expected      bool
		expectedValue interface{}
	}{
		{"default used for nil", Array(String()).Default([]string{"default"}), nil, true, []string{"default"}},
		{"default not used for valid array", Array(String()).Default([]string{"default"}), []string{"actual"}, true, []string{"actual"}},
		{"optional with default", Array(String()).Optional().Default([]string{"default"}), nil, true, []string{"default"}},
		{"nullable with default keeps nil", Array(String()).Nullable().Default([]string{"default"}), nil, true, nil}, // nullable schemas return nil directly
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
				return
			}
			if result.Valid && tt.expectedValue != nil {
				// Convert both values to strings for comparison to avoid reflection issues
				expectedStr := fmt.Sprintf("%v", tt.expectedValue)
				actualStr := fmt.Sprintf("%v", result.Value)
				if expectedStr != actualStr {
					t.Errorf("Array.Parse(%v) value = %v (%T), want %v (%T)", tt.value, result.Value, result.Value, tt.expectedValue, tt.expectedValue)
				}
			}
		})
	}
}

// Test Array Complex Combinations
func TestArraySchema_ComplexCombinations(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *ArraySchema
		value    interface{}
		expected bool
	}{
		{
			"string array with all constraints",
			Array(String().MinLength(3)).MinItems(2).MaxItems(4).UniqueItems(),
			[]string{"hello", "world", "test"},
			true,
		},
		{
			"string array too short item",
			Array(String().MinLength(3)).MinItems(2).MaxItems(4).UniqueItems(),
			[]string{"hello", "hi", "test"},
			false,
		},
		{
			"string array too few items",
			Array(String().MinLength(3)).MinItems(2).MaxItems(4).UniqueItems(),
			[]string{"hello"},
			false,
		},
		{
			"string array not unique",
			Array(String().MinLength(3)).MinItems(2).MaxItems(4).UniqueItems(),
			[]string{"hello", "world", "hello"},
			false,
		},
		{
			"integer array with constraints",
			Array(Int().Min(0).Max(100)).Length(3),
			[]int{10, 50, 90},
			true,
		},
		{
			"integer array out of range",
			Array(Int().Min(0).Max(100)).Length(3),
			[]int{10, 150, 90},
			false,
		},
		{
			"nested array validation",
			Array(Array(String().MinLength(2)).MinItems(1)),
			[][]string{{"hello", "world"}, {"test", "case"}},
			true,
		},
		{
			"nested array invalid inner",
			Array(Array(String().MinLength(2)).MinItems(1)),
			[][]string{{"hello", "world"}, {"a"}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Array JSON Schema Generation
func TestArraySchema_JSON(t *testing.T) {
	tests := []struct {
		name           string
		schema         *ArraySchema
		expectedFields map[string]interface{}
	}{
		{
			"basic array",
			Array(String()),
			map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{"type": "string"},
			},
		},
		{
			"array with constraints",
			Array(String()).MinItems(2).MaxItems(5).UniqueItems(),
			map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"minItems":    2,
				"maxItems":    5,
				"uniqueItems": true,
			},
		},
		{
			"nullable array",
			Array(String()).Nullable(),
			map[string]interface{}{
				"type":  []string{"array", "null"},
				"items": map[string]interface{}{"type": "string"},
			},
		},
		{
			"array with title and description",
			Array(String()).Title("Names").Description("List of names"),
			map[string]interface{}{
				"type":        "array",
				"title":       "Names",
				"description": "List of names",
				"items":       map[string]interface{}{"type": "string"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.JSON()

			for key, expected := range tt.expectedFields {
				actual, exists := result[key]
				if !exists {
					t.Errorf("Expected field %s not found in JSON", key)
					continue
				}

				// Special handling for slice comparisons
				if expectedSlice, ok := expected.([]interface{}); ok {
					if actualSlice, ok := actual.([]interface{}); ok {
						if len(expectedSlice) != len(actualSlice) {
							t.Errorf("Field %s length = %d, want %d", key, len(actualSlice), len(expectedSlice))
							continue
						}
						allMatch := true
						for i, expectedItem := range expectedSlice {
							if actualSlice[i] != expectedItem {
								allMatch = false
								break
							}
						}
						if !allMatch {
							t.Errorf("Field %s = %v, want %v", key, actual, expected)
						}
						continue
					} else {
						t.Errorf("Field %s: expected slice but got %T: %v", key, actual, actual)
						continue
					}
				}

				// Handle []string comparisons
				if expectedStringSlice, ok := expected.([]string); ok {
					if actualStringSlice, ok := actual.([]string); ok {
						if !reflect.DeepEqual(actualStringSlice, expectedStringSlice) {
							t.Errorf("Field %s = %v, want %v", key, actual, expected)
						}
						continue
					} else {
						t.Errorf("Field %s: expected string slice but got %T: %v", key, actual, actual)
						continue
					}
				}

				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Field %s = %v, want %v", key, actual, expected)
				}
			}
		})
	}
}

// Test Array Edge Cases
func TestArraySchema_EdgeCases(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name     string
		schema   *ArraySchema
		value    interface{}
		expected bool
	}{
		// Empty arrays
		{"empty array with min items 0", Array(String()).MinItems(0), []string{}, true},
		{"empty array with unique items", Array(String()).UniqueItems(), []string{}, true},

		// Large arrays
		{"large array", Array(String()), make([]string, 1000), true},

		// Mixed type handling
		{"interface slice with consistent types", Array(Any()), []interface{}{"string", 123, true, nil}, true},

		// Nil item schemas edge case
		{"array with nil items allowed", Array(Any()), []interface{}{"valid", nil, "items"}, true},

		// Zero constraints
		{"max items zero", Array(String()).MaxItems(0), []string{}, true},
		{"max items zero with content", Array(String()).MaxItems(0), []string{"item"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create large array for the large array test
			if tt.name == "large array" {
				largeArray := make([]string, 1000)
				for i := 0; i < 1000; i++ {
					largeArray[i] = fmt.Sprintf("item%d", i)
				}
				tt.value = largeArray
			}

			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid != tt.expected {
				t.Errorf("Array.Parse(%v) = %v, want %v", tt.value, result.Valid, tt.expected)
				if len(result.Errors) > 0 {
					t.Errorf("Error: %s", result.Errors[0].Message)
				}
			}
		})
	}
}

// Test Array Error Messages
func TestArraySchema_ErrorMessages(t *testing.T) {
	ctx := DefaultValidationContext()

	tests := []struct {
		name            string
		schema          *ArraySchema
		value           interface{}
		expectedContain string
	}{
		{"type error", Array(String()), "not an array", "must be an array"},
		{"min items error", Array(String()).MinItems(3), []string{"a", "b"}, "at least 3 items"},
		{"max items error", Array(String()).MaxItems(2), []string{"a", "b", "c"}, "at most 2 items"},
		{"unique items error", Array(String()).UniqueItems(), []string{"a", "b", "a"}, "unique items"},
		{"item validation error", Array(String().MinLength(3)), []string{"hello", "hi"}, "invalid"},
		{"required error", Array(String()), nil, "required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Parse(tt.value, ctx)
			if result.Valid {
				t.Errorf("Expected validation to fail for %s", tt.name)
				return
			}
			if len(result.Errors) == 0 {
				t.Errorf("Expected error message for %s", tt.name)
				return
			}
			errorMsg := result.Errors[0].Message
			if !contains(errorMsg, tt.expectedContain) {
				t.Errorf("Error message '%s' should contain '%s'", errorMsg, tt.expectedContain)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
