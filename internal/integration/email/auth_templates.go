package email

import (
	"strings"

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
		base = "http://localhost:5000"
	}
	return base + path
}

func buildLogoURL() string {
	base := config.Envs.FrontendURL.ClientBaseURL
	if isLocalhostURL(base) {
		return presenseLogoDataURI
	}

	return buildAssetURL("/favicon.png")
}

func isLocalhostURL(raw string) bool {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	return normalized == "" ||
		strings.Contains(normalized, "localhost") ||
		strings.Contains(normalized, "127.0.0.1")
}

const presenseLogoDataURI = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABmJLR0QA/wD/AP+gvaeTAAAG8klEQVRYhZWXa2xU1xHHf3Pu3V17/VjbBBYb4wePgjDYILWENpWSQF4oioqEoC1SvpBIVaUkTaNUivqJT3ypokZpK7UUhNRIrQr0U5uk4ZG2iUTSEihExCAoNgnUT/zAsOtd43umH+7d+1gTqTnyes6ec/bOf2b+Z2auUDWaX9i1QR153lrdJiJdQF31ma84Cqp6HTEnHTg49eYfL8Y3JZy9uD2TM40/R/UHgJGUg5NJgZHolMal+n+qoP4c1cR3VQULqrbyW0/g19Mt+gr7js5FAF7cnslJw7vAoybjkmrKYmpSiAhIhNFXHiixitqKtKFUTyGQai14FrWKzltAUQSE928363b2HZ0zADlpeAN41KnLkFnSgKlJI8aAMYgxiBN8Ymu+FDDiywBsCFpIzk1kAMrWpgnzOoA0v7BrgzVy3mQck1ncCKGi4GEm5oWK5Rqz3Ku23FatJ+fqWTAGRDzH8fqMOvI8YNxcFsSEqEUkBCOBlWKMv7bAWvExCoSTSuQkPg+84lkAZ95znjMoj4lrMCk3fEDF8lAGIUisxdwbKqLi+vhanO4Rm9VaBB53FTqcdIVwscNh7MR3GSBYVJWMm6Jz8QN05/KsbVnKwOQY71w8Q9krV10wDf4vHOIvdrpAvThynyMLx6K6erau6KU918LArREcK7Rk61mxKM/uvi3M3/PYf+IY/V8M/h9PU4AGN/yqiqhEe6qgglpFsPS1dvHwyh5OXr7AkbMfhmQ7EhDLUeHZzY/wkyd2cvrqZxx4/+1IzZeoF8BNnIolE7Gg4su+ti42L1/Fb08fp1AugfXvdCaVQo1StiU8z+Pw6RO8feEMh5/9EVjlN6f+UmWwJrUDJtrTRCZTVfAsLdk6Hl7Zw1tn/k6hNItYpautjY1r19LVvozujuVs2rCe7q5ODMLo7Un2/u4NtnxtHWvalkdZMlCkGgNRAZBQnshwymOrN3Ky/zzF2VnEwvpVqykUZzl3sZ/LV69x6cpVzp3/lGKxQG9vD0aE4akJPr7Sz0+f+V6QsYP0HMU7nJpwzcaV+/HNGIf2phYu3ryOepbO1lZGxscYHR8PzkS/GRkeY3h4lM4VnaDKwb+9y2i5QDqVItRfMTTGhcgDVZbjKcubl3BtfBhrPVBLU2MDo6O3wFN68x384fsv8/s9L9O7tMMHMTTCouYmFPA8j08Gr9CdbwsV+xwjERaTiE1geUXmarJM3b0DnpJ20hQLs35KtZbXHtnBykV5Vj2wlNe27QxDWCgUyaTTqCqTM7fJ1dcnq2Oo3Efg+g4ILA9ykVZYai2iEuRzz4fs2dCgOMUrFTHuUVG/FEeKNRaKqhDELa/Mpwt3aMrWoZ6lNFuipqYm3Nt/4hj/uTXM1fEh9r93NAxdNltLqVgCVZobGpmeuRMVsDjZgxHlAU9R489F/BwwODLMjm98G6xiUaYnp8kvWczI0Cif3hjku4d+lqiQrW1LmRif8IuNwopl7fzp1InQarVJEFLtgaQXlHKpzI3xMTZ0dINnGbw2QGt+Cfn84gUea23L09qWZ+DSAGqVvlWr+Xx4iLnyXHRjYl1SwgMaoBM/+4JoWJz++slH7H3yGa7evEmxVOT82Qt0rezm65s3USzMgkK2vpbJsUnOffRvrOfRkKnlqW8+xIFjRyLXx2+ZKhKEQXIv7VZxBOOaqq4mgCjCplVreGh9H4ff+zOFYrGyQTqTBlXKsyWfqKo0ZLLs/c4O/nH2DOf6Ly1s2YIrLpkUpjaNU/Ngzz4J6rbEGaoRo0cmbjE3f489W5/gbrHA6MQE6lm8uXt4c/NBwYKNq9ew58mnOfWvj++j3Lc+9IBrkJQbeCB0eayHi3c1wViUy/HUlm/RkW9l4L83mZqZ8dnemGPFsnY+Hx7inQ8/YGJqKiLdlzSvpjaDyWaQ3Eu7Z1AaBI11RETdTHWroJBOpehubaOxvh4UZu7cZeDmDcpzcwnFSRBxT1hMfQ2mNjPjCtxQ0XXq2cALikrQx8j9EZRLZS4PDt4nw2nEdqugdgH5CBKTSbkAX7gWOSGwDqt+1YpZLxUg1RBi/Ih6CIjf8Xj3HIJQn4TiGHAdFDluHDgIeLhOjKUBci+65+En2AtfOMKqGMuiXlTQNHR/5Rw4DVkAz6o95JT++dlY7YM9SxDZDMC8R1DECeto9cdGUi2hldiK6+PhiHMCnFzlrYtfzrx59C0XYLpFX8lNslZcZxtGsOV5P37hqH45ZMFVXRgCkunXMbiNWSTtApycnm16NcmwfbvSTRPmdRX9IeBU3usSihI8iCHRuIwBCtp6cR3EJ50nwq+mZ5te5cCBe/ejOC0/3tkz7znPCTwOdAH1C7R/tXEXuK7Icav20N1fHO2Pb/4PMkoFFvRnWOwAAAAASUVORK5CYII="
