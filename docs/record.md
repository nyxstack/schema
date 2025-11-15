# Record Schema

The `RecordSchema` provides validation for key-value maps where both keys and values follow specific schemas. This is ideal for dynamic objects with unknown property names.

## Creating a Record Schema

```go
import "github.com/nyxstack/schema"

// Basic record schema
metadataSchema := schema.Record(
    schema.String(), // Key schema
    schema.String(), // Value schema
)

// With validation constraints
labelsSchema := schema.Record(
    schema.String().Pattern("^[a-z0-9-]+$"), // Keys: lowercase alphanumeric
    schema.String().MaxLength(100),           // Values: max 100 chars
).MinProperties(1).MaxProperties(10)
```

## Methods

### Type Configuration

#### `Required(messages ...ErrorMessage) *RecordSchema`
Marks the record as required (cannot be nil or omitted).

```go
schema.Record(schema.String(), schema.String()).Required()
schema.Record(schema.String(), schema.String()).Required(i18n.S("record is required"))
```

#### `Optional() *RecordSchema`
Marks the record as optional (can be nil or omitted).

```go
schema.Record(schema.String(), schema.String()).Optional()
```

#### `Nullable() *RecordSchema`
Allows the record value to be explicitly null.

```go
schema.Record(schema.String(), schema.String()).Nullable()
```

#### `Default(value interface{}) *RecordSchema`
Sets a default value when the input is nil.

```go
schema.Record(schema.String(), schema.String()).Default(map[string]interface{}{})
```

### Schema Configuration

#### `Keys(keySchema Parseable) *RecordSchema`
Sets or changes the schema for keys.

```go
schema.Record(schema.String(), schema.String()).
    Keys(schema.String().Pattern("^[a-z]+$"))
```

#### `Values(valueSchema Parseable) *RecordSchema`
Sets or changes the schema for values.

```go
schema.Record(schema.String(), schema.String()).
    Values(schema.String().MinLength(1).MaxLength(100))
```

### Property Constraints

#### `MinProperties(min int, messages ...ErrorMessage) *RecordSchema`
Sets the minimum number of key-value pairs.

```go
schema.Record(schema.String(), schema.String()).MinProperties(1)
schema.Record(schema.String(), schema.String()).MinProperties(1, i18n.F("must have at least %d properties", 1))
```

#### `MaxProperties(max int, messages ...ErrorMessage) *RecordSchema`
Sets the maximum number of key-value pairs.

```go
schema.Record(schema.String(), schema.String()).MaxProperties(10)
```

### Metadata

#### `Title(title string) *RecordSchema`
Sets a title for documentation.

```go
schema.Record(schema.String(), schema.String()).Title("Metadata")
```

#### `Description(description string) *RecordSchema`
Sets a description for documentation.

```go
schema.Record(schema.String(), schema.String()).Description("Key-value metadata pairs")
```

## Usage Examples

### Basic Record Validation

```go
metadataSchema := schema.Record(
    schema.String().MinLength(1),
    schema.String().MinLength(1),
)

ctx := schema.DefaultValidationContext()
result := metadataSchema.Parse(map[string]interface{}{
    "author": "Alice",
    "version": "1.0.0",
}, ctx)
```

### Labels with Pattern

```go
labelsSchema := schema.Record(
    schema.String().Pattern("^[a-z0-9-]+$"),
    schema.String().MaxLength(100),
).MinProperties(1).MaxProperties(10)

result := labelsSchema.Parse(map[string]interface{}{
    "env": "production",
    "team": "backend",
    "region": "us-west-2",
}, ctx)
```

### Configuration Map

```go
configSchema := schema.Record(
    schema.String().Pattern("^[A-Z_]+$"), // UPPERCASE_KEYS
    schema.String(),                       // Any string value
).MinProperties(1)

result := configSchema.Parse(map[string]interface{}{
    "DATABASE_URL": "postgres://localhost",
    "API_KEY": "secret",
}, ctx)
```

### Integer Values

```go
scoresSchema := schema.Record(
    schema.String(), // Player names
    schema.Int().Min(0).Max(100), // Scores
)

result := scoresSchema.Parse(map[string]interface{}{
    "alice": 95,
    "bob": 87,
    "charlie": 92,
}, ctx)
```

### Nested Objects as Values

```go
usersSchema := schema.Record(
    schema.String().Pattern("^[a-z0-9]+$"), // User IDs
    schema.Object().
        Property("name", schema.String().Required()).
        Property("email", schema.String().Email().Required()),
)

result := usersSchema.Parse(map[string]interface{}{
    "user1": map[string]interface{}{
        "name": "Alice",
        "email": "alice@example.com",
    },
    "user2": map[string]interface{}{
        "name": "Bob",
        "email": "bob@example.com",
    },
}, ctx)
```

### Feature Flags

```go
featuresSchema := schema.Record(
    schema.String().Pattern("^feature_[a-z_]+$"),
    schema.Bool(),
).MinProperties(1)

result := featuresSchema.Parse(map[string]interface{}{
    "feature_dark_mode": true,
    "feature_new_ui": false,
}, ctx)
```

## Related

- [Object Schema](object.md) - For fixed property names
- [String Schema](string.md) - For record keys/values
- [Array Schema](array.md) - For arrays of records
