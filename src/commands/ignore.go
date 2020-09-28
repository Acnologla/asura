package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"ignore", "ignorar","privar_rinha"},
		Run:       runIgnore,
		Available: true,
		Cooldown:  3,
		Usage:     "j!ignore",
		Help:      "Ignore automaticamente todos os pedidos de rinha!",
	})
}

func runIgnore(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author
	
	galo, _ := rinha.GetGaloDB(user.ID)

	if !galo.Ignore {
		msg.Reply(context.Background(), session, "Voce não receberá mais pedidos de rinha!")
	} else {
		msg.Reply(context.Background(), session, "Voce poderá receber pedidos de rinha agora!")
	}

	galo.Ignore = !galo.Ignore
	rinha.SaveGaloDB(user.ID, galo)
}
