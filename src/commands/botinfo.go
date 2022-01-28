package commands

import (
	"asura/src/handler"
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "botinfo",
		Description: translation.T("BotInfoHelp", "pt"),
		Run:         runBotInfo,
		Cooldown:    15,
	})
}

func runBotInfo(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	guildSize := len(handler.Client.GetConnectedGuilds())
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	ramUsage := memory.Alloc / 1000 / 1000
	guild, err := handler.Client.Guild(itc.GuildID).Get()
	guildUsage := 0
	if err == nil {
		guildUsage = int(unsafe.Sizeof(*guild))
		for _, member := range guild.Members {
			guildUsage += int(unsafe.Sizeof(*(member.User)))
			guildUsage += int(unsafe.Sizeof(*member))
		}
		for _, role := range guild.Roles {
			guildUsage += int(unsafe.Sizeof(*role))
		}
		for _, channel := range guild.Channels {
			guildUsage += int(unsafe.Sizeof(*channel))
		}
		for _, emoji := range guild.Emojis {
			guildUsage += int(unsafe.Sizeof(*emoji))
		}
		for _, presence := range guild.Presences {
			guildUsage += int(unsafe.Sizeof(*(presence.Game)))
			guildUsage += int(unsafe.Sizeof(*presence))
		}
		for _, voiceState := range guild.VoiceStates {
			guildUsage += int(unsafe.Sizeof(*voiceState))
		}
	}
	guildUsageText := fmt.Sprintf("Este servidor está usando %.2fKB de memória ram", float32(guildUsage)/1000)
	if guildUsage >= 1000000 {
		guildUsageText = fmt.Sprintf("Este servidor está usando %.2fMB de memória ram", float32(guildUsage)/1000/1000)
	}
	ping, _ := handler.Client.HeartbeatLatencies()
	if err != nil {
		return nil
	}
	botInfo, err := handler.Client.Gateway().GetBot()
	if err != nil {
		return nil
	}
	shard := disgord.ShardID(itc.GuildID, botInfo.Shards)
	myself, _ := handler.Client.Cache().GetCurrentUser()

	avatar, _ := myself.AvatarURL(512, true)
	date := ((uint64(myself.ID) >> 22) + 1420070400000) / 1000
	readyAt := int(time.Since(handler.ReadyAt).Minutes())
	freeWorkers := handler.GetFreeWorkers()
	description := translation.T("BotInfoDescription", translation.GetLocale(itc), map[string]interface{}{
		"workers":     handler.Workers,
		"freeWorkers": freeWorkers,
		"createdTime": int((uint64(time.Now().Unix()) - date) / 60 / 60 / 24),
		"guilds":      guildSize,
		"ram":         ramUsage,
		"ping":        ping[shard].Milliseconds(),
		"shards":      botInfo.Shards,
		"days":        readyAt / 60 / 24,
		"hours":       readyAt / 60 % 24,
		"minutes":     readyAt % 60,
	})
	res := &disgord.InteractionApplicationCommandCallbackData{
		Embeds: []*disgord.Embed{{
			Title: fmt.Sprintf("Asura (Shard %d)", disgord.ShardID(itc.GuildID, botInfo.Shards)),
			Color: 65535,
			Thumbnail: &disgord.EmbedThumbnail{
				URL: avatar,
			},
			Footer: &disgord.EmbedFooter{
				Text: guildUsageText,
			},
			Description: description,
		}},
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:  disgord.MessageComponentButton,
						Label: "Invite",
						Style: disgord.Link,
						Url:   "https://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%20bot&permissions=8",
					},
					{
						Type:  disgord.MessageComponentButton,
						Label: "Website",
						Style: disgord.Link,
						Url:   "https://acnologla.github.io/asura-site/",
					},
					{
						Type:  disgord.MessageComponentButton,
						Label: "Support",
						Style: disgord.Link,
						Url:   "https://discord.gg/tdVWQGV",
					},
					{
						Type:  disgord.MessageComponentButton,
						Label: "Vote",
						Style: disgord.Link,
						Url:   "https://top.gg/bot/470684281102925844",
					},
				},
			},
		},
	}
	return &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: res,
	}
}
