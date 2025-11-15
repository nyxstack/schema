# Transform Schema

The `TransformSchema` validates input data, applies a transformation function, then validates the transformed output. This is useful for data coercion, normalization, and complex type conversions.

## Creating a Transform Schema

```go
import "github.com/nyxstack/schema"

// Transform validates input → applies function → validates output
transformSchema := schema.Transform(
    inputSchema,    // Schema to validate input
    outputSchema,   // Schema to validate output
    transformFunc,  // Function to transform input to output
)
```

## Transform Function Signature

```go
type TransformFunc func(input interface{}) (interface{}, error)
```

The transform function:
- Takes validated input value
- Returns transformed output value or error
- Error stops validation and returns transform error

## Methods

### Core Methods

#### `Transform(inputSchema, outputSchema Parseable, transformFunc TransformFunc, errorMessage ...interface{}) *TransformSchema`
Creates a new transform schema.

```go
schema.Transform(
    schema.String(),           // Input: string
    schema.Int(),              // Output: int
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        return strconv.Atoi(str)
    },
)
```

### Required/Optional/Nullable

#### `Required(errorMessage ...interface{}) *TransformSchema`
Marks the schema as required.

```go
schema.Transform(inputSchema, outputSchema, transformFunc).
    Required(i18n.S("value is required"))
```

#### `Optional() *TransformSchema`
Marks the schema as optional (default).

```go
schema.Transform(inputSchema, outputSchema, transformFunc).Optional()
```

#### `Nullable() *TransformSchema`
Allows null values.

```go
schema.Transform(inputSchema, outputSchema, transformFunc).Nullable()
```

### Error Customization

#### `WithTransformError(errorMessage ...interface{}) *TransformSchema`
Sets custom error message for transformation failures.

```go
schema.Transform(inputSchema, outputSchema, transformFunc).
    WithTransformError(i18n.S("failed to convert value"))
```

### Metadata

#### `Title(title string) *TransformSchema`
Sets the schema title.

```go
schema.Transform(inputSchema, outputSchema, transformFunc).
    Title("Age from String")
```

#### `Description(description string) *TransformSchema`
Sets the schema description.

```go
schema.Transform(inputSchema, outputSchema, transformFunc).
    Description("Converts string age to integer")
```

#### `Default(value interface{}) *TransformSchema`
Sets a default value.

```go
schema.Transform(inputSchema, outputSchema, transformFunc).
    Default("0")
```

## Usage Examples

### String to Integer

```go
import "strconv"

stringToIntSchema := schema.Transform(
    schema.String().Pattern("^[0-9]+$"), // Input: numeric string
    schema.Int().Min(0),                  // Output: positive integer
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        return strconv.Atoi(str)
    },
)

ctx := schema.DefaultValidationContext()

result := stringToIntSchema.Parse("42", ctx)
// result.Value = 42 (int)

result = stringToIntSchema.Parse("abc", ctx) // Invalid: not numeric string
result = stringToIntSchema.Parse("-5", ctx)  // Invalid: output validation fails (Min(0))
```

### Lowercase Transformation

```go
import "strings"

lowercaseSchema := schema.Transform(
    schema.String().MinLength(1),
    schema.String().Pattern("^[a-z]+$"),
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        return strings.ToLower(str), nil
    },
)

result := lowercaseSchema.Parse("HELLO", ctx)
// result.Value = "hello"
```

### Trim Whitespace

```go
trimSchema := schema.Transform(
    schema.String(),
    schema.String().MinLength(1),
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        return strings.TrimSpace(str), nil
    },
)

result := trimSchema.Parse("  hello  ", ctx)
// result.Value = "hello"
```

### Parse JSON String

```go
import "encoding/json"

jsonParseSchema := schema.Transform(
    schema.String().MinLength(1),
    schema.Object().
        Property("name", schema.String()).
        Property("age", schema.Int()),
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        var result map[string]interface{}
        err := json.Unmarshal([]byte(str), &result)
        return result, err
    },
)

result := jsonParseSchema.Parse(`{"name":"Alice","age":30}`, ctx)
// result.Value = map[string]interface{}{"name": "Alice", "age": 30}
```

### Unix Timestamp to Date

```go
import "time"

timestampToDateSchema := schema.Transform(
    schema.Int().Min(0),
    schema.String().Pattern("^\\d{4}-\\d{2}-\\d{2}$"),
    func(input interface{}) (interface{}, error) {
        timestamp := input.(int)
        t := time.Unix(int64(timestamp), 0)
        return t.Format("2006-01-02"), nil
    },
)

result := timestampToDateSchema.Parse(1700234400, ctx)
// result.Value = "2023-11-17"
```

