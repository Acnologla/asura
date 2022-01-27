package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"asura/src/utils/rinha/engine"
	"context"
	"fmt"
	"sync"

	"github.com/andersfylling/disgord"
)

var (
	currentBattles               = map[disgord.Snowflake]string{}
	battleMutex    *sync.RWMutex = &sync.RWMutex{}
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"rinha", "brigadegalo", "rinhadegalo"},
		Run:       runRinha,
		Available: true,
		Cooldown:  10,
		Usage:     "j!rinha <user>",
		Help:      "Briga",
		Category:  1,
	})
}

func initGalo(galo *rinha.Galo, user *disgord.User) {
	if galo.Type == 0 {
		galoType := rinha.GetCommonOrRare()
		galo.Type = galoType
		rinha.SaveGaloDB(user.ID, *galo)
	}
}

func lockBattle(authorID disgord.Snowflake, advID disgord.Snowflake, authorNick string, advNick string) {
	battleMutex.Lock()
	currentBattles[authorID] = advNick
	currentBattles[advID] = authorNick
	battleMutex.Unlock()
}

func unlockBattle(authorID disgord.Snowflake, advID disgord.Snowflake) {
	battleMutex.Lock()
	delete(currentBattles, authorID)
	delete(currentBattles, advID)
	battleMutex.Unlock()
}

func LockEvent(authorID disgord.Snowflake, authorNick string) {
	battleMutex.Lock()
	currentBattles[authorID] = authorNick
	battleMutex.Unlock()
}

func UnlockEvent(authorID disgord.Snowflake) {
	battleMutex.Lock()
	delete(currentBattles, authorID)
	battleMutex.Unlock()
}

func updateGaloWin(id disgord.Snowflake, xp int, win int) {
	rinha.UpdateGaloDB(id, func(galo rinha.Galo) (rinha.Galo, error) {
		galo.Xp += xp
		galo.Win = win
		return galo, nil
	})
}

func runRinha(session disgord.Session, msg *disgord.Message, args []string) {
	if len(msg.Mentions) != 0 {
		if msg.Mentions[0].ID == msg.Author.ID {
			msg.Reply(context.Background(), session, "Você não pode lutar contra si mesmo!")
			return
		}
		if msg.Mentions[0].Bot {
			msg.Reply(context.Background(), session, "Você não pode lutar contra um bot!")
			return
		}
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Você já está lutando com o "+currentBattles[msg.Author.ID])
			return
		}
		if currentBattles[msg.Mentions[0].ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Este usuário já está lutando com o "+currentBattles[msg.Mentions[0].ID])
			return
		}

		battleMutex.RUnlock()
		text := fmt.Sprintf("**%s** você foi convidado para um duelo com %s", msg.Mentions[0].Username, msg.Author.Username)
		utils.Confirm(text, msg.ChannelID, msg.Mentions[0].ID, func() {
			battleMutex.RLock()
			if currentBattles[msg.Author.ID] != "" {
				battleMutex.RUnlock()
				msg.Reply(context.Background(), session, "Você já esta lutando com o "+currentBattles[msg.Author.ID])
				return
			}
			if currentBattles[msg.Mentions[0].ID] != "" {
				battleMutex.RUnlock()
				msg.Reply(context.Background(), session, "Este usuário já esta lutando com o "+currentBattles[msg.Mentions[0].ID])
				return
			}
			battleMutex.RUnlock()
			rinhaTest := false
			if len(args) > 0 {
				rinhaTest = args[len(args)-1] == "rinhatest"
			}
			executePVP(msg, session, rinhaTest)
		})

	} else {
		msg.Reply(context.Background(), session, "Você precisa mencionar alguém")
	}
}

func sendLevelUpEmbed(msg *disgord.Message, session disgord.Session, galoWinner *rinha.Galo, user *disgord.User, xpOb int) {
	if rinha.CalcLevel(galoWinner.Xp) > rinha.CalcLevel(galoWinner.Xp-xpOb) {
		nextLevel := rinha.CalcLevel(galoWinner.Xp)
		nextSkill := rinha.GetNextSkill(*galoWinner)
		nextSkillStr := ""
		if len(nextSkill) != 0 {
			nextSkillStr = fmt.Sprintf("E desbloqueando a habilidade %s", nextSkill[0].Name)
		}
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       "Galo upou de nivel",
				Color:       65535,
				Description: fmt.Sprintf("O galo de %s upou para o nivel %d %s", user.Username, nextLevel, nextSkillStr),
			},
		})
	}
}

func executePVP(msg *disgord.Message, session disgord.Session, newRinhaEngine bool) {
	lockBattle(msg.Author.ID, msg.Mentions[0].ID, msg.Author.Username, msg.Mentions[0].Username)
	defer unlockBattle(msg.Author.ID, msg.Mentions[0].ID)
	author := msg.Author
	adv := msg.Mentions[0]

	galoAdv, _ := rinha.GetGaloDB(adv.ID)
	galoAuthor, _ := rinha.GetGaloDB(author.ID)

	initGalo(&galoAuthor, author)
	initGalo(&galoAdv, adv)

	authorLevel := rinha.CalcLevel(galoAuthor.Xp)
	advLevel := rinha.CalcLevel(galoAdv.Xp)

	authorName := author.Username
	advName := adv.Username

	if galoAuthor.Name != "" {
		authorName = galoAuthor.Name
	}

	if galoAdv.Name != "" {
		advName = galoAdv.Name
	}
	whoWin, battle := engine.ExecuteRinha(msg, session, engine.RinhaOptions{
		GaloAuthor:  galoAuthor,
		GaloAdv:     galoAdv,
		AuthorName:  authorName,
		AdvName:     advName,
		AuthorLevel: authorLevel,
		AdvLevel:    advLevel,
		IDs:         [2]disgord.Snowflake{author.ID, adv.ID},
	}, newRinhaEngine)

	if whoWin == -1 {
		return
	}

	if 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life {
		winnerTurn := whoWin
		turn := 1
		winner := author
		loser := adv

		if whoWin == 1 {
			winner = adv
			loser = author
			turn = 0
		}

		galoWinner := battle.Fighters[winnerTurn].Galo
		galoLoser := battle.Fighters[turn].Galo

		xpOb := 0
		money := 0
		clanMsg := ""
		if 2 >= rinha.CalcLevel(galoWinner.Xp)-rinha.CalcLevel(galoLoser.Xp) {
			galoLoser.Lose++
			rinha.UpdateGaloDB(loser.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				galo.Lose = galoLoser.Lose
				return galo, nil
			})
			galoWinner.Win++
		}
		updateGaloWin(winner.ID, xpOb, galoWinner.Win)
		sendLevelUpEmbed(msg, session, galoWinner, winner, xpOb)
		embed := &disgord.Embed{Title: "Briga de galo", Color: 16776960, Description: ""}

		if winnerTurn == 1 {
			embed.Description = fmt.Sprintf("\n**%s** venceu a batalha, ganhou %d de dinheiro e %d de XP%s", advName, money, xpOb, clanMsg)
		} else {
			embed.Description = fmt.Sprintf("\n**%s** venceu a batalha, ganhou %d de dinheiro e %d de XP%s", authorName, money, xpOb, clanMsg)
		}

		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{Embed: embed})
	}

}
