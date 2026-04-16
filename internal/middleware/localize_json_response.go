package middleware

import (
	"encoding/json"
	"strings"

	"codebase-app/pkg/errmsg"

	"github.com/gofiber/fiber/v2"
)

func LocalizeJSONResponse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}

		contentType := strings.ToLower(string(c.Response().Header.ContentType()))
		if !strings.Contains(contentType, fiber.MIMEApplicationJSON) {
			return nil
		}

		body := c.Response().Body()
		if len(body) == 0 {
			return nil
		}

		var payload map[string]any
		err = json.Unmarshal(body, &payload)
		if err != nil {
			return nil
		}

		lang := errmsg.LanguageFromContext(c.UserContext())

		if message, ok := payload["message"].(string); ok {
			payload["message"] = errmsg.TranslateText(lang, message)
		}

		if errorsValue, ok := payload["errors"]; ok {
			payload["errors"] = localizeErrors(lang, errorsValue)
		}

		localizedBody, err := json.Marshal(payload)
		if err != nil {
			return nil
		}

		c.Response().SetBodyRaw(localizedBody)

		return nil
	}
}

func localizeErrors(lang errmsg.Language, value any) any {
	switch typed := value.(type) {
	case string:
		return errmsg.TranslateText(lang, typed)
	case []any:
		for i, item := range typed {
			typed[i] = localizeErrors(lang, item)
		}
		return typed
	case map[string]any:
		for key, item := range typed {
			typed[key] = localizeErrors(lang, item)
		}
		return typed
	default:
		return value
	}
}
