# Ref Schema

The `RefSchema` enables schema reuse through references, allowing you to define schemas once and reference them multiple times. This is essential for recursive schemas and reducing duplication.

## Creating a Ref Schema

```go
import "github.com/nyxstack/schema"

// Create a registry
registry := schema.NewSchemaRegistry()

// Define schemas
registry.Define("User", userSchema)
registry.Define("Address", addressSchema)

// Reference schemas
userRef := schema.Ref("#/User", registry)
addressRef := schema.Ref("#/Address", registry)
```

## Schema Registry

### `NewSchemaRegistry() *SchemaRegistry`
Creates a new schema registry for storing definitions.

```go
registry := schema.NewSchemaRegistry()
```

### `Define(name string, schema Parseable)`
Adds a schema definition to the registry.

```go
registry.Define("User", schema.Object().
    Property("id", schema.Int()).
    Property("name", schema.String()))
```

### `Get(name string) (Parseable, bool)`
Retrieves a schema definition by name.

```go
userSchema, exists := registry.Get("User")
```

### `Clear()`
Removes all definitions from the registry.

```go
registry.Clear()
```

## Methods

### Core Methods

#### `Ref(ref string, registry *SchemaRegistry) *RefSchema`
Creates a new reference schema pointing to a definition.

```go
userRef := schema.Ref("#/User", registry)
```

**Reference Format**: Must start with `#/` followed by the definition name.

### Error Customization

#### `RefError(err ErrorMessage) *RefSchema`
Sets custom error message for reference resolution failures.

```go
schema.Ref("#/User", registry).
    RefError(i18n.S("user schema reference not found"))
```

## Usage Examples

### Basic Schema Reuse

```go
registry := schema.NewSchemaRegistry()

// Define address schema
addressSchema := schema.Object().
    Property("street", schema.String().Required()).
    Property("city", schema.String().Required()).
    Property("zipCode", schema.String().Pattern("^[0-9]{5}$"))

registry.Define("Address", addressSchema)

// Use address schema in multiple places
userSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("homeAddress", schema.Ref("#/Address", registry)).
    Property("workAddress", schema.Ref("#/Address", registry))

ctx := schema.DefaultValidationContext()

result := userSchema.Parse(map[string]interface{}{
    "name": "Alice",
    "homeAddress": map[string]interface{}{
        "street":  "123 Main St",
        "city":    "Springfield",
        "zipCode": "12345",
    },
    "workAddress": map[string]interface{}{
        "street":  "456 Office Blvd",
        "city":    "Springfield",
        "zipCode": "12346",
    },
}, ctx)
```

### Recursive Schema (Tree Structure)

```go
registry := schema.NewSchemaRegistry()

// Define node schema (references itself)
nodeSchema := schema.Object().
    Property("value", schema.Int().Required()).
    Property("left", schema.OneOf(
        schema.Ref("#/Node", registry),
        schema.Null(),
    )).
    Property("right", schema.OneOf(
        schema.Ref("#/Node", registry),
        schema.Null(),
    ))

registry.Define("Node", nodeSchema)

// Parse tree structure
treeSchema := schema.Ref("#/Node", registry)

result := treeSchema.Parse(map[string]interface{}{
    "value": 10,
    "left": map[string]interface{}{
        "value": 5,
        "left":  nil,
        "right": nil,
    },
    "right": map[string]interface{}{
        "value": 15,
        "left":  nil,
        "right": nil,
    },
}, ctx)
```

### Linked List

```go
registry := schema.NewSchemaRegistry()

listNodeSchema := schema.Object().
    Property("data", schema.Any().Required()).
    Property("next", schema.OneOf(
        schema.Ref("#/ListNode", registry),
        schema.Null(),
    ))

registry.Define("ListNode", listNodeSchema)

// Parse linked list
result := listNodeSchema.Parse(map[string]interface{}{
    "data": "first",
    "next": map[string]interface{}{
        "data": "second",
        "next": map[string]interface{}{
            "data": "third",
            "next": nil,
        },
    },
}, ctx)
```

### Comment Thread (Nested Comments)

```go
registry := schema.NewSchemaRegistry()

commentSchema := schema.Object().
    Property("id", schema.Int().Required()).
    Property("author", schema.String().Required()).
    Property("text", schema.String().Required()).
    Property("replies", schema.Array(
        schema.Ref("#/Comment", registry),
    ).Optional())

registry.Define("Comment", commentSchema)

result := commentSchema.Parse(map[string]interface{}{
    "id":     1,
    "author": "Alice",
    "text":   "Great article!",
    "replies": []interface{}{
        map[string]interface{}{
            "id":     2,
            "author": "Bob",
            "text":   "I agree!",
            "replies": []interface{}{
                map[string]interface{}{
                    "id":      3,
                    "author":  "Charlie",
                    "text":    "Me too!",
                    "replies": []interface{}{},
                },
            },
        },
    },
}, ctx)
```

### API with Shared Components

