package handler

import (
	"asura/src/interaction"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const url = "https://discord.com/api/v8/applications/%s/commands"

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

func isInApplicationCommandSlice(command string, commands []*interaction.ApplicationCommand) bool {
	for _, c := range commands {
		if c.Name == command {
			return true
		}
	}
	return false
}

func Init(appID, token string) {
	resp, err := http.Get(fmt.Sprintf(url, appID))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	var commands []*interaction.ApplicationCommand
	err = json.NewDecoder(resp.Body).Decode(&commands)
	if err != nil {
		log.Fatal(err)
		return
	}
	for name, command := range Commands {
		if isInApplicationCommandSlice(name, commands) {
			var newCommand interaction.ApplicationCommand
			newCommand.Name = command.Name
			newCommand.Type = interaction.CHAT_INPUT
			newCommand.Options = command.Options
			newCommand.Description = command.Description
			val, _ := json.Marshal(newCommand)
			reader := bytes.NewReader(val)
			client := &http.Client{}
			req, _ := http.NewRequest("POST", url, reader)
			req.Header.Add("Authorization", "Bot "+token)
			req.Header.Add("Content-Type", "application/json")
			_, err = client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
