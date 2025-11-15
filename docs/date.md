# Date, DateTime, and Time Schemas

The Nyx Schema library provides specialized schemas for validating date, datetime, and time values with support for multiple formats, range constraints, and timezone handling.

## Creating Date/Time Schemas

```go
import "github.com/nyxstack/schema"

// Date only (YYYY-MM-DD)
dateSchema := schema.Date()

// Date and time with timezone (RFC3339)
datetimeSchema := schema.DateTime()

// Time only (HH:MM:SS)
timeSchema := schema.Time()
```

## Supported Formats

| Format | Constant | Example | Description |
|--------|----------|---------|-------------|
| Date | `FormatDate` | `2025-11-17` | YYYY-MM-DD format |
| DateTime | `FormatDateTime` | `2025-11-17T14:30:00Z` | RFC3339 with timezone |
| Time | `FormatTime` | `14:30:00` | HH:MM:SS format |
| ISO 8601 | `FormatISO8601` | `2025-11-17T14:30:00+00:00` | ISO 8601 format |
| RFC 3339 | `FormatRFC3339` | `2025-11-17T14:30:00Z` | RFC 3339 format |
| Unix | `FormatUnix` | `1700234400` | Unix timestamp (as string) |

## Methods

### Core Methods

#### `Date(errorMessage ...interface{}) *DateSchema`
Creates a new date-only schema (YYYY-MM-DD).

```go
birthDate := schema.Date()
birthDate := schema.Date(i18n.S("birth date is required"))
```

#### `DateTime(errorMessage ...interface{}) *DateSchema`
Creates a new datetime schema with RFC3339 format.

```go
createdAt := schema.DateTime()
createdAt := schema.DateTime(i18n.S("timestamp is required"))
```

#### `Time(errorMessage ...interface{}) *DateSchema`
Creates a new time-only schema (HH:MM:SS).

```go
openingTime := schema.Time()
openingTime := schema.Time(i18n.S("opening time is required"))
```

### Format Validation

#### `Format(format DateFormat) *DateSchema`
Sets the date format to validate against.

```go
isoDate := schema.Date().Format(schema.FormatISO8601)
unixTimestamp := schema.Date().Format(schema.FormatUnix)
```

### Range Constraints

#### `MinDate(min time.Time, errorMessage ...interface{}) *DateSchema`
Sets the minimum date/time constraint.

```go
import "time"

futureDate := schema.Date().MinDate(time.Now())
minBirthDate := schema.Date().
    MinDate(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
        i18n.S("birth date must be after 1900"))
```

#### `MaxDate(max time.Time, errorMessage ...interface{}) *DateSchema`
Sets the maximum date/time constraint.

```go
pastDate := schema.Date().MaxDate(time.Now())
maxAge := schema.Date().
    MaxDate(time.Now().AddDate(-18, 0, 0),
        i18n.S("must be at least 18 years old"))
```

#### `DateRange(min, max time.Time) *DateSchema`
Sets both minimum and maximum date constraints.

```go
validRange := schema.Date().DateRange(
    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
    time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
)
```

### Required/Optional/Nullable

#### `Required(errorMessage ...interface{}) *DateSchema`
Marks the schema as required (default behavior).

```go
required := schema.Date().Required()
required := schema.Date().Required(i18n.S("date is required"))
```

#### `Optional() *DateSchema`
Marks the schema as optional.

```go
optional := schema.Date().Optional()
```

#### `Nullable() *DateSchema`
Allows null values.

```go
nullable := schema.Date().Nullable()
```

### Enum and Const

#### `Enum(values []string, errorMessage ...interface{}) *DateSchema`
Restricts values to a specific set of dates.

```go
allowedDates := schema.Date().Enum([]string{
    "2025-11-17",
    "2025-11-18",
    "2025-11-19",
}, i18n.S("must be one of the allowed dates"))
```

#### `Const(value string, errorMessage ...interface{}) *DateSchema`
Restricts to a single exact date.

```go
deadline := schema.Date().Const("2025-12-31", i18n.S("deadline is December 31, 2025"))
```

### Error Customization

#### `TypeError(message string) *DateSchema`
Sets custom error message for type mismatch.

```go
schema.Date().TypeError("value must be a date string")
```

#### `FormatError(message string) *DateSchema`
Sets custom error message for format validation.

```go
schema.Date().FormatError("date must be in YYYY-MM-DD format")
```

## Usage Examples

### Basic Date Validation

```go
birthDateSchema := schema.Date()

ctx := schema.DefaultValidationContext()

// Valid
result := birthDateSchema.Parse("1990-05-15", ctx) // Valid
result := birthDateSchema.Parse("2025-11-17", ctx) // Valid

// Invalid
result = birthDateSchema.Parse("11/17/2025", ctx)   // Invalid format
result = birthDateSchema.Parse("2025-13-01", ctx)   // Invalid month
result = birthDateSchema.Parse("not-a-date", ctx)   // Invalid
```

### DateTime with Timezone

