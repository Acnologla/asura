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
		Aliases:   []string{"items"},
		Run:       runItem,
		Available: true,
		Cooldown:  2,
		Usage:     "j!items",
		Help:      "Equipe items",
		Category:  1,
	})
}

func runItem(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if len(args) == 0 {
		text := ""
		if len(galo.Items) > 0 {
			equipped := rinha.Items[galo.Items[0]]
			text = fmt.Sprintf("**Item equipado**\n %s (Raridade: **%s**)\n%s\n\n**Inventory**\n", equipped.Name, rinha.LevelToString(equipped.Level), rinha.ItemToString(equipped))
		}
		for i, item := range galo.Items {
			name := rinha.Items[item]
			text += fmt.Sprintf("[%d] - %s (Raridade: **%s**) \n%s\n", i, name.Name, rinha.LevelToString(name.Level), rinha.ItemToString(name))
		}
		if text == "" {
			text = "Você não tem nenhum item, para conseguir items use j!dungeon"
		}
		avatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       65535,
			Title:       "Items",
			Description: text,
			Footer: &disgord.EmbedFooter{
				IconURL: avatar,
				Text:    "Use j!item <numero do item> para equipar um item",
			},
		})
	} else {
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere a sua rinha terminar para equipar items")
			return
		}
		battleMutex.RUnlock()
		value, err := strconv.Atoi(args[0])
		if err != nil {
			msg.Reply(context.Background(), session, "Use j!item <numero do item> para equipar um item")
			return
		}
		if value >= 0 && len(galo.Items) > value {
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				newItem := galo.Items[value]
				old := galo.Items[0]
				galo.Items[0] = newItem
				galo.Items[value] = old
				msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce equipou o item %s", msg.Author.Mention(), rinha.Items[newItem].Name))
				return galo, nil
			})
		} else {
			msg.Reply(context.Background(), session, "Numero inválido")
		}
	}
}
