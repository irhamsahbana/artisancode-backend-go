package postgres

import (
	"fmt"
	"time"

	"codebase-app/internal/entity/common"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	routerMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/jmoiron/sqlx"
)

func NewConsumerRouter(
	db *sqlx.DB,
	cfg Config,
) (*watermillMessage.Router, func() error, error) {
	if db == nil {
		return nil, nil, fmt.Errorf("postgres db is required")
	}

	cfg = cfg.normalized()
	maxRetries := cfg.MaxAttempts - 1
	if maxRetries < 0 {
		maxRetries = 0
	}

	router, err := watermillMessage.NewRouter(
		watermillMessage.RouterConfig{
			CloseTimeout: cfg.ProcessingTimeout,
		},
		watermill.NopLogger{},
	)
	if err != nil {
		return nil, nil, err
	}

	poisonPublisher, err := newWatermillPublisher(
		watermillSQL.BeginnerFromStdSQL(db.DB),
	)
	if err != nil {
		_ = router.Close()
		return nil, nil, err
	}

	poisonQueueMiddleware, err := routerMiddleware.PoisonQueue(
		poisonPublisher,
		common.MessageTopicDeadLetter,
	)
	if err != nil {
		_ = poisonPublisher.Close()
		_ = router.Close()
		return nil, nil, err
	}

	retryMiddleware := routerMiddleware.Retry{
		MaxRetries:          maxRetries,
		InitialInterval:     cfg.RetryDelay,
		MaxInterval:         cfg.RetryDelay,
		Multiplier:          1,
		MaxElapsedTime:      time.Duration(cfg.MaxAttempts) * cfg.ProcessingTimeout,
		RandomizationFactor: 0,
		ResetContextOnRetry: true,
		Logger:              watermill.NopLogger{},
	}

	router.AddMiddleware(
		poisonQueueMiddleware,
		retryMiddleware.Middleware,
		routerMiddleware.Recoverer,
	)

	return router, poisonPublisher.Close, nil
}

func AddConsumerHandler(
	router *watermillMessage.Router,
	db *sqlx.DB,
	cfg Config,
	def common.MessageBusSubscriptionConfig,
	handlerName string,
	topic string,
	handler watermillMessage.NoPublishHandlerFunc,
) error {
	if router == nil {
		return fmt.Errorf("watermill router is required")
	}
	if db == nil {
		return fmt.Errorf("postgres db is required")
	}

	subscriber, err := newWatermillSubscriber(
		watermillSQL.BeginnerFromStdSQL(db.DB),
		cfg.normalized(),
		def,
	)
	if err != nil {
		return err
	}

	router.AddConsumerHandler(
		handlerName,
		topic,
		subscriber,
		handler,
	)

	return nil
}
