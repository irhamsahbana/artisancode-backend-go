package postgres

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"codebase-app/internal/entity/common"
	secondarypostgres "codebase-app/internal/framework/secondary/db/postgres"
	integrationPorts "codebase-app/internal/ports/integration"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
)

func TestPublishJSONInsertsDefaultPostgresSchemaMessage(t *testing.T) {
	db := openTestDB(t)
	prepareWatermillTables(t, db)
	truncateWatermillTables(t, db)

	publisher := NewPublisher(db)
	err := publisher.PublishJSON(context.Background(), common.MessageTopicEmailVerification, map[string]any{
		"hello": "world",
	})
	if err != nil {
		t.Fatalf("expected no publish error, got %v", err)
	}

	var row struct {
		UUID          string `db:"uuid"`
		Payload       string `db:"payload"`
		Metadata      string `db:"metadata"`
		TransactionID string `db:"transaction_id"`
	}

	query := `
		SELECT uuid,
			payload::text AS payload,
			metadata::text AS metadata,
			transaction_id::text AS transaction_id
		FROM "watermill_email_verification"
		ORDER BY transaction_id DESC, "offset" DESC
		LIMIT 1
	`
	err = db.Get(&row, query)
	if err != nil {
		t.Fatalf("expected queued message row, got %v", err)
	}

	if row.UUID == "" {
		t.Fatal("expected Watermill message uuid to be stored")
	}
	if row.TransactionID == "" {
		t.Fatal("expected Watermill transaction_id to be stored")
	}
	if !strings.Contains(row.Payload, `"hello": "world"`) && !strings.Contains(row.Payload, `"hello":"world"`) {
		t.Fatalf("expected payload json to contain hello world, got %s", row.Payload)
	}
	if !strings.Contains(row.Metadata, `"topic": "email_verification"`) &&
		!strings.Contains(row.Metadata, `"topic":"email_verification"`) {
		t.Fatalf("expected metadata topic to be stored, got %s", row.Metadata)
	}
}

