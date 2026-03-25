package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Locals struct {
	UserId     string
	Role       string
	TenantId   string
	TenantName string
	UserName   string
	IsVerified bool
}

func GetLocals(c *fiber.Ctx) Locals {
	l := Locals{}
	userId, ok := c.Locals("user_id").(string)
	if ok {
		l.UserId = userId
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_id from locals")
	}

	role, ok := c.Locals("role").(string)
	if ok {
		l.Role = role
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get role from locals")
	}

	tenantId, ok := c.Locals("tenant_id").(string)
	if ok {
		l.TenantId = tenantId
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get tenant_id from locals")
	}

	tenantName, ok := c.Locals("tenant_name").(string)
	if ok {
		l.TenantName = tenantName
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get tenant_name from locals")
	}

	userName, ok := c.Locals("user_name").(string)
	if ok {
		l.UserName = userName
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_name from locals")
	}

	return l
}

func (l *Locals) GetUserId() string {
	return l.UserId
}

func (l *Locals) GetRole() string {
	return l.Role
}

func (l *Locals) GetTenantId() string {
	return l.TenantId
}

func (l *Locals) GetTenantName() string {
	return l.TenantName
}

func (l *Locals) GetCompanyId() string   { return l.GetTenantId() }
func (l *Locals) GetCompanyName() string { return l.GetTenantName() }

func (l *Locals) GetUserName() string {
	return l.UserName
}

func (l *Locals) GetIsVerified() bool {
	return l.IsVerified
}
