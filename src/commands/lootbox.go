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

func GenerateBuyOptions(arr []string) (opts []*disgord.SelectMenuOption) {
	for _, name := range arr {
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
								Options:     GenerateBuyOptions(rinha.LootNames[:]),
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
					database.Client.InsertItem(itc.Member.UserID, u.Items, i, entities.LootboxType)
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
	case "open":
		lootbox := rinha.GetLootboxes(&user)
		handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Lootbox open",
						Color:       65535,
						Description: fmt.Sprintf("Lootbox comum: **%d**\nLootbox normal: **%d**\nLootbox rara: **%d**\nLootbox epica: **%d**\nLootbox lendaria: **%d**\nLootbox items: **%d**\nLootbox cosmetica: **%d**", lootbox.Common, lootbox.Normal, lootbox.Rare, lootbox.Epic, lootbox.Legendary, lootbox.Items, lootbox.Cosmetic),
					},
				},
				Components: []*disgord.MessageComponent{
					{
						Type: disgord.MessageComponentActionRow,
						Components: []*disgord.MessageComponent{
							{
								Type:        disgord.MessageComponentButton + 1,
								Options:     GenerateBuyOptions(rinha.GetUserLootboxes(&user)),
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
			image := ""
			color := 65535
			name := ""
			newVal := 0
			winType := ""
			pity := 0
			database.Client.UpdateUser(itc.Member.UserID, func(u entities.User) entities.User {
				lbID, ok := rinha.GetLbID(u.Items, i)
				if !ok {
					return u
				}
				database.Client.RemoveItem(u.Items, lbID)
				newVal, pity = rinha.Open(i, &user)
				u.Pity = pity
				if lb == "cosmetica" {
					database.Client.InsertItem(itc.Member.UserID, u.Items, newVal, entities.CosmeticType)
				} else if lb == "items" {
					database.Client.InsertItem(itc.Member.UserID, u.Items, newVal, entities.NormalType)
				} else {
					database.Client.InsertRooster(&entities.Rooster{
						UserID: itc.Member.UserID,
						Type:   newVal,
					})
				}
				return u
			}, "Items")
			if pity == 0 {
				return
			}
			if lb == "cosmetica" {
				cosmetic := rinha.Cosmetics[newVal]
				image = cosmetic.Value
				winType = "cosmetico"
				name = cosmetic.Name
				color = cosmetic.Rarity.Color()
			} else if lb == "items" {
				item := rinha.Items[newVal]
				winType = "item"
				name = item.Name
			} else {
				galo := rinha.Classes[newVal]
				color = galo.Rarity.Color()
				image = rinha.Sprites[0][newVal-1]
				name = galo.Name
				winType = "galo"
			}
			handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Embeds: []*disgord.Embed{
						{
							Color: color,
							Image: &disgord.EmbedImage{
								URL: image,
							},
							Description: translation.T("LootboxOpen", translation.GetLocale(interaction), map[string]interface{}{
								"name":    name,
								"type":    winType,
								"lootbox": lb,
							}),
							Title: "Lootbox open",
						},
					},
				},
			})
		}, 120)
	}
	return nil
}
