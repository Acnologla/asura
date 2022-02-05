package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"context"
	"fmt"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "lootbox",
		Description: translation.T("LootboxHelp", "pt"),
		Run:         runLootbox,
		Cooldown:    8,
		Aliases:     []string{"money"},
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "view",
				Description: translation.T("LootboxViewHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "open",
				Description: translation.T("LootboxOpenHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "buy",
				Description: translation.T("LootboxBuyHelp", "pt"),
			},
		),
	})

}

func GenerateBuyOptions() (opts []*disgord.SelectMenuOption) {
	for _, name := range rinha.LootNames {
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       name,
			Description: fmt.Sprintf("%s lootbox", name),
			Value:       name,
		})
	}
	return
}

func runLootbox(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	command := itc.Data.Options[0].Name
	user := database.Client.GetUser(itc.Member.UserID, "Items")
	switch command {
	case "view":
		text := translation.T("LootboxView", "pt", rinha.VipMessage(&user))
		lootbox := rinha.GetLootboxes(&user)
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Lootbox",
						Color:       65535,
						Description: fmt.Sprintf("Money: **%d**\nAsuraCoins: **%d**\nPity: **%d%%**\n\nLootbox comum: **%d**\nLootbox normal: **%d**\nLootbox rara: **%d**\nLootbox epica: **%d**\nLootbox lendaria: **%d**\nLootbox items: **%d**\nLootbox cosmetica: **%d**\n\n%s", user.Money, user.AsuraCoin, user.Pity*rinha.PityMultiplier, lootbox.Common, lootbox.Normal, lootbox.Rare, lootbox.Epic, lootbox.Legendary, lootbox.Items, lootbox.Cosmetic, text),
					},
				},
			},
		}
	case "buy":
		handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Lootbox buy",
						Color:       65535,
						Description: fmt.Sprintf("Money: **%d**\nAsuraCoins: **%d**\n\n%s", user.Money, user.AsuraCoin, rinha.GenerateLootPrices()),
					},
				},
				Components: []*disgord.MessageComponent{
					{
						Type: disgord.MessageComponentActionRow,
						Components: []*disgord.MessageComponent{
							{
								Type:        disgord.MessageComponentButton + 1,
								Options:     GenerateBuyOptions(),
								CustomID:    "type",
								MaxValues:   1,
								Placeholder: "Select lootbox",
							},
						},
					},
				},
			},
		})
		handler.RegisterHandler(itc, func(interaction *disgord.InteractionCreate) {
			if len(interaction.Data.Values) == 0 {
				return
			}
			lb := interaction.Data.Values[0]
			i := rinha.GetLbIndex(lb)
			price := rinha.Prices[i]
			done := false
			database.Client.UpdateUser(itc.Member.UserID, func(u entities.User) entities.User {
				if price[0] == 0 {
					if u.AsuraCoin >= price[1] {
						done = true
						u.AsuraCoin -= price[1]
					}
				} else if u.Money >= price[0] {
					u.Money -= price[0]
					done = true
				}
				if done {
					database.Client.InsertLootbox(itc.Member.UserID, u.Items, i)
				}
				return u
			}, "Items")
			if done {
				handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.InteractionApplicationCommandCallbackData{
						Content: translation.T("LootboxBuyDone", translation.GetLocale(interaction), lb),
					},
				})
			} else {
				handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.InteractionApplicationCommandCallbackData{
						Content: translation.T("NoMoney", translation.GetLocale(interaction)),
					},
				})
			}

		}, 120)
	}
	ping, _ := handler.Client.HeartbeatLatencies()
	botInfo, _ := handler.Client.Gateway().GetBot()
	shard := disgord.ShardID(itc.GuildID, botInfo.Shards)
	return &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Content: fmt.Sprintf("%dms", ping[shard].Milliseconds()),
		},
	}
}
