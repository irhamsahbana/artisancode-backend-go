package jwthandler

import (
	"codebase-app/internal/infrastructure/config"
	"context"
	"errors"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func GenerateTokenString(ctx context.Context, payload CostumClaimsPayload) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "pkg:jwthandler:jwt_http:GenerateTokenString")
	defer span.End()

	claims := CustomClaims{
		UserID:      payload.UserID,
		TenantID:    payload.TenantID,
		TenantName:  payload.TenantName,
		UserName:    payload.UserName,
		Roles:       payload.Roles,
		CompanyID:   payload.CompanyID,
		CompanyName: payload.CompanyName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user",
			Issuer:    "codebase-app",
			ExpiresAt: jwt.NewNumericDate(payload.TokenExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString([]byte(config.Envs.Guard.JwtPrivateKey))
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]any{
				"subject":    claims.Subject,
				"user_id":    payload.UserID,
				"tenant_id":  payload.TenantID,
				"company_id": payload.CompanyID,
			}).
			Msg("Failed to sign JWT token")
		return "", err
	}

	return tokenString, nil
}

func ParseTokenString(ctx context.Context, tokenString string) (*CustomClaims, error) {
	ctx, span := tracing.StartSpan(ctx, "pkg:jwthandler:jwt_http:ParseTokenString")
	defer span.End()

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Envs.Guard.JwtPrivateKey), nil
	})
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Warn().
			Err(err).
			Any(common.LogKeyPayload, map[string]any{
				"subject": "user",
			}).
			Msg("Failed to parse JWT token")
		return nil, err
	}

	if !token.Valid {
		err = errors.New("invalid jwt token")
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Warn().
			Any(common.LogKeyPayload, map[string]any{
				"subject": "user",
			}).
			Msg("JWT token is invalid")
		return nil, err
	}

	return claims, nil
}

func GenerateInternalUserTokenString(
	ctx context.Context,
	payload InternalUserClaimsPayload,
) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "pkg:jwthandler:jwt_http:GenerateInternalUserTokenString")
	defer span.End()

	claims := InternalUserClaims{
		UserID:   payload.UserID,
		UserName: payload.UserName,
		Roles:    payload.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "internal_user",
			Issuer:    "codebase-app",
			ExpiresAt: jwt.NewNumericDate(payload.TokenExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString([]byte(config.Envs.Guard.JwtPrivateKey))
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]any{
				"subject": claims.Subject,
				"user_id": payload.UserID,
			}).
			Msg("Failed to sign internal user JWT token")
		return "", err
	}

	return tokenString, nil
}

func ParseInternalUserTokenString(
	ctx context.Context,
	tokenString string,
) (*InternalUserClaims, error) {
	ctx, span := tracing.StartSpan(ctx, "pkg:jwthandler:jwt_http:ParseInternalUserTokenString")
	defer span.End()

	claims := &InternalUserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Envs.Guard.JwtPrivateKey), nil
	})
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Warn().
			Err(err).
			Any(common.LogKeyPayload, map[string]any{
				"subject": "internal_user",
			}).
			Msg("Failed to parse internal user JWT token")
		return nil, err
	}

	if !token.Valid || claims.Subject != "internal_user" {
		err = errors.New("invalid internal user jwt token")
		tracing.RecordError(span, err)
		log.Ctx(ctx).
			Warn().
			Any(common.LogKeyPayload, map[string]any{
				"subject":          claims.Subject,
				"expected_subject": "internal_user",
			}).
			Msg("Internal user JWT token is invalid")
		return nil, err
	}

	return claims, nil
}
