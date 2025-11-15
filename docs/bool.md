# Boolean Schema

The `BoolSchema` provides validation for boolean (true/false) values with support for const values and enum constraints.

## Creating a Boolean Schema

```go
import "github.com/nyxstack/schema"

// Basic bool schema
activeSchema := schema.Bool()

// With validation constraints
termsSchema := schema.Bool().
    Const(true, "Must accept terms").
    Required("Terms acceptance is required")
```

## Methods

### Type Configuration

#### `Required(messages ...ErrorMessage) *BoolSchema`
Marks the boolean as required (cannot be nil or omitted).

```go
schema.Bool().Required()
schema.Bool().Required("Field is required")
schema.Bool().Required(i18n.S("field is required"))
```

#### `Optional() *BoolSchema`
Marks the boolean as optional (can be nil or omitted).

```go
schema.Bool().Optional()
```

#### `Nullable() *BoolSchema`
Allows the boolean value to be explicitly null.

```go
schema.Bool().Nullable()
```

#### `Default(value interface{}) *BoolSchema`
Sets a default value when the input is nil.

```go
schema.Bool().Default(false)
schema.Bool().Default(true)
```

### Value Constraints

#### `Const(value bool, messages ...ErrorMessage) *BoolSchema`
Requires the boolean to match an exact value. Useful for "must be true" validations.

```go
// Must be true
schema.Bool().Const(true, "Must accept")

// Must be false
schema.Bool().Const(false, "Must decline")

// With i18n
schema.Bool().Const(true, i18n.S("must accept terms and conditions"))
```

#### `Enum(values []bool, messages ...ErrorMessage) *BoolSchema`
Restricts the boolean to a set of allowed values (typically `[]bool{true}` or `[]bool{false}`).

```go
// Only true is valid
schema.Bool().Enum([]bool{true})

// Only false is valid
schema.Bool().Enum([]bool{false}, "Must be false")
```

### Metadata

#### `Title(title string) *BoolSchema`
Sets a title for documentation and JSON Schema generation.

```go
schema.Bool().Title("Active Status")
```

#### `Description(description string) *BoolSchema`
Sets a description for documentation and JSON Schema generation.

```go
schema.Bool().Description("Whether the user account is active")
```

## Usage Examples

### Basic Boolean Validation

```go
activeSchema := schema.Bool().
    Required("Active status is required")

ctx := schema.DefaultValidationContext()
result := activeSchema.Parse(true, ctx)

if result.Valid {
    fmt.Printf("Active: %v\n", result.Value)
}
```

### Terms and Conditions

```go
termsSchema := schema.Bool().
    Const(true, "You must accept the terms and conditions").
    Required()

result := termsSchema.Parse(true, ctx)  // Valid
result := termsSchema.Parse(false, ctx) // Invalid - must be true
```

### Newsletter Subscription (Optional)

```go
newsletterSchema := schema.Bool().
    Default(false).
    Optional().
    Title("Newsletter Subscription")

result := newsletterSchema.Parse(nil, ctx) // Uses default: false
```

### Feature Flag

```go
featureSchema := schema.Bool().
    Default(false).
    Title("Feature Enabled").
    Description("Whether this feature is enabled for the user")

result := featureSchema.Parse(nil, ctx) // Returns false
```

### Email Verified Status

```go
verifiedSchema := schema.Bool().
    Default(false).
    Title("Email Verified")

result := verifiedSchema.Parse(false, ctx)
```

### Privacy Settings

```go
publicProfileSchema := schema.Bool().
    Default(true).
    Optional().
    Title("Public Profile").
    Description("Whether the profile is visible to other users")

result := publicProfileSchema.Parse(nil, ctx) // Returns true
```

### Only True Allowed (Enum)

```go
agreeSchema := schema.Bool().
    Enum([]bool{true}, "You must agree to continue").
    Required()

result := agreeSchema.Parse(true, ctx)  // Valid
result := agreeSchema.Parse(false, ctx) // Invalid
```

### Complex Validation with i18n

```go
termsSchema := schema.Bool().
    Const(true, i18n.S("must accept terms and conditions")).
    Required(i18n.S("terms acceptance is required"))

result := termsSchema.Parse(true, ctx)
```

### Nullable Boolean (Three States)

```go
// Represents: Yes (true), No (false), or Not Answered (nil)
answerSchema := schema.Bool().
    Nullable().
    Optional()

result := answerSchema.Parse(nil, ctx)   // Valid - nil
result := answerSchema.Parse(true, ctx)  // Valid - true
result := answerSchema.Parse(false, ctx) // Valid - false
```

### Type Coercion

The boolean schema does NOT automatically coerce non-boolean types. Values must be actual booleans:

```go
schema := schema.Bool()

// These work
schema.Parse(true, ctx)
schema.Parse(false, ctx)

// These fail (no automatic coercion)
schema.Parse(1, ctx)      // Invalid
schema.Parse(0, ctx)      // Invalid
schema.Parse("true", ctx) // Invalid
schema.Parse("false", ctx) // Invalid
```

## Common Use Cases

### Checkboxes in Forms

```go
checkboxSchema := schema.Bool().
    Default(false).
    Optional()
```

### Required Agreements

```go
agreementSchema := schema.Bool().
    Const(true, "Must agree to proceed").
    Required()
```

### Account Settings

```go
notificationsSchema := schema.Bool().
    Default(true).
    Optional().
    Title("Email Notifications")
```

### Admin Flags

```go
isAdminSchema := schema.Bool().
    Default(false).
    Title("Administrator")
```

## Internationalization

All error messages support i18n through the `github.com/nyxstack/i18n` package:

```go
schema.Bool().
    Const(true, i18n.S("must accept terms")).
    Required(i18n.S("acceptance is required"))
```

## JSON Schema Generation

```go
schema := schema.Bool().
    Title("Email Verified").
    Description("Whether the email has been verified").
    Default(false).
    Required()

jsonSchema := schema.JSON()
// Outputs:
// {
//   "type": "boolean",
//   "default": false,
//   "title": "Email Verified",
//   "description": "Whether the email has been verified"
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

- [String Schema](string.md) - For string values
- [Int Schema](int.md) - For integer values
- [Object Schema](object.md) - For objects with boolean properties
- [Array Schema](array.md) - For arrays of booleans
