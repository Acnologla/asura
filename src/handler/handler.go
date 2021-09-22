package handler

import (
	"asura/src/database"
	"asura/src/telemetry"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	Run       func(disgord.Session, *disgord.Message, []*disgord.ApplicationCommandDataOption) (*disgord.Message, func(disgord.Snowflake))
	Help      string
	Usage     string
	NoSlash   bool
	Options   []*disgord.ApplicationCommandOption
	Cooldown  int
	Available bool
	Category  int
}

type Entry struct {
	Msg         *disgord.Message
	IsSlash     bool
	Command     Command
	Interaction *disgord.InteractionCreate
}

// The place that will be stored all the commands
var Commands []Command = make([]Command, 0)
var CommandChannel = make(chan *Entry)
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
func checkCooldown(session disgord.Session, msg *disgord.Message, command string, id disgord.Snowflake, cooldown int) (bool, float64) {
	CooldownMutexes[command].RLock()
	if val, ok := Cooldowns[command][id]; ok {
		time_until := float64(cooldown) - time.Since(val).Seconds()
		CooldownMutexes[command].RUnlock()
		return true, time_until
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
	return false, 0
}

func ExecuteCommand(command Command, msg *disgord.Message, interaction *disgord.InteractionCreate) {
	fmt.Println(1)
	if len(command.Aliases) == 0 {
		return
	}
	args := interaction.Data.Options
	if DevMod {
		if !isDev(msg.Author.ID) {
			return
		}
	}
	if database.IsBanned(msg.Author.ID) {
		return
	}
	cooldown, time_until := checkCooldown(Client, msg, command.Aliases[0], msg.Author.ID, command.Cooldown)
	if cooldown {
		content := fmt.Sprintf("Voce tem que esperar %.1f segundos para usar esse comando denovo!", time_until)
		if interaction.ID == 0 {
			msg.Reply(context.Background(), Client, content)
		} else {
			Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
				Type: disgord.ChannelMessageWithSource,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Content: content,
				},
			})
		}
		return
	}
	tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
	telemetry.Info(fmt.Sprintf("Command %s used by %s", command.Aliases[0], tag), map[string]string{
		"guild":   strconv.FormatUint(uint64(msg.GuildID), 10),
		"user":    strconv.FormatUint(uint64(msg.Author.ID), 10),
		"command": command.Aliases[0],
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
					"command": command.Aliases[0],
					"channel": strconv.FormatUint(uint64(msg.ChannelID), 10),
				})
			}
		}()
	}
	Client.Channel(msg.ChannelID).TriggerTypingIndicator()
	message, followup := command.Run(Client, msg, args)
	if message == nil {
		return
	}
	if interaction.ID == 0 {
		newmsg, err := msg.Reply(context.Background(), Client, message)
		if followup != nil && err == nil {
			followup(newmsg.ID)
		}
	} else {
		Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
			Type: disgord.ChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Tts:        message.Tts,
				Content:    message.Content,
				Embeds:     message.Embeds,
				Components: message.Components,
			},
		})
		if followup != nil {
			var msg *disgord.Message
			endpoint := fmt.Sprintf("%s/webhooks/%d/%s/messages/@original", endpoint, AppID, interaction.Token)
			resp, err := http.Get(endpoint)
			if err == nil {
				defer resp.Body.Close()
				json.NewDecoder(resp.Body).Decode(&msg)
				followup(msg.ID)
			}
		}
	}
}

func getArgs(command Command, args []string) []*disgord.ApplicationCommandDataOption {
	val := []*disgord.ApplicationCommandDataOption{}
	for i, arg := range args {
		val = append(val, &disgord.ApplicationCommandDataOption{
			Name:  getOptionName(command, i),
			Type:  disgord.STRING,
			Value: arg,
		})
	}
	return val
}
func HandleCommand(session disgord.Session, msg *disgord.Message) {
	if msg.Author == nil {
		return
	}
	if msg.Author.Bot || msg.GuildID == 0 {
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

		realCommand := FindCommand(command)
		args := getArgs(realCommand, splited[1:])
		if len(realCommand.Aliases) > 0 {
			ExecuteCommand(realCommand, msg, &disgord.InteractionCreate{
				Data: &disgord.ApplicationCommandInteractionData{
					Options: args,
				},
			})
		}
	}
}

//Handles messages and call the functions that they have to execute.
func OnMessage(session disgord.Session, evt *disgord.MessageCreate) {
	go SendMsg(evt.Message)
	CommandChannel <- &Entry{
		Msg: evt.Message,
	}
}

//If you want to edit a message to make the command work again
func OnMessageUpdate(old *disgord.Message, new *disgord.Message) {
	if len(new.Embeds) == 0 && old.Pinned == new.Pinned {
		CommandChannel <- &Entry{
			Msg: new,
		}
	}
}

func Worker(id int) {
	for entry := range CommandChannel {
		fmt.Println(1)
		if id >= len(WorkersArray) {
			if entry.IsSlash {
				ExecuteCommand(entry.Command, entry.Msg, entry.Interaction)
			} else {
				HandleCommand(Client, entry.Msg)
			}
			return
		}
		WorkersArray[id] = true
		if entry.IsSlash {
			ExecuteCommand(entry.Command, entry.Msg, entry.Interaction)
		} else {
			HandleCommand(Client, entry.Msg)
		}
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
