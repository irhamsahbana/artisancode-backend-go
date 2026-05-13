package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"codebase-app/pkg/errmsg"

	"github.com/gofiber/fiber/v3"
)

func TestWithRequestLanguageResolvesAcceptLanguage(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Use(WithRequestLanguage())
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"language": string(errmsg.LanguageFromContext(c.Context())),
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	var payload map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["language"] != "en" {
		t.Fatalf("language = %q, want en", payload["language"])
	}
}

func TestLocalizeJSONResponseLocalizesMessageAndNestedErrors(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Use(WithRequestLanguage())
	app.Use(LocalizeJSONResponse())
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": errmsg.MessageCompanyNotFound,
			"errors": fiber.Map{
				"file": []string{errmsg.MessageFileIsRequired},
			},
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "id-ID")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	var payload struct {
		Message string              `json:"message"`
		Errors  map[string][]string `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Message != "Perusahaan tidak ditemukan" {
		t.Fatalf("message = %q, want localized company message", payload.Message)
	}

	fileErrors := payload.Errors["file"]
	if len(fileErrors) != 1 || fileErrors[0] != "file wajib diisi" {
		t.Fatalf("file errors = %#v, want localized file required error", fileErrors)
	}
}

func TestLocalizeJSONResponseLeavesNonJSONResponse(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Use(WithRequestLanguage())
	app.Use(LocalizeJSONResponse())
	app.Get("/", func(c fiber.Ctx) error {
		c.Type("text")
		return c.SendString(errmsg.MessageCompanyNotFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "id-ID")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}

	if string(body) != errmsg.MessageCompanyNotFound {
		t.Fatalf("body = %q, want original text", string(body))
	}
}
