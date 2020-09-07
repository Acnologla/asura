package commands

import (
	"asura/src/handler"
	"asura/src/interpreter"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strings"
)

func init() {
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
