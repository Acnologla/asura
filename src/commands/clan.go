package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"asura/src/translation"
	"asura/src/utils"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "clan",
		Description: translation.T("ClanHelp", "pt"),
		Run:         runClan,
		Category:    handler.Rinha,
		Cooldown:    6,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "view",
				Description: translation.T("ClanViewHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "create",
				Description: translation.T("ClanCreate", "pt"),
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeString,
					Required:    true,
					Name:        "name",
					Description: "clan name",
				}),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "mission",
				Description: "Clan mission",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "membros",
				Description: "Veja os membros do seu clan",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "battle",
				Description: "Batalhe um chefe",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "invite",
				Description: translation.T("ClanInvite", "pt"),
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeUser,
					Required:    true,
					Name:        "user",
					Description: "User to invite",
				}),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "banco",
				Description: translation.T("ClanBank", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "remove",
				Description: translation.T("ClanRemove", "pt"),
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeUser,
					Required:    false,
					Name:        "user",
					Description: "User to remove",
				}, &disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeString,
					Required:    false,
					Name:        "user_id",
					Description: "User id to remove",
				}),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "depositar",
				Description: translation.T("ClanMoney", "pt"),
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeNumber,
					Required:    true,
					Name:        "money",
					MinValue:    0,
					MaxValue:    50000,
					Description: "money to deposit",
				}),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "background",
				Description: translation.T("ClanBackground", "pt"),
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeString,
					Required:    true,
					Name:        "background",
					Description: "background url",
				}),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "upgrade",
				Description: translation.T("ClanUpgrade", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "leave",
				Description: translation.T("ClanLeave", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "admin",
				Description: translation.T("ClanAdmin", "pt"),
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeUser,
					Required:    true,
					Name:        "user",
					Description: "User admin",
				}),
			},
		),
	})

}

const PRIZE_XP = 2500
const PRIZE_MONEY = 3200

func completeClanBoss(ctx context.Context, clan *entities.Clan, channel disgord.Snowflake) {
	var xpPrize, moneyPrize int
	database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
		originalHp := rinha.CalcBossLife(&c)
		percentange := utils.Min(float64(originalHp-c.BossLife)/float64(originalHp), 1)
		if c.BossLife > originalHp {
			percentange = 0
		}
		xpPrize = int(float64(PRIZE_XP) * percentange)
		moneyPrize = int(float64(PRIZE_MONEY) * percentange)
		for _, member := range c.Members {
			database.User.UpdateUser(ctx, member.ID, func(u entities.User) entities.User {
				u.Money += moneyPrize
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					r.Xp += xpPrize
					return r
				})
				return u
			}, "Galos")
		}
		return c
	}, "Members")
	handler.Client.Channel(channel).CreateMessage(&disgord.CreateMessage{
		Embeds: []*disgord.Embed{
			{
				Title:       "Clan boss",
				Color:       65535,
				Description: fmt.Sprintf("O boss foi derrotado, todos os membros receberam **%d** de xp e **%d** de dinheiro", xpPrize, moneyPrize),
			},
		},
	})

}

func generateUpgradesOptions() (opts []*disgord.SelectMenuOption) {
	for _, name := range [...]string{"membros", "missão"} {
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       name,
			Description: fmt.Sprintf("%s upgrade", name),
			Value:       name,
		})
	}
	return
}

func getBattleEmbed(users []*disgord.User, title string, minutes int) *disgord.Embed {
	msg := "Membros que iram participar da batalha: \n\n"
	for _, user := range users {
		msg += fmt.Sprintf("%s\n", user.Mention())
	}
	return &disgord.Embed{
		Title:       title,
		Description: msg,
		Color:       65535,
		Footer: &disgord.EmbedFooter{
			Text: fmt.Sprintf("Em %d minutos a batalha ira começar", minutes),
		},
	}
}

