package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/translation"
	"asura/src/utils"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "clan",
		Description: translation.T("ClanHelp", "pt"),
		Run:         runClan,
		Category:    handler.Rinha,
		Cooldown:    8,
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
func runClan(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	command := itc.Data.Options[0].Name
	user := database.User.GetUser(itc.Member.UserID)
	userClan := database.Clan.GetUserClan(user.ID, "Members")
	clan := userClan.Clan
	maxMembers := rinha.GetMaxMembers(clan)
	member := userClan.Member
	ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))
	var msg string
	if command != "create" && clan.Name == "" {
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T("NoClan", translation.GetLocale(itc)),
			},
		}
	}
	if (command == "invite" || command == "remove" || command == "admin" || command == "background" || command == "upgrade") && member.Role == entities.Member {
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T("MissingPermissions", translation.GetLocale(itc)),
			},
		}

	}
	switch command {
	case "create":
		name := itc.Data.Options[0].Options[0].Value.(string)
		name = rinha.Format(name)
		if name == "_test" {
			return nil
		}
		if len(name) >= 25 || len(name) <= 4 {
			return &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: translation.T("ClanNameLength", translation.GetLocale(itc)),
				},
			}
		}
		clan := database.Clan.GetClan(name)
		if clan.Name == "" {
			return &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: translation.T("ClanArleadyExists", translation.GetLocale(itc)),
				},
			}
		}
		database.User.UpdateUser(user.ID, func(u entities.User) entities.User {
			clan := database.Clan.GetClan(name)
			if clan.Name == "" {
				msg = "ClanAlreadyIn"
				return u
			}
			if 1000 > u.Money {
				msg = "ClanNotEnoughMoney"
				return u
			}
			u.Money -= 1000
			database.Clan.CreateClan(entities.Clan{
				Name: name,
			}, user.ID)
			return u
		})
		if msg != "" {
			return &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: translation.T(msg, translation.GetLocale(itc)),
				},
			}
		}
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T("ClanCreated", translation.GetLocale(itc), name),
			},
		}
	case "view":
		level := rinha.ClanXpToLevel(clan.Xp)
		memberMsg := ""
		for _, member := range clan.Members {
			user, err := handler.Client.User(disgord.Snowflake(member.ID)).Get()
			if err == nil {
				memberMsg += fmt.Sprintf("[**%s**] %s (**%d** XP)\n", member.Role.ToString(), user.Username, member.Xp)
			}
		}
		if len(memberMsg) > 1500 {
			memberMsg = memberMsg[:1500]
		}
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
				"membersText": memberMsg,
			}),
		}
		if bg != "" {
			embed.Image = &disgord.EmbedImage{
				URL: bg,
			}
		}
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Embeds: []*disgord.Embed{embed},
			},
		}
	case "mission":
		clan = rinha.PopulateClanMissions(clan)
		database.Clan.UpdateClan(clan, func(clanUpdate entities.Clan) entities.Clan {
			clanUpdate.Mission = clan.Mission
			clanUpdate.MissionProgress = clan.MissionProgress
			return clanUpdate
		})
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Embeds: []*disgord.Embed{{
					Title:       "Clan mission",
					Description: rinha.MissionToString(clan),
					Color:       65535,
				}},
			},
		}
	case "invite":
		user := utils.GetOptionsUser(itc.Data.Options[0].Options, itc, 0)
		if user == nil {
			return &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: "Invalid user",
				},
			}
		}
		text := translation.T("InviteToClan", translation.GetLocale(itc), map[string]interface{}{
			"clan":     clan.Name,
			"username": user.Username,
		})
		utils.Confirm(text, itc, user.ID, func() {
			database.Clan.UpdateClan(clan, func(c entities.Clan) entities.Clan {
				uClan := database.Clan.GetUserClan(user.ID)
				if uClan.Clan.Name != "" {
					msg = "UserArleadyInClan"
				} else {
					if len(c.Members) >= maxMembers {
						msg = "ClanMaxMembers"
					} else {
						database.User.GetUser(user.ID)
						database.Clan.InsertMember(&c, &entities.ClanMember{
							ID: user.ID,
						})
						msg = "SucessInvite"
					}
				}
				return c
			})
		})

		ch.CreateMessage(&disgord.CreateMessage{
			Content: translation.T(msg, translation.GetLocale(itc), user.Username),
		})
	case "remove":
		if len(itc.Data.Options[0].Options) == 0 {
			return &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: "Invalid user",
				},
			}
		}
		user := utils.GetOptionsUser(itc.Data.Options[0].Options, itc, 0)
		database.Clan.UpdateClan(clan, func(c entities.Clan) entities.Clan {
			uClan := database.Clan.GetUserClan(user.ID)
			if uClan.Clan.Name != clan.Name {
				msg = "UserNotInClan"
			} else {
				if member.Role <= uClan.Member.Role {
					msg = "NoPermission"
				} else {
					msg = "SucessRemove"
					database.Clan.RemoveMember(&c, user.ID)
				}
			}
			return c
		})

		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T(msg, translation.GetLocale(itc), user.Username),
			},
		}
	case "admin":
		if len(itc.Data.Options[0].Options) == 0 {
			return &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: "Invalid user",
				},
			}
		}
		user := utils.GetOptionsUser(itc.Data.Options[0].Options, itc, 0)
		database.Clan.UpdateClan(clan, func(c entities.Clan) entities.Clan {
			uClan := database.Clan.GetUserClan(user.ID)
			if uClan.Clan.Name != clan.Name {
				msg = "UserNotInClan"
			} else {
				if member.Role != entities.Owner || uClan.Member.Role == entities.Owner {
					msg = "NoPermission"
				} else {
					uClan.Member.Role = entities.Administrator
					database.Clan.UpdateMember(&c, uClan.Member)
					msg = "SucessAdmin"
				}
			}
			return c
		})
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T(msg, translation.GetLocale(itc), user.Username),
			},
		}
	case "banco":
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
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
		database.User.UpdateUser(user.ID, func(u entities.User) entities.User {
			if u.Money < money {
				msg = "NoMoney"
			} else {
				u.Money -= money
				database.Clan.UpdateClan(clan, func(c entities.Clan) entities.Clan {
					c.Money += money
					return c
				})
				msg = "ClanMoneyAdd"
			}
			return u
		})
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T(msg, translation.GetLocale(itc), money),
			},
		}
	case "upgrade":
		handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
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
		handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
			if len(ic.Data.Values) == 0 {
				return
			}
			upgrade := ic.Data.Values[0]
			if ic.Member.UserID == itc.Member.UserID {
				database.Clan.UpdateClan(clan, func(c entities.Clan) entities.Clan {
					var money int
					if upgrade == "membros" {
						money = rinha.CalcClanUpgrade(clan.MembersUpgrade, 1)
					}
					if upgrade == "missão" {
						money = rinha.CalcClanUpgrade(clan.MissionsUpgrade, 1)
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
				handler.Client.SendInteractionResponse(context.Background(), ic, &disgord.InteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.InteractionApplicationCommandCallbackData{
						Content: translation.T(msg, translation.GetLocale(itc), upgrade),
					},
				})
			}
		}, 100)

	case "background":
		img := itc.Data.Options[0].Options[0].Value.(string)
		if !utils.CheckImage(img) {
			return &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: translation.T("InvalidImage", translation.GetLocale(itc)),
				},
			}
		}
		database.Clan.UpdateClan(clan, func(c entities.Clan) entities.Clan {
			if c.Money < 10000 {
				msg = "ClanNoMoney"
			} else {
				c.Money -= 10000
				c.Background = img
				msg = "ChangedClanBackground"
			}
			return c
		})
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T(msg, translation.GetLocale(itc)),
			},
		}
	}
	return nil
}
