package core

import (
	"fmt"
	"time"

	"codebase-app/pkg/errmsg"
)

func localizeExportText(language, englishText, indonesianText string) string {
	if errmsg.ResolveLanguage(language) == errmsg.LanguageEnglish {
		return englishText
	}

	return indonesianText
}

func buildGoogleMapsURL(latitude, longitude *float64) string {
	if latitude == nil || longitude == nil {
		return ""
	}

	return fmt.Sprintf("https://www.google.com/maps?q=%.7f,%.7f", *latitude, *longitude)
}

func formatAttendanceTimestamp(value string) string {
	if value == "" {
		return ""
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}

	return parsed.Format("2006-01-02 15:04:05 MST")
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func intValue(value *int) string {
	if value == nil {
		return ""
	}

	return fmt.Sprintf("%d", *value)
}
