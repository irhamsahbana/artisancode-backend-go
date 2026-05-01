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
  - Loads the embedded go-i18n bundle.
- `pkg/errmsg/locales/*.toml`
  - Contains message-code based translations.
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

1. Add the message code to `pkg/errmsg/locales/id.toml` and `pkg/errmsg/locales/en.toml`.
2. Regenerate/add the matching constant in `pkg/errmsg/message_code.go`.
3. Return `errmsg.NewCustomErrors(code).SetMessage(errmsg.MessageSomeCode)`.

For validation errors:

- Use validator tags and `errmsg.Errors(ctx, err, req)`.
- Validation messages use `pkg/errmsg/err_validator.go` plus message templates in `pkg/errmsg/locales/id.toml` and `pkg/errmsg/locales/en.toml`.

For database and framework errors:

- Keep routing through `errmsg.Errors(...)`.
- Avoid handcrafting translated maps in modules.

## Rules

- Do not read `Accept-Language` directly in feature handlers.
- Do not branch feature logic by language unless behavior truly differs.
- Do not return ad-hoc localized JSON shapes outside the shared response pattern.
- Use stable message-code constants in code, not raw strings or user-facing English sentences.
- If a custom message is user-facing and reused, register it in both locale files instead of duplicating string literals.
- `pkg/errmsg/message_code_test.go` checks locale parity and verifies constants exist in every supported locale.

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
4. verify custom message code exists in `pkg/errmsg/locales/id.toml` and `pkg/errmsg/locales/en.toml`
