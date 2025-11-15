package schema

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/nyxstack/i18n"
)

// BinaryFormat represents the encoding format for binary data
type BinaryFormat int

const (
	BinaryFormatBase64    BinaryFormat = 0 // Standard base64 encoding
	BinaryFormatBase64URL BinaryFormat = 1 // URL-safe base64 encoding
	BinaryFormatHex       BinaryFormat = 2 // Hexadecimal encoding
)

// Default error messages for binary validation
var (
	binaryRequiredError         = i18n.S("value is required")
	binaryTypeError             = i18n.S("value must be a string")
	binaryInvalidBase64Error    = i18n.S("must be valid base64 encoded data")
	binaryInvalidBase64URLError = i18n.S("must be valid base64url encoded data")
	binaryInvalidHexError       = i18n.S("must be valid hexadecimal encoded data")
	binaryHexLengthError        = i18n.S("hex string must have even length")
)

func binaryInvalidFormatError(format string) i18n.TranslatedFunc {
	return i18n.F("must be valid %s encoded binary data", format)
}

func binarySizeTooSmallError(actual, min int) i18n.TranslatedFunc {
	return i18n.F("binary data size %d bytes is less than minimum %d bytes", actual, min)
}

func binarySizeTooLargeError(actual, max int) i18n.TranslatedFunc {
	return i18n.F("binary data size %d bytes exceeds maximum %d bytes", actual, max)
}

// BinarySchema represents binary data validation schema
type BinarySchema struct {
	Schema
	format      BinaryFormat
	minSize     *int
	maxSize     *int
	formatError ErrorMessage
	sizeError   ErrorMessage
}

// Binary creates a new binary schema with base64 encoding
func Binary() *BinarySchema {
	return &BinarySchema{
		format: BinaryFormatBase64,
	}
}

// Base64 creates a new binary schema with standard base64 encoding
func Base64() *BinarySchema {
	return &BinarySchema{
		format: BinaryFormatBase64,
	}
}

// Base64URL creates a new binary schema with URL-safe base64 encoding
func Base64URL() *BinarySchema {
	return &BinarySchema{
		format: BinaryFormatBase64URL,
	}
}

// Hex creates a new binary schema with hexadecimal encoding
func Hex() *BinarySchema {
	return &BinarySchema{
		format: BinaryFormatHex,
	}
}

// Format sets the binary encoding format
func (s *BinarySchema) Format(format BinaryFormat) *BinarySchema {
	s.format = format
	return s
}

// MinSize sets the minimum size constraint in bytes
func (s *BinarySchema) MinSize(min int) *BinarySchema {
	s.minSize = &min
	return s
}

// MaxSize sets the maximum size constraint in bytes
func (s *BinarySchema) MaxSize(max int) *BinarySchema {
	s.maxSize = &max
	return s
}

// Size sets both minimum and maximum size constraints in bytes
func (s *BinarySchema) Size(min, max int) *BinarySchema {
	s.minSize = &min
	s.maxSize = &max
	return s
}

// FormatError sets custom error message for format validation
func (s *BinarySchema) FormatError(err ErrorMessage) *BinarySchema {
	s.formatError = err
	return s
}

// SizeError sets custom error message for size validation
func (s *BinarySchema) SizeError(err ErrorMessage) *BinarySchema {
	s.sizeError = err
	return s
}

// Required marks the binary data as required (non-empty)
func (s *BinarySchema) Required() *BinarySchema {
	s.Schema.required = true
	return s
}

// Parse validates binary data
func (s *BinarySchema) Parse(value interface{}, ctx *ValidationContext) ParseResult {
	var errors []ValidationError

	// Convert to string
	binaryStr, ok := value.(string)
	if !ok {
		message := binaryTypeError(ctx.Locale)
		errors = append(errors, NewPrimitiveError(value, message, "invalid_type"))
		return ParseResult{Valid: false, Value: value, Errors: errors}
	}

	// Required validation
	if s.Schema.required && binaryStr == "" {
		message := binaryRequiredError(ctx.Locale)
		errors = append(errors, NewPrimitiveError(binaryStr, message, "required"))
		return ParseResult{Valid: false, Value: value, Errors: errors}
	}

	// If empty and not required, return early
	if binaryStr == "" {
		return ParseResult{Valid: true, Value: binaryStr, Errors: nil}
	}

	// Decode and validate format
	decodedData, err := s.validateAndDecode(binaryStr, ctx)
	if err != nil {
		// err is already a localized error message
		errors = append(errors, NewPrimitiveError(binaryStr, err.Error(), "format"))
		return ParseResult{Valid: false, Value: value, Errors: errors}
	}

	// Validate size constraints
	dataSize := len(decodedData)

	if s.minSize != nil && dataSize < *s.minSize {
		message := binarySizeTooSmallError(dataSize, *s.minSize)(ctx.Locale)
		if !isEmptyErrorMessage(s.sizeError) {
			message = resolveErrorMessage(s.sizeError, ctx)
		}
		errors = append(errors, NewPrimitiveError(binaryStr, message, "min_size"))
	}

	if s.maxSize != nil && dataSize > *s.maxSize {
		message := binarySizeTooLargeError(dataSize, *s.maxSize)(ctx.Locale)
		if !isEmptyErrorMessage(s.sizeError) {
			message = resolveErrorMessage(s.sizeError, ctx)
		}
		errors = append(errors, NewPrimitiveError(binaryStr, message, "max_size"))
	}

	// Return result
	if len(errors) > 0 {
		return ParseResult{Valid: false, Value: value, Errors: errors}
	}

	return ParseResult{Valid: true, Value: binaryStr, Errors: nil}
}

