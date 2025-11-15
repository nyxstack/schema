# UUID Schema

The `UUIDSchema` provides comprehensive UUID (Universally Unique Identifier) validation with support for different UUID versions, formats, and case sensitivity.

## Creating a UUID Schema

```go
import "github.com/nyxstack/schema"

// Basic UUID validation (any version, any format)
basicUUID := schema.UUID()

// UUID v4 only
v4UUID := schema.UUID().Version(schema.UUIDVersion4)

// Specific format
hyphenatedUUID := schema.UUID().Format(schema.UUIDFormatHyphenated)

// Case-sensitive lowercase
lowercaseUUID := schema.UUID().Lowercase()
```

## UUID Versions

The library supports all standard UUID versions:

| Version | Constant | Description |
|---------|----------|-------------|
| Any | `UUIDVersionAny` | Accept any valid UUID format (default) |
| 1 | `UUIDVersion1` | Time-based UUID |
| 2 | `UUIDVersion2` | DCE Security UUID |
| 3 | `UUIDVersion3` | Name-based using MD5 hash |
| 4 | `UUIDVersion4` | Random UUID (most common) |
| 5 | `UUIDVersion5` | Name-based using SHA-1 hash |
| 6 | `UUIDVersion6` | Reordered time-based UUID |
| 7 | `UUIDVersion7` | Unix timestamp-based UUID |
| 8 | `UUIDVersion8` | Custom/vendor-specific UUID |

## UUID Formats

| Format | Constant | Example |
|--------|----------|---------|
| Hyphenated | `UUIDFormatHyphenated` | `550e8400-e29b-41d4-a716-446655440000` |
| Compact | `UUIDFormatCompact` | `550e8400e29b41d4a716446655440000` |
| Braced | `UUIDFormatBraced` | `{550e8400-e29b-41d4-a716-446655440000}` |
| URN | `UUIDFormatURN` | `urn:uuid:550e8400-e29b-41d4-a716-446655440000` |
| Any | `UUIDFormatAny` | Accept any of the above (default) |

## Methods

### Version Validation

#### `Version(version UUIDVersion) *UUIDSchema`
Specifies the required UUID version.

```go
v4UUID := schema.UUID().Version(schema.UUIDVersion4)
timeBasedUUID := schema.UUID().Version(schema.UUIDVersion1)
```

### Format Validation

#### `Format(format UUIDFormat) *UUIDSchema`
Specifies the required UUID format.

```go
hyphenatedOnly := schema.UUID().Format(schema.UUIDFormatHyphenated)
compactOnly := schema.UUID().Format(schema.UUIDFormatCompact)
```

### Case Validation

#### `CaseSensitive() *UUIDSchema`
Enables case-sensitive validation (UUID must match exact case).

```go
caseSensitive := schema.UUID().CaseSensitive()
```

#### `Lowercase() *UUIDSchema`
Forces UUID to be lowercase. Automatically transforms output to lowercase.

```go
lowercase := schema.UUID().Lowercase()
```

#### `Uppercase() *UUIDSchema`
Forces UUID to be uppercase. Automatically transforms output to uppercase.

```go
uppercase := schema.UUID().Uppercase()
```

### Error Messages

#### `FormatError(err ErrorMessage) *UUIDSchema`
Sets custom error message for format validation failures.

```go
schema.UUID().
    Format(schema.UUIDFormatHyphenated).
    FormatError(i18n.S("UUID must use standard hyphenated format"))
```

#### `VersionError(err ErrorMessage) *UUIDSchema`
Sets custom error message for version validation failures.

```go
schema.UUID().
    Version(schema.UUIDVersion4).
    VersionError(i18n.S("only random UUIDs (version 4) are accepted"))
```

#### `CaseError(err ErrorMessage) *UUIDSchema`
Sets custom error message for case validation failures.

```go
schema.UUID().
    Lowercase().
    CaseError(i18n.S("UUID must be in lowercase"))
```

## Usage Examples

### Basic UUID Validation

```go
uuidSchema := schema.UUID()

ctx := schema.DefaultValidationContext()

// Valid UUIDs (any format, any version)
result := uuidSchema.Parse("550e8400-e29b-41d4-a716-446655440000", ctx) // Valid
result = uuidSchema.Parse("550e8400e29b41d4a716446655440000", ctx)       // Valid (compact)
result = uuidSchema.Parse("{550e8400-e29b-41d4-a716-446655440000}", ctx) // Valid (braced)

// Invalid
result = uuidSchema.Parse("not-a-uuid", ctx)        // Invalid
result = uuidSchema.Parse("550e8400-e29b-41d4", ctx) // Invalid (incomplete)
```

### UUID Version 4 Only

