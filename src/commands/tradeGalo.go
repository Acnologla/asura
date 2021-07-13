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
		Aliases:   []string{"trocagalos", "tradegalos", "tradeG"},
		Run:       runTradeGal,
		Available: true,
		Cooldown:  3,
		Usage:     "j!trocagalo <@user>",
		Help:      "Troque seus galos",
		Category:  1,
	})
}

func galoToText(galo rinha.Galo, filter rinha.Rarity, galoType int) string {
	text := ""
	_galo := galo
	for i, galo := range galo.Galos {
		name := rinha.Classes[galo.Type]
		if filter == -1 || name.Rarity == filter {
			if !rinha.HaveGalo(galoType, _galo.Galos) {
				text += fmt.Sprintf("[%d] - %s (Raridade: **%s**) \n", i, name.Name, name.Rarity.String())
			}
		}
	}
	return text
}

func runTradeGal(session disgord.Session, msg *disgord.Message, args []string) {
	if len(msg.Mentions) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", mencione alguem para trocar galos")
		return
	}
	user := msg.Mentions[0]
	if user.Bot || user.ID == msg.Author.ID {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Usuario invalido")
		return
	}
	text := fmt.Sprintf("**%s** voce foi convidado para trocar galos com %s", user.Username, msg.Author.Username)
	utils.Confirm(text, msg.ChannelID, msg.Mentions[0].ID, func() {
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
		galos := galoToText(galo, -1, -1)
		newMsg := &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       "Troca de galos",
				Color:       65535,
				Description: galos,
				Footer: &disgord.EmbedFooter{
					Text: "Digite no chat o numero do galo que deseja trocar",
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
			msg.Reply(context.Background(), session, msg.Author.Mention()+", numero invalido.\nTroca cancelada")
			return
		}
		if 0 > i || i >= len(galo.Galos) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", numero invalido.\nTroca cancelada")
			return
		}
		firstGalo := galo.Galos[i]
		advGalos := galoToText(galoAdv, rinha.Classes[firstGalo.Type].Rarity, firstGalo.Type)
		newMsg.Content = user.Mention()
		session.Channel(tradeMsg.ChannelID).Message(tradeMsg.ID).Delete()
		newMsg.Embed.Description = advGalos + "\n\nGalo que voce ira receber: **" + rinha.Classes[firstGalo.Type].Name + "**"
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
		if 0 > j || j >= len(galoAdv.Galos) {
			msg.Reply(context.Background(), session, message.Author.Mention()+", numero invalido.\nTroca cancelada")
			return
		}
		session.Channel(secondTradeMsg.ChannelID).Message(secondTradeMsg.ID).Delete()
		secondGalo := galoAdv.Galos[j]
		if rinha.Classes[secondGalo.Type].Rarity != rinha.Classes[firstGalo.Type].Rarity {
			msg.Reply(context.Background(), session, message.Author.Mention()+", numero invalido.\nTroca cancelada")
			return
		}
		if rinha.HaveGalo(secondGalo.Type, galo.Galos) || rinha.HaveGalo(firstGalo.Type, galoAdv.Galos) {
			msg.Reply(context.Background(), session, message.Author.Mention()+", numero invalido.\nTroca cancelada")
			return
		}
		galoAdv.Galos[j] = firstGalo
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			secondGalo.Xp = 0
			secondGalo.GaloReset = 0
			galo.Galos[i] = secondGalo
			return galo, nil
		})
		rinha.UpdateGaloDB(user.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			firstGalo.Xp = 0
			firstGalo.GaloReset = 0
			galo.Galos[j] = firstGalo
			return galo, nil
		})
		msg.Reply(context.Background(), session, fmt.Sprintf("%s voce trocou o galo **%s** pelo galo **%s** com sucesso", msg.Author.Mention(), rinha.Classes[firstGalo.Type].Name, rinha.Classes[secondGalo.Type].Name))
	})

}
