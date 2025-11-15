# Number Schema (Float/Number)

The `FloatSchema` and `NumberSchema` provide comprehensive validation for floating-point and decimal values with support for range constraints, precision control, and more.

## Creating a Number Schema

```go
import "github.com/nyxstack/schema"

// Float32 schema
priceSchema := schema.Float()

// Float64 schema (Number)
preciseSchema := schema.Number()

// With validation constraints
priceSchema := schema.Number().
    Min(0.01).
    Max(9999.99).
    MultipleOf(0.01).
    Required("Price is required")
```

## Methods

Both `Float` (float32) and `Number` (float64) schemas share the same methods.

### Type Configuration

#### `Required(messages ...ErrorMessage) *NumberSchema`
Marks the number as required (cannot be nil or omitted).

```go
schema.Number().Required()
schema.Number().Required("Price is required")
schema.Number().Required(i18n.S("price is required"))
```

#### `Optional() *NumberSchema`
Marks the number as optional (can be nil or omitted).

```go
schema.Number().Optional()
```

#### `Nullable() *NumberSchema`
Allows the number value to be explicitly null.

```go
schema.Number().Nullable()
```

#### `Default(value interface{}) *NumberSchema`
Sets a default value when the input is nil.

```go
schema.Number().Default(0.0)
schema.Float().Default(3.14)
```

### Range Constraints

#### `Min(min float64, messages ...ErrorMessage) *NumberSchema`
Sets the minimum value (inclusive).

```go
schema.Number().Min(0.0)
schema.Number().Min(0.01, "Price must be at least 0.01")
schema.Number().Min(0.0, i18n.S("value must be positive"))
```

#### `Max(max float64, messages ...ErrorMessage) *NumberSchema`
Sets the maximum value (inclusive).

```go
schema.Number().Max(100.0)
schema.Number().Max(999.99, "Price cannot exceed 999.99")
```

#### `Range(min, max float64, messages ...ErrorMessage) *NumberSchema`
Sets both minimum and maximum values in one call.

```go
schema.Number().Range(0.0, 100.0)
schema.Number().Range(0.01, 9999.99, "Price must be between 0.01 and 9999.99")
```

### Precision Control

#### `MultipleOf(multiple float64, messages ...ErrorMessage) *NumberSchema`
Requires the value to be a multiple of the specified number.

```go
// Must be in increments of 0.01 (cents)
schema.Number().MultipleOf(0.01)

// Must be in increments of 0.25
schema.Number().MultipleOf(0.25, "Value must be in increments of 0.25")

// Must be whole dollars
schema.Number().MultipleOf(1.0)
```

### Value Constraints

#### `Enum(values []float64, messages ...ErrorMessage) *NumberSchema`
Restricts the number to a set of allowed values.

```go
schema.Number().Enum([]float64{0.0, 0.5, 1.0, 1.5, 2.0})
schema.Number().Enum([]float64{9.99, 19.99, 29.99}, "Invalid price point")
```

#### `Const(value float64, messages ...ErrorMessage) *NumberSchema`
Requires the number to match an exact value.

```go
schema.Number().Const(3.14159)
schema.Number().Const(9.99, "Price must be 9.99")
```

### Metadata

#### `Title(title string) *NumberSchema`
Sets a title for documentation and JSON Schema generation.

```go
schema.Number().Title("Product Price")
```

#### `Description(description string) *NumberSchema`
Sets a description for documentation and JSON Schema generation.

```go
schema.Number().Description("The price in USD")
```

## Float vs Number

| Type | Size | Precision | Use Case |
|------|------|-----------|----------|
| `Float` | 32-bit | ~7 decimal digits | Performance-critical, less precision needed |
| `Number` | 64-bit | ~15 decimal digits | Financial calculations, high precision |

**Recommendation:** Use `Number` (float64) for most cases, especially financial data.

## Usage Examples

### Basic Number Validation

