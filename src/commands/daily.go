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
	topGGCalc := (uint64(time.Now().Unix()) - galo.Daily.Voted) / 60 / 60 / 12
	voted := rinha.HasVoted(msg.Author.ID)
	if voted && topGGCalc >= 1 {
		strike := 0
		money := 0
		xp := 0
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			galo.Daily.Last = uint64(time.Now().Unix())
			if topGGCalc >= 2 {
				galo.Daily.Strike = 0
			}
			if topGGCalc >= 1 && voted {
				galo.Daily.Voted = uint64(time.Now().Unix())
			}
			if rinha.IsVip(galo) {
				money += 15
				xp += 35
			}
			money = 25 + galo.Daily.Strike/5
			xp = 50 + galo.Daily.Strike
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
		need := uint64(time.Now().Unix()) - galo.Daily.Voted
		if topGGCalc >= 1 && !voted {
			msg.Reply(context.Background(), session, fmt.Sprintf("%s, para pegar o bonus diario voce precisa votar em mim.\n Voce pode pegar o bonus diario a cada 12 horas\nLink para votar:\nhttps://top.gg/bot/470684281102925844", msg.Author.Mention()))
		} else {
			msg.Reply(context.Background(), session, fmt.Sprintf("%s, faltam **%d** horas e **%d** minutos para voce poder usar o daily novamente", msg.Author.Mention(), 11-(need/60/60), 59-(need/60%60)))
		}
	}
}
