package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"time"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"missoes", "daily", "mission"},
		Run:       runMission,
		Available: true,
		Cooldown:  3,
		Usage:     "j!mission",
		Help:      "Veja suas missoes",
		Category:  1,
	})
}

func runMission(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author
	galo, _ := rinha.GetGaloDB(user.ID)
	text := rinha.MissionsToString(user.ID, galo)
	galo, _ = rinha.GetGaloDB(user.ID)
	embed := &disgord.Embed{
		Color:       65535,
		Title:       fmt.Sprintf("Missoes (%d/3)", len(galo.Missions)),
		Description: text,
	}
	if len(galo.Missions) != 3 {
		need := uint64(time.Now().Unix()) - galo.LastMission
		embed.Footer = &disgord.EmbedFooter{
			Text: fmt.Sprintf("Faltam %d horas e %d minutos para voce receber uma nova missão", 23-(need/60/60), 59-(need/60%60)),
		}
	}
	msg.Reply(context.Background(), session, embed)
}
