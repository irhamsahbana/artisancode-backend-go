package errmsg

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestMessageCodesExistInAllLocales(t *testing.T) {
	t.Parallel()

	idMessages := loadLocaleMessages(t, "locales/id.toml")
	enMessages := loadLocaleMessages(t, "locales/en.toml")

	for _, code := range AllMessageCodes {
		if _, ok := idMessages[code]; !ok {
			t.Fatalf("message code %q missing from id locale", code)
		}
		if _, ok := enMessages[code]; !ok {
			t.Fatalf("message code %q missing from en locale", code)
		}
	}
}

func TestLocaleKeysMatch(t *testing.T) {
	t.Parallel()

	idMessages := loadLocaleMessages(t, "locales/id.toml")
	enMessages := loadLocaleMessages(t, "locales/en.toml")

	for code := range idMessages {
		if _, ok := enMessages[code]; !ok {
			t.Fatalf("message code %q exists in id locale but not en locale", code)
		}
	}

	for code := range enMessages {
		if _, ok := idMessages[code]; !ok {
			t.Fatalf("message code %q exists in en locale but not id locale", code)
		}
	}
}

func loadLocaleMessages(t *testing.T, path string) map[string]string {
	t.Helper()

	raw, err := localeFiles.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	messages := make(map[string]string)
	if err := toml.Unmarshal(raw, &messages); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	return messages
}
