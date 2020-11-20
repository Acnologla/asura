package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"changename", "galoname", "setname"},
		Run:       runChangeName,
		Available: true,
		Cooldown:  3,
		Usage:     "j!changename",
		Help:      "Troque o nome do seu galo",
		Category:  1,
	})
}

func runChangeName(session disgord.Session, msg *disgord.Message, args []string) {
	if len(args) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Use j!changename <novo nome>")
		return
	}
	if len(args[0]) > 25 || 3 > len(args[0]) {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", O nome do seu galo deve ter entre 3 e 25 caracteres")
		return
	}
	galo, _ := rinha.GetGaloDB(msg.Author.ID)

	if galo.Money >= 50 {
		rinha.ChangeMoney(msg.Author.ID, -50)
		rinha.UpdateGaloDB(msg.Author.ID, map[string]interface{}{
			"name": args[0],
		})
		msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce trocou o nome do seu galo para **%s** com sucesso\nCustou 50 de dinheiro", msg.Author.Mention(), args[0]))
	} else {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce precisa ter 100 de dinheiro para trocar o nome do seu galo\nUse j!lootbox para ver seu dinheiro")
		return
	}
}
