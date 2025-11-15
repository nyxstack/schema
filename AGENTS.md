# AI Agents Guide for Nyx Schema

This document provides guidance for AI agents (like GitHub Copilot, ChatGPT, Claude, etc.) working with the Nyx Schema library. It contains architectural patterns, common use cases, and implementation guidelines to help AI assistants provide accurate and helpful code suggestions.

## Library Architecture Overview

Nyx Schema is a fluent Go library for schema validation and JSON Schema generation. Key architectural principles:

### Core Design Patterns

1. **Fluent Interface**: All schema types support method chaining for intuitive API usage
2. **Immutable Schemas**: Schema configuration methods return new instances, allowing safe reuse
3. **Parse-Don't-Validate**: The `Parse` method both validates and transforms data, returning structured results
4. **Internationalization First**: Built-in i18n support through the `github.com/nyxstack/i18n` package
5. **Type Safety**: Strong typing with interfaces for different schema behaviors

### Schema Hierarchy

```
Schema (base)
├── StringSchema
├── IntSchema / FloatSchema / NumberSchema
├── BoolSchema
├── ArraySchema
├── ObjectSchema
├── UnionSchema (AnyOf/OneOf/AllOf)
├── ConditionalSchema
└── Specialized (UUIDSchema, DateSchema, etc.)
```

## Common Usage Patterns

### 1. Basic Schema Creation

```go
// Primitive types
stringSchema := schema.String().Required().MinLength(1)
intSchema := schema.Int().Min(0).Max(100)
boolSchema := schema.Bool().Required()

// With custom error messages
emailSchema := schema.String().
    Email().
    Required("Email is required")
```

### 2. Object Schema Patterns

```go
// Method 1: Fluent API
userSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("email", schema.String().Email().Required()).
    Property("age", schema.Int().Min(0).Optional())

// Method 2: Shape pattern (more concise)
userSchema := schema.Shape{
    "name":  schema.String().Required(),
    "email": schema.String().Email().Required(),
    "age":   schema.Int().Min(0).Optional(),
}.AsObject()
```

### 3. Array Schema Patterns

```go
// Simple array
tagsSchema := schema.Array(schema.String().MinLength(1))

// Complex array with constraints
usersSchema := schema.Array(
    schema.Object().
        Property("id", schema.Int().Required()).
        Property("name", schema.String().Required()),
).MinItems(1).MaxItems(100).UniqueItems()
```

### 4. Validation Context Usage

```go
// Always use validation context for parsing
ctx := schema.DefaultValidationContext()

// For i18n applications
ctx := schema.NewValidationContext("es") // Spanish locale

// Parse and handle results
result := mySchema.Parse(data, ctx)
if !result.Valid {
    // Handle validation errors
    for _, err := range result.Errors {
        log.Printf("Validation error: %s", err.Message)
    }
}
```

## Implementation Guidelines for AI Agents

### 1. Schema Construction Recommendations

When suggesting schema construction, prefer:

```go
// ✅ Good: Fluent, readable
userSchema := schema.String().
    MinLength(3).
    MaxLength(50).
    Pattern("^[a-zA-Z\\s]+$").
    Required()

// ❌ Avoid: Multiple separate assignments
userSchema := schema.String()
userSchema = userSchema.MinLength(3)
userSchema = userSchema.MaxLength(50)
// ... etc
```

### 2. Error Message Patterns

Always suggest custom error messages for user-facing validation:

```go
// ✅ Good: User-friendly messages
passwordSchema := schema.String().
    MinLength(8, "Password must be at least 8 characters long").
    Pattern(".*[A-Z].*", "Password must contain an uppercase letter").
    Required("Password is required")

// ✅ Also good: Default messages for internal validation
internalIdSchema := schema.Int().Min(1).Required()
```

### 3. Validation Result Handling

Always suggest proper error handling:

```go
// ✅ Good: Complete error handling
result := schema.Parse(data, ctx)
if !result.Valid {
    var errorMessages []string
    for _, err := range result.Errors {
        errorMessages = append(errorMessages, err.Message)
    }
    return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
}
// Use result.Value for validated data
```

### 4. Common Validation Scenarios

#### User Input Validation
```go
// Forms, APIs, user registration
userInputSchema := schema.Object().
    Property("username", schema.String().
        Pattern("^[a-zA-Z0-9_]{3,20}$").
        Required("Username is required")).
    Property("email", schema.String().
        Email().
        Required("Email address is required"))
```

#### Configuration Validation
```go
// Config files, environment variables
configSchema := schema.Object().
    Property("port", schema.Int().Min(1024).Max(65535).Default(8080)).
    Property("host", schema.String().Default("localhost")).
    Property("debug", schema.Bool().Default(false))
```

#### API Request/Response Validation
```go
// REST API payloads
createUserRequest := schema.Object().
    Property("user", schema.Object().
        Property("name", schema.String().MinLength(1).Required()).
        Property("email", schema.String().Email().Required())).
    AdditionalProperties(false)
```

