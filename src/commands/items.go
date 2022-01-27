package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"items"},
		Run:       runItem,
		Available: true,
		Cooldown:  10,
		Usage:     "j!items",
		Help:      "Equipe items",
		Category:  2,
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
			msg.Reply(context.Background(), session, "Espere a sua rinha terminar para equipar items ou vender items")
			return
		}
		battleMutex.RUnlock()
		value, err := strconv.Atoi(args[0])
		if err != nil {
			msg.Reply(context.Background(), session, "Use j!item <numero do item> para equipar um item ou j!item <numero do item> vender para vender um item")
			return
		}
		if !(value >= 0 && len(galo.Items) > value) {
			msg.Reply(context.Background(), session, "Número inválido")
			return
		}
		var text string
		if len(args) > 1 {
			if args[1] == "vender" {
				utils.Confirm(fmt.Sprintf("Você deseja vender o item **%s**?", rinha.Items[galo.Items[value]].Name), msg.ChannelID, msg.Author.ID, func() {
					rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
						if value >= 0 && len(galo.Items) > value {
							newItem := galo.Items[value]
							item := rinha.Items[newItem]
							price := rinha.LevelToPrice(*item)
							for i := value; i < len(galo.Items)-1; i++ {
								galo.Items[i] = galo.Items[i+1]
							}
							galo.Money += price
							galo.Items = galo.Items[0 : len(galo.Items)-1]
							text = fmt.Sprintf("%s, vendeu o item %s por %d com sucesso", msg.Author.Mention(), item.Name, price)
							return galo, nil
						} else {
							text = "Numero inválido"
							return galo, nil
						}
					})

				})
				msg.Reply(context.Background(), session, text)
				return
			}
		}
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			if value >= 0 && len(galo.Items) > value {
				newItem := galo.Items[value]
				item := rinha.Items[newItem]
				old := galo.Items[0]
				galo.Items[0] = newItem
				galo.Items[value] = old
				text = fmt.Sprintf("%s, Voc equipou o item %s", msg.Author.Mention(), item.Name)
				return galo, nil
			} else {
				text = "Numero inválido"
				return galo, nil
			}
		})
		msg.Reply(context.Background(), session, text)
	}
}
