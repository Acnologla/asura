package handler

import (
	"asura/src/database"
	"asura/src/telemetry"
	"context"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/andersfylling/disgord"
)

const BaseWorkers = 128

var Workers = BaseWorkers

var DevMod = false
var DevIds = []disgord.Snowflake{365948625676795904, 395542409959964672, 364614407721844737}

func isDev(id disgord.Snowflake) bool {
	for _, devID := range DevIds {
		if id == devID {
			return true
		}
	}
	return false
}

var WorkersArray = make([]bool, Workers)

// A struct that stores the information of a single "command" of the bot
type Command struct {
	Aliases   []string
	Run       func(disgord.Session, *disgord.Message, []string)
	Help      string
	Usage     string
	Cooldown  int
	Available bool
	Category  int
}

// The place that will be stored all the commands
var Commands []Command = make([]Command, 0)
var CommandChannel = make(chan *disgord.Message)
var Client *disgord.Client
var ReadyAt = time.Now()
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
		msg.Reply(context.Background(), session, fmt.Sprintf("Voce tem que esperar %.1f segundos para usar esse comando denovo!", time_until))
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

func HandleCommand(session disgord.Session, msg *disgord.Message) {
	if msg.Author == nil {
		return
	}
	if msg.Author.Bot || uint64(msg.GuildID) == 0 {
		return
	}

	myself, err := Client.Cache().GetCurrentUser()
	if err != nil {
		fmt.Printf("Unable to get current user info")
		return
	}
	onlyBotMention := regexp.MustCompile(fmt.Sprintf(`<@(\!?)%d>`, uint64(myself.ID)))
	botMention := onlyBotMention.FindString(msg.Content)
	if strings.HasPrefix(strings.ToLower(msg.Content), "j!") || strings.HasPrefix(strings.ToLower(msg.Content), "asura ") || (botMention != "" && strings.HasPrefix(msg.Content, botMention)) {
		// I dont know if it's efficient but this is the easiest way to remove one of the three prefixes from the command.
		var raw string
		switch {
		case strings.HasPrefix(strings.ToLower(msg.Content), "j!"):
			raw = msg.Content[2:]
		case strings.HasPrefix(strings.ToLower(msg.Content), "asura "):
			raw = msg.Content[6:]
		case strings.HasPrefix(msg.Content, botMention):
			raw = strings.TrimPrefix(msg.Content, botMention)
			msg.Mentions = msg.Mentions[1:]
		}
		splited := strings.Fields(raw)
		if len(splited) == 0 {
			if strings.HasPrefix(msg.Content, botMention) {
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Meu prefix é **j!** use j!comandos para ver meus comandos")
			}
			return
		}
		command := strings.ToLower(splited[0])
		args := splited[1:]
		realCommand := FindCommand(command)

		if len(realCommand.Aliases) > 0 {
			if DevMod {
				if !isDev(msg.Author.ID) {
					return
				}
			}
			if database.IsBanned(msg.Author.ID) {
				return
			}
			if checkCooldown(session, msg, realCommand.Aliases[0], msg.Author.ID, realCommand.Cooldown) {
				return
			}
			tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
			telemetry.Info(fmt.Sprintf("Command %s used by %s", realCommand.Aliases[0], tag), map[string]string{
				"guild":   strconv.FormatUint(uint64(msg.GuildID), 10),
				"user":    strconv.FormatUint(uint64(msg.Author.ID), 10),
				"command": realCommand.Aliases[0],
				"content": msg.Content,
				"channel": strconv.FormatUint(uint64(msg.ChannelID), 10),
			})
			if !DevMod {
				defer func() {
					err := recover()
					if err != nil {
						stringError := string(debug.Stack())
						telemetry.Error(stringError, map[string]string{
							"guild":   strconv.FormatUint(uint64(msg.GuildID), 10),
							"user":    strconv.FormatUint(uint64(msg.Author.ID), 10),
							"content": msg.Content,
							"command": realCommand.Aliases[0],
							"channel": strconv.FormatUint(uint64(msg.ChannelID), 10),
						})
					}
				}()
			}
			Client.Channel(msg.ChannelID).TriggerTypingIndicator()
			realCommand.Run(session, msg, args)
		}
	}
}

//Handles messages and call the functions that they have to execute.
func OnMessage(session disgord.Session, evt *disgord.MessageCreate) {
	go SendMsg(evt.Message)
	CommandChannel <- evt.Message
}

//If you want to edit a message to make the command work again
func OnMessageUpdate(old *disgord.Message, new *disgord.Message) {
	if len(new.Embeds) == 0 && old.Pinned == new.Pinned {
		CommandChannel <- new
	}
}

func Worker(id int) {
	for msg := range CommandChannel {
		if id >= len(WorkersArray) {
			HandleCommand(Client, msg)
			return
		}
		WorkersArray[id] = true
		HandleCommand(Client, msg)
		WorkersArray[id] = false
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

func init() {
	if os.Getenv("PRODUCTION") == "" {
		DevMod = true
	}
	for i := 0; i < Workers; i++ {
		go Worker(i)
	}
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			freeWorkers := GetFreeWorkers()
			if int(Workers/5) >= freeWorkers {
				telemetry.Debug("Added workers...", map[string]string{})
				num := int(Workers / 2)
				WorkersArray = append(WorkersArray, make([]bool, num)...)
				for i := 0; i < num; i++ {
					go Worker(Workers + i)
				}
				Workers += num
			} else if Workers != BaseWorkers && (freeWorkers*100)/Workers >= 80 {
				num := int(Workers / 3)
				WorkersArray = WorkersArray[:Workers-num]
				Workers -= num
				telemetry.Debug("Removed extra workers...", map[string]string{})
			}
		}
	}()
}
