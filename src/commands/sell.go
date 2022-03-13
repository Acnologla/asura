package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
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
		Cooldown:    15,
		Category:    handler.Profile,
	})
}

func genSellOptions(user *entities.User, isRooster bool) (opts []*disgord.SelectMenuOption) {
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
			if item.Type == entities.NormalType {
				item := rinha.Items[item.ItemID]
				itemName = item.Name
				price = rinha.SellItem(*item)
			} else if item.Type == entities.CosmeticType {
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
	return
}

func runSell(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(itc.Member.UserID, "Galos", "Items")
	optsGalos := genSellOptions(&galo, true)
	optsItems := genSellOptions(&galo, false)
	handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
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
			},
		},
	})
	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
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
		database.User.UpdateUser(userIC.ID, func(u entities.User) entities.User {
			if isInRinha(userIC) != "" {
				msg = "IsInRinha"
				return u
			}
			if name == "itemSell" {
				item := rinha.GetItemByID(u.Items, itemID)
				if item != nil {
					database.User.RemoveItem(u.Items, itemID)
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
					database.User.RemoveRooster(itemID)
					msg = "SellGalo"
					money, asuraCoins := rinha.Sell(class.Rarity, galo.Xp, galo.Resets)
					price = money
					if price == 0 {
						price = asuraCoins
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
		if msg != "" {
			handler.Client.SendInteractionResponse(context.Background(), ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T(msg, translation.GetLocale(ic), price),
				},
			})
		}
	}, 120)
	return nil
}
