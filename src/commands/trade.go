package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"troca", "trocaritems", "trade"},
		Run:       runTrade,
		Available: true,
		Cooldown:  3,
		Usage:     "j!trade <@user>",
		Help:      "Troque seus itens",
		Category:  1,
	})
}

func itemsToText(galo rinha.Galo, filter int) string {
	text := ""
	for i, item := range galo.Items {
		name := rinha.Items[item]
		if filter == -1 || name.Level == filter {
			text += fmt.Sprintf("[%d] - %s (Raridade: **%s**) \n%s\n", i, name.Name, rinha.LevelToString(name.Level), rinha.ItemToString(name))
		}
	}
	return text
}

func runTrade(session disgord.Session, msg *disgord.Message, args []string) {
	if len(msg.Mentions) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", mencione alguem para trocar items")
		return
	}
	user := msg.Mentions[0]
	if user.Bot || user.ID == msg.Author.ID {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Usuario invalido")
		return
	}
	confirmMsg, confirmErr := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: msg.Mentions[0].Mention(),
		Embed: &disgord.Embed{
			Color:       65535,
			Description: fmt.Sprintf("**%s** clique na reação abaixo para aceitar a troca", msg.Mentions[0].Username),
		},
	})
	if confirmErr == nil {
		utils.Confirm(confirmMsg, msg.Mentions[0].ID, func() {
			battleMutex.RLock()
			if currentBattles[msg.Author.ID] != "" {
				battleMutex.RUnlock()
				msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
				return
			}
			if currentBattles[msg.Mentions[0].ID] != "" {
				battleMutex.RUnlock()
				msg.Reply(context.Background(), session, "Este usuario ja esta lutando com o "+currentBattles[msg.Mentions[0].ID])
				return
			}
			battleMutex.RUnlock()
			lockBattle(msg.Author.ID, msg.Mentions[0].ID, msg.Author.Username, msg.Mentions[0].Username)
			defer unlockBattle(msg.Author.ID, msg.Mentions[0].ID)
			galo, _ := rinha.GetGaloDB(msg.Author.ID)
			galoAdv, _ := rinha.GetGaloDB(user.ID)
			items := itemsToText(galo, -1)
			newMsg := &disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title:       "Troca de items",
					Color:       65535,
					Description: items,
					Footer: &disgord.EmbedFooter{
						Text: "Digite no chat o numero do item que deseja trocar",
					},
				},
				Content: msg.Author.Mention(),
			}
			tradeMsg, err := msg.Reply(context.Background(), session, newMsg)
			if err != nil {
				return
			}
			message := handler.CreateMessageCollector(msg.ChannelID, func(message *disgord.Message) bool {
				return message.Author.ID == msg.Author.ID
			})
			if message == nil{
				return
			}
			i, err := strconv.Atoi(message.Content)
			if err != nil {
				msg.Reply(context.Background(), session, msg.Author.Mention()+", numero invalido.\nTroca cancelada")
				return
			}
			if 0 > i || i >= len(galo.Items) {
				msg.Reply(context.Background(), session, msg.Author.Mention()+", numero invalido.\nTroca cancelada")
				return
			}
			firstItem := galo.Items[i]
			advItems := itemsToText(galoAdv, rinha.Items[firstItem].Level)
			newMsg.Content = user.Mention()
			session.Channel(tradeMsg.ChannelID).Message(tradeMsg.ID).Delete()
			newMsg.Embed.Description = advItems + "\n\nItem que voce ira receber: **" + rinha.Items[firstItem].Name + "**"
			secondTradeMsg, err := msg.Reply(context.Background(), session, newMsg)
			if err != nil {
				return
			}
			message = handler.CreateMessageCollector(msg.ChannelID, func(message *disgord.Message) bool {
				return message.Author.ID == user.ID
			})
			if message == nil{
				return
			}
			j, err := strconv.Atoi(message.Content)
			if err != nil {
				msg.Reply(context.Background(), session, message.Author.Mention()+", numero invalido.\nTroca cancelada")
				return
			}
			if 0 > j || j >= len(galoAdv.Items) {
				msg.Reply(context.Background(), session, message.Author.Mention()+", numero invalido.\nTroca cancelada")
				return
			}
			session.Channel(secondTradeMsg.ChannelID).Message(secondTradeMsg.ID).Delete()
			secondItem := galoAdv.Items[j]
			if rinha.Items[secondItem].Level != rinha.Items[firstItem].Level{
				msg.Reply(context.Background(), session, message.Author.Mention()+", numero invalido.\nTroca cancelada")
				return
			}
			galo.Items[i] = secondItem
			galoAdv.Items[j] = firstItem
			rinha.UpdateGaloDB(msg.Author.ID, map[string]interface{}{
				"items": galo.Items,
			})
			rinha.UpdateGaloDB(user.ID, map[string]interface{}{
				"items": galoAdv.Items,
			})
			msg.Reply(context.Background(), session, fmt.Sprintf("%s voce trocou o item **%s** pelo item **%s** com sucesso", msg.Author.Mention(), rinha.Items[firstItem].Name, rinha.Items[secondItem].Name))
		})
	}
}
