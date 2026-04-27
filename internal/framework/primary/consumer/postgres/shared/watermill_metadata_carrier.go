package shared

import watermillMessage "github.com/ThreeDotsLabs/watermill/message"

type WatermillMetadataCarrier watermillMessage.Metadata

func (c WatermillMetadataCarrier) Get(key string) string {
	return watermillMessage.Metadata(c).Get(key)
}

func (c WatermillMetadataCarrier) Set(key, value string) {
	watermillMessage.Metadata(c).Set(key, value)
}

func (c WatermillMetadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}

	return keys
}
