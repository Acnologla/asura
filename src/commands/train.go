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
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "train",
		Description: translation.T("TrainHelp", "pt"),
		AliasesMsg:  []string{"tn"},
		Run:         runTrain,
		Cooldown:    6,
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

func completeAchievement(ctx context.Context, itc *disgord.InteractionCreate, achievementID int) {
	database.User.UpdateUser(ctx, itc.Member.User.ID, func(u entities.User) entities.User {
		if rinha.HasAchievement(&u, achievementID) {
			return u
		}
		ach := rinha.Achievements[achievementID]
		database.User.InsertItem(ctx, u.ID, u.Items, achievementID, entities.AchievementType)
		u.Money += ach.Money
		u.AsuraCoin += ach.AsuraCoin
		msg := fmt.Sprintf("**%d** Money", ach.Money)
		if ach.AsuraCoin > 0 {
			msg = fmt.Sprintf("**%d** AsuraCoin", ach.AsuraCoin)
		}
		go handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
			Content: translation.T("AcvComplete", translation.GetLocale(itc), map[string]interface{}{
				"name": ach.Name,
				"msg":  msg,
				"id":   u.ID,
			}),
		})
		return u
	}, "Items")
}

func completeMission(ctx context.Context, user *entities.User, galoAdv *entities.Rooster, winner bool, itc *disgord.InteractionCreate, battleType string) {
	tempUser := database.User.GetUser(ctx, user.ID, "Missions")
	user.Missions = tempUser.Missions
	user.LastMission = tempUser.LastMission
	if len(user.Missions) == 4 {
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
				if mission.Progress >= 3 {
					xp += 260
					money += 110
					done = true
				}
			}
		case entities.FightGalo:
			if galoAdv.Type == mission.Adv {
				mission.Progress++
				if mission.Progress >= 6 {
					xp += 260
					money += 110
					done = true
				}
			}
		case entities.PlayTrial:
			if battleType == "trial" {
				mission.Progress++
				xp += 100
				money += 50
				done = true
			}
		case entities.WinRaid:
			if battleType == "raid" {
				mission.Progress++
				xp += 285
				money += 125
				done = true
			}
		case entities.FightTower:
			if battleType == "tower" {
				mission.Progress++
				if mission.Progress >= (mission.Level+1)*3 {
					xp += 60 * (mission.Level + 1)
					money += 50 + (5 * mission.Level)
					done = true
				}
			}
		case entities.WinDungeon:
			if battleType == "dungeon" {
				if winner {
					mission.Progress++
					if mission.Progress >= (mission.Level+1)*2 {
						xp += 60 * (mission.Level + 1)
						money += 55 + (5 * mission.Level)
						done = true
					}
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
		money += 7
	}
	if rinha.HasUpgrade(user.Upgrades, 0, 1, 0, 1) {
		xp += 30
	}
	if rinha.HasUpgrade(user.Upgrades, 0, 1, 1, 1) {
		money += 7
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
		database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
			u.Money += money
			database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
				r.Xp += xp
				return r
			})
			return u
		}, "Galos")
	}

}

