package commands

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
	"fmt"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"invite","convite","convidar"},
		Run:       runInvite,
		Available: true,
		Cooldown:  1,
		Usage:     "j!invite",
		Help:      "Me convide para o seu servidor",
	})
}

func runInvite(session disgord.Session, msg *disgord.Message, args []string) {
	msg.Reply(context.Background(), session, fmt.Sprintf("%s, https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8",msg.Author.Mention(),"470684281102925844"))
}
