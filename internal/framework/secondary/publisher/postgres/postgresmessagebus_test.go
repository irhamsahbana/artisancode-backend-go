package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	secondarypostgres "codebase-app/internal/framework/secondary/db/postgres"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
)

func TestPublishJSONInsertsMessage(t *testing.T) {
	db := openTestDB(t)
	prepareMessageQueueTable(t, db)
	truncateMessageQueueTables(t, db)

	publisher := NewPublisher(db)
	err := publisher.PublishJSON(context.Background(), "test.publish.insert", map[string]any{
		"hello": "world",
	})
	if err != nil {
		t.Fatalf("expected no publish error, got %v", err)
	}

	var row struct {
		Subject     string `db:"subject"`
		Status      string `db:"status"`
		PayloadJSON string `db:"payload_json"`
		HeadersJSON string `db:"headers_json"`
	}

	query := `
		SELECT subject, status, payload_json::text AS payload_json, headers_json::text AS headers_json
		FROM message_queue
		WHERE subject = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	err = db.Get(&row, query, "test.publish.insert")
	if err != nil {
		t.Fatalf("expected queued message row, got %v", err)
	}

	if row.Subject != "test.publish.insert" {
		t.Fatalf("expected subject test.publish.insert, got %s", row.Subject)
	}
	if row.Status != "pending" {
		t.Fatalf("expected pending status, got %s", row.Status)
	}
	if !strings.Contains(row.PayloadJSON, `"hello": "world"`) && !strings.Contains(row.PayloadJSON, `"hello":"world"`) {
		t.Fatalf("expected payload json to contain hello world, got %s", row.PayloadJSON)
	}
	if row.HeadersJSON == "" || row.HeadersJSON == "null" {
		t.Fatalf("expected headers json to be present, got %s", row.HeadersJSON)
	}
}

func TestConsumeAndAckMarksProcessed(t *testing.T) {
	db := openTestDB(t)
	prepareMessageQueueTable(t, db)
	truncateMessageQueueTables(t, db)

	publisher := NewPublisher(db)
	manager := NewSubscriptionManager(db, Config{
		PollInterval: 10 * time.Millisecond,
		BatchSize:    1,
		RetryDelay:   50 * time.Millisecond,
		MaxAttempts:  3,
	})

	subscription, err := manager.CreateSubscription(context.Background(), integrationPorts.MessageBusSubscriptionConfig{
		ConsumerName: "test-consumer-ack",
		Subjects:     []string{"test.consume.ack"},
	})
	if err != nil {
		t.Fatalf("expected consumer creation to succeed, got %v", err)
	}

	received := make(chan string, 1)
	consumeCtx, err := subscription.Consume(func(msg integrationPorts.MessageBusMessage) {
		received <- string(msg.Data())
		if ackErr := msg.Ack(); ackErr != nil {
			t.Errorf("expected ack to succeed, got %v", ackErr)
		}
	})
	if err != nil {
		t.Fatalf("expected consume to start, got %v", err)
	}
	defer consumeCtx.Stop()

	err = publisher.PublishJSON(context.Background(), "test.consume.ack", map[string]any{
		"job": "ack",
	})
	if err != nil {
		t.Fatalf("expected publish to succeed, got %v", err)
	}

	select {
	case <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("expected consumer to receive queued message")
	}

	var row struct {
		Status       string     `db:"status"`
		ProcessedAt  *time.Time `db:"processed_at"`
		AttemptCount int        `db:"attempt_count"`
	}

	query := `
		SELECT status, processed_at, attempt_count
		FROM message_queue
		WHERE subject = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	err = db.Get(&row, query, "test.consume.ack")
	if err != nil {
		t.Fatalf("expected processed row, got %v", err)
	}

	if row.Status != "processed" {
		t.Fatalf("expected processed status, got %s", row.Status)
	}
	if row.ProcessedAt == nil {
		t.Fatal("expected processed_at to be set")
	}
	if row.AttemptCount != 0 {
		t.Fatalf("expected attempt_count 0 after ack, got %d", row.AttemptCount)
	}
}

