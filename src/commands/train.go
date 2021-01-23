package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"train", "treino", "treinar"},
		Run:       runTrain,
		Available: true,
		Cooldown:  5,
		Usage:     "j!train",
		Category:  1,
		Help:      "Batalhe contra um galo",
	})
}

func runTrain(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Type == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem um galo, use j!galo para criar um")
		return
	}
	battleMutex.RLock()
	if currentBattles[msg.Author.ID] != "" {
		battleMutex.RUnlock()
		msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
		return
	}
	battleMutex.RUnlock()
	galoAdv := rinha.Galo{
		Xp:   galo.Xp,
		Type: rinha.GetRand(),
	}
	LockEvent(msg.Author.ID, "Clone de "+rinha.Classes[galoAdv.Type].Name)
	defer UnlockEvent(msg.Author.ID)
	winner, _ := ExecuteRinha(msg, session, rinhaOptions{
		galoAuthor:  &galo,
		galoAdv:     &galoAdv,
		authorName:  rinha.GetName(msg.Author.Username, galo),
		advName:     "Clone de " + rinha.Classes[galoAdv.Type].Name,
		authorLevel: rinha.CalcLevel(galo.Xp),
		advLevel:    rinha.CalcLevel(galoAdv.Xp),
	})
	rinha.CompleteMission(msg.Author.ID, galo, galoAdv, winner == 0, msg)
	if winner == 0 {
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       16776960,
			Title:       "Train",
			Description: fmt.Sprintf("Parabens %s, voce venceu", msg.Author.Username),
		})
	} else {
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       16711680,
			Title:       "Train",
			Description: fmt.Sprintf("Parabens %s, voce perdeu. Use j!train para treinar novamente", msg.Author.Username),
		})
	}

}
