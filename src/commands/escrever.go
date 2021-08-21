package commands

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
	"strings"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"escrever", "escrita", "mock"},
		Run:       runEscrever,
		Available: true,
		Cooldown:  1,
		Usage:     "j!escrever <texto>",
		Help:      "Escreva igual autista",
	})
}

func runEscrever(session disgord.Session, msg *disgord.Message, args []string) {
	if len(args) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+" diga algo para eu escrever!")
		return
	}
	text := ""
	str := strings.Join(args, " ")
	for i := 0; i < len(str); i++ {
		if i%2 == 0 {
			text += strings.ToUpper(string(str[i]))
		} else {
			text += strings.ToLower(string(str[i]))
		}
	}
	msg.Reply(context.Background(), session, msg.Author.Mention()+", "+text)
}
