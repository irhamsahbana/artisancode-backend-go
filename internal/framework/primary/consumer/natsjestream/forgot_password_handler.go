package consumer

import postgresconsumer "codebase-app/internal/framework/primary/consumer/postgres"

type ForgotPasswordEventPayload = postgresconsumer.ForgotPasswordEventPayload

var ForgotPasswordHandler = postgresconsumer.ForgotPasswordHandler
