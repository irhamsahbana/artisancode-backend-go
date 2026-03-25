package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type EmailSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

type EmailPayload struct {
	To      []string
	Subject string
	Body    string
}

var (
	sender *EmailSender
	once   sync.Once
)

func NewEmailSender(host, port, username, password, from string) *EmailSender {
	return &EmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func InitEmailSender(host, port, username, password, from string) {
	once.Do(func() {
		sender = NewEmailSender(host, port, username, password, from)
	})
}

func GetEmailSender() *EmailSender {
	return sender
}

func (s *EmailSender) buildMessage(payload EmailPayload) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\r\n", s.from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(payload.To, ", ")))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", payload.Subject))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(payload.Body)
	return []byte(sb.String())
}

func (s *EmailSender) Send(payload EmailPayload) {
	go func() {
		if err := s.send(payload); err != nil {
			log.Error().Err(err).
				Str("to", strings.Join(payload.To, ", ")).
				Str("subject", payload.Subject).
				Msg("Failed to send email")
		} else {
			log.Info().
				Str("to", strings.Join(payload.To, ", ")).
				Str("subject", payload.Subject).
				Msg("Email sent successfully")
		}
	}()
}

func (s *EmailSender) send(payload EmailPayload) error {
	address := fmt.Sprintf("%s:%s", s.host, s.port)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.host,
	}

	conn, err := tls.Dial("tcp", address, tlsconfig)
	if err != nil {
		return fmt.Errorf("failed to dial TLS connection: %w", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer c.Close()

	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err = c.Mail(s.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, addr := range payload.To {
		if err = c.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", addr, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	message := s.buildMessage(payload)
	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return c.Quit()
}

func SendEmail(payload EmailPayload) {
	if sender == nil {
		log.Warn().Msg("Email sender not initialized, skipping email")
		return
	}
	sender.Send(payload)
}
