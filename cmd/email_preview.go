package cmd

import (
	"flag"
	"os"
	"path/filepath"

	emailint "codebase-app/internal/integration/email"

	"github.com/rs/zerolog/log"
)

func RunEmailPreview(cmd *flag.FlagSet, args []string) {
	outputDir := cmd.String("out", "./tmp/email-previews", "directory for rendered email preview files")
	language := cmd.String("lang", "id", "preview language: id or en")
	tenantName := cmd.String("tenant-name", "Presense Demo Workspace", "tenant name used in preview")
	userName := cmd.String("user-name", "Rizky Pratama", "user name used in preview")

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing email-preview flags")
	}

	if err := os.MkdirAll(*outputDir, 0o755); err != nil {
		log.Fatal().Err(err).Str("output_dir", *outputDir).Msg("Failed to create email preview output directory")
	}

	verificationPreview, err := emailint.BuildVerificationEmail(emailint.AuthTemplateInput{
		UserName:          *userName,
		TenantName:        *tenantName,
		ActionLink:        "https://artisanco.de/preview/email-verification",
		PreferredLanguage: *language,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build verification email preview")
	}

	resetPreview, err := emailint.BuildPasswordResetEmail(emailint.AuthTemplateInput{
		UserName:          *userName,
		TenantName:        *tenantName,
		ActionLink:        "https://artisanco.de/preview/reset-password",
		PreferredLanguage: *language,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build password reset email preview")
	}

	verificationPath := filepath.Join(*outputDir, "verification.html")
	resetPath := filepath.Join(*outputDir, "password_reset.html")

	if err := os.WriteFile(verificationPath, []byte(verificationPreview.Body), 0o644); err != nil {
		log.Fatal().Err(err).Str("path", verificationPath).Msg("Failed to write verification email preview")
	}

	if err := os.WriteFile(resetPath, []byte(resetPreview.Body), 0o644); err != nil {
		log.Fatal().Err(err).Str("path", resetPath).Msg("Failed to write password reset email preview")
	}

	log.Info().
		Str("verification_preview", verificationPath).
		Str("password_reset_preview", resetPath).
		Str("language", *language).
		Msg("Email previews generated successfully")
}
