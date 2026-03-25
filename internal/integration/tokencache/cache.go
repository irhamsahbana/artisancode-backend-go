package tokencache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type TokenCacheContract interface {
	SetRefreshToken(token string, data RefreshTokenData, expiration time.Duration)
	GetRefreshToken(token string) (RefreshTokenData, bool)
	DeleteRefreshToken(token string)
	DeleteUserRefreshTokens(userID string)
}

type RefreshTokenData struct {
	UserID   string
	TenantID string
}

type TokenCache struct {
	cache *cache.Cache
}

func NewTokenCache(defaultExpiration, cleanupInterval time.Duration) *TokenCache {
	return &TokenCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (t *TokenCache) SetRefreshToken(token string, data RefreshTokenData, expiration time.Duration) {
	t.cache.Set(token, data, expiration)
}

func (t *TokenCache) GetRefreshToken(token string) (RefreshTokenData, bool) {
	data, found := t.cache.Get(token)
	if !found {
		return RefreshTokenData{}, false
	}
	return data.(RefreshTokenData), true
}

func (t *TokenCache) DeleteRefreshToken(token string) {
	t.cache.Delete(token)
}

func (t *TokenCache) DeleteUserRefreshTokens(userID string) {
	items := t.cache.Items()
	for token, item := range items {
		if data, ok := item.Object.(RefreshTokenData); ok {
			if data.UserID == userID {
				t.cache.Delete(token)
			}
		}
	}
}