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

func isInRinha(ctx context.Context, user *disgord.User) string {
	val := cache.Client.Get(ctx, "/rinha/"+user.ID.String())
	result, _ := val.Result()
	return result
}

func lockEvent(ctx context.Context, author disgord.Snowflake, evtUsername string) {
	cache.Client.Set(ctx, "/rinha/"+author.String(), evtUsername, time.Minute*4)
}

func unlockEvent(ctx context.Context, author disgord.Snowflake) {
	cache.Client.Del(ctx, "/rinha/"+author.String())
}

func lockBattle(ctx context.Context, author, adv disgord.Snowflake, authorUsername, advUsername string) {
	cache.Client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Set(ctx, "/rinha/"+author.String(), advUsername, time.Minute*4)
		p.Set(ctx, "/rinha/"+adv.String(), authorUsername, time.Minute*4)
		return nil
	})
}

func unlockBattle(ctx context.Context, author, adv disgord.Snowflake) {
	cache.Client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Del(ctx, "/rinha/"+author.String())
		p.Del(ctx, "/rinha/"+adv.String())
		return nil
	})
}

func sendLevelUpEmbed(ctx context.Context, itc *disgord.InteractionCreate, galoWinner *entities.Rooster, user *disgord.User, xpOb int) {
	if rinha.CalcLevel(galoWinner.Xp) > rinha.CalcLevel(galoWinner.Xp-xpOb) {
		nextLevel := rinha.CalcLevel(galoWinner.Xp)
		nextSkill := rinha.GetNextSkill(*galoWinner)
		nextSkillStr := ""
		if len(nextSkill) != 0 {
			nextSkillStr = fmt.Sprintf("e desbloqueando a habilidade %s", nextSkill[0].Name)
		}
		ch, _ := handler.Client.Channel(itc.ChannelID).Get()
		if ch != nil {
			ch.SendMsg(ctx, handler.Client, &disgord.Message{
				Embeds: []*disgord.Embed{{
					Title:       "Galo upou de nivel",
					Color:       65535,
					Description: fmt.Sprintf("O galo de %s upou para o nivel %d %s", user.Username, nextLevel, nextSkillStr),
				}},
			})
		}
	}
}

func rinhaMessage(username, text string) *disgord.CreateInteractionResponse {
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: translation.T("RinhaMessage", "pt", map[string]interface{}{"username": username, "text": text}),
		},
	}
}

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "rinha",
		Description: translation.T("RinhaHelp", "pt"),
		Run:         runRinha,
		Cooldown:    8,
		Category:    handler.Rinha,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user",
			Required:    true,
		}),
	})
}

func runRinha(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := utils.GetUser(itc, 0)
	if user.ID == itc.Member.User.ID || user.Bot {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Invalid user",
			},
		}
	}
	author := itc.Member.User
	text := translation.T("RinhaDuel", translation.GetLocale(itc), map[string]interface{}{"author": author.Username, "adv": user.Username})
	utils.Confirm(ctx, text, itc, user.ID, func() {
		userRinha := isInRinha(ctx, user)
		if userRinha != "" {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(user.Username, userRinha).Data.Content,
			})
			return
		}
		authorRinha := isInRinha(ctx, itc.Member.User)
		if authorRinha != "" {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: rinhaMessage(author.Username, authorRinha).Data.Content,
			})
			return
		}
		executePVP(ctx, itc)
	})
	return nil
}

func executePVP(ctx context.Context, itc *disgord.InteractionCreate) {
	author := itc.Member.User
	user := utils.GetUser(itc, 0)
	lockBattle(ctx, author.ID, user.ID, author.Username, user.Username)
	defer unlockBattle(ctx, author.ID, user.ID)
	authorDb := database.User.GetUser(ctx, author.ID, "Galos", "Items", "Trials")
	advDb := database.User.GetUser(ctx, user.ID, "Galos", "Items", "Trials")

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
			database.User.UpdateUser(ctx, loser.ID, func(u entities.User) entities.User {
				u.Lose++
				return u
			})
		}
		database.User.UpdateUser(ctx, winner.ID, func(u entities.User) entities.User {
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
			ch.SendMsg(ctx, handler.Client, &disgord.Message{
				Embeds: []*disgord.Embed{embed},
			})
		}
	}

}
