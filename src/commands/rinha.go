package commands

import (
	"asura/src/cache"
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"asura/src/translation"
	"asura/src/utils"
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/go-redis/redis/v8"
)

func isInRinha(user *disgord.User) string {
	val := cache.Client.Get(context.Background(), "/rinha/"+user.ID.String())
	result, _ := val.Result()
	return result
}

func lockBattle(author, adv disgord.Snowflake, authorUsername, advUsername string) {
	cache.Client.TxPipelined(context.Background(), func(p redis.Pipeliner) error {
		p.Set(context.Background(), "/rinha/"+author.String(), advUsername, time.Minute*4)
		p.Set(context.Background(), "/rinha/"+adv.String(), authorUsername, time.Minute*4)
		return nil
	})
}

func sendLevelUpEmbed(itc *disgord.InteractionCreate, galoWinner *entities.Rooster, user *disgord.User, xpOb int) {
	if rinha.CalcLevel(galoWinner.Xp) > rinha.CalcLevel(galoWinner.Xp-xpOb) {
		nextLevel := rinha.CalcLevel(galoWinner.Xp)
		nextSkill := rinha.GetNextSkill(*galoWinner)
		nextSkillStr := ""
		if len(nextSkill) != 0 {
			nextSkillStr = fmt.Sprintf("e desbloqueando a habilidade %s", nextSkill[0].Name)
		}
		ch, _ := handler.Client.Channel(itc.ChannelID).Get()
		if ch != nil {
			ch.SendMsg(context.Background(), handler.Client, &disgord.Message{
				Embeds: []*disgord.Embed{{
					Title:       "Galo upou de nivel",
					Color:       65535,
					Description: fmt.Sprintf("O galo de %s upou para o nivel %d %s", user.Username, nextLevel, nextSkillStr),
				}},
			})
		}
	}
}

func unlockBattle(author, adv disgord.Snowflake) {
	cache.Client.TxPipelined(context.Background(), func(p redis.Pipeliner) error {
		p.Del(context.Background(), "/rinha/"+author.String())
		p.Del(context.Background(), "/rinha/"+adv.String())
		return nil
	})
}
func rinhaMessage(username, text string) *disgord.InteractionResponse {
	return &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Content: translation.T("RinhaMessage", "pt", map[string]interface{}{"username": username, "text": text}),
		},
	}
}

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "rinha",
		Description: translation.T("RinhaHelp", "pt"),
		Run:         runRinha,
		Cooldown:    15,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user",
			Required:    true,
		}),
	})
}

func runRinha(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	user := utils.GetUser(itc, 0)
	if user.ID == itc.Member.User.ID || user.Bot {
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: "Invalid user",
			},
		}
	}
	author := itc.Member.User
	text := translation.T("RinhaDuel", translation.GetLocale(itc), map[string]interface{}{"author": author.Username, "adv": user.Username})
	utils.Confirm(text, itc, user.ID, func() {
		userRinha := isInRinha(user)
		if userRinha != "" {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(user.Username, userRinha).Data.Content,
			})
			return
		}
		authorRinha := isInRinha(itc.Member.User)
		if authorRinha != "" {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(author.Username, userRinha).Data.Content,
			})
			return
		}
		executePVP(itc)
	})
	return nil
}

func executePVP(itc *disgord.InteractionCreate) {
	author := itc.Member.User
	user := utils.GetUser(itc, 0)
	lockBattle(author.ID, user.ID, author.Username, user.Username)
	defer unlockBattle(author.ID, user.ID)
	authorDb := database.Client.GetUser(author.ID, "Galos", "Items")
	advDb := database.Client.GetUser(user.ID, "Galos", "Items")

	galoAuthor := rinha.GetEquippedGalo(&authorDb)
	galoAdv := rinha.GetEquippedGalo(&advDb)
	authorLevel := rinha.CalcLevel(galoAuthor.Xp)
	advLevel := rinha.CalcLevel(galoAdv.Xp)

	authorName := rinha.GetName(author.Username, *galoAuthor)
	advName := rinha.GetName(user.Username, *galoAdv)
	whoWin, battle := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
		GaloAuthor:  &authorDb,
		GaloAdv:     &advDb,
		AuthorName:  authorName,
		AdvName:     advName,
		AuthorLevel: authorLevel,
		AdvLevel:    advLevel,
		IDs:         [2]disgord.Snowflake{author.ID, user.ID},
	}, false)

	if whoWin == -1 {
		return
	}

	if 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life {
		winnerTurn := whoWin
		turn := 1
		winner := author
		loser := user

		if whoWin == 1 {
			winner = user
			loser = author
			turn = 0
		}

		galoWinner := battle.Fighters[winnerTurn].Galo
		galoLoser := battle.Fighters[turn].Galo
		if 2 >= rinha.CalcLevel(galoWinner.Xp)-rinha.CalcLevel(galoLoser.Xp) {
			database.Client.UpdateUser(loser.ID, func(u entities.User) entities.User {
				u.Lose++
				return u
			})
		}
		database.Client.UpdateUser(winner.ID, func(u entities.User) entities.User {
			u.Win++
			return u
		})
		embed := &disgord.Embed{Title: "Briga de galo", Color: 16776960, Description: ""}

		if winnerTurn == 1 {
			embed.Description = "\n" + translation.T("WinDuel", translation.GetLocale(itc), advName)
		} else {
			embed.Description = "\n" + translation.T("WinDuel", translation.GetLocale(itc), authorName)
		}
		ch, _ := handler.Client.Channel(itc.ChannelID).Get()
		if ch != nil {
			ch.SendMsg(context.Background(), handler.Client, &disgord.Message{
				Embeds: []*disgord.Embed{embed},
			})
		}
	}

}
