# Object Schema

The `ObjectSchema` provides comprehensive validation for object/map values with support for property schemas, required fields, additional properties control, and more.

## Creating an Object Schema

```go
import "github.com/nyxstack/schema"

// Basic object schema
userSchema := schema.Object()

// Using Property method
userSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("email", schema.String().Email().Required()).
    Property("age", schema.Int().Min(0).Optional())

// Using Shape (concise syntax)
userSchema := schema.Shape{
    "name":  schema.String().Required(),
    "email": schema.String().Email().Required(),
    "age":   schema.Int().Min(0).Optional(),
}.AsObject()
```

## Methods

### Type Configuration

#### `Required(messages ...ErrorMessage) *ObjectSchema`
Marks the object as required (cannot be nil or omitted).

```go
schema.Object().Required()
schema.Object().Required("Object is required")
schema.Object().Required(i18n.S("object is required"))
```

#### `Optional() *ObjectSchema`
Marks the object as optional (can be nil or omitted).

```go
schema.Object().Optional()
```

#### `Nullable() *ObjectSchema`
Allows the object value to be explicitly null.

```go
schema.Object().Nullable()
```

#### `Default(value interface{}) *ObjectSchema`
Sets a default value when the input is nil.

```go
schema.Object().Default(map[string]interface{}{
    "status": "active",
})
```

### Property Definition

#### `Property(name string, schema interface{}) *ObjectSchema`
Adds a property with automatic required/optional detection from the schema.

```go
schema.Object().
    Property("name", schema.String().Required()).    // Required
    Property("age", schema.Int().Optional())         // Optional
```

#### `RequiredProperty(name string, schema interface{}) *ObjectSchema`
Explicitly adds a required property.

```go
schema.Object().
    RequiredProperty("id", schema.Int()).
    RequiredProperty("name", schema.String())
```

#### `OptionalProperty(name string, schema interface{}) *ObjectSchema`
Explicitly adds an optional property.

```go
schema.Object().
    OptionalProperty("nickname", schema.String()).
    OptionalProperty("avatar", schema.String().URL())
```

### Property Constraints

#### `MinProperties(min int, messages ...ErrorMessage) *ObjectSchema`
Sets the minimum number of properties.

```go
schema.Object().MinProperties(1)
schema.Object().MinProperties(2, "Must have at least 2 properties")
schema.Object().MinProperties(1, i18n.F("must have at least %d properties", 1))
```

#### `MaxProperties(max int, messages ...ErrorMessage) *ObjectSchema`
Sets the maximum number of properties.

```go
schema.Object().MaxProperties(10)
schema.Object().MaxProperties(5, "Cannot exceed 5 properties")
```

#### `PropertyRange(min, max int, messages ...ErrorMessage) *ObjectSchema`
Sets both minimum and maximum property counts.

```go
schema.Object().PropertyRange(2, 10)
```

### Additional Properties Control

#### `AdditionalProperties(allowed bool, messages ...ErrorMessage) *ObjectSchema`
Controls whether additional properties beyond defined ones are allowed.

```go
// Disallow additional properties
schema.Object().AdditionalProperties(false)

// Allow additional properties
schema.Object().AdditionalProperties(true)

// With custom error
schema.Object().AdditionalProperties(false, "Extra fields not allowed")
```

#### `Strict() *ObjectSchema`
Disallows additional properties (default behavior).

```go
schema.Object().
    Property("name", schema.String().Required()).
    Strict() // Only "name" is allowed
```

#### `Passthrough() *ObjectSchema`
Allows additional properties.

```go
schema.Object().
    Property("name", schema.String().Required()).
    Passthrough() // Additional properties allowed
```

### Metadata

#### `Title(title string) *ObjectSchema`
Sets a title for documentation and JSON Schema generation.

```go
schema.Object().Title("User Profile")
```

#### `Description(description string) *ObjectSchema`
Sets a description for documentation and JSON Schema generation.

```go
schema.Object().Description("User profile information")
```

## Usage Examples

### Basic Object Validation

```go
userSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("email", schema.String().Email().Required()).
    Property("age", schema.Int().Min(0).Optional())

ctx := schema.DefaultValidationContext()
result := userSchema.Parse(map[string]interface{}{
    "name":  "Alice",
    "email": "alice@example.com",
    "age":   30,
}, ctx)

if result.Valid {
    fmt.Printf("Valid user: %v\n", result.Value)
}
```

### Using Shape Syntax

```go
userSchema := schema.Shape{
    "username": schema.String().MinLength(3).Required(),
    "email":    schema.String().Email().Required(),
    "age":      schema.Int().Min(18).Optional(),
}.AsObject()
```

### Nested Objects

```go
addressSchema := schema.Shape{
    "street":  schema.String().Required(),
    "city":    schema.String().Required(),
    "zip":     schema.String().Pattern("^\\d{5}$").Required(),
}.AsObject()

userSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("address", addressSchema.Required())

result := userSchema.Parse(map[string]interface{}{
    "name": "Bob",
    "address": map[string]interface{}{
        "street": "123 Main St",
        "city":   "Springfield",
        "zip":    "12345",
    },
}, ctx)
```

### Strict vs Passthrough

