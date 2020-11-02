package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"math/rand"
	"sync"
	"time"
)

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
	} else {
		return rinha.Sprites[turn^1][first.Type-1]
	}
}

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"rinha", "brigadegalo", "rinhadegalo"},
		Run:       runRinha,
		Available: true,
		Cooldown:  10,
		Usage:     "j!rinha <user>",
		Help:      "Briga",
	})
}

func effectToStr(effect *rinha.Result, affected *disgord.User, author *disgord.User, battle *rinha.Battle) string {
	if effect.Effect == rinha.Damaged {
		if effect.Skill.Self {
			return fmt.Sprintf("%s **%s** Usou **%s** em si mesmo\n", rinhaEmojis[battle.GetReverseTurn()], author.Username, effect.Skill.Name)
		} else {
			return fmt.Sprintf("%s **%s** Usou **%s** causando **%d** de dano\n", rinhaEmojis[battle.GetReverseTurn()], author.Username, effect.Skill.Name, effect.Damage)
		}
	} else if effect.Effect == rinha.Effected {
		effect_literal := rinha.GetEffectFromSkill(effect.Skill)
		if effect.Self {
			return fmt.Sprintf(effect_literal.Phrase+"\n", author.Username, effect.Damage)
		} else {
			return fmt.Sprintf(effect_literal.Phrase+"\n", affected.Username, effect.Damage)
		}
	} else if effect.Effect == rinha.SideEffected {
		effect_literal := rinha.GetEffectFromSkill(effect.Skill)
		if effect.Self {
			return fmt.Sprintf("**%s** Tomou **%d** de dano de '%s'\n", affected.Username, effect.Damage, effect_literal.Name)
		} else {
			return fmt.Sprintf("**%s** Tomou **%d** de dano de '%s'\n", author.Username, effect.Damage, effect_literal.Name)
		}
	} else if effect.Effect == rinha.NotEffective {
		return fmt.Sprintf("**reduzido**\n")
	}
	return ""
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
		utils.Try(func () error {
			return confirmMsg.React(context.Background(),session,"✅")
		},5)
		if confirmErr == nil {
			say := false
			handler.RegisterHandler(confirmMsg, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
				if !removed {
					if emoji.Name == "✅" && u == msg.Mentions[0].ID {
						battleMutex.RLock()
						if currentBattles[msg.Author.ID] != "" {
							battleMutex.RUnlock()
							if !say{
								say = true
								msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
							}
							return
						}
						if currentBattles[msg.Mentions[0].ID] != "" {
							battleMutex.RUnlock()
							if !say{
								say = true
								msg.Reply(context.Background(), session, "Este usuario ja esta lutando com o "+currentBattles[msg.Mentions[0].ID])
							}
							return
						}
						battleMutex.RUnlock()
						go session.Channel(confirmMsg.ChannelID).Message(confirmMsg.ID).Delete()
						executeRinha(msg,session)
					}
				}
			}, 120)
		}

	} else {
		msg.Reply(context.Background(), session, "Voce precisa mencionar alguem")
	}
}

