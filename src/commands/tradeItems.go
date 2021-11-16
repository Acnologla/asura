package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"trocaitems", "tradeitems", "trade"},
		Run:       runTrade,
		Available: true,
		Cooldown:  3,
		Usage:     "j!trocaitems <@user>",
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
	text := fmt.Sprintf("**%s** você foi convidado para trocar itens com %s", user.Username, msg.Author.Username)
	utils.Confirm(text, msg.ChannelID, msg.Mentions[0].ID, func() {
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Você ja esta lutando com o "+currentBattles[msg.Author.ID])
			return
		}
		if currentBattles[msg.Mentions[0].ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Este usuário ja esta lutando com o "+currentBattles[msg.Mentions[0].ID])
			return
		}
		battleMutex.RUnlock()
		lockBattle(msg.Author.ID, msg.Mentions[0].ID, msg.Author.Username, msg.Mentions[0].Username)
		defer unlockBattle(msg.Author.ID, msg.Mentions[0].ID)
		galo, _ := rinha.GetGaloDB(msg.Author.ID)
		galoAdv, _ := rinha.GetGaloDB(user.ID)
		now := uint64(time.Now().Unix())
		if !rinha.CanTrade(galo) {
			calc := 72 - ((now - galo.Cooldowns.TradeItem) / 60 / 60)
			msg.Reply(context.Background(), session, fmt.Sprintf("%s voce precisa esperar mais %d horas para trocar outro item", msg.Author.Username, calc))
			return
		}
		if !rinha.CanTrade(galoAdv) {
			calc := 72 - ((now - galoAdv.Cooldowns.TradeItem) / 60 / 60)
			msg.Reply(context.Background(), session, fmt.Sprintf("%s voce precisa esperar mais %d horas para trocar outro item", user.Username, calc))
			return
		}
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
		if message == nil {
			return
		}
		i, err := strconv.Atoi(message.Content)
		if err != nil {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", numero inválido.\nTroca cancelada")
			return
		}
		if 0 > i || i >= len(galo.Items) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", numero inválido.\nTroca cancelada")
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
		if message == nil {
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
		if rinha.Items[secondItem].Level != rinha.Items[firstItem].Level {
			msg.Reply(context.Background(), session, message.Author.Mention()+", numero invalido.\nTroca cancelada")
			return
		}
		galoAdv.Items[j] = firstItem
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			galo.Items[i] = secondItem
			galo.Cooldowns.TradeItem = uint64(time.Now().Unix())
			return galo, nil
		})
		rinha.UpdateGaloDB(user.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			galo.Items[j] = firstItem
			galo.Cooldowns.TradeItem = uint64(time.Now().Unix())
			return galo, nil
		})
		msg.Reply(context.Background(), session, fmt.Sprintf("%s voce trocou o item **%s** pelo item **%s** com sucesso", msg.Author.Mention(), rinha.Items[firstItem].Name, rinha.Items[secondItem].Name))
	})

}
