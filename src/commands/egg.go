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

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "egg",
		Description: translation.T("EggHelp", "pt"),
		Run:         runEgg,
		Cooldown:    5,
		Category:    handler.Profile,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "view",
				Description: translation.T("EggViewHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "feed",
				Description: translation.T("EggFeedHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "hatch",
				Description: translation.T("EggHatchHelp", "pt"),
			},
		),
	})
}

func generateShardOptions(items []*entities.Item) []*disgord.SelectMenuOption {
	var opts []*disgord.SelectMenuOption
	for _, item := range items {
		if item.Type == entities.ShardType {
			rarity := rinha.Rarity(item.ItemID)
			opts = append(opts, &disgord.SelectMenuOption{
				Description: fmt.Sprintf("[%d] Shard %s", item.Quantity, rarity.String()),
				Label:       strconv.Itoa(rinha.ShardToPrice(rarity)) + "XP",
				Value:       item.ID.String(),
			})
		}
	}
	return opts
}

const EGG_PRICE = 1

func runEgg(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	discordUser := itc.Member.User
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items")
	command := itc.Data.Options[0].Name
	if !rinha.HasEgg(&user) {
		r := &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("EggNotHave", translation.GetLocale(itc)),
				Components: []*disgord.MessageComponent{{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Primary,
						Label:    translation.T("EggBuy", translation.GetLocale(itc)),
						CustomID: "eggBuy",
					}},
				}},
			},
		}
		itcID, err := handler.SendInteractionResponse(ctx, itc, r)
		if err != nil {
			return nil
		}
		done := false
		handler.RegisterHandler(itcID, func(ic *disgord.InteractionCreate) {
			if ic.Member.UserID == user.ID && !done {
				database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
					if u.AsuraCoin >= EGG_PRICE {
						u.AsuraCoin -= EGG_PRICE
						u.Egg = 0
						done = true
					}
					return u
				})
				msg := "EggNoMoney"
				if done {
					msg = "EggBought"
				}
				ic.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: translation.T(msg, translation.GetLocale(itc)),
					},
				})
			}
		}, 100)
		return nil
	}
	switch command {
	case "view":
		authorAvatar, _ := discordUser.AvatarURL(512, true)
		level := rinha.CalcLevel(user.Egg)
		curLevelXP := rinha.CalcXP(level)
		nextLevelXp := rinha.CalcXP(level + 1)
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title: "Egg",
						Footer: &disgord.EmbedFooter{
							Text:    discordUser.Username,
							IconURL: authorAvatar,
						},
						Color: 65535,
						Description: translation.T("EggView", translation.GetLocale(itc), map[string]string{
							"level":       strconv.Itoa(level),
							"curXp":       strconv.Itoa(user.Egg - curLevelXP),
							"nextLevelXp": strconv.Itoa(nextLevelXp - curLevelXP),
						}),
					},
				},
			},
		}
	case "feed":
		opts := generateShardOptions(user.Items)
		if len(opts) == 0 {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("EggNoShards", translation.GetLocale(itc)),
				},
			}
		}
		itcID, err := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Egg",
						Color:       65535,
						Description: translation.T("EggFeed", translation.GetLocale(itc)),
					},
				},
				Components: []*disgord.MessageComponent{{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{
						{
							Type:        disgord.MessageComponentSelectMenu,
							CustomID:    "eggFeedShard",
							Placeholder: translation.T("EggSelectShard", translation.GetLocale(itc)),
							Options:     opts,
							MaxValues:   1,
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
		handler.RegisterHandler(itcID, func(ic *disgord.InteractionCreate) {
			userIC := ic.Member.User
			if userIC.ID != itc.Member.UserID {
				return
			}
			if len(ic.Data.Values) == 0 {
				ic.Data.Values = []string{lastValue}
				if lastValue == "" {
					return
				}
			}
			val := ic.Data.Values[0]
			if val == "nil" {
				return
			}
			lastValue = val
			ItemID := uuid.MustParse(val)
			msg := "InvalidFeedEgg"
			xp := 0
			database.User.UpdateUser(ctx, userIC.ID, func(u entities.User) entities.User {
				item := rinha.GetItemByID(u.Items, ItemID)
				if item == nil {
					return u
				}
				if item.Quantity == 0 {
					return u
				}
				rarity := rinha.Rarity(item.ItemID)
				xp = rinha.ShardToPrice(rarity)
				u.Egg += xp
				database.User.RemoveItem(ctx, u.Items, item.ID)
				msg = "FeedEgg"
				return u
			}, "Items")
			ic.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T(msg, translation.GetLocale(ic), map[string]interface{}{
						"xp": xp,
					}),
				},
			})
		}, 180)

	case "hatch":
		utils.Confirm(ctx, translation.T("EggHatchConfirm", translation.GetLocale(itc)), itc, user.ID, func() {
			rooster := 0
			var rarity rinha.Rarity
			database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
				if !rinha.HasEgg(&u) {
					return u
				}
				rarity = rinha.GetRoosterFromEgg(u.Egg)
				rooster = rinha.GetRandByType(rarity)
				database.User.InsertRooster(ctx, &entities.Rooster{
					Type:   rooster,
					UserID: u.ID,
				})
				u.Egg = -1
				return u
			}, "Galos")

			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{
					{
						Color: rarity.Color(),
						Title: "Egg",
						Description: translation.T("EggHatchDesc", translation.GetLocale(itc), map[string]string{
							"rooster": rinha.Classes[rooster].Name,
							"rarity":  rarity.String(),
						}),
						Image: &disgord.EmbedImage{
							URL: rinha.Sprites[0][rooster-1],
						},
					},
				},
			})

		})
	}
	return nil
}
