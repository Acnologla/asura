package handler

import (
	"asura/src/translation"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

const BOT_PREFIX = "j!"
const BOT_NAME_PREFIX = "asura"
const INTERACTION_FAKE_TOKEN = "m"

func extractOnlyNumbersFromString(str string) string {
	result := ""
	for _, char := range str {
		if char >= '0' && char <= '9' {
			result += string(char)
		}
	}
	return result
}

func getCommandNameAndArgsFromPrefix(m *disgord.Message) (string, []string) {
	splited := strings.Split(m.Content, " ")
	return strings.ToLower(splited[0][len(BOT_PREFIX):]), splited[1:]
}

func getCommandNameAndArgsFromBotName(m *disgord.Message) (string, []string) {
	splited := strings.Split(m.Content, " ")

	if len(splited) < 2 {
		return "", []string{}
	}

	return strings.ToLower(splited[1]), splited[2:]
}

func getCommandNameAndArgs(m *disgord.Message) (string, []string) {
	lowerContent := strings.ToLower(m.Content)

	if strings.HasPrefix(lowerContent, BOT_PREFIX) {
		return getCommandNameAndArgsFromPrefix(m)
	} else if strings.HasPrefix(lowerContent, BOT_NAME_PREFIX) {
		return getCommandNameAndArgsFromBotName(m)
	}

	return "", []string{}
}

func parseOptionMessage(option *disgord.ApplicationCommandOption, value string) (interface{}, error) {
	switch option.Type {
	case disgord.OptionTypeUser:
		id := extractOnlyNumbersFromString(value)
		_, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil, errors.New("invalid user")
		}
		return id, nil
	case disgord.OptionTypeSubCommand, disgord.OptionTypeSubCommandGroup:
		value := strings.ToLower(value)
		if value != option.Name {
			return nil, errors.New("invalid subcommand")
		}
		return value, nil
	case disgord.OptionTypeNumber:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, errors.New("invalid number")
		}
		if num < option.MinValue || num > option.MaxValue {
			return nil, errors.New("number out of range")
		}
		return num, nil
	case disgord.OptionTypeString:
		return value, nil
	case disgord.OptionTypeChannel:
		id := extractOnlyNumbersFromString(value)
		_, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil, errors.New("invalid channel")
		}
		return id, nil
	}
	return nil, nil
}

func GetOptionError(option *disgord.ApplicationCommandOption, err error) error {
	return fmt.Errorf("%s|%s", option.Name, err.Error())
}

func GetOptionsFromMessage(m *disgord.Message, opts []*disgord.ApplicationCommandOption, args []string) ([]*disgord.ApplicationCommandDataOption, error) {
	options := []*disgord.ApplicationCommandDataOption{}
	isSub := false
	i := 0
	for _, option := range opts {
		required := option.Required
		isSub = option.Type == disgord.OptionTypeSubCommand || option.Type == disgord.OptionTypeSubCommandGroup
		if required && len(args) <= i {
			return nil, GetOptionError(option, errors.New("missing required option"))
		}
		if len(args) > i {
			arg := args[i]
			value, err := parseOptionMessage(option, arg)
			if err != nil && required {
				return nil, GetOptionError(option, err)
			}
			if err == nil {
				if isSub && len(option.Options) > 0 {
					optNews, errNew := GetOptionsFromMessage(m, option.Options, args[i+1:])
					if errNew != nil {
						return nil, errNew
					}
					options = append(options, &disgord.ApplicationCommandDataOption{
						Name:    option.Name,
						Type:    option.Type,
						Options: optNews,
						Value:   value,
					})

				} else {
					options = append(options, &disgord.ApplicationCommandDataOption{
						Name:  option.Name,
						Type:  option.Type,
						Value: value,
					})
				}
				i++
			}
		}
	}
	if isSub && len(options) == 0 {
		return nil, GetOptionError(opts[0], errors.New("invalid subcommand"))
	}
	return options, nil
}

func createInteractionFromMessage(m *disgord.Message, command Command, commandName string, args []string) (*disgord.InteractionCreate, error) {
	options, err := GetOptionsFromMessage(m, command.Options, args)
	if err != nil {
		return nil, err
	}
	m.Member.User = m.Author
	return &disgord.InteractionCreate{
		ID:          m.ID,
		ChannelID:   m.ChannelID,
		GuildID:     m.GuildID,
		User:        m.Author,
		Member:      m.Member,
		Message:     m,
		Locale:      m.Author.Locale,
		GuildLocale: m.Author.Locale,
		Token:       INTERACTION_FAKE_TOKEN,
		Type:        disgord.InteractionApplicationCommand,
		Data: &disgord.ApplicationCommandInteractionData{
			Name:    commandName,
			Options: options,
		},
	}, nil
}

func optionsToString(options []*disgord.ApplicationCommandOption) string {
	result := ""
	for i, option := range options {
		isLast := i == len(options)-1
		if option.Type == disgord.OptionTypeSubCommandGroup {
			result += "<" + option.Name
			if len(option.Options) > 0 {
				result += " " + optionsToString(option.Options)
			}

			result += "> "

		} else if option.Type == disgord.OptionTypeSubCommand {
			result += option.Name
			if len(option.Options) > 0 {
				result += " " + optionsToString(option.Options)
			}
			if !isLast {
				result += "|"
			}
		} else {
			result += option.Name
			if !isLast {
				result += ", "
			}
		}
	}
	return result
}

func SendErrorMessage(m *disgord.Message, command Command, err error) {
	errorStr := err.Error()
	splited := strings.Split(errorStr, "|")
	optionName := splited[0]
	errorStr = strings.ReplaceAll(splited[1], " ", "")
	options := optionsToString(command.Options)
	msg := translation.T(errorStr, "pt", map[string]string{
		"option":  optionName,
		"options": options,
	})
	usageExample := translation.T("ExampleUsage", "pt", map[string]string{
		"commandName": command.Name,
		"prefix":      BOT_PREFIX,
		"options":     strings.ReplaceAll(options, ",", ""),
	})
	Client.Channel(m.ChannelID).CreateMessage(&disgord.CreateMessage{
		Content: msg + "\n" + usageExample,
	})
}

func ProcessMessage(m *disgord.Message) {
	if m.GuildID == 597089324114116635 || m.GuildID == 710179373860519997 {
		commandName, args := getCommandNameAndArgs(m)
		if commandName == "" {
			return
		}
		command := GetCommand(commandName)
		if command.Run == nil {
			return
		}
		interaction, err := createInteractionFromMessage(m, command, commandName, args)
		if err != nil {
			SendErrorMessage(m, command, err)
			return
		}

		InteractionChannel <- interaction
	}
}
