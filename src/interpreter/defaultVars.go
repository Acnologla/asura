package interpreter

import (
	"github.com/andersfylling/disgord"
	"asura/src/handler"
	"asura/src/utils"
	"reflect"
	"context"
	"fmt"
)


var defaultVars = map[string]interface{}{
	"false" : false,
	"true": true,
	"currentUserGuilds":&disgord.GetCurrentUserGuildsParams{},
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
		if ok{
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
		s := reflect.ValueOf(values)
		if s.Kind() != reflect.Slice {
			print("InterfaceSlice() given a non-slice type")
			return nil
		}
	
		arr := make([]interface{}, s.Len())
	
		for i:=0; i<s.Len(); i++ {
			arr[i] = s.Index(i).Interface()
		}
		return float64(len(arr)) 
	},
	"getClient": func() interface{} {
		return handler.Client
	},
}