package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

var categorys = [...]string{translation.T("GeneralCommands", "pt"), translation.T("RinhaCommands", "pt"), translation.T("ProfileCommands", "pt"), translation.T("GameCommands", "pt")}

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "help",
		Description: translation.T("HelpHelp", "pt"),
		Run:         runHelp,
		Cooldown:    3,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "command",
			Type:        disgord.OptionTypeString,
			Description: translation.T("HelpOptionCommand", "pt"),
		}),
	})
}

func runHelp(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	if len(itc.Data.Options) > 0 {
		commandText := itc.Data.Options[0].Value.(string)
		command, ok := handler.Commands[commandText]
		if ok && !command.Dev {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Embeds: []*disgord.Embed{{
						Description: fmt.Sprintf("**%s**\n\nCooldown:\n **%d Segundos**", command.Description, command.Cooldown),
						Color:       65535,
						Title:       command.Name,
					}},
				},
			}
		} else {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("HelpCommandNotFound", translation.GetLocale(itc)),
				},
			}
		}
	} else {
		commandText := ""
		for i, category := range categorys {
			if i != 0 {
				commandText += "\n\n"
			}
			commandText += fmt.Sprintf("**%s**\n", category)
			for _, command := range handler.Commands {
				if command.Category == handler.CommandCategory(i) && !command.Dev {
					commandText += fmt.Sprintf("`%s` ", command.Name)
				}
			}
		}
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{{
					Description: commandText + "\n\nUse **/rinhahelp** para ver o tutorial de rinha de galo\n\n[Support](https://discord.gg/CfkBZyVsd7)\n[Vote](https://top.gg/bot/470684281102925844)\n[Website](https://acnologla.github.io/asura-site/)",
					Color:       65535,
					Footer: &disgord.EmbedFooter{
						Text: "Se tem alguma duvida não exite em entrar no meu servidor de suporte",
					},
				}},
			},
		}
	}
}
