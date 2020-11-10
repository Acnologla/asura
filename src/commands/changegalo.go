package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"asura/src/database"
	"math/rand"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"changegalo", "galotype"},
		Run:       runChangeGalo,
		Available: true,
		Cooldown:  5,
		Usage:     "j!changegalo",
		Help:      "Mude o tipo do seu galo",
	})
}

func runChangeGalo(session disgord.Session, msg *disgord.Message, args []string) {
	galo := rinha.Galo{
		Changes: -1,
	}
	database.Database.NewRef(fmt.Sprintf("galo/%d",msg.Author.ID)).Get(context.Background(), &galo)
	if galo.Changes == 0 {
		galo.Changes = 3
	}
	if galo.Changes == 1{
		msg.Reply(context.Background(),session,msg.Author.Mention()+", Acabaram suas fichas de troca")
		return
	}
	galo.Changes--
	galoType := rand.Intn(len(rinha.Classes)-1) + 1
	galo.Type = galoType
	galo.Equipped = []int{}
	rinha.SaveGaloDB(msg.Author.ID, galo)
	msg.Reply(context.Background(),session,fmt.Sprintf("%s, Seu galo virou **%s**, lhe restam **%d** fichas de troca",msg.Author.Mention(),rinha.Classes[galo.Type].Name,galo.Changes - 1))
}

