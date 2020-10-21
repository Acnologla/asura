package interpreter

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"time"
	"github.com/andersfylling/disgord"
)

var defaultVars = map[string]interface{}{
	"false":             false,
	"true":              true,
	"commands":          &handler.Commands,
	"currentUserGuilds": &disgord.GetCurrentUserGuildsParams{},
	"getContext": func() interface{} {
		return context.Background()
	},
	"reply": func(msg interface{},val interface{}) interface{}{
		message := msg.(*disgord.Message)
		newMsg,_ := message.Reply(context.Background(),handler.Client,val)
		return newMsg
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
	"sleep": func(values interface{}) interface{}{
		time.Sleep(time.Duration(values.(float64)) * time.Second)
		return nil
	},
	"len": func(values interface{}) interface{} {
		arr, ok := toArrInterface(values)
		if !ok {
			return float64(0)
		}
		return float64(len(arr))
	},
	"append": func (values interface{},item interface{}) interface{}{
		arr, ok := toArrInterface(values)
		if !ok{
			return nil
		}
		return append(arr,item)
	},
	"getClient": func() interface{} {
		return handler.Client
	},
}
