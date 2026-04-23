package common

import "time"

const (
	MessageSubjectEmailVerification   = "email.verification"
	MessageSubjectEmailForgotPassword = "email.forgot-password"
	MessageSubjectEmailInvitation     = "email.invitation"
	MessageSubjectEmailAll            = "email.>"
	MessageSubjectExportJobRequested  = "export.job.requested"
)

const (
	MessageStreamEmailService       = "email-service"
	MessageConsumerEmailService     = "email-service-consumer"
	MessageStreamExportJobService   = "export-job-service"
	MessageConsumerExportJobService = "export-job-service-consumer"
)

type MessageBusSubscriptionConfig struct {
	StreamName          string
	StreamDescription   string
	Subjects            []string
	MaxBytes            int64
	MaxAge              time.Duration
	ConsumerName        string
	Durable             string
	ConsumerDescription string
}
