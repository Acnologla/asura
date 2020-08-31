package commands

import (
	"asura/src/handler"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"galoupgrades", "upgrades"},
		Run:       runUpgrades,
		Available: true,
		Cooldown:  3,
		Usage:     "j!upgrades",
		Help:      "Veja os upgrades apra o seu galo",
	})
}

func runUpgrades(session disgord.Session, msg *disgord.Message, args []string) {
	var text string
	for key, value := range upgrades {
		text += fmt.Sprintf("[Level %d] - %s\n", value, key)
	}
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Title:       "Upgrades",
			Color:       65535,
			Description: text,
		},
	})
}