func runTrain(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	discordUser := itc.Member.User
	user := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items", "Trials")
	galo := rinha.GetEquippedGalo(&user)
	//text := translation.T("TrainMessage", translation.GetLocale(itc), discordUser.Username)
	//utils.Confirm(ctx, text, itc, discordUser.ID, func() {

	authorRinha := isInRinha(ctx, discordUser)
	if authorRinha != "" {
		handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
			Content: rinhaMessage(discordUser.Username, authorRinha).Data.Content,
		})
		return nil
	}
	galoAdv := entities.Rooster{
		Xp:    galo.Xp,
		Type:  rinha.GetRand(),
		Equip: true,
	}

	if rinha.CalcLevel(galo.Xp) > 1 {
		galoAdv.Xp = rinha.CalcXP(rinha.CalcLevel(galo.Xp) - 1)
	}

	advClass := rinha.Classes[galoAdv.Type]
	authorClass := rinha.Classes[galo.Type]
	if rinha.Epic >= authorClass.Rarity {
		sub := 0
		if rinha.Epic == authorClass.Rarity {
			sub = 1
		}
		if advClass.Rarity == rinha.Legendary {
			galoAdv.Xp = rinha.CalcXP(rinha.CalcLevel(galo.Xp) - 3 - sub)
		}
		if advClass.Rarity == rinha.Mythic {
			galoAdv.Xp = rinha.CalcXP(rinha.CalcLevel(galo.Xp) - 4 - sub)
		}
	}

	if galoAdv.Xp < 0 {
		galoAdv.Xp = 5
	}
	lockEvent(ctx, discordUser.ID, "Clone de "+advClass.Name)
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
		return nil
	}
	ch := handler.Client.Channel(itc.ChannelID)

	completeMission(ctx, &user, &galoAdv, winner == 0, itc, "train")
	isLimit := rinha.IsInLimit(&user)
	resetLimit := user.TrainLimit == 0 || 1 <= ((uint64(time.Now().Unix())-user.TrainLimitReset)/60/60/24)
	if isLimit && !resetLimit {
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
		return nil
	}
	if winner == 0 {
		xpOb := utils.RandInt(13) + 12
		eggXpOb := 11 + utils.RandInt(12)
		if rinha.HasUpgrade(user.Upgrades, 0) {
			xpOb += 3
			if rinha.HasUpgrade(user.Upgrades, 0, 1, 1) {
				xpOb += 3
			}
			if rinha.HasUpgrade(user.Upgrades, 0, 1, 1, 0) {
				xpOb += 4
			}
		}
		calc := int(rinha.Classes[galoAdv.Type].Rarity - rinha.Classes[galo.Type].Rarity)
		if calc > 0 {
			calc++
		}
		if calc < 0 {
			if rinha.Classes[galo.Type].Rarity >= rinha.Legendary {
				calc -= 2
			}
		}
		xpOb += calc
		money := 7 + utils.RandInt(2)

		if rinha.HasUpgrade(user.Upgrades, 0, 1, 0) {
			money += 3
		}
		clanMsg := ""
		if resetLimit {
			user.TrainLimit = 0
			user.TrainLimitReset = uint64(time.Now().Unix())
			database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
				u.TrainLimit = 0
				u.TrainLimitReset = user.TrainLimitReset
				return u
			})
		}
		bpXP := 0
		dropKey := -1
		if galo.Resets > 0 {
			for i := 0; i < galo.Resets; i++ {
				xpOb = int(float64(xpOb) * 0.75)
			}
		}
		xpTotal := 0

		database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {

			item := rinha.GetItem(&u)
			if item != nil {
				if item.Effect == 8 {
					xpOb += int(float64(xpOb) * item.Payload)
				}
				if item.Effect == 9 {
					xpOb += 2
				}
				if item.Effect == 13 {
					money += int(item.Payload)
				}
			}

			if rinha.IsVip(&u) {
				xpOb += 12
				money += 3
				eggXpOb += 6
			}

			u.UserXp++
			u.TrainLimit++
			clanUser := database.Clan.GetUserClan(ctx, discordUser.ID, "Members")
			clan := clanUser.Clan

			var clanLevel = 0
			if clan.Name != "" {
				xpOb++
				level := rinha.ClanXpToLevel(clan.Xp)
				clanLevel = level
				if level >= 2 {
					xpOb++
				}
				if level >= 4 {
					money++
				}
				if level >= 6 {
					eggXpOb++
				}
				if level >= 8 {
					u.UserXp++
				}
				if level >= 10 {
					eggXpOb++
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
			if rinha.HasEgg(&user) {
				u.Egg += eggXpOb
			}

			database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
				bpXP = database.User.UpdateBp(ctx, &u, &r)
				r.Xp += xpOb
				xpTotal = r.Xp
				return r
			})

			u.Money += money
			if rinha.DropKey(u.UserXp, rinha.IsVip(&u), clanLevel) {
				key := rinha.GetKeyRarity()
				dropKey = int(key)
				database.User.InsertItem(ctx, u.ID, u.Items, int(key), entities.KeyType)
				telemetry.Debug(fmt.Sprintf("%s get key %s", discordUser.Username, key.String()), map[string]string{
					"user":      strconv.FormatUint(uint64(u.ID), 10),
					"keyRarity": key.String(),
				})
			}
			return u
		}, "Galos", "Items")
		bpLevel := rinha.CalcBPLevel(user.BattlePass + bpXP)
		w := user.Win + 1
		if w >= 100 {
			completeAchievement(ctx, itc, 0)
			if w >= 1000 {
				completeAchievement(ctx, itc, 1)
				if w >= 10000 {
					completeAchievement(ctx, itc, 2)
				}
				if w >= 100000 {
					completeAchievement(ctx, itc, 3)
				}
			}
		}
		if bpLevel >= len(rinha.BattlePass)/2 {
			completeAchievement(ctx, itc, 10)
			if bpLevel >= len(rinha.BattlePass) {
				completeAchievement(ctx, itc, 11)
			}
		}
		if dropKey != -1 {
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{
					{
						Color: 65535,
						Title: "Train",
						Description: translation.T("TrainWinKey", translation.GetLocale(itc), map[string]interface{}{
							"username":  discordUser.Username,
							"xp":        xpOb,
							"money":     money,
							"clanMsg":   clanMsg,
							"bpXP":      bpXP,
							"keyRarity": rinha.Rarity(dropKey).String(),
						}),
					},
				},
			})
		} else {
			level := rinha.CalcLevel(xpTotal)
			curLevelXP := rinha.CalcXP(level)
			nextLevelXp := rinha.CalcXP(level + 1)
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{
					{
						Color: 16776960,
						Title: "Train",
						Description: translation.T("TrainWin", translation.GetLocale(itc), map[string]interface{}{
							"username":  discordUser.Username,
							"xp":        xpOb,
							"money":     money,
							"clanMsg":   clanMsg,
							"bpXP":      bpXP,
							"totalXp":   xpTotal - curLevelXP,
							"xpToLevel": nextLevelXp - curLevelXP,
						}),
					},
				},
			})
		}
		galo.Xp += xpOb
		sendLevelUpEmbed(ctx, itc, galo, discordUser, xpOb)
	} else {
		xpO := utils.RandInt(13) + 3
		moneyO := utils.RandInt(3) + 1
		if galo.Resets > 0 {
			for i := 0; i < galo.Resets; i++ {
				xpO = int(float64(xpO) * 0.75)
			}
		}
		if user.TrainLimit >= 100 {
			xpO = 0
			moneyO = 0
		} else {
			database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {
				u.Lose++
				u.Money += moneyO

				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					r.Xp += xpO
					return r
				})

				return u
			}, "Galos", "Items")

		}
		if user.Lose+1 >= 1000 {
			completeAchievement(ctx, itc, 4)
			if user.Lose+1 >= 10000 {
				completeAchievement(ctx, itc, 5)
			}
		}
		ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{
				{
					Color: 16711680,
					Title: "Train",
					Description: translation.T("TrainLose", translation.GetLocale(itc), map[string]interface{}{
						"username": discordUser.Username,
						"xp":       xpO,
						"money":    moneyO,
					}),
				},
			},
		})

		galo.Xp += xpO
		sendLevelUpEmbed(ctx, itc, galo, discordUser, xpO)
	}
	//	})
	return nil
}
