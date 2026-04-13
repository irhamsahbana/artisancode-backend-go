package adapter

import (
	"context"
	"time"

	infra "codebase-app/internal/infrastructure/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func WithStorage() Option {
	return func(a *Adapter) {
		if infra.Envs.Storage.Bucket == "" {
			log.Fatal().Msg("storage bucket is required")
		}

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

		verifyCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = s3Client.HeadBucket(verifyCtx, &s3.HeadBucketInput{
			Bucket: aws.String(infra.Envs.Storage.Bucket),
		})
		if err != nil {
			log.Fatal().Err(err).Str("bucket", infra.Envs.Storage.Bucket).Msg("failed to verify S3 storage connection")
		}

		a.Storage = s3Client

		log.Info().Msg("S3 storage connected")
	}
}
