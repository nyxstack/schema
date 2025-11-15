# Nyx Schema Documentation

Complete documentation for all schema types in the Nyx Schema validation library.

## Table of Contents

- [Getting Started](#getting-started)
- [Core Schemas](#core-schemas)
- [Specialized Schemas](#specialized-schemas)
- [Advanced Schemas](#advanced-schemas)
- [Utility Schemas](#utility-schemas)

## Getting Started

Before diving into specific schema types, familiarize yourself with the core concepts:

- **Fluent API**: All schemas support method chaining for intuitive configuration
- **Parse-Don't-Validate**: Use `Parse(value, ctx)` to both validate and transform data
- **Internationalization**: Built-in i18n support through error messages
- **JSON Schema**: Generate JSON Schema output with `.JSON()`

```go
import "github.com/nyxstack/schema"

// Basic validation pattern
ctx := schema.DefaultValidationContext()
result := mySchema.Parse(data, ctx)

if !result.Valid {
    // Handle validation errors
    for _, err := range result.Errors {
        log.Printf("Error: %s", err.Message)
    }
}
```

## Core Schemas

These are the fundamental building blocks for data validation:

### Primitive Types

| Schema | Description | Documentation |
|--------|-------------|---------------|
| **[String](string.md)** | Text validation with formats, patterns, and length constraints | [View →](string.md) |
| **[Int](int.md)** | Integer validation (Int, Int8, Int16, Int32, Int64) with range constraints | [View →](int.md) |
| **[Number](number.md)** | Floating-point validation (Float, Number) with precision control | [View →](number.md) |
| **[Bool](bool.md)** | Boolean validation with const and enum support | [View →](bool.md) |

### Composite Types

| Schema | Description | Documentation |
|--------|-------------|---------------|
| **[Object](object.md)** | Structured data with defined properties and validation rules | [View →](object.md) |
| **[Array](array.md)** | Collections with item validation, length, and uniqueness constraints | [View →](array.md) |
| **[Record](record.md)** | Dynamic key-value maps with validated keys and values | [View →](record.md) |
| **[Tuple](tuple.md)** | Fixed-position arrays where each position has a specific type | [View →](tuple.md) |

## Specialized Schemas

Type-specific validators for common data formats:

| Schema | Description | Documentation |
|--------|-------------|---------------|
| **[UUID](uuid.md)** | UUID validation with version and format support | [View →](uuid.md) |
| **[Date](date.md)** | Date, DateTime, and Time validation with range constraints | [View →](date.md) |
| **[Binary](binary.md)** | Binary data validation (base64, base64url, hex encoding) | [View →](binary.md) |

## Advanced Schemas

Complex validation patterns for sophisticated use cases:

| Schema | Description | Documentation |
|--------|-------------|---------------|
| **[Union](union.md)** | Multiple schema options (OneOf, AnyOf, AllOf) | [View →](union.md) |
| **[Conditional](conditional.md)** | If/then/else validation logic based on conditions | [View →](conditional.md) |
| **[Transform](transform.md)** | Validate input → transform → validate output pipeline | [View →](transform.md) |
| **[Ref](ref.md)** | Schema references for reuse and recursive structures | [View →](ref.md) |

## Utility Schemas

Special-purpose validators:

| Schema | Description | Documentation |
|--------|-------------|---------------|
| **[Any](any.md)** | Accept any value type without validation | [View →](any.md) |
| **[Not](not.md)** | Inverse validation - reject values matching a schema | [View →](not.md) |
| **[Null](null.md)** | Explicit null value validation | [View →](null.md) |

## Quick Reference by Use Case

### User Input Validation
- Form validation → [String](string.md), [Int](int.md), [Bool](bool.md)
- Email addresses → [String](string.md#email-validation)
- Phone numbers → [String](string.md#pattern-validation)
- User registration → [Object](object.md)

### API Validation
- Request bodies → [Object](object.md)
- Query parameters → [String](string.md), [Int](int.md)
- Dynamic responses → [Any](any.md)
- Polymorphic data → [Union](union.md)

### Database Models
- Primary keys → [UUID](uuid.md), [Int](int.md)
- Foreign keys → [Union](union.md) with [Null](null.md)
- Timestamps → [Date](date.md)
- Nullable columns → [Union](union.md)

### File Processing
- File uploads → [Binary](binary.md)
- Image data → [Binary](binary.md)
- CSV parsing → [Transform](transform.md)
- Configuration files → [Object](object.md), [Any](any.md)

### Complex Structures
- Nested objects → [Object](object.md)
- Trees and graphs → [Ref](ref.md)
- Linked lists → [Ref](ref.md)
- Polymorphic types → [Union](union.md)

## Common Patterns

### Nullable Fields

```go
// String OR null
nullableString := schema.OneOf(
    schema.String(),
    schema.Null(),
)
```

[Learn more →](union.md#nullable-fields)

### Conditional Validation

```go
// Different rules based on account type
schema.Conditional(
    schema.Object().Property("type", schema.String().Const("premium")),
).Then(premiumSchema).Else(freeSchema)
```

[Learn more →](conditional.md)

### Data Transformation

```go
// Convert string to integer
schema.Transform(
    schema.String().Pattern("^[0-9]+$"),
    schema.Int().Min(0),
    func(input interface{}) (interface{}, error) {
        return strconv.Atoi(input.(string))
    },
)
```

[Learn more →](transform.md)

### Reusable Schemas

```go
registry := schema.NewSchemaRegistry()
registry.Define("Address", addressSchema)

// Use reference
schema.Ref("#/Address", registry)
```

[Learn more →](ref.md)

## Internationalization

All error messages support internationalization:

```go
// Static message
schema.String().Required(i18n.S("field is required"))

// Formatted message with placeholders
schema.Int().Min(18, i18n.F("age must be at least %d", 18))
```

Each schema documentation page includes i18n examples.

## JSON Schema Generation

All schemas can generate JSON Schema output:

```go
jsonSchema := mySchema.JSON()
jsonBytes, _ := json.MarshalIndent(jsonSchema, "", "  ")
fmt.Println(string(jsonBytes))
```

Perfect for:
- OpenAPI/Swagger documentation
- API documentation generation
- Schema sharing across languages
- Client-side validation

## Error Handling

Consistent error handling across all schemas:

```go
result := schema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Path: %s\n", err.Path)
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        fmt.Printf("Value: %v\n", err.Value)
    }
}
```

## Navigation Tips

- **By Type**: Use the tables above to find schema types
- **By Use Case**: Check the "Quick Reference by Use Case" section
- **Related Schemas**: Each documentation page has a "Related" section linking to similar schemas
- **Search**: Use Ctrl+F to search for specific features or methods

## Additional Resources

- [Main README](../README.md) - Library overview and installation
- [Examples](../examples/) - Working code examples
- [AGENTS.md](../AGENTS.md) - AI agent integration guide

## Contributing

Found an issue or want to improve the documentation? Please submit a PR or open an issue on GitHub.

---

**Need help?** Start with the schema type closest to your use case, or check out the [examples](../examples/) for working code.
