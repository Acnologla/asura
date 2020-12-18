package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/acnologla/interpreter"
	"github.com/andersfylling/disgord"
	"strings"
)

func init() {
	interpreter.Init(map[string]interface{}{
		"reply": func(msg interface{}, val interface{}) interface{} {
			message := msg.(*disgord.Message)
			newMsg, _ := message.Reply(context.Background(), handler.Client, val)
			return newMsg
		},
		"getClient": func() interface{} {
			return handler.Client
		},
		"commands": &handler.Commands,
	})
	handler.Register(handler.Command{
		Aliases:   []string{"eval", "evl"},
		Run:       runEval,
		Available: false,
		Cooldown:  1,
		Usage:     "j!eval",
		Help:      "eval",
	})
}

func runEval(session disgord.Session, msg *disgord.Message, args []string) {
	if msg.Author.ID != utils.StringToID("365948625676795904") && msg.Author.ID != utils.StringToID("395542409959964672") {
		return
	}
	if len(args) > 0 {
		defer func() {
			err := recover()
			if err != nil {
				msg.Reply(context.Background(), session, fmt.Sprintf("```js\n%v ```", err.(error).Error()))
			}
		}()
		code := strings.Join(args, " ")
		if strings.HasPrefix(code, "```") && strings.HasSuffix(code, "```") {
			code = code[3:]
			code = code[:len(code)-3]
			if strings.HasPrefix(code, "rust") {
				code = code[4:]
			}
		}
		eval := interpreter.Run(code, map[string]interface{}{
			"msg":     msg,
			"session": session,
		})
		msg.Reply(context.Background(), session, fmt.Sprintf("```js\n%v ```", eval))
	}
}
