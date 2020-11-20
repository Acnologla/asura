package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"sync"
	"time"
	"math/rand"
	"github.com/andersfylling/disgord"
)

type rinhaOptions struct {
	galoAuthor *rinha.Galo
	galoAdv *rinha.Galo
	authorLevel int
	advLevel int
	authorName string
	advName string
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
		return fmt.Sprintf("%s **%s** Usou **%s** causando **%d** de dano\n", rinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name, effect.Damage)
	} else if effect.Effect == rinha.Effected {
		effectLiteral := rinha.GetEffectFromSkill(effect.Skill)
		if effect.Self {
			return fmt.Sprintf(effectLiteral.Phrase+"\n", author, effect.Damage)
		}
		return fmt.Sprintf(effectLiteral.Phrase+"\n", affected, effect.Damage)
	} else if effect.Effect == rinha.SideEffected {
		effectLiteral := rinha.GetEffectFromSkill(effect.Skill)
		if effect.Self {
			return fmt.Sprintf("**%s** Tomou **%d** de dano de '%s'\n", affected, effect.Damage, effectLiteral.Name)
		}
		return fmt.Sprintf("**%s** Tomou **%d** de dano de '%s'\n", author, effect.Damage, effectLiteral.Name)
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

func updateGaloWin(id disgord.Snowflake, galo rinha.Galo){
	updatedGalo,_ := rinha.GetGaloDB(id)
	if updatedGalo.Type == galo.Type && galo.Name != updatedGalo.Name {
		galo.Name = updatedGalo.Name
	}
	rinha.UpdateGaloDB(id, map[string]interface{}{
		"name": galo.Name,
		"xp":   galo.Xp,
		"type": galo.Type,
		"galos": galo.Galos,
		"equipped": galo.Equipped,
		"win": 	galo.Win,
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
			utils.Try(func() error {
				return confirmMsg.React(context.Background(), session, "✅")
			}, 5)
			say := false
			handler.RegisterHandler(confirmMsg, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
				if !removed {
					if emoji.Name == "✅" && u == msg.Mentions[0].ID {
						battleMutex.RLock()
						if currentBattles[msg.Author.ID] != "" {
							battleMutex.RUnlock()
							if !say {
								say = true
								msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
							}
							return
						}
						if currentBattles[msg.Mentions[0].ID] != "" {
							battleMutex.RUnlock()
							if !say {
								say = true
								msg.Reply(context.Background(), session, "Este usuario ja esta lutando com o "+currentBattles[msg.Mentions[0].ID])
							}
							return
						}
						battleMutex.RUnlock()
						go session.Channel(confirmMsg.ChannelID).Message(confirmMsg.ID).Delete()

						executePVP(msg, session)
					}
				}
			}, 120)
		}
	} else {
		msg.Reply(context.Background(), session, "Voce precisa mencionar alguem")
	}
}

func sendLevelUpEmbed(msg *disgord.Message, session disgord.Session, galoWinner *rinha.Galo, user *disgord.User, xpOb int){
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

	
	lockBattle(msg.Author.ID, msg.Mentions[0].ID, msg.Author.Username, msg.Mentions[0].Username)
	winner, battle := ExecuteRinha(msg, session, rinhaOptions{
		galoAuthor: &galoAuthor,
		galoAdv: &galoAdv,
		authorName: authorName,
		advName: advName,
		authorLevel: authorLevel,
		advLevel: advLevel,
	})
	unlockBattle(msg.Author.ID, msg.Mentions[0].ID)
	
	if winner == -1 {
		return
	}


	if 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life {
		turn := battle.GetTurn()
		winner := author
		loser := adv
		winnerTurn := battle.GetReverseTurn()
		
		if 0 >= battle.Fighters[battle.GetReverseTurn()].Life {
			winner = adv
			loser = author
			winnerTurn = turn
			turn = battle.GetReverseTurn()
		}

		galoWinner := battle.Fighters[winnerTurn].Galo
		galoLoser  := battle.Fighters[turn].Galo

		xpOb := (rand.Intn(15) + 15) - (3 * (rinha.CalcLevel(galoWinner.Xp) - rinha.CalcLevel(galoLoser.Xp)))

		if 0 > xpOb {
			xpOb = 0
		}

		money := 0

		if 2 >= rinha.CalcLevel(galoWinner.Xp)-rinha.CalcLevel(galoLoser.Xp) {
			money = 5
			rinha.ChangeMoney(winner.ID,5)
		}

		if xpOb > 9 {
			galoLoser.Lose++
			rinha.UpdateGaloDB(loser.ID, map[string]interface{}{
				"lose": galoLoser.Lose,
			})
			galoWinner.Win++
		}

		galoWinner.Xp += xpOb
		updateGaloWin(winner.ID, *galoWinner)
		sendLevelUpEmbed(msg, session, galoWinner, winner, xpOb)
		embed := &disgord.Embed{ Title: "Briga de galo", Color: 16776960, Description: ""}

		if winnerTurn == 1 {
			embed.Description = fmt.Sprintf("\n**%s** venceu a batalha, ganhou %d de dinheiro e %d de XP", advName, money, xpOb)
		} else {
			embed.Description = fmt.Sprintf("\n**%s** venceu a batalha, ganhou %d de dinheiro e %d de XP", authorName, money, xpOb)
		}
		
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{Embed:   embed})
		unlockBattle(msg.Author.ID, msg.Mentions[0].ID)
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
		battle := rinha.CreateBattle(options.galoAuthor, options.galoAdv)
		var lastEffects string

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
			time.Sleep(4 * time.Second)
		}
	} else {
		return -1, nil
	}
}