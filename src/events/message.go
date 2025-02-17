package events

import (
	"asura/src/commands"
	"asura/src/database"
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

var MessageChannel = make(chan *disgord.MessageCreate)

const NUMBER_OF_PROCESSING_MESSAGE_WORKERS = 1024

type GuildInfo struct {
	sync.Mutex
	NewLootBoxTime int64
	LastUser       string
}

var cache = map[string]*GuildInfo{}

func GetGuildInfo(guildID string) *GuildInfo {
	if cache[guildID] == nil {
		cache[guildID] = &GuildInfo{}
	}
	return cache[guildID]
}

func IsFlood(msg *disgord.Message, cache *GuildInfo) bool {
	if cache.LastUser == msg.Author.ID.String() {
		return true
	}
	cache.LastUser = msg.Author.ID.String()
	return false
}

func setNewLootboxTime(cache *GuildInfo, now int64) {
	randomMinutes := utils.RandInt(500)
	cache.NewLootBoxTime = now + 60*60*10 + int64(randomMinutes)*60
}

const MIN_MEMBERS = 15

func RecieveLootbox(msg *disgord.Message) {
	guildDb := database.Guild.GetGuild(context.Background(), msg.GuildID)
	if guildDb.DisableLootbox || (guildDb.LootBoxChannel != 0 && guildDb.LootBoxChannel != msg.ChannelID) {
		return
	}
	cache := GetGuildInfo(msg.GuildID.String())
	now := time.Now().Unix()
	guild, _ := handler.Client.Cache().GetGuild(msg.GuildID)
	members := guild.MemberCount
	randomNumber := rand.Intn(100 + int(members/100))
	if msg.GuildID.String() == "710179373860519997" {
		randomNumber = rand.Intn(20)
	}
	if members > MIN_MEMBERS {
		isFlood := IsFlood(msg, cache)
		if !isFlood && randomNumber < 3 && now > cache.NewLootBoxTime {
			setNewLootboxTime(cache, now)
			go commands.SendLootbox(msg)
		}
	}
}

func HandleMessage(h *disgord.MessageCreate) {
	msg := h.Message
	appID := os.Getenv("APP_ID")
	if !msg.Author.Bot {
		if msg.GuildID != 0 {
			for _, user := range msg.Mentions {
				if user.ID.String() == appID {
					msg.Reply(context.Background(), handler.Client, "Meu prefix é **j!**\n Use **/help** ou **j!help** para ver meus comandos\nUse **/rinhahelp** para ver o **tutorial de rinha**")
					break
				}
			}

			RecieveLootbox(msg)

			handler.ProcessMessage(msg)

		}
	}
}

func Worker(id int) {
	for message := range MessageChannel {
		HandleMessage(message)
	}
}

func init() {
	for i := 0; i < NUMBER_OF_PROCESSING_MESSAGE_WORKERS; i++ {
		go Worker(i)
	}
}
