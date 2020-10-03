package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"topwins", "topgalowins"},
		Run:       runTopWins,
		Available: true,
		Cooldown:  15,
		Usage:     "j!topwins",
		Help:      "Veja os galos com mais vitorias",
	})
}

func runTopWins(session disgord.Session, msg *disgord.Message, args []string) {
	uid := strconv.FormatUint(uint64(msg.Author.ID), 10)
	q := database.Database.NewRef("galo").OrderByChild("win")
	result, err := q.GetOrdered(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var text string
	var myPos int
	for i := len(result) - 1; 0 <= i; i-- {
		if i > len(result)-11 {
			var gal rinha.Galo
			if err := result[i].Unmarshal(&gal); err != nil {
				continue
			}
			if result[i].Key() == uid {
				myPos = len(result) - i
			}
			converted, _ := strconv.Atoi(result[i].Key())
			user, err := session.GetUser(context.Background(), disgord.NewSnowflake(uint64(converted)))
			var username string
			if err != nil {
				username = "Anonimo"
			} else {
				username = user.Username + "#" + user.Discriminator.String()
			}
			text += fmt.Sprintf("[%d] - %s\nVitorias: **%d** | Derrotas: **%d**\n", len(result)-i, username, gal.Win, gal.Lose)
		} else {
			if result[i].Key() == uid {
				myPos = len(result) - i
			}
			if myPos != 0 {
				i = -1
			}
		}
	}
	avatar, _ := msg.Author.AvatarURL(128, true)
	var footer string
	if myPos == 0 {
		footer = `Voce não jogou nenhuma partida de brigadegalo`
	} else {
		footer = fmt.Sprintf("Voce esta na posição %d", myPos)
	}
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Description: text,
			Footer: &disgord.EmbedFooter{
				Text:    footer,
				IconURL: avatar,
			},
			Color: 65535,
			Title: "Topwins",
		},
	})

}
