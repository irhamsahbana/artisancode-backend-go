package consumer

import (
	"codebase-app/internal/infrastructure/config"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type EmailVerificationEventPayload struct {
	Email   string   `json:"email"`
	Token   string   `json:"token"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	IsHTML  bool     `json:"is_html"`
}

func EmailVerificationHandler(ctx context.Context, msg integrationPorts.MessageBusMessage) {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.EmailVerificationHandler")
	defer span.End()

	payload := &EmailVerificationEventPayload{}

	err := json.Unmarshal(msg.Data(), payload)
	if err != nil {
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("consumer::EmailVerificationHandler Error while unmarshalling payload")
		if err := msg.Ack(); err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(ctx).Error().Err(err).Msg("consumer::EmailVerificationHandler Error while rejecting message")
		}
		return
	}

	log.Ctx(ctx).Debug().Any("payload", payload).Msg("consumer::EmailVerificationHandler Received message")

	err = sendEmailVerification(ctx, []string{payload.Email}, payload)
	if err != nil {
		sendErr := err
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Any("payload", payload).Msg("consumer::EmailVerificationHandler Error while sending email")
		if nakErr := msg.Nak(sendErr.Error()); nakErr != nil {
			infraTracing.RecordError(span, nakErr)
			log.Ctx(ctx).Error().Err(nakErr).Any("payload", payload).Msg("consumer::EmailVerificationHandler Error while negative-acknowledging message")
		}
		return
	}
	if err := msg.Ack(); err != nil {
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Any("payload", payload).Msg("consumer::EmailVerificationHandler Error while acknowledging message")
	}
}

func sendEmailVerification(ctx context.Context, to []string, payload *EmailVerificationEventPayload) error {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.SendEmailVerification")
	defer span.End()

	if len(payload.To) > 0 && payload.Subject != "" && payload.Body != "" {
		return sendQueuedEmail(ctx, emailint.EmailPayload{
			To:      payload.To,
			Subject: payload.Subject,
			Body:    payload.Body,
			IsHTML:  payload.IsHTML,
		})
	}

	var (
		mail             = config.Envs.Mail
		feURL            = config.Envs.FrontendURL
		mailSMTPHost     = mail.Host
		mailSMTPPort     = mail.Port
		mailSMTPUsername = mail.Username
		mailSMTPPassword = mail.Password
	)

	body := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Email Verification</title>
		</head>
		<body>
			<p>Click the link below to verify your email</p>
			<a href="` + feURL.ClientBaseURL + feURL.EmailVerification + `?` + payload.Token + `">Verify Email</a>
		</body>
		</html>
	`

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "Crowners <"+mailSMTPUsername+">")
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", "Email Verification")
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(mailSMTPHost, mailSMTPPort, mailSMTPUsername, mailSMTPPassword)
	err := dialer.DialAndSend(mailer)
	if err != nil {
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("consumer::sendEmailVerification Error while sending email")
		return err
	}

	return nil
}