```go
registry := schema.NewSchemaRegistry()

// Define common schemas
registry.Define("Error", schema.Object().
    Property("code", schema.String().Required()).
    Property("message", schema.String().Required()))

registry.Define("User", schema.Object().
    Property("id", schema.Int().Required()).
    Property("username", schema.String().Required()))

registry.Define("Post", schema.Object().
    Property("id", schema.Int().Required()).
    Property("title", schema.String().Required()).
    Property("author", schema.Ref("#/User", registry)))

// API response schemas
successResponseSchema := schema.Object().
    Property("success", schema.Bool().Const(true)).
    Property("data", schema.Any())

errorResponseSchema := schema.Object().
    Property("success", schema.Bool().Const(false)).
    Property("error", schema.Ref("#/Error", registry))

apiResponseSchema := schema.OneOf(
    successResponseSchema,
    errorResponseSchema,
)
```

### JSON Schema Definitions

```go
registry := schema.NewSchemaRegistry()

// Define all component schemas
registry.Define("User", userSchema)
registry.Define("Address", addressSchema)
registry.Define("Order", orderSchema)

// Create main schema with definitions
mainSchema := schema.WithDefinitions(
    schema.Object().
        Property("users", schema.Array(schema.Ref("#/User", registry))).
        Property("addresses", schema.Array(schema.Ref("#/Address", registry))),
    registry,
)

// Generate JSON Schema with $defs section
jsonSchema := mainSchema.JSON()
// Outputs:
// {
//   "type": "object",
//   "properties": {
//     "users": {"type": "array", "items": {"$ref": "#/User"}},
//     "addresses": {"type": "array", "items": {"$ref": "#/Address"}}
//   },
//   "$defs": {
//     "User": {...},
//     "Address": {...}
//   }
// }
```

### File System (Directories and Files)

```go
registry := schema.NewSchemaRegistry()

fileSchema := schema.Object().
    Property("type", schema.String().Const("file")).
    Property("name", schema.String().Required()).
    Property("size", schema.Int().Min(0))

directorySchema := schema.Object().
    Property("type", schema.String().Const("directory")).
    Property("name", schema.String().Required()).
    Property("children", schema.Array(
        schema.Ref("#/FileSystemNode", registry),
    ))

fileSystemNodeSchema := schema.OneOf(
    fileSchema,
    directorySchema,
)

registry.Define("FileSystemNode", fileSystemNodeSchema)
```

### Organization Hierarchy

```go
registry := schema.NewSchemaRegistry()

employeeSchema := schema.Object().
    Property("id", schema.Int().Required()).
    Property("name", schema.String().Required()).
    Property("title", schema.String().Required()).
    Property("reports", schema.Array(
        schema.Ref("#/Employee", registry),
    ).Optional())

registry.Define("Employee", employeeSchema)

// CEO with nested reports
orgChartSchema := schema.Ref("#/Employee", registry)

result := orgChartSchema.Parse(map[string]interface{}{
    "id":    1,
    "name":  "Alice",
    "title": "CEO",
    "reports": []interface{}{
        map[string]interface{}{
            "id":    2,
            "name":  "Bob",
            "title": "CTO",
            "reports": []interface{}{
                map[string]interface{}{
                    "id":      3,
                    "name":    "Charlie",
                    "title":   "Developer",
                    "reports": []interface{}{},
                },
            },
        },
    },
}, ctx)
```

### Circular Reference Detection

```go
registry := schema.NewSchemaRegistry()

// Schema that references itself
registry.Define("Node", schema.Object().
    Property("id", schema.Int()).
    Property("next", schema.Ref("#/Node", registry)))

nodeRef := schema.Ref("#/Node", registry)

// The library detects circular references during parsing
// to prevent infinite loops
```

## Error Handling

```go
result := refSchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes:
        // - "invalid_ref_format" - Reference doesn't start with "#/"
        // - "ref_not_found" - Referenced schema not in registry
        // - "circular_ref" - Circular reference detected
        // - Plus any errors from the referenced schema
    }
}
```

## Important Notes

### Reference Format

References must use the format `#/DefinitionName`:

```go
// ✅ Correct
schema.Ref("#/User", registry)

// ❌ Wrong - missing #/
schema.Ref("User", registry)

// ❌ Wrong - wrong prefix
schema.Ref("$/User", registry)
```

### Circular Reference Protection

The library automatically detects and prevents infinite loops:

```go
// This is safe - library detects circular references
registry.Define("A", schema.Object().
    Property("b", schema.Ref("#/B", registry)))
registry.Define("B", schema.Object().
    Property("a", schema.Ref("#/A", registry)))
```

## When to Use

Ref schemas are ideal for:

- **Schema Reuse**: Define once, use many times
- **Recursive Structures**: Trees, linked lists, nested data
- **Component Libraries**: Shared API schemas
- **OpenAPI/Swagger**: Reusable components
- **Complex Domain Models**: Avoiding duplication
- **Graph Structures**: Nodes with references
- **Organizational Hierarchies**: Employee reporting structures

## Internationalization

```go
refSchema := schema.Ref("#/User", registry).
    RefError(i18n.S("user schema reference not found"))
```

## JSON Schema Generation

```go
registry := schema.NewSchemaRegistry()
registry.Define("User", userSchema)

refSchema := schema.Ref("#/User", registry)

jsonSchema := refSchema.JSON()
// Outputs:
// {
//   "$ref": "#/User"
// }
```

## Related

- [Object Schema](object.md) - Primary type for definitions
- [Array Schema](array.md) - For arrays of references
- [Union Schema](union.md) - For polymorphic references
- [Null Schema](null.md) - For nullable references