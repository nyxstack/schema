# Binary Schema

The `BinarySchema` validates binary data encoded as strings in various formats including base64, base64url, and hexadecimal.

## Creating a Binary Schema

```go
import "github.com/nyxstack/schema"

// Base64 encoded binary data (default)
binarySchema := schema.Binary()

// Explicitly specify base64
base64Schema := schema.Base64()

// URL-safe base64
base64URLSchema := schema.Base64URL()

// Hexadecimal encoding
hexSchema := schema.Hex()
```

## Supported Formats

| Format | Constant | Description | Example |
|--------|----------|-------------|---------|
| Base64 | `BinaryFormatBase64` | Standard base64 encoding (default) | `SGVsbG8gV29ybGQ=` |
| Base64URL | `BinaryFormatBase64URL` | URL-safe base64 (no padding) | `SGVsbG8gV29ybGQ` |
| Hex | `BinaryFormatHex` | Hexadecimal encoding | `48656c6c6f20576f726c64` |

## Methods

### Format Selection

#### `Binary() *BinarySchema`
Creates a new binary schema with base64 encoding (default).

```go
binary := schema.Binary()
```

#### `Base64() *BinarySchema`
Creates a new binary schema with standard base64 encoding.

```go
base64 := schema.Base64()
```

#### `Base64URL() *BinarySchema`
Creates a new binary schema with URL-safe base64 encoding.

```go
base64URL := schema.Base64URL()
```

#### `Hex() *BinarySchema`
Creates a new binary schema with hexadecimal encoding.

```go
hex := schema.Hex()
```

#### `Format(format BinaryFormat) *BinarySchema`
Sets the binary encoding format explicitly.

```go
schema.Binary().Format(schema.BinaryFormatHex)
```

### Size Constraints

#### `MinSize(min int) *BinarySchema`
Sets the minimum size constraint in bytes (decoded data).

```go
schema.Binary().MinSize(10) // At least 10 bytes
```

#### `MaxSize(max int) *BinarySchema`
Sets the maximum size constraint in bytes (decoded data).

```go
schema.Binary().MaxSize(1024) // At most 1KB
```

#### `Size(min, max int) *BinarySchema`
Sets both minimum and maximum size constraints.

```go
schema.Binary().Size(100, 1000) // Between 100 and 1000 bytes
```

### Required Validation

#### `Required() *BinarySchema`
Marks the binary data as required (non-empty).

```go
schema.Binary().Required()
```

### Error Customization

#### `FormatError(err ErrorMessage) *BinarySchema`
Sets custom error message for format validation failures.

```go
schema.Binary().
    FormatError(i18n.S("file must be base64 encoded"))
```

#### `SizeError(err ErrorMessage) *BinarySchema`
Sets custom error message for size validation failures.

```go
schema.Binary().
    MaxSize(1048576).
    SizeError(i18n.S("file size must not exceed 1MB"))
```

## Usage Examples

### Basic Base64 Validation

```go
base64Schema := schema.Base64()

ctx := schema.DefaultValidationContext()

// Valid base64
result := base64Schema.Parse("SGVsbG8gV29ybGQ=", ctx)  // "Hello World"
result = base64Schema.Parse("Zm9vYmFy", ctx)            // "foobar"

// Invalid
result = base64Schema.Parse("not-base64!", ctx)         // Invalid format
result = base64Schema.Parse("SGVsbG8=", ctx)            // Invalid padding
```

### File Upload Validation

```go
imageSchema := schema.Base64().
    MinSize(100).
    MaxSize(5 * 1024 * 1024). // 5MB max
    FormatError(i18n.S("image must be base64 encoded")).
    SizeError(i18n.S("image size must be between 100 bytes and 5MB"))

uploadSchema := schema.Object().
    Property("filename", schema.String().MinLength(1)).
    Property("data", imageSchema)

result := uploadSchema.Parse(map[string]interface{}{
    "filename": "photo.jpg",
    "data":     "iVBORw0KGgoAAAANSUhEUgAAAAUA...",
}, ctx)
```

### Hexadecimal Hash Validation

