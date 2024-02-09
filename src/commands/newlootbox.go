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
		Name:        "newlootbox",
		Description: translation.T("LootboxHelp", "pt"),
		Run:         runNewLootbox,
		Cooldown:    8,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "view",
				Description: translation.T("LootboxViewHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "buy",
				Description: "compre lootboxes",
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeString,
					Name:        "type",
					Required:    true,
					Description: "tipo da lootbox (lendaria, items, epica, rara, cosmetica, normal, comum)",
					Choices: []*disgord.ApplicationCommandOptionChoice{
						{
							Name:  "comum",
							Value: "comum",
						},
						{
							Name:  "normal",
							Value: "normal",
						},
						{
							Name:  "cosmetica",
							Value: "cosmetica",
						},
						{
							Name:  "rara",
							Value: "rara",
						},
						{
							Name:  "epica",
							Value: "epica",
						},
						{
							Name:  "items",
							Value: "items",
						},
						{
							Name:  "lendaria",
							Value: "lendaria",
						},
					},
				},
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeNumber,
						Name:        "quantity",
						Required:    true,
						Description: "quantidade para comprar",
						MinValue:    1,
						MaxValue:    100,
					}),
			}),
	})
}

func runNewLootbox(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
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
		quantity := int(itc.Data.Options[0].Options[1].Value.(float64))
		lb := itc.Data.Options[0].Options[0].Value.(string)
		i := rinha.GetLbIndex(lb)
		prices := rinha.Prices[i]
		done := false

		database.User.UpdateUser(ctx, itc.Member.UserID, func(u entities.User) entities.User {
			price := prices[0] * quantity
			if prices[0] == 0 {
				if u.AsuraCoin >= prices[1]*quantity {
					done = true
					price = prices[1] * quantity
					u.AsuraCoin -= price
				}
			} else if u.Money >= price {
				done = true
				u.Money -= price
			}

			if done {
				database.User.InsertManyItems(ctx, itc.Member.UserID, u.Items, i, entities.LootboxType, quantity)
			}

			return u
		}, "Items")

		if done {
			handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("LootboxBuyDone", translation.GetLocale(itc), map[string]interface{}{
						"lb": lb,
						"q":  quantity,
					}),
				},
			})
		} else {
			handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("NoMoney", translation.GetLocale(itc)),
				},
			})
		}
	}

	return nil
}
