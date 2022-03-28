package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "chifresize",
		Description: translation.T("ChifreSizeHelp", "pt"),
		Run:         runChifreSize,
		Cooldown:    3,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.OptionTypeUser,
			Name:        "user",
			Description: translation.T("ChifreSizeUser", "pt"),
		}),
	})
}

func runChifreSize(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := utils.GetUserOrAuthor(itc, 0)
	idString := strconv.FormatUint(uint64(user.ID), 10)
	result, _ := strconv.Atoi(string(idString[3:4]))
	random, _ := strconv.Atoi(string(idString[5]))

	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Description: translation.T("ChifreSizeDescription", "pt", map[string]interface{}{
						"username":      user.Username,
						"height":        result * 3,
						"circumference": result + random,
					}),
					Color: 65535,
					Title: ":ox: " + user.Username,
				},
			},
		},
	}
}
