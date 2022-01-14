package commands

import (
	"asura/src/handler"
	"context"
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"botinfo", "infobot", "bi"},
		Run:       runBotinfo,
		Available: true,
		Cooldown:  15,
		Usage:     "j!botinfo",
		Help:      "Veja as informaçoes do bot",
	})
}

func runBotinfo(session disgord.Session, msg *disgord.Message, args []string) {
	guildSize := len(session.GetConnectedGuilds())
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	ramUsage := memory.Alloc / 1000 / 1000
	guild, err := session.Guild(msg.GuildID).Get()
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
	shard := disgord.ShardID(msg.GuildID, 1)
	myself, err := handler.Client.Cache().GetCurrentUser()
	if err != nil {
		return
	}
	botInfo, err := handler.Client.Gateway().GetBot()
	if err != nil {
		return
	}
	avatar, _ := myself.AvatarURL(512, true)
	date := ((uint64(myself.ID) >> 22) + 1420070400000) / 1000
	readyAt := int(time.Since(handler.ReadyAt).Minutes())
	freeWorkers := handler.GetFreeWorkers()
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Title: fmt.Sprintf("Asura (Shard %d)", disgord.ShardID(msg.GuildID, botInfo.Shards)),
			Color: 65535,
			Thumbnail: &disgord.EmbedThumbnail{
				URL: avatar,
			},
			Footer: &disgord.EmbedFooter{
				Text: guildUsageText,
			},
			Description: fmt.Sprintf(`
			Workers: **%d** (**%d** Livres)
			Bot criado a **%d** dias
			Servidores: **%d**
			Ram usada: **%d**MB
			Ping: **%dms**
			Shards: **%d**
			Bot online a %d dias %d horas %d minutos
			`, handler.Workers, freeWorkers, int((uint64(time.Now().Unix())-date)/60/60/24), guildSize, ramUsage, ping[shard].Milliseconds(), botInfo.Shards, readyAt/60/24, readyAt/60%24, readyAt%60),
		},
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:  disgord.MessageComponentButton,
						Label: "Convite",
						Style: disgord.Link,
						Url:   " https://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%20bot&permissions=8",
					},
					{
						Type:  disgord.MessageComponentButton,
						Label: "Website",
						Style: disgord.Link,
						Url:   "https://acnologla.github.io/asura-site/",
					},
					{
						Type:  disgord.MessageComponentButton,
						Label: "Suporte",
						Style: disgord.Link,
						Url:   "https://discord.gg/tdVWQGV",
					},
					{
						Type:  disgord.MessageComponentButton,
						Label: "Vote em mim",
						Style: disgord.Link,
						Url:   "https://top.gg/bot/470684281102925844",
					},
				},
			},
		},
	})
}
