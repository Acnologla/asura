package handler

import (
	"asura/src/translation"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/andersfylling/disgord"
)

const Workers = 128

var WorkersArray = make([]bool, Workers)

var InteractionChannel = make(chan *disgord.InteractionCreate)

var ReadyAt = time.Now()

type CommandCategory int

const (
	General CommandCategory = iota
	Rinha
	Profile
	Games
)

const apiURL = "https://discord.com/api/v8"

var Client *disgord.Client
var client = &http.Client{}

type Command struct {
	Name        string
	Cooldown    int
	Description string
	Options     []*disgord.ApplicationCommandOption
	Run         func(*disgord.InteractionCreate) *disgord.InteractionResponse
	Category    CommandCategory
	Dev         bool
	Aliases     []string
}

var Commands = map[string]Command{}

func RegisterCommand(command Command) {
	Commands[command.Name] = command
}

func ExecuteInteraction(interaction *disgord.InteractionCreate) *disgord.InteractionResponse {
	if interaction.Type == disgord.InteractionApplicationCommand {
		return Run(interaction)
	}
	return nil
}

func Run(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	command := Commands[itc.Data.Name]
	if command.Name == "" {
		for _, cmd := range Commands {
			for _, alias := range cmd.Aliases {
				if alias == itc.Data.Name {
					command = cmd
					break
				}
			}
		}
	}
	locale := translation.GetLocale(itc)
	if cooldown, ok := GetCooldown(itc.Member.User.ID, command); ok {
		needTime := command.Cooldown - int(time.Since(cooldown).Seconds())
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T("Cooldown", locale, needTime),
			},
		}
	}
	SetCooldown(itc.Member.User.ID, command)
	return command.Run(itc)
}

func findCommand(command string, commands []*disgord.ApplicationCommand) *disgord.ApplicationCommand {
	for _, c := range commands {
		if c.Name == command {
			return c
		}
	}
	return nil
}

func HasChanged(command *disgord.ApplicationCommand, realCommand Command) bool {
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
		/*
			if option.MaxValue != realCommand.Options[i].MaxValue {
				return true
			}
			if option.MinValue != realCommand.Options[i].MinValue {

				return true
			}
			if option.AutoComplete != realCommand.Options[i].AutoComplete {
				return true
			}
		*/
	}
	return false
}

func Init(appID, token string, session *disgord.Client) {
	var commands []*disgord.ApplicationCommand
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
		request := func(method string, name string) {
			var newCommand disgord.ApplicationCommand
			newCommand.Name = name
			newCommand.DefaultPermission = true
			newCommand.Type = disgord.ApplicationCommandChatInput
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
			request("POST", command.Name)
			for _, alias := range command.Aliases {
				request("POST", alias)
			}
		} else if HasChanged(commandR, command) {
			request("PATCH", command.Name)
		}
	}
	Client = session
}

func HandleInteraction(itc *disgord.InteractionCreate) {
	if itc.Member == nil {
		return
	}
	if itc.Type == disgord.InteractionApplicationCommand {
		response := ExecuteInteraction(itc)
		if response != nil {
			Client.SendInteractionResponse(context.Background(), itc, response)
		}
	} else if itc.Type == disgord.InteractionMessageComponent {
		ComponentInteraction(Client, itc)
	}
}

func Worker(id int) {
	for interaction := range InteractionChannel {
		WorkersArray[id] = true
		HandleInteraction(interaction)
		WorkersArray[id] = false
	}
}

func init() {
	for i := 0; i < Workers; i++ {
		go Worker(i)
	}
}

func GetFreeWorkers() int {
	freeWorkers := 0
	for _, worker := range WorkersArray {
		if !worker {
			freeWorkers++
		}
	}
	return freeWorkers
}
