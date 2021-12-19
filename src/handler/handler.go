package handler

import "asura/src/interaction"

type Command struct {
	Name        string
	Description string
	Options     []*interaction.ApplicationCommandInteractionDataOption
	Run         func(interaction.Interaction) *interaction.InteractionResponse
}

var Commands = map[string]Command{}

func RegisterCommand(command Command) {
	Commands[command.Name] = command
}

func Run(itc interaction.Interaction) *interaction.InteractionResponse {
	return Commands[itc.Data.Name].Run(itc)
}
