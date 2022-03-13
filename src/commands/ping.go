package commands

import (
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
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: fmt.Sprintf("%dms", ping[shard].Milliseconds()),
		},
	}
}
