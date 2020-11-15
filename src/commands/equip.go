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
			text += fmt.Sprintf("[%d] - %s (Level: **%d**) \n", i, rinha.Classes[otherGalo.Type].Name,rinha.CalcLevel(otherGalo.Xp))
		}
		avatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       65535,
			Title:       "Galos",
			Description: text,
			Footer: &disgord.EmbedFooter{
				IconURL: avatar,
				Text:    "Use j!equipar <numero do galo> para equipar um galo | use j!equipar <numero do galo> remove para remover um galo",
			},
		})
	} else {
		value, err := strconv.Atoi(args[0])
		if err != nil {
			msg.Reply(context.Background(), session, "Use j!equipar <numero do galo> para equipar um galo | use j!equipar <numero do galo> remove para remover um galo")
			return
		}
		if value >= 0 && len(galo.Galos) > value {
			if len(args) >= 2{
				if args[1] == "remove"{
					gal := galo.Galos[value]
					for i:= value; i < len(galo.Galos)-1;i++{
						galo.Galos[i] = galo.Galos[i+1]
					}
					galo.Galos = galo.Galos[0:len(galo.Galos)-1]
					rinha.SaveGaloDB(msg.Author.ID, galo)
					msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce removeu o galo **%s** com sucesso", msg.Author.Mention(), rinha.Classes[gal.Type].Name))
					return
				}
			}
			newGalo := galo.Galos[value]
			old := rinha.SubGalo{
				Xp: galo.Xp,
				Type: galo.Type,
			}
			galo.Type = newGalo.Type
			galo.Xp  = newGalo.Xp
			galo.Galos[value] = old
			galo.Equipped = []int{}
			rinha.SaveGaloDB(msg.Author.ID, galo)
			msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce trocou seu galo **%s** por **%s**", msg.Author.Mention(), rinha.Classes[old.Type].Name, rinha.Classes[galo.Type].Name))
		} else {
			msg.Reply(context.Background(), session, "Numero invalido")

		}
	}
}
