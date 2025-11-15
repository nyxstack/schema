# Array Schema

The `ArraySchema` provides comprehensive validation for array/slice values with support for item validation, length constraints, uniqueness checking, and more.

## Creating an Array Schema

```go
import "github.com/nyxstack/schema"

// Array of strings
tagsSchema := schema.Array(schema.String())

// Array of integers
scoresSchema := schema.Array(schema.Int())

// With validation constraints
tagsSchema := schema.Array(schema.String().MinLength(1)).
    MinItems(1).
    MaxItems(10).
    UniqueItems().
    Required("Tags are required")
```

## Methods

### Type Configuration

#### `Required(messages ...ErrorMessage) *ArraySchema`
Marks the array as required (cannot be nil or omitted).

```go
schema.Array(schema.String()).Required()
schema.Array(schema.String()).Required("List is required")
schema.Array(schema.String()).Required(i18n.S("list is required"))
```

#### `Optional() *ArraySchema`
Marks the array as optional (can be nil or omitted).

```go
schema.Array(schema.String()).Optional()
```

#### `Nullable() *ArraySchema`
Allows the array value to be explicitly null.

```go
schema.Array(schema.String()).Nullable()
```

#### `Default(value interface{}) *ArraySchema`
Sets a default value when the input is nil.

```go
schema.Array(schema.String()).Default([]string{})
schema.Array(schema.Int()).Default([]int{1, 2, 3})
```

### Length Constraints

#### `MinItems(min int, messages ...ErrorMessage) *ArraySchema`
Sets the minimum number of items in the array.

```go
schema.Array(schema.String()).MinItems(1)
schema.Array(schema.String()).MinItems(3, "Must have at least 3 items")
schema.Array(schema.String()).MinItems(1, i18n.F("must have at least %d items", 1))
```

#### `MaxItems(max int, messages ...ErrorMessage) *ArraySchema`
Sets the maximum number of items in the array.

```go
schema.Array(schema.String()).MaxItems(10)
schema.Array(schema.String()).MaxItems(5, "Cannot exceed 5 items")
```

#### `Length(exact int) *ArraySchema`
Requires an exact number of items.

```go
// Exactly 3 items required
schema.Array(schema.String()).Length(3)
```

### Uniqueness

#### `UniqueItems(messages ...ErrorMessage) *ArraySchema`
Requires all items in the array to be unique (no duplicates).

```go
schema.Array(schema.String()).UniqueItems()
schema.Array(schema.Int()).UniqueItems("Items must be unique")
schema.Array(schema.String()).UniqueItems(i18n.S("items must be unique"))
```

### Item Schema

#### `Items(itemSchema Parseable) *ArraySchema`
Sets or changes the schema used to validate each array item.

```go
// Define array first, then set item schema
arraySchema := schema.Array(schema.String())
arraySchema.Items(schema.String().MinLength(3))
```

### Metadata

#### `Title(title string) *ArraySchema`
Sets a title for documentation and JSON Schema generation.

```go
schema.Array(schema.String()).Title("User Tags")
```

#### `Description(description string) *ArraySchema`
Sets a description for documentation and JSON Schema generation.

```go
schema.Array(schema.String()).Description("List of tags associated with the user")
```

## Usage Examples

### Basic Array Validation

```go
tagsSchema := schema.Array(schema.String()).
    MinItems(1).
    MaxItems(10).
    Required("Tags are required")

ctx := schema.DefaultValidationContext()
result := tagsSchema.Parse([]string{"go", "schema", "validation"}, ctx)

if result.Valid {
    fmt.Printf("Valid tags: %v\n", result.Value)
}
```

### Array of Integers with Range

```go
scoresSchema := schema.Array(schema.Int().Min(0).Max(100)).
    MinItems(1).
    MaxItems(10)

result := scoresSchema.Parse([]int{85, 92, 78, 95}, ctx)
```

### Unique Items

```go
uniqueTagsSchema := schema.Array(schema.String().MinLength(1)).
    UniqueItems("Tags must be unique").
    MinItems(1)

result := uniqueTagsSchema.Parse([]string{"go", "rust", "python"}, ctx) // Valid
result := uniqueTagsSchema.Parse([]string{"go", "go", "python"}, ctx) // Invalid - duplicate "go"
```

### Array of Objects

