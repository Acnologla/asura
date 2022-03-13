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
		if !command.Dev {
			commands = append(commands, Command{
				Aliases:   command.Aliases,
				Help:      command.Description,
				Cooldown:  command.Cooldown,
				Available: !command.Dev,
				Category:  int(command.Category),
			})
		}
	}
	result, _ := json.Marshal(commands)
	file, _ := os.OpenFile("generate/commands.json", os.O_RDWR, 0644)
	defer file.Close()
	file.Write(result)
}