```go
priceSchema := schema.Number().
    Min(0.01).
    Max(9999.99).
    Required("Price is required")

ctx := schema.DefaultValidationContext()
result := priceSchema.Parse(19.99, ctx)

if result.Valid {
    fmt.Printf("Valid price: %.2f\n", result.Value)
}
```

### Price with Precision

```go
priceSchema := schema.Number().
    Min(0.01, "Price must be at least $0.01").
    Max(9999.99, "Price cannot exceed $9999.99").
    MultipleOf(0.01, "Price must be in cents").
    Required()

result := priceSchema.Parse(19.99, ctx) // Valid
result := priceSchema.Parse(19.999, ctx) // Invalid - not a multiple of 0.01
```

### Percentage Validation

```go
percentageSchema := schema.Number().
    Range(0.0, 100.0, "Percentage must be between 0 and 100").
    MultipleOf(0.1)

result := percentageSchema.Parse(75.5, ctx)
```

### Temperature with Negatives

```go
tempSchema := schema.Number().
    Min(-273.15, "Cannot be below absolute zero").
    Max(1000.0)

result := tempSchema.Parse(-40.0, ctx) // Valid
result := tempSchema.Parse(-300.0, ctx) // Invalid
```

### Discount Rate

```go
discountSchema := schema.Number().
    Enum([]float64{0.0, 0.05, 0.10, 0.15, 0.20, 0.25}).
    Default(0.0).
    Title("Discount Rate")

result := discountSchema.Parse(nil, ctx) // Uses default: 0.0
```

### Float for Performance

```go
coordinateSchema := schema.Float().
    Range(-180.0, 180.0).
    Required()

result := coordinateSchema.Parse(45.5, ctx)
```

### Complex Validation

```go
priceSchema := schema.Number().
    Min(0.01, i18n.F("price must be at least $%g", 0.01)).
    Max(99999.99, i18n.F("price cannot exceed $%g", 99999.99)).
    MultipleOf(0.01, i18n.S("price must be in cents")).
    Required(i18n.S("price is required"))

result := priceSchema.Parse(1299.99, ctx)
```

### Optional with Default

```go
taxRateSchema := schema.Number().
    Min(0.0).
    Max(1.0).
    Default(0.08).
    Optional().
    Title("Tax Rate")

result := taxRateSchema.Parse(nil, ctx) // Returns 0.08
```

### Type Coercion

The number schemas support automatic type coercion from compatible types:

```go
schema := schema.Number()

// These all work
schema.Parse(3.14, ctx)         // float64
schema.Parse(float32(3.14), ctx) // float32
schema.Parse(42, ctx)           // int (converted to float)
schema.Parse(int64(42), ctx)    // int64 (converted to float)
```

## Floating-Point Precision

⚠️ **Important:** Due to floating-point representation, exact comparisons may fail:

```go
// This might fail due to precision
schema.Number().Const(0.1 + 0.2) // 0.30000000000000004

// Use MultipleOf for precision control
schema.Number().MultipleOf(0.01) // Better for currency
```

## Internationalization

All error messages support i18n through the `github.com/nyxstack/i18n` package:

```go
schema.Number().
    Min(0.01, i18n.F("price must be at least $%g", 0.01)).
    Required(i18n.S("price is required"))
```

## JSON Schema Generation

```go
schema := schema.Number().
    Title("Product Price").
    Description("Price in USD").
    Min(0.01).
    Max(9999.99).
    MultipleOf(0.01).
    Required()

jsonSchema := schema.JSON()
// Outputs:
// {
//   "type": "number",
//   "minimum": 0.01,
//   "maximum": 9999.99,
//   "multipleOf": 0.01,
//   "title": "Product Price",
//   "description": "Price in USD"
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

- [Int Schema](int.md) - For whole numbers
- [String Schema](string.md) - For string values
- [Object Schema](object.md) - For objects with number properties
- [Array Schema](array.md) - For arrays of numbers
