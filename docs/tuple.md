# Tuple Schema

The `TupleSchema` provides validation for fixed-length arrays where each position has a specific type. Unlike regular arrays, tuples validate each position independently.

## Creating a Tuple Schema

```go
import "github.com/nyxstack/schema"

// RGB color tuple [red, green, blue]
rgbSchema := schema.Tuple(
    schema.Int().Min(0).Max(255),
    schema.Int().Min(0).Max(255),
    schema.Int().Min(0).Max(255),
)

// Coordinate pair [latitude, longitude]
coordinateSchema := schema.Tuple(
    schema.Number().Min(-90).Max(90),
    schema.Number().Min(-180).Max(180),
)

// Mixed types [name, age, active]
userTupleSchema := schema.Tuple(
    schema.String().MinLength(1),
    schema.Int().Min(0),
    schema.Bool(),
)
```

## Methods

### Type Configuration

#### `Required(messages ...ErrorMessage) *TupleSchema`
Marks the tuple as required.

```go
schema.Tuple(schema.String(), schema.Int()).Required()
schema.Tuple(schema.String(), schema.Int()).Required(i18n.S("tuple is required"))
```

#### `Optional() *TupleSchema`
Marks the tuple as optional.

```go
schema.Tuple(schema.String(), schema.Int()).Optional()
```

#### `Nullable() *TupleSchema`
Allows null values.

```go
schema.Tuple(schema.String(), schema.Int()).Nullable()
```

### Length Constraints

#### `Strict() *TupleSchema`
Requires exact length matching (default behavior).

```go
schema.Tuple(schema.String(), schema.Int()).Strict()
// Must have exactly 2 items
```

#### `AllowAdditionalItems() *TupleSchema`
Allows extra items beyond defined positions.

```go
schema.Tuple(schema.String(), schema.Int()).AllowAdditionalItems()
// Can have 2 or more items
```

### Uniqueness

#### `UniqueItems(messages ...ErrorMessage) *TupleSchema`
Requires all items to be unique.

```go
schema.Tuple(schema.Int(), schema.Int(), schema.Int()).UniqueItems()
schema.Tuple(schema.String(), schema.String()).UniqueItems(i18n.S("items must be unique"))
```

### Metadata

#### `Title(title string) *TupleSchema`
Sets a title.

```go
schema.Tuple(schema.Number(), schema.Number()).Title("Coordinates")
```

#### `Description(description string) *TupleSchema`
Sets a description.

```go
schema.Tuple(schema.Number(), schema.Number()).Description("Latitude and longitude pair")
```

## Usage Examples

### RGB Color

```go
rgbSchema := schema.Tuple(
    schema.Int().Min(0).Max(255),
    schema.Int().Min(0).Max(255),
    schema.Int().Min(0).Max(255),
)

ctx := schema.DefaultValidationContext()
result := rgbSchema.Parse([]interface{}{255, 128, 0}, ctx)
```

### Geographic Coordinates

```go
coordinateSchema := schema.Tuple(
    schema.Number().Min(-90).Max(90),   // Latitude
    schema.Number().Min(-180).Max(180), // Longitude
).Title("GPS Coordinates")

result := coordinateSchema.Parse([]interface{}{37.7749, -122.4194}, ctx)
```

### Version Triple

```go
versionSchema := schema.Tuple(
    schema.Int().Min(0), // Major
    schema.Int().Min(0), // Minor
    schema.Int().Min(0), // Patch
)

result := versionSchema.Parse([]interface{}{1, 2, 3}, ctx) // v1.2.3
```

### Mixed Types

```go
userTupleSchema := schema.Tuple(
    schema.String().MinLength(1), // Name
    schema.Int().Min(0),           // Age
    schema.Bool(),                 // Active
)

result := userTupleSchema.Parse([]interface{}{"Alice", 30, true}, ctx)
```

### With Additional Items Allowed

```go
commandSchema := schema.Tuple(
    schema.String().MinLength(1), // Command name
    schema.String(),               // First arg
).AllowAdditionalItems()

result := commandSchema.Parse([]interface{}{"git", "commit", "-m", "message"}, ctx)
```

### Unique Items Constraint

```go
uniqueTripleSchema := schema.Tuple(
    schema.Int(),
    schema.Int(),
    schema.Int(),
).UniqueItems(i18n.S("all values must be different"))

result := uniqueTripleSchema.Parse([]interface{}{1, 2, 3}, ctx) // Valid
result := uniqueTripleSchema.Parse([]interface{}{1, 1, 3}, ctx) // Invalid
```

## Related

- [Array Schema](array.md) - For variable-length arrays with uniform types
- [Object Schema](object.md) - For named properties
