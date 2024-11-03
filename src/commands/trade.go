package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"
	"strings"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "trade",
		Description: translation.T("TradeHelp", "pt"),
		Run:         runTrade,
		Cooldown:    15,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.OptionTypeUser,
			Required:    true,
			Name:        "user",
			Description: "user trade",
		}),
		Category: handler.Rinha,
	})
}

type tradeItemType int

const (
	roosterTradeType tradeItemType = iota
	itemTradeType
)

type ItemTrade struct {
	Type tradeItemType
	ID   string
	Name string
}

func isInItemTrade(id string, arr []*ItemTrade) bool {
	for _, item := range arr {
		if item.ID == id {
			return true
		}
	}
	return false
}

func removeItemFromItemTrade(id string, arr []*ItemTrade) []*ItemTrade {
	for i, item := range arr {
		if item.ID == id {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr

}

func findOptionById(id string, options []*disgord.SelectMenuOption) *disgord.SelectMenuOption {
	for _, opt := range options {
		if strings.Contains(opt.Value, id) {
			return opt
		}
	}
	return nil
}

func editEmbed(authorItems, userItems []*ItemTrade, authorUsername, userUsername string, minLevel int, extraMsg string) *disgord.Embed {
	authorVal := ""
	userVal := ""
	for _, item := range authorItems {
		if item.Type == roosterTradeType {
			authorVal += item.Name + "\n"
		} else {
			authorVal += item.Name + "\n"
		}
	}
	for _, item := range userItems {
		if item.Type == roosterTradeType {
			userVal += item.Name + "\n"
		} else {
			userVal += item.Name + "\n"
		}
	}
	if authorVal == "" {
		authorVal = "0 Items"
	}
	if userVal == "" {
		userVal = "0 Items"
	}
	return &disgord.Embed{
		Title: "Trade",
		Color: 65535,
		Description: translation.T("TradeDesc", "pt", map[string]interface{}{
			"minLvl":   minLevel,
			"extraMsg": extraMsg,
		}),
		Fields: []*disgord.EmbedField{{
			Name:  authorUsername + " Items",
			Value: authorVal,
		}, {
			Name:  userUsername + " Items",
			Value: userVal,
		}},
	}
}

func itemsToOptions(user *entities.User, minLevel *int) (opts []*disgord.SelectMenuOption) {
	for _, rooster := range user.Galos {
		if !rooster.Equip {
			galo := rinha.Classes[rooster.Type]
			lvl := int(galo.Rarity) * 75
			if galo.Rarity >= rinha.Legendary {
				lvl += 550
			}
			if lvl > *minLevel {
				*minLevel = lvl
			}
			opts = append(opts, &disgord.SelectMenuOption{
				Label:       "Galo " + galo.Name,
				Value:       fmt.Sprintf("rooster|%s|%d", rooster.ID.String(), lvl),
				Description: "Adicionar ou remover da troca",
			})
		}
	}
	for _, item := range user.Items {
		if !item.Equip {
			lvl := 0
			name := ""
			switch item.Type {
			case entities.NormalType:
				_item := rinha.Items[item.ItemID]
				if 6 > _item.Level {
					lvl = int(_item.Level) * 75
					if _item.Level >= 3 {
						lvl += 550
					}
					name = _item.Name
				}
			case entities.CosmeticType:
				_item := rinha.Cosmetics[item.ItemID]
				lvl = int(_item.Rarity) * 40
				name = _item.Name

			}
			if lvl > *minLevel {
				*minLevel = lvl
			}
			if name != "" {
				opts = append(opts, &disgord.SelectMenuOption{
					Label:       "Item " + name,
					Description: "Adicionar ou remover da troca",
					Value:       fmt.Sprintf("item|%s|%d", item.ID.String(), lvl),
				})
			}

		}
	}
	if len(opts) == 0 {
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       "Nenhum item",
			Description: "Nenhum item para troca",
			Value:       "nil",
		})
	}
	if len(opts) >= 25 {
		opts = opts[:25]
	}
	return
}

