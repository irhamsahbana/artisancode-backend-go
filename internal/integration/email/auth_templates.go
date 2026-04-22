package email

import (
	"codebase-app/internal/infrastructure/config"
)

type AuthTemplateInput struct {
	UserName          string
	TenantName        string
	ActionLink        string
	PreferredLanguage string
}

type RenderedEmailTemplate struct {
	Subject string
	Body    string
}

func BuildVerificationEmail(input AuthTemplateInput) (*RenderedEmailTemplate, error) {
	data := buildBaseTemplateData(input)

	if input.PreferredLanguage == "en" {
		data.Eyebrow = "Account verification"
		data.Title = "Activate your " + input.TenantName + " workspace"
		data.Greeting = "Hello " + input.UserName + ","
		data.Body = "Thanks for getting started with Presense. Confirm your email address to finish activating your account and continue with your sign in."
		data.ActionLabel = "Verify Email"
		data.ExpiryText = "7 days"
		data.ExpiryLabel = "This verification link will expire in 7 days."
		data.HelpText = "If you did not create this account, you can safely ignore this message."
		data.SupportLabel = "Need help?"
		data.FooterNote = "Presense is a product by artisanco.de built for calm, reliable attendance workflows."
		body, err := RenderTemplate("verification.html", data)
		if err != nil {
			return nil, err
		}
		return &RenderedEmailTemplate{
			Subject: "Verify your email - " + input.TenantName,
			Body:    body,
		}, nil
	}

	data.Eyebrow = "Verifikasi akun"
	data.Title = "Aktifkan workspace " + input.TenantName + " Anda"
	data.Greeting = "Halo " + input.UserName + ","
	data.Body = "Terima kasih sudah memulai dengan Presense. Konfirmasi alamat email Anda untuk menyelesaikan aktivasi akun dan melanjutkan proses masuk."
	data.ActionLabel = "Verifikasi Email"
	data.ExpiryText = "7 hari"
	data.ExpiryLabel = "Tautan verifikasi ini akan kedaluwarsa dalam 7 hari."
	data.HelpText = "Jika Anda tidak membuat akun ini, Anda bisa mengabaikan email ini dengan aman."
	data.SupportLabel = "Butuh bantuan?"
	data.FooterNote = "Presense adalah produk dari artisanco.de yang dirancang untuk alur absensi yang tenang dan andal."

	body, err := RenderTemplate("verification.html", data)
	if err != nil {
		return nil, err
	}

	return &RenderedEmailTemplate{
		Subject: "Verifikasi email Anda - " + input.TenantName,
		Body:    body,
	}, nil
}

func BuildPasswordResetEmail(input AuthTemplateInput) (*RenderedEmailTemplate, error) {
	data := buildBaseTemplateData(input)

	if input.PreferredLanguage == "en" {
		data.Eyebrow = "Password reset"
		data.Title = "Reset your password"
		data.Greeting = "Hello " + input.UserName + ","
		data.Body = "We received a request to reset your Presense password. Use the button below to continue securely."
		data.ActionLabel = "Reset Password"
		data.ExpiryText = "30 minutes"
		data.ExpiryLabel = "This password reset link will expire in 30 minutes."
		data.IgnoreText = "If you did not request this change, you can safely ignore this email."
		data.HelpText = "If you keep seeing this unexpectedly, please contact your administrator or support team."
		data.SupportLabel = "Need help?"
		data.FooterNote = "Presense is a product by artisanco.de built for calm, reliable attendance workflows."
		body, err := RenderTemplate("password_reset.html", data)
		if err != nil {
			return nil, err
		}
		return &RenderedEmailTemplate{
			Subject: "Reset your password - " + input.TenantName,
			Body:    body,
		}, nil
	}

	data.Eyebrow = "Reset password"
	data.Title = "Reset password Anda"
	data.Greeting = "Halo " + input.UserName + ","
	data.Body = "Kami menerima permintaan untuk mereset password Presense Anda. Gunakan tombol di bawah ini untuk melanjutkan dengan aman."
	data.ActionLabel = "Reset Password"
	data.ExpiryText = "30 menit"
	data.ExpiryLabel = "Tautan reset password ini akan kedaluwarsa dalam 30 menit."
	data.IgnoreText = "Jika Anda tidak meminta perubahan ini, Anda bisa mengabaikan email ini dengan aman."
	data.HelpText = "Jika hal ini terus muncul tanpa Anda minta, silakan hubungi administrator atau tim support Anda."
	data.SupportLabel = "Butuh bantuan?"
	data.FooterNote = "Presense adalah produk dari artisanco.de yang dirancang untuk alur absensi yang tenang dan andal."

	body, err := RenderTemplate("password_reset.html", data)
	if err != nil {
		return nil, err
	}

	return &RenderedEmailTemplate{
		Subject: "Reset password Anda - " + input.TenantName,
		Body:    body,
	}, nil
}

