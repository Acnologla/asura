package utils

import "time"

func FormatDate(ts int64) string {
	return time.Unix(ts, 0).Format("02/01/2006")
}
