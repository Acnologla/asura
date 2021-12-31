package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"asura/src/utils/rinha/engine"
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
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
	text := fmt.Sprintf("**%s** gostaria de treinar seu galo?", msg.Author.Username)
	utils.Confirm(text, msg.ChannelID, msg.Author.ID, func() {
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Você ja esta lutando com o "+currentBattles[msg.Author.ID])
			return
		}
		battleMutex.RUnlock()
		galoAdv := rinha.Galo{
			Xp:   galo.Xp,
			Type: rinha.GetRand(),
		}
		if rinha.Classes[galo.Type].Rarity == rinha.Common {
			galoAdv.Type = rinha.GetRarityPlusOne(rinha.Common)
		}
		if len(galo.Items) > 0 {
			randItem := rinha.GetItemByLevel(rinha.Items[galo.Items[0]].Level)
			galoAdv.Items = []int{randItem}
		}
		LockEvent(msg.Author.ID, "Clone de "+rinha.Classes[galoAdv.Type].Name)
		defer UnlockEvent(msg.Author.ID)
		winner, _ := engine.ExecuteRinha(msg, session, engine.RinhaOptions{
			GaloAuthor:  galo,
			GaloAdv:     galoAdv,
			AuthorName:  rinha.GetName(msg.Author.Username, galo),
			AdvName:     "Clone de " + rinha.Classes[galoAdv.Type].Name,
			AuthorLevel: rinha.CalcLevel(galo.Xp),
			AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
		}, false)
		if winner == -1 {
			return
		}
		rinha.CompleteMission(msg.Author.ID, galo, galoAdv, winner == 0, msg)
		if winner == 0 {
			xpOb := utils.RandInt(10) + 10
			if rinha.HasUpgrade(galo.Upgrades, 0) {
				xpOb++
				if rinha.HasUpgrade(galo.Upgrades, 0, 1, 1) {
					xpOb += 2
				}
				if rinha.HasUpgrade(galo.Upgrades, 0, 1, 1, 0) {
					xpOb += 3
				}
			}
			xpOb += int(rinha.Classes[galoAdv.Type].Rarity-rinha.Classes[galo.Type].Rarity) * 2
			money := 5
			if len(galo.Items) > 0 {
				item := rinha.GetItem(galo)
				if item.Effect == 8 {
					xpOb += xpOb * int(item.Payload)
				}
				if item.Effect == 9 {
					xpOb += 3
					money++
				}
			}
			if galo.GaloReset > 0 {
				for i := 0; i < galo.GaloReset; i++ {
					xpOb = int(float64(xpOb) * 0.75)
				}
			}
			if rinha.HasUpgrade(galo.Upgrades, 0, 1, 0) {
				money++
			}
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
					"user": fmt.Sprintf("%d", msg.Author.ID),
				})
				return
			}
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				if rinha.IsVip(galo) {
					xpOb += 9
					money++
				}
				galo.UserXp++
				galo.TrainLimit.Times++
				if galo.Clan != "" {
					clan := rinha.GetClan(galo.Clan)
					xpOb++
					level := rinha.ClanXpToLevel(clan.Xp)
					if level >= 2 {
						xpOb++
					}
					if level >= 4 {
						money++
					}
					if level >= 6 {
						money++
					}
					if level >= 8 {
						galo.UserXp++
					}
					go rinha.CompleteClanMission(galo.Clan, msg.Author.ID)
					clanMsg = "\nGanhou **1** de xp para seu clan"
				}
				galo.Win++
				galo.Xp += xpOb
				galo.Money += money
				return galo, nil
			})
			msg.Reply(context.Background(), session, &disgord.Embed{
				Color:       16776960,
				Title:       "Train",
				Description: fmt.Sprintf("Parabens %s, você venceu\nGanhou **%d** de dinheiro e **%d** de xp%s", msg.Author.Username, money, xpOb, clanMsg),
			})
			galo.Xp += xpOb
			sendLevelUpEmbed(msg, session, &galo, msg.Author, xpOb)
		} else {
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				galo.Lose++
				return galo, nil
			})
			msg.Reply(context.Background(), session, &disgord.Embed{
				Color:       16711680,
				Title:       "Train",
				Description: fmt.Sprintf("Parabens %s, você perdeu. Use j!train para treinar novamente", msg.Author.Username),
			})
		}
	})

}
