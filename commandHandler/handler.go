package commandHandler

import "asura/interaction"

type Command struct {
	Name string
	Run  func(interaction.Interaction)
}

var Commands = map[string]Command{}

func RegisterCommand(command Command) {
	Commands[command.Name] = command
}

func Run(itc interaction.Interaction) {
	Commands[itc.Data.Name].Run(itc)
}
