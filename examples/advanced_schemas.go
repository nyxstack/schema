package main

import (
	"fmt"

	"github.com/nyxstack/schema"
)

func main() {
	ctx := schema.DefaultValidationContext()

	fmt.Println("=== Union Types (AnyOf) ===")

	// Field that accepts either string or number
	idSchema := schema.AnyOf(
		schema.String().Pattern("^[A-Z]{3}\\d{3}$"),
		schema.Int().Min(1000).Max(9999),
	)

	result := idSchema.Parse("ABC123", ctx)
	fmt.Printf("Valid string ID: %v, Value: %v\n", result.Valid, result.Value)

	result = idSchema.Parse(5000, ctx)
	fmt.Printf("Valid numeric ID: %v, Value: %v\n", result.Valid, result.Value)

	result = idSchema.Parse("invalid", ctx)
	fmt.Printf("Invalid ID: %v\n", result.Valid)

	fmt.Println("\n=== AllOf Validation ===")

	// Must satisfy ALL constraints
	strictStringSchema := schema.AllOf(
		schema.String().MinLength(8),
		schema.String().MaxLength(20),
		schema.String().Pattern("^[a-zA-Z0-9]+$"),
	)

	result = strictStringSchema.Parse("Valid123", ctx)
	fmt.Printf("Valid strict string: %v\n", result.Valid)

	result = strictStringSchema.Parse("short", ctx)
	fmt.Printf("Too short: %v\n", result.Valid)

	fmt.Println("\n=== Conditional Validation ===")

	// Different validation based on type
	paymentSchema := schema.Conditional(
		schema.Object().Property("type", schema.String().Const("credit_card")),
	).
		Then(schema.Object().
			Property("cardNumber", schema.String().MinLength(16).Required()).
			Property("cvv", schema.String().MinLength(3).Required())).
		Else(schema.Object().
			Property("accountNumber", schema.String().Required()))

	creditCard := map[string]interface{}{
		"type":       "credit_card",
		"cardNumber": "1234567890123456",
		"cvv":        "123",
	}
	result = paymentSchema.Parse(creditCard, ctx)
	fmt.Printf("Valid credit card: %v\n", result.Valid)

	bankAccount := map[string]interface{}{
		"type":          "bank_transfer",
		"accountNumber": "12345678",
	}
	result = paymentSchema.Parse(bankAccount, ctx)
	fmt.Printf("Valid bank account: %v\n", result.Valid)

	fmt.Println("\n=== Tuple Validation ===")

	// Fixed-length array with position-specific types
	coordinateSchema := schema.Tuple(
		schema.Float().Min(-90).Max(90),   // latitude
		schema.Float().Min(-180).Max(180), // longitude
	)

	result = coordinateSchema.Parse([]interface{}{40.7128, -74.0060}, ctx)
	fmt.Printf("Valid coordinates: %v, Value: %v\n", result.Valid, result.Value)

	result = coordinateSchema.Parse([]interface{}{100.0, -74.0}, ctx)
	fmt.Printf("Invalid latitude: %v\n", result.Valid)

	fmt.Println("\n=== Record (Map) Validation ===")

	// Key-value map with schema for both
	configSchema := schema.Record(
		schema.String().Pattern("^[a-z_]+$"), // keys must be lowercase_snake_case
		schema.String().MinLength(1),         // values must be non-empty strings
	).MinProperties(1).MaxProperties(10)

	config := map[string]interface{}{
		"api_key":  "secret123",
		"base_url": "https://api.example.com",
		"timeout":  "30s",
	}
	result = configSchema.Parse(config, ctx)
	fmt.Printf("Valid config: %v\n", result.Valid)

	invalidConfig := map[string]interface{}{
		"ApiKey":  "bad", // uppercase not allowed
		"baseUrl": "ok",
	}
	result = configSchema.Parse(invalidConfig, ctx)
	fmt.Printf("Invalid config keys: %v\n", result.Valid)

	fmt.Println("\n=== Transform Schema ===")

	// Transform and validate data
	upperCaseTransform := schema.Transform(
		schema.String().MinLength(3),
		schema.String().Pattern("^UPPER_"),
		func(input interface{}) (interface{}, error) {
			if str, ok := input.(string); ok {
				return fmt.Sprintf("UPPER_%s", str), nil
			}
			return nil, fmt.Errorf("not a string")
		},
	)

	result = upperCaseTransform.Parse("hello", ctx)
	fmt.Printf("Transformed: %v, Value: %v\n", result.Valid, result.Value)

	fmt.Println("\n=== Nested Complex Schema ===")

	// Real-world example: API request validation
	createPostSchema := schema.Object().
		Property("title", schema.String().MinLength(5).MaxLength(100).Required()).
		Property("content", schema.String().MinLength(20).Required()).
		Property("tags", schema.Array(schema.String().MinLength(2)).MaxItems(5).Optional()).
		Property("author", schema.Object().
			Property("id", schema.Int().Min(1).Required()).
			Property("name", schema.String().Required()).
			AdditionalProperties(false)).
		Property("published", schema.Bool().Default(false)).
		AdditionalProperties(false)

	validPost := map[string]interface{}{
		"title":   "Introduction to Go Schema Validation",
		"content": "This is a detailed post about schema validation in Go...",
		"tags":    []string{"go", "validation", "schema"},
		"author": map[string]interface{}{
			"id":   123,
			"name": "John Doe",
		},
	}

	result = createPostSchema.Parse(validPost, ctx)
	fmt.Printf("\nValid post: %v\n", result.Valid)
	if result.Valid {
		fmt.Printf("Parsed value: %v\n", result.Value)
	}
}
