package commands

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:    []string{"ping"},
		Run:        runPing,
		Available:  true,
	})
}

func runPing(session disgord.Session, msg *disgord.Message, args []string) {
	msg.Reply(context.Background(), session, "pong")
}
