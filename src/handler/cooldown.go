package handler

import (
	"asura/src/cache"
	"context"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

func SetCooldown(ctx context.Context, userId disgord.Snowflake, command Command) {
	cache.Client.SetNX(ctx, command.Name+"_"+userId.String(), time.Now().Unix(), time.Duration(command.Cooldown)*time.Second)
}

func GetCooldown(ctx context.Context, userId disgord.Snowflake, command Command) (time.Time, bool) {
	val := cache.Client.Get(ctx, command.Name+"_"+userId.String())
	result, _ := val.Result()
	if result == "" {
		return time.Now(), false
	}
	n, _ := strconv.ParseInt(result, 10, 64)
	return time.Unix(n, 0), true
}
