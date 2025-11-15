# Null Schema

The `NullSchema` validates that a value is exactly `nil` (null). It's useful for explicitly accepting only null values or in union types.

## Creating a Null Schema

```go
import "github.com/nyxstack/schema"

// Accept only null
nullSchema := schema.Null()

// With custom error message
nullSchema := schema.Null(i18n.S("value must be null"))
```

## Methods

### Core Methods

#### `Null(errorMessage ...interface{}) *NullSchema`
Creates a new schema that only accepts null values.

```go
nullValue := schema.Null()
```

### Required/Optional

#### `Required(errorMessage ...interface{}) *NullSchema`
Marks the null value as required (default).

```go
schema.Null().Required()
schema.Null().Required(i18n.S("value must be null"))
```

#### `Optional() *NullSchema`
Marks the null value as optional.

```go
schema.Null().Optional()
```

### Error Customization

#### `TypeError(message string) *NullSchema`
Sets custom error message for non-null values.

```go
schema.Null().TypeError("only null is accepted")
```

### Metadata

#### `Title(title string) *NullSchema`
Sets the schema title.

```go
schema.Null().Title("Null Value")
```

#### `Description(description string) *NullSchema`
Sets the schema description.

```go
schema.Null().Description("Explicitly null field")
```

## Usage Examples

### Basic Null Validation

```go
nullSchema := schema.Null()

ctx := schema.DefaultValidationContext()

// Valid
result := nullSchema.Parse(nil, ctx) // Valid

// Invalid
result = nullSchema.Parse(0, ctx)        // Invalid
result = nullSchema.Parse("", ctx)       // Invalid (empty string is not null)
result = nullSchema.Parse(false, ctx)    // Invalid
```

### Nullable Field (OneOf)

```go
// String OR null
nullableStringSchema := schema.OneOf(
    schema.String(),
    schema.Null(),
)

result := nullableStringSchema.Parse("hello", ctx) // Valid (string)
result = nullableStringSchema.Parse(nil, ctx)      // Valid (null)
result = nullableStringSchema.Parse(42, ctx)       // Invalid (not string or null)
```

### Optional Nullable Field

```go
userSchema := schema.Object().
    Property("id", schema.Int().Required()).
    Property("name", schema.String().Required()).
    Property("deletedAt", schema.OneOf(
        schema.DateTime(),
        schema.Null(),
    ).Optional())

// Valid: with null deletedAt
result := userSchema.Parse(map[string]interface{}{
    "id":        1,
    "name":      "Alice",
    "deletedAt": nil,
}, ctx)

// Valid: without deletedAt field
result = userSchema.Parse(map[string]interface{}{
    "id":   1,
    "name": "Alice",
}, ctx)

// Valid: with datetime deletedAt
result = userSchema.Parse(map[string]interface{}{
    "id":        1,
    "name":      "Alice",
    "deletedAt": "2025-11-17T14:30:00Z",
}, ctx)
```

### Explicitly Disabled Feature

```go
featureSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("enabled", schema.Bool()).
    Property("config", schema.OneOf(
        schema.Object(), // Feature config
        schema.Null(),   // Explicitly disabled (no config)
    ))

// Enabled with config
result := featureSchema.Parse(map[string]interface{}{
    "name":    "EmailNotifications",
    "enabled": true,
    "config":  map[string]interface{}{"smtp": "smtp.example.com"},
}, ctx)

// Disabled (null config)
result = featureSchema.Parse(map[string]interface{}{
    "name":    "EmailNotifications",
    "enabled": false,
    "config":  nil,
}, ctx)
```

### API Response with Null Data

```go
apiResponseSchema := schema.Object().
    Property("success", schema.Bool().Required()).
    Property("data", schema.OneOf(
        schema.Object(),
        schema.Array(schema.Any()),
        schema.Null(), // No data available
    ))

// Success with data
result := apiResponseSchema.Parse(map[string]interface{}{
    "success": true,
    "data":    map[string]interface{}{"id": 1},
}, ctx)

// Success but no data
result = apiResponseSchema.Parse(map[string]interface{}{
    "success": true,
    "data":    nil,
}, ctx)
```

### Database Nullable Column

