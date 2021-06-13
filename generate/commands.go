package main

import (
	_ "asura/src/commands"
	"asura/src/handler"
	"encoding/json"
	"os"
)

type Command struct {
	Aliases   []string `json:"aliases"`
	Help      string   `json:"help"`
	Usage     string   `json:"usage"`
	Cooldown  int      `json:"cooldown"`
	Available bool     `json:"available"`
	Category  int      `json:"category"`
}

func main() {
	commands := []Command{}
	for _, command := range handler.Commands {
		commands = append(commands, Command{
			Aliases:   command.Aliases,
			Help:      command.Help,
			Usage:     command.Usage,
			Cooldown:  command.Cooldown,
			Available: command.Available,
			Category:  command.Category,
		})
	}
	result, _ := json.Marshal(commands)
	file, _ := os.OpenFile("generate/commands.json", os.O_RDWR, 0644)
	defer file.Close()
	file.Write(result)
}
