# Auth Email Flow

## Scope

This project separates normal auth sessions from one-time email actions.

Current covered flows:

- email verification after register
- resend verification email
- forgot password
- reset password
- local email template preview generation

## Core Rules

- Never reuse access token or refresh token as email action token.
- Email verification and password reset must use separate one-time tokens.
- Store action tokens in `user_action_tokens`.
- Login must reject users whose `email_verified_at` is still `NULL`.
- Forgot-password responses must stay generic enough to avoid leaking whether an email exists.

## Database

Relevant schema:

- `users.email_verified_at`
- `user_action_tokens`

`user_action_tokens` stores:

- `user_id`
- `purpose`
- `token_hash`
- `expires_at`
- `used_at`

Purposes currently used:

- `email_verification`
- `password_reset`

## HTTP Endpoints

Public auth routes:

- `POST /users/register`
- `POST /users/login`
- `POST /users/refresh-token`
- `POST /users/verify-email`
- `POST /users/resend-verification-email`
- `POST /users/forgot-password`
- `POST /users/reset-password`

## Rate Limiting

Sensitive email-triggering routes are limited by email + client IP:

- resend verification email:
  - cooldown: 1 request per 60 seconds
  - burst: 3 requests per 15 minutes
- forgot password:
  - cooldown: 1 request per 60 seconds
  - burst: 3 requests per 15 minutes

When blocked, handlers return:

- HTTP `429`
- `Retry-After` header
- localized error message through `pkg/errmsg`

## Email Templates

Templates live in:

- `internal/integration/email/templates/verification.html`
- `internal/integration/email/templates/password_reset.html`

Shared render/build logic lives in:

- `internal/integration/email/render_template.go`
- `internal/integration/email/auth_templates.go`

Brand direction:

- product name: `Presense`
- signature: `by artisanco.de`
- palette: deep teal + mint

Template language follows tenant preferred language when available.

Environment config:

- `APP_PRODUCT_NAME`
- `APP_WEBSITE_URL`
- `APP_SUPPORT_EMAIL`

## Preview Command

Generate local HTML previews without sending an email:

```bash
go run ./cmd/bin/main.go email-preview
```

Optional flags:

```bash
go run ./cmd/bin/main.go email-preview \
  -lang=en \
  -tenant-name="Presense Demo Workspace" \
  -user-name="Rizky Pratama" \
  -out=./tmp/email-previews
```

Outputs:

- `verification.html`
- `password_reset.html`

## Frontend Expectations

Frontend auth pages should support:

- check-email state after register
- verify-email result page
- resend verification with countdown tied to `Retry-After`
- forgot-password page
- reset-password page

The web UI currently implements this through dedicated public auth routes instead of a login-only inline notice.
