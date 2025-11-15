# NYXStack Schema Examples

This directory contains practical examples demonstrating how to use the NYXStack Schema validation library.

## Examples

### 1. Basic Validation (`basic_validation.go`)

Demonstrates fundamental validation patterns:
- String validation with email format and length constraints
- Integer validation with range constraints
- Object validation with required and optional properties
- Array validation with uniqueness and size constraints

**Run:**
```bash
go run basic_validation.go
```

### 2. Advanced Schemas (`advanced_schemas.go`)

Shows complex schema patterns:
- **Union Types (AnyOf)**: Accept multiple different types (string OR number)
- **AllOf Validation**: Must satisfy ALL constraints simultaneously
- **Conditional Validation**: Different rules based on field values
- **Tuple Validation**: Fixed-length arrays with position-specific types
- **Record/Map Validation**: Key-value maps with schema for both keys and values
- **Transform Schema**: Validate, transform, and re-validate data
- **Nested Complex Schema**: Real-world API request validation

**Run:**
```bash
go run advanced_schemas.go
```

### 3. Internationalization (`internationalization.go`)

Demonstrates i18n support:
- Using different locales for error messages
- Custom error messages for better UX
- Default values and optional fields
- Enum validation
- Nullable vs Optional fields

**Run:**
```bash
go run internationalization.go
```

### 4. JSON Schema Generation (`json_schema_generation.go`)

Shows how to generate JSON Schema output:
- Converting Go schemas to JSON Schema format
- Nested object schemas
- Array schemas with constraints
- Union types (anyOf)
- Tuple schemas

**Run:**
```bash
go run json_schema_generation.go
```

## Building All Examples

```bash
# Build all examples
go build -o build/ ./...

# Or build individually
go build -o build/basic_validation basic_validation.go
go build -o build/advanced_schemas advanced_schemas.go
go build -o build/internationalization internationalization.go
go build -o build/json_schema_generation json_schema_generation.go
```

## Common Patterns

### Creating a Schema
```go
userSchema := schema.Object().
    Property("name", schema.String().MinLength(2).Required()).
    Property("email", schema.String().Email().Required()).
    Property("age", schema.Int().Min(18).Optional())
```

### Validating Data
```go
ctx := schema.DefaultValidationContext()
result := userSchema.Parse(userData, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Error at %v: %s\n", err.Path, err.Message)
    }
}
```

### Using Custom Error Messages
```go
emailSchema := schema.String().
    Email().
    Required("Email address is required").
    MinLength(5, "Email must be at least 5 characters")
```

### Working with Different Locales
```go
// English
ctxEN := schema.NewValidationContext("en")
result := schema.Parse(data, ctxEN)

// French (if translations available)
ctxFR := schema.NewValidationContext("fr")
result := schema.Parse(data, ctxFR)
```

## See Also

- [Main Documentation](../README.md)
- [AI Agents Guide](../AGENTS.md)
- [API Reference](https://pkg.go.dev/github.com/nyxstack/schema)
