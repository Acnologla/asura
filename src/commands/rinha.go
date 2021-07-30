package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"sync"
	"time"

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
		galoType := rinha.GetRandByType(rinha.Common)
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
			msg.Reply(context.Background(), session, "Voce não pode lutar contra si mesmo!")
			return
		}
		if msg.Mentions[0].Bot {
			msg.Reply(context.Background(), session, "Voce não pode lutar contra um bot!")
			return
		}
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
			return
		}
		if currentBattles[msg.Mentions[0].ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Este usuario ja esta lutando com o "+currentBattles[msg.Mentions[0].ID])
			return
		}

		battleMutex.RUnlock()
		text := fmt.Sprintf("**%s** voce foi convidado para um duelo com %s", msg.Mentions[0].Username, msg.Author.Username)
		utils.Confirm(text, msg.ChannelID, msg.Mentions[0].ID, func() {
			battleMutex.RLock()
			if currentBattles[msg.Author.ID] != "" {
				battleMutex.RUnlock()
				msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
				return
			}
			if currentBattles[msg.Mentions[0].ID] != "" {
				battleMutex.RUnlock()
				msg.Reply(context.Background(), session, "Este usuario ja esta lutando com o "+currentBattles[msg.Mentions[0].ID])
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
		msg.Reply(context.Background(), session, "Voce precisa mencionar alguem")
	}
}

func sendLevelUpEmbed(msg *disgord.Message, session disgord.Session, galoWinner *rinha.Galo, user *disgord.User, xpOb int) {
	if rinha.CalcLevel(galoWinner.Xp) > rinha.CalcLevel(galoWinner.Xp-xpOb) {
		nextLevel := rinha.CalcLevel(galoWinner.Xp)
		nextSkill := rinha.GetNextSkill(*galoWinner)
		nextSkillStr := ""
		if len(nextSkill) != 0 {
			nextSkillStr = fmt.Sprintf("e desbloqueando a habilidade %s", nextSkill[0].Name)
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
	whoWin, battle := ExecuteRinha(msg, session, rinha.RinhaOptions{
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

		xpOb := (utils.RandInt(6) + 3) - (3 * (rinha.CalcLevel(galoWinner.Xp) - rinha.CalcLevel(galoLoser.Xp)))

		if 0 > xpOb {
			xpOb = 0
		}

		money := 0
		clanMsg := ""
		if galoWinner.Clan != "" {
			clan := rinha.GetClan(galoWinner.Clan)
			xpOb += int(xpOb / 10)
			level := rinha.ClanXpToLevel(clan.Xp)
			if level >= 2 {
				xpOb += int(xpOb / 10)
			}
			if 2 >= rinha.CalcLevel(galoWinner.Xp)-rinha.CalcLevel(galoLoser.Xp) {
				if level >= 5 {
					money++
				}
				go rinha.CompleteClanMission(galoWinner.Clan, winner.ID)
				clanMsg = "\nGanhou **1** de xp para seu clan"
			}
		}
		if 2 >= rinha.CalcLevel(galoWinner.Xp)-rinha.CalcLevel(galoLoser.Xp) {
			money += 2
			rinha.ChangeMoney(winner.ID, money, 0)
			galoLoser.Lose++
			rinha.UpdateGaloDB(loser.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				galo.Lose = galoLoser.Lose
				return galo, nil
			})
			galoWinner.Win++
		}

		vip := rinha.IsVip(*galoWinner)
		if xpOb > 38 && !vip {
			xpOb = 38
		}
		if vip {
			xpOb += int(xpOb / 4)
		}
		galoWinner.Xp += xpOb
		updateGaloWin(winner.ID, xpOb, galoWinner.Win)
		rinha.CompleteMission(winner.ID, *galoWinner, *galoLoser, true, msg)
		rinha.CompleteMission(loser.ID, *galoLoser, *galoWinner, false, msg)
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

func RinhaEngine(battle *rinha.Battle, options *rinha.RinhaOptions, message *disgord.Message, embed *disgord.Embed) (int, *rinha.Battle) {
	var lastEffects string
	round := 0
	for {
		effects := battle.Play(-1)
		var text string

		authorName := options.AuthorName
		affectedName := options.AdvName

		turn := battle.GetTurn()

		if turn == 0 {
			authorName = options.AdvName
			affectedName = options.AuthorName
		}

		for _, effect := range effects {
			text += rinha.EffectToStr(effect, affectedName, authorName, battle)
		}
		if round >= 35 {
			if battle.Fighters[1].Life >= battle.Fighters[0].Life {
				text += "\n" + options.AuthorName + " Foi executado"
				battle.Fighters[0].Life = 0
				battle.Turn = false
			} else {
				battle.Fighters[1].Life = 0
				text += "\n" + options.AdvName + " Foi executado"
				battle.Turn = true
			}
		}
		embed.Color = rinha.RinhaColors[battle.GetReverseTurn()]
		embed.Description = lastEffects + "\n" + text

		embed.Fields = []*disgord.EmbedField{
			{
				Name:   fmt.Sprintf("%s Level %d", options.AuthorName, options.AuthorLevel),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, battle.Fighters[0].MaxLife),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s Level %d", options.AdvName, options.AdvLevel),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, battle.Fighters[1].MaxLife),
				Inline: true,
			},
		}

		embed.Image = &disgord.EmbedImage{
			URL: rinha.GetImageTile(&options.GaloAuthor, &options.GaloAdv, turn),
		}

		if 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life {
			winnerTurn := battle.GetReverseTurn()

			if winnerTurn == 1 {
				embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.AdvName)
			} else {
				embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.AuthorName)
			}
			rinha.EditRinhaEmbed(message, embed)
			return winnerTurn, battle
		}

		rinha.EditRinhaEmbed(message, embed)
		lastEffects = text
		round++
		time.Sleep(4 * time.Second)
	}
}

func ExecuteRinha(msg *disgord.Message, session disgord.Session, options rinha.RinhaOptions, newEngine bool) (int, *rinha.Battle) {
	if rinha.HasUpgrade(options.GaloAuthor.Upgrades, 2, 1, 0) {
		options.GaloAdv.Xp = rinha.CalcXP(rinha.CalcLevel(options.GaloAdv.Xp)) - 1
		if rinha.HasUpgrade(options.GaloAuthor.Upgrades, 2, 1, 0, 0) {
			options.GaloAdv.Xp = rinha.CalcXP(rinha.CalcLevel(options.GaloAdv.Xp)) - 1
		}
	}
	if rinha.HasUpgrade(options.GaloAdv.Upgrades, 2, 1, 0) {
		options.GaloAuthor.Xp = rinha.CalcXP(rinha.CalcLevel(options.GaloAuthor.Xp)) - 1
		if rinha.HasUpgrade(options.GaloAdv.Upgrades, 2, 1, 0, 0) {
			options.GaloAuthor.Xp = rinha.CalcXP(rinha.CalcLevel(options.GaloAuthor.Xp)) - 1
		}
	}
	if options.GaloAuthor.Xp < 0 {
		options.GaloAuthor.Xp = 0
	}
	if options.GaloAdv.Xp < 0 {
		options.GaloAdv.Xp = 0
	}
	options.AdvLevel = rinha.CalcLevel(options.GaloAdv.Xp)
	options.AuthorLevel = rinha.CalcLevel(options.GaloAuthor.Xp)
	u, _ := handler.Client.User(options.IDs[0]).Get()
	avatar, _ := u.AvatarURL(128, true)
	embed := &disgord.Embed{
		Title: "Briga de galo",
		Color: rinha.RinhaColors[0],
		Footer: &disgord.EmbedFooter{
			Text:    u.Username,
			IconURL: avatar,
		},
		Image: &disgord.EmbedImage{
			URL: rinha.GetImageTile(&options.GaloAuthor, &options.GaloAdv, 1),
		},
	}
	if newEngine {
		if 3 > options.AdvLevel || 3 > options.AuthorLevel {
			msg.Reply(context.Background(), session, "Tem que ser pelomenos nivel 3 para batalhar na rinha com botoes")
			return -1, nil
		}
	}
	message, err := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed:   embed,
	})

	if err == nil {
		battle := rinha.CreateBattle(options.GaloAuthor, options.GaloAdv, options.NoItems, options.IDs[0], options.IDs[1])
		embed.Fields = []*disgord.EmbedField{
			{
				Name:   fmt.Sprintf("%s Level %d", options.AuthorName, options.AuthorLevel),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, battle.Fighters[0].MaxLife),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s Level %d", options.AdvName, options.AdvLevel),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, battle.Fighters[1].MaxLife),
				Inline: true,
			},
		}
		rinha.EditRinhaEmbed(message, embed)
		if newEngine {
			return rinha.RinhaEngine(&battle, &options, message, embed)
		}
		return RinhaEngine(&battle, &options, message, embed)
	} else {
		return -1, nil
	}
}
