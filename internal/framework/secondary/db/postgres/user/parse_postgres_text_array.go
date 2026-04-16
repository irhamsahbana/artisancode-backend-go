package repository

import "github.com/lib/pq"

func parsePostgresTextArray(value string) ([]string, error) {
	if value == "" {
		return []string{}, nil
	}

	items := make([]string, 0)
	if err := pq.Array(&items).Scan(value); err != nil {
		return nil, err
	}

	return items, nil
}