func TestConsumeAndAckStoresConsumerGroupOffset(t *testing.T) {
	db := openTestDB(t)
	prepareWatermillTables(t, db)
	truncateWatermillTables(t, db)

	publisher := NewPublisher(db)
	manager := NewSubscriptionManager(db, Config{
		PollInterval:      10 * time.Millisecond,
		BatchSize:         1,
		RetryDelay:        10 * time.Millisecond,
		ProcessingTimeout: time.Second,
	})

	subscription, err := manager.CreateSubscription(context.Background(), common.MessageBusSubscriptionConfig{
		ConsumerGroup: common.MessageConsumerGroupNotificationService,
		Topics:        []string{common.MessageTopicEmailVerification},
	})
	if err != nil {
		t.Fatalf("expected consumer creation to succeed, got %v", err)
	}

	received := make(chan string, 1)
	consumeCtx, err := subscription.Consume(context.Background(), func(ctx context.Context, msg integrationPorts.MessageBusMessage) {
		received <- string(msg.Data())
		if msg.Subject() != common.MessageTopicEmailVerification {
			t.Errorf("expected topic %s, got %s", common.MessageTopicEmailVerification, msg.Subject())
		}
		if ackErr := msg.Ack(ctx); ackErr != nil {
			t.Errorf("expected ack to succeed, got %v", ackErr)
		}
	})
	if err != nil {
		t.Fatalf("expected consume to start, got %v", err)
	}
	defer consumeCtx.Stop()

	err = publisher.PublishJSON(context.Background(), common.MessageTopicEmailVerification, map[string]any{
		"job": "ack",
	})
	if err != nil {
		t.Fatalf("expected publish to succeed, got %v", err)
	}

	select {
	case payload := <-received:
		if !strings.Contains(payload, `"job": "ack"`) && !strings.Contains(payload, `"job":"ack"`) {
			t.Fatalf("expected ack payload, got %s", payload)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected consumer to receive queued message")
	}

	var messageCount int
	err = db.Get(&messageCount, `SELECT COUNT(*) FROM "watermill_email_verification"`)
	if err != nil {
		t.Fatalf("expected message count query to succeed, got %v", err)
	}
	if messageCount != 1 {
		t.Fatalf("expected default pub/sub message to remain stored, got %d", messageCount)
	}

	var offset struct {
		OffsetAcked                int    `db:"offset_acked"`
		LastProcessedTransactionID string `db:"last_processed_transaction_id"`
	}
	err = db.Get(&offset, `
		SELECT offset_acked,
			last_processed_transaction_id::text AS last_processed_transaction_id
		FROM "watermill_offsets_email_verification"
		WHERE consumer_group = $1
	`, common.MessageConsumerGroupNotificationService)
	if err != nil {
		t.Fatalf("expected consumer group offset row, got %v", err)
	}
	if offset.OffsetAcked <= 0 {
		t.Fatalf("expected positive offset_acked, got %d", offset.OffsetAcked)
	}
	if offset.LastProcessedTransactionID == "" || offset.LastProcessedTransactionID == "0" {
		t.Fatalf("expected last_processed_transaction_id to advance, got %s", offset.LastProcessedTransactionID)
	}
}

func TestSubscriptionConsumesMultipleTopics(t *testing.T) {
	db := openTestDB(t)
	prepareWatermillTables(t, db)
	truncateWatermillTables(t, db)

	publisher := NewPublisher(db)
	manager := NewSubscriptionManager(db, Config{
		PollInterval:      10 * time.Millisecond,
		BatchSize:         1,
		RetryDelay:        10 * time.Millisecond,
		ProcessingTimeout: time.Second,
	})

	subscription, err := manager.CreateSubscription(context.Background(), common.MessageBusSubscriptionConfig{
		ConsumerGroup: common.MessageConsumerGroupNotificationService,
		Topics: []string{
			common.MessageTopicEmailVerification,
			common.MessageTopicEmailInvitation,
		},
	})
	if err != nil {
		t.Fatalf("expected consumer creation to succeed, got %v", err)
	}

	received := make(chan string, 2)
	consumeCtx, err := subscription.Consume(context.Background(), func(ctx context.Context, msg integrationPorts.MessageBusMessage) {
		received <- msg.Subject()
		if ackErr := msg.Ack(ctx); ackErr != nil {
			t.Errorf("expected ack to succeed, got %v", ackErr)
		}
	})
	if err != nil {
		t.Fatalf("expected consume to start, got %v", err)
	}
	defer consumeCtx.Stop()

	if err = publisher.PublishJSON(context.Background(), common.MessageTopicEmailVerification, map[string]any{"job": "verification"}); err != nil {
		t.Fatalf("expected verification publish to succeed, got %v", err)
	}
	if err = publisher.PublishJSON(context.Background(), common.MessageTopicEmailInvitation, map[string]any{"job": "invitation"}); err != nil {
		t.Fatalf("expected invitation publish to succeed, got %v", err)
	}

	seen := map[string]bool{}
	for range 2 {
		select {
		case subject := <-received:
			seen[subject] = true
		case <-time.After(3 * time.Second):
			t.Fatal("expected multi-topic consumer to receive messages")
		}
	}

	if !seen[common.MessageTopicEmailVerification] {
		t.Fatalf("expected %s to be consumed, got %+v", common.MessageTopicEmailVerification, seen)
	}
	if !seen[common.MessageTopicEmailInvitation] {
		t.Fatalf("expected %s to be consumed, got %+v", common.MessageTopicEmailInvitation, seen)
	}
}

func TestConsumerRouterConsumesAndStoresOffset(t *testing.T) {
	db := openTestDB(t)
	prepareWatermillTables(t, db)
	truncateWatermillTables(t, db)

	publisher := NewPublisher(db)
	router, routerClose, err := NewConsumerRouter(db, Config{
		PollInterval:      10 * time.Millisecond,
		BatchSize:         1,
		RetryDelay:        10 * time.Millisecond,
		ProcessingTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("expected router creation to succeed, got %v", err)
	}
	defer func() {
		_ = router.Close()
		if routerClose != nil {
			_ = routerClose()
		}
	}()

	received := make(chan string, 1)
	err = AddConsumerHandler(
		router,
		db,
		Config{
			PollInterval:      10 * time.Millisecond,
			BatchSize:         1,
			RetryDelay:        10 * time.Millisecond,
			ProcessingTimeout: time.Second,
		},
		common.MessageBusSubscriptionConfig{
			ConsumerGroup: common.MessageConsumerGroupNotificationService,
			Topics:        []string{common.MessageTopicEmailVerification},
		},
		"router_email_verification",
		common.MessageTopicEmailVerification,
		func(msg *message.Message) error {
			received <- string(msg.Payload)
			return nil
		},
	)
	if err != nil {
		t.Fatalf("expected router handler registration to succeed, got %v", err)
	}

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- router.Run(runCtx)
	}()

	if err = publisher.PublishJSON(context.Background(), common.MessageTopicEmailVerification, map[string]any{
		"job": "router",
	}); err != nil {
		t.Fatalf("expected publish to succeed, got %v", err)
	}

	select {
	case payload := <-received:
		if !strings.Contains(payload, `"job": "router"`) && !strings.Contains(payload, `"job":"router"`) {
			t.Fatalf("expected router payload, got %s", payload)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected router consumer to receive queued message")
	}

	cancel()
	if closeErr := router.Close(); closeErr != nil {
		t.Fatalf("expected router close to succeed, got %v", closeErr)
	}

	select {
	case runErr := <-runErrCh:
		if runErr != nil {
			t.Fatalf("expected router run to exit cleanly, got %v", runErr)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected router run to stop after close")
	}

	var offset struct {
		OffsetAcked                int    `db:"offset_acked"`
		LastProcessedTransactionID string `db:"last_processed_transaction_id"`
	}
	err = db.Get(&offset, `
		SELECT offset_acked,
			last_processed_transaction_id::text AS last_processed_transaction_id
		FROM "watermill_offsets_email_verification"
		WHERE consumer_group = $1
	`, common.MessageConsumerGroupNotificationService)
	if err != nil {
		t.Fatalf("expected consumer group offset row, got %v", err)
	}
	if offset.OffsetAcked <= 0 {
		t.Fatalf("expected positive offset_acked, got %d", offset.OffsetAcked)
	}
	if offset.LastProcessedTransactionID == "" || offset.LastProcessedTransactionID == "0" {
		t.Fatalf("expected last_processed_transaction_id to advance, got %s", offset.LastProcessedTransactionID)
	}
}

func TestConsumerRouterRetriesThenPublishesToDeadLetter(t *testing.T) {
	db := openTestDB(t)
	prepareWatermillTables(t, db)
	truncateWatermillTables(t, db)

	publisher := NewPublisher(db)
	router, routerClose, err := NewConsumerRouter(db, Config{
		PollInterval:      10 * time.Millisecond,
		BatchSize:         1,
		RetryDelay:        10 * time.Millisecond,
		MaxAttempts:       3,
		ProcessingTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("expected router creation to succeed, got %v", err)
	}
	defer func() {
		_ = router.Close()
		if routerClose != nil {
			_ = routerClose()
		}
	}()

	attempts := make(chan int32, 3)
	var attemptCounter int32
	err = AddConsumerHandler(
		router,
		db,
		Config{
			PollInterval:      10 * time.Millisecond,
			BatchSize:         1,
			RetryDelay:        10 * time.Millisecond,
			MaxAttempts:       3,
			ProcessingTimeout: time.Second,
		},
		common.MessageBusSubscriptionConfig{
			ConsumerGroup: common.MessageConsumerGroupNotificationService,
			Topics:        []string{common.MessageTopicEmailVerification},
		},
		"router_email_verification_retry",
		common.MessageTopicEmailVerification,
		func(msg *message.Message) error {
			attempts <- atomic.AddInt32(&attemptCounter, 1)
			return errors.New("forced failure")
		},
	)
	if err != nil {
		t.Fatalf("expected router handler registration to succeed, got %v", err)
	}

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- router.Run(runCtx)
	}()

	if err = publisher.PublishJSON(context.Background(), common.MessageTopicEmailVerification, map[string]any{
		"job": "dead-letter",
	}); err != nil {
		t.Fatalf("expected publish to succeed, got %v", err)
	}

	deadline := time.After(3 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case <-attempts:
		case <-deadline:
			t.Fatalf("expected retry attempt %d to run", i+1)
		}
	}

	var deadLetter struct {
		Payload  string `db:"payload"`
		Metadata string `db:"metadata"`
	}
	for {
		err = db.Get(&deadLetter, `
			SELECT payload::text AS payload,
				metadata::text AS metadata
			FROM "watermill_dead_letter"
			ORDER BY transaction_id DESC, "offset" DESC
			LIMIT 1
		`)
		if err == nil {
			break
		}

		select {
		case <-deadline:
			t.Fatalf("expected dead-letter row, got %v", err)
		default:
			time.Sleep(20 * time.Millisecond)
		}
	}

	if !strings.Contains(deadLetter.Payload, `"job": "dead-letter"`) && !strings.Contains(deadLetter.Payload, `"job":"dead-letter"`) {
		t.Fatalf("expected dead-letter payload, got %s", deadLetter.Payload)
	}
	if !strings.Contains(deadLetter.Metadata, `"reason_poisoned"`) {
		t.Fatalf("expected dead-letter metadata to include poison reason, got %s", deadLetter.Metadata)
	}
	if !strings.Contains(deadLetter.Metadata, `"topic_poisoned":"email_verification"`) &&
		!strings.Contains(deadLetter.Metadata, `"topic_poisoned": "email_verification"`) {
		t.Fatalf("expected dead-letter metadata to include source topic, got %s", deadLetter.Metadata)
	}

	cancel()
	if closeErr := router.Close(); closeErr != nil {
		t.Fatalf("expected router close to succeed, got %v", closeErr)
	}

	select {
	case runErr := <-runErrCh:
		if runErr != nil {
			t.Fatalf("expected router run to exit cleanly, got %v", runErr)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected router run to stop after close")
	}
}

func openTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	host := envOrDefault("POSTGRES_HOST", "")
	user := envOrDefault("POSTGRES_USER", "")
	password := envOrDefault("POSTGRES_PASSWORD", "")
	database := envOrDefault("POSTGRES_DB", "")
	port := envOrDefault("POSTGRES_PORT", "5432")
	sslMode := envOrDefault("POSTGRES_SSL_MODE", "require")
	channelBinding := envOrDefault("POSTGRES_CHANNEL_BINDING", "require")

	if host == "" || user == "" || database == "" {
		t.Skip("postgres integration env is not configured")
	}

	db, err := secondarypostgres.New(secondarypostgres.Config{
		Username:       user,
		Password:       password,
		Database:       database,
		Host:           host,
		Port:           port,
		SSLMode:        sslMode,
		ChannelBinding: channelBinding,
		MaxOpenConns:   4,
		MaxIdleConns:   4,
	})
	if err != nil {
		t.Skipf("postgres connection unavailable: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func prepareWatermillTables(t *testing.T, db *sqlx.DB) {
	t.Helper()

	migration := "20250101070000_create_watermill_pubsub_tables.sql"
	path := filepath.Join("..", "..", "..", "..", "..", "db", "migrations", migration)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected migration file %s to be readable, got %v", migration, err)
	}

	dropWatermillTables(t, db)

	query := extractGooseUpSQL(string(content))
	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("expected migration %s to be prepared, got %v", migration, err)
	}
}

func truncateWatermillTables(t *testing.T, db *sqlx.DB) {
	t.Helper()

	_, err := db.Exec(`
		TRUNCATE TABLE
			"watermill_email_verification",
			"watermill_offsets_email_verification",
			"watermill_email_forgot_password",
			"watermill_offsets_email_forgot_password",
			"watermill_email_invitation",
			"watermill_offsets_email_invitation",
			"watermill_export_job_requested",
			"watermill_offsets_export_job_requested",
			"watermill_dead_letter",
			"watermill_offsets_dead_letter"
		RESTART IDENTITY
	`)
	if err != nil {
		t.Fatalf("expected Watermill tables truncate to succeed, got %v", err)
	}
}

func dropWatermillTables(t *testing.T, db *sqlx.DB) {
	t.Helper()

	_, err := db.Exec(`
			DROP TABLE IF EXISTS "watermill_offsets_dead_letter";
			DROP TABLE IF EXISTS "watermill_dead_letter";
			DROP TABLE IF EXISTS "watermill_offsets_export_job_requested";
			DROP TABLE IF EXISTS "watermill_export_job_requested";
		DROP TABLE IF EXISTS "watermill_offsets_email_invitation";
		DROP TABLE IF EXISTS "watermill_email_invitation";
		DROP TABLE IF EXISTS "watermill_offsets_email_forgot_password";
		DROP TABLE IF EXISTS "watermill_email_forgot_password";
		DROP TABLE IF EXISTS "watermill_offsets_email_verification";
		DROP TABLE IF EXISTS "watermill_email_verification";
		DROP TABLE IF EXISTS message_queue;
	`)
	if err != nil {
		t.Fatalf("expected old Watermill tables drop to succeed, got %v", err)
	}
}

func extractGooseUpSQL(content string) string {
	var lines []string
	inUp := false

	for _, line := range strings.Split(content, "\n") {
		switch {
		case strings.HasPrefix(line, "-- +goose Up"):
			inUp = true
			continue
		case strings.HasPrefix(line, "-- +goose Down"):
			inUp = false
		}

		if !inUp {
			continue
		}

		if strings.HasPrefix(line, "-- +goose") {
			continue
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}
