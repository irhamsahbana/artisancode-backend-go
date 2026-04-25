package coreentity

import (
	"codebase-app/internal/entity/common"
)

type Company struct {
	UserCtx common.UserContext

	ID        string
	TenantID  string
	Name      string
	Code      string
	Config    CompanyConfig
	CreatedAt string
	UpdatedAt string
}

type CompanyConfig struct {
	Logo               *string  `json:"logo"`
	PreferredLanguage  string   `json:"preferred_language"`
	SupportedLanguages []string `json:"supported_languages"`
	Timezone           string   `json:"timezone"`
	DateFormat         string   `json:"date_format"`
	TimeFormat         string   `json:"time_format"`
}

func DefaultCompanyConfig() CompanyConfig {
	return CompanyConfig{
		Logo:               nil,
		PreferredLanguage:  "id",
		SupportedLanguages: []string{"id", "en"},
		Timezone:           "Asia/Makassar",
		DateFormat:         "YYYY-MM-DD",
		TimeFormat:         "HH:mm:ss",
	}
}

type CompanyListFilter struct {
	TenantID string
	Q        string
	Page     int
	Paginate int
}

type CompanyDeleteFilter struct {
	TenantID string
	ID       string
}
