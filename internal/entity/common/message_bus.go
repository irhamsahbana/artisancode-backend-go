package common

const (
	MessageTopicEmailVerification   = "email_verification"
	MessageTopicEmailForgotPassword = "email_forgot_password"
	MessageTopicEmailInvitation     = "email_invitation"
	MessageTopicExportJobRequested  = "export_job_requested"
	MessageTopicDeadLetter          = "dead_letter"
)

const (
	MessageConsumerGroupNotificationService = "notification_service"
	MessageConsumerGroupExportJobService    = "export_job_service"
)

type MessageBusSubscriptionConfig struct {
	Topics              []string
	ConsumerGroup       string
	ConsumerDescription string
}
