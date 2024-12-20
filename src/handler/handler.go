package handler

import (
	"asura/src/cache"
	"asura/src/telemetry"
	"asura/src/translation"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	Run         func(context.Context, *disgord.InteractionCreate) *disgord.CreateInteractionResponse
	Category    CommandCategory
	Dev         bool
	Aliases     []string
	Cache       int
}

var Commands = map[string]Command{}

func RegisterCommand(command Command) {
	Commands[command.Name] = command
}

func ExecuteInteraction(ctx context.Context, interaction *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	if interaction.Type == disgord.InteractionApplicationCommand {
		return Run(ctx, interaction)
	}
	return nil
}

func GetCommand(name string) Command {
	command := Commands[name]
	if command.Name == "" {
		for _, cmd := range Commands {
			for _, alias := range cmd.Aliases {
				if alias == name {
					command = cmd
					break
				}
			}
		}
	}
	return command
}

func Run(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	cacheCommand := cache.GetCachedCommand(ctx, itc)
	if cacheCommand != nil {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: cacheCommand,
		}
	}
	command := GetCommand(itc.Data.Name)
	if command.Run == nil {
		return nil
	}
	locale := translation.GetLocale(itc)
	if cooldown, ok := GetCooldown(ctx, itc.Member.User.ID, command); ok {
		needTime := command.Cooldown - int(time.Since(cooldown).Seconds())
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("Cooldown", locale, needTime),
			},
		}
	}
	SetCooldown(ctx, itc.Member.User.ID, command)
	res := command.Run(ctx, itc)
	if command.Cache != 0 && res != nil {
		cache.CacheCommand(ctx, itc, res, command.Cache)
	}
	return res
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

		if option.MaxValue != realCommand.Options[i].MaxValue {
			return true
		}
		if option.MinValue != realCommand.Options[i].MinValue {

			return true
		}

		if option.Autocomplete != realCommand.Options[i].Autocomplete {
			return true
		}
		if len(option.Options) != len(realCommand.Options[i].Options) {
			return true
		}

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
			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			x, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(x))
		}
		if commandR == nil {
			fmt.Println("criando")

			request("POST", command.Name)
			for _, alias := range command.Aliases {
				request("POST", alias)
			}
		} else if HasChanged(commandR, command) {
			fmt.Println("atualizando")
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
		ctx := context.Background()
		response := ExecuteInteraction(ctx, itc)
		if response != nil {
			err := Client.SendInteractionResponse(ctx, itc, response)
			if err != nil {
				fmt.Println(err)
			}
		}
		author := itc.Member.User
		tag := author.Username + "#" + author.Discriminator.String()
		name := GetCommand(itc.Data.Name).Name
		telemetry.Info(fmt.Sprintf("Command %s used by %s", name, tag), map[string]string{
			"guild":   strconv.FormatUint(uint64(itc.GuildID), 10),
			"user":    strconv.FormatUint(uint64(author.ID), 10),
			"command": name,
			"channel": strconv.FormatUint(uint64(itc.ChannelID), 10),
		})
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
