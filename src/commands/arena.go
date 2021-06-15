package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"arena", "coliseu"},
		Run:       runArena,
		Available: true,
		Cooldown:  5,
		Usage:     "j!arena",
		Category:  1,
		Help:      "Batalhe na arena",
	})
}

func logArenaFinish(user *disgord.User, xp, money int) {
	tag := user.Username + "#" + user.Discriminator.String()
	telemetry.Debug(tag+" Finish Arena", map[string]string{
		"user":  strconv.FormatUint(uint64(user.ID), 10),
		"xp":    strconv.Itoa(xp),
		"money": strconv.Itoa(money),
	})
}

func runArena(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Type == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem um galo, use j!galo para criar um")
		return
	}

	if len(args) == 0 {
		authorAvatar, _ := msg.Author.AvatarURL(512, true)
		text := "Use **j!arena ingresso** para comprar um ingresso para a arena custa **500** de dinheiro"
		if galo.Arena.Active {
			text = fmt.Sprintf("Vitorias: **%d/12**\nDerrotas: **%d/3**\nUse **j!arena batalha** para batalhar na arena", galo.Arena.Win, galo.Arena.Lose)
		}
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title: "Arena",
			Footer: &disgord.EmbedFooter{
				Text:    msg.Author.Username,
				IconURL: authorAvatar,
			},
			Color:       65535,
			Description: text,
		})
		return
	} else if args[0] == "ingresso" && !galo.Arena.Active {
		rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
			if gal.Money >= 500 {
				gal.Money -= 500
				gal.Arena.Active = true
				msg.Reply(context.Background(), session, "Voce comprou um ingresso para a arena use **j!arena batalha** para batalhar")
			} else {
				msg.Reply(context.Background(), session, "Voce precisa ter **500** de dinheiro para comprar um ingresso na arena")
			}
			return gal, nil
		})
		return
	}
	if !galo.Arena.Active {
		msg.Reply(context.Background(), session, "Voce precisa ter um ingresso para batalhar na arena, use **j!arena** para comprar um")
		return
	}
	battleMutex.RLock()
	if currentBattles[msg.Author.ID] != "" {
		battleMutex.RUnlock()
		msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
		return
	}
	battleMutex.RUnlock()
	advType := galo.Type
	if advType == int(rinha.Epic) {
		advType = rinha.GetEpicOrLegendary()
	} else {
		rinha.GetRandByType(rinha.Classes[galo.Type].Rarity)
	}
	galoAdv := rinha.Galo{
		Xp:   rinha.CalcXP(rinha.CalcLevel(galo.Xp) + 1),
		Type: advType,
	}
	LockEvent(msg.Author.ID, "Arena "+rinha.Classes[galoAdv.Type].Name)
	defer UnlockEvent(msg.Author.ID)
	galoAdv.Items = galo.Items
	winner, _ := ExecuteRinha(msg, session, rinhaOptions{
		galoAuthor:  galo,
		galoAdv:     galoAdv,
		authorName:  rinha.GetName(msg.Author.Username, galo),
		advName:     "Arena " + rinha.Classes[galoAdv.Type].Name,
		authorLevel: rinha.CalcLevel(galo.Xp),
		advLevel:    rinha.CalcLevel(galoAdv.Xp),
	})
	if winner == 0 {
		rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
			gal.Arena.Win++
			if gal.Arena.Win >= 12 {
				var xp, money int
				gal, xp, money = rinha.CalcArena(gal)
				msg.Reply(context.Background(), session, &disgord.Embed{
					Color:       16776960,
					Title:       "Arena",
					Description: fmt.Sprintf("Parabens %s, voce atingiu o limite de vitorias na arena\nPremios:\nXp: **%d**\nMoney: **%d**", msg.Author.Username, xp, money),
				})
				logArenaFinish(msg.Author, xp, money)
				return gal, nil
			} else {
				msg.Reply(context.Background(), session, &disgord.Embed{
					Color:       16776960,
					Title:       "Arena",
					Description: fmt.Sprintf("Parabens %s, voce venceu\n %d/12 Vitorias", msg.Author.Username, gal.Arena.Win),
				})
				return gal, nil
			}
		})
	} else {
		rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
			gal.Arena.Lose++
			if gal.Arena.Lose >= 3 {
				var money, xp int
				gal, xp, money = rinha.CalcArena(gal)
				msg.Reply(context.Background(), session, &disgord.Embed{
					Color:       16711680,
					Title:       "Arena",
					Description: fmt.Sprintf("Parabens %s, voce atingiu o limite de derrotas na arena\nPremios:\nXp: **%d**\nMoney: **%d**", msg.Author.Username, xp, money),
				})
				logArenaFinish(msg.Author, xp, money)
				return gal, nil
			} else {
				msg.Reply(context.Background(), session, &disgord.Embed{
					Color:       16711680,
					Title:       "Arena",
					Description: fmt.Sprintf("Parabens %s, voce perdeu. %d/3 Derrotas", msg.Author.Username, gal.Arena.Lose),
				})
				return gal, nil
			}
		})
	}
}
