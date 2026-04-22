package pkg

import "strings"

// Explode splits a delimited string, trims each segment, skips empty values,
// and converts every item into the requested type through mapper.
func Explode[T any](raw string, separator string, mapper func(string) (T, error)) ([]T, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, separator)
	items := make([]T, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}

		mapped, err := mapper(value)
		if err != nil {
			return nil, err
		}

		items = append(items, mapped)
	}

	if len(items) == 0 {
		return nil, nil
	}

	return items, nil
}
