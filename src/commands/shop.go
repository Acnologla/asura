package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/translation"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

var shop = rinha.GenerateShop()
var lastTime = time.Now()

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "shop",
		Description: "Comprar coisas",
		Run:         runShop,
		Cooldown:    5,
		Category:    handler.Profile,
	})
	go func() {
		for range time.NewTicker(time.Hour * 8).C {
			shop = rinha.GenerateShop()
			lastTime = time.Now()
		}
	}()
}

func getEmbedFromShopItem(item *entities.ShopItem, itc *disgord.InteractionCreate) *disgord.Embed {
	moneyPrice, asuraCoinPrice := item.Price()
	price := utils.Ternary(moneyPrice > 0, fmt.Sprintf("%d Dinheiro", moneyPrice), fmt.Sprintf("%d AsuraCoins", asuraCoinPrice))
	userID := itc.Member.UserID
	embed := &disgord.Embed{
		Title: rinha.GetShopItemName(item),
		Color: utils.Ternary(!item.CanBuy(userID), 9868951, 65535),
		Footer: &disgord.EmbedFooter{
			Text: translation.T("ShopFooter", translation.GetLocale(itc)),
		},
	}
	if item.Discount == 1 {
		embed.Description = translation.T("ShopItemDescription", translation.GetLocale(itc), map[string]interface{}{
			"price":  price,
			"rarity": rinha.GetShopItemRarity(item),
		})
	} else {
		moneyPO, asuraCoinPO := item.OriginalPrice()
		originalPrice := utils.Ternary(moneyPO > 0, fmt.Sprintf("%d Dinheiro", moneyPO), fmt.Sprintf("%d AsuraCoins", asuraCoinPO))
		embed.Description = translation.T("ShopItemDescriptionDiscount", translation.GetLocale(itc), map[string]interface{}{
			"price":         price,
			"originalPrice": originalPrice,
			"rarity":        rinha.GetShopItemRarity(item),
		})
	}

	if item.Type == entities.Roosters {
		embed.Thumbnail = &disgord.EmbedThumbnail{
			URL: rinha.Sprites[0][item.Value-1],
		}
	}
	if item.Type == entities.Cosmetics {
		embed.Thumbnail = &disgord.EmbedThumbnail{
			URL: rinha.Cosmetics[item.Value].Value,
		}
	}

	return embed
}

func runShop(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID)
	targetTime := lastTime.Add(8 * time.Hour)

	now := time.Now()
	remaining := targetTime.Sub(now)
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	currentItem := -1
	defaultEmbed := &disgord.Embed{
		Title: "Shop",
		Color: 65535,
		Description: translation.T("ShopDefaultMessage", translation.GetLocale(itc), map[string]interface{}{
			"minutes": minutes,
			"hours":   hours,
		}),
		Footer: &disgord.EmbedFooter{
			Text: translation.T("ShopFooter", translation.GetLocale(itc)),
		},
	}

	component := disgord.MessageComponent{
		Type: disgord.MessageComponentActionRow,
		Components: []*disgord.MessageComponent{
			{
				Type:        disgord.MessageComponentSelectMenu,
				MaxValues:   1,
				Placeholder: "Selecione o que comprar",
			},
		},
	}

	for i, item := range shop {
		component.Components[0].Options = append(component.Components[0].Options, &disgord.SelectMenuOption{
			Label:       rinha.GetShopItemName(item),
			Description: rinha.GetShopItemRarity(item),
			Value:       fmt.Sprintf("%d", i),
		})
	}

	itcID, err := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds:     []*disgord.Embed{defaultEmbed},
			Components: []*disgord.MessageComponent{&component},
		},
	})

	if err != nil {
		return nil
	}

	handler.RegisterHandler(itcID, func(ic *disgord.InteractionCreate) {
		if ic.Member.UserID == user.ID {
			customID := ic.Data.CustomID
			if customID == "buy" {
				if currentItem == -1 {
					return
				}
				item := shop[currentItem]

				msg := ""
				database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
					if item.Type == entities.Roosters && len(u.Galos) >= 10 {
						msg = "MaxGalos"
						return u
					}
					if !item.CanBuy(u.ID) {
						msg = "ArleadyBought"
						return u
					}
					money, asuraCoin := item.Price()
					if u.Money < money || u.AsuraCoin < asuraCoin {
						msg = "NotEnoughMoney"
						return u
					}
					u.Money -= money
					u.AsuraCoin -= asuraCoin
					item.Buy(u.ID)
					switch item.Type {
					case entities.Shards:
						database.User.InsertItem(ctx, user.ID, u.Items, item.Value, entities.ShardType)
					case entities.Cosmetics:
						database.User.InsertItem(ctx, user.ID, u.Items, item.Value, entities.CosmeticType)
					case entities.Roosters:
						database.User.InsertRooster(ctx, &entities.Rooster{
							UserID: u.ID,
							Type:   item.Value,
						})
					case entities.Items:
						database.User.InsertItem(ctx, user.ID, u.Items, item.Value, entities.NormalType)
					case entities.AsuraCoin:
						u.AsuraCoin += item.Value
					case entities.Xp:
						database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
							r.Xp += item.Value / (r.Resets + 1)
							return r
						})
					}
					msg = "BoughtShop"
					return u
				}, "Galos", "Items")
				handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: translation.T(msg, translation.GetLocale(ic), rinha.GetShopItemName(item)),
					},
				})
			} else {
				if len(ic.Data.Values) == 0 {
					return
				}
				val := ic.Data.Values[0]
				v, err := strconv.Atoi(val)
				if err != nil {
					return
				}
				currentItem = v
				embed := getEmbedFromShopItem(shop[currentItem], ic)
				handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackUpdateMessage,
					Data: &disgord.CreateInteractionResponseData{
						Embeds: []*disgord.Embed{embed},
						Components: []*disgord.MessageComponent{&component, {
							Type: disgord.MessageComponentActionRow,
							Components: []*disgord.MessageComponent{
								{
									Type:     disgord.MessageComponentButton,
									Style:    disgord.Success,
									Label:    "Comprar",
									CustomID: "buy",
								},
							},
						}},
					},
				})
			}
		}
	}, 200)

	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title:       "Money",
					Color:       65535,
					Description: fmt.Sprintf("Dinheiro: **%d**\nAsuraCoin: **%d**\nUserXP: **%d**\nTreinos: **%d/%d**\n\nUse `/lootbox view` para visualizar ou comprar lootboxes\nUse `/givemoney` para doar dinheiro\nUse `/daily` para pegar o bonus diario\nUse `/trial status` para fortificar seu galo\n\n**Entre no meu [Servidor](https://discord.gg/CfkBZyVsd7) para ganhar um bonus no daily**\n**[Comprar Moedas e XP](https://acnologla.github.io/asura-site/donate)**", user.Money, user.AsuraCoin, user.UserXp, user.TrainLimit, rinha.CalcLimit(&user)),
				},
			},
		},
	}
}
