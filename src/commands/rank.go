package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"fmt"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

const DEFAULT_RANK_LIMIT = 15

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "rank",
		Description: translation.T("RankHelp", "pt"),
		Run:         runRank,
		Cooldown:    20,
		Category:    handler.Profile,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommandGroup,
				Name:        "usuario",
				Description: translation.T("RankUserHelp", "pt"),
				Options: utils.GenerateOptions(
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeSubCommand,
						Name:        "players",
						Description: translation.T("RankUserPlayersHelp", "pt"),
					},
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeSubCommand,
						Name:        "money",
						Description: translation.T("RankUserMoneyHelp", "pt"),
					},
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeSubCommand,
						Name:        "level",
						Description: translation.T("RankUserLevelHelp", "pt"),
					},
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeSubCommand,
						Name:        "vitorias",
						Description: translation.T("RankUserWinsHelp", "pt"),
					},
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeSubCommand,
						Name:        "derrotas",
						Description: translation.T("RankUserDefeatsHelp", "pt"),
					},
				),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommandGroup,
				Name:        "clan",
				Description: translation.T("RankClanHelp", "pt"),
				Options: utils.GenerateOptions(
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeSubCommand,
						Name:        "level",
						Description: translation.T("RankClanPlayersHelp", "pt"),
					},
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeSubCommand,
						Name:        "money",
						Description: translation.T("RankClanMoneyHelp", "pt"),
					}),
			},
		),
	})
}

var rankNameToType = map[string]string{
	"players":  "Andar da dungeon",
	"money":    "Dinheiro",
	"level":    "Nível",
	"vitorias": "Vitórias",
	"derrotas": "Derrotas",
}

func runRank(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	rankType := itc.Data.Options[0].Name
	rankName := itc.Data.Options[0].Options[0].Name
	var text string
	if rankType == "usuario" {
		var users []*entities.User
		var data func(u *entities.User) int
		switch rankName {
		case "players":
			users = database.User.SortUsers(DEFAULT_RANK_LIMIT, "dungeonreset", "dungeon")
			data = func(u *entities.User) int {
				return (u.DungeonReset * len(rinha.Dungeon)) + u.Dungeon
			}
		case "money":
			users = database.User.SortUsers(DEFAULT_RANK_LIMIT, "money")
			data = func(u *entities.User) int {
				return u.Money
			}
		case "level":
			users = database.User.SortUsersByRooster(DEFAULT_RANK_LIMIT, "resets", "xp")
			data = func(u *entities.User) int {
				return rinha.CalcLevel(u.Galos[0].Xp)
			}
		case "vitorias":
			users = database.User.SortUsers(DEFAULT_RANK_LIMIT, "win")
			data = func(u *entities.User) int {
				return u.Win
			}
		case "derrotas":
			users = database.User.SortUsers(DEFAULT_RANK_LIMIT, "lose")
			data = func(u *entities.User) int {
				return u.Lose
			}
		}
		for i, user := range users {
			u, err := handler.Client.User(user.ID).Get()
			if err == nil {
				tag := u.Username + "#" + u.Discriminator.String()
				extraMsg := ""
				if rankName == "level" {
					extraMsg = fmt.Sprintf(" (%d Reset)", user.Galos[0].Resets)
				}
				text += fmt.Sprintf("[**%d**] - %s\n%s: %d%s\n", i+1, tag, rankNameToType[rankName], data(user), extraMsg)
			}
		}

	} else {
		var clans []*entities.Clan
		var data func(u *entities.Clan) int
		switch rankName {
		case "level":
			clans = database.Clan.SortClan("xp", DEFAULT_RANK_LIMIT)
			data = func(u *entities.Clan) int {
				return rinha.ClanXpToLevel(u.Xp)
			}
		case "money":
			clans = database.Clan.SortClan("money", DEFAULT_RANK_LIMIT)
			data = func(u *entities.Clan) int {
				return u.Money
			}
		}
		for i, clan := range clans {
			text += fmt.Sprintf("[**%d**] - %s\n%s: %d\n", i+1, clan.Name, rankNameToType[rankName], data(clan))
		}
	}
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title:       rankName + " Rank",
					Color:       65535,
					Description: text,
				},
			},
		},
	}
}
