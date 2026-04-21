package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"codebase-app/internal/entity/coreentity"
)

const invitationExpiryDuration = 7 * 24 * time.Hour

func generateInvitationToken() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}

	rawToken := hex.EncodeToString(buf)
	hashed := sha256.Sum256([]byte(rawToken))

	return rawToken, hex.EncodeToString(hashed[:]), nil
}

func hashInvitationToken(token string) string {
	hashed := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashed[:])
}

func normalizeInvitationEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func deriveUsernameFromEmail(email string) string {
	email = normalizeInvitationEmail(email)
	parts := strings.SplitN(email, "@", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "user"
	}

	return parts[0]
}

func isExpiredInvitation(item *coreentity.UserInvitation) bool {
	return item.Status == coreentity.UserInvitationStatusPending && item.ExpiresAt.Before(time.Now())
}
