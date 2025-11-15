package schema

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nyxstack/i18n"
)

// UUIDVersion represents the version of UUID to validate
type UUIDVersion int

const (
	UUIDVersionAny UUIDVersion = 0 // Accept any valid UUID format
	UUIDVersion1   UUIDVersion = 1 // Time-based UUID
	UUIDVersion2   UUIDVersion = 2 // DCE Security UUID
	UUIDVersion3   UUIDVersion = 3 // Name-based (MD5)
	UUIDVersion4   UUIDVersion = 4 // Random UUID
	UUIDVersion5   UUIDVersion = 5 // Name-based (SHA-1)
	UUIDVersion6   UUIDVersion = 6 // Reordered time-based
	UUIDVersion7   UUIDVersion = 7 // Unix timestamp-based
	UUIDVersion8   UUIDVersion = 8 // Custom/vendor-specific
)

// UUIDFormat represents the expected format of the UUID string
type UUIDFormat int

const (
	UUIDFormatHyphenated UUIDFormat = 0 // Standard format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	UUIDFormatCompact    UUIDFormat = 1 // No hyphens: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	UUIDFormatBraced     UUIDFormat = 2 // With braces: {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	UUIDFormatURN        UUIDFormat = 3 // URN format: urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	UUIDFormatAny        UUIDFormat = 4 // Accept any valid format
)

// Default error message functions for UUID validation
func uuidInvalidFormatError(expected string) i18n.TranslatedFunc {
	return i18n.F("must be a valid UUID in %s format", expected)
}

func uuidInvalidVersionError(version int, actual string) i18n.TranslatedFunc {
	return i18n.F("must be a UUID version %d, got version %s", version, actual)
}

var uuidInvalidVersionAnyError = i18n.S("must be a valid UUID")

func uuidInvalidCaseError(expected string) i18n.TranslatedFunc {
	return i18n.F("UUID must be in %s case", expected)
}

// UUIDs defines error message functions
var UUIDs = struct {
	InvalidFormat  func(string) i18n.TranslatedFunc
	InvalidVersion func(int, string) i18n.TranslatedFunc
	InvalidCase    func(string) i18n.TranslatedFunc
}{
	InvalidFormat: uuidInvalidFormatError,
	InvalidVersion: func(version int, actual string) i18n.TranslatedFunc {
		if version == 0 {
			return uuidInvalidVersionAnyError
		}
		return uuidInvalidVersionError(version, actual)
	},
	InvalidCase: uuidInvalidCaseError,
}

// UUIDSchema represents a UUID validation schema
type UUIDSchema struct {
	version        UUIDVersion
	format         UUIDFormat
	caseSensitive  bool
	forceLowercase bool
	forceUppercase bool
	formatError    ErrorMessage
	versionError   ErrorMessage
	caseError      ErrorMessage
}

// UUID creates a new UUID schema
func UUID() *UUIDSchema {
	return &UUIDSchema{
		version:       UUIDVersionAny,
		format:        UUIDFormatAny,
		caseSensitive: false,
	}
}

// Version specifies the required UUID version
func (s *UUIDSchema) Version(version UUIDVersion) *UUIDSchema {
	s.version = version
	return s
}

// Format specifies the required UUID format
func (s *UUIDSchema) Format(format UUIDFormat) *UUIDSchema {
	s.format = format
	return s
}

// CaseSensitive enables case-sensitive validation
func (s *UUIDSchema) CaseSensitive() *UUIDSchema {
	s.caseSensitive = true
	return s
}

// Lowercase forces UUID to be lowercase
func (s *UUIDSchema) Lowercase() *UUIDSchema {
	s.forceLowercase = true
	s.forceUppercase = false
	return s
}

// Uppercase forces UUID to be uppercase
func (s *UUIDSchema) Uppercase() *UUIDSchema {
	s.forceUppercase = true
	s.forceLowercase = false
	return s
}

// FormatError sets custom error message for format validation
func (s *UUIDSchema) FormatError(err ErrorMessage) *UUIDSchema {
	s.formatError = err
	return s
}

// VersionError sets custom error message for version validation
func (s *UUIDSchema) VersionError(err ErrorMessage) *UUIDSchema {
	s.versionError = err
	return s
}

// CaseError sets custom error message for case validation
func (s *UUIDSchema) CaseError(err ErrorMessage) *UUIDSchema {
	s.caseError = err
	return s
}