func isInUsers(users []*disgord.User, user *disgord.User) bool {
	for _, u := range users {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}

func runClan(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	command := itc.Data.Options[0].Name
	user := database.User.GetUser(ctx, itc.Member.UserID)
	userClan := database.Clan.GetUserClan(ctx, user.ID, "Members")
	clan := userClan.Clan
	maxMembers := rinha.GetMaxMembers(clan)
	member := userClan.Member
	ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))
	clanLevel := rinha.ClanXpToLevel(clan.Xp)
	var msg string
	if command != "create" && clan.Name == "" {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("NoClan", translation.GetLocale(itc)),
			},
		}
	}
	if (command == "invite" || command == "remove" || command == "admin" || command == "background" || command == "upgrade") && member.Role == entities.Member {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("MissingPermissions", translation.GetLocale(itc)),
			},
		}

	}
	switch command {
	case "create":
		if clan.Name != "" {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("ClanArleadyIn", translation.GetLocale(itc)),
				},
			}
		}
		name := itc.Data.Options[0].Options[0].Value.(string)
		name = rinha.Format(name)
		if name == "_test" {
			return nil
		}
		if len(name) >= 25 || len(name) <= 4 {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("ClanNameLength", translation.GetLocale(itc)),
				},
			}
		}
		clan := database.Clan.GetClan(ctx, name)
		if clan.Name != "" {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("ClanArleadyExists", translation.GetLocale(itc)),
				},
			}
		}
		database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
			clan := database.Clan.GetClan(ctx, name)
			if clan.Name != "" {
				msg = "ClanAlreadyIn"
				return u
			}
			if 1000 > u.Money {
				msg = "ClanNotEnoughMoney"
				return u
			}
			u.Money -= 1000
			database.Clan.CreateClan(ctx, entities.Clan{
				Name: name,
			}, user.ID)
			return u
		})
		if msg != "" {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T(msg, translation.GetLocale(itc)),
				},
			}
		}
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("ClanCreated", translation.GetLocale(itc), name),
			},
		}
	case "membros":
		msgID, _ := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Carregando...",
			},
		})
		memberMsg := ""
		sort.SliceStable(clan.Members, func(i, j int) bool {
			return clan.Members[i].Xp > clan.Members[j].Xp
		})
		for _, member := range clan.Members {
			user, err := handler.Client.User(disgord.Snowflake(member.ID)).Get()
			if err == nil {
				username := user.Username
				if strings.Contains(strings.ToLower(username), "deleted_user") {
					username = user.ID.String()
				}
				memberMsg += fmt.Sprintf("[**%s**] `%s` (**%d** XP)\n", member.Role.ToString(), username, member.Xp)
			}
		}

		if len(memberMsg) > 1650 {
			memberMsg = memberMsg[:1650]
		}

		embed := &disgord.Embed{
			Title: clan.Name,
			Color: 65535,
			Footer: &disgord.EmbedFooter{
				Text: translation.T("ClanFooter", translation.GetLocale(itc)),
			},
			Description: fmt.Sprintf("Membros (**%d/%d**):\n\n%s", len(clan.Members), maxMembers, memberMsg),
		}
		str := ""
		handler.EditInteractionResponse(ctx, msgID, itc, &disgord.UpdateMessage{
			Embeds:  &([]*disgord.Embed{embed}),
			Content: &str,
		})

		return nil
	case "view":
		level := rinha.ClanXpToLevel(clan.Xp)
		benefits := rinha.GetBenefits(clan.Xp)
		bg := clan.Background
		embed := &disgord.Embed{
			Title: clan.Name,
			Color: 65535,
			Footer: &disgord.EmbedFooter{
				Text: translation.T("ClanFooter", translation.GetLocale(itc)),
			},
			Description: translation.T("ClanDescription", translation.GetLocale(itc), map[string]interface{}{
				"level":       level,
				"xp":          clan.Xp,
				"needXp":      rinha.ClanLevelToXp(level),
				"benefits":    benefits,
				"members":     len(clan.Members),
				"maxMembers":  maxMembers,
				"membersText": "",
				"bossHP":      clan.BossLife,
				"bossHPMax":   rinha.CalcBossLife(clan),
			}),
		}
		if bg != "" {
			embed.Image = &disgord.EmbedImage{
				URL: bg,
			}
		}

		handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: ([]*disgord.Embed{embed}),
			},
		})

		return nil
	case "mission":
		clan = rinha.PopulateClanMissions(clan)
		database.Clan.UpdateClan(ctx, clan, func(clanUpdate entities.Clan) entities.Clan {
			clanUpdate.Mission = clan.Mission
			clanUpdate.MissionProgress = clan.MissionProgress
			return clanUpdate
		})
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{{
					Title:       "Clan mission",
					Description: rinha.MissionToString(clan),
					Color:       65535,
				}},
			},
		}
	case "invite":
		completeAchievement(ctx, itc, 7)
		user := utils.GetOptionsUser(itc.Data.Options[0].Options, itc, 0)
		if user == nil {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Invalid user",
				},
			}
		}
		text := translation.T("InviteToClan", translation.GetLocale(itc), map[string]interface{}{
			"clan":     clan.Name,
			"username": user.Username,
		})
		utils.Confirm(ctx, text, itc, user.ID, func() {
			database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
				uClan := database.Clan.GetUserClan(ctx, user.ID)
				if uClan.Clan.Name != "" {
					msg = "UserArleadyInClan"
				} else {
					if len(c.Members) >= maxMembers {
						msg = "ClanMaxMembers"
					} else {
						database.User.GetUser(ctx, user.ID)
						database.Clan.InsertMember(ctx, &c, &entities.ClanMember{
							ID: user.ID,
						})
						msg = "SucessInvite"
						itc.Member.User.ID = user.ID
						completeAchievement(ctx, itc, 7)
					}
				}
				return c
			}, "Members")
		})

		ch.CreateMessage(&disgord.CreateMessage{
			Content: translation.T(msg, translation.GetLocale(itc), user.Username),
		})
	case "remove":
		if len(itc.Data.Options[0].Options) == 0 {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Invalid user",
				},
			}
		}
		user := utils.GetOptionsUser(itc.Data.Options[0].Options, itc, 0)
		if user == nil {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Invalid user",
				},
			}
		}
		database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
			uClan := database.Clan.GetUserClan(ctx, user.ID)
			if uClan.Clan.Name != clan.Name {
				msg = "UserNotInClan"
			} else {
				if member.Role <= uClan.Member.Role {
					msg = "NoPermission"
				} else {
					msg = "SucessRemove"
					database.Clan.RemoveMember(ctx, &c, user.ID, false)
				}
			}
			return c
		})

		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T(msg, translation.GetLocale(itc), user.Username),
			},
		}
	case "admin":
		if len(itc.Data.Options[0].Options) == 0 {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Invalid user",
				},
			}
		}
		user := utils.GetOptionsUser(itc.Data.Options[0].Options, itc, 0)
		database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
			uClan := database.Clan.GetUserClan(ctx, user.ID)
			if uClan.Clan.Name != clan.Name {
				msg = "UserNotInClan"
			} else {
				if member.Role != entities.Owner || uClan.Member.Role == entities.Owner {
					msg = "NoPermission"
				} else {
					uClan.Member.Role = entities.Administrator
					database.Clan.UpdateMember(ctx, &c, uClan.Member)
					msg = "SucessAdmin"
				}
			}
			return c
		})
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T(msg, translation.GetLocale(itc), user.Username),
			},
		}
	case "battle":
		if len(clan.Members) < 7 || clanLevel < 3 {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Seu clan precisa de pelo menos 7 membros para batalhar e ser level 3",
				},
			}
		}
		daysPassed := (uint64(time.Now().Unix()) - clan.BossDate) / 60 / 60 / 24
		if daysPassed >= 3 && clan.BossDate != 0 {
			database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
				c.BossDate = 0
				return c
			})
			completeClanBoss(ctx, clan, itc.ChannelID)
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Use o comando novamente para iniciar a batalha",
				},
			}
		}

		if clan.BossDate == 0 {
			database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
				c.BossDate = uint64(time.Now().Unix())
				c.BossLife = rinha.CalcBossLife(&c)
				c.BossBattles = []disgord.Snowflake{}
				return c
			}, "Members")
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "O chefe chegou, vocês tem 3 dias para derrota-lo",
				},
			}
		}

		if 0 >= clan.BossLife {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Aguarde 3 dias para lutar com o chefe novamente",
				},
			}
		}

		if utils.Has(clan.BossBattles, itc.Member.User.ID) {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Você ja batalhou com o chefe",
				},
			}
		}

		database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
			c.BossBattles = append(c.BossBattles, itc.Member.User.ID)
			return c
		})
		user := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items", "Trials")
		userGalo := rinha.GetEquippedGalo(&user)

		gAdv := &entities.Rooster{
			Xp:     rinha.CalcXP(utils.Max((clan.BossLife-100)/5, 7)) + 1,
			Type:   53,
			Equip:  true,
			Resets: userGalo.Resets,
		}

		userAdv := entities.User{
			Galos:      []*entities.Rooster{gAdv},
			Attributes: [6]int{0, user.Attributes[1], 0, 0, 50, 350},
		}
		handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "A batalha esta iniciando",
			},
		})
		winner, battle := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
			GaloAuthor:  &user,
			GaloAdv:     &userAdv,
			IDs:         [2]disgord.Snowflake{itc.Member.UserID},
			AuthorName:  rinha.GetName(itc.Member.User.Username, *userGalo),
			AdvName:     "Chefe do clan",
			AuthorLevel: rinha.CalcLevel(userGalo.Xp),
			AdvLevel:    rinha.CalcLevel(gAdv.Xp),
		}, false)
		database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
			c.BossLife = utils.Max(battle.Fighters[1].Life, 0)
			return c
		})
		if winner == 0 {
			completeClanBoss(ctx, clan, itc.ChannelID)
		} else {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: "O boss venceu",
			})
		}
	case "banco":
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title: "Banco",
						Color: 65535,
						Description: translation.T("ClanBankDesc", translation.GetLocale(itc), map[string]interface{}{
							"money":          clan.Money,
							"membersUpgrade": rinha.CalcClanUpgrade(clan.MembersUpgrade, 1),
							"missionUpgrade": rinha.CalcClanUpgrade(clan.MissionsUpgrade, 1),
						}),
					},
				},
			},
		}
	case "depositar":
		money := int(itc.Data.Options[0].Options[0].Value.(float64))
		database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
			if u.Money < money {
				msg = "NoMoney"
			} else {
				u.Money -= money
				database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
					c.Money += money
					return c
				})
				msg = "ClanMoneyAdd"
			}
			return u
		})
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T(msg, translation.GetLocale(itc), money),
			},
		}
	case "upgrade":
		msgID, err := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title: "Upgrades",
						Color: 65535,
						Description: translation.T("ClanBankMinified", translation.GetLocale(itc), map[string]interface{}{
							"money":          clan.Money,
							"membersUpgrade": rinha.CalcClanUpgrade(clan.MembersUpgrade, 1),
							"missionUpgrade": rinha.CalcClanUpgrade(clan.MissionsUpgrade, 1),
						}),
					},
				},
				Components: []*disgord.MessageComponent{
					{
						Type: disgord.MessageComponentActionRow,
						Components: []*disgord.MessageComponent{
							{
								Type:        disgord.MessageComponentButton + 1,
								Options:     generateUpgradesOptions(),
								CustomID:    "upgrade",
								MaxValues:   1,
								Placeholder: "Select upgrade",
							},
						},
					},
				},
			},
		})
		if err != nil {
			return nil
		}
		handler.RegisterHandler(msgID, func(ic *disgord.InteractionCreate) {
			if len(ic.Data.Values) == 0 {
				return
			}
			upgrade := ic.Data.Values[0]
			if ic.Member.UserID == itc.Member.UserID {
				database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
					var money int
					if upgrade == "membros" {
						money = rinha.CalcClanUpgrade(c.MembersUpgrade, 1)
					}
					if upgrade == "missão" {
						money = rinha.CalcClanUpgrade(c.MissionsUpgrade, 1)
					}
					if money > c.Money {
						msg = "ClanNoMoney"
					} else {
						msg = "ClanUpgradePuchased"
						c.Money -= money
						if upgrade == "membros" {
							c.MembersUpgrade++
						}
						if upgrade == "missão" {
							c.MissionsUpgrade++
						}
					}
					return c
				})
				handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: translation.T(msg, translation.GetLocale(itc), upgrade),
					},
				})
			}
		}, 100)
	case "leave":
		var msg string
		database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
			err := database.Clan.RemoveMember(ctx, &c, user.ID, true)
			if err == nil {
				msg = "sucessLeaveClan"
			}
			return c
		}, "Members")
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T(msg, translation.GetLocale(itc), clan.Name),
			},
		}
	case "background":
		img := itc.Data.Options[0].Options[0].Value.(string)
		if !utils.CheckImage(img) {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("InvalidImage", translation.GetLocale(itc)),
				},
			}
		}
		database.Clan.UpdateClan(ctx, clan, func(c entities.Clan) entities.Clan {
			if c.Money < 10000 {
				msg = "ClanNoMoney"
			} else {
				c.Money -= 10000
				c.Background = img
				msg = "ChangedClanBackground"
			}
			return c
		})
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T(msg, translation.GetLocale(itc)),
			},
		}
	}
	return nil
}
