package handler

import (
	"asura/src/interaction"
	"asura/src/translation"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/andersfylling/disgord"
)

const apiURL = "https://discord.com/api/v8"

var client = &http.Client{}

type Command struct {
	Name        string
	Cooldown    int
	Description string
	Options     []*interaction.ApplicationCommandOption
	Run         func(interaction.Interaction) *interaction.InteractionResponse
}

var Commands = map[string]Command{}

func RegisterCommand(command Command) {
	Commands[command.Name] = command
}

func Run(itc interaction.Interaction) *interaction.InteractionResponse {
	command := Commands[itc.Data.Name]
	if cooldown, ok := GetCooldown(itc.Member.User.ID, command); ok {
		needTime := command.Cooldown - int(time.Since(cooldown).Seconds())
		return &interaction.InteractionResponse{
			Type: interaction.CHANNEL_MESSAGE_WITH_SOURCE,
			Data: &interaction.InteractionCallbackData{
				Content: translation.T("Cooldown", itc.GuildLocale, needTime),
			},
		}
	}
	SetCooldown(itc.Member.User.ID, command)
	return command.Run(itc)
}

func findCommand(command string, commands []*interaction.ApplicationCommand) *interaction.ApplicationCommand {
	for _, c := range commands {
		if c.Name == command {
			return c
		}
	}
	return nil
}

func HasChanged(command *interaction.ApplicationCommand, realCommand Command) bool {
	if command.Description != realCommand.Description {
		return true
	}
	if len(command.Options) != len(realCommand.Options) {
		return true
	}
	for i, option := range command.Options {
		if option.Name != realCommand.Options[i].Name {
			return true
		}
		if option.Type != realCommand.Options[i].Type {
			return true
		}
		if option.Description != realCommand.Options[i].Description {
			return true
		}
		if option.Required != realCommand.Options[i].Required {
			return true
		}
		if option.MaxValue != realCommand.Options[i].MaxValue {
			return true
		}
		if option.MinValue != realCommand.Options[i].MinValue {

			return true
		}
		if option.AutoComplete != realCommand.Options[i].AutoComplete {
			return true
		}
	}
	return false
}

func Init(appID, token string, session *disgord.Client) {
	var commands []*interaction.ApplicationCommand
	endpoint := fmt.Sprintf("%s/applications/%s/commands", apiURL, appID)
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
		commandR := findCommand(name, commands)
		request := func(method string) {
			var newCommand interaction.ApplicationCommand
			newCommand.Name = command.Name
			newCommand.DefaultPermission = true
			newCommand.Type = interaction.CHAT_INPUT
			newCommand.Options = command.Options
			newCommand.Description = command.Description
			val, _ := json.Marshal(newCommand)
			reader := bytes.NewBuffer(val)
			_endpoint := endpoint
			if method == "PATCH" {
				_endpoint += fmt.Sprintf("/%d", commandR.ID)
			}
			req, _ := http.NewRequest(method, _endpoint, reader)
			req.Header.Add("Authorization", "Bot "+token)
			req.Header.Add("Content-Type", "application/json")
			_, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
		}
		if commandR == nil {
			request("POST")
		} else if HasChanged(commandR, command) {
			request("PATCH")
		}
	}
}
