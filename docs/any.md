# Any Schema

The `AnySchema` accepts any value type without validation. It's useful for dynamic data, configuration objects, or when you need to accept anything.

## Creating an Any Schema

```go
import "github.com/nyxstack/schema"

// Accepts any value
anySchema := schema.Any()

// With custom error message
anySchema := schema.Any(i18n.S("value is required"))
```

## Methods

### Core Methods

#### `Any(errorMessage ...interface{}) *AnySchema`
Creates a new schema that accepts any value.

```go
anyValue := schema.Any()
```

### Constraints

#### `Enum(values []interface{}) *AnySchema`
Restricts to specific allowed values (any types).

```go
schema.Any().Enum([]interface{}{1, "two", true, nil})
```

#### `Const(value interface{}) *AnySchema`
Restricts to a single exact value.

```go
schema.Any().Const(42)
```

### Required/Optional/Nullable

#### `Required(errorMessage ...interface{}) *AnySchema`
Marks the value as required.

```go
schema.Any().Required()
schema.Any().Required(i18n.S("field is required"))
```

#### `Optional() *AnySchema`
Marks the value as optional (default).

```go
schema.Any().Optional()
```

#### `Nullable() *AnySchema`
Allows null values (default for Any).

```go
schema.Any().Nullable()
```

### Metadata

#### `Title(title string) *AnySchema`
Sets the schema title.

```go
schema.Any().Title("Dynamic Value")
```

#### `Description(description string) *AnySchema`
Sets the schema description.

```go
schema.Any().Description("Accepts any type of value")
```

#### `Default(value interface{}) *AnySchema`
Sets a default value.

```go
schema.Any().Default(0)
```

#### `Example(example interface{}) *AnySchema`
Adds an example value.

```go
schema.Any().Example("example value").Example(42)
```

## Usage Examples

### Dynamic Configuration

```go
configSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("settings", schema.Any()) // Accept any settings structure

ctx := schema.DefaultValidationContext()

// Valid with different settings types
result := configSchema.Parse(map[string]interface{}{
    "name":     "MyApp",
    "settings": map[string]interface{}{"debug": true, "port": 8080},
}, ctx)

result = configSchema.Parse(map[string]interface{}{
    "name":     "MyApp",
    "settings": []string{"option1", "option2"},
}, ctx)

result = configSchema.Parse(map[string]interface{}{
    "name":     "MyApp",
    "settings": "simple-mode",
}, ctx)
```

### Metadata Fields

```go
userSchema := schema.Object().
    Property("id", schema.Int().Required()).
    Property("name", schema.String().Required()).
    Property("metadata", schema.Any()) // Flexible metadata

result := userSchema.Parse(map[string]interface{}{
    "id":   1,
    "name": "Alice",
    "metadata": map[string]interface{}{
        "lastLogin": "2025-11-17",
        "preferences": []string{"email", "sms"},
        "customData": 42,
    },
}, ctx)
```

### Mixed Array

```go
mixedArraySchema := schema.Array(schema.Any())

// Array can contain any types
result := mixedArraySchema.Parse([]interface{}{
    1,
    "string",
    true,
    map[string]interface{}{"key": "value"},
    []int{1, 2, 3},
}, ctx)
```

### API Response Data

```go
apiResponseSchema := schema.Object().
    Property("status", schema.String().Required()).
    Property("data", schema.Any()) // Data can be anything

// Success with object
result := apiResponseSchema.Parse(map[string]interface{}{
    "status": "success",
    "data":   map[string]interface{}{"id": 1, "name": "Alice"},
}, ctx)

// Success with array
result = apiResponseSchema.Parse(map[string]interface{}{
    "status": "success",
    "data":   []interface{}{1, 2, 3},
}, ctx)

// Success with string
result = apiResponseSchema.Parse(map[string]interface{}{
    "status": "success",
    "data":   "Operation completed",
}, ctx)
```

### Plugin System

```go
pluginSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("enabled", schema.Bool().Required()).
    Property("config", schema.Any()) // Plugin-specific config

// Each plugin has different config structure
result := pluginSchema.Parse(map[string]interface{}{
    "name":    "EmailPlugin",
    "enabled": true,
    "config": map[string]interface{}{
        "smtp": "smtp.example.com",
        "port": 587,
    },
}, ctx)
```

