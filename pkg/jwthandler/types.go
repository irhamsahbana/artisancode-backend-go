package jwthandler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID      string   `json:"user_id"`
	UserName    string   `json:"user_name"`
	TenantID    string   `json:"tenant_id"`
	TenantName  string   `json:"tenant_name"`
	Roles       []string `json:"roles"`
	CompanyID   *string  `json:"company_id,omitempty"`
	CompanyName *string  `json:"company_name,omitempty"`
	jwt.RegisteredClaims
}

type CostumClaimsWs struct {
	UserID      string   `json:"user_id"`
	Roles       []string `json:"roles"`
	CompanyID   *string  `json:"company_id,omitempty"`
	CompanyName *string  `json:"company_name,omitempty"`
	jwt.RegisteredClaims
}

type CostumClaimsPayload struct {
	UserID          string    `json:"user_id"`
	UserName        string    `json:"user_name"`
	TenantID        string    `json:"tenant_id"`
	TenantName      string    `json:"tenant_name"`
	Roles           []string  `json:"roles"`
	CompanyID       *string   `json:"company_id,omitempty"`
	CompanyName     *string   `json:"company_name,omitempty"`
	TokenExpiration time.Time `json:"token_expiration"`
}

type CostumClaimsPayloadWs struct {
	UserID          string    `json:"user_id"`
	Roles           []string  `json:"roles"`
	CompanyID       *string   `json:"company_id,omitempty"`
	CompanyName     *string   `json:"company_name,omitempty"`
	TokenExpiration time.Time `json:"token_expiration"`
}

type InternalUserClaims struct {
	UserID   string   `json:"user_id"`
	UserName string   `json:"user_name"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

type InternalUserClaimsPayload struct {
	UserID          string
	UserName        string
	Roles           []string
	TokenExpiration time.Time
}