func BuildInvitationEmail(input AuthTemplateInput) (*RenderedEmailTemplate, error) {
	data := buildBaseTemplateData(input)

	if input.PreferredLanguage == "en" {
		data.Eyebrow = "Workspace access"
		data.Title = "Your access to " + input.TenantName + " is ready"
		data.Greeting = "Hello " + input.UserName + ","
		data.Body = "Your administrator has invited you to Presense. Open the link below to create your password and start signing in."
		data.ActionLabel = "Create Password"
		data.ExpiryText = "7 days"
		data.ExpiryLabel = "This activation link will expire in 7 days."
		data.HelpText = "If you were not expecting this invite, you can ignore this email and contact your administrator."
		data.SupportLabel = "Need help?"
		data.FooterNote = "Presense is a product by artisanco.de built for calm, reliable attendance workflows."
		body, err := RenderTemplate("invitation.html", data)
		if err != nil {
			return nil, err
		}
		return &RenderedEmailTemplate{
			Subject: "Your access is ready - " + input.TenantName,
			Body:    body,
		}, nil
	}

	data.Eyebrow = "Akses workspace"
	data.Title = "Akses Anda ke " + input.TenantName + " sudah siap"
	data.Greeting = "Halo " + input.UserName + ","
	data.Body = "Administrator Anda telah mengundang Anda ke Presense. Buka tautan di bawah untuk membuat kata sandi dan mulai masuk."
	data.ActionLabel = "Buat Kata Sandi"
	data.ExpiryText = "7 hari"
	data.ExpiryLabel = "Tautan aktivasi ini akan kedaluwarsa dalam 7 hari."
	data.HelpText = "Jika Anda tidak mengharapkan undangan ini, abaikan email ini dan hubungi administrator Anda."
	data.SupportLabel = "Butuh bantuan?"
	data.FooterNote = "Presense adalah produk dari artisanco.de yang dirancang untuk alur absensi yang tenang dan andal."

	body, err := RenderTemplate("invitation.html", data)
	if err != nil {
		return nil, err
	}

	return &RenderedEmailTemplate{
		Subject: "Akses Anda sudah siap - " + input.TenantName,
		Body:    body,
	}, nil
}

func buildBaseTemplateData(input AuthTemplateInput) AuthEmailTemplateData {
	appName := config.Envs.App.ProductName
	if appName == "" {
		appName = "Presense"
	}

	appWebsiteURL := config.Envs.App.WebsiteURL
	if appWebsiteURL == "" {
		appWebsiteURL = "https://artisanco.de"
	}

	supportEmail := config.Envs.App.SupportEmail
	if supportEmail == "" {
		supportEmail = "support@artisanco.de"
	}

	return AuthEmailTemplateData{
		AppName:            appName,
		AppSignaturePrefix: "by",
		AppSignatureLabel:  "artisanco.de",
		AppSignatureURL:    appWebsiteURL,
		AppWebsiteURL:      appWebsiteURL,
		SupportEmail:       supportEmail,
		TenantName:         input.TenantName,
		UserName:           input.UserName,
		ActionLink:         input.ActionLink,
		LogoURL:            buildLogoURL(),
	}
}

func buildAssetURL(path string) string {
	base := config.Envs.FrontendURL.ClientBaseURL
	if base == "" {
		base = "http://localhost:3030"
	}
	return base + path
}

func buildLogoURL() string {
	if config.Envs != nil {
		customLogoURL := config.Envs.App.EmailLogoURL
		if customLogoURL != "" {
			return customLogoURL
		}
	}

	return buildAssetURL("/favicon.png")
}
