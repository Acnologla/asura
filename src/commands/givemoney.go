package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"dardinheiro", "givemoney"},
		Run:       runGiveMoney,
		Available: true,
		Cooldown:  5,
		Usage:     "j!givemoney <user>",
		Help:      "De dinheiro para alguem",
		Category:  1,
	})
}

func runGiveMoney(session disgord.Session, msg *disgord.Message, args []string) {
	if 1 >= len(args) || len(msg.Mentions) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Use j!givemoney <user> <quantia>")
		return
	}
	value, err := strconv.Atoi(args[1])
	if err != nil {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Quantia invalida")
		return
	}
	if 0 >= value {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Quantia invalida")
		return
	}
	if msg.Mentions[0].ID == msg.Author.ID || msg.Mentions[0].Bot {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Usuario invalido")
		return
	}
	err = rinha.ChangeMoney(msg.Author.ID, -value, value)
	if err == nil {
		rinha.ChangeMoney(msg.Mentions[0].ID, value, 0)
		msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce deu **%d** de dinheiro a %s com sucesso", msg.Author.Mention(), value, msg.Mentions[0].Username))
	} else {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem essa quantia use j!lootbox para ver seu dinheiro")
		return
	}
}
