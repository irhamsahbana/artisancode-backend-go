package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"

	"github.com/rs/zerolog/log"
)

func marshalMetadata(metadata map[string]any) ([]byte, error) {
	if metadata == nil {
		metadata = map[string]any{}
	}

	return json.Marshal(metadata)
}

func parseMetadata(ctx context.Context, raw json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return map[string]any{}, nil
	}

	var metadata map[string]any
	if err := json.Unmarshal(raw, &metadata); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, string(raw)).Msg("Failed to unmarshal metadata")
		return nil, err
	}
	if metadata == nil {
		metadata = map[string]any{}
	}

	return metadata, nil
}
