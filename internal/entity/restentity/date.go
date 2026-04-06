package restentity

import "time"

func TodayDate() string {
	return time.Now().Format("2006-01-02")
}
