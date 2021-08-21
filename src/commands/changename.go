package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strings"
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
	text := strings.Join(args, " ")
	if len(text) > 25 || 3 > len(text) {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", O nome do seu galo deve ter entre 3 e 25 caracteres")
		return
	}
	price := 100
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if rinha.IsVip(galo) {
		price = 0
	}
	err := rinha.ChangeMoney(msg.Author.ID, -price, price)
	if err == nil {
		rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
			gal.Name = text
			return gal, nil
		})
		msg.Reply(context.Background(), session, fmt.Sprintf("%s, Você trocou o nome do seu galo para `%s` com sucesso\nCustou %d de dinheiro", msg.Author.Mention(), text, price))
	} else {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Você precisa ter 100 de dinheiro para trocar o nome do seu galo\nUse j!lootbox para ver seu dinheiro")
		return
	}
}
