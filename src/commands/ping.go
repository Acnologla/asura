package commands

import (
	"asura/src/handler"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"ping"},
		Run:       runPing,
		Available: true,
		Cooldown:  1,
		Usage:     "j!ping",
		Help:      "Veja meu ping",
	})
}

func runPing(session disgord.Session, msg *disgord.Message, args []*disgord.ApplicationCommandDataOption) (*disgord.Message, func(disgord.Snowflake)) {
	fmt.Println(6)
	ping, _ := handler.Client.HeartbeatLatencies()
	botInfo, _ := handler.Client.Gateway().GetBot()
	shard := disgord.ShardID(msg.GuildID, botInfo.Shards)
	return handler.CreateMessageContent(fmt.Sprintf("%dms", ping[shard].Milliseconds())), nil
}
