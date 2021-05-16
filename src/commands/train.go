package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"time"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"train", "treino", "treinar"},
		Run:       runTrain,
		Available: true,
		Cooldown:  12,
		Usage:     "j!train",
		Category:  1,
		Help:      "Batalhe contra um galo",
	})
}

func runTrain(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Type == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem um galo, use j!galo para criar um")
		return
	}
	confirmMsg, confirmErr := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Color:       65535,
			Description: fmt.Sprintf("**%s** clique na reação abaixo para lutar", msg.Author.Username),
		},
	})
	if confirmErr == nil {
		utils.Confirm(confirmMsg, msg.Author.ID, func() {
			battleMutex.RLock()
			if currentBattles[msg.Author.ID] != "" {
				battleMutex.RUnlock()
				msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
				return
			}
			battleMutex.RUnlock()
			galoAdv := rinha.Galo{
				Xp:   galo.Xp,
				Type: rinha.GetRand(),
			}
			if len(galo.Items) > 0 {
				randItem := rinha.GetItemByLevel(rinha.Items[galo.Items[0]].Level)
				galoAdv.Items = []int{randItem}
			}
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
			rinha.CompleteMission(msg.Author.ID, galo, galoAdv, winner == 0, msg)
			if winner == 0 {
				xpOb := utils.RandInt(11) + 12
				money := 4
				clanMsg := ""
				isLimit := rinha.IsInLimit(galo, msg.Author.ID)
				if isLimit {
					need := uint64(time.Now().Unix()) - galo.TrainLimit.LastReset
					msg.Reply(context.Background(), session, &disgord.Embed{
						Color:       16776960,
						Title:       "Train",
						Description: fmt.Sprintf("Voce excedeu seu limite de trains ganhos por dia. Faltam %d horas e %d minutos para voce poder usar mais", 23-(need/60/60), 59-(need/60%60)),
					})
					telemetry.Debug(fmt.Sprintf("%s in rinha limit", msg.Author.Username), map[string]string{
						"id": fmt.Sprintf("%d", msg.Author.ID),
					})
					return
				}
				rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
					if rinha.IsVip(galo) {
						xpOb += 8
					}
					galo.TrainLimit.Times++
					if galo.Clan != "" {
						clan := rinha.GetClan(galo.Clan)
						xpOb++
						level := rinha.ClanXpToLevel(clan.Xp)
						if level >= 2 {
							xpOb += 2
						}
						if level >= 4 {
							money++
						}
						if level >= 5 {
							money++
						}
						go rinha.CompleteClanMission(galo.Clan, msg.Author.ID)
						clanMsg = "\nGanhou **1** de xp para seu clan"
					}
					galo.Xp += xpOb
					galo.Money += money
					return galo, nil
				})
				msg.Reply(context.Background(), session, &disgord.Embed{
					Color:       16776960,
					Title:       "Train",
					Description: fmt.Sprintf("Parabens %s, voce venceu\nGanhou **%d** de dinheiro e **%d** de xp%s", msg.Author.Username, money, xpOb, clanMsg),
				})
				galo.Xp += xpOb
				sendLevelUpEmbed(msg, session, &galo, msg.Author, xpOb)
			} else {
				msg.Reply(context.Background(), session, &disgord.Embed{
					Color:       16711680,
					Title:       "Train",
					Description: fmt.Sprintf("Parabens %s, voce perdeu. Use j!train para treinar novamente", msg.Author.Username),
				})
			}
		})
	}

}
