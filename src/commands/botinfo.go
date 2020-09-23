package commands

import (
	"asura/src/handler"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"runtime"
	"time"
	"unsafe"
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
	guild, err := session.GetGuild(context.Background(), msg.GuildID)
	guildUsage := 0
	if err == nil {
		guildUsage = int(unsafe.Sizeof(*guild))
		for _, member := range guild.Members {
			guildUsage += int(unsafe.Sizeof(*member))
		}
		for _, role := range guild.Roles {
			guildUsage += int(unsafe.Sizeof(*role))
		}
		for _, emoji := range guild.Emojis {
			guildUsage += int(unsafe.Sizeof(*emoji))
		}
		for _, presence := range guild.Presences {
			guildUsage += int(unsafe.Sizeof(*presence))
		}
		for _, voiceState := range guild.VoiceStates {
			guildUsage += int(unsafe.Sizeof(*voiceState))
		}
	}
	guildUsageText := fmt.Sprintf("Este servidor esta usando %.2fKB de memoria ram", float32(guildUsage)/1000)
	if guildUsage >= 1000000 {
		guildUsageText = fmt.Sprintf("Este servidor esta usando %.2fMB de memoria ram", float32(guildUsage)/1000/1000)
	}
	ping, _ := handler.Client.HeartbeatLatencies()
	shard := disgord.ShardID(msg.GuildID, 1)
	myself, err := session.GetCurrentUser(context.Background())
	if err != nil {
		return
	}
	avatar, _ := myself.AvatarURL(512, true)
	date := ((uint64(myself.ID) >> 22) + 1420070400000) / 1000
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Title: "Asura",
			Color: 65535,
			Thumbnail: &disgord.EmbedThumbnail{
				URL: avatar,
			},
			Footer: &disgord.EmbedFooter{
				Text: guildUsageText,
			},
			Description: fmt.Sprintf(`
			Bot criado a **%d** dias
			Servidores: **%d**
			Ram usada: **%d**MB
			Ping: **%s**
			
			[**Convite**](https://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=bot&permissions=8)
			[**Servidor de suporte**](https://discord.gg/tdVWQGV)
			`, int((uint64(time.Now().Unix())-date)/60/60/24), guildSize, ramUsage, ping[shard].String()),
		},
	})
}