func executeRinha(msg *disgord.Message,session disgord.Session) {
	galoAdv, _ := rinha.GetGaloDB(msg.Mentions[0].ID)
	battleMutex.Lock()
	currentBattles[msg.Author.ID] = msg.Mentions[0].Username
	currentBattles[msg.Mentions[0].ID] = msg.Author.Username
	battleMutex.Unlock()
	user := msg.Mentions[0]

	galoAuthor, _ := rinha.GetGaloDB(msg.Author.ID)
	authorLevel := rinha.CalcLevel(galoAuthor.Xp)

	AdvLevel := rinha.CalcLevel(galoAdv.Xp)

	if galoAuthor.Type == 0 {
		galoType := rand.Intn(len(rinha.Classes)-1) + 1
		galoAuthor.Type = galoType
		rinha.SaveGaloDB(msg.Author.ID, galoAuthor)
	}
	if galoAdv.Type == 0 {
		galoType := rand.Intn(len(rinha.Classes)-1) + 1
		galoAdv.Type = galoType
		rinha.SaveGaloDB(user.ID, galoAdv)
	}

	embed := &disgord.Embed{
		Title: "Briga de galo",
		Color: 16776960,
		Footer: &disgord.EmbedFooter{
			Text: "Use j!galo para ver informaçoes sobre seu galo",
		},
		Image: &disgord.EmbedImage{
			URL: getImageTile(&galoAuthor, &galoAdv, 0),
		},
		Description: "Iniciando a briga de galo	",
		Fields: []*disgord.EmbedField{
			&disgord.EmbedField{
				Name:   fmt.Sprintf("%s Level %d", msg.Author.Username, authorLevel),
				Value:  fmt.Sprintf("%d/%d", 100, 100),
				Inline: true,
			},
			&disgord.EmbedField{
				Name:   fmt.Sprintf("%s Level %d", user.Username, AdvLevel),
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
		battle := rinha.CreateBattle(&galoAuthor, &galoAdv)
		var lastEffects string
		for {
			effects := battle.Play()
			var text string
			author := msg.Author
			affected := user
			turn := battle.GetTurn()
			if turn == 0 {
				author = user
				affected = msg.Author
			}

			for _, effect := range effects {
				text += effectToStr(effect, affected, author, &battle)
			}

			embed.Color = rinhaColors[battle.GetReverseTurn()]
			embed.Description = lastEffects + "\n" + text
			embed.Fields = []*disgord.EmbedField{
				&disgord.EmbedField{
					Name:   fmt.Sprintf("%s Level %d", msg.Author.Username, authorLevel),
					Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, battle.Fighters[0].MaxLife),
					Inline: true,
				},
				&disgord.EmbedField{
					Name:   fmt.Sprintf("%s Level %d", user.Username, AdvLevel),
					Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, battle.Fighters[1].MaxLife),
					Inline: true,
				},
			}
			embed.Image = &disgord.EmbedImage{
				URL: getImageTile(&galoAuthor, &galoAdv, turn),
			}
			if 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life {
				winner := author
				loser := affected
				winnerTurn := battle.GetReverseTurn()
				if 0 >= battle.Fighters[battle.GetReverseTurn()].Life {
					winner = affected
					loser = author
					winnerTurn = turn
					turn = battle.GetReverseTurn()
				}
				xpOb := (rand.Intn(10) + 5) - (2 * (rinha.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp) - rinha.CalcLevel(battle.Fighters[turn].Galo.Xp)))
				if 0 > xpOb {
					xpOb = 0
				}
				if xpOb > 3 {
					battle.Fighters[turn].Galo.Lose++
					rinha.SaveGaloDB(loser.ID, *battle.Fighters[turn].Galo)
					battle.Fighters[winnerTurn].Galo.Win++
				}
				battle.Fighters[winnerTurn].Galo.Xp += xpOb
				rinha.SaveGaloDB(winner.ID, *battle.Fighters[winnerTurn].Galo)
				if rinha.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp) > rinha.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp-xpOb) {
					nextLevel := rinha.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp)
					nextSkill := rinha.GetNextSkill(*battle.Fighters[winnerTurn].Galo)
					nextSkillStr := ""
					if len(nextSkill) != 0 {
						nextSkillStr = fmt.Sprintf("e desbloqueando a habilidade %s", nextSkill[0].Name)
					}
					msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "Galo upou de nivel",
							Color:       65535,
							Description: fmt.Sprintf("O galo de %s upou para o nivel %d %s", winner.Username, nextLevel, nextSkillStr),
						},
					})
				}
				embed.Description += fmt.Sprintf("\nO galo de **%s** venceu e ganhou %d de XP (%d/%d)", winner.Username, xpOb, battle.Fighters[winnerTurn].Galo.Xp, rinha.CalcXP(rinha.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp)+1))
				edit(message, embed)
				battleMutex.Lock()
				delete(currentBattles, msg.Author.ID)
				delete(currentBattles, msg.Mentions[0].ID)
				battleMutex.Unlock()
				break
			}
			edit(message, embed)
			lastEffects = text
			time.Sleep(4 * time.Second)
		}
	} else {
		battleMutex.Lock()
		delete(currentBattles, msg.Author.ID)
		delete(currentBattles, msg.Mentions[0].ID)
		battleMutex.Unlock()
	}
}
