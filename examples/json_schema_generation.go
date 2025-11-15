package main

import (
	"encoding/json"
	"fmt"

	"github.com/nyxstack/schema"
)

func main() {
	fmt.Println("=== JSON Schema Generation ===\n")

	// Create a user schema
	userSchema := schema.Object().
		Title("User").
		Description("A user in the system").
		Property("id", schema.Int().Min(1).Required()).
		Property("username", schema.String().
			MinLength(3).
			MaxLength(20).
			Pattern("^[a-zA-Z0-9_]+$").
			Required()).
		Property("email", schema.String().
			Email().
			Required()).
		Property("age", schema.Int().
			Min(13).
			Max(120).
			Optional()).
		Property("role", schema.String().
			Enum([]string{"admin", "user", "guest"}).
			Default("user")).
		AdditionalProperties(false)

	// Generate JSON Schema
	jsonSchema := userSchema.JSON()

	// Pretty print the JSON Schema
	jsonBytes, err := json.MarshalIndent(jsonSchema, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("User Schema (JSON Schema format):")
	fmt.Println(string(jsonBytes))

	fmt.Println("\n=== Nested Object Schema ===\n")

	// More complex example with nested objects
	postSchema := schema.Object().
		Title("Blog Post").
		Property("title", schema.String().MinLength(5).MaxLength(200).Required()).
		Property("content", schema.String().MinLength(20).Required()).
		Property("author", schema.Object().
			Property("id", schema.Int().Required()).
			Property("name", schema.String().Required()).
			Property("email", schema.String().Email().Required())).
		Property("tags", schema.Array(schema.String()).
			MinItems(1).
			MaxItems(10).
			UniqueItems()).
		Property("published", schema.Bool().Default(false)).
		Property("createdAt", schema.Date().Format(schema.FormatDateTime).Required())

	postJSON := postSchema.JSON()
	postBytes, _ := json.MarshalIndent(postJSON, "", "  ")
	fmt.Println("Blog Post Schema:")
	fmt.Println(string(postBytes))

	fmt.Println("\n=== Array Schema ===\n")

	numbersSchema := schema.Array(
		schema.Int().Min(0).Max(100),
	).MinItems(1).MaxItems(10).UniqueItems().
		Title("Number List").
		Description("A list of unique numbers between 0 and 100")

	numbersJSON := numbersSchema.JSON()
	numbersBytes, _ := json.MarshalIndent(numbersJSON, "", "  ")
	fmt.Println("Numbers Array Schema:")
	fmt.Println(string(numbersBytes))

	fmt.Println("\n=== Union Type Schema (AnyOf) ===\n")

	flexibleIDSchema := schema.AnyOf(
		schema.String().Pattern("^[A-Z]{3}\\d{6}$"),
		schema.Int().Min(100000).Max(999999),
	)

	flexibleJSON := flexibleIDSchema.JSON()
	flexibleBytes, _ := json.MarshalIndent(flexibleJSON, "", "  ")
	fmt.Println("Flexible ID Schema (string or number):")
	fmt.Println(string(flexibleBytes))

	fmt.Println("\n=== Tuple Schema ===\n")

	coordinateSchema := schema.Tuple(
		schema.Float().Min(-90).Max(90),
		schema.Float().Min(-180).Max(180),
	).Title("GPS Coordinate").
		Description("Latitude and longitude pair")

	coordJSON := coordinateSchema.JSON()
	coordBytes, _ := json.MarshalIndent(coordJSON, "", "  ")
	fmt.Println("Coordinate Tuple Schema:")
	fmt.Println(string(coordBytes))
}
