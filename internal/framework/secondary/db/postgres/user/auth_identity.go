package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) FindAuthIdentityByProviderSubject(
	ctx context.Context,
	provider string,
	providerSubject string,
) (*coreentity.UserAuthIdentity, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:auth_identity:FindAuthIdentityByProviderSubject",
	)
	defer span.End()

	query := `
		SELECT
			id,
			tenant_id,
			user_id,
			provider,
			provider_subject,
			email,
			email_verified,
			display_name,
			picture_url
		FROM user_auth_identities
		WHERE provider = ? AND provider_subject = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var row struct {
		ID              string  `db:"id"`
		TenantID        string  `db:"tenant_id"`
		UserID          string  `db:"user_id"`
		Provider        string  `db:"provider"`
		ProviderSubject string  `db:"provider_subject"`
		Email           string  `db:"email"`
		EmailVerified   bool    `db:"email_verified"`
		DisplayName     string  `db:"display_name"`
		PictureURL      *string `db:"picture_url"`
	}

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &row, exec.Rebind(query), provider, providerSubject)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"provider": provider}).
			Msg("Failed to find auth identity")
		return nil, err
	}

	return &coreentity.UserAuthIdentity{
		ID:              row.ID,
		TenantID:        row.TenantID,
		UserID:          row.UserID,
		Provider:        row.Provider,
		ProviderSubject: row.ProviderSubject,
		Email:           row.Email,
		EmailVerified:   row.EmailVerified,
		DisplayName:     row.DisplayName,
		PictureURL:      row.PictureURL,
	}, nil
}

func (r *userRepo) CreateAuthIdentity(ctx context.Context, identity coreentity.UserAuthIdentity) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:auth_identity:CreateAuthIdentity",
	)
	defer span.End()

	query := `
		INSERT INTO user_auth_identities (
			tenant_id,
			user_id,
			provider,
			provider_subject,
			email,
			email_verified,
			display_name,
			picture_url,
			last_login_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	exec := r.executor(ctx)
	_, err := exec.ExecContext(
		ctx,
		exec.Rebind(query),
		identity.TenantID,
		identity.UserID,
		identity.Provider,
		identity.ProviderSubject,
		identity.Email,
		identity.EmailVerified,
		identity.DisplayName,
		identity.PictureURL,
	)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{
				"tenant_id": identity.TenantID,
				"user_id":   identity.UserID,
				"provider":  identity.Provider,
			}).
			Msg("Failed to create auth identity")
		return err
	}

	return nil
}

func (r *userRepo) UpdateAuthIdentityLastLogin(ctx context.Context, identityID string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:auth_identity:UpdateAuthIdentityLastLogin",
	)
	defer span.End()

	query := `
		UPDATE user_auth_identities
		SET last_login_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	if _, err := exec.ExecContext(ctx, exec.Rebind(query), identityID); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"identity_id": identityID}).
			Msg("Failed to update auth identity last login")
		return err
	}

	return nil
}
