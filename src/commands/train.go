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
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "train",
		Description: translation.T("TrainHelp", "pt"),
		Run:         runTrain,
		Cooldown:    20,
		Category:    handler.Rinha,
	})
}

func getMissions(ctx context.Context, user *entities.User) []*entities.Mission {
	missions := rinha.PopulateMissions(user)
	for _, mission := range missions {
		database.User.InsertMission(ctx, user.ID, mission)
	}
	if len(missions) > len(user.Missions) {
		user.Missions = append(user.Missions, missions...)
		user.LastMission = uint64(time.Now().Unix())
		database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
			u.LastMission = uint64(time.Now().Unix())
			return u
		})
	}
	return user.Missions
}

func completeMission(ctx context.Context, user *entities.User, galoAdv *entities.Rooster, winner bool, itc *disgord.InteractionCreate) {
	tempUser := database.User.GetUser(ctx, user.ID, "Missions")
	user.Missions = tempUser.Missions
	user.LastMission = tempUser.LastMission
	if len(user.Missions) == 3 {
		database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
			u.LastMission = uint64(time.Now().Unix())
			return u
		})
	}
	user.Missions = getMissions(ctx, user)
	xp := 0
	money := 0
	toRemove := []int{}
	for i, mission := range user.Missions {
		old := mission.Progress
		done := false
		switch mission.Type {
		case entities.Win:
			if winner {
				mission.Progress++
				if mission.Progress == (mission.Level+1)*3 {
					xp += 55 * (mission.Level + 1)
					money += 45 + (5 * mission.Level)
					done = true
				}
			}
		case entities.Fight:
			mission.Progress++
			if mission.Progress == (mission.Level+1)*6 {
				xp += 55 * (mission.Level + 1)
				money += 45 + (5 * mission.Level)
				done = true
			}
		case entities.WinGalo:
			if winner && galoAdv.Type == mission.Adv {
				mission.Progress++
				if mission.Progress == 4 {
					xp += 200
					money += 80
					done = true
				}
			}
		case entities.FightGalo:
			if galoAdv.Type == mission.Adv {
				mission.Progress++
				if mission.Progress == 8 {
					xp += 200
					money += 80
					done = true
				}
			}
		}
		if mission.Progress != old {
			database.User.UpdateMissions(ctx, user.ID, mission, done)
			if done {
				toRemove = append(toRemove, i)
			} else {
				user.Missions[i] = mission
			}
		}
	}
	galo := rinha.GetEquippedGalo(user)
	if galo.Resets > 0 {
		xp = xp / (galo.Resets + 1)
	}
	if rinha.HasUpgrade(user.Upgrades, 0, 1) {
		money += 3
	}
	if rinha.HasUpgrade(user.Upgrades, 0, 1, 0, 1) {
		xp += 8
	}
	if rinha.HasUpgrade(user.Upgrades, 0, 1, 1, 1) {
		money += 3
	}
	if len(toRemove) > 0 {
		text := "missão"
		if len(toRemove) > 1 {
			text = "missoes"
		}
		handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
			Content: translation.T("MissionComplete", translation.GetLocale(itc), map[string]interface{}{
				"money":    money,
				"xp":       xp,
				"quantity": len(toRemove),
				"text":     text,
				"id":       user.ID,
			}),
		})
	}
	database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
		u.Money += money
		database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
			r.Xp += xp
			return r
		})
		return u
	}, "Galos")
}