```go
v4Schema := schema.UUID().Version(schema.UUIDVersion4)

// Valid UUID v4
result := v4Schema.Parse("550e8400-e29b-41d4-a716-446655440000", ctx) // Valid

// Invalid: UUID v1 (time-based)
result = v4Schema.Parse("c232ab00-9414-11ec-b3c8-9f68deced2fc", ctx) // Invalid version
```

### Strict Format Requirements

```go
hyphenatedSchema := schema.UUID().Format(schema.UUIDFormatHyphenated)

// Valid
result := hyphenatedSchema.Parse("550e8400-e29b-41d4-a716-446655440000", ctx) // Valid

// Invalid formats
result = hyphenatedSchema.Parse("550e8400e29b41d4a716446655440000", ctx)       // Invalid (compact)
result = hyphenatedSchema.Parse("{550e8400-e29b-41d4-a716-446655440000}", ctx) // Invalid (braced)
```

### Case Enforcement

```go
lowercaseSchema := schema.UUID().Lowercase()

// Valid
result := lowercaseSchema.Parse("550e8400-e29b-41d4-a716-446655440000", ctx)
fmt.Println(result.Value) // "550e8400-e29b-41d4-a716-446655440000"

// Invalid
result = lowercaseSchema.Parse("550E8400-E29B-41D4-A716-446655440000", ctx) // Invalid (uppercase)
```

### Automatic Case Transformation

```go
uppercaseSchema := schema.UUID().Uppercase()

// Input is lowercase, output is transformed to uppercase
result := uppercaseSchema.Parse("550e8400-e29b-41d4-a716-446655440000", ctx)
if result.Valid {
    fmt.Println(result.Value) // "550E8400-E29B-41D4-A716-446655440000"
}
```

### Database ID Validation

```go
// User ID: UUID v4, hyphenated format, lowercase
userIDSchema := schema.UUID().
    Version(schema.UUIDVersion4).
    Format(schema.UUIDFormatHyphenated).
    Lowercase()

userSchema := schema.Object().
    Property("id", userIDSchema).
    Property("username", schema.String().MinLength(3))

result := userSchema.Parse(map[string]interface{}{
    "id":       "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
}, ctx)
```

### API Request Validation

```go
createResourceSchema := schema.Object().
    Property("resourceId", schema.UUID().
        Version(schema.UUIDVersion4).
        FormatError(i18n.S("resource ID must be a valid UUID v4"))).
    Property("name", schema.String().MinLength(1))

// Valid request
result := createResourceSchema.Parse(map[string]interface{}{
    "resourceId": "123e4567-e89b-12d3-a456-426614174000",
    "name":       "My Resource",
}, ctx)
```

### Multiple UUID Fields with Different Requirements

```go
orderSchema := schema.Object().
    Property("orderId", schema.UUID().
        Version(schema.UUIDVersion7).
        VersionError(i18n.S("order ID must use timestamp-based UUID v7"))).
    Property("customerId", schema.UUID().
        Version(schema.UUIDVersion4)).
    Property("trackingId", schema.UUID().
        Format(schema.UUIDFormatHyphenated))
```

### Compact Format for URLs

```go
// Compact UUIDs for shorter URLs
urlParamSchema := schema.UUID().
    Format(schema.UUIDFormatCompact).
    Lowercase()

// Valid: 550e8400e29b41d4a716446655440000
// Invalid: 550e8400-e29b-41d4-a716-446655440000 (has hyphens)
```

## When to Use

UUID schemas are ideal for:

- **Database Primary Keys**: Validating UUID-based identifiers
- **API Resources**: Ensuring resource IDs are valid UUIDs
- **Distributed Systems**: Validating globally unique identifiers
- **Session IDs**: Validating session tokens
- **File Names**: Ensuring unique file identifiers
- **Tracking IDs**: Order tracking, shipment tracking, etc.

## Error Handling

```go
result := uuidSchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes: "format", "version", "case"
    }
}
```

## Internationalization

```go
uuidSchema := schema.UUID().
    Version(schema.UUIDVersion4).
    Format(schema.UUIDFormatHyphenated).
    FormatError(i18n.S("must be a hyphenated UUID")).
    VersionError(i18n.S("must be a random UUID (version 4)")).
    Lowercase().
    CaseError(i18n.S("UUID must be lowercase"))
```

## JSON Schema Generation

```go
uuidSchema := schema.UUID().
    Version(schema.UUIDVersion4).
    Format(schema.UUIDFormatHyphenated)

jsonSchema := uuidSchema.JSON()
// Outputs:
// {
//   "type": "string",
//   "format": "uuid",
//   "pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
// }
```

## Related

- [String Schema](string.md) - For general string validation (includes UUID format)
- [Object Schema](object.md) - For validating objects with UUID properties
- [Transform Schema](transform.md) - For transforming UUIDs to different formats