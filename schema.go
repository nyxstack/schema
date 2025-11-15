package schema

// Schema represents the base fields for all JSON Schema types
type Schema struct {
	// Core JSON Schema fields (private - use getters to access)
	schemaType   string        // JSON Schema type
	title        string        // Schema title
	description  string        // Schema description
	defaultValue interface{}   // Default value
	examples     []interface{} // Example values

	// Schema composition
	ref         string             // $ref
	id          string             // $id
	schema      string             // $schema
	definitions map[string]*Schema // definitions

	// Validation - common to all types
	enum     []interface{} // enum values
	constVal interface{}   // const value

	// Required flag (internal for builder logic)
	required bool // Not serialized, used for validation
}

// Base getters for all schema types

// GetType returns the JSON Schema type
func (s *Schema) GetType() string {
	return s.schemaType
}

// GetTitle returns the schema title
func (s *Schema) GetTitle() string {
	return s.title
}

// GetDescription returns the schema description
func (s *Schema) GetDescription() string {
	return s.description
}

// GetDefault returns the default value
func (s *Schema) GetDefault() interface{} {
	return s.defaultValue
}

// GetExamples returns the example values
func (s *Schema) GetExamples() []interface{} {
	return s.examples
}

// GetRef returns the $ref value
func (s *Schema) GetRef() string {
	return s.ref
}

// GetId returns the $id value
func (s *Schema) GetId() string {
	return s.id
}

// GetSchema returns the $schema value
func (s *Schema) GetSchema() string {
	return s.schema
}

// GetDefinitions returns the definitions
func (s *Schema) GetDefinitions() map[string]*Schema {
	return s.definitions
}

// GetEnum returns the enum values
func (s *Schema) GetEnum() []interface{} {
	return s.enum
}

// GetConst returns the const value
func (s *Schema) GetConst() interface{} {
	return s.constVal
}

// IsRequired returns whether the schema is required
func (s *Schema) IsRequired() bool {
	return s.required
}