```go
// SHA-256 hash (64 hex characters = 32 bytes)
sha256Schema := schema.Hex().
    Size(32, 32).
    Required()

// Valid SHA-256 hash
result := sha256Schema.Parse(
    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", 
    ctx,
)

// Invalid length
result = sha256Schema.Parse("abc123", ctx) // Too short
```

### URL-Safe Base64 (JWT Payload)

```go
jwtPayloadSchema := schema.Base64URL()

// Valid base64url (no padding)
result := jwtPayloadSchema.Parse(
    "eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
    ctx,
)
```

### Certificate Validation

```go
certificateSchema := schema.Base64().
    MinSize(256).
    Required().
    FormatError(i18n.S("certificate must be PEM base64 encoded"))

configSchema := schema.Object().
    Property("cert", certificateSchema).
    Property("key", certificateSchema)
```

### API Key Validation

```go
// API keys as hex strings (16 bytes = 32 hex chars)
apiKeySchema := schema.Hex().
    Size(16, 16).
    Required()

// Valid API key
result := apiKeySchema.Parse("a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6", ctx)
```

### Document Storage

```go
documentSchema := schema.Object().
    Property("title", schema.String().MinLength(1)).
    Property("content", schema.Base64().
        MaxSize(10 * 1024 * 1024). // 10MB
        Required()).
    Property("checksum", schema.Hex().Size(32, 32)) // SHA-256

result := documentSchema.Parse(map[string]interface{}{
    "title":    "Report.pdf",
    "content":  "JVBERi0xLjQKJeLjz9MKMyAw...",
    "checksum": "abc123...",
}, ctx)
```

### Binary Attachment in Email

```go
attachmentSchema := schema.Object().
    Property("filename", schema.String().MinLength(1)).
    Property("mimeType", schema.String()).
    Property("data", schema.Base64().
        MaxSize(25 * 1024 * 1024). // 25MB email limit
        Required())

emailSchema := schema.Object().
    Property("to", schema.String().Email()).
    Property("subject", schema.String().MinLength(1)).
    Property("attachments", schema.Array(attachmentSchema).MaxItems(10))
```

### Encryption Key Validation

```go
// AES-256 key (32 bytes)
aes256KeySchema := schema.Hex().
    Size(32, 32).
    Required().
    SizeError(i18n.S("AES-256 requires a 32-byte key"))

// IV (16 bytes)
ivSchema := schema.Hex().
    Size(16, 16).
    Required()

encryptionSchema := schema.Object().
    Property("key", aes256KeySchema).
    Property("iv", ivSchema)
```

### Image Thumbnail

```go
thumbnailSchema := schema.Base64().
    MaxSize(100 * 1024). // 100KB
    Required().
    SizeError(i18n.S("thumbnail must be less than 100KB"))

profileSchema := schema.Object().
    Property("username", schema.String().MinLength(3)).
    Property("avatar", thumbnailSchema)
```

## When to Use

Binary schemas are ideal for:

- **File Uploads**: Validating base64-encoded files
- **Cryptographic Data**: Hashes, signatures, keys in hex format
- **API Tokens**: Binary tokens encoded as strings
- **Image Data**: Inline images in JSON/API requests
- **Certificates**: PEM-encoded certificates and keys
- **Document Storage**: Binary documents encoded for transport
- **Email Attachments**: File attachments in email APIs

## Error Handling

```go
result := binarySchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes: "invalid_type", "required", "format", "min_size", "max_size"
    }
}
```

## Internationalization

```go
binarySchema := schema.Base64().
    Required().
    MaxSize(5 * 1024 * 1024).
    FormatError(i18n.S("must be valid base64 encoded data")).
    SizeError(i18n.F("file size must not exceed %d MB", 5))
```

## JSON Schema Generation

```go
binarySchema := schema.Base64().
    MinSize(100).
    MaxSize(1024)

jsonSchema := binarySchema.JSON()
// Outputs:
// {
//   "type": "string",
//   "contentEncoding": "base64",
//   "minLength": 100,
//   "maxLength": 1024
// }
```

## Related

- [String Schema](string.md) - For text validation
- [Object Schema](object.md) - For validating objects with binary properties
- [Array Schema](array.md) - For arrays of binary data