# Conditional Schema

The `ConditionalSchema` provides if-then-else validation logic, allowing different validation rules based on a condition.

## Creating a Conditional Schema

```go
import "github.com/nyxstack/schema"

// If type is "premium", then features array must have at least 5 items
// Otherwise, features array can have at most 2 items
conditionalSchema := schema.Conditional(
    schema.Object().Property("type", schema.String().Const("premium")),
).Then(
    schema.Object().Property("features", schema.Array(schema.String()).MinItems(5)),
).Else(
    schema.Object().Property("features", schema.Array(schema.String()).MaxItems(2)),
)
```

## Methods

### Condition

#### `Conditional(ifSchema Parseable) *ConditionalSchema`
Creates a new conditional schema with an if condition.

```go
schema.Conditional(
    schema.Object().Property("type", schema.String().Const("admin")),
)
```

### Then Branch

#### `Then(thenSchema Parseable) *ConditionalSchema`
Sets the schema to validate when the if condition is true.

```go
schema.Conditional(ifSchema).
    Then(schema.Object().Property("permissions", schema.Array(schema.String()).MinItems(1)))
```

### Else Branch

#### `Else(elseSchema Parseable) *ConditionalSchema`
Sets the schema to validate when the if condition is false.

```go
schema.Conditional(ifSchema).
    Then(thenSchema).
    Else(schema.Object().Property("permissions", schema.Array(schema.String()).MaxItems(3)))
```

### Error Messages

#### `ThenError(err ErrorMessage) *ConditionalSchema`
Sets a custom error message for when the then validation fails.

```go
schema.Conditional(ifSchema).
    Then(thenSchema).
    ThenError(i18n.S("premium accounts must have more features"))
```

#### `ElseError(err ErrorMessage) *ConditionalSchema`
Sets a custom error message for when the else validation fails.

```go
schema.Conditional(ifSchema).
    Then(thenSchema).
    Else(elseSchema).
    ElseError(i18n.S("free accounts cannot have too many features"))
```

## Usage Examples

### Account Type Validation

```go
accountSchema := schema.Conditional(
    schema.Object().Property("type", schema.String().Const("premium")),
).Then(
    schema.Object().
        Property("maxProjects", schema.Int().Min(10)).
        Property("storage", schema.Int().Min(100)),
).Else(
    schema.Object().
        Property("maxProjects", schema.Int().Max(3)).
        Property("storage", schema.Int().Max(10)),
)

ctx := schema.DefaultValidationContext()

// Premium account
result := accountSchema.Parse(map[string]interface{}{
    "type": "premium",
    "maxProjects": 50,
    "storage": 500,
}, ctx) // Valid

// Free account
result := accountSchema.Parse(map[string]interface{}{
    "type": "free",
    "maxProjects": 2,
    "storage": 5,
}, ctx) // Valid
```

### Age-Based Validation

```go
userSchema := schema.Conditional(
    schema.Object().Property("age", schema.Int().Min(18)),
).Then(
    schema.Object().
        Property("canVote", schema.Bool().Const(true)),
).Else(
    schema.Object().
        Property("canVote", schema.Bool().Const(false)).
        Property("guardian", schema.String().Required()),
)

// Adult
result := userSchema.Parse(map[string]interface{}{
    "age": 25,
    "canVote": true,
}, ctx)

// Minor
result := userSchema.Parse(map[string]interface{}{
    "age": 15,
    "canVote": false,
    "guardian": "Parent Name",
}, ctx)
```

### Discount Logic

```go
orderSchema := schema.Conditional(
    schema.Object().Property("total", schema.Number().Min(100)),
).Then(
    schema.Object().
        Property("discount", schema.Number().Min(0.1).Max(0.5)),
).Else(
    schema.Object().
        Property("discount", schema.Number().Max(0.1)),
)

// Large order (discount 10-50%)
result := orderSchema.Parse(map[string]interface{}{
    "total": 150.00,
    "discount": 0.20,
}, ctx)

// Small order (discount max 10%)
result := orderSchema.Parse(map[string]interface{}{
    "total": 50.00,
    "discount": 0.05,
}, ctx)
```

### Shipping Method

```go
shippingSchema := schema.Conditional(
    schema.Object().Property("priority", schema.Bool().Const(true)),
).Then(
    schema.Object().
        Property("deliveryDate", schema.String().Date().Required()).
        Property("cost", schema.Number().Min(20)),
).Else(
    schema.Object().
        Property("estimatedDays", schema.Int().Min(5).Max(14)).
        Property("cost", schema.Number().Min(5)),
)

// Priority shipping
result := shippingSchema.Parse(map[string]interface{}{
    "priority": true,
    "deliveryDate": "2025-11-20",
    "cost": 25.00,
}, ctx)
```

### Membership Level Features

```go
membershipSchema := schema.Conditional(
    schema.Object().Property("level", schema.String().Enum([]string{"gold", "platinum"})),
).Then(
    schema.Object().
        Property("perks", schema.Array(schema.String()).MinItems(5)),
).Else(
    schema.Object().
        Property("perks", schema.Array(schema.String()).MaxItems(2)),
)
```

### Nested Conditionals

```go
// Outer conditional: Check if corporate account
outerSchema := schema.Conditional(
    schema.Object().Property("accountType", schema.String().Const("corporate")),
).Then(
    // Inner conditional: Check team size
    schema.Conditional(
        schema.Object().Property("teamSize", schema.Int().Min(50)),
    ).Then(
        schema.Object().Property("dedicatedSupport", schema.Bool().Const(true)),
    ).Else(
        schema.Object().Property("dedicatedSupport", schema.Bool().Const(false)),
    ),
)
```

### With Custom Error Messages

```go
conditionalSchema := schema.Conditional(
    schema.Object().Property("type", schema.String().Const("premium")),
).Then(
    schema.Object().Property("features", schema.Array(schema.String()).MinItems(5)),
).ThenError(i18n.S("premium accounts require at least 5 features")).
Else(
    schema.Object().Property("features", schema.Array(schema.String()).MaxItems(2)),
).ElseError(i18n.S("free accounts cannot have more than 2 features"))
```

## When to Use

Conditional schemas are ideal for:

- **Account Tiers**: Different validation rules for premium vs free accounts
- **Age Restrictions**: Different requirements based on age
- **Role-Based Validation**: Admin vs user permissions
- **Dynamic Pricing**: Validation rules based on order size
- **Feature Flags**: Different schemas when features are enabled/disabled

## Error Handling

```go
result := conditionalSchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes: "then_failed", "else_failed"
    }
}
```

## Internationalization

```go
schema.Conditional(ifSchema).
    Then(thenSchema).
    ThenError(i18n.S("then validation failed")).
    Else(elseSchema).
    ElseError(i18n.S("else validation failed"))
```

## JSON Schema Generation

```go
schema := schema.Conditional(
    schema.Object().Property("type", schema.String().Const("premium")),
).Then(
    schema.Object().Property("limit", schema.Int().Min(100)),
).Else(
    schema.Object().Property("limit", schema.Int().Max(10)),
)

jsonSchema := schema.JSON()
// Outputs:
// {
//   "if": {
//     "properties": {"type": {"const": "premium"}}
//   },
//   "then": {
//     "properties": {"limit": {"type": "integer", "minimum": 100}}
//   },
//   "else": {
//     "properties": {"limit": {"type": "integer", "maximum": 10}}
//   }
// }
```

## Related

- [Union Schema](union.md) - For either/or validation without conditions
- [Object Schema](object.md) - For structured data
- [Transform Schema](transform.md) - For data transformation