### 5. Advanced Pattern Suggestions

#### Union Types for Flexible APIs
```go
// ID can be string or integer
idSchema := schema.AnyOf(
    schema.String().Pattern("^[a-zA-Z0-9_]+$"),
    schema.Int().Min(1),
)
```

#### Conditional Validation
```go
// Different validation based on type field
conditionalSchema := schema.Conditional().
    If(schema.Object().Property("type", schema.String().Const("premium"))).
    Then(schema.Object().Property("features", schema.Array(schema.String()).MinItems(5))).
    Else(schema.Object().Property("features", schema.Array(schema.String()).MaxItems(2)))
```

## Testing Patterns

Suggest comprehensive test patterns:

```go
func TestUserSchema(t *testing.T) {
    ctx := schema.DefaultValidationContext()
    userSchema := createUserSchema()
    
    tests := []struct {
        name     string
        input    interface{}
        expected bool
    }{
        {"valid user", map[string]interface{}{
            "name": "John Doe",
            "email": "john@example.com",
        }, true},
        {"missing name", map[string]interface{}{
            "email": "john@example.com",
        }, false},
        {"invalid email", map[string]interface{}{
            "name": "John Doe",
            "email": "invalid-email",
        }, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := userSchema.Parse(tt.input, ctx)
            if result.Valid != tt.expected {
                t.Errorf("Expected %v, got %v. Errors: %v", 
                    tt.expected, result.Valid, result.Errors)
            }
        })
    }
}
```

## JSON Schema Generation

For applications needing JSON Schema output:

```go
// Generate JSON Schema for documentation/tooling
jsonSchema := userSchema.JSON()
jsonBytes, _ := json.MarshalIndent(jsonSchema, "", "  ")

// Can be used with OpenAPI, Swagger, etc.
fmt.Println(string(jsonBytes))
```

## Performance Considerations

1. **Schema Reuse**: Create schemas once and reuse them
2. **Validation Context**: Reuse validation contexts when possible
3. **Error Allocation**: For high-throughput validation, consider error pooling

```go
// ✅ Good: Reuse schemas
var userSchema = schema.Object().
    Property("name", schema.String().Required()).
    Property("email", schema.String().Email().Required())

func validateUser(data interface{}) error {
    result := userSchema.Parse(data, ctx)
    // ... handle result
}

// ❌ Avoid: Recreating schemas repeatedly
func validateUser(data interface{}) error {
    userSchema := schema.Object().
        Property("name", schema.String().Required()).
        Property("email", schema.String().Email().Required())
    // ... validation
}
```

## Common Pitfalls to Avoid

### 1. Not Using Validation Context
```go
// ❌ Wrong: Missing context
result := schema.Parse(data) // This won't work

// ✅ Correct: Always provide context
result := schema.Parse(data, schema.DefaultValidationContext())
```

### 2. Ignoring Parse Results
```go
// ❌ Wrong: Not checking validity
result := schema.Parse(data, ctx)
processData(result.Value) // Value might be nil if invalid

// ✅ Correct: Check validity first
result := schema.Parse(data, ctx)
if result.Valid {
    processData(result.Value)
} else {
    handleErrors(result.Errors)
}
```

### 3. Mixing Required/Optional Incorrectly
```go
// ❌ Confusing: Required is default for most types
optionalField := schema.String().Required().Optional() // Contradictory

// ✅ Clear: Be explicit about requirements
requiredField := schema.String().Required()
optionalField := schema.String().Optional()
```

## Integration Examples

### With HTTP Handlers
```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var requestData interface{}
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    result := createUserSchema.Parse(requestData, schema.DefaultValidationContext())
    if !result.Valid {
        var errors []string
        for _, err := range result.Errors {
            errors = append(errors, err.Message)
        }
        http.Error(w, strings.Join(errors, "; "), http.StatusBadRequest)
        return
    }
    
    // Use result.Value which contains validated/parsed data
    createUser(result.Value)
}
```

### With Configuration Loading
```go
func loadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var configData interface{}
    if err := json.Unmarshal(data, &configData); err != nil {
        return nil, err
    }
    
    result := configSchema.Parse(configData, schema.DefaultValidationContext())
    if !result.Valid {
        return nil, fmt.Errorf("config validation failed: %v", result.Errors)
    }
    
    // Convert validated result to typed config struct
    return convertToConfig(result.Value), nil
}
```

## Version Compatibility

- Requires Go 1.24.2 or later
- Depends on `github.com/nyxstack/i18n v1.0.0`
- Follows semantic versioning for API stability

## Future Considerations

When suggesting code, be aware that future versions may include:
- Additional format validators
- Performance optimizations
- Enhanced error reporting
- More union type combinations
- Custom validation function support

Always check the latest documentation for current capabilities and best practices.