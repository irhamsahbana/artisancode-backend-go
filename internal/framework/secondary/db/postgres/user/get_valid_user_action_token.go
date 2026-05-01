package repository

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) GetValidUserActionToken(
	ctx context.Context,
	tokenHash, purpose string,
) (*coreentity.UserActionToken, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:get_valid_user_action_token:GetValidUserActionToken",
	)
	defer span.End()

	query := `
		SELECT
			uat.id,
			uat.user_id,
			uat.purpose,
			uat.token_hash,
			uat.expires_at,
			uat.used_at,
			u.email,
			u.username,
			u.tenant_id,
			t.name AS tenant_name
		FROM user_action_tokens uat
		JOIN users u ON u.id = uat.user_id
		JOIN tenants t ON t.id = u.tenant_id
		WHERE
			uat.token_hash = ?
			AND uat.purpose = ?
			AND uat.deleted_at IS NULL
			AND uat.used_at IS NULL
			AND uat.expires_at > NOW()
			AND u.deleted_at IS NULL
		LIMIT 1
	`

	var row struct {
		ID         string       `db:"id"`
		UserID     string       `db:"user_id"`
		Purpose    string       `db:"purpose"`
		TokenHash  string       `db:"token_hash"`
		ExpiresAt  time.Time    `db:"expires_at"`
		UsedAt     sql.NullTime `db:"used_at"`
		Email      string       `db:"email"`
		UserName   string       `db:"username"`
		TenantID   string       `db:"tenant_id"`
		TenantName string       `db:"tenant_name"`
	}

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &row, exec.Rebind(query), tokenHash, purpose)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
				"purpose": purpose,
			}).Msg(errmsg.MessageUserActionTokenNotFound)
			switch purpose {
			case coreentity.UserActionTokenPurposeEmailVerification:
				return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvalidOrExpiredEmailVerificationToken)
			case coreentity.UserActionTokenPurposePasswordReset:
				return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvalidOrExpiredPasswordResetToken)
			default:
				return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageUserActionTokenNotFound)
			}
		}

		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"purpose": purpose,
		}).Msg("Failed to get user action token")
		return nil, err
	}

	var usedAt *time.Time
	if row.UsedAt.Valid {
		usedAt = &row.UsedAt.Time
	}

	return &coreentity.UserActionToken{
		ID:         row.ID,
		UserID:     row.UserID,
		TenantID:   row.TenantID,
		Email:      row.Email,
		UserName:   row.UserName,
		TenantName: row.TenantName,
		Purpose:    row.Purpose,
		TokenHash:  row.TokenHash,
		ExpiresAt:  row.ExpiresAt,
		UsedAt:     usedAt,
	}, nil
}