### Array Length to Integer

```go
arrayLengthSchema := schema.Transform(
    schema.Array(schema.Any()),
    schema.Int().Min(0),
    func(input interface{}) (interface{}, error) {
        arr := input.([]interface{})
        return len(arr), nil
    },
)

result := arrayLengthSchema.Parse([]interface{}{"a", "b", "c"}, ctx)
// result.Value = 3
```

### Normalize Phone Number

```go
import "regexp"

normalizePhoneSchema := schema.Transform(
    schema.String().Pattern("^[0-9\\-\\s\\(\\)\\+]+$"),
    schema.String().Pattern("^[0-9]{10}$"),
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        // Remove all non-digit characters
        re := regexp.MustCompile("[^0-9]")
        normalized := re.ReplaceAllString(str, "")
        return normalized, nil
    },
)

result := normalizePhoneSchema.Parse("(555) 123-4567", ctx)
// result.Value = "5551234567"

result = normalizePhoneSchema.Parse("+1-555-123-4567", ctx)
// result.Value = "15551234567" (11 digits, fails output validation)
```

### CSV String to Array

```go
import "strings"

csvToArraySchema := schema.Transform(
    schema.String().MinLength(1),
    schema.Array(schema.String().MinLength(1)).MinItems(1),
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        parts := strings.Split(str, ",")
        result := make([]interface{}, len(parts))
        for i, part := range parts {
            result[i] = strings.TrimSpace(part)
        }
        return result, nil
    },
)

result := csvToArraySchema.Parse("apple, banana, cherry", ctx)
// result.Value = []interface{}{"apple", "banana", "cherry"}
```

### Boolean String to Bool

```go
stringToBoolSchema := schema.Transform(
    schema.String().Enum([]string{"true", "false", "yes", "no", "1", "0"}),
    schema.Bool(),
    func(input interface{}) (interface{}, error) {
        str := input.(string)
        switch strings.ToLower(str) {
        case "true", "yes", "1":
            return true, nil
        case "false", "no", "0":
            return false, nil
        default:
            return nil, errors.New("invalid boolean string")
        }
    },
)

result := stringToBoolSchema.Parse("yes", ctx)
// result.Value = true
```

### Form Data Normalization

```go
normalizeFormSchema := schema.Object().
    Property("email", schema.Transform(
        schema.String().Email(),
        schema.String().Email(),
        func(input interface{}) (interface{}, error) {
            return strings.ToLower(strings.TrimSpace(input.(string))), nil
        },
    )).
    Property("username", schema.Transform(
        schema.String().MinLength(3),
        schema.String().Pattern("^[a-z0-9_]+$"),
        func(input interface{}) (interface{}, error) {
            return strings.ToLower(input.(string)), nil
        },
    ))
```

## When to Use

Transform schemas are ideal for:

- **Data Normalization**: Lowercase, trim, format conversions
- **Type Coercion**: String to int, timestamp to date
- **Parsing**: JSON strings, CSV, serialized data
- **Sanitization**: Remove unwanted characters, normalize formats
- **Computed Values**: Calculate derived values from input
- **Legacy API Compatibility**: Transform old formats to new
- **User Input Processing**: Normalize and validate user-submitted data

## Error Handling

```go
result := transformSchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes:
        // - "required" - Value is required
        // - "input_*" - Input validation failed (prefixed)
        // - "transform" - Transformation function failed
        // - "output_*" - Output validation failed (prefixed)
    }
}
```

## Internationalization

```go
transformSchema := schema.Transform(
    schema.String(),
    schema.Int(),
    func(input interface{}) (interface{}, error) {
        return strconv.Atoi(input.(string))
    },
).WithTransformError(i18n.S("failed to convert string to number"))
```

## JSON Schema Generation

```go
transformSchema := schema.Transform(
    schema.String(),
    schema.Int(),
    transformFunc,
)

jsonSchema := transformSchema.JSON()
// Outputs:
// {
//   "type": "transform",
//   "inputSchema": {"type": "string"},
//   "outputSchema": {"type": "integer"}
// }
```

## Related

- [String Schema](string.md) - Common input type for transforms
- [Object Schema](object.md) - For transforming object properties
- [Array Schema](array.md) - For transforming array items
- [Conditional Schema](conditional.md) - For conditional transformations