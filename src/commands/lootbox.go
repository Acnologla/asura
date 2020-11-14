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
		Aliases:   []string{"lootbox", "lb"},
		Run:       runLootbox,
		Available: true,
		Cooldown:  3,
		Usage:     "j!lootbox",
		Help:      "Abra lootboxs",
		Category:  1,
	})
}

func runLootbox(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	normal := func() {
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title:       "Lootbox",
			Color:       65535,
			Description: fmt.Sprintf("Lootbox: **%d**\nMoney: **%d**\n\nUse `j!lootbox buy` para comprar lootbox\nUse `j!lootbox open` para abrir lootbox", galo.Lootbox, galo.Money),
		})
	}
	if len(args) == 0 {
		normal()
		return
	}
	if args[0] == "open" || args[0] == "abrir" {
		if galo.Lootbox == 0 {
			msg.Reply(context.Background(), session, "Voce tem 0 lootboxs, use `j!lootbox buy` para comprar uma")
			return
		}
		result := rinha.Open()
		if !rinha.IsIntInList(result, galo.Galos) {
			galo.Galos = append(galo.Galos, result)
		}
		galo.Lootbox--
		rinha.SaveGaloDB(msg.Author.ID, galo)
		newGalo := rinha.Classes[result]
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       newGalo.Rarity.Color(),
			Title:       "Lootbox open",
			Description: "Voce abriu uma lootbox e ganhou o galo **" + newGalo.Name + "**\nRaridade: " + newGalo.Rarity.String(),
			Thumbnail: &disgord.EmbedThumbnail{
				URL: rinha.Sprites[0][result-1],
			},
		})
	} else if args[0] == "buy" || args[0] == "comprar" {
		if 500 > galo.Money {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce precisa ter 500 de dinheiro para comprar uma lootbox, use `j!lootbox` para ver seu dinheiro")
			return
		}
		galo.Money -= 500
		galo.Lootbox++
		rinha.SaveGaloDB(msg.Author.ID, galo)
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce comprou uma lootbox use `j!lootbox open` para abrir")
	} else {
		normal()
	}
}