func runTrade(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := utils.GetUser(itc, 0)
	if user.Bot || user.ID == itc.Member.UserID {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "invalid user",
			},
		}
	}
	authorUser := itc.Member.User

	utils.Confirm(ctx, translation.T("TradeMsg", translation.GetLocale(itc), map[string]interface{}{
		"mentionedUsername": user.Username,
		"username":          authorUser.Username,
	}), itc, user.ID, func() {
		ch := handler.Client.Channel(itc.ChannelID)
		userRinha := isInRinha(ctx, user)
		if userRinha != "" {
			ch.CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(user.Username, userRinha).Data.Content,
			})
			return
		}
		authorRinha := isInRinha(ctx, authorUser)
		if authorRinha != "" {
			ch.CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(authorUser.Username, userRinha).Data.Content,
			})
			return
		}
		lockBattle(ctx, itc.Member.UserID, user.ID, authorUser.Username, user.Username)
		defer unlockBattle(ctx, itc.Member.UserID, user.ID)
		authorGalo := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items")
		userGalo := database.User.GetUser(ctx, user.ID, "Galos", "Items")
		minLevel := 0
		optsUser := itemsToOptions(&userGalo, &minLevel)
		optsAuthor := itemsToOptions(&authorGalo, &minLevel)
		minLevel = 0
		components := []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Disabled: true,
						Style:    disgord.Primary,
						Label:    authorUser.Username,
						CustomID: authorUser.Username + "Disabled",
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:        disgord.MessageComponentButton + 1,
						Options:     optsAuthor,
						Placeholder: "Select items",
						MaxValues:   1,
						CustomID:    itc.Member.UserID.String(),
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Disabled: true,
						Style:    disgord.Primary,
						Label:    user.Username,
						CustomID: user.Username + "Disabled",
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:        disgord.MessageComponentButton + 1,
						Options:     optsUser,
						Placeholder: "Select items",
						CustomID:    user.ID.String(),
						MaxValues:   1,
					},
				},
			},
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Success,
						Label:    "Confirm",
						CustomID: "done",
					},
					{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Danger,
						Label:    "Reject",
						CustomID: "reject",
					},
				},
			},
		}
		msg, err := ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{{
				Title: "Trade",
				Color: 65535,
				Fields: []*disgord.EmbedField{{
					Name:  authorUser.Username + " Items",
					Value: "0 Items",
				}, {
					Name:  user.Username + " Items",
					Value: "0 Items",
				}},
			}},
			Components: components,
		})
		if err == nil {
			itemsAuthor := []*ItemTrade{}
			userDone := false
			authorDone := false
			itemsUser := []*ItemTrade{}
			handler.RegisterHandler(msg.ID, func(interaction *disgord.InteractionCreate) {
				userIC := interaction.Member.User
				name := interaction.Data.CustomID
				if userIC.ID != user.ID && userIC.ID != authorUser.ID {
					return
				}
				switch name {
				case "nil":
				case "done":
					if userIC.ID == user.ID {
						userDone = !userDone
					}
					if userIC.ID == authorUser.ID {
						authorDone = !authorDone
					}
					if userDone && authorDone {
						if minLevel > userGalo.UserXp || minLevel > authorGalo.UserXp {
							handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
								Type: disgord.InteractionCallbackUpdateMessage,
								Data: &disgord.CreateInteractionResponseData{
									Embeds:     []*disgord.Embed{editEmbed(itemsAuthor, itemsUser, authorUser.Username, user.Username, minLevel, translation.T("UserMinLevelTrade", translation.GetLocale(itc), minLevel))},
									Components: components,
								},
							})
							return
						}
						database.User.UpdateUser(ctx, authorUser.ID, func(a entities.User) entities.User {
							database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
								for _, item := range itemsAuthor {
									itemID := uuid.MustParse(item.ID)
									if item.Type == roosterTradeType {
										galo := rinha.GetGaloByID(a.Galos, itemID)
										if galo != nil {
											database.User.RemoveRooster(ctx, itemID)
											database.User.InsertRooster(ctx, &entities.Rooster{
												UserID: user.ID,
												Type:   galo.Type,
											})
										}
									} else {
										item := rinha.GetItemByID(a.Items, itemID)
										if item != nil {
											database.User.RemoveItem(ctx, a.Items, itemID)
											database.User.InsertItem(ctx, u.ID, u.Items, item.ItemID, item.Type)
										}
									}
								}
								for _, item := range itemsUser {
									itemID := uuid.MustParse(item.ID)
									if item.Type == roosterTradeType {
										galo := rinha.GetGaloByID(u.Galos, itemID)
										if galo != nil {
											database.User.RemoveRooster(ctx, itemID)
											database.User.InsertRooster(ctx, &entities.Rooster{
												UserID: authorUser.ID,
												Type:   galo.Type,
											})
										}
									} else {
										item := rinha.GetItemByID(u.Items, itemID)
										if item != nil {
											database.User.RemoveItem(ctx, u.Items, itemID)
											database.User.InsertItem(ctx, u.ID, a.Items, item.ItemID, item.Type)
										}
									}
								}
								return u
							}, "Galos", "Items")
							return a
						}, "Galos", "Items")
						handler.DeleteHandler(msg.ID)
						handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
							Type: disgord.InteractionCallbackChannelMessageWithSource,
							Data: &disgord.CreateInteractionResponseData{
								Content: translation.T("TradeDone", translation.GetLocale(itc)),
							},
						})
						handler.Client.Channel(msg.ChannelID).Message(msg.ID).Delete()

					} else {
						acceptUsername := authorUser.Username
						if userDone {
							acceptUsername = user.Username
						}
						handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
							Type: disgord.InteractionCallbackUpdateMessage,
							Data: &disgord.CreateInteractionResponseData{
								Embeds:     []*disgord.Embed{editEmbed(itemsAuthor, itemsUser, authorUser.Username, user.Username, minLevel, translation.T("UserAcceptTrade", translation.GetLocale(itc), acceptUsername))},
								Components: components,
							},
						})
					}
				case "reject":
					handler.Client.Channel(interaction.ChannelID).Message(interaction.Message.ID).Delete()
					handler.DeleteHandler(msg.ID)
				default:
					if len(interaction.Data.Values) == 0 {
						return
					}
					item := interaction.Data.Values[0]
					_id, _ := strconv.ParseUint(name, 10, 64)
					id := disgord.Snowflake(_id)
					split := strings.Split(item, "|")
					if len(split) != 3 {
						return
					}
					_itemType := split[0]
					itemID := split[1]
					lvl, _ := strconv.Atoi(split[2])
					var itemType tradeItemType
					if _itemType == "rooster" {
						itemType = roosterTradeType
					} else {
						itemType = itemTradeType
					}
					if userIC.ID == id {
						if id == user.ID {
							item := findOptionById(itemID, optsUser)
							if item != nil {
								if isInItemTrade(itemID, itemsUser) {
									itemsUser = removeItemFromItemTrade(itemID, itemsUser)
								} else {
									if lvl > minLevel {
										minLevel = lvl
									}
									itemsUser = append(itemsUser, &ItemTrade{
										Type: itemType,
										Name: item.Label,
										ID:   itemID,
									})
								}
							}
						} else {
							item := findOptionById(itemID, optsAuthor)
							if item != nil {
								if isInItemTrade(itemID, itemsAuthor) {
									itemsAuthor = removeItemFromItemTrade(itemID, itemsAuthor)
								} else {
									if lvl > minLevel {
										minLevel = lvl
									}
									itemsAuthor = append(itemsAuthor, &ItemTrade{
										Type: itemType,
										Name: item.Label,
										ID:   itemID,
									})
								}
							}

						}
						handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
							Type: disgord.InteractionCallbackUpdateMessage,
							Data: &disgord.CreateInteractionResponseData{
								Embeds:     []*disgord.Embed{editEmbed(itemsAuthor, itemsUser, authorUser.Username, user.Username, minLevel, "")},
								Components: components,
							},
						})
					}
				}

			}, 220)
		}
	})
	return nil
}
