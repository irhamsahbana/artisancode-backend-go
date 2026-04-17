package consumer

import postgresconsumer "codebase-app/internal/framework/primary/consumer/postgres"

type EmailVerificationEventPayload = postgresconsumer.EmailVerificationEventPayload

var EmailVerificationHandler = postgresconsumer.EmailVerificationHandler
