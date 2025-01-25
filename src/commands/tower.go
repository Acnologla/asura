package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"asura/src/utils"
	"context"
	"fmt"
	"math"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "tower",
		Description: translation.T("TowerHelp", "pt"),
		Run:         runTower,
		Cooldown:    6,
		Category:    handler.Rinha,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "rank",
				Description: translation.T("TowerRankHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "info",
				Description: translation.T("TowerInfoHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "battle",
				Description: translation.T("TowerBattleHelp", "pt"),
			},
		),
	})
}

func runTower(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	discordUser := itc.Member.User
	user := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items", "Trials")
	galo := rinha.GetEquippedGalo(&user)
	tower := database.User.GetTower(ctx, itc.Member.UserID)
	command := itc.Data.Options[0].Name
	switch command {
	case "info":
		authorAvatar, _ := discordUser.AvatarURL(512, true)
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title: "Tower",
						Footer: &disgord.EmbedFooter{
							Text:    discordUser.Username,
							IconURL: authorAvatar,
						},
						Color: 65535,
						Description: translation.T("TowerInfo", translation.GetLocale(itc), map[string]string{
							"timeUntilReset": utils.TimeUntilNextSunday(),
							"floor":          strconv.Itoa(tower.Floor),
						}),
					},
				},
			},
		}
	case "rank":
		towers := database.User.SortTowers(ctx, 16)
		msgID, _ := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Carregando...",
			},
		})
		var rank string
		for i, tower := range towers {
			u, err := handler.Client.User(tower.UserID).Get()
			if err == nil {
				rank += fmt.Sprintf("[**%d**] %s - **%d**\n", i+1, u.Username, tower.Floor)
			}
		}

		str := ""
		embeds := []*disgord.Embed{
			{
				Title:       "Tower Rank",
				Color:       65535,
				Description: rank,
			},
		}
		handler.EditInteractionResponse(ctx, msgID, itc, &disgord.UpdateMessage{
			Embeds:  &embeds,
			Content: &str,
		})
	case "battle":
		authorRinha := isInRinha(ctx, discordUser)
		if authorRinha != "" {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(discordUser.Username, authorRinha).Data.Content,
			})
			return rinhaMessage(discordUser.Username, authorRinha)
		}
		if tower.Floor > rinha.MAXIMUM_FLOOR {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("TowerFinish", translation.GetLocale(itc)),
				},
			}
		}

		if tower.Lose {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("TowerLoseWait", translation.GetLocale(itc), utils.TimeUntilNextSunday()),
				},
			}
		}
		lockEvent(ctx, discordUser.ID, "Tower")
		defer unlockEvent(ctx, discordUser.ID)
		level := 5 + tower.Floor
		galoAdv := &entities.Rooster{
			Xp:     rinha.CalcXP(level) + 1,
			Type:   rinha.GetRandByType(rinha.GetFloorRarity(tower.Floor)),
			Equip:  true,
			Resets: tower.Floor / 50,
		}
		userAdv := &entities.User{
			Galos: []*entities.Rooster{galoAdv},
		}

		/*
			if tower.Floor >= 60 {
				userAdv.Items = user.Items
			}
		*/

		if tower.Floor >= 90 {
			galoAdv.Evolved = true

		}

		if tower.Floor >= 130 {
			userAdv.Attributes = user.Attributes
		}

		handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "A batalha esta iniciando",
			},
		})

		winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
			GaloAuthor: &user,
			GaloAdv:    userAdv,
			IDs:        [2]disgord.Snowflake{discordUser.ID},

			AuthorName:  rinha.GetName(discordUser.Username, *galo),
			AdvName:     "Tower",
			AuthorLevel: rinha.CalcLevel(galo.Xp),
			AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
			NoItems:     false,
		}, false)
		if winner == -1 {
			return nil
		}
		ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))
		go completeMission(ctx, &user, galoAdv, winner == 0, itc, "tower")
		if winner == 0 {

			xp, money := rinha.CalcTowerReward(tower.Floor)
			lootbox := rinha.CalcFloorReward(tower.Floor)
			if tower.Floor >= 150 {
				completeAchievement(ctx, itc, 17)
			}
			database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {

				if rinha.IsVip(&u) {
					money += int(float64(money) * 0.3)
					xp += int(float64(xp) * 0.3)
				}

				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					resets := int(math.Min(float64(r.Resets), 6))
					xp = xp / (resets + 1)
					r.Xp += xp
					return r
				})
				u.Money += money
				if lootbox != -1 {
					database.User.InsertItem(ctx, user.ID, u.Items, lootbox, entities.LootboxType)
				}
				tower.Floor++
				database.User.UpdateTower(ctx, tower)

				return u
			}, "Items", "Galos")
			description := translation.T("TowerWin", translation.GetLocale(itc), tower.Floor)
			if xp != 0 {
				description += "\n" + translation.T("TowerWinBonus", translation.GetLocale(itc), map[string]string{
					"xp":    strconv.Itoa(xp),
					"money": strconv.Itoa(money),
				})
			}
			if lootbox != -1 {
				lb := rinha.LootNames[lootbox]
				description += "\n" + translation.T("TowerWinReward", translation.GetLocale(itc), lb)
			}
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{{
					Color:       16776960,
					Title:       "Tower",
					Description: description,
				}},
			})

		} else {

			database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {
				tower.Lose = true
				database.User.UpdateTower(ctx, tower)
				return u
			})

			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{{
					Color:       16711680,
					Title:       "Tower",
					Description: translation.T("TowerLose", translation.GetLocale(itc)),
				}},
			})
		}
	}
	return nil
}
