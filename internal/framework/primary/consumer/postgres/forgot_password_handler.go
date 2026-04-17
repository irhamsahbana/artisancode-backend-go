package consumer

import (
	"codebase-app/internal/infrastructure/config"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type ForgotPasswordEventPayload struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

func ForgotPasswordHandler(ctx context.Context, msg integrationPorts.MessageBusMessage) {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.ForgotPasswordHandler")
	defer span.End()

	payload := &ForgotPasswordEventPayload{}

	err := json.Unmarshal(msg.Data(), payload)
	if err != nil {
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("consumer::ForgotPasswordHandler Error while unmarshalling payload")
		if err := msg.Ack(); err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(ctx).Error().Err(err).Msg("consumer::ForgotPasswordHandler Error while rejecting message")
		}
		return
	}

	err = sendForgotPassword(ctx, []string{payload.Email}, payload)
	if err != nil {
		sendErr := err
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Any("payload", payload).Msg("consumer::ForgotPasswordHandler Error while sending email")
		if nakErr := msg.Nak(sendErr.Error()); nakErr != nil {
			infraTracing.RecordError(span, nakErr)
			log.Ctx(ctx).Error().Err(nakErr).Any("payload", payload).Msg("consumer::ForgotPasswordHandler Error while negative-acknowledging message")
		}
		return
	}

	if err := msg.Ack(); err != nil {
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Any("payload", payload).Msg("consumer::ForgotPasswordHandler Error while acknowledging message")
	}
}

func sendForgotPassword(ctx context.Context, to []string, payload *ForgotPasswordEventPayload) error {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.SendForgotPassword")
	defer span.End()

	var (
		mail             = config.Envs.Mail
		feURL            = config.Envs.FrontendURL
		mailSMTPHost     = mail.Host
		mailSMTPPort     = mail.Port
		mailSMTPUsername = mail.Username
		mailSMTPPassword = mail.Password
		baseURL          string
	)

	if payload.Role == "client" {
		baseURL = feURL.ClientBaseURL
	} else {
		baseURL = feURL.AdminBaseURL
	}
	body := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Password Reset</title>
		</head>
		<body>
			<p>Click the link below to reset your password</p>
			<a href="` + baseURL + feURL.PasswordReset + `?` + payload.Token + `">Reset Password</a>
		</body>
		</html>
	`

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "Crowners <"+mailSMTPUsername+">")
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", "Password Reset")
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(mailSMTPHost, mailSMTPPort, mailSMTPUsername, mailSMTPPassword)
	err := dialer.DialAndSend(mailer)
	if err != nil {
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("consumer::sendForgotPassword Error while sending email")
		return err
	}

	return nil
}
