package adapter

import (
	"context"

	infra "codebase-app/internal/infrastructure/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func WithStorage() Option {
	return func(a *Adapter) {
		s3Config, err := config.LoadDefaultConfig(
			context.Background(),
			config.WithRegion(infra.Envs.Storage.Region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					infra.Envs.Storage.Key,
					infra.Envs.Storage.Secret,
					"",
				),
			),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("error while loading s3 config")
		}

		s3Client := s3.NewFromConfig(s3Config, func(o *s3.Options) {
			if infra.Envs.Storage.Endpoint != "" {
				o.BaseEndpoint = aws.String(infra.Envs.Storage.Endpoint)
			}
			if infra.Envs.Storage.IsUsePathStyle {
				o.UsePathStyle = true
			}
		})

		a.Storage = s3Client

		log.Info().Msg("S3 storage connected")
	}
}
