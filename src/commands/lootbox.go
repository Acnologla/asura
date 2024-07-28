package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/telemetry"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "lootbox",
		Description: translation.T("LootboxHelp", "pt"),
		Run:         runLootbox,
		Cooldown:    8,
		Category:    handler.Profile,
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

func runLootbox(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	command := itc.Data.Options[0].Name
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items")
	switch command {
	case "view":
		text := translation.T("LootboxView", "pt", rinha.VipMessage(&user))
		lootbox := rinha.GetLootboxes(&user)
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Lootbox",
						Color:       65535,
						Description: fmt.Sprintf("Money: **%d**\nAsuraCoins: **%d**\nPity: **%d%%**\nTreinos diarios: **%d/%d** \n\nLootbox comum: **%d**\nLootbox normal: **%d**\nLootbox rara: **%d**\nLootbox epica: **%d**\nLootbox lendaria: **%d**\nLootbox items: **%d**\nLootbox cosmetica: **%d**\n\n%s", user.Money, user.AsuraCoin, user.Pity*rinha.PityMultiplier, user.TrainLimit, rinha.CalcLimit(&user), lootbox.Common, lootbox.Normal, lootbox.Rare, lootbox.Epic, lootbox.Legendary, lootbox.Items, lootbox.Cosmetic, text),
					},
				},
			},
		}
	case "buy":
		handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
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
		handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
			if interaction.Member.UserID != itc.Member.UserID {
				return
			}
			if len(interaction.Data.Values) == 0 {
				return
			}
			lb := interaction.Data.Values[0]
			i := rinha.GetLbIndex(lb)
			price := rinha.Prices[i]
			done := false
			database.User.UpdateUser(ctx, itc.Member.UserID, func(u entities.User) entities.User {
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
					database.User.InsertItem(ctx, itc.Member.UserID, u.Items, i, entities.LootboxType)
				}
				return u
			}, "Items")
			if done {
				handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: translation.T("LootboxBuyDone", translation.GetLocale(interaction), lb),
					},
				})
			} else {
				handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: translation.T("NoMoney", translation.GetLocale(interaction)),
					},
				})
			}

		}, 120)
	case "open":
		lootbox := rinha.GetLootboxes(&user)
		loots := rinha.GetUserLootboxes(&user)
		if len(loots) == 0 {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("NoLootbox", translation.GetLocale(itc)),
				},
			}
		}
		handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
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
								Options:     GenerateBuyOptions(loots),
								CustomID:    "type",
								MaxValues:   1,
								Placeholder: "Select lootbox",
							},
						},
					},
				},
			},
		})
		handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
			if interaction.Member.UserID != itc.Member.UserID {
				return
			}
			if len(interaction.Data.Values) == 0 {
				return
			}
			lb := interaction.Data.Values[0]
			i := rinha.GetLbIndex(lb)
			image := ""
			name := ""
			newVal := -1
			winType := ""
			pity := 0
			extraMsg := ""
			var rarity rinha.Rarity
			database.User.UpdateUser(ctx, itc.Member.UserID, func(u entities.User) entities.User {
				lbID, ok := rinha.GetLbID(u.Items, i)
				if !ok {
					return u
				}
				if lb != "cosmetica" && lb != "items" {
					if len(u.Galos) >= 10 {
						handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
							Type: disgord.InteractionCallbackChannelMessageWithSource,
							Data: &disgord.CreateInteractionResponseData{
								Content: "Voce ja chegou no limite de galos (10)",
							},
						})
						return u

					}
				}
				database.User.RemoveItem(ctx, u.Items, lbID)
				newVal, pity = rinha.Open(i, &u)
				u.Pity = pity
				if lb == "cosmetica" {
					database.User.InsertItem(ctx, itc.Member.UserID, u.Items, newVal, entities.CosmeticType)
				} else if lb == "items" {
					database.User.InsertItem(ctx, itc.Member.UserID, u.Items, newVal, entities.NormalType)
				} else {
					if !rinha.HaveRooster(u.Galos, newVal) {
						database.User.InsertRooster(ctx, &entities.Rooster{
							UserID: itc.Member.UserID,
							Type:   newVal,
						})
					} else if rinha.Classes[newVal].Rarity != rinha.Mythic {
						gal := rinha.Classes[newVal]
						money, _ := rinha.Sell(gal.Rarity, 0, 0)
						u.Money += money
						extraMsg = translation.T("SellRepeated", translation.GetLocale(interaction), money)
					}

				}
				return u
			}, "Items", "Galos")
			if newVal == -1 {
				return
			}
			if lb == "cosmetica" {
				cosmetic := rinha.Cosmetics[newVal]
				image = cosmetic.Value
				winType = "cosmetico"
				name = cosmetic.Name
				rarity = cosmetic.Rarity
			} else if lb == "items" {
				item := rinha.Items[newVal]
				if item.Level == 4 {
					rarity = rinha.Mythic
				} else {
					rarity = rinha.Legendary
				}
				winType = "item"
				name = item.Name
			} else {
				galo := rinha.Classes[newVal]
				image = rinha.Sprites[0][newVal-1]
				name = galo.Name
				rarity = galo.Rarity
				winType = "galo"
			}
			handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Embeds: []*disgord.Embed{
						{
							Color: rarity.Color(),
							Image: &disgord.EmbedImage{
								URL: image,
							},
							Description: translation.T("LootboxOpen", translation.GetLocale(interaction), map[string]interface{}{
								"name":     name,
								"type":     winType,
								"lootbox":  lb,
								"extraMsg": extraMsg,
							}),
							Title: fmt.Sprintf("Lootbox open (%s)", rarity.String()),
						},
					},
				},
			})
			author := itc.Member.User
			tag := author.Username + "#" + author.Discriminator.String()
			telemetry.Debug(fmt.Sprintf("%s wins %s", tag, name), map[string]string{
				"value":    name,
				"user":     strconv.FormatUint(uint64(author.ID), 10),
				"rarity":   rarity.String(),
				"lootType": lb,
				"pity":     strconv.Itoa(pity),
			})
		}, 120)
	}
	return nil
}
