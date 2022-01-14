package handler

import (
	"asura/src/interaction"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andersfylling/disgord"
)

const url = "https://discord.com/api/v8/applications/%s/commands"
const apiURL = "https://discord.com/api/v8"

var client = &http.Client{}

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

func Init(appID, token string, session *disgord.Client) {
	var commands []*interaction.ApplicationCommand
	endpoint := fmt.Sprintf("%s/applications/%s/commands", apiURL, appID)
	client := &http.Client{}
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("authorization", "Bot "+token)
	resp, err := client.Do(request)
	json.NewDecoder(resp.Body).Decode(&commands)
	if err != nil {
		log.Fatal(err)
	}
	for name, command := range Commands {
		if !isInApplicationCommandSlice(name, commands) {
			var newCommand interaction.ApplicationCommand
			newCommand.Name = command.Name
			newCommand.DefaultPermission = true
			newCommand.Type = interaction.CHAT_INPUT
			newCommand.Options = command.Options
			newCommand.Description = command.Description
			val, _ := json.Marshal(newCommand)
			reader := bytes.NewBuffer(val)
			req, _ := http.NewRequest("POST", endpoint, reader)
			req.Header.Add("Authorization", "Bot "+token)
			req.Header.Add("Content-Type", "application/json")
			_, err = client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
