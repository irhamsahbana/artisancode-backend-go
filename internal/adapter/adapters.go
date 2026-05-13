package adapter

import (
	"fmt"
	"net/http"
	"strings"

	emailint "codebase-app/internal/integration/email"
	integrationPorts "codebase-app/internal/ports/integration"

	// import "codebase-app/internal/pkg/validator"
	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/openai/openai-go/v3"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

var (
	Adapters *Adapter
)

type Option func(adapter *Adapter)

type Validator interface {
	Validate(i any) error
}

type Adapter struct {
	// Driving Adapters / Primary Adapters
	RestServer *fiber.App
	WsServer   *http.Server

	//Driven Adapters / Secondary Adapters
	Postgres                   *sqlx.DB
	MessagePublisher           integrationPorts.MessagePublisher
	MessageSubscriptionManager integrationPorts.MessageSubscriptionManager
	Validator                  Validator // *validator.Validator
	Storage                    *s3.Client
	EmailSender                *emailint.EmailSender
	VenamonGolog               *tele.Bot
	FirebaseSDK                *firebase.App
	OpenAISDK                  *openai.Client
	DropboxFiles               files.Client
}

func (a *Adapter) Sync(opts ...Option) {
	for o := range opts {
		opt := opts[o]
		opt(a)
	}
}

func (a *Adapter) Unsync() error {
	var errs []string

	if a.RestServer != nil {
		if err := a.RestServer.Shutdown(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Rest server disconnected")
	}

	if a.WsServer != nil {
		if err := a.WsServer.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Ws server disconnected")
	}

	if a.MessagePublisher != nil {
		if err := a.MessagePublisher.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Message publisher disconnected")
	}

	if a.MessageSubscriptionManager != nil {
		if err := a.MessageSubscriptionManager.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Message subscription manager disconnected")
	}

	if a.Postgres != nil {
		if err := a.Postgres.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Postgres disconnected")
	}

	// if a.VenamonGolog != nil {
	// 	a.VenamonGolog.Stop()
	// 	log.Info().Msg("Venamon Golog disconnected")
	// }

	if len(errs) > 0 {
		err := fmt.Errorf("%s", strings.Join(errs, "\n"))
		log.Error().Msgf("Error while disconnecting adapters: %v", err)
		return err
	}

	return nil
}
