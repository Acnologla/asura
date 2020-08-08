package commands

import (
	"asura/src/handler"
	"fmt"
	"context"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:    []string{"help","ajuda","comandos","cmds"},
		Run:        runHelp,
		Available:  true,
		Cooldown: 1,
		Usage: "j!help <comando>",
		Help: "Veja meus comandos ou informações sobre um comando",
	})
}

func runHelp(session disgord.Session, msg *disgord.Message, args []string) {
	if len(args) > 0{
		command := handler.FindCommand(args[0])
		aliasesText := ""
		for _,aliase := range command.Aliases{
			aliasesText+= fmt.Sprintf("`%s` ",aliase)
		}
		if (len(command.Aliases) >0 ){
			msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Description: fmt.Sprintf("**%s**\n\nUso:\n **%s**\n\nOutras alternativas:\n %s",command.Help, command.Usage,aliasesText),
					Color: 65535,
					Title: command.Aliases[0],
				},
			})
		}else {
			msg.Reply(context.Background(),session, "Não achei esse comando, use j!comandos para ver meus comandos")
		}
	}else{
		commandText := ""
		for _,command := range handler.Commands{
			commandText += fmt.Sprintf("`%s` ",command.Aliases[0])
		}
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Description: commandText,
				Footer: &disgord.EmbedFooter{
					Text: "Use j!help <comando> para ver informaçoes sobre um comando",
				},
				Color: 65535,
				Title: "Comandos",
			},
		})
	}
}
