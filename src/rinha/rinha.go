package rinha

import (
	"asura/src/entities"
	"fmt"
	"time"
)

const PityMultiplier = 1

func VipMessage(user *entities.User) string {
	now := uint64(time.Now().Unix())
	if now >= user.Vip {
		return ""
	}
	return fmt.Sprintf("Vip por **%d** dias", (user.Vip-now)/60/60/24)
}
