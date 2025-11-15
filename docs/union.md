# Union Schema (AnyOf/OneOf/AllOf)

Union schemas provide validation for values that can match one or multiple schemas. Nyx Schema supports three types of union validation.

## Union Types

### OneOf (Union)
Value must match **exactly one** of the provided schemas.

### AnyOf  
Value must match **at least one** of the provided schemas.

### AllOf
Value must match **all** of the provided schemas.

## Creating Union Schemas

```go
import "github.com/nyxstack/schema"

// OneOf: string OR number
idSchema := schema.OneOf(
    schema.String().Pattern("^[a-zA-Z0-9]+$"),
    schema.Int().Min(1),
)

// AnyOf: matches at least one
flexibleSchema := schema.AnyOf(
    schema.String().MinLength(1),
    schema.Int().Min(0),
)

// AllOf: must match all
strictStringSchema := schema.AllOf(
    schema.String().MinLength(8),
    schema.String().Pattern(".*[A-Z].*"),
    schema.String().Pattern(".*[0-9].*"),
)
```

## OneOf / Union Methods

### Type Configuration

#### `Required(messages ...ErrorMessage) *UnionSchema`
Marks the union as required.

```go
schema.OneOf(schema.String(), schema.Int()).Required()
schema.OneOf(schema.String(), schema.Int()).Required(i18n.S("value is required"))
```

#### `Optional() *UnionSchema`
Marks the union as optional.

```go
schema.OneOf(schema.String(), schema.Int()).Optional()
```

#### `Nullable() *UnionSchema`
Allows null values.

```go
schema.OneOf(schema.String(), schema.Int()).Nullable()
```

### Schema Manipulation

#### `Add(schemas ...Parseable) *UnionSchema`
Adds additional schemas to the union.

```go
idSchema := schema.OneOf(schema.String(), schema.Int())
idSchema.Add(schema.String().UUID())
```

### Metadata

#### `Title(title string) *UnionSchema`
Sets a title.

```go
schema.OneOf(schema.String(), schema.Int()).Title("User ID")
```

#### `Description(description string) *UnionSchema`
Sets a description.

```go
schema.OneOf(schema.String(), schema.Int()).Description("ID can be string or integer")
```

## Usage Examples

### OneOf: String or Integer ID

```go
idSchema := schema.OneOf(
    schema.String().Pattern("^[a-zA-Z0-9_]+$"),
    schema.Int().Min(1),
)

ctx := schema.DefaultValidationContext()
result := idSchema.Parse("user_123", ctx) // Valid - matches string
result := idSchema.Parse(456, ctx)        // Valid - matches int
result := idSchema.Parse("", ctx)         // Invalid - matches neither
```

### AnyOf: Multiple Valid Formats

```go
dateSchema := schema.AnyOf(
    schema.String().DateTime(), // ISO 8601
    schema.Int().Min(0),        // Unix timestamp
    schema.String().Date(),     // YYYY-MM-DD
)

result := dateSchema.Parse("2025-11-16T10:30:00Z", ctx) // Valid
result := dateSchema.Parse(1700000000, ctx)             // Valid
result := dateSchema.Parse("2025-11-16", ctx)           // Valid
```

### AllOf: Combine Multiple Constraints

```go
passwordSchema := schema.AllOf(
    schema.String().MinLength(8),
    schema.String().Pattern(".*[A-Z].*"), // Has uppercase
    schema.String().Pattern(".*[a-z].*"), // Has lowercase
    schema.String().Pattern(".*[0-9].*"), // Has number
)

result := passwordSchema.Parse("SecurePass123", ctx) // Valid - matches all
result := passwordSchema.Parse("weak", ctx)          // Invalid - fails length and patterns
```

### Phone Number Formats

```go
phoneSchema := schema.AnyOf(
    schema.String().Pattern("^\\d{10}$"),            // 1234567890
    schema.String().Pattern("^\\+1\\d{10}$"),        // +11234567890
    schema.String().Pattern("^\\(\\d{3}\\) \\d{3}-\\d{4}$"), // (123) 456-7890
)

result := phoneSchema.Parse("5551234567", ctx)
result := phoneSchema.Parse("+15551234567", ctx)
result := phoneSchema.Parse("(555) 123-4567", ctx)
```

### Polymorphic API Responses

```go
responseSchema := schema.OneOf(
    schema.Object().
        Property("success", schema.Bool().Const(true)).
        Property("data", schema.Any()),
    schema.Object().
        Property("success", schema.Bool().Const(false)).
        Property("error", schema.String()),
)

// Success response
result := responseSchema.Parse(map[string]interface{}{
    "success": true,
    "data": map[string]interface{}{"id": 1},
}, ctx)

// Error response
result := responseSchema.Parse(map[string]interface{}{
    "success": false,
    "error": "Not found",
}, ctx)
```

### Flexible Configuration Values

```go
configValueSchema := schema.AnyOf(
    schema.String(),
    schema.Int(),
    schema.Bool(),
    schema.Array(schema.String()),
)

result := configValueSchema.Parse("localhost", ctx)
result := configValueSchema.Parse(8080, ctx)
result := configValueSchema.Parse(true, ctx)
result := configValueSchema.Parse([]string{"opt1", "opt2"}, ctx)
```

### Strong Password Validation (AllOf)

```go
strongPasswordSchema := schema.AllOf(
    schema.String().MinLength(12, i18n.F("must be at least %d characters", 12)),
    schema.String().Pattern(".*[A-Z].*", i18n.S("must contain uppercase letter")),
    schema.String().Pattern(".*[a-z].*", i18n.S("must contain lowercase letter")),
    schema.String().Pattern(".*[0-9].*", i18n.S("must contain number")),
    schema.String().Pattern(".*[!@#$%^&*].*", i18n.S("must contain special character")),
)

result := strongPasswordSchema.Parse("MyS3cur3P@ssword!", ctx)
```

### Mixed Identifier Types

```go
identifierSchema := schema.OneOf(
    schema.String().UUID(),               // UUID
    schema.String().Pattern("^[A-Z]{3}\\d{6}$"), // ABC123456
    schema.Int().Min(1),                  // Numeric ID
).Title("Identifier")

result := identifierSchema.Parse("550e8400-e29b-41d4-a716-446655440000", ctx)
result := identifierSchema.Parse("ABC123456", ctx)
result := identifierSchema.Parse(42, ctx)
```

### Nullable Union

```go
optionalIdSchema := schema.OneOf(
    schema.String(),
    schema.Int(),
).Nullable()

result := optionalIdSchema.Parse(nil, ctx) // Valid
result := optionalIdSchema.Parse("abc", ctx) // Valid
result := optionalIdSchema.Parse(123, ctx) // Valid
```

## Error Handling

```go
result := schema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // OneOf errors: "value matches multiple schemas" or "value does not match any"
        // AnyOf errors: "value does not match any of the allowed schemas"
        // AllOf errors: Shows specific validation failures from each schema
    }
}
```

## Internationalization

```go
schema.OneOf(
    schema.String(),
    schema.Int(),
).Required(i18n.S("value is required"))
```

## JSON Schema Generation

```go
schema := schema.OneOf(
    schema.String(),
    schema.Int(),
).Title("User ID")

jsonSchema := schema.JSON()
// Outputs:
// {
//   "oneOf": [
//     {"type": "string"},
//     {"type": "integer"}
//   ],
//   "title": "User ID"
// }
```

## Related

- [String Schema](string.md) - For string union options
- [Int Schema](int.md) - For integer union options
- [Object Schema](object.md) - For polymorphic objects
- [Conditional Schema](conditional.md) - For if/then/else logic
