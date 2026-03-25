package middleware

import (
	"codebase-app/pkg/jwthandler"
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

func AuthWs(next http.Handler) http.Handler {
	unauthorizedResponse := `
	{
		"message": "Unauthorized",
		"success": false
	}`

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token := r.URL.Query().Get("token")

		if token == "" {
			log.Error().Msg("Unauthorized - Token not set")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(unauthorizedResponse))
			return
		}

		claims, err := jwthandler.ParseAphemeralToken(token)
		if err != nil {
			log.Error().Err(err).Msg("Error while parsing token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(unauthorizedResponse))
			return
		}

		ctx := r.Context()
		claimsMap := map[string]any{
			"user_id": claims.UserID,
			"roles":   claims.Roles,
		}

		ctx = context.WithValue(ctx, "claims", claimsMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(ctx context.Context) (claims map[string]any, err error) {
	claims, ok := ctx.Value("claims").(map[string]any)
	if !ok {
		return nil, errors.New("claims not found in context")
	}

	return claims, nil
}