```go
createdAtSchema := schema.DateTime()

// Valid RFC3339 formats
result := createdAtSchema.Parse("2025-11-17T14:30:00Z", ctx)           // UTC
result = createdAtSchema.Parse("2025-11-17T14:30:00+05:30", ctx)       // IST
result = createdAtSchema.Parse("2025-11-17T14:30:00.123Z", ctx)        // With milliseconds

// Invalid
result = createdAtSchema.Parse("2025-11-17 14:30:00", ctx) // Wrong format
```

### Time-Only Validation

```go
openingTimeSchema := schema.Time()

// Valid
result := openingTimeSchema.Parse("09:00:00", ctx)     // Valid
result = openingTimeSchema.Parse("23:59:59", ctx)      // Valid

// Invalid
result = openingTimeSchema.Parse("9:00", ctx)          // Missing seconds
result = openingTimeSchema.Parse("25:00:00", ctx)      // Invalid hour
```

### Age Restriction (18+)

```go
eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)

birthDateSchema := schema.Date().
    MaxDate(eighteenYearsAgo, i18n.S("you must be at least 18 years old"))

userSchema := schema.Object().
    Property("name", schema.String().MinLength(1)).
    Property("birthDate", birthDateSchema)

// Valid: Person over 18
result := userSchema.Parse(map[string]interface{}{
    "name":      "Alice",
    "birthDate": "2000-01-01",
}, ctx)

// Invalid: Person under 18
result = userSchema.Parse(map[string]interface{}{
    "name":      "Bob",
    "birthDate": "2010-01-01",
}, ctx) // Validation error
```

### Event Date Range

```go
eventStart := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)
eventEnd := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)

eventDateSchema := schema.Date().
    DateRange(eventStart, eventEnd)

// Valid dates in November 2025
result := eventDateSchema.Parse("2025-11-15", ctx) // Valid

// Invalid dates
result = eventDateSchema.Parse("2025-10-31", ctx) // Before range
result = eventDateSchema.Parse("2025-12-01", ctx) // After range
```

### Scheduled Task with Future Date

```go
scheduledDateSchema := schema.DateTime().
    MinDate(time.Now(), i18n.S("scheduled time must be in the future"))

taskSchema := schema.Object().
    Property("name", schema.String().MinLength(1)).
    Property("scheduledFor", scheduledDateSchema)
```

### Working Hours

```go
workHoursSchema := schema.Object().
    Property("start", schema.Time()).
    Property("end", schema.Time())

result := workHoursSchema.Parse(map[string]interface{}{
    "start": "09:00:00",
    "end":   "17:00:00",
}, ctx)
```

### Appointment Booking

```go
today := time.Now().Truncate(24 * time.Hour)
thirtyDaysFromNow := today.AddDate(0, 0, 30)

appointmentSchema := schema.Object().
    Property("date", schema.Date().
        DateRange(today, thirtyDaysFromNow)).
    Property("time", schema.Time().
        Enum([]string{"09:00:00", "10:00:00", "11:00:00", "14:00:00", "15:00:00"},
            i18n.S("please select an available time slot")))
```

### Historical Date Range

```go
historicalDateSchema := schema.Date().
    MinDate(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)).
    MaxDate(time.Now())

// For historical records
recordSchema := schema.Object().
    Property("eventDate", historicalDateSchema).
    Property("description", schema.String())
```

### Nullable Optional Date

```go
optionalDateSchema := schema.Date().Optional().Nullable()

profileSchema := schema.Object().
    Property("name", schema.String().Required()).
    Property("lastLoginDate", optionalDateSchema)

// Valid with null
result := profileSchema.Parse(map[string]interface{}{
    "name":          "Alice",
    "lastLoginDate": nil,
}, ctx)

// Valid without the field
result = profileSchema.Parse(map[string]interface{}{
    "name": "Alice",
}, ctx)
```

### Default Date

```go
registrationSchema := schema.Object().
    Property("registeredAt", schema.DateTime().Default(time.Now().Format(time.RFC3339)))
```

## When to Use

Date/time schemas are ideal for:

- **Birth Dates**: Age verification and restrictions
- **Timestamps**: Created/updated timestamps
- **Event Scheduling**: Appointments, bookings, meetings
- **Deadlines**: Task due dates, expiration dates
- **Working Hours**: Business hours, shift times
- **Historical Records**: Past events with date validation
- **Future Planning**: Scheduled tasks, reservations

## Error Handling

```go
result := dateSchema.Parse(data, ctx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Message: %s\n", err.Message)
        fmt.Printf("Code: %s\n", err.Code)
        // Possible codes: "required", "invalid_type", "format", 
        //                 "min_date", "max_date", "enum", "const"
    }
}
```

## Internationalization

```go
birthDateSchema := schema.Date().
    Required(i18n.S("birth date is required")).
    MaxDate(time.Now().AddDate(-18, 0, 0),
        i18n.S("you must be at least 18 years old")).
    FormatError("date must be in YYYY-MM-DD format")
```

## JSON Schema Generation

```go
dateSchema := schema.Date().
    MinDate(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)).
    MaxDate(time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC))

jsonSchema := dateSchema.JSON()
// Outputs:
// {
//   "type": "string",
//   "format": "date"
// }
```

## Related

- [String Schema](string.md) - String schema includes date/time format validators
- [Object Schema](object.md) - For validating objects with date properties
- [Conditional Schema](conditional.md) - For date-dependent validation logic