package commands

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"suporte", "servidor"},
		Run:       runSupport,
		Available: true,
		Cooldown:  1,
		Usage:     "j!suporte",
		Help:      "Veja meu servidor de suporte",
	})
}

func runSupport(session disgord.Session, msg *disgord.Message, args []string) {
	msg.Reply(context.Background(), session, "https://discord.gg/tdVWQGV")
}
