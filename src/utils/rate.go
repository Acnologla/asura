package utils

import "time"

func Try(callback func() error, limit int) {
	for i := 0; i < limit; i++ {
		err := callback()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
}
