package errmsg

import (
	"context"
	"embed"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Language string

const (
	LanguageIndonesian Language = "id"
	LanguageEnglish    Language = "en"
)

const DefaultLanguage = LanguageIndonesian

type contextLanguageKey string

const languageContextKey contextLanguageKey = "request_language"

//go:embed locales/*.toml
var localeFiles embed.FS

var (
	bundleOnce sync.Once
	bundle     *goi18n.Bundle
)

func ResolveLanguage(value string) Language {
	normalized := strings.TrimSpace(strings.ToLower(value))
	switch {
	case strings.HasPrefix(normalized, "en"):
		return LanguageEnglish
	case strings.HasPrefix(normalized, "id"):
		return LanguageIndonesian
	default:
		return DefaultLanguage
	}
}

func LanguageFromContext(ctx context.Context) Language {
	if ctx == nil {
		return DefaultLanguage
	}

	if language, ok := ctx.Value(languageContextKey).(string); ok && language != "" {
		return ResolveLanguage(language)
	}

	return DefaultLanguage
}

func ContextWithLanguage(ctx context.Context, language string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, languageContextKey, string(ResolveLanguage(language)))
}

func TranslateText(lang Language, message string) string {
	return TranslateTextWithData(lang, message, nil)
}

func TranslateTextWithData(lang Language, message string, data any) string {
	if message == "" {
		return ""
	}

	if localized, ok := localizeMessage(lang, message, data); ok {
		return localized
	}

	trimmedMessage := strings.TrimSuffix(message, ".")
	if trimmedMessage != message {
		if localized, ok := localizeMessage(lang, trimmedMessage, data); ok {
			return localized
		}
	}

	const timeFormatSuffix = " format must be HH:mm"
	if strings.HasSuffix(message, timeFormatSuffix) {
		fieldName := strings.TrimSuffix(message, timeFormatSuffix)
		if lang == LanguageEnglish {
			return message
		}
		return LocalizeFieldName(lang, fieldName) + " harus berformat HH:mm"
	}

	return message
}

func LocalizeFieldName(lang Language, field string) string {
	if field == "" {
		return ""
	}

	field = strings.ReplaceAll(field, "_", " ")
	if lang == LanguageEnglish {
		return field
	}

	return field
}

func localizeMessage(lang Language, message string, data any) (string, bool) {
	localizer := goi18n.NewLocalizer(messageBundle(), string(lang))
	localized, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    message,
		TemplateData: data,
	})
	if err != nil {
		return "", false
	}

	return localized, true
}

func messageBundle() *goi18n.Bundle {
	bundleOnce.Do(func() {
		b := goi18n.NewBundle(language.Indonesian)
		b.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		mustLoadMessageFileFS(b, "locales/id.toml")
		mustLoadMessageFileFS(b, "locales/en.toml")
		bundle = b
	})

	return bundle
}

func mustLoadMessageFileFS(bundle *goi18n.Bundle, path string) {
	if _, err := bundle.LoadMessageFileFS(localeFiles, path); err != nil {
		panic(err)
	}
}
