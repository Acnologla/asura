package handler

import (
	"asura/src/telemetry"
	"context"
	"fmt"
	"github.com/agnivade/levenshtein"
	"github.com/andersfylling/disgord"
	"strconv"
	"strings"
)

// A struct that stores the information of a single "command" of the bot
type Command struct {
	Aliases    []string
	Run        func(disgord.Session, *disgord.Message, []string)
	Available  bool
}

// The place that will be stored all the commands
var Commands []Command = make([]Command, 0)
var Client *disgord.Client

//Register a new command into the big array of commands
func Register(command Command) {
	Commands = append(Commands, command)
}

// This function simply calculates the levenshtein distance between two strings
func CompareStrings(first string, second string) int {
	return levenshtein.ComputeDistance(first, second)
}

func handleCommand(session disgord.Session, msg *disgord.Message) {

	if msg.Author.Bot {
		return
	}

	myself, err := session.GetCurrentUser(context.Background())
	if err != nil {
		fmt.Printf("Unable to get current user info")
		return
	}

	botMention := fmt.Sprintf("<@!%d> ", uint64(myself.ID))

	if strings.HasPrefix(msg.Content, "j!") || strings.HasPrefix(msg.Content, "asura ") || strings.HasPrefix(msg.Content, botMention) {
		var raw string
		switch {
		case strings.HasPrefix(msg.Content, "j!"):
			raw = msg.Content[2:]
		case strings.HasPrefix(msg.Content, "asura "):
			raw = msg.Content[6:]
		case strings.HasPrefix(msg.Content, botMention):
			raw = strings.TrimPrefix(msg.Content, botMention)
		}
		splited := strings.Fields(raw)
		command := strings.ToLower(splited[0])
		args := splited[1:]
		var realCommand Command
		for _, cmd := range Commands {
			for _, alias := range cmd.Aliases {
				if alias == command {
					realCommand = cmd
					continue
				}
			}
		}
		for _, cmd := range Commands {
			if 1 >= CompareStrings(cmd.Aliases[0], command) {
				realCommand = cmd
				continue
			}
		}
		if len(realCommand.Aliases) > 0 {
			realCommand.Run(session, msg, args)
			tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
			telemetry.Info(fmt.Sprintf("Command %s used by %s", realCommand.Aliases[0], tag), map[string]string{
				"guild":   strconv.FormatUint(uint64(msg.GuildID), 10),
				"user":    strconv.FormatUint(uint64(msg.Author.ID), 10),
				"command": realCommand.Aliases[0],
				"channel": strconv.FormatUint(uint64(msg.ChannelID), 10),
			})
		}
	}
}

//Handles messages and call the functions that they have to execute.
func OnMessage(session disgord.Session, evt *disgord.MessageCreate) {
	msg := evt.Message
	handleCommand(session, msg)
}
