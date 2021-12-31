package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"asura/src/utils/rinha/engine"
	"context"
	"fmt"

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
	ids := []disgord.Snowflake{msg.Author.ID}
	usernames := []string{msg.Author.Username}
	if len(msg.Mentions) > 5 {
		return
	}
	for _, user := range msg.Mentions {
		if user.ID != msg.Author.ID && !user.Bot {
			ids = append(ids, user.ID)
			usernames = append(usernames, user.Username)
		}

	}
	arr := []utils.IDUsername{}
	for i, id := range ids {
		if i != 0 {
			arr = append(arr, utils.IDUsername{ID: id, Username: usernames[i]})
		}
	}
	callback := func() {
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
		advLevel := 200 + (100 * len(ids))
		advReset := len(ids) - 1
		ngaloAdv := rinha.Galo{
			Xp:        rinha.CalcXP(advLevel) + 1,
			Type:      32,
			GaloReset: advReset,
		}
		galo := galos[0]
		winner, _ := engine.ExecuteRinha(msg, session, engine.RinhaOptions{
			GaloAuthor:  galo,
			GaloAdv:     ngaloAdv,
			AuthorName:  rinha.GetName(msg.Author.Username, galo),
			AdvName:     "Boss",
			AuthorLevel: rinha.CalcLevel(galo.Xp),
			AdvLevel:    advLevel,
			Waiting:     galos,
			Usernames:   usernames,
			NoItems:     true,
		}, false)
		if winner == 0 {
			text := ""
			selected := -1
			num := utils.RandInt(100)
			if 1 > num {
				selected = utils.RandInt(len(ids))
				text = fmt.Sprintf("O %s ganhou o evento e ganhou um item", usernames[selected])
				rinha.UpdateGaloDB(ids[selected], func(galo rinha.Galo) (rinha.Galo, error) {
					if !rinha.IsIntInList(43, galo.Items) {
						galo.Items = append(galo.Items, 43)
					}
					return galo, nil
				})
			}
			msg.Reply(context.Background(), session, &disgord.Embed{
				Color:       16776960,
				Title:       "Boss",
				Description: fmt.Sprintf("Parabens voces venceram\n%s", text),
			})
		} else {
			msg.Reply(context.Background(), session, &disgord.Embed{
				Color:       16776960,
				Title:       "Boss",
				Description: "Parabens voces perderam",
			})
		}
	}
	if len(arr) == 0 {
		callback()
	} else {
		utils.ConfirmArray("Batalhar contra o boss", msg.ChannelID, arr, callback)
	}
}