// decodeBinary decodes binary data according to the specified format
func (s *BinarySchema) decodeBinary(data string) ([]byte, error) {
	switch s.format {
	case BinaryFormatBase64:
		return base64.StdEncoding.DecodeString(data)
	case BinaryFormatBase64URL:
		return base64.RawURLEncoding.DecodeString(data)
	case BinaryFormatHex:
		// Simple hex validation and decoding
		if len(data)%2 != 0 {
			// Use the same error message as for format validation
			return nil, errors.New(binaryInvalidHexError("en"))
		}
		decoded := make([]byte, len(data)/2)
		for i := 0; i < len(data); i += 2 {
			var b byte
			_, err := fmt.Sscanf(data[i:i+2], "%02x", &b)
			if err != nil {
				return nil, err
			}
			decoded[i/2] = b
		}
		return decoded, nil
	default:
		return base64.StdEncoding.DecodeString(data)
	}
}

// validateAndDecode validates the format and returns decoded data with localized error messages
func (s *BinarySchema) validateAndDecode(data string, ctx *ValidationContext) ([]byte, error) {
	switch s.format {
	case BinaryFormatBase64:
		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			message := binaryInvalidBase64Error(ctx.Locale)
			if !isEmptyErrorMessage(s.formatError) {
				message = resolveErrorMessage(s.formatError, ctx)
			}
			return nil, errors.New(message)
		}
		return decoded, nil
	case BinaryFormatBase64URL:
		decoded, err := base64.RawURLEncoding.DecodeString(data)
		if err != nil {
			message := binaryInvalidBase64URLError(ctx.Locale)
			if !isEmptyErrorMessage(s.formatError) {
				message = resolveErrorMessage(s.formatError, ctx)
			}
			return nil, errors.New(message)
		}
		return decoded, nil
	case BinaryFormatHex:
		// Check hex string length first
		if len(data)%2 != 0 {
			message := binaryHexLengthError(ctx.Locale)
			if !isEmptyErrorMessage(s.formatError) {
				message = resolveErrorMessage(s.formatError, ctx)
			}
			return nil, errors.New(message)
		}
		// Decode hex
		decoded := make([]byte, len(data)/2)
		for i := 0; i < len(data); i += 2 {
			var b byte
			_, err := fmt.Sscanf(data[i:i+2], "%02x", &b)
			if err != nil {
				message := binaryInvalidHexError(ctx.Locale)
				if !isEmptyErrorMessage(s.formatError) {
					message = resolveErrorMessage(s.formatError, ctx)
				}
				return nil, errors.New(message)
			}
			decoded[i/2] = b
		}
		return decoded, nil
	default:
		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			message := binaryInvalidBase64Error(ctx.Locale)
			if !isEmptyErrorMessage(s.formatError) {
				message = resolveErrorMessage(s.formatError, ctx)
			}
			return nil, errors.New(message)
		}
		return decoded, nil
	}
}

// getFormatErrorMessage returns the appropriate error message for format validation
func (s *BinarySchema) getFormatErrorMessage(ctx *ValidationContext) string {
	if !isEmptyErrorMessage(s.formatError) {
		return resolveErrorMessage(s.formatError, ctx)
	}

	switch s.format {
	case BinaryFormatBase64:
		return binaryInvalidBase64Error(ctx.Locale)
	case BinaryFormatBase64URL:
		return binaryInvalidBase64URLError(ctx.Locale)
	case BinaryFormatHex:
		return binaryInvalidHexError(ctx.Locale)
	default:
		return binaryInvalidBase64Error(ctx.Locale)
	}
}

// getFormatName returns the format name for JSON Schema
func (s *BinarySchema) getFormatName() string {
	switch s.format {
	case BinaryFormatBase64:
		return "base64"
	case BinaryFormatBase64URL:
		return "base64url"
	case BinaryFormatHex:
		return "hex"
	default:
		return "base64"
	}
}

// JSON generates JSON Schema for binary validation
func (s *BinarySchema) JSON() map[string]interface{} {
	schema := map[string]interface{}{
		"type": "string",
	}

	// Add format
	format := s.getFormatName()
	if format == "base64" {
		schema["contentEncoding"] = "base64"
	} else {
		schema["format"] = format
	}

	// Add size constraints (these apply to the decoded binary data)
	if s.minSize != nil {
		schema["minLength"] = *s.minSize
	}
	if s.maxSize != nil {
		schema["maxLength"] = *s.maxSize
	}

	return schema
}
