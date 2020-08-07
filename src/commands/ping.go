package commands

import (
	"context"
	"github.com/andersfylling/disgord"
	"asura/src/handler"
)

func init() {
	handler.Register(handler.Command {
		Aliases: []string{"ping"},
		Run: runPing,
	})
}

func runPing(session disgord.Session, evt *disgord.MessageCreate, args []string){
	msg := evt.Message
	msg.Reply(context.Background(), session, "pong")
}