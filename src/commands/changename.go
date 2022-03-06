package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "changename",
		Description: translation.T("ChangeNameHelp", "pt"),
		Run:         runChangeName,
		Cooldown:    10,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.OptionTypeString,
			Name:        "name",
			Required:    true,
			Description: translation.T("ChangeNameName", "pt"),
		}),
		Category: handler.Profile,
	})
}

func runChangeName(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	name := itc.Data.Options[0].Value.(string)
	if len(name) > 25 || 3 > len(name) {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("InvalidName", translation.GetLocale(itc)),
			},
		}
	}
	price := 100
	var msg string
	database.User.UpdateUser(itc.Member.UserID, func(u entities.User) entities.User {
		if rinha.IsVip(&u) {
			price = 0
		}
		if price > u.Money {
			msg = "NoMoney"
			return u
		}
		database.User.UpdateEquippedRooster(u, func(r entities.Rooster) entities.Rooster {
			r.Name = name
			return r
		})
		msg = "ChangeName"
		u.Money -= price
		return u
	}, "Galos")
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: translation.T(msg, translation.GetLocale(itc), name),
		},
	}
}
