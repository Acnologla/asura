package utils

import (
	"strings"
	"time"
)

func Try(callback func() error, limit int) {
	for i := 0; i < limit; i++ {
		err := callback()
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "rate limite"){
			break
		}
		time.Sleep(time.Second)
	}
}
