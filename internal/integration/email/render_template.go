package email

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

type AuthEmailTemplateData struct {
	AppName            string
	AppSignaturePrefix string
	AppSignatureLabel  string
	AppSignatureURL    string
	AppWebsiteURL      string
	SupportEmail       string
	TenantName         string
	UserName           string
	ActionLink         string
	LogoURL            string
	Eyebrow            string
	Title              string
	Greeting           string
	Body               string
	ActionLabel        string
	ExpiryText         string
	ExpiryLabel        string
	IgnoreText         string
	HelpText           string
	SupportLabel       string
	FooterNote         string
}

func RenderTemplate(name string, data AuthEmailTemplateData) (string, error) {
	tpl, err := template.ParseFS(templateFS, "templates/"+name)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
