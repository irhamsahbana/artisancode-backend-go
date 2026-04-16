package postgres

import "testing"

func TestNormalizeConnectionSetting(t *testing.T) {
	t.Run("defaults to require when empty", func(t *testing.T) {
		got, err := normalizeConnectionSetting("sslmode", "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != "require" {
			t.Fatalf("expected require, got %s", got)
		}
	})

	t.Run("accepts require", func(t *testing.T) {
		got, err := normalizeConnectionSetting("channel_binding", "require")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != "require" {
			t.Fatalf("expected require, got %s", got)
		}
	})

	t.Run("accepts disable", func(t *testing.T) {
		got, err := normalizeConnectionSetting("sslmode", "disable")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != "disable" {
			t.Fatalf("expected disable, got %s", got)
		}
	})

	t.Run("rejects unsupported values", func(t *testing.T) {
		_, err := normalizeConnectionSetting("sslmode", "verify-full")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