### Enum of Mixed Types

```go
statusSchema := schema.Any().Enum([]interface{}{
    "pending",
    "approved",
    "rejected",
    1,    // Numeric status
    2,
    3,
    nil,  // Unknown status
})

result := statusSchema.Parse("pending", ctx)  // Valid
result = statusSchema.Parse(1, ctx)           // Valid
result = statusSchema.Parse(nil, ctx)         // Valid
result = statusSchema.Parse("unknown", ctx)   // Invalid: not in enum
```

### Logging Context

```go
logSchema := schema.Object().
    Property("level", schema.String().Enum([]string{"info", "warn", "error"})).
    Property("message", schema.String().Required()).
    Property("context", schema.Any()) // Any additional context data

result := logSchema.Parse(map[string]interface{}{
    "level":   "error",
    "message": "Database connection failed",
    "context": map[string]interface{}{
        "host":    "db.example.com",
        "attempt": 3,
        "error":   "timeout",
    },
}, ctx)
```

### Optional Dynamic Field

```go
documentSchema := schema.Object().
    Property("title", schema.String().Required()).
    Property("content", schema.String().Required()).
    Property("extra", schema.Any().Optional()) // Optional dynamic field

// Valid without extra
result := documentSchema.Parse(map[string]interface{}{
    "title":   "My Document",
    "content": "Document content",
}, ctx)

// Valid with extra
result = documentSchema.Parse(map[string]interface{}{
    "title":   "My Document",
    "content": "Document content",
    "extra":   []string{"tag1", "tag2"},
}, ctx)
```

### Migration Support

```go
// Support old and new API versions
migrationSchema := schema.Object().
    Property("version", schema.Int().Required()).
    Property("data", schema.Any()) // Different structure per version

// V1 format
result := migrationSchema.Parse(map[string]interface{}{
    "version": 1,
    "data":    "simple string data",
}, ctx)

// V2 format
result = migrationSchema.Parse(map[string]interface{}{
    "version": 2,
    "data":    map[string]interface{}{"structured": "data"},
}, ctx)
```

### Webhook Payload

```go
webhookSchema := schema.Object().
    Property("event", schema.String().Required()).
    Property("timestamp", schema.DateTime().Required()).
    Property("payload", schema.Any().Required()) // Event-specific data

result := webhookSchema.Parse(map[string]interface{}{
    "event":     "user.created",
    "timestamp": "2025-11-17T14:30:00Z",
    "payload": map[string]interface{}{
        "userId":   123,
        "username": "johndoe",
        "email":    "john@example.com",
    },
}, ctx)
```

## When to Use

Any schemas are ideal for:

- **Dynamic Configuration**: Flexible settings with unknown structure
- **Plugin Systems**: Plugin-specific configuration
- **API Responses**: Variable response data structures
- **Metadata Fields**: Arbitrary additional data
- **Migration/Versioning**: Supporting multiple data formats
- **Logging Context**: Dynamic logging information
- **Webhook Payloads**: Event-specific data
- **Mixed Collections**: Arrays with heterogeneous types

## When NOT to Use

Avoid Any schemas when:
- You know the expected type - use specific schema instead
- You need type safety and validation
- Data structure is predictable and consistent
- You want to enforce data contracts

## Error Handling

```go
result := anySchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes: "required", "enum", "const"
    }
}
```

## Internationalization

```go
anySchema := schema.Any().
    Required(i18n.S("value is required")).
    Enum([]interface{}{1, 2, 3})
```

## JSON Schema Generation

```go
anySchema := schema.Any().
    Title("Dynamic Value").
    Description("Accepts any type")

jsonSchema := anySchema.JSON()
// Outputs:
// {
//   "title": "Dynamic Value",
//   "description": "Accepts any type"
// }
// Note: No "type" field means accepts anything
```

## Related

- [Object Schema](object.md) - For structured data with known properties
- [Union Schema](union.md) - For multiple specific type options
- [Transform Schema](transform.md) - For type conversion with validation