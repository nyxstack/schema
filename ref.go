package schema

import (
	"strings"

	"github.com/nyxstack/i18n"
)

// Default error message functions for reference validation
func refNotFoundError(ref string) i18n.TranslatedFunc {
	return i18n.F("schema reference '%s' not found", ref)
}

func refCircularError(ref string) i18n.TranslatedFunc {
	return i18n.F("circular reference detected: '%s'", ref)
}

var refInvalidFormatError = i18n.S("invalid reference format - must start with '#/'")

// RefErrors defines error message functions
var RefErrors = struct {
	NotFound      func(string) i18n.TranslatedFunc
	CircularRef   func(string) i18n.TranslatedFunc
	InvalidFormat i18n.TranslatedFunc
}{
	NotFound:      refNotFoundError,
	CircularRef:   refCircularError,
	InvalidFormat: refInvalidFormatError,
}

// SchemaRegistry manages schema definitions for references
type SchemaRegistry struct {
	definitions map[string]Parseable
	resolving   map[string]bool // Track schemas currently being resolved to detect circular refs
}

// NewSchemaRegistry creates a new schema registry
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		definitions: make(map[string]Parseable),
		resolving:   make(map[string]bool),
	}
}

// Define adds a schema definition to the registry
func (r *SchemaRegistry) Define(name string, schema Parseable) {
	r.definitions[name] = schema
}

// Get retrieves a schema definition by name
func (r *SchemaRegistry) Get(name string) (Parseable, bool) {
	schema, exists := r.definitions[name]
	return schema, exists
}

// Clear removes all definitions
func (r *SchemaRegistry) Clear() {
	r.definitions = make(map[string]Parseable)
	r.resolving = make(map[string]bool)
}

// RefSchema represents a JSON Schema reference ($ref)
type RefSchema struct {
	ref      string
	registry *SchemaRegistry
	refError ErrorMessage
}

// Ref creates a new reference schema that points to a definition in the registry
func Ref(ref string, registry *SchemaRegistry) *RefSchema {
	return &RefSchema{
		ref:      ref,
		registry: registry,
	}
}

// RefError sets a custom error message for reference resolution failures
func (s *RefSchema) RefError(err ErrorMessage) *RefSchema {
	s.refError = err
	return s
}

// Parse resolves the reference and validates using the referenced schema
func (s *RefSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	// Validate reference format
	if !strings.HasPrefix(s.ref, "#/") {
		message := RefErrors.InvalidFormat(ctx.Locale)
		if !isEmptyErrorMessage(s.refError) {
			message = resolveErrorMessage(s.refError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  value,
			Errors: []ValidationError{NewPrimitiveError(value, message, "invalid_ref_format")},
		}
	}

	// Extract definition name (remove "#/" prefix)
	defName := s.ref[2:]

	// Check for circular reference
	if s.registry.resolving[s.ref] {
		message := RefErrors.CircularRef(s.ref)(ctx.Locale)
		if !isEmptyErrorMessage(s.refError) {
			message = resolveErrorMessage(s.refError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  value,
			Errors: []ValidationError{NewPrimitiveError(value, message, "circular_ref")},
		}
	}

	// Look up the referenced schema
	referencedSchema, exists := s.registry.Get(defName)
	if !exists {
		message := RefErrors.NotFound(s.ref)(ctx.Locale)
		if !isEmptyErrorMessage(s.refError) {
			message = resolveErrorMessage(s.refError, ctx)
		}
		return ParseResult{
			Valid:  false,
			Value:  value,
			Errors: []ValidationError{NewPrimitiveError(value, message, "ref_not_found")},
		}
	}

	// Mark this reference as currently being resolved
	s.registry.resolving[s.ref] = true
	defer func() {
		delete(s.registry.resolving, s.ref)
	}()

	// Validate using the referenced schema
	return referencedSchema.Parse(value, ctx)
}

// JSON generates JSON Schema for reference
func (s *RefSchema) JSON() map[string]interface{} {
	return map[string]interface{}{
		"$ref": s.ref,
	}
}

// CreateDefinitionSchema creates a schema that includes definitions for use with Ref
type DefinitionSchema struct {
	schema      Parseable
	registry    *SchemaRegistry
	definitions map[string]Parseable
}

// WithDefinitions creates a schema wrapper that includes schema definitions
func WithDefinitions(schema Parseable, registry *SchemaRegistry) *DefinitionSchema {
	return &DefinitionSchema{
		schema:      schema,
		registry:    registry,
		definitions: registry.definitions,
	}
}

// Parse validates using the main schema (definitions are just metadata)
func (s *DefinitionSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	return s.schema.Parse(value, ctx)
}

// JSON generates JSON Schema with definitions
func (s *DefinitionSchema) JSON() map[string]interface{} {
	schema := map[string]interface{}{}

	// Add the main schema
	if mainSchema, ok := s.schema.(interface{ JSON() map[string]interface{} }); ok {
		for k, v := range mainSchema.JSON() {
			schema[k] = v
		}
	}

	// Add definitions section
	if len(s.definitions) > 0 {
		definitions := make(map[string]interface{})
		for name, defSchema := range s.definitions {
			if jsonSchema, ok := defSchema.(interface{ JSON() map[string]interface{} }); ok {
				definitions[name] = jsonSchema.JSON()
			} else {
				definitions[name] = map[string]interface{}{"type": "unknown"}
			}
		}
		schema["$defs"] = definitions
	}

	return schema
}
