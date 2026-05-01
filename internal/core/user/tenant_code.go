package core

import (
	"regexp"
	"strings"

	"codebase-app/pkg/errmsg"
)

var (
	tenantCodePattern  = regexp.MustCompile(`^[A-Z0-9]{3,5}$`)
	tenantCodeReserved = map[string]struct{}{
		"ADMIN": {},
		"OWNER": {},
		"LOGIN": {},
		"ROOT":  {},
		"TEST":  {},
		"NULL":  {},
		"API":   {},
		"APP":   {},
		"WWW":   {},
	}
)

const (
	errorCodeTenantCodeInvalid     = "tenant_code_invalid"
	errorCodeTenantCodeReserved    = "tenant_code_reserved"
	errorCodeTenantCodeAlreadyUsed = "tenant_code_already_used"
)

func normalizeAndValidateTenantCode(code string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if !tenantCodePattern.MatchString(normalized) {
		return "", codedError(400, errmsg.MessageTenantCodeIsInvalid, errorCodeTenantCodeInvalid)
	}
	if _, ok := tenantCodeReserved[normalized]; ok {
		return "", codedError(400, errmsg.MessageTenantCodeIsReserved, errorCodeTenantCodeReserved)
	}
	return normalized, nil
}
