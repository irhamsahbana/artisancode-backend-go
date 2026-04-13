package core

import "codebase-app/internal/module/z_template_v2/ports"

var _ ports.XxxService = &xxxCore{}

type xxxCore struct {
	repo ports.XxxRepository
}

func NewXxxCore(repo ports.XxxRepository) *xxxCore {
	return &xxxCore{
		repo: repo,
	}
}
