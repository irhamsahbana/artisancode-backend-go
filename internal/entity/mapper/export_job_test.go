package mapper

import (
	"encoding/json"
	"testing"

	"codebase-app/internal/entity/restentity"
)

func TestExportJobParamsToJSON_IncludesLanguage(t *testing.T) {
	raw, err := ExportJobParamsToJSON(restentity.CreateExportJobReq{Q: "budi"}, "en")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var payload map[string]any
	if err = json.Unmarshal([]byte(raw), &payload); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if got := payload["language"]; got != "en" {
		t.Fatalf("expected language to be persisted, got %#v", got)
	}
}
