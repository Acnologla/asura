package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"private", "privar"},
		Run:       runPrivate,
		Available: true,
		Cooldown:  1,
		Usage:     "j!private",
		Help:      "Prive seu historico de avatares",
	})
}

func runPrivate(session disgord.Session, msg *disgord.Message, args []string) {
	var acc bool
	id := strconv.FormatUint(uint64(msg.Author.ID), 10)
	ctx := context.Background()
	if err := database.Database.NewRef("private/"+id).Get(ctx, &acc); err != nil {
		fmt.Println(err)
		return
	}
	newVal := !acc
	if err := database.Database.NewRef("private/"+id).Set(ctx, &(newVal)); err != nil {
		fmt.Println(err)
		return
	}
	if acc {
		msg.Reply(ctx, session, msg.Author.Mention()+", Seu perfil nao esta mais privado")
	} else {
		msg.Reply(ctx, session, msg.Author.Mention()+", Seu perfil foi privado som sucesso")
	}
}
