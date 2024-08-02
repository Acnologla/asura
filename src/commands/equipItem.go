package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "item",
		Description: translation.T("RunEquipItem", "pt"),
		Run:         runEquipItem,
		Cooldown:    15,
		Category:    handler.Profile,
		Aliases:     []string{"equipitem"},
	})
}

func getEquipItem(items []*entities.Item, itemType entities.ItemType) *entities.Item {
	for _, item := range items {
		if item.Type == itemType && item.Equip {
			if itemType == entities.CosmeticType {
				_item := rinha.Cosmetics[item.ItemID]
				if _item.Type == rinha.Background {
					return item
				}
			} else {
				return item
			}
		}
	}
	return nil
}

func genEquipItemsOptions(user *entities.User, itemType entities.ItemType) (opts []*disgord.SelectMenuOption) {
	for _, item := range user.Items {
		if item.Type == itemType {
			if itemType == entities.CosmeticType {
				cosmetic := rinha.Cosmetics[item.ItemID]
				if cosmetic.Type == rinha.Background {
					opts = append(opts, &disgord.SelectMenuOption{
						Label:       cosmetic.Name,
						Value:       item.ID.String(),
						Description: "Equipar background " + cosmetic.Name,
						Default:     item.Equip,
					})
				}
			} else if itemType == entities.NormalType {
				_item := rinha.Items[item.ItemID]
				opts = append(opts, &disgord.SelectMenuOption{
					Label:       _item.Name,
					Value:       item.ID.String(),
					Description: rinha.ItemToString(_item),
					Default:     item.Equip,
				})
			}
		}

	}
	if len(opts) == 0 {
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       "Nenhum item",
			Value:       "nil",
			Description: "Nenhum item para equipar",
		})
	}
	if len(opts) >= 25 {
		opts = opts[:25]
	}
	return
}

func runEquipItem(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(ctx, itc.Member.UserID, "Items")
	optsItems := genEquipItemsOptions(&galo, entities.NormalType)
	optsBackground := genEquipItemsOptions(&galo, entities.CosmeticType)

	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title: "Equip items",
					Color: 65535,
				},
			},
			Components: []*disgord.MessageComponent{
				{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{
						{
							Type:        disgord.MessageComponentButton + 1,
							Style:       disgord.Primary,
							Placeholder: translation.T("EquipItem", translation.GetLocale(itc)),
							CustomID:    "equipItems",
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
							Placeholder: translation.T("EquipBackground", translation.GetLocale(itc)),
							CustomID:    "equipCosmetics",
							Options:     optsBackground,
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
		database.User.UpdateUser(ctx, userIC.ID, func(u entities.User) entities.User {
			if isInRinha(ctx, userIC) != "" {
				msg = "IsInRinha"
				return u
			}
			t := entities.NormalType
			if name != "equipItems" {
				t = entities.CosmeticType
			}
			equipedItem := getEquipItem(u.Items, t)
			if equipedItem != nil {
				database.User.UpdateItem(ctx, &u, equipedItem.ID, func(i entities.Item) entities.Item {
					i.Equip = false
					return i
				})
			}
			database.User.UpdateItem(ctx, &u, itemID, func(i entities.Item) entities.Item {
				i.Equip = true
				return i
			})
			msg = "EquipItemSuccess"
			return u
		}, "Items")
		if msg != "" {
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T(msg, translation.GetLocale(ic)),
				},
			})
		}
	}, 120)
	return nil
}
