package commands

import (
	"asura/src/entities"
	"asura/src/handler"
	"fmt"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ping",
		Description: translation.T("PingHelp", "pt"),
		Run:         runPing,
		Cooldown:    3,
	})
}

func runPing(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	ping, _ := handler.Client.HeartbeatLatencies()
	botInfo, _ := handler.Client.Gateway().GetBot()
	shard := disgord.ShardID(itc.GuildID, botInfo.Shards)
	return entities.CreateMsg().
		Content(fmt.Sprintf("%dms", ping[shard].Milliseconds())).
		Res()
}
