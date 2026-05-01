# Auth Email Flow

## Scope

This project separates normal auth sessions from one-time email actions.

Current covered flows:

- email verification after register
- Google register/login session creation
- resend verification email
- employee or admin invitation email
- forgot password
- reset password
- local email template preview generation

## Core Rules

- Never reuse access token or refresh token as email action token.
- Email verification and password reset must use separate one-time tokens.
- Store action tokens in `user_action_tokens`.
- Login must reject users whose `email_verified_at` is still `NULL`.
- Google register may set `email_verified_at` immediately when Google `email_verified` is true.
- Google access tokens, refresh tokens, authorization codes, and raw `id_token` values must not be stored or logged.
- Forgot-password responses must stay generic enough to avoid leaking whether an email exists.

## Database

Relevant schema:

- `users.email_verified_at`
- `user_action_tokens`
- `user_auth_identities`

`user_action_tokens` stores:

- `user_id`
- `purpose`
- `token_hash`
- `expires_at`
- `used_at`

Purposes currently used:

- `email_verification`
- `password_reset`

Invitation acceptance currently uses `user_invitations.accept_token_hash` and is delivered through the invitation flow, not `user_action_tokens`.

`user_auth_identities` links external providers to active users. For Google SSO V1:

- provider is `google`
- `(provider, provider_subject)` is globally unique for active identities
- one Google account can be linked to one active Presense user only
- auto-link is allowed only when the Google email matches exactly one active user globally

## HTTP Endpoints

Public auth routes:

- `POST /users/register`
- `POST /users/login`
- `POST /users/refresh-token`
- `POST /users/verify-email`
- `POST /users/resend-verification-email`
- `POST /users/forgot-password`
- `POST /users/reset-password`
- `POST /users/google/register`
- `POST /users/google/login`

Protected tenant route:

- `GET /tenant/profile`

Tenant-aware request payloads:

- `POST /users/login` requires `email`, `password`, and `tenant_code`
- `POST /users/resend-verification-email` requires `email` and `tenant_code`
- `POST /users/forgot-password` requires `email` and `tenant_code`
- `POST /users/google/register` requires `id_token`, `tenant_name`, `tenant_code`, and `confirm_tenant_setup=true`
- `POST /users/google/login` requires `id_token`

Google register requires a user-chosen `tenant_code` before tenant creation. Backend normalizes it to uppercase and enforces 3-5 alphanumeric characters, reserved codes, and active tenant collisions.

## Rate Limiting

Sensitive email-triggering routes are limited by tenant code + email + client IP:

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
- `internal/integration/email/templates/invitation.html`

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
- `APP_EMAIL_LOGO_URL` (optional, overrides default email logo asset URL)
- `APP_SUPPORT_EMAIL`
- `FRONTEND_CLIENT_BASE_URL` (used to build the default logo asset URL when `APP_EMAIL_LOGO_URL` is empty)
- `FRONTEND_INVITATION_URL` (used to build employee/admin invitation links in outgoing email)
- `GOOGLE_WEB_CLIENT_ID`
- `GOOGLE_IOS_CLIENT_ID`
- `GOOGLE_ANDROID_CLIENT_ID`
- `GOOGLE_CLIENT_ID` (fallback compatibility)

## Invitation Delivery

Employee and admin access invitation now uses email as the primary delivery path.

Trigger points:

- `POST /user-invitations`
- `POST /user-invitations/:id/resend`

Behavior:

- handler creates or refreshes invitation data
- core renders the invitation email with tenant language preference
- message is published to `email_invitation`
- consumer sends the HTML email through the configured SMTP transport
- HTTP response includes `email_sent` so UI can hide manual fallback details when delivery is queued successfully

Manual token/link sharing should be treated as a fallback only when `email_sent=false`.

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
- `invitation.html`

## Frontend Expectations

Frontend auth pages should support:

- check-email state after register
- verify-email result page
- resend verification with countdown tied to `Retry-After`
- forgot-password page
- reset-password page

The web UI currently implements this through dedicated public auth routes instead of a login-only inline notice.