func TestConsumeAndNakRequeuesMessage(t *testing.T) {
	db := openTestDB(t)
	prepareMessageQueueTable(t, db)
	truncateMessageQueueTables(t, db)

	publisher := NewPublisher(db)
	manager := NewSubscriptionManager(db, Config{
		PollInterval: 10 * time.Millisecond,
		BatchSize:    1,
		RetryDelay:   50 * time.Millisecond,
		MaxAttempts:  3,
	})

	subscription, err := manager.CreateSubscription(context.Background(), integrationPorts.MessageBusSubscriptionConfig{
		ConsumerName: "test-consumer-nak",
		Subjects:     []string{"test.consume.nak"},
	})
	if err != nil {
		t.Fatalf("expected consumer creation to succeed, got %v", err)
	}

	received := make(chan struct{}, 1)
	consumeCtx, err := subscription.Consume(func(msg integrationPorts.MessageBusMessage) {
		select {
		case received <- struct{}{}:
		default:
		}
		if nakErr := msg.Nak("temporary failure"); nakErr != nil {
			t.Errorf("expected nak to succeed, got %v", nakErr)
		}
	})
	if err != nil {
		t.Fatalf("expected consume to start, got %v", err)
	}
	defer consumeCtx.Stop()

	err = publisher.PublishJSON(context.Background(), "test.consume.nak", map[string]any{
		"job": "nak",
	})
	if err != nil {
		t.Fatalf("expected publish to succeed, got %v", err)
	}

	select {
	case <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("expected consumer to receive queued message")
	}

	time.Sleep(100 * time.Millisecond)

	var row struct {
		Status       string         `db:"status"`
		AvailableAt  time.Time      `db:"available_at"`
		AttemptCount int            `db:"attempt_count"`
		LastError    sql.NullString `db:"last_error"`
	}

	query := `
		SELECT status, available_at, attempt_count, last_error
		FROM message_queue
		WHERE subject = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	err = db.Get(&row, query, "test.consume.nak")
	if err != nil {
		t.Fatalf("expected requeued row, got %v", err)
	}

	if row.Status != "pending" {
		t.Fatalf("expected pending status after nak, got %s", row.Status)
	}
	if row.AttemptCount != 1 {
		t.Fatalf("expected attempt_count 1 after nak, got %d", row.AttemptCount)
	}
	if !row.LastError.Valid || row.LastError.String != "temporary failure" {
		t.Fatalf("expected last_error to be stored, got %+v", row.LastError)
	}
	if row.AvailableAt.Before(time.Now().UTC().Add(-1 * time.Second)) {
		t.Fatalf("expected available_at to be rescheduled, got %v", row.AvailableAt)
	}
}

func TestConsumeAndNakMovesMessageToDeadLetterAfterMaxAttempts(t *testing.T) {
	db := openTestDB(t)
	prepareMessageQueueTable(t, db)
	truncateMessageQueueTables(t, db)

	publisher := NewPublisher(db)
	manager := NewSubscriptionManager(db, Config{
		PollInterval: 10 * time.Millisecond,
		BatchSize:    1,
		RetryDelay:   50 * time.Millisecond,
		MaxAttempts:  1,
	})

	subscription, err := manager.CreateSubscription(context.Background(), integrationPorts.MessageBusSubscriptionConfig{
		ConsumerName: "test-consumer-dead-letter",
		Subjects:     []string{"test.consume.dead-letter"},
	})
	if err != nil {
		t.Fatalf("expected consumer creation to succeed, got %v", err)
	}

	received := make(chan struct{}, 1)
	consumeCtx, err := subscription.Consume(func(msg integrationPorts.MessageBusMessage) {
		select {
		case received <- struct{}{}:
		default:
		}
		if nakErr := msg.Nak("permanent failure"); nakErr != nil {
			t.Errorf("expected nak to succeed, got %v", nakErr)
		}
	})
	if err != nil {
		t.Fatalf("expected consume to start, got %v", err)
	}
	defer consumeCtx.Stop()

	err = publisher.PublishJSON(context.Background(), "test.consume.dead-letter", map[string]any{
		"job": "dead-letter",
	})
	if err != nil {
		t.Fatalf("expected publish to succeed, got %v", err)
	}

	select {
	case <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("expected consumer to receive queued message")
	}

	time.Sleep(100 * time.Millisecond)

	var queueCount int
	err = db.Get(&queueCount, `SELECT COUNT(*) FROM message_queue WHERE subject = $1`, "test.consume.dead-letter")
	if err != nil {
		t.Fatalf("expected queue count query to succeed, got %v", err)
	}
	if queueCount != 0 {
		t.Fatalf("expected queued message to be removed after dead-lettering, got %d", queueCount)
	}

	var deadLetter struct {
		Subject        string    `db:"subject"`
		AttemptCount   int       `db:"attempt_count"`
		FailureReason  string    `db:"failure_reason"`
		DeadLetteredAt time.Time `db:"dead_lettered_at"`
	}

	err = db.Get(&deadLetter, `
		SELECT subject, attempt_count, failure_reason, dead_lettered_at
		FROM message_queue_dead_letters
		WHERE subject = $1
		ORDER BY dead_lettered_at DESC
		LIMIT 1
	`, "test.consume.dead-letter")
	if err != nil {
		t.Fatalf("expected dead-letter row to exist, got %v", err)
	}

	if deadLetter.AttemptCount != 1 {
		t.Fatalf("expected dead-letter attempt_count 1, got %d", deadLetter.AttemptCount)
	}
	if deadLetter.FailureReason != "permanent failure" {
		t.Fatalf("expected dead-letter failure reason to be stored, got %s", deadLetter.FailureReason)
	}
	if deadLetter.DeadLetteredAt.IsZero() {
		t.Fatal("expected dead_lettered_at to be set")
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

func prepareMessageQueueTable(t *testing.T, db *sqlx.DB) {
	t.Helper()

	migrations := []string{
		"20250101070000_create_message_queue_table.sql",
		"20260418093000_create_message_queue_dead_letters_table.sql",
	}

	for _, migration := range migrations {
		path := filepath.Join("..", "..", "..", "..", "..", "db", "migrations", migration)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("expected migration file %s to be readable, got %v", migration, err)
		}

		query := extractGooseUpSQL(string(content))
		_, err = db.Exec(query)
		if err != nil {
			t.Fatalf("expected migration %s to be prepared, got %v", migration, err)
		}
	}
}

func truncateMessageQueueTables(t *testing.T, db *sqlx.DB) {
	t.Helper()

	_, err := db.Exec(`TRUNCATE TABLE message_queue_dead_letters, message_queue`)
	if err != nil {
		t.Fatalf("expected message queue tables truncate to succeed, got %v", err)
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
