package consumer

import postgresconsumer "codebase-app/internal/framework/primary/consumer/postgres"

type App = postgresconsumer.App

var NewApp = postgresconsumer.NewApp
