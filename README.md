# NYX Schema

A powerful, fluent Go library for schema validation and JSON Schema generation with internationalization support.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24.2-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Features

- **Fluent API**: Chainable method calls for intuitive schema construction
- **Type Safety**: Strong typing with validation for primitives, objects, arrays, and unions
- **JSON Schema Generation**: Export schemas as standard JSON Schema format
- **Internationalization**: Built-in i18n support for error messages
- **Comprehensive Validation**: Support for all common validation constraints
- **Advanced Schema Types**: AllOf, AnyOf, OneOf, Not, and conditional schemas
- **Format Validation**: Built-in support for email, URI, UUID, date-time, and more
- **Custom Error Messages**: Override default validation messages with custom ones
- **Nullable & Optional**: Fine-grained control over required fields and null values

## Installation

```bash
go get github.com/nyxstack/schema
```

## Documentation

For comprehensive documentation on all schema types, including detailed examples and API references:

**[📚 View Full Documentation →](docs/README.md)**

Quick links to specific schema types:
- [String Schema](docs/string.md) - Text validation with formats and patterns
- [Int/Number Schema](docs/int.md) - Integer and floating-point validation
- [Object Schema](docs/object.md) - Structured data validation
- [Array Schema](docs/array.md) - Collection validation
- [Union Schema](docs/union.md) - Multiple schema options (OneOf/AnyOf/AllOf)
- [Conditional Schema](docs/conditional.md) - If/then/else validation logic
- [UUID Schema](docs/uuid.md) - UUID validation with versions
- [Date Schema](docs/date.md) - Date, DateTime, Time validation
- [Transform Schema](docs/transform.md) - Input transformation and validation
- [Ref Schema](docs/ref.md) - Schema references and reuse

[View all schema types →](docs/README.md)


## Quick Start

```go
package main

import (
    "fmt"
    "github.com/nyxstack/schema"
)

func main() {
    // Create a string schema with validation
    nameSchema := schema.String().
        Title("Full Name").
        Description("User's full name").
        MinLength(2).
        MaxLength(50).
        Pattern("^[a-zA-Z\\s]+$").
        Required()

    // Validate data
    ctx := schema.DefaultValidationContext()
    result := nameSchema.Parse("John Doe", ctx)
    
    if result.Valid {
        fmt.Printf("Valid name: %v\n", result.Value)
    } else {
        for _, err := range result.Errors {
            fmt.Printf("Error: %s\n", err.Message)
        }
    }
}
```

## Core Schema Types

### String Schema

```go
// Basic string validation
emailSchema := schema.String().
    Email().
    Required("Email is required").
    MinLength(5, "Email too short")

// Pattern validation
usernameSchema := schema.String().
    Pattern("^[a-zA-Z0-9_]+$").
    MinLength(3).
    MaxLength(20)

// Enum validation
statusSchema := schema.String().
    Enum([]string{"active", "inactive", "pending"})
```

### Number Schema

```go
// Integer validation
ageSchema := schema.Int().
    Min(0).
    Max(120).
    Required()

// Float validation
priceSchema := schema.Float().
    Min(0.01).
    Max(9999.99).
    MultipleOf(0.01)
```

### Array Schema

```go
// Array of strings
tagsSchema := schema.Array(schema.String().MinLength(1)).
    MinItems(1).
    MaxItems(10).
    UniqueItems()

// Complex array validation
usersSchema := schema.Array(
    schema.Object().
        Property("name", schema.String().Required()).
        Property("email", schema.String().Email().Required()),
).MinItems(1)
```

### Object Schema

```go
// Structured object validation
userSchema := schema.Object().
    Property("id", schema.Int().Min(1).Required()).
    Property("name", schema.String().MinLength(2).Required()).
    Property("email", schema.String().Email().Required()).
    Property("age", schema.Int().Min(13).Max(120).Optional()).
    Property("tags", schema.Array(schema.String()).Optional()).
    AdditionalProperties(false)

// Using Shape for concise object creation
user := schema.Shape{
    "name":  schema.String().Required(),
    "email": schema.String().Email().Required(),
    "age":   schema.Int().Min(0).Optional(),
}.AsObject()
```

### Record Schema

```go
// Record schema for key-value maps with dynamic keys
metadataSchema := schema.Record(
    schema.String().MinLength(1), // Key schema
    schema.String().MinLength(1), // Value schema
).MinProperties(1).MaxProperties(10)

// Example: validate a map of string labels
labelsSchema := schema.Record(
    schema.String().Pattern("^[a-z0-9-]+$"), // Keys must be lowercase with dashes
    schema.String().MaxLength(100),           // Values max 100 chars
)
```

