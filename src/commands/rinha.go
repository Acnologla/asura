package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"sync"
	"time"
)

type rinhaOptions struct {
	galoAuthor  *rinha.Galo
	galoAdv     *rinha.Galo
	authorLevel int
	advLevel    int
	authorName  string
	advName     string
	noItems     bool
}

var (
	currentBattles               = map[disgord.Snowflake]string{}
	battleMutex    *sync.RWMutex = &sync.RWMutex{}
	rinhaEmojis                  = [2]string{"<:sverde:744682222644363296>", "<:svermelha:744682249408217118>"}
	rinhaColors                  = [2]int{65280, 16711680}
)

func edit(message *disgord.Message, embed *disgord.Embed) {
	utils.Try(func() error {
		msgUpdater := handler.Client.Channel(message.ChannelID).Message(message.ID).Update()
		msgUpdater.SetEmbed(embed)
		_, err := msgUpdater.Execute()
		return err
	}, 5)
}

func getImageTile(first *rinha.Galo, sec *rinha.Galo, turn int) string {
	if turn == 0 {
		return rinha.Sprites[turn^1][sec.Type-1]
	}
	return rinha.Sprites[turn^1][first.Type-1]
}

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

func effectToStr(effect *rinha.Result, affected string, author string, battle *rinha.Battle) string {
	if effect.Effect == rinha.Damaged {
		if effect.Skill.Self {
			return fmt.Sprintf("%s **%s** Usou **%s** em si mesmo\n", rinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name)
		}
		if effect.Reflected {
			return fmt.Sprintf("%s **%s** Refletiu o ataque **%s** causando **%d** de dano\n", rinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name, effect.Damage)
		}
		return fmt.Sprintf("%s **%s** Usou **%s** causando **%d** de dano\n", rinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name, effect.Damage)
	} else if effect.Effect == rinha.Effected {
		effectLiteral := rinha.Effects[effect.EffectID]
		if effect.Self {
			return fmt.Sprintf(effectLiteral.Phrase+"\n", author, effect.Damage)
		}
		return fmt.Sprintf(effectLiteral.Phrase+"\n", affected, effect.Damage)
	} else if effect.Effect == rinha.NotEffective {
		return fmt.Sprintf("**reduzido**\n")
	}
	return ""
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
		confirmMsg, confirmErr := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Content: msg.Mentions[0].Mention(),
			Embed: &disgord.Embed{
				Color:       65535,
				Description: fmt.Sprintf("**%s** clique na reação abaixo para aceitar o duelo", msg.Mentions[0].Username),
			},
		})
		if confirmErr == nil {
			utils.Confirm(confirmMsg, msg.Mentions[0].ID, func() {
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
				executePVP(msg, session)
			})
		}
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

func executePVP(msg *disgord.Message, session disgord.Session) {
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

	whoWin, battle := ExecuteRinha(msg, session, rinhaOptions{
		galoAuthor:  &galoAuthor,
		galoAdv:     &galoAdv,
		authorName:  authorName,
		advName:     advName,
		authorLevel: authorLevel,
		advLevel:    advLevel,
	})

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

		xpOb := (utils.RandInt(10) + 4) - (3 * (rinha.CalcLevel(galoWinner.Xp) - rinha.CalcLevel(galoLoser.Xp)))

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
				if level >= 4 {
					money++
				}
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

func ExecuteRinha(msg *disgord.Message, session disgord.Session, options rinhaOptions) (int, *rinha.Battle) {

	embed := &disgord.Embed{
		Title: "Briga de galo",
		Color: 16776960,
		Footer: &disgord.EmbedFooter{
			Text: "Use j!galo para ver informaçoes sobre seu galo",
		},
		Image: &disgord.EmbedImage{
			URL: getImageTile(options.galoAuthor, options.galoAdv, 0),
		},
		Description: "Iniciando a briga de galo	",
		Fields: []*disgord.EmbedField{
			&disgord.EmbedField{
				Name:   fmt.Sprintf("%s Level %d", options.authorName, options.authorLevel),
				Value:  fmt.Sprintf("%d/%d", 100, 100),
				Inline: true,
			},
			&disgord.EmbedField{
				Name:   fmt.Sprintf("%s Level %d", options.advName, options.advLevel),
				Value:  fmt.Sprintf("%d/%d", 100, 100),
				Inline: true,
			},
		},
	}

	message, err := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed:   embed,
	})

	if err == nil {
		battle := rinha.CreateBattle(options.galoAuthor, options.galoAdv, options.noItems)
		var lastEffects string
		round := 0
		for {
			effects := battle.Play()
			var text string

			authorName := options.authorName
			affectedName := options.advName

			turn := battle.GetTurn()

			if turn == 0 {
				authorName = options.advName
				affectedName = options.authorName
			}

			for _, effect := range effects {
				text += effectToStr(effect, affectedName, authorName, &battle)
			}
			if round >= 35 {
				if battle.Fighters[1].Life >= battle.Fighters[0].Life {
					text += "\n" + options.authorName + " Foi executado"
					battle.Fighters[0].Life = 0
					battle.Turn = false
				} else {
					battle.Fighters[1].Life = 0
					text += "\n" + options.advName + " Foi executado"
					battle.Turn = true
				}
			}
			embed.Color = rinhaColors[battle.GetReverseTurn()]
			embed.Description = lastEffects + "\n" + text

			embed.Fields = []*disgord.EmbedField{
				&disgord.EmbedField{
					Name:   fmt.Sprintf("%s Level %d", options.authorName, options.authorLevel),
					Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, battle.Fighters[0].MaxLife),
					Inline: true,
				},
				&disgord.EmbedField{
					Name:   fmt.Sprintf("%s Level %d", options.advName, options.advLevel),
					Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, battle.Fighters[1].MaxLife),
					Inline: true,
				},
			}

			embed.Image = &disgord.EmbedImage{
				URL: getImageTile(options.galoAuthor, options.galoAdv, turn),
			}

			if 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life {
				winnerTurn := battle.GetReverseTurn()

				if winnerTurn == 1 {
					embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.advName)
				} else {
					embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.authorName)
				}
				edit(message, embed)
				return winnerTurn, &battle
			}

			edit(message, embed)
			lastEffects = text
			round++
			time.Sleep(5 * time.Second)
		}
	} else {
		return -1, nil
	}
}
