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
	"github.com/google/uuid"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "sell",
		Description: translation.T("RunSell", "pt"),
		Run:         runSell,
		Cooldown:    8,
		Category:    handler.Profile,
	})
}

func findItemNameInOptions(options []*disgord.SelectMenuOption, id uuid.UUID) string {
	for _, option := range options {
		if option.Value == id.String() {
			return option.Label
		}
	}
	return ""
}

func genSellOptions(user *entities.User, isRooster bool, isCosmetic bool) (opts []*disgord.SelectMenuOption) {
	if isRooster {
		for _, galo := range user.Galos {
			if !galo.Equip {
				class := rinha.Classes[galo.Type]
				money, asuracoins := rinha.Sell(class.Rarity, galo.Xp, galo.Resets)
				label := fmt.Sprintf("[%d Money]", money)
				if money == 0 {
					label = fmt.Sprintf("[%d Asuracoins]", asuracoins)
				}
				label += " - " + class.Name
				opts = append(opts, &disgord.SelectMenuOption{
					Label:       label,
					Value:       galo.ID.String(),
					Description: "Vender galo " + class.Name,
				})
			}
		}
	} else {
		for _, item := range user.Items {
			var label string
			var itemName string
			var price int
			if item.Type == entities.NormalType && !isCosmetic {
				item := rinha.Items[item.ItemID]
				itemName = item.Name
				price = rinha.SellItem(*item)
			} else if item.Type == entities.CosmeticType && isCosmetic {
				cosmetic := rinha.Cosmetics[item.ItemID]
				itemName = cosmetic.Name
				price = rinha.SellCosmetic(*cosmetic)
			}
			if itemName != "" {
				label = fmt.Sprintf("[%d money] - %s", price, itemName)
				opts = append(opts, &disgord.SelectMenuOption{
					Label:       label,
					Value:       item.ID.String(),
					Description: "Vender item " + itemName,
				})
			}
		}
	}
	if len(opts) == 0 {
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       "Nenhum item",
			Value:       "nil",
			Description: "Nenhum item para vender",
		})
	}
	if len(opts) >= 25 {
		opts = opts[:25]
	}
	return
}

func runSell(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items")
	optsGalos := genSellOptions(&galo, true, false)
	optsItems := genSellOptions(&galo, false, false)
	optCosmetics := genSellOptions(&galo, false, true)
	data := &disgord.CreateInteractionResponseData{
		Embeds: []*disgord.Embed{
			{
				Title: translation.T("SellTitle", translation.GetLocale(itc)),
				Color: 65535,
			},
		},
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Primary,
						Label:    "Galos",
						CustomID: "GalosDisabled",
						Disabled: true,
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:        disgord.MessageComponentButton + 1,
						Style:       disgord.Primary,
						Placeholder: translation.T("SellGaloPlaceholder", translation.GetLocale(itc)),
						CustomID:    "galoSell",
						Options:     optsGalos,
						MaxValues:   1,
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Primary,
						Label:    "Items",
						CustomID: "ItemsDisabled",
						Disabled: true,
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:        disgord.MessageComponentButton + 1,
						Style:       disgord.Primary,
						Placeholder: translation.T("SellItemPlaceholder", translation.GetLocale(itc)),
						CustomID:    "itemSell",
						Options:     optsItems,
						MaxValues:   1,
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:        disgord.MessageComponentButton + 1,
						Style:       disgord.Primary,
						Placeholder: "Selecione os cosmeticos para vender",
						CustomID:    "cosmeticSell",
						Options:     optCosmetics,
						MaxValues:   1,
					},
				},
			},
		},
	}
	optsMap := map[string][]*disgord.SelectMenuOption{
		"itemSell":     optsItems,
		"cosmeticSell": optCosmetics,
		"galoSell":     optsGalos,
	}

	itcID, err := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: data,
	})
	if err != nil {
		return nil
	}
	handler.RegisterHandler(itcID, func(ic *disgord.InteractionCreate) {
		userIC := ic.Member.User
		name := ic.Data.CustomID
		if userIC.ID != itc.Member.UserID {
			return
		}
		if len(ic.Data.Values) == 0 {
			return
		}
		val := ic.Data.Values[0]
		if val == "nil" {
			return
		}
		itemID := uuid.MustParse(val)
		msg := ""
		price := 0
		isAsuraCoins := false
		handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackUpdateMessage,
			Data: data,
		})
		opts := optsMap[name]
		itemName := findItemNameInOptions(opts, itemID)
		utils.ConfirmMessage(ctx, fmt.Sprintf("Deseja mesmo vender?\n\n%s", itemName), itc, itc.Member.UserID, func() {
			database.User.UpdateUser(ctx, userIC.ID, func(u entities.User) entities.User {
				if isInRinha(ctx, userIC) != "" {
					msg = "IsInRinha"
					return u
				}

				if name == "itemSell" || name == "cosmeticSell" {
					item := rinha.GetItemByID(u.Items, itemID)
					if item != nil {
						database.User.RemoveItem(ctx, u.Items, itemID)
						msg = "SellItem"
						if item.Type == entities.NormalType {
							price = rinha.SellItem(*rinha.Items[item.ItemID])
						} else {
							price = rinha.SellCosmetic(*rinha.Cosmetics[item.ItemID])
						}
					}
				}
				if name == "galoSell" {
					galo := rinha.GetGaloByID(u.Galos, itemID)
					if galo != nil && !galo.Equip {
						class := rinha.Classes[galo.Type]
						database.User.RemoveRooster(ctx, itemID)
						msg = "SellGalo"
						money, asuraCoins := rinha.Sell(class.Rarity, galo.Xp, galo.Resets)
						price = money
						if price == 0 {
							price = asuraCoins
							if rinha.IsVip(&u) {
								price++
							}
							isAsuraCoins = true
						}
					}

				}
				if isAsuraCoins {
					msg = "SellGaloAc"
					u.AsuraCoin += price
				} else {
					u.Money += price
				}
				return u
			}, "Items", "Galos")
		})
		if msg != "" {
			handler.Client.Channel(ic.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: translation.T(msg, translation.GetLocale(ic), price),
			})

		}
	}, 120)
	return nil
}
