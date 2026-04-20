package adapter

import (
	"strconv"

	"codebase-app/internal/infrastructure/config"
	emailint "codebase-app/internal/integration/email"

	"github.com/rs/zerolog/log"
)

func WithEmailSender() Option {
	return func(a *Adapter) {
		mailCfg := config.Envs.Mail
		if mailCfg.Host == "" || mailCfg.Port == 0 || mailCfg.Username == "" || mailCfg.Password == "" || mailCfg.From == "" {
			log.Warn().Any("mail_config", map[string]any{
				"host":         mailCfg.Host,
				"port":         mailCfg.Port,
				"username_set": mailCfg.Username != "",
				"password_set": mailCfg.Password != "",
				"from_set":     mailCfg.From != "",
			}).Msg("SMTP config is incomplete, email sender will not be initialized")
			return
		}

		emailint.InitEmailSender(
			mailCfg.Host,
			strconv.Itoa(mailCfg.Port),
			mailCfg.Username,
			mailCfg.Password,
			mailCfg.From,
		)

		a.EmailSender = emailint.GetEmailSender()
		log.Info().Msg("Email sender initialized")
	}
}
