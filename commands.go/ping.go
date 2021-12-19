package commands

import (
	"asura/commandHandler"
	"asura/interaction"
)

func init() {
	commandHandler.RegisterCommand(commandHandler.Command{
		Name: "ping",
		Run:  runPing,
	})
}
func runPing(itc interaction.Interaction) {

}