// UUID regex patterns for different formats
var uuidPatterns = map[UUIDFormat]*regexp.Regexp{
	UUIDFormatHyphenated: regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
	UUIDFormatCompact:    regexp.MustCompile(`^[0-9a-fA-F]{32}$`),
	UUIDFormatBraced:     regexp.MustCompile(`^\{[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\}$`),
	UUIDFormatURN:        regexp.MustCompile(`^urn:uuid:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
}

// Parse validates a UUID value
func (s *UUIDSchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	// Convert to string
	uuidStr, ok := value.(string)
	if !ok {
		message := UUIDs.InvalidFormat("string")(ctx.Locale)
		if !isEmptyErrorMessage(s.formatError) {
			message = resolveErrorMessage(s.formatError, ctx)
		}
		errors = append(errors, NewPrimitiveError(value, message, "format"))
		return ParseResult{Valid: false, Value: value, Errors: errors}
	}

	// Validate format
	normalizedUUID := s.normalizeUUID(uuidStr)
	if !s.validateFormat(uuidStr) {
		formatName := s.getFormatName()
		message := UUIDs.InvalidFormat(formatName)(ctx.Locale)
		if !isEmptyErrorMessage(s.formatError) {
			message = resolveErrorMessage(s.formatError, ctx)
		}
		errors = append(errors, NewPrimitiveError(uuidStr, message, "format"))
		return ParseResult{Valid: false, Value: value, Errors: errors}
	}

	// Validate version if specified
	if s.version != UUIDVersionAny {
		actualVersion := s.extractVersion(normalizedUUID)
		if actualVersion != int(s.version) {
			message := UUIDs.InvalidVersion(int(s.version), fmt.Sprintf("%d", actualVersion))(ctx.Locale)
			if !isEmptyErrorMessage(s.versionError) {
				message = resolveErrorMessage(s.versionError, ctx)
			}
			errors = append(errors, NewPrimitiveError(uuidStr, message, "version"))
		}
	}

	// Validate case if required
	if s.caseSensitive || s.forceLowercase || s.forceUppercase {
		if !s.validateCase(uuidStr) {
			expected := "mixed"
			if s.forceLowercase {
				expected = "lowercase"
			} else if s.forceUppercase {
				expected = "uppercase"
			}
			message := UUIDs.InvalidCase(expected)(ctx.Locale)
			if !isEmptyErrorMessage(s.caseError) {
				message = resolveErrorMessage(s.caseError, ctx)
			}
			errors = append(errors, NewPrimitiveError(uuidStr, message, "case"))
		}
	}

	// Return result
	if len(errors) > 0 {
		return ParseResult{Valid: false, Value: value, Errors: errors}
	}

	// Transform output based on case requirements
	result := uuidStr
	if s.forceLowercase {
		result = strings.ToLower(uuidStr)
	} else if s.forceUppercase {
		result = strings.ToUpper(uuidStr)
	}

	return ParseResult{Valid: true, Value: result, Errors: nil}
}

// normalizeUUID converts UUID to hyphenated format for internal processing
func (s *UUIDSchema) normalizeUUID(uuid string) string {
	// Remove braces and urn prefix
	normalized := uuid
	if strings.HasPrefix(normalized, "{") && strings.HasSuffix(normalized, "}") {
		normalized = normalized[1 : len(normalized)-1]
	}
	normalized = strings.TrimPrefix(normalized, "urn:uuid:")

	// Add hyphens if compact format
	if len(normalized) == 32 && !strings.Contains(normalized, "-") {
		normalized = fmt.Sprintf("%s-%s-%s-%s-%s",
			normalized[0:8],
			normalized[8:12],
			normalized[12:16],
			normalized[16:20],
			normalized[20:32])
	}

	return normalized
}

// validateFormat checks if the UUID matches the required format
func (s *UUIDSchema) validateFormat(uuid string) bool {
	if s.format == UUIDFormatAny {
		// Check all patterns
		for _, pattern := range uuidPatterns {
			if pattern.MatchString(uuid) {
				return true
			}
		}
		return false
	}

	pattern, exists := uuidPatterns[s.format]
	if !exists {
		return false
	}

	return pattern.MatchString(uuid)
}

// validateCase checks if the UUID matches case requirements
func (s *UUIDSchema) validateCase(uuid string) bool {
	if s.forceLowercase {
		return uuid == strings.ToLower(uuid)
	}
	if s.forceUppercase {
		return uuid == strings.ToUpper(uuid)
	}
	return true // No case requirements
}

// extractVersion extracts the version number from a normalized UUID
func (s *UUIDSchema) extractVersion(normalizedUUID string) int {
	// Version is the first character of the third group
	if len(normalizedUUID) >= 15 {
		versionChar := normalizedUUID[14] // Position of version in xxxxxxxx-xxxx-Vxxx-xxxx-xxxxxxxxxxxx
		switch versionChar {
		case '1':
			return 1
		case '2':
			return 2
		case '3':
			return 3
		case '4':
			return 4
		case '5':
			return 5
		case '6':
			return 6
		case '7':
			return 7
		case '8':
			return 8
		}
	}
	return 0 // Unknown version
}

// getFormatName returns human-readable format name
func (s *UUIDSchema) getFormatName() string {
	switch s.format {
	case UUIDFormatHyphenated:
		return "hyphenated"
	case UUIDFormatCompact:
		return "compact"
	case UUIDFormatBraced:
		return "braced"
	case UUIDFormatURN:
		return "URN"
	case UUIDFormatAny:
		return "valid UUID"
	default:
		return "valid UUID"
	}
}

// JSON generates JSON Schema for UUID validation
func (s *UUIDSchema) JSON() map[string]interface{} {
	schema := map[string]interface{}{
		"type":   "string",
		"format": "uuid",
	}

	// Add pattern if specific format is required
	if s.format != UUIDFormatAny {
		if pattern, exists := uuidPatterns[s.format]; exists {
			schema["pattern"] = pattern.String()
		}
	}

	return schema
}
