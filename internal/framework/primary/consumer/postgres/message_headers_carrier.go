package consumer

import "go.opentelemetry.io/otel/propagation"

type MessageHeadersCarrier map[string][]string

var _ propagation.TextMapCarrier = MessageHeadersCarrier{}

func (c MessageHeadersCarrier) Get(key string) string {
	values := c[key]
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (c MessageHeadersCarrier) Set(key, value string) {
	c[key] = []string{value}
}

func (c MessageHeadersCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}

	return keys
}
