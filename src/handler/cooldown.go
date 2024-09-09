package handler

import (
	"asura/src/cache"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

func SetCooldown(ctx context.Context, userId disgord.Snowflake, command Command) {
	key := fmt.Sprintf("%s_%s", command.Name, userId.String())

	cache.Client.SetNX(ctx, key, time.Now().Unix(), time.Duration(command.Cooldown)*time.Second)
}

func GetCooldown(ctx context.Context, userId disgord.Snowflake, command Command) (time.Time, bool) {
	key := fmt.Sprintf("%s_%s", command.Name, userId.String())
	value := cache.Client.Get(ctx, key)
	result, _ := value.Result()

	if result == "" {
		return time.Now(), false
	}

	remaining, _ := strconv.ParseInt(result, 10, 64)

	return time.Unix(remaining, 0), true
}
