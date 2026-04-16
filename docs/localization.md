# Localization

## Current State

Backend multi-language is already implemented for Indonesian (`id`) and English (`en`) in the HTTP layer and error pipeline.

Main building blocks:

- `internal/middleware/language.go`
  - Reads `Accept-Language`.
  - Normalizes it with `errmsg.ResolveLanguage(...)`.
  - Stores it in request context using `errmsg.ContextWithLanguage(...)`.
- `pkg/errmsg/i18n.go`
  - Defines supported languages.
  - Contains translation catalog and language helpers.
- `pkg/errmsg/errors.go`
  - Central entry point for localized error mapping.
- `internal/middleware/localize_json_response.go`
  - Localizes JSON response `message` and nested `errors` payloads before sending response.

## Supported Locales

- `id` default
- `en` supported

Normalization is prefix-based:

- values starting with `en` resolve to English
- values starting with `id` resolve to Indonesian
- everything else falls back to Indonesian

That means headers such as `en-US` and `id-ID` are already handled.

## Request Flow

1. Client sends `Accept-Language`.
2. `WithRequestLanguage()` resolves the effective language.
3. Language is stored in request context.
4. Handlers/core/repository return normal errors.
5. `errmsg.Errors(...)` or `errmsg.NewCustomErrors(...)` builds response payload.
6. `LocalizeJSONResponse()` translates `message` and `errors` before response is written.

## How To Add Localized Errors

For business errors:

1. Return `errmsg.NewCustomErrors(code).SetMessage("Exact message key")`.
2. Add the same exact message key to `pkg/errmsg/i18n.go` when bilingual output is required.

For validation errors:

- Use validator tags and `errmsg.Errors(ctx, err, req)`.
- Validation messages are already generated in `pkg/errmsg/err_validator.go` for both languages.

For database and framework errors:

- Keep routing through `errmsg.Errors(...)`.
- Avoid handcrafting translated maps in modules.

## Rules

- Do not read `Accept-Language` directly in feature handlers.
- Do not branch feature logic by language unless behavior truly differs.
- Do not return ad-hoc localized JSON shapes outside the shared response pattern.
- Prefer English message keys in code, then map them in `pkg/errmsg/i18n.go`.
- If a custom message is user-facing and reused, register it in the catalog instead of duplicating string literals.

## Company Language Config

Company config already contains:

- `preferred_language`
- `supported_languages`

Validation for these fields already exists in company core. This config is tenant data, not the transport-level request language mechanism.

## Practical Check

If localization looks broken:

1. verify client sends `Accept-Language`
2. verify route is behind `WithRequestLanguage()`
3. verify response uses `errmsg`
4. verify custom message exists in `pkg/errmsg/i18n.go`
