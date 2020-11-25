package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"math/rand"
	"time"
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
	battleMutex.RLock()
	if currentBattles[msg.Author.ID] != "" {
		battleMutex.RUnlock()
		msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
		return
	}
	battleMutex.RUnlock()
	galoAdv := rinha.Galo{
		Xp:   galo.Xp,
		Type: rinha.DailyGalo,
	}
	if (uint64(time.Now().Unix())-galo.Daily)/60/60/24 >= 1 {
		galo.Daily = uint64(time.Now().Unix())
		LockEvent(msg.Author.ID, "Clone de "+rinha.Classes[galoAdv.Type].Name)
		defer UnlockEvent(msg.Author.ID)
		winner, _ := ExecuteRinha(msg, session, rinhaOptions{
			galoAuthor:  &galo,
			galoAdv:     &galoAdv,
			authorName:  rinha.GetName(msg.Author.Username, galo),
			advName:     "Clone de " + rinha.Classes[galoAdv.Type].Name,
			authorLevel: rinha.CalcLevel(galo.Xp),
			advLevel:    rinha.CalcLevel(galoAdv.Xp),
		})
		if winner == 0 {
			xpOb := rand.Intn(50) + 30
			money := 20
			galo.Xp += xpOb
			rinha.UpdateGaloDB(msg.Author.ID, map[string]interface{}{
				"daily": galo.Daily,
			})
			rinha.ChangeMoney(msg.Author.ID, money, 0)
			updateGaloWin(msg.Author.ID, galo)
			sendLevelUpEmbed(msg, session, &galo, msg.Author, xpOb)
			msg.Reply(context.Background(), session, &disgord.Embed{
				Color:       16776960,
				Title:       "Daily",
				Description: fmt.Sprintf("Parabens %s, voce venceu e ganhou **%d** de dinheiro e **%d** de XP", msg.Author.Username, money, xpOb),
			})
		} else {
			msg.Reply(context.Background(), session, &disgord.Embed{
				Color:       16711680,
				Title:       "Daily",
				Description: fmt.Sprintf("Parabens %s, voce perdeu. Use j!daily para tentar novamente", msg.Author.Username),
			})
		}
	} else {
		need := uint64(time.Now().Unix()) - galo.Daily
		msg.Reply(context.Background(), session, fmt.Sprintf("%s, faltam **%d** horas e **%d** minutos para voce poder usar o daily novamente", msg.Author.Mention(), 23-(need/60/60), 59-(need/60%60)))
	}
}