func runTrain(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	discordUser := itc.Member.User
	user := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items")
	galo := rinha.GetEquippedGalo(&user)
	text := translation.T("TrainMessage", translation.GetLocale(itc), discordUser.Username)
	utils.Confirm(ctx, text, itc, discordUser.ID, func() {
		authorRinha := isInRinha(ctx, discordUser)
		if authorRinha != "" {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(discordUser.Username, authorRinha).Data.Content,
			})
			return
		}
		galoAdv := entities.Rooster{
			Xp:    galo.Xp,
			Type:  rinha.GetRand(),
			Equip: true,
		}
		if rinha.CalcLevel(galo.Xp) > 1 {
			galoAdv.Xp = rinha.CalcXP(rinha.CalcLevel(galo.Xp) - 1)
		}
		lockEvent(ctx, discordUser.ID, "Clone de "+rinha.Classes[galoAdv.Type].Name)
		defer unlockEvent(ctx, discordUser.ID)
		userAdv := entities.User{
			Galos: []*entities.Rooster{&galoAdv},
		}
		winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
			GaloAuthor:  &user,
			GaloAdv:     &userAdv,
			IDs:         [2]disgord.Snowflake{discordUser.ID},
			AuthorName:  rinha.GetName(discordUser.Username, *galo),
			AdvName:     "Clone de " + rinha.Classes[galoAdv.Type].Name,
			AuthorLevel: rinha.CalcLevel(galo.Xp),
			AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
		}, false)
		if winner == -1 {
			return
		}
		ch := handler.Client.Channel(itc.ChannelID)

		completeMission(ctx, &user, &galoAdv, winner == 0, itc)

		if winner == 0 {
			xpOb := utils.RandInt(20) + 11
			if rinha.HasUpgrade(user.Upgrades, 0) {
				xpOb++
				if rinha.HasUpgrade(user.Upgrades, 0, 1, 1) {
					xpOb += 2
				}
				if rinha.HasUpgrade(user.Upgrades, 0, 1, 1, 0) {
					xpOb += 3
				}
			}
			calc := int(rinha.Classes[galoAdv.Type].Rarity - rinha.Classes[galo.Type].Rarity)
			if calc > 0 {
				calc++
			}
			if calc < 0 {
				if rinha.Classes[galo.Type].Rarity == rinha.Legendary {
					calc -= 2
				}
			}
			xpOb += calc
			money := 8

			if galo.Resets > 0 {
				for i := 0; i < galo.Resets; i++ {
					xpOb = int(float64(xpOb) * 0.75)
				}
			}
			if rinha.HasUpgrade(user.Upgrades, 0, 1, 0) {
				money++
			}
			clanMsg := ""
			if user.TrainLimit == 0 || 1 <= ((uint64(time.Now().Unix())-user.TrainLimitReset)/60/60/24) {
				user.TrainLimit = 0
				user.TrainLimitReset = uint64(time.Now().Unix())
				database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
					u.TrainLimit = 0
					u.TrainLimitReset = user.TrainLimitReset
					return u
				})
			}

			isLimit := rinha.IsInLimit(&user)
			if isLimit {
				need := uint64(time.Now().Unix()) - user.TrainLimitReset
				embed := &disgord.Embed{
					Color: 16776960,
					Title: "Train",
					Description: translation.T("TrainLimit", translation.GetLocale(itc), map[string]interface{}{
						"hours":   23 - (need / 60 / 60),
						"minutes": 59 - (need / 60 % 60),
					}),
				}
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{embed},
				})
				telemetry.Debug(fmt.Sprintf("%s in rinha limit", discordUser.Username), map[string]string{
					"user": fmt.Sprintf("%d", discordUser.ID),
				})
				return
			}
			database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {

				if rinha.IsVip(&u) {
					xpOb += 9
					money++
				}
				item := rinha.GetItem(&u)
				if item != nil {
					if item.Effect == 8 {
						xpOb += xpOb * int(item.Payload)
					}
					if item.Effect == 9 {
						xpOb += 3
						money++
					}
				}
				u.UserXp++
				u.TrainLimit++
				clanUser := database.Clan.GetUserClan(ctx, discordUser.ID)
				clan := clanUser.Clan

				if clan.Name != "" {
					xpOb++
					level := rinha.ClanXpToLevel(clan.Xp)
					if level >= 2 {
						xpOb++
					}
					if level >= 4 {
						money++
					}
					if level >= 6 {
						money++
					}
					if level >= 8 {
						u.UserXp++
					}
					clanXpOb := 2
					if item != nil {
						if item.Effect == 10 {
							if utils.RandInt(101) <= 30 {
								clanXpOb++
							}
						}
					}
					go database.Clan.CompleteClanMission(ctx, clan, discordUser.ID, clanXpOb)
					clanMsg = fmt.Sprintf("\nGanhou **%d** de xp para seu clan", clanXpOb)
				}
				u.Win++
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					r.Xp += xpOb
					return r
				})
				u.Money += money
				return u
			}, "Galos", "Items")
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{
					{
						Color: 16776960,
						Title: "Train",
						Description: translation.T("TrainWin", translation.GetLocale(itc), map[string]interface{}{
							"username": discordUser.Username,
							"xp":       xpOb,
							"money":    money,
							"clanMsg":  clanMsg,
						}),
					},
				},
			})
			galo.Xp += xpOb
			sendLevelUpEmbed(ctx, itc, galo, discordUser, xpOb)
		} else {

			database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {
				u.Lose++
				return u
			})
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{
					{
						Color:       16711680,
						Title:       "Train",
						Description: translation.T("TrainLose", translation.GetLocale(itc), discordUser.Username),
					},
				},
			})

		}
	})
	return nil
}