```go
// Strict: Only defined properties allowed (default)
strictSchema := schema.Object().
    Property("name", schema.String().Required()).
    Strict()

// Invalid - "age" not defined
strictSchema.Parse(map[string]interface{}{
    "name": "Alice",
    "age":  30, // Error: additional property not allowed
}, ctx)

// Passthrough: Additional properties allowed
passthroughSchema := schema.Object().
    Property("name", schema.String().Required()).
    Passthrough()

// Valid - "age" is passed through
passthroughSchema.Parse(map[string]interface{}{
    "name": "Alice",
    "age":  30, // OK
}, ctx)
```

### Required vs Optional Properties

```go
schema := schema.Object().
    Property("id", schema.Int().Required()).           // Required
    Property("name", schema.String().Required()).      // Required
    Property("email", schema.String().Optional()).     // Optional
    Property("phone", schema.String().Optional())      // Optional

// Valid - optional fields can be omitted
schema.Parse(map[string]interface{}{
    "id":   1,
    "name": "Alice",
}, ctx)
```

### Property Count Constraints

```go
metadataSchema := schema.Object().
    MinProperties(1, "At least one property required").
    MaxProperties(10, "Maximum 10 properties allowed").
    Passthrough()

result := metadataSchema.Parse(map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
}, ctx)
```

### Complex Nested Structure

```go
orderSchema := schema.Object().
    Property("orderId", schema.Int().Required()).
    Property("customer", schema.Object().
        Property("id", schema.Int().Required()).
        Property("name", schema.String().Required()).
        Property("email", schema.String().Email().Required())).
    Property("items", schema.Array(
        schema.Object().
            Property("productId", schema.Int().Required()).
            Property("quantity", schema.Int().Min(1).Required()).
            Property("price", schema.Number().Min(0).Required()),
    ).MinItems(1)).
    Property("total", schema.Number().Min(0).Required())

result := orderSchema.Parse(map[string]interface{}{
    "orderId": 12345,
    "customer": map[string]interface{}{
        "id":    101,
        "name":  "Alice",
        "email": "alice@example.com",
    },
    "items": []map[string]interface{}{
        {"productId": 1, "quantity": 2, "price": 19.99},
        {"productId": 2, "quantity": 1, "price": 49.99},
    },
    "total": 89.97,
}, ctx)
```

### API Request Validation

```go
createUserRequest := schema.Object().
    Property("user", schema.Object().
        Property("username", schema.String().MinLength(3).Required()).
        Property("email", schema.String().Email().Required()).
        Property("password", schema.String().MinLength(8).Required())).
    AdditionalProperties(false)

result := createUserRequest.Parse(requestData, ctx)
```

### With i18n Support

```go
userSchema := schema.Object().
    Property("name", schema.String().
        MinLength(2, i18n.F("name must be at least %d characters", 2)).
        Required(i18n.S("name is required"))).
    Property("age", schema.Int().
        Min(18, i18n.F("must be at least %d years old", 18)).
        Required(i18n.S("age is required"))).
    AdditionalProperties(false, i18n.S("additional properties not allowed"))

result := userSchema.Parse(data, ctx)
```

### Configuration Object

```go
configSchema := schema.Object().
    Property("port", schema.Int().Min(1024).Max(65535).Default(8080)).
    Property("host", schema.String().Default("localhost")).
    Property("debug", schema.Bool().Default(false)).
    Property("timeout", schema.Int().Min(0).Default(30)).
    Passthrough()

result := configSchema.Parse(nil, ctx) // Uses all defaults
```

### Partial Updates

```go
updateSchema := schema.Object().
    OptionalProperty("name", schema.String().MinLength(1)).
    OptionalProperty("email", schema.String().Email()).
    OptionalProperty("age", schema.Int().Min(0)).
    MinProperties(1, "At least one field must be provided").
    AdditionalProperties(false)

result := updateSchema.Parse(map[string]interface{}{
    "email": "newemail@example.com",
}, ctx) // Valid - only updating email
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

Property errors include the property name in the path:

```
Path: name
Message: value must be at least 2 characters long

Path: address.zip
Message: value format is invalid
```

## Internationalization

All error messages support i18n through the `github.com/nyxstack/i18n` package:

```go
schema.Object().
    Property("name", schema.String().Required(i18n.S("name is required"))).
    MinProperties(1, i18n.F("must have at least %d properties", 1)).
    AdditionalProperties(false, i18n.S("additional properties not allowed"))
```

## JSON Schema Generation

```go
schema := schema.Object().
    Title("User").
    Description("User profile").
    Property("id", schema.Int().Required()).
    Property("name", schema.String().Required()).
    Property("email", schema.String().Email().Optional()).
    AdditionalProperties(false)

jsonSchema := schema.JSON()
// Outputs:
// {
//   "type": "object",
//   "properties": {
//     "id": {"type": "integer"},
//     "name": {"type": "string"},
//     "email": {"type": "string", "format": "email"}
//   },
//   "required": ["id", "name"],
//   "additionalProperties": false,
//   "title": "User",
//   "description": "User profile"
// }
```

## Related

- [Record Schema](record.md) - For dynamic key-value maps
- [Array Schema](array.md) - For arrays of objects
- [String Schema](string.md) - For object properties
- [Shape](object.md#using-shape-syntax) - Concise object syntax
