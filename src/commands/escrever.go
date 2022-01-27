package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"strings"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "escrever",
		Description: translation.T("EscreverHelp", "pt"),
		Run:         runEscrever,
		Cooldown:    3,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "text",
			Type:        disgord.OptionTypeString,
			Description: "message text",
			Required:    true,
		}),
	})
}

func runEscrever(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	str := itc.Data.Options[0].Value.(string)
	text := ""
	for i := 0; i < len(str); i++ {
		if i%2 == 0 {
			text += strings.ToUpper(string(str[i]))
		} else {
			text += strings.ToLower(string(str[i]))
		}
	}
	return &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Content: text,
		},
	}
}
