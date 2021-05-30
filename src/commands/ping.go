package commands

import (
	"asura/src/handler"
	"context"
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

func runPing(session disgord.Session, msg *disgord.Message, args []string) {
	ping, _ := handler.Client.HeartbeatLatencies()
	botInfo, _ := handler.Client.Gateway().GetBot()
	shard := disgord.ShardID(msg.GuildID, botInfo.Shards)
	msg.Reply(context.Background(), session, fmt.Sprintf("%dms", ping[shard].Milliseconds()))

}
