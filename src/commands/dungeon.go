package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
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
		Name:        "dungeon",
		Description: translation.T("DungeonHelp", "pt"),
		Run:         runDungeon,
		Cooldown:    15,
		Category:    handler.Rinha,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "info",
				Description: translation.T("DungeonInfoHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "battle",
				Description: translation.T("DungeonBattleHelp", "pt"),
			},
		),
	})
}

func runDungeon(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	discordUser := itc.Member.User
	user := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items")
	galo := rinha.GetEquippedGalo(&user)
	command := itc.Data.Options[0].Name
	switch command {
	case "info":
		authorAvatar, _ := discordUser.AvatarURL(512, true)
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title: "Dungeon",
						Footer: &disgord.EmbedFooter{
							Text:    discordUser.Username,
							IconURL: authorAvatar,
						},
						Color:       65535,
						Description: translation.T("DungeonFloor", translation.GetLocale(itc), user.Dungeon),
					},
				},
			},
		}
	case "battle":
		authorRinha := isInRinha(ctx, discordUser)
		if authorRinha != "" {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(discordUser.Username, authorRinha).Data.Content,
			})
			return rinhaMessage(discordUser.Username, authorRinha)
		}
		if len(rinha.Dungeon) == user.Dungeon {
			database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
				u.Dungeon = 0
				u.DungeonReset++
				return u
			})
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("DungeonFinish", translation.GetLocale(itc)),
				},
			}
		}
		dungeon := rinha.Dungeon[user.Dungeon]
		galoAdv := dungeon.Boss
		lockEvent(ctx, discordUser.ID, "Boss "+rinha.Classes[galoAdv.Type].Name)
		defer unlockEvent(ctx, discordUser.ID)
		multiplier := 1 + user.DungeonReset

		AdvLVL := rinha.CalcLevel(galoAdv.Xp) * multiplier

		ngaloAdv := &entities.Rooster{
			Xp:    rinha.CalcXP(AdvLVL) + 1,
			Type:  galoAdv.Type,
			Equip: true,
		}
		itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "A batalha esta iniciando",
			},
		})
		winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
			GaloAuthor: &user,
			GaloAdv: &entities.User{
				Galos: []*entities.Rooster{ngaloAdv},
			},
			IDs: [2]disgord.Snowflake{discordUser.ID},

			AuthorName:  rinha.GetName(discordUser.Username, *galo),
			AdvName:     "Boss " + rinha.Classes[galoAdv.Type].Name,
			AuthorLevel: rinha.CalcLevel(galo.Xp),
			AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
			NoItems:     false,
		}, false)
		if winner == -1 {
			return nil
		}
		ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))

		if winner == 0 {

			if user.DungeonReset != 0 && user.Dungeon+1 != len(rinha.Dungeon) {
				database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {
					u.Dungeon++
					return u
				})
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{{
						Color: 16776960,
						Title: "Dungeon",
						Description: translation.T("DungeonWin", translation.GetLocale(itc), map[string]interface{}{
							"floor":    user.Dungeon + 1,
							"username": discordUser.Username,
							"msg":      "",
						}),
					}},
				})
				return nil
			}
			var endMsg string
			percents := rinha.DungeonsPercentages[dungeon.Level]
			value := utils.RandInt(101)
			var selected rinha.DungeonWin
			for _, v := range percents {
				if v.Percentage >= value || v.Percentage == 0 {
					selected = v
					break
				}
			}
			database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {
				if selected.PrizeType == entities.LootboxType {
					database.User.InsertItem(ctx, u.ID, u.Items, selected.PrizeRarity, entities.LootboxType)
					endMsg = translation.T("DungeonWinLootbox", translation.GetLocale(itc), rinha.LootNames[selected.PrizeRarity])
				} else {
					item := rinha.GetItemByLevel(selected.PrizeRarity)
					_item := rinha.Items[item]
					database.User.InsertItem(ctx, u.ID, u.Items, item, entities.NormalType)
					endMsg = translation.T("DungeonWinItem", translation.GetLocale(itc), map[string]interface{}{
						"rarity": rinha.LevelToString(_item.Level),
						"name":   _item.Name,
					})
				}
				u.Dungeon++
				return u
			}, "Items")
			tag := discordUser.Username + "#" + discordUser.Discriminator.String()
			telemetry.Debug(fmt.Sprintf("%s %s", tag, endMsg), map[string]string{
				"user":         strconv.FormatUint(uint64(discordUser.ID), 10),
				"dungeonLevel": fmt.Sprintf("%d", user.Dungeon),
			})
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{{
					Color: 16776960,
					Title: "Dungeon",
					Description: translation.T("DungeonWin", translation.GetLocale(itc), map[string]interface{}{
						"floor":    user.Dungeon + 1,
						"username": discordUser.Username,
						"msg":      endMsg,
					}),
				}},
			})
		} else {
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{{
					Color:       16711680,
					Title:       "Dungeon",
					Description: translation.T("DungeonLose", translation.GetLocale(itc), discordUser.Username),
				}},
			})
		}
	}
	return nil
}
