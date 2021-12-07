package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"asura/src/utils/rinha/engine"
	"context"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"evento"},
		Run:       runNatal,
		Available: true,
		Cooldown:  5,
		Usage:     "j!evento",
		Category:  1,
		Help:      "Evento",
	})
}

func runNatal(session disgord.Session, msg *disgord.Message, args []string) {
	ids := []disgord.Snowflake{}
	usernames := []string{}
	if len(msg.Mentions) > 5 {
		return
	}
	for _, user := range msg.Mentions {
		if user.ID != msg.Author.ID && !user.Bot {
			ids = append(ids, user.ID)
			usernames = append(usernames, user.Username)
		}

	}

	battleMutex.Lock()
	for _, id := range ids {
		if currentBattles[id] != "" {
			battleMutex.Unlock()
			msg.Reply(context.Background(), session, "Um dos jogadores já esta lutando com o "+currentBattles[id])
			return
		}
	}
	galos := []rinha.Galo{}
	for _, id := range ids {
		gal, _ := rinha.GetGaloDB(id)
		if gal.Type == 0 {
			battleMutex.Unlock()
			return
		}
		galos = append(galos, gal)
	}
	for _, id := range ids {
		currentBattles[id] = "Boss"
	}
	battleMutex.Unlock()
	defer func() {
		for _, id := range ids {
			battleMutex.Lock()
			delete(currentBattles, id)
			battleMutex.Unlock()
		}
	}()

	ngaloAdv := rinha.Galo{
		Xp:   rinha.CalcXP(100) + 1,
		Type: 32,
	}
	galo := galos[0]
	winner, _ := engine.ExecuteRinha(msg, session, engine.RinhaOptions{
		GaloAuthor:  galo,
		GaloAdv:     ngaloAdv,
		AuthorName:  rinha.GetName(msg.Author.Username, galo),
		AdvName:     "Boss",
		AuthorLevel: rinha.CalcLevel(galo.Xp),
		AdvLevel:    100,
		Waiting:     galos,
		Usernames:   usernames,
		NoItems:     true,
	}, false)
	if winner == -1 {
		return
	}
}
