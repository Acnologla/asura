package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"daily", "bonusdiario", "diario"},
		Run:       runDaily,
		Available: true,
		Cooldown:  5,
		Usage:     "j!daily",
		Category:  1,
		Help:      "Receba seu bonus diario",
	})
}

func runDaily(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Type == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem um galo, use j!galo para criar um")
		return
	}
	calc := (uint64(time.Now().Unix()) - galo.Daily.Last) / 60 / 60 / 24
	if calc >= 1 {
		strike := 0
		money := 0
		xp := 0
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			galo.Daily.Last = uint64(time.Now().Unix())
			if calc >= 2 {
				galo.Daily.Strike = 0
			}
			money = 20 + galo.Daily.Strike/5
			xp = 40 + galo.Daily.Strike
			galo.Daily.Strike++
			galo.Money += money
			galo.Xp += xp
			strike = galo.Daily.Strike
			return galo, nil
		})
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       65535,
			Title:       "Daily",
			Description: fmt.Sprintf("Voce ganhou **%d** de dinheiro, **%d** de xp\n\nStrike: **%d**", money, xp, strike),
		})
	} else {
		need := uint64(time.Now().Unix()) - galo.Daily.Last
		msg.Reply(context.Background(), session, fmt.Sprintf("%s, faltam **%d** horas e **%d** minutos para voce poder usar o daily novamente", msg.Author.Mention(), 23-(need/60/60), 59-(need/60%60)))
	}
}