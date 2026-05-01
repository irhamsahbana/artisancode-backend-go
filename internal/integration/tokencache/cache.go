package tokencache

import (
	"context"
	"time"

	"codebase-app/internal/infrastructure/tracing"

	"github.com/patrickmn/go-cache"
)

type TokenCacheContract interface {
	SetRefreshToken(ctx context.Context, token string, data RefreshTokenData, expiration time.Duration)
	GetRefreshToken(ctx context.Context, token string) (RefreshTokenData, bool)
	DeleteRefreshToken(ctx context.Context, token string)
	DeleteUserRefreshTokens(ctx context.Context, userID string)
	SetGoogleRegistration(ctx context.Context, token string, data GoogleRegistrationData, expiration time.Duration)
	GetGoogleRegistration(ctx context.Context, token string) (GoogleRegistrationData, bool)
	DeleteGoogleRegistration(ctx context.Context, token string)
}

type RefreshTokenData struct {
	UserID   string
	TenantID string
}

type GoogleRegistrationData struct {
	Subject       string
	Email         string
	EmailVerified bool
	DisplayName   string
	PictureURL    *string
	Nonce         string
}

type TokenCache struct {
	cache *cache.Cache
}

func NewTokenCache(defaultExpiration, cleanupInterval time.Duration) *TokenCache {
	return &TokenCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (t *TokenCache) SetRefreshToken(
	ctx context.Context,
	token string,
	data RefreshTokenData,
	expiration time.Duration,
) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:tokencache:cache:SetRefreshToken")
	defer span.End()

	t.cache.Set(token, data, expiration)
}

func (t *TokenCache) GetRefreshToken(ctx context.Context, token string) (RefreshTokenData, bool) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:tokencache:cache:GetRefreshToken")
	defer span.End()

	data, found := t.cache.Get(token)
	if !found {
		return RefreshTokenData{}, false
	}
	return data.(RefreshTokenData), true
}

func (t *TokenCache) DeleteRefreshToken(ctx context.Context, token string) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:tokencache:cache:DeleteRefreshToken")
	defer span.End()

	t.cache.Delete(token)
}

func (t *TokenCache) DeleteUserRefreshTokens(ctx context.Context, userID string) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:tokencache:cache:DeleteUserRefreshTokens")
	defer span.End()

	items := t.cache.Items()
	for token, item := range items {
		if data, ok := item.Object.(RefreshTokenData); ok {
			if data.UserID == userID {
				t.cache.Delete(token)
			}
		}
	}
}

func (t *TokenCache) SetGoogleRegistration(
	ctx context.Context,
	token string,
	data GoogleRegistrationData,
	expiration time.Duration,
) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:tokencache:cache:SetGoogleRegistration")
	defer span.End()

	t.cache.Set(token, data, expiration)
}

func (t *TokenCache) GetGoogleRegistration(ctx context.Context, token string) (GoogleRegistrationData, bool) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:tokencache:cache:GetGoogleRegistration")
	defer span.End()

	data, found := t.cache.Get(token)
	if !found {
		return GoogleRegistrationData{}, false
	}
	return data.(GoogleRegistrationData), true
}

func (t *TokenCache) DeleteGoogleRegistration(ctx context.Context, token string) {
	ctx, span := tracing.StartSpan(ctx, "internal:integration:tokencache:cache:DeleteGoogleRegistration")
	defer span.End()

	t.cache.Delete(token)
}
