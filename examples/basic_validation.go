package main

import (
	"fmt"

	"github.com/nyxstack/schema"
)

func main() {
	ctx := schema.DefaultValidationContext()

	fmt.Println("=== Basic String Validation ===")

	// Create a string schema with constraints
	emailSchema := schema.String().
		Email().
		Required("Email is required").
		MinLength(5, "Email must be at least 5 characters")

	// Valid email
	result := emailSchema.Parse("user@example.com", ctx)
	fmt.Printf("Valid email: %v, Value: %v\n", result.Valid, result.Value)

	// Invalid email
	result = emailSchema.Parse("invalid-email", ctx)
	fmt.Printf("Invalid email: %v\n", result.Valid)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("  Error: %s (code: %s)\n", err.Message, err.Code)
		}
	}

	// Required field missing
	result = emailSchema.Parse(nil, ctx)
	fmt.Printf("Missing email: %v\n", result.Valid)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("  Error: %s\n", err.Message)
		}
	}

	fmt.Println("\n=== Integer Validation ===")

	// Integer with range constraints
	ageSchema := schema.Int().
		Min(0, "Age must be positive").
		Max(150, "Age must be realistic").
		Required()

	result = ageSchema.Parse(25, ctx)
	fmt.Printf("Valid age: %v, Value: %v\n", result.Valid, result.Value)

	result = ageSchema.Parse(200, ctx)
	fmt.Printf("Invalid age: %v\n", result.Valid)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("  Error: %s\n", err.Message)
		}
	}

	fmt.Println("\n=== Object Validation ===")

	// Object schema with properties
	userSchema := schema.Object().
		Property("name", schema.String().MinLength(2).Required()).
		Property("email", schema.String().Email().Required()).
		Property("age", schema.Int().Min(18).Optional()).
		AdditionalProperties(false)

	// Valid user
	validUser := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}
	result = userSchema.Parse(validUser, ctx)
	fmt.Printf("Valid user: %v\n", result.Valid)

	// Invalid user - missing required field
	invalidUser := map[string]interface{}{
		"name": "Jane",
	}
	result = userSchema.Parse(invalidUser, ctx)
	fmt.Printf("Invalid user: %v\n", result.Valid)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("  Error at %v: %s\n", err.Path, err.Message)
		}
	}

	// Invalid user - extra field not allowed
	userWithExtra := map[string]interface{}{
		"name":    "Bob",
		"email":   "bob@example.com",
		"unknown": "field",
	}
	result = userSchema.Parse(userWithExtra, ctx)
	fmt.Printf("User with extra fields: %v\n", result.Valid)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("  Error: %s\n", err.Message)
		}
	}

	fmt.Println("\n=== Array Validation ===")

	// Array of strings with constraints
	tagsSchema := schema.Array(
		schema.String().MinLength(2),
	).MinItems(1).MaxItems(5).UniqueItems()

	result = tagsSchema.Parse([]string{"go", "validation", "schema"}, ctx)
	fmt.Printf("Valid tags: %v, Value: %v\n", result.Valid, result.Value)

	result = tagsSchema.Parse([]string{"go", "go", "duplicate"}, ctx)
	fmt.Printf("Duplicate tags: %v\n", result.Valid)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("  Error: %s\n", err.Message)
		}
	}
}
