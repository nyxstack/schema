# Not Schema

The `NotSchema` validates that a value does NOT match a specified schema. It's the inverse of normal validation - succeeding when the inner schema fails.

## Creating a Not Schema

```go
import "github.com/nyxstack/schema"

// Reject values matching the schema
notSchema := schema.Not(innerSchema)

// With custom error message
notSchema := schema.Not(innerSchema).
    NotError(i18n.S("value must not match the pattern"))
```

## Methods

### Core Methods

#### `Not(schema Parseable) *NotSchema`
Creates a new NOT schema that rejects values matching the given schema.

```go
notString := schema.Not(schema.String())
```

### Error Customization

#### `NotError(err ErrorMessage) *NotSchema`
Sets custom error message for when the value matches (and shouldn't).

```go
schema.Not(schema.String()).
    NotError(i18n.S("value must not be a string"))
```

## Usage Examples

### Reject Strings

```go
notStringSchema := schema.Not(schema.String())

ctx := schema.DefaultValidationContext()

// Valid: not strings
result := notStringSchema.Parse(42, ctx)         // Valid (int)
result = notStringSchema.Parse(true, ctx)        // Valid (bool)
result = notStringSchema.Parse([]int{1, 2}, ctx) // Valid (array)

// Invalid: strings
result = notStringSchema.Parse("hello", ctx)     // Invalid
```

### Reject Empty Strings

```go
notEmptySchema := schema.Not(schema.String().Const(""))

// Valid: non-empty values
result := notEmptySchema.Parse("hello", ctx)     // Valid
result = notEmptySchema.Parse(42, ctx)           // Valid

// Invalid: empty string
result = notEmptySchema.Parse("", ctx)           // Invalid
```

### Reject Specific Values

```go
notZeroSchema := schema.Not(schema.Int().Const(0))

result := notZeroSchema.Parse(1, ctx)            // Valid
result = notZeroSchema.Parse(-5, ctx)            // Valid
result = notZeroSchema.Parse(0, ctx)             // Invalid
```

### Reject Null

```go
notNullSchema := schema.Not(schema.Null())

result := notNullSchema.Parse("value", ctx)      // Valid
result = notNullSchema.Parse(42, ctx)            // Valid
result = notNullSchema.Parse(nil, ctx)           // Invalid
```

### Reject Pattern

```go
// Reject emails
notEmailSchema := schema.Not(schema.String().Email())

result := notEmailSchema.Parse("username", ctx)           // Valid
result = notEmailSchema.Parse("not-an-email", ctx)        // Valid
result = notEmailSchema.Parse("user@example.com", ctx)    // Invalid
```

### Reject Number Range

```go
// Reject numbers between 10 and 20
notInRangeSchema := schema.Not(
    schema.Int().Min(10).Max(20),
)

result := notInRangeSchema.Parse(5, ctx)         // Valid (< 10)
result = notInRangeSchema.Parse(25, ctx)         // Valid (> 20)
result = notInRangeSchema.Parse(15, ctx)         // Invalid (in range)
```

### Username Validation (No Reserved Words)

```go
// Reject reserved usernames
notReservedSchema := schema.Not(
    schema.String().Enum([]string{"admin", "root", "system"}),
)

usernameSchema := schema.String().
    MinLength(3).
    MaxLength(20)

// Combine both validations
result := notReservedSchema.Parse("john", ctx)   // Valid
result = notReservedSchema.Parse("admin", ctx)   // Invalid (reserved)
```

### Reject Objects with Specific Property

```go
// Reject objects that have an "internal" property
notInternalSchema := schema.Not(
    schema.Object().Property("internal", schema.Any()),
)

result := notInternalSchema.Parse(map[string]interface{}{
    "name": "Public Data",
}, ctx) // Valid

result = notInternalSchema.Parse(map[string]interface{}{
    "name":     "Internal Data",
    "internal": true,
}, ctx) // Invalid
```

### Reject Boolean True

```go
notTrueSchema := schema.Not(schema.Bool().Const(true))

result := notTrueSchema.Parse(false, ctx)        // Valid
result = notTrueSchema.Parse(nil, ctx)           // Valid
result = notTrueSchema.Parse(true, ctx)          // Invalid
```

### Complex Rejection (Nested)

```go
// Reject arrays of strings
notStringArraySchema := schema.Not(
    schema.Array(schema.String()),
)

result := notStringArraySchema.Parse([]int{1, 2, 3}, ctx)              // Valid (int array)
result = notStringArraySchema.Parse(42, ctx)                            // Valid (not array)
result = notStringArraySchema.Parse([]string{"a", "b"}, ctx)           // Invalid (string array)
```

### API Version Validation

```go
// Reject deprecated API versions
notDeprecatedSchema := schema.Not(
    schema.String().Enum([]string{"v1", "v2"}),
)

apiRequestSchema := schema.Object().
    Property("version", notDeprecatedSchema).
    Property("data", schema.Any())

result := apiRequestSchema.Parse(map[string]interface{}{
    "version": "v3",
    "data":    "request",
}, ctx) // Valid

result = apiRequestSchema.Parse(map[string]interface{}{
    "version": "v1",
    "data":    "request",
}, ctx) // Invalid (deprecated version)
```

### Port Range Validation

```go
// Reject well-known ports (0-1023)
notWellKnownPortSchema := schema.Not(
    schema.Int().Min(0).Max(1023),
)

portSchema := schema.Int().Min(0).Max(65535)

// Valid ports
result := notWellKnownPortSchema.Parse(8080, ctx)  // Valid
result = notWellKnownPortSchema.Parse(3000, ctx)   // Valid

// Invalid (well-known port)
result = notWellKnownPortSchema.Parse(80, ctx)     // Invalid
result = notWellKnownPortSchema.Parse(443, ctx)    // Invalid
```

## When to Use

Not schemas are ideal for:

- **Blacklisting**: Reject specific values or patterns
- **Reserved Keywords**: Prevent use of reserved names
- **Deprecated Values**: Reject old/deprecated data
- **Exclusion Rules**: Must not be X
- **Security**: Block sensitive patterns
- **Port Ranges**: Exclude reserved/restricted ports
- **Type Exclusion**: Accept anything except specific type

## Validation Logic

```
Inner Schema Valid → Not Schema INVALID (rejects match)
Inner Schema Invalid → Not Schema VALID (accepts non-match)
```

## Error Handling

```go
result := notSchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Code: "not_match" - value matched when it shouldn't
    }
}
```

## Internationalization

```go
notSchema := schema.Not(
    schema.String().Enum([]string{"admin", "root"}),
).NotError(i18n.S("username cannot be a reserved keyword"))
```

## JSON Schema Generation

```go
notSchema := schema.Not(schema.String())

jsonSchema := notSchema.JSON()
// Outputs:
// {
//   "not": {
//     "type": "string"
//   }
// }
```

## Combining with Other Schemas

### AllOf: Must be X but NOT Y

```go
// Must be string, but not email
stringButNotEmailSchema := schema.AllOf(
    schema.String(),
    schema.Not(schema.String().Email()),
)
```

### Object with Exclusions

```go
// User object but not admin
userNotAdminSchema := schema.AllOf(
    schema.Object().
        Property("username", schema.String()).
        Property("role", schema.String()),
    schema.Not(
        schema.Object().Property("role", schema.String().Const("admin")),
    ),
)
```

## Related

- [Union Schema](union.md) - For either/or validation
- [Conditional Schema](conditional.md) - For if/then/else logic
- [Any Schema](any.md) - For accepting anything
- [String Schema](string.md) - For string pattern validation