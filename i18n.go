package schema

import (
	"fmt"

	"github.com/nyxstack/i18n"
)

// ErrorMessage represents either a static string or an i18n translatable message
type ErrorMessage interface {
	Resolve(ctx *ValidationContext) string
}

// StaticMessage represents a static error message
type StaticMessage string

func (s StaticMessage) Resolve(ctx *ValidationContext) string {
	return string(s)
}

// I18nMessage represents an i18n translatable message
type I18nMessage i18n.TranslatedFunc

func (i I18nMessage) Resolve(ctx *ValidationContext) string {
	return i18n.TranslatedFunc(i)(ctx.Locale)
}

// Msg creates a static error message
func Msg(message string) ErrorMessage {
	return StaticMessage(message)
}

// I18n creates an i18n error message
func I18n(translatedFunc i18n.TranslatedFunc) ErrorMessage {
	return I18nMessage(translatedFunc)
}

// Helper function to convert string or i18n function to ErrorMessage
func toErrorMessage(input interface{}) ErrorMessage {
	switch v := input.(type) {
	case string:
		if v == "" {
			return nil
		}
		return Msg(v)
	case ErrorMessage:
		return v
	case i18n.TranslatedFunc:
		return I18n(v)
	case nil:
		return nil
	default:
		return Msg(fmt.Sprintf("%v", v))
	}
}

// Helper function to check if ErrorMessage is empty/nil
func isEmptyErrorMessage(em ErrorMessage) bool {
	return em == nil
}

// Helper function to resolve ErrorMessage to string
func resolveErrorMessage(em ErrorMessage, ctx *ValidationContext) string {
	if em == nil {
		return ""
	}
	return em.Resolve(ctx)
}