### Union Types

```go
// AnyOf - value matches at least one schema
stringOrNumber := schema.AnyOf(
    schema.String().MinLength(1),
    schema.Int().Min(0),
)

// OneOf - value matches exactly one schema
idSchema := schema.OneOf(
    schema.String().Pattern("^[a-zA-Z0-9_]+$"),
    schema.Int().Min(1),
)

// AllOf - value matches all schemas
restrictiveString := schema.AllOf(
    schema.String().MinLength(8),
    schema.String().Pattern(".*[A-Z].*"),
    schema.String().Pattern(".*[0-9].*"),
)
```

## Advanced Features

### Conditional Validation

```go
conditionalSchema := schema.Conditional().
    If(schema.Object().Property("type", schema.String().Const("premium"))).
    Then(schema.Object().Property("features", schema.Array(schema.String()).MinItems(3))).
    Else(schema.Object().Property("features", schema.Array(schema.String()).MaxItems(1)))
```

### Custom Error Messages

```go
schema := schema.String().
    MinLength(8, "Password must be at least 8 characters").
    Pattern(".*[A-Z].*", "Password must contain uppercase letter").
    Required("Password is mandatory")
```

### Nullable and Optional

```go
// Optional field (can be omitted)
optionalField := schema.String().Optional()

// Nullable field (can be null)
nullableField := schema.String().Nullable()

// Both optional and nullable
flexibleField := schema.String().Optional().Nullable()
```

### Format Validation

```go
// Built-in formats
emailSchema := schema.String().Email()
urlSchema := schema.String().URL()
uuidSchema := schema.String().UUID()
dateSchema := schema.String().DateTime()

// Custom format with pattern
phoneSchema := schema.String().
    Format(schema.StringFormat("phone")).
    Pattern("^\\+?[1-9]\\d{1,14}$")
```

## Validation Context

```go
// Default English context
ctx := schema.DefaultValidationContext()

// Custom locale for internationalization
ctx := schema.NewValidationContext("es") // Spanish

// With Go context
ctx := schema.DefaultValidationContext().
    WithContext(context.Background())
```

## JSON Schema Generation

```go
userSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("email", schema.String().Email().Required())

// Generate JSON Schema
jsonSchema := userSchema.JSON()

// Convert to JSON
jsonBytes, _ := json.Marshal(jsonSchema)
fmt.Println(string(jsonBytes))
```

Output:
```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "email": {
      "type": "string",
      "format": "email"
    }
  },
  "required": ["name", "email"],
  "additionalProperties": false
}
```

## Error Handling

```go
result := schema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Field: %v\n", err.Path)
        fmt.Printf("Value: %s\n", err.Value)
        fmt.Printf("Error: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
    }
}
```

## Real-World Example

```go
// User registration schema with i18n support
registrationSchema := schema.Object().
    Property("username", schema.String().
        Pattern("^[a-zA-Z0-9_]+$").
        MinLength(3).
        MaxLength(20).
        Required(i18n.S("username is required"))).
    Property("email", schema.String().
        Email().
        Required(i18n.S("email address is required"))).
    Property("password", schema.String().
        MinLength(8, i18n.F("password must be at least %d characters", 8)).
        Pattern(".*[A-Z].*", i18n.S("password must contain uppercase")).
        Pattern(".*[0-9].*", i18n.S("password must contain number")).
        Required(i18n.S("password is required"))).
    Property("age", schema.Int().
        Min(13, i18n.F("must be at least %d years old", 13)).
        Max(120).
        Optional()).
    Property("terms", schema.Bool().
        Const(true, i18n.S("must accept terms and conditions")).
        Required()).
    AdditionalProperties(false)

// Validate registration data
data := map[string]interface{}{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123",
    "age": 25,
    "terms": true,
}

result := registrationSchema.Parse(data, schema.DefaultValidationContext())
if result.Valid {
    fmt.Println("Registration data is valid!")
} else {
    for _, err := range result.Errors {
        fmt.Printf("Validation error: %s\n", err.Message)
    }
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Commit your changes (`git commit -am 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [nyxstack/i18n](https://github.com/nyxstack/i18n) - Internationalization support used by this library
- [nyxstack/validator](https://github.com/nyxstack/validator) - Struct tag based validation using this schema library. Enables validation through struct tags for seamless integration with Go structs.

## Support

- Create an issue for bug reports or feature requests
- Check existing issues before creating new ones
- Provide minimal reproducible examples for bugs

---

Made with ❤️ by the Nyx team