package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"asura/src/telemetry"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"equipar", "equip", "equipgalo"},
		Run:       runEquip,
		Available: true,
		Cooldown:  2,
		Usage:     "j!equipar",
		Help:      "Equipe outro galo",
		Category:  1,
	})
}

func runEquip(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if len(args) == 0 {
		text := ""
		for i, otherGalo := range galo.Galos {
			name := rinha.Classes[otherGalo.Type].Name
			if otherGalo.Name != "" {
				name = otherGalo.Name
			}
			text += fmt.Sprintf("[%d] - %s (Level: **%d**) \n", i, name, rinha.CalcLevel(otherGalo.Xp))
		}
		if text == "" {
			text = "Voce não tem nenhum galo, para conseguir galos compre lootboxs usando j!lootbox"
		}
		avatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       65535,
			Title:       "Galos",
			Description: text,
			Footer: &disgord.EmbedFooter{
				IconURL: avatar,
				Text:    "Use j!equipar <numero do galo> para equipar um galo | use j!equipar <numero do galo> remove para vender um galo",
			},
		})
	} else {
		value, err := strconv.Atoi(args[0])
		if err != nil {
			msg.Reply(context.Background(), session, "Use j!equipar <numero do galo> para equipar um galo | use j!equipar <numero do galo> remove para vender um galo")
			return
		}
		if value >= 0 && len(galo.Galos) > value {
			if len(args) >= 2 {
				if args[1] == "remove" || args[1] == "vender" {
					gal := galo.Galos[value]
					for i := value; i < len(galo.Galos)-1; i++ {
						galo.Galos[i] = galo.Galos[i+1]
					}
					galo.Galos = galo.Galos[0 : len(galo.Galos)-1]
					rinha.UpdateGaloDB(msg.Author.ID, map[string]interface{}{
						"galos": galo.Galos,
					})
					price := rinha.Sell(rinha.Classes[gal.Type].Rarity,gal.Xp)
					rinha.ChangeMoney(msg.Author.ID,price, 0)
					newGalo := rinha.Classes[gal.Type]
					tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
					telemetry.Debug(fmt.Sprintf("%s Sell %s", tag, newGalo.Name), map[string]string{
						"galo": newGalo.Name,
						"user": strconv.FormatUint(uint64(msg.Author.ID), 10),
						"rarity": newGalo.Rarity.String(),
					})
					msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce vendeu o galo **%s** por **%d** de dinheiro com sucesso", msg.Author.Mention(), rinha.Classes[gal.Type].Name, price))
					return
				}
			}
			newGalo := galo.Galos[value]
			old := rinha.SubGalo{
				Xp:   galo.Xp,
				Type: galo.Type,
				Name: galo.Name,
			}
			galo.Type = newGalo.Type
			galo.Xp = newGalo.Xp
			galo.Name = newGalo.Name
			galo.Galos[value] = old
			galo.Equipped = []int{}
			rinha.UpdateGaloDB(msg.Author.ID, map[string]interface{}{
				"galos":    galo.Galos,
				"equipped": galo.Equipped,
				"xp":       galo.Xp,
				"type":     galo.Type,
				"name":     galo.Name,
			})
			msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce trocou seu galo **%s** por **%s**", msg.Author.Mention(), rinha.Classes[old.Type].Name, rinha.Classes[galo.Type].Name))
		} else {
			msg.Reply(context.Background(), session, "Numero invalido")

		}
	}
}
