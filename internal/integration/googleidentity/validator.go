package googleidentity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	integrationPorts "codebase-app/internal/ports/integration"

	"google.golang.org/api/idtoken"
)

var _ integrationPorts.GoogleIDTokenValidator = &Validator{}

type Validator struct {
	allowedClientIDs []string
}

func NewValidator(allowedClientIDs []string) *Validator {
	seen := make(map[string]struct{}, len(allowedClientIDs))
	ids := make([]string, 0, len(allowedClientIDs))
	for _, id := range allowedClientIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return &Validator{allowedClientIDs: ids}
}

func NewValidatorFromEnv() *Validator {
	if config.Envs == nil {
		return NewValidator(nil)
	}

	googleCfg := config.Envs.Oauth.Google
	return NewValidator([]string{
		googleCfg.WebClientID,
		googleCfg.IOSClientID,
		googleCfg.AndroidClientID,
		googleCfg.ClientId,
	})
}

func (v *Validator) Validate(ctx context.Context, rawIDToken string, nonce string) (*coreentity.GoogleIdentity, error) {
	if strings.TrimSpace(rawIDToken) == "" {
		return nil, errors.New("google id token is required")
	}
	if len(v.allowedClientIDs) == 0 {
		return nil, errors.New("google client id is not configured")
	}

	var lastErr error
	for _, audience := range v.allowedClientIDs {
		payload, err := idtoken.Validate(ctx, rawIDToken, audience)
		if err != nil {
			lastErr = err
			continue
		}

		identity, err := identityFromPayload(payload)
		if err != nil {
			return nil, err
		}

		if nonce != "" && identity.Nonce != nonce {
			return nil, errors.New("google nonce is invalid")
		}

		return identity, nil
	}

	if lastErr == nil {
		lastErr = errors.New("google id token audience is invalid")
	}
	return nil, lastErr
}

func identityFromPayload(payload *idtoken.Payload) (*coreentity.GoogleIdentity, error) {
	subject := strings.TrimSpace(payload.Subject)
	email := claimString(payload.Claims, "email")
	if subject == "" {
		return nil, errors.New("google subject is required")
	}
	if email == "" {
		return nil, errors.New("google email is required")
	}

	emailVerified, ok := claimBool(payload.Claims, "email_verified")
	if !ok || !emailVerified {
		return nil, errors.New("google email is not verified")
	}

	picture := claimString(payload.Claims, "picture")
	var pictureURL *string
	if picture != "" {
		pictureURL = &picture
	}

	return &coreentity.GoogleIdentity{
		Subject:       subject,
		Email:         strings.ToLower(strings.TrimSpace(email)),
		EmailVerified: emailVerified,
		DisplayName:   claimString(payload.Claims, "name"),
		PictureURL:    pictureURL,
		Nonce:         claimString(payload.Claims, "nonce"),
	}, nil
}

func claimString(claims map[string]any, key string) string {
	value, ok := claims[key]
	if !ok || value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return strings.TrimSpace(str)
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func claimBool(claims map[string]any, key string) (bool, bool) {
	value, ok := claims[key]
	if !ok || value == nil {
		return false, false
	}

	switch v := value.(type) {
	case bool:
		return v, true
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true":
			return true, true
		case "false":
			return false, true
		}
	}

	return false, false
}