```go
usersSchema := schema.Array(
    schema.Object().
        Property("id", schema.Int().Required()).
        Property("name", schema.String().Required()),
).MinItems(1)

result := usersSchema.Parse([]map[string]interface{}{
    {"id": 1, "name": "Alice"},
    {"id": 2, "name": "Bob"},
}, ctx)
```

### Nested Arrays

```go
matrixSchema := schema.Array(
    schema.Array(schema.Int()).MinItems(3).MaxItems(3),
).MinItems(3).MaxItems(3)

result := matrixSchema.Parse([][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}, ctx)
```

### Optional Array with Default

```go
rolesSchema := schema.Array(schema.String()).
    Default([]string{"user"}).
    Optional()

result := rolesSchema.Parse(nil, ctx) // Returns []string{"user"}
```

### Email List

```go
emailsSchema := schema.Array(schema.String().Email()).
    MinItems(1, "At least one email is required").
    MaxItems(5, "Maximum 5 emails allowed").
    UniqueItems("Duplicate emails not allowed")

result := emailsSchema.Parse([]string{
    "alice@example.com",
    "bob@example.com",
}, ctx)
```

### Fixed-Length Array

```go
rgbSchema := schema.Array(schema.Int().Min(0).Max(255)).
    Length(3).
    Title("RGB Color")

result := rgbSchema.Parse([]int{255, 0, 128}, ctx) // Valid
result := rgbSchema.Parse([]int{255, 0}, ctx) // Invalid - needs exactly 3 items
```

### Product IDs

```go
productIdsSchema := schema.Array(schema.Int().Min(1)).
    MinItems(1, i18n.F("must select at least %d product", 1)).
    MaxItems(50, i18n.F("cannot select more than %d products", 50)).
    UniqueItems(i18n.S("duplicate product ids not allowed"))

result := productIdsSchema.Parse([]int{101, 102, 103}, ctx)
```

### Complex Item Validation

```go
todoSchema := schema.Array(
    schema.Object().
        Property("id", schema.Int().Required()).
        Property("title", schema.String().MinLength(1).Required()).
        Property("completed", schema.Bool().Default(false)),
).MinItems(1).MaxItems(100)

result := todoSchema.Parse([]map[string]interface{}{
    {"id": 1, "title": "Buy groceries", "completed": false},
    {"id": 2, "title": "Write docs", "completed": true},
}, ctx)
```

### Array with Enum Items

```go
statusesSchema := schema.Array(
    schema.String().Enum([]string{"pending", "active", "completed"}),
).UniqueItems()

result := statusesSchema.Parse([]string{"pending", "active"}, ctx)
```

## Type Coercion

The array schema accepts any slice type and validates each item:

```go
schema := schema.Array(schema.Int())

// These all work
schema.Parse([]int{1, 2, 3}, ctx)
schema.Parse([]interface{}{1, 2, 3}, ctx)
```

## Empty Arrays

Empty arrays are valid by default unless `MinItems` is set:

```go
// Empty array is valid
schema.Array(schema.String()).Parse([]string{}, ctx)

// Empty array is invalid
schema.Array(schema.String()).MinItems(1).Parse([]string{}, ctx)
```

## Internationalization

All error messages support i18n through the `github.com/nyxstack/i18n` package:

```go
schema.Array(schema.String()).
    MinItems(1, i18n.F("must have at least %d items", 1)).
    MaxItems(10, i18n.F("cannot exceed %d items", 10)).
    UniqueItems(i18n.S("items must be unique")).
    Required(i18n.S("list is required"))
```

## JSON Schema Generation

```go
schema := schema.Array(schema.String()).
    Title("Tags").
    Description("User tags").
    MinItems(1).
    MaxItems(10).
    UniqueItems().
    Required()

jsonSchema := schema.JSON()
// Outputs:
// {
//   "type": "array",
//   "items": {
//     "type": "string"
//   },
//   "minItems": 1,
//   "maxItems": 10,
//   "uniqueItems": true,
//   "title": "Tags",
//   "description": "User tags"
// }
```

## Error Handling

```go
result := schema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Path: %v\n", err.Path)
        fmt.Printf("Value: %v\n", err.Value)
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
    }
}
```

Errors on individual items will include the array index in the path:

```
Path: [1]
Message: value must be at least 3 characters long
```

## Related

- [Tuple Schema](tuple.md) - For fixed-position arrays with different types
- [Object Schema](object.md) - For objects/maps
- [String Schema](string.md) - For array items
- [Int Schema](int.md) - For integer arrays