```go
// Nullable integer (database NULL)
nullableIntSchema := schema.OneOf(
    schema.Int(),
    schema.Null(),
)

recordSchema := schema.Object().
    Property("id", schema.Int().Required()).
    Property("userId", nullableIntSchema).      // Foreign key (nullable)
    Property("parentId", nullableIntSchema)     // Self-reference (nullable)
```

### Three-State Boolean

```go
// Boolean with explicit null (Yes/No/Unknown)
threeStateSchema := schema.OneOf(
    schema.Bool(),
    schema.Null(),
)

surveySchema := schema.Object().
    Property("question", schema.String().Required()).
    Property("answer", threeStateSchema) // true, false, or null (not answered)

// Answered: Yes
result := surveySchema.Parse(map[string]interface{}{
    "question": "Do you agree?",
    "answer":   true,
}, ctx)

// Answered: No
result = surveySchema.Parse(map[string]interface{}{
    "question": "Do you agree?",
    "answer":   false,
}, ctx)

// Not answered
result = surveySchema.Parse(map[string]interface{}{
    "question": "Do you agree?",
    "answer":   nil,
}, ctx)
```

### JSON Patch Operations

```go
// JSON Patch "remove" operation uses null
patchSchema := schema.Object().
    Property("op", schema.String().Enum([]string{"add", "remove", "replace"})).
    Property("path", schema.String().Required()).
    Property("value", schema.OneOf(
        schema.Any(),
        schema.Null(), // For "remove" operations
    ))

result := patchSchema.Parse(map[string]interface{}{
    "op":    "remove",
    "path":  "/name",
    "value": nil,
}, ctx)
```

### Conditional Null

```go
// If status is "deleted", deletedAt must be datetime, otherwise null
deletedSchema := schema.Conditional(
    schema.Object().Property("status", schema.String().Const("deleted")),
).Then(
    schema.Object().Property("deletedAt", schema.DateTime().Required()),
).Else(
    schema.Object().Property("deletedAt", schema.Null()),
)
```

### Graph Node (Nullable Edges)

```go
nodeSchema := schema.Object().
    Property("id", schema.Int().Required()).
    Property("value", schema.Any()).
    Property("left", schema.OneOf(
        schema.Ref("#/Node", registry),
        schema.Null(), // Leaf node (no left child)
    )).
    Property("right", schema.OneOf(
        schema.Ref("#/Node", registry),
        schema.Null(), // Leaf node (no right child)
    ))
```

## When to Use

Null schemas are ideal for:

- **Nullable Fields**: Database NULL columns
- **Optional Values**: Explicitly absent data
- **Union Types**: String OR null, Int OR null
- **Three-State Logic**: true/false/null (yes/no/unknown)
- **Soft Deletes**: deletedAt timestamp or null
- **Feature Flags**: Enabled with config, or null (disabled)
- **API Responses**: Data available or null
- **Graph Structures**: Nullable references/pointers

## Important Notes

### Empty String vs Null

```go
// These are DIFFERENT
emptyString := ""  // Not null
nullValue := nil   // Null

nullSchema := schema.Null()
result := nullSchema.Parse("", ctx)  // INVALID - empty string is not null
result = nullSchema.Parse(nil, ctx)  // VALID
```

### Zero Values vs Null

```go
// These are DIFFERENT
zero := 0         // Not null
false := false    // Not null
nullValue := nil  // Null

nullSchema := schema.Null()
result := nullSchema.Parse(0, ctx)     // INVALID
result = nullSchema.Parse(false, ctx)  // INVALID
result = nullSchema.Parse(nil, ctx)    // VALID
```

## Error Handling

```go
result := nullSchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes: "required", "invalid_type"
    }
}
```

## Internationalization

```go
nullSchema := schema.Null().
    Required(i18n.S("value must be null")).
    TypeError("only null values are accepted")
```

## JSON Schema Generation

```go
nullSchema := schema.Null().
    Title("Null Value").
    Description("Explicitly null field")

jsonSchema := nullSchema.JSON()
// Outputs:
// {
//   "type": "null",
//   "title": "Null Value",
//   "description": "Explicitly null field"
// }
```

## Related

- [Union Schema](union.md) - For nullable types (OneOf with Null)
- [Any Schema](any.md) - For accepting anything including null
- [Conditional Schema](conditional.md) - For conditional null values
- [Not Schema](not.md) - For rejecting null values