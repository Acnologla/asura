package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/andersfylling/disgord"
)

const endpoint = "https://discord.com/api/v8"

var AppID uint64

func getAppCommands() []*disgord.ApplicationCommand {
	var arr []*disgord.ApplicationCommand
	endpoint := fmt.Sprintf("%s/applications/%d/commands", endpoint, AppID)
	client := &http.Client{}
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return arr
	}
	request.Header.Set("authorization", "Bot "+os.Getenv("TOKEN"))
	resp, err := client.Do(request)
	if err != nil {
		return arr
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&arr)
	return arr
}

func HaveCommand(commands []*disgord.ApplicationCommand, command Command) bool {
	for _, appCommand := range commands {
		if appCommand.Name == command.Aliases[0] {
			return true
		}
	}
	return false
}

func GetChoices(choices ...*disgord.ApplicationCommandOptionChoice) []*disgord.ApplicationCommandOptionChoice {
	return choices
}

func CreateMessage(message *disgord.CreateMessageParams) *disgord.Message {
	msg := &disgord.Message{}
	msg.Embeds = []*disgord.Embed{message.Embed}
	msg.Content = message.Content
	msg.Components = message.Components
	return msg
}

func CreateMessageContent(content string) *disgord.Message {
	msg := &disgord.Message{}
	msg.Content = content
	return msg
}

func CreateMessageEmbed(embed *disgord.Embed) *disgord.Message {
	msg := &disgord.Message{}
	msg.Embeds = []*disgord.Embed{embed}
	return msg
}

func GetUserArg(args []*disgord.ApplicationCommandDataOption, i int) (*disgord.User, error) {
	if i >= len(args) {
		return nil, errors.New("too many arguments")
	}
	val := args[i]
	UserValue, ok := val.Value.(*disgord.User)
	if ok {
		return UserValue, nil
	}
	return nil, errors.New("invalid user")
}
func GetStringArg(args []*disgord.ApplicationCommandDataOption, i int) (string, error) {
	if i >= len(args) {
		return "", errors.New("too many arguments")
	}
	val := args[i]
	stringValue, ok := val.Value.(string)
	if ok {
		return stringValue, nil
	}
	return "", errors.New("invalid string")
}
func GetIntArg(args []*disgord.ApplicationCommandDataOption, i int) (int, error) {
	if i >= len(args) {
		return 0, errors.New("too many arguments")
	}
	val := args[i]
	intValue, ok := val.Value.(int)
	if ok {
		return intValue, nil
	}
	return 0, errors.New("invalid number")
}
func GetOptions(options ...*disgord.ApplicationCommandOption) []*disgord.ApplicationCommandOption {
	return options
}

func getOptionName(command Command, i int) string {
	if i >= len(command.Options) {
		return "No name"
	}
	return command.Options[i].Name
}

func handleSlash(interaction *disgord.InteractionCreate) {
	fmt.Println(2)
	if interaction.Member != nil {
		fmt.Println(3)
		if !interaction.Member.User.Bot {
			fmt.Println(4)
			msg := &disgord.Message{
				Author:  interaction.Member.User,
				GuildID: interaction.GuildID,
			}
			CommandChannel <- &Entry{
				Msg:         msg,
				IsSlash:     true,
				Interaction: interaction,
				Command:     FindCommand(interaction.Data.Name),
			}
		}
	}
}

func InitSlashCommands() {
	appid := os.Getenv("APP_ID")
	val, err := strconv.ParseUint(appid, 10, 64)
	if err != nil {
		return
	}
	AppID = val
	commands := getAppCommands()
	/*if lettern(commands) == 0 {
		return
	}*/
	for _, command := range Commands {
		if !HaveCommand(commands, command) && !command.NoSlash {
			err := Client.ApplicationCommand(disgord.Snowflake(AppID)).Global().Create(&disgord.ApplicationCommand{
				Type:        disgord.ApplicationCommandChatInput,
				Name:        command.Aliases[0],
				Options:     command.Options,
				Description: command.Help,
			})
			fmt.Println(err)
		}
	}
}
