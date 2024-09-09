package handler

import (
	"asura/src/cache"
	"asura/src/telemetry"
	"asura/src/translation"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

const apiURL = "https://discord.com/api/v10"

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

func Run(ctx context.Context, interaction *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	cached := cache.GetCachedCommand(ctx, interaction)

	if cached != nil {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: cached,
		}
	}

	command := GetCommand(interaction.Data.Name)

	if command.Run == nil {
		return nil
	}

	locale := translation.GetLocale(interaction)

	if cooldown, ok := GetCooldown(ctx, interaction.Member.User.ID, command); ok {
		since := int(time.Since(cooldown).Seconds())
		needTime := command.Cooldown - since

		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("Cooldown", locale, needTime),
			},
		}
	}

	SetCooldown(ctx, interaction.Member.User.ID, command)

	resp := command.Run(ctx, interaction)

	if command.Cache != 0 && resp != nil {
		cache.CacheCommand(ctx, interaction, resp, command.Cache)
	}

	return resp
}

func HasChanged(command *disgord.ApplicationCommand, realCommand Command) bool {
	if command.Description != realCommand.Description || len(command.Options) != len(realCommand.Options) {
		return true
	}

	for i, option := range command.Options {
		realOption := realCommand.Options[i]
		if option.Name != realOption.Name || option.Type != realOption.Type ||
			option.Description != realOption.Description || option.Required != realOption.Required ||
			option.MaxValue != realOption.MaxValue || option.MinValue != realOption.MinValue ||
			option.Autocomplete != realOption.Autocomplete {
			return true
		}
	}

	return false
}

func createDiscordCommand(name string, command Command) disgord.ApplicationCommand {
	return disgord.ApplicationCommand{
		Name:              name,
		DefaultPermission: true,
		Type:              disgord.ApplicationCommandChatInput,
		Options:           command.Options,
		Description:       command.Description,
	}
}

func getExistingCommands(endpoint, token string) []*disgord.ApplicationCommand {
	var commands []*disgord.ApplicationCommand
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

	return commands
}

func uploadCommands(endpoint, token string, commands []disgord.ApplicationCommand) {
	cmds, _ := json.Marshal(commands)
	reader := bytes.NewBuffer(cmds)
	req, _ := http.NewRequest("PUT", endpoint, reader)

	req.Header.Add("Authorization", "Bot "+token)
	req.Header.Add("Content-Type", "application/json")

	_, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
}

func clone[K comparable, V any](og map[K]V) map[K]V {
	cloned := make(map[K]V)

	for key, value := range og {
		cloned[key] = value
	}

	return cloned
}

func Init(appID, token string, session *disgord.Client) {
	var newCommands []disgord.ApplicationCommand

	endpoint := fmt.Sprintf("%s/applications/%s/commands", apiURL, appID)
	commands := getExistingCommands(endpoint, token)
	cached := clone(Commands)
	needUpload := false

	for _, command := range commands {
		realCommand, hasCommand := Commands[command.Name]

		if hasCommand {
			if HasChanged(command, realCommand) {
				needUpload = true

				newCommands = append(newCommands, createDiscordCommand(command.Name, realCommand))
			}

			delete(cached, command.Name)
		}
	}

	if len(cached) != 0 {
		needUpload = true

		for name, command := range Commands {
			newCommands = append(newCommands, createDiscordCommand(name, command))
		}
	}

	if needUpload {
		fmt.Println("Upando comandos novos!")
		uploadCommands(endpoint, token, newCommands)
	}

	Client = session
}

func HandleInteraction(interaction *disgord.InteractionCreate) {
	if interaction.Member == nil {
		return
	}

	if interaction.Type == disgord.InteractionApplicationCommand {
		ctx := context.Background()
		resp := ExecuteInteraction(ctx, interaction)

		if resp != nil {
			err := Client.SendInteractionResponse(ctx, interaction, resp)

			if err != nil {
				fmt.Println(err)
			}
		}

		author := interaction.Member.User
		command := GetCommand(interaction.Data.Name)

		tag := fmt.Sprintf("%s#%s", author.Username, author.Discriminator.String())
		message := fmt.Sprintf("Command %s used by %s", command.Name, tag)

		telemetry.Info(message, map[string]string{
			"guild":   strconv.FormatUint(uint64(interaction.GuildID), 10),
			"channel": strconv.FormatUint(uint64(interaction.ChannelID), 10),
			"user":    strconv.FormatUint(uint64(author.ID), 10),
			"command": command.Name,
		})
	} else if interaction.Type == disgord.InteractionMessageComponent {
		ComponentInteraction(Client, interaction)
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
