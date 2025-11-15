# String Schema

The `StringSchema` provides comprehensive validation for string values with support for length constraints, pattern matching, format validation, and more.

## Creating a String Schema

```go
import "github.com/nyxstack/schema"

// Basic string schema
nameSchema := schema.String()

// With validation constraints
emailSchema := schema.String().
    Email().
    Required("Email is required")
```

## Methods

### Type Configuration

#### `Required(messages ...ErrorMessage) *StringSchema`
Marks the string as required (cannot be nil or omitted).

```go
schema.String().Required()
schema.String().Required("Username is required")
schema.String().Required(i18n.S("username is required"))
```

#### `Optional() *StringSchema`
Marks the string as optional (can be nil or omitted).

```go
schema.String().Optional()
```

#### `Nullable() *StringSchema`
Allows the string value to be explicitly null.

```go
schema.String().Nullable()
```

#### `Default(value interface{}) *StringSchema`
Sets a default value when the input is nil.

```go
schema.String().Default("guest")
```

### Length Constraints

#### `MinLength(min int, messages ...ErrorMessage) *StringSchema`
Sets the minimum length for the string.

```go
schema.String().MinLength(3)
schema.String().MinLength(8, "Password must be at least 8 characters")
schema.String().MinLength(8, i18n.F("password must be at least %d characters", 8))
```

#### `MaxLength(max int, messages ...ErrorMessage) *StringSchema`
Sets the maximum length for the string.

```go
schema.String().MaxLength(100)
schema.String().MaxLength(20, "Username too long")
```

#### `Length(exact int, messages ...ErrorMessage) *StringSchema`
Requires an exact string length.

```go
schema.String().Length(10, "Must be exactly 10 characters")
```

### Pattern Matching

#### `Pattern(pattern string, messages ...ErrorMessage) *StringSchema`
Validates the string against a regular expression.

```go
// Alphanumeric only
schema.String().Pattern("^[a-zA-Z0-9]+$", "Must be alphanumeric")

// Phone number
schema.String().Pattern("^\\+?[1-9]\\d{1,14}$", "Invalid phone number")
```

### Format Validation

#### `Email(messages ...ErrorMessage) *StringSchema`
Validates email format.

```go
schema.String().Email()
schema.String().Email("Invalid email address")
```

#### `URL(messages ...ErrorMessage) *StringSchema`
Validates URL format.

```go
schema.String().URL()
schema.String().URL("Invalid URL")
```

#### `UUID(messages ...ErrorMessage) *StringSchema`
Validates UUID format.

```go
schema.String().UUID()
schema.String().UUID("Invalid UUID")
```

#### `DateTime(messages ...ErrorMessage) *StringSchema`
Validates ISO 8601 date-time format.

```go
schema.String().DateTime()
```

#### `Date(messages ...ErrorMessage) *StringSchema`
Validates date format (YYYY-MM-DD).

```go
schema.String().Date()
```

#### `Time(messages ...ErrorMessage) *StringSchema`
Validates time format (HH:MM:SS).

```go
schema.String().Time()
```

#### `Format(format StringFormat) *StringSchema`
Applies a format validator.

**Available Formats:**
- `StringFormatEmail` - Email address validation
- `StringFormatURI` - URI validation  
- `StringFormatURL` - URL validation
- `StringFormatDateTime` - ISO 8601 date-time (e.g., "2025-11-16T10:30:00Z")
- `StringFormatDate` - ISO 8601 date (e.g., "2025-11-16")
- `StringFormatTime` - ISO 8601 time (e.g., "10:30:00")
- `StringFormatUUID` - UUID validation (v1-v5)
- `StringFormatHostname` - Hostname validation
- `StringFormatIPv4` - IPv4 address validation
- `StringFormatIPv6` - IPv6 address validation
- `StringFormatPassword` - Password format (metadata only)
- `StringFormatBinary` - Binary data format
- `StringFormatByte` - Base64 encoded byte data

```go
schema.String().Format(schema.StringFormatIPv4)
schema.String().Format(schema.StringFormatIPv6)
schema.String().Format(schema.StringFormatHostname)
schema.String().Format(schema.StringFormatBinary)
schema.String().Format(schema.StringFormatByte)
```

### Value Constraints

#### `Enum(values []string, messages ...ErrorMessage) *StringSchema`
Restricts the string to a set of allowed values.

```go
schema.String().Enum([]string{"active", "inactive", "pending"})
schema.String().Enum([]string{"red", "green", "blue"}, "Invalid color")
```

#### `Const(value string, messages ...ErrorMessage) *StringSchema`
Requires the string to match an exact value.

```go
schema.String().Const("accepted")
schema.String().Const("yes", "Must agree to terms")
```

### Metadata

#### `Title(title string) *StringSchema`
Sets a title for documentation and JSON Schema generation.

```go
schema.String().Title("User Email")
```

#### `Description(description string) *StringSchema`
Sets a description for documentation and JSON Schema generation.

```go
schema.String().Description("The user's primary email address")
```

## Usage Examples

### Basic String Validation

```go
nameSchema := schema.String().
    MinLength(2).
    MaxLength(50).
    Required("Name is required")

ctx := schema.DefaultValidationContext()
result := nameSchema.Parse("John Doe", ctx)

if result.Valid {
    fmt.Printf("Valid name: %s\n", result.Value)
}
```

### Email Validation

```go
emailSchema := schema.String().
    Email().
    Required("Email is required")

result := emailSchema.Parse("user@example.com", ctx)
```

### Pattern Validation

```go
usernameSchema := schema.String().
    Pattern("^[a-zA-Z0-9_]{3,20}$", "Username must be 3-20 alphanumeric characters").
    Required()

result := usernameSchema.Parse("john_doe123", ctx)
```

### Enum Validation

```go
statusSchema := schema.String().
    Enum([]string{"draft", "published", "archived"}).
    Default("draft")

result := statusSchema.Parse(nil, ctx) // Uses default: "draft"
```

### Complex Validation

```go
passwordSchema := schema.String().
    MinLength(8, i18n.F("password must be at least %d characters", 8)).
    Pattern(".*[A-Z].*", i18n.S("password must contain uppercase")).
    Pattern(".*[0-9].*", i18n.S("password must contain number")).
    Pattern(".*[!@#$%^&*].*", i18n.S("password must contain special char")).
    Required(i18n.S("password is required"))

result := passwordSchema.Parse("SecurePass123!", ctx)
```

### Optional with Default

```go
roleSchema := schema.String().
    Enum([]string{"admin", "user", "guest"}).
    Default("guest").
    Optional()

result := roleSchema.Parse(nil, ctx) // Returns "guest"
```

## Internationalization

All error messages support i18n through the `github.com/nyxstack/i18n` package:

```go
schema.String().
    MinLength(3, i18n.F("username must be at least %d characters", 3)).
    Required(i18n.S("username is required"))
```

## JSON Schema Generation

```go
schema := schema.String().
    Title("User Email").
    Description("Primary email address").
    Email().
    Required()

jsonSchema := schema.JSON()
// Outputs:
// {
//   "type": "string",
//   "format": "email",
//   "title": "User Email",
//   "description": "Primary email address"
// }
```

## Error Handling

```go
result := schema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Path: %v\n", err.Path)
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
    }
}
```

## Related

- [Object Schema](object.md) - For validating objects with string properties
- [Array Schema](array.md) - For arrays of strings
- [Format Validation](formats.md) - Custom format validators
- [i18n Support](internationalization.md) - Internationalization guide
