package handler

import (
	"asura/src/telemetry"
	"context"
	"fmt"
	"github.com/agnivade/levenshtein"
	"github.com/andersfylling/disgord"
	"strconv"
	"strings"
	"sync"
	"time"
)

// A struct that stores the information of a single "command" of the bot
type Command struct {
	Aliases   []string
	Run       func(disgord.Session, *disgord.Message, []string)
	Help      string
	Usage     string
	Cooldown  int
	Available bool
}

// The place that will be stored all the commands
var Commands []Command = make([]Command, 0)
var Client *disgord.Client

var Cooldowns = map[string]map[disgord.Snowflake]time.Time{}
var CooldownMutexes = map[string]*sync.RWMutex{}

//Register a new command into the big array of commands
func Register(command Command) {
	Cooldowns[command.Aliases[0]] = map[disgord.Snowflake]time.Time{}
	CooldownMutexes[command.Aliases[0]] = &sync.RWMutex{}
	Commands = append(Commands, command)
}

// This function simply calculates the levenshtein distance between two strings
func CompareStrings(first string, second string) int {
	return levenshtein.ComputeDistance(first, second)
}

func FindCommand(command string) (realCommand Command) {
	for _, cmd := range Commands {
		for _, alias := range cmd.Aliases {
			if alias == command {
				realCommand = cmd
				return
			}
		}
	}
	for _, cmd := range Commands {
		if 1 >= CompareStrings(cmd.Aliases[0], command) {
			realCommand = cmd
			continue
		}
	}
	return
}

// This hell of Mutexes are a simple way to sincronize all the goroutines around the Cooldown array.
func checkCooldown(session disgord.Session, msg *disgord.Message, command string, id disgord.Snowflake, cooldown int) bool {
	CooldownMutexes[command].RLock()
	if val, ok := Cooldowns[command][id]; ok {
		time_until := float64(cooldown) - time.Since(val).Seconds()
		msg.Reply(context.Background(), session, fmt.Sprintf("You have to wait %.2f seconds to use this command again", time_until))
		CooldownMutexes[command].RUnlock()
		return true
	}
	CooldownMutexes[command].RUnlock()
	CooldownMutexes[command].Lock()
	Cooldowns[command][id] = time.Now()
	CooldownMutexes[command].Unlock()
	go func() {
		time.Sleep(time.Duration(cooldown) * time.Second)
		CooldownMutexes[command].Lock()
		delete(Cooldowns[command], id)
		CooldownMutexes[command].Unlock()
	}()
	return false
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
		// I dont know if it's efficient but this is the easiest way to remove one of the three prefixes from the command.
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
		realCommand := FindCommand(command)

		if len(realCommand.Aliases) > 0 {
			if checkCooldown(session, msg, realCommand.Aliases[0], msg.Author.ID, realCommand.Cooldown) {
				return
			}
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

//If you want to edit a message to make the command work again
func OnMessageUpdate(session disgord.Session, evt *disgord.MessageUpdate) {
	if len(evt.Message.Embeds) == 0 {
		msg := evt.Message
		handleCommand(session, msg)
	}
}
