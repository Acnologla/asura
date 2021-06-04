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
		Aliases:   []string{"dungeon", "dg", "boss"},
		Run:       runDungeon,
		Available: true,
		Cooldown:  5,
		Usage:     "j!dungeon",
		Category:  1,
		Help:      "Adentre na dungeon",
	})
}

func runDungeon(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Type == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem um galo, use j!galo para criar um")
		return
	}
	if len(args) == 0 {
		authorAvatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title: "Dungeon",
			Footer: &disgord.EmbedFooter{
				Text:    msg.Author.Username,
				IconURL: authorAvatar,
			},
			Color:       65535,
			Description: fmt.Sprintf("Voce esta no andar **%d**\nUse j!dungeon battle para batalhar contra o chefe", galo.Dungeon),
		})
		return
	}
	battleMutex.RLock()
	if currentBattles[msg.Author.ID] != "" {
		battleMutex.RUnlock()
		msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
		return
	}
	battleMutex.RUnlock()
	if len(rinha.Dungeon) == galo.Dungeon {
		galo.Dungeon = 0
		galo.DungeonReset += 1
		msg.Reply(context.Background(), session, "Parabens, você terminou a dungeon e agora pode recomeçar (cuidado)!")
	}

	dungeon := rinha.Dungeon[galo.Dungeon]
	galoAdv := dungeon.Boss
	LockEvent(msg.Author.ID, "Boss "+rinha.Classes[galoAdv.Type].Name)
	defer UnlockEvent(msg.Author.ID)
	multiplier := 1 + galo.DungeonReset

	AdvLVL := rinha.CalcLevel(galoAdv.Xp) * multiplier

	ngaloAdv := rinha.Galo{
		Xp:   rinha.CalcXP(AdvLVL) + 1,
		Type: galoAdv.Type,
	}
	winner, _ := ExecuteRinha(msg, session, rinhaOptions{
		galoAuthor:  galo,
		galoAdv:     ngaloAdv,
		authorName:  rinha.GetName(msg.Author.Username, galo),
		advName:     "Boss " + rinha.Classes[galoAdv.Type].Name,
		authorLevel: rinha.CalcLevel(galo.Xp),
		advLevel:    AdvLVL,
		noItems:     2 > galo.DungeonReset,
	})
	if winner == 0 {
		if galo.DungeonReset != 0 && galo.Dungeon+1 != len(rinha.Dungeon) {
			rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
				gal.Dungeon = galo.Dungeon + 1
				gal.DungeonReset = galo.DungeonReset
				return gal, nil
			})
			msg.Reply(context.Background(), session, &disgord.Embed{
				Color:       16776960,
				Title:       "Dungeon",
				Description: fmt.Sprintf("Parabens %s voce consegiu derrotar o boss e avançar para o andar **%d**", msg.Author.Username, galo.Dungeon+1),
			})
			return
		}
		var endMsg string
		rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
			diffGalo, endMsg2 := rinha.DungeonWin(dungeon.Level, gal)
			endMsg = endMsg2
			diffGalo.Dungeon = gal.Dungeon + 1
			return diffGalo, nil
		})
		tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
		telemetry.Debug(fmt.Sprintf("%s %s", tag, endMsg), map[string]string{
			"user":         strconv.FormatUint(uint64(msg.Author.ID), 10),
			"dungeonLevel": string(galo.Dungeon),
		})
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       16776960,
			Title:       "Dungeon",
			Description: fmt.Sprintf("Parabens %s voce consegiu derrotar o boss e avançar para o andar **%d** %s", msg.Author.Username, galo.Dungeon+1, endMsg),
		})
	} else {
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       16711680,
			Title:       "Dungeon",
			Description: fmt.Sprintf("Parabens %s, voce perdeu. Use j!dungeon battle para tentar novamente", msg.Author.Username),
		})
	}
}
