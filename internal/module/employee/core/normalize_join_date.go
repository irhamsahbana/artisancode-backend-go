package core

import (
	"fmt"
	"time"

	"codebase-app/internal/entity/coreentity"
)

func normalizeEmployeeJoinDate(data *coreentity.Employee) error {
	if data.JoinDate == nil || *data.JoinDate == "" {
		data.JoinDate = nil
		data.JoinDateTimezone = nil
		return nil
	}

	if data.JoinDateTimezone == nil || *data.JoinDateTimezone == "" {
		return fmt.Errorf("join date timezone is required")
	}

	location, err := time.LoadLocation(*data.JoinDateTimezone)
	if err != nil {
		return err
	}

	parsed, err := time.ParseInLocation("2006-01-02", *data.JoinDate, location)
	if err != nil {
		return err
	}

	normalizedJoinDate := parsed.Format("2006-01-02")
	data.JoinDate = &normalizedJoinDate

	normalizedTimezone := location.String()
	data.JoinDateTimezone = &normalizedTimezone

	return nil
}
