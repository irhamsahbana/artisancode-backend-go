package adapter

import (
	"codebase-app/internal/infrastructure/config"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/rs/zerolog/log"
)

func WithOpenAISDK() Option {
	return func(a *Adapter) {
		client := openai.NewClient(
			option.WithAPIKey(config.Envs.OpenAI.APIKey),
			option.WithBaseURL(config.Envs.OpenAI.BaseURL),
		)
		log.Info().Msg("OpenAI SDK client created")
		a.OpenAISDK = &client
	}
}
