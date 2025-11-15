# Integer Schema

The `IntSchema` and related integer schemas (`Int8Schema`, `Int16Schema`, `Int32Schema`, `Int64Schema`) provide comprehensive validation for integer values with support for range constraints, multiple-of validation, and more.

## Creating an Integer Schema

```go
import "github.com/nyxstack/schema"

// Basic int schema
ageSchema := schema.Int()

// Int8 schema (-128 to 127)
byteSchema := schema.Int8()

// Int16 schema (-32768 to 32767)
shortSchema := schema.Int16()

// Int32 schema (-2147483648 to 2147483647)
intSchema := schema.Int32()

// Int64 schema (full 64-bit range)
longSchema := schema.Int64()

// With validation constraints
ageSchema := schema.Int().
    Min(0).
    Max(120).
    Required("Age is required")
```

## Methods

All integer schema types (`Int`, `Int8`, `Int16`, `Int32`, `Int64`) share the same methods.

### Type Configuration

#### `Required(messages ...ErrorMessage) *IntSchema`
Marks the integer as required (cannot be nil or omitted).

```go
schema.Int().Required()
schema.Int().Required("Age is required")
schema.Int().Required(i18n.S("age is required"))
```

#### `Optional() *IntSchema`
Marks the integer as optional (can be nil or omitted).

```go
schema.Int().Optional()
```

#### `Nullable() *IntSchema`
Allows the integer value to be explicitly null.

```go
schema.Int().Nullable()
```

#### `Default(value interface{}) *IntSchema`
Sets a default value when the input is nil.

```go
schema.Int().Default(0)
schema.Int().Default(18)
```

### Range Constraints

#### `Min(min int, messages ...ErrorMessage) *IntSchema`
Sets the minimum value (inclusive).

```go
schema.Int().Min(0)
schema.Int().Min(18, "Must be at least 18 years old")
schema.Int().Min(1, i18n.S("value must be positive"))
```

#### `Max(max int, messages ...ErrorMessage) *IntSchema`
Sets the maximum value (inclusive).

```go
schema.Int().Max(100)
schema.Int().Max(120, "Age cannot exceed 120")
```

#### `Range(min, max int, messages ...ErrorMessage) *IntSchema`
Sets both minimum and maximum values in one call.

```go
schema.Int().Range(0, 100)
schema.Int().Range(18, 65, "Age must be between 18 and 65")
```

### Multiple Validation

#### `MultipleOf(multiple int, messages ...ErrorMessage) *IntSchema`
Requires the value to be a multiple of the specified number.

```go
// Must be even
schema.Int().MultipleOf(2)

// Must be multiple of 5
schema.Int().MultipleOf(5, "Quantity must be in multiples of 5")

// Must be multiple of 10
schema.Int().MultipleOf(10)
```

### Value Constraints

#### `Enum(values []int, messages ...ErrorMessage) *IntSchema`
Restricts the integer to a set of allowed values.

```go
schema.Int().Enum([]int{1, 2, 3, 5, 8, 13})
schema.Int().Enum([]int{18, 21, 25, 30}, "Invalid age bracket")
```

#### `Const(value int, messages ...ErrorMessage) *IntSchema`
Requires the integer to match an exact value.

```go
schema.Int().Const(1)
schema.Int().Const(42, "Value must be 42")
```

### Metadata

#### `Title(title string) *IntSchema`
Sets a title for documentation and JSON Schema generation.

```go
schema.Int().Title("User Age")
```

#### `Description(description string) *IntSchema`
Sets a description for documentation and JSON Schema generation.

```go
schema.Int().Description("The user's age in years")
```

## Integer Type Ranges

Each integer schema type has specific range limits:

| Type | Size | Minimum | Maximum |
|------|------|---------|---------|
| `Int8` | 8-bit | -128 | 127 |
| `Int16` | 16-bit | -32,768 | 32,767 |
| `Int32` | 32-bit | -2,147,483,648 | 2,147,483,647 |
| `Int64` | 64-bit | -9,223,372,036,854,775,808 | 9,223,372,036,854,775,807 |
| `Int` | Platform | Platform dependent | Platform dependent |

Values outside these ranges will fail type validation.

## Usage Examples

### Basic Integer Validation

```go
ageSchema := schema.Int().
    Min(0).
    Max(120).
    Required("Age is required")

ctx := schema.DefaultValidationContext()
result := ageSchema.Parse(25, ctx)

if result.Valid {
    fmt.Printf("Valid age: %d\n", result.Value)
}
```

### Range Validation

```go
percentageSchema := schema.Int().
    Range(0, 100, "Percentage must be between 0 and 100")

result := percentageSchema.Parse(75, ctx)
```

### Multiple Of Validation

```go
evenNumberSchema := schema.Int().
    MultipleOf(2, "Must be an even number").
    Min(0)

result := evenNumberSchema.Parse(42, ctx) // Valid
result := evenNumberSchema.Parse(43, ctx) // Invalid
```

### Enum Validation

```go
prioritySchema := schema.Int().
    Enum([]int{1, 2, 3, 4, 5}).
    Default(3).
    Title("Priority Level")

result := prioritySchema.Parse(nil, ctx) // Uses default: 3
```

### Int8 for Small Ranges

```go
dayOfWeekSchema := schema.Int8().
    Range(1, 7, "Day must be 1-7").
    Required()

result := dayOfWeekSchema.Parse(5, ctx)
```

### Int64 for Large Numbers

```go
timestampSchema := schema.Int64().
    Min(0).
    Title("Unix Timestamp")

result := timestampSchema.Parse(1700000000, ctx)
```

### Complex Validation

```go
quantitySchema := schema.Int().
    Min(1, i18n.F("quantity must be at least %d", 1)).
    Max(1000, i18n.F("quantity cannot exceed %d", 1000)).
    MultipleOf(5, i18n.F("quantity must be in multiples of %d", 5)).
    Required(i18n.S("quantity is required"))

result := quantitySchema.Parse(25, ctx)
```

### Optional with Default

```go
timeoutSchema := schema.Int().
    Min(0).
    Max(300).
    Default(30).
    Optional().
    Title("Timeout in seconds")

result := timeoutSchema.Parse(nil, ctx) // Returns 30
```

### Type Coercion

The integer schemas support automatic type coercion from compatible types:

```go
schema := schema.Int()

// These all work
schema.Parse(42, ctx)           // int
schema.Parse(int8(42), ctx)     // int8
schema.Parse(int16(42), ctx)    // int16
schema.Parse(int32(42), ctx)    // int32
schema.Parse(int64(42), ctx)    // int64
schema.Parse(float64(42.0), ctx) // float64 (whole number only)
```

## Internationalization

All error messages support i18n through the `github.com/nyxstack/i18n` package:

```go
schema.Int().
    Min(18, i18n.F("must be at least %d years old", 18)).
    Required(i18n.S("age is required"))
```

## JSON Schema Generation

```go
schema := schema.Int().
    Title("User Age").
    Description("Age in years").
    Min(0).
    Max(120).
    Required()

jsonSchema := schema.JSON()
// Outputs:
// {
//   "type": "integer",
//   "minimum": 0,
//   "maximum": 120,
//   "title": "User Age",
//   "description": "Age in years"
// }
```

## Error Handling

```go
result := schema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Value: %v\n", err.Value)
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
    }
}
```

## Related

- [Float/Number Schema](number.md) - For decimal numbers
- [String Schema](string.md) - For string values
- [Object Schema](object.md) - For objects with integer properties
- [Array Schema](array.md) - For arrays of integers
