package events

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/andersfylling/disgord"
)

var lastLootbox int64 = 0

func SendLootbox(msg *disgord.Message) {
	rarity, lootbox := rinha.MessageRandomLootbox()
	embeds := []*disgord.Embed{
		{
			Color: rarity.Color(),
			Image: &disgord.EmbedImage{
				URL: "https://assetsio.gnwcdn.com/overwatch_loot_box.jpg?width=1200&height=1200&fit=bounds&quality=70&format=jpg&auto=webp",
			},
			Description: fmt.Sprintf("Uma lootbox de raridade **%s** apareceu, clique no botão abaixo para adquirir", rinha.LootNames[lootbox]),
			Title:       "Lootbox",
		},
	}
	message := disgord.CreateMessage{
		Embeds: embeds,
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Pegar Lootbox",
						CustomID: "GetLoobox",
						Style:    disgord.Primary,
					},
				},
			},
		},
	}
	newMessage, err := msg.Reply(context.Background(), handler.Client, message)
	if err == nil {
		handler.RegisterHandler(newMessage.ID, func(ic *disgord.InteractionCreate) {
			done := false
			if done {
				return
			}
			done = true
			u := database.User.GetUser(context.Background(), ic.Member.UserID, "Items")
			database.User.InsertItem(context.Background(), ic.Member.UserID, u.Items, lootbox, entities.LootboxType)
			embeds[0].Description = fmt.Sprintf("<@%s> Você adquiriu uma lootbox **%s**", ic.Member.UserID, rinha.LootNames[lootbox])
			embeds[0].Color = 16776960
			message.Components[0].Components[0].Disabled = true
			message.Components[0].Components[0].Label = "Lootbox pega"
			handler.Client.SendInteractionResponse(context.Background(), ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.CreateInteractionResponseData{
					Embeds:     embeds,
					Components: message.Components,
				},
			})

		}, 100)
	}
}

var lastUser = ""

func IsFlood(msg *disgord.Message) bool {
	if lastUser == msg.Author.ID.String() {
		return true
	}
	lastUser = msg.Author.ID.String()
	return false
}

func RecieveLootbox(msg *disgord.Message) {
	now := time.Now().Unix()
	if rand.Intn(100) < 15 && (now-int64(lastLootbox))/60 > 90 && !IsFlood(msg) {
		lastLootbox = now
		SendLootbox(msg)
	}
}

func HandleMessage(s disgord.Session, h *disgord.MessageCreate) {
	msg := h.Message
	appID := os.Getenv("APP_ID")
	if !msg.Author.Bot {
		if msg.GuildID != 0 {
			for _, user := range msg.Mentions {
				if user.ID.String() == appID {
					msg.Reply(context.Background(), s, "Use /help para ver meus comandos\nCaso meus comandos não aparecam me readicione no servidor com este link:\nhttps://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%%20bot&permissions=8")
					break
				}
			}
			if msg.GuildID.String() == "710179373860519997" {
				RecieveLootbox(msg)
			}
		}
	}
}
