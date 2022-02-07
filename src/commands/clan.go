package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/translation"
	"asura/src/utils"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "clan",
		Description: translation.T("ClanHelp", "pt"),
		Run:         runClan,
		Cooldown:    10,
		Options: utils.GenerateOptions(
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
				Name:        "view",
				Description: translation.T("ClanViewHelp", "pt"),
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
				Name:        "admin",
				Description: translation.T("ClanAdmin", "pt"),
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeUser,
					Required:    true,
					Name:        "user",
					Description: "User admin",
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
				Name:        "upgrade",
				Description: translation.T("ClanUpgrade", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "banco",
				Description: translation.T("ClanBank", "pt"),
			},
		),
	})

}

func runClan(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	command := itc.Data.Options[0].Name
	user := database.User.GetUser(itc.Member.UserID)
	switch command {
	case "create":
		var msg string
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
	}
	return nil
}
