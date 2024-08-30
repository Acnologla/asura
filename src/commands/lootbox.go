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
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "lootbox",
		Description: translation.T("LootboxHelp", "pt"),
		Run:         runLootbox,
		Cooldown:    3,
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

func generateQuantityOptions() (opts []*disgord.SelectMenuOption) {
	for i := 1; i <= 10; i++ {
		opts = append(opts, &disgord.SelectMenuOption{
			Label: fmt.Sprintf("Comprar %d lootboxes", i),
			//		Description: fmt.Sprintf("Quantidade de lootboxes"),
			Value:   strconv.Itoa(i),
			Default: i == 1,
		})
	}
	return
}

func GenerateBuyOptions(arr []string) (opts []*disgord.SelectMenuOption) {
	for i, name := range arr {
		price := rinha.Prices[i]
		if price[0] != -1 {
			opts = append(opts, &disgord.SelectMenuOption{
				Label:       name,
				Description: fmt.Sprintf("%s lootbox", name),
				Value:       name,
			})
		}
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
						Description: fmt.Sprintf("Money: **%d**\nAsuraCoins: **%d**\nPity: **%d%%**\n\n<:lt_comum:1271148114102714479> Lootbox comum: **%d**\n<:lt_normal:1271148156414984286> Lootbox normal: **%d**\n<:lt_rara:1271148187725594644> Lootbox rara: **%d**\n<:lt_epica:1271148219623145493> Lootbox epica: **%d**\n<:lt_lendaria:1271148244767998115> Lootbox lendaria: **%d**\n<:lt_itens:1271148268969267201> Lootbox items: **%d**\n<:lt_cosmetica:1271148294952980614> Lootbox cosmetica: **%d**\n<:lt_mitica:1271148323344093246> Lootbox mitica: **%d**\n<:lt_itens_mitica:1271148347700412578>  Lootbox items mitica: **%d**\n%s\n**Entre no meu [Servidor](https://discord.gg/tdVWQGV) para ganhar um bonus no daily**", user.Money, user.AsuraCoin, user.Pity*rinha.PityMultiplier, lootbox.Common, lootbox.Normal, lootbox.Rare, lootbox.Epic, lootbox.Legendary, lootbox.Items, lootbox.Cosmetic, lootbox.Mythic, lootbox.ItemsMythic, text),
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
								Options:     generateQuantityOptions(),
								CustomID:    "quantity",
								MaxValues:   1,
								Placeholder: "Selecione a quantidade",
							},
						},
					},
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
		quantity := 1
		lastValue := ""
		handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
			if interaction.Member.UserID != itc.Member.UserID {
				return
			}
			if interaction.Data.CustomID == "quantity" {
				quantity, _ = strconv.Atoi(interaction.Data.Values[0])
				handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: fmt.Sprintf("Quantidade selecionada: %d", quantity),
					},
				})
				return
			}
			if len(interaction.Data.Values) == 0 {
				interaction.Data.Values = []string{lastValue}
				if lastValue == "" {
					return
				}
			}
			lastValue = interaction.Data.Values[0]
			lb := interaction.Data.Values[0]
			i := rinha.GetLbIndex(lb)
			price := rinha.Prices[i]
			done := false
			database.User.UpdateUser(ctx, itc.Member.UserID, func(u entities.User) entities.User {
				if price[0] == 0 {
					if u.AsuraCoin >= price[1]*quantity {
						done = true
						u.AsuraCoin -= price[1] * quantity
					}
				} else if u.Money >= price[0]*quantity {
					u.Money -= price[0] * quantity
					done = true
				}
				if done {
					database.User.InsertItemQuantity(ctx, itc.Member.UserID, u.Items, i, entities.LootboxType, quantity)
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
		err := handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Lootbox open",
						Color:       65535,
						Description: fmt.Sprintf("Use **/lootbox buy** para comprar lootboxes\n\n<:lt_comum:1271148114102714479> Lootbox comum: **%d**\n<:lt_normal:1271148156414984286> Lootbox normal: **%d**\n<:lt_rara:1271148187725594644> Lootbox rara: **%d**\n<:lt_epica:1271148219623145493> Lootbox epica: **%d**\n<:lt_lendaria:1271148244767998115> Lootbox lendaria: **%d**\n<:lt_itens:1271148268969267201> Lootbox items: **%d**\n<:lt_cosmetica:1271148294952980614> Lootbox cosmetica: **%d**\n<:lt_mitica:1271148323344093246> Lootbox mistica: **%d**\n<:lt_itens_mitica:1271148347700412578>  Lootbox items mitica: **%d**", lootbox.Common, lootbox.Normal, lootbox.Rare, lootbox.Epic, lootbox.Legendary, lootbox.Items, lootbox.Cosmetic, lootbox.Mythic, lootbox.ItemsMythic),
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
		if err != nil {
			return nil
		}
		lastValue := ""
		handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
			if interaction.Member.UserID != itc.Member.UserID {
				return
			}
			if len(interaction.Data.Values) == 0 {
				interaction.Data.Values = []string{lastValue}
				if lastValue == "" {
					return
				}
			}
			lastValue = interaction.Data.Values[0]
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
				if lb != "cosmetica" && lb != "items" && lb != "items mistica" {
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
				} else if lb == "items" || lb == "items mistica" {
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
			} else if lb == "items" || lb == "items mistica" {
				item := rinha.Items[newVal]
				if item.Level >= 4 {
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

			embed := disgord.Embed{
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
			}

			if winType == "galo" {
				rand := rinha.GetRand()
				class := rinha.Classes[rand]
				openEmbed := disgord.Embed{
					Color: class.Rarity.Color(),
					Image: &disgord.EmbedImage{
						URL: rinha.Sprites[0][rand-1],
					},
					Title: "Abrindo lootbox",
				}
				err := handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Embeds: []*disgord.Embed{
							&openEmbed,
						},
					},
				})
				if err != nil {
					return
				}
				for i := 0; i < 4; i++ {
					time.Sleep(3 * time.Second)
					rand := rinha.GetRand()
					openEmbed.Color = rinha.Classes[rand].Rarity.Color()
					openEmbed.Image = &disgord.EmbedImage{
						URL: rinha.Sprites[0][rand-1],
					}
					handler.Client.EditInteractionResponse(ctx, interaction, &disgord.UpdateMessage{
						Embeds: &([]*disgord.Embed{&openEmbed}),
					})
				}
				time.Sleep(3 * time.Second)
				handler.Client.EditInteractionResponse(ctx, interaction, &disgord.UpdateMessage{
					Embeds: &([]*disgord.Embed{&embed}),
				})
			} else {
				handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Embeds: []*disgord.Embed{
							&embed,
						},
					},
				})
			}
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
