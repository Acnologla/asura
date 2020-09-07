package interpreter

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

var defaultVars = map[string]interface{}{
	"false":             false,
	"true":              true,
	"commands": handler.Commands,
	"currentUserGuilds": &disgord.GetCurrentUserGuildsParams{},
	"getContext": func() interface{} {
		return context.Background()
	},
	"getUser": func(values interface{}) interface{} {
		str, ok := values.(string)
		if ok {
			user, _ := handler.Client.GetUser(context.Background(), utils.StringToID(str))
			return user
		}
		id, ok := values.(disgord.Snowflake)
		if ok {
			user, _ := handler.Client.GetUser(context.Background(), id)
			return user
		}
		return nil
	},
	"print": func(values ...interface{}) interface{} {
		fmt.Println(values...)
		return nil
	},
	"len": func(values interface{}) interface{} {
		arr,ok := toArrInterface(values)
		if !ok{
			return float64(0)
		}
		return float64(len(arr))
	},
	"getClient": func() interface{} {
		return handler.Client
	},
}
