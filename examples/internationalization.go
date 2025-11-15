package main

import (
	"fmt"

	"github.com/nyxstack/schema"
)

func main() {
	// Example showing i18n support for different locales

	emailSchema := schema.String().
		Email().
		Required("Email is required").
		MinLength(5)

	fmt.Println("=== English Validation Errors ===")
	ctxEN := schema.NewValidationContext("en")
	result := emailSchema.Parse("bad", ctxEN)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("EN Error: %s\n", err.Message)
		}
	}

	fmt.Println("\n=== With Custom Error Messages ===")
	customSchema := schema.Object().
		Property("username", schema.String().
			MinLength(3, "Username must be at least 3 characters").
			Pattern("^[a-zA-Z0-9_]+$", "Username can only contain letters, numbers, and underscores").
			Required("Username is required")).
		Property("password", schema.String().
			MinLength(8, "Password must be at least 8 characters").
			Pattern(".*[A-Z].*", "Password must contain at least one uppercase letter").
			Required("Password is required"))

	invalidData := map[string]interface{}{
		"username": "ab",
		"password": "short",
	}

	result = customSchema.Parse(invalidData, ctxEN)
	if !result.Valid {
		fmt.Println("Validation errors:")
		for _, err := range result.Errors {
			if len(err.Path) > 0 {
				fmt.Printf("  Field '%s': %s\n", err.Path[0], err.Message)
			} else {
				fmt.Printf("  %s\n", err.Message)
			}
		}
	}

	fmt.Println("\n=== Default Values and Optional Fields ===")
	configSchema := schema.Object().
		Property("host", schema.String().Default("localhost")).
		Property("port", schema.Int().Min(1024).Max(65535).Default(8080)).
		Property("debug", schema.Bool().Default(false)).
		Property("apiKey", schema.String().Optional())

	// Empty object gets defaults
	emptyConfig := map[string]interface{}{}
	result = configSchema.Parse(emptyConfig, ctxEN)
	fmt.Printf("Valid with defaults: %v\n", result.Valid)
	if result.Valid {
		if configMap, ok := result.Value.(map[string]interface{}); ok {
			fmt.Printf("  host: %v\n", configMap["host"])
			fmt.Printf("  port: %v\n", configMap["port"])
			fmt.Printf("  debug: %v\n", configMap["debug"])
		}
	}

	fmt.Println("\n=== Enum Validation ===")
	statusSchema := schema.String().
		Enum([]string{"pending", "approved", "rejected"}).
		Required()

	result = statusSchema.Parse("approved", ctxEN)
	fmt.Printf("Valid status: %v\n", result.Valid)

	result = statusSchema.Parse("invalid", ctxEN)
	fmt.Printf("Invalid status: %v\n", result.Valid)
	if !result.Valid {
		for _, err := range result.Errors {
			fmt.Printf("  Error: %s\n", err.Message)
		}
	}

	fmt.Println("\n=== Nullable Fields ===")
	nullableSchema := schema.Object().
		Property("requiredField", schema.String().Required()).
		Property("nullableField", schema.String().Nullable()).
		Property("optionalField", schema.String().Optional())

	dataWithNull := map[string]interface{}{
		"requiredField": "value",
		"nullableField": nil, // This is OK
		// optionalField is missing - also OK
	}

	result = nullableSchema.Parse(dataWithNull, ctxEN)
	fmt.Printf("Valid with null: %v\n", result.Valid)
}
