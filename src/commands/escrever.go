package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"strings"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"escrever", "escrita", "mock"},
		Run:       runEscrever,
		Available: true,
		Cooldown:  1,
		Options: handler.GetOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.STRING,
			Name:        "texto",
			Description: "Texto para escrever",
			Required:    true,
		}),
		Usage: "j!escrever <texto>",
		Help:  "Escreva igual autista",
	})
}

func runEscrever(session disgord.Session, msg *disgord.Message, args []*disgord.ApplicationCommandDataOption) (*disgord.Message, func(disgord.Snowflake)) {
	if len(args) == 0 {
		return handler.CreateMessageContent(msg.Author.Mention() + " diga algo para eu escrever!"), nil
	}
	text := ""
	str := strings.Join(utils.GetAllStringArgs(args), " ")
	for i := 0; i < len(str); i++ {
		if i%2 == 0 {
			text += strings.ToUpper(string(str[i]))
		} else {
			text += strings.ToLower(string(str[i]))
		}
	}
	return handler.CreateMessageContent(msg.Author.Mention() + ", " + text), nil
}
