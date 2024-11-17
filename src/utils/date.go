package utils

import (
	"fmt"
	"time"
)

func FormatDate(ts int64) string {
	return time.Unix(ts, 0).Format("02/01/2006")
}

func TimeUntilNextSunday() string {
	now := time.Now()
	daysUntilSunday := (7 - int(now.Weekday())) % 7

	nextSunday := time.Date(
		now.Year(), now.Month(), now.Day()+daysUntilSunday,
		0, 0, 0, 0, time.Local,
	)

	durationUntilSunday := nextSunday.Sub(now)

	days := int(durationUntilSunday.Hours()) / 24
	hours := int(durationUntilSunday.Hours()) % 24
	minutes := int(durationUntilSunday.Minutes()) % 60

	return fmt.Sprintf("%d dias, %d horas, %d minutos", days, hours, minutes)

}
