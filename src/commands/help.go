package commands

import (
	"asura/src/handler"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"help", "ajuda", "comandos", "cmds"},
		Run:       runHelp,
		Available: true,
		Cooldown:  1,
		Usage:     "j!help <comando>",
		Help:      "Veja meus comandos ou informações sobre um comando",
	})
}

func runHelp(session disgord.Session, msg *disgord.Message, args []string) {
	if len(args) > 0 {
		command := handler.FindCommand(args[0])
		aliasesText := ""
		if len(command.Aliases) > 1 {
			for _, aliase := range command.Aliases[1:] {
				aliasesText += fmt.Sprintf("`%s` ", aliase)
			}
		}
		if len(command.Aliases) > 0 {
			msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Description: fmt.Sprintf("**%s**\n\nCooldown:\n **%d Segundos**\n\nUso:\n **%s**\n\nOutras alternativas:\n %s", command.Help, command.Cooldown, command.Usage, aliasesText),
					Color:       65535,
					Title:       command.Aliases[0],
				},
			})
		} else {
			msg.Reply(context.Background(), session, "Não achei esse comando, use j!comandos para ver meus comandos")
		}
	} else {
		commandText := ""
		rinhaText := ""
		for _, command := range handler.Commands {
			if command.Available {
				if command.Category == 0 {
					commandText += fmt.Sprintf("`%s` ", command.Aliases[0])
				} else {
					rinhaText += fmt.Sprintf("`%s` ", command.Aliases[0])
				}
			}
		}
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Description: commandText + "\n\n**Comandos de rinha**\n" + rinhaText + "\n\n[**Servidor de Suporte**](https://discord.gg/tdVWQGV)\n[**Website**](https://acnologla.github.io/asura-site/)",
				Footer: &disgord.EmbedFooter{
					Text: "Use j!help <comando> para ver informações sobre um comando!",
				},
				Color: 65535,
				Title: "Comandos",
			},
		})
	}
}
