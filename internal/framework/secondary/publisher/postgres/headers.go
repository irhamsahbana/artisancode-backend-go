package postgres

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/propagation"
)

type mapHeaderCarrier map[string][]string

var _ propagation.TextMapCarrier = mapHeaderCarrier{}

func (c mapHeaderCarrier) Get(key string) string {
	values := c[key]
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (c mapHeaderCarrier) Set(key, value string) {
	c[key] = []string{value}
}

func (c mapHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}

	return keys
}

type metadataHeaderCarrier message.Metadata

var _ propagation.TextMapCarrier = metadataHeaderCarrier{}

func (c metadataHeaderCarrier) Get(key string) string {
	return message.Metadata(c).Get(key)
}

func (c metadataHeaderCarrier) Set(key, value string) {
	message.Metadata(c).Set(key, value)
}

func (c metadataHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}

	return keys
}
