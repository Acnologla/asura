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
			text += fmt.Sprintf("[%d] - %s\n", i, rinha.Classes[otherGalo].Name)
		}
		avatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       65535,
			Title:       "Galos",
			Description: text,
			Footer: &disgord.EmbedFooter{
				IconURL: avatar,
				Text:    "Use j!equipar <numero do galo> para equipar um galo",
			},
		})
	} else {
		value, err := strconv.Atoi(args[0])
		if err != nil {
			msg.Reply(context.Background(), session, "Use j!equipar <numero do galo> para equipar um galo")
			return
		}
		if value >= 0 && len(galo.Galos) > value {
			newGalo := galo.Galos[value]
			old := galo.Type
			galo.Type = newGalo
			galo.Galos[value] = old
			galo.Equipped = []int{}
			rinha.SaveGaloDB(msg.Author.ID, galo)
			msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce trocou seu galo **%s** por **%s**", msg.Author.Mention(), rinha.Classes[old].Name, rinha.Classes[galo.Type].Name))
		} else {
			msg.Reply(context.Background(), session, "Numero invalido")

		}
	}
}
