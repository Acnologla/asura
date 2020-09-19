package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"math/rand"
	"sync"
	"time"
	"fmt"
	"github.com/andersfylling/disgord"
)

var (
	currentBattles               = map[disgord.Snowflake]string{}
	battleMutex    *sync.RWMutex = &sync.RWMutex{}
	rinhaEmojis                  = [2]string{"<:sverde:744682222644363296>", "<:svermelha:744682249408217118>"}
	rinhaColors                  = [2]int{65280, 16711680}
)

func edit(message *disgord.Message, embed *disgord.Embed) {
	utils.Try(func() error {
		msgUpdater := handler.Client.UpdateMessage(context.Background(), message.ChannelID, message.ID)
		msgUpdater.SetEmbed(embed)
		_, err := msgUpdater.Execute()
		return err
	}, 5)
}

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"rinhanova", "brigadegalo","rinhadegalo"},
		Run:       runRinha,
		Available: true,
		Cooldown:  10,
		Usage:     "j!rinha <user>",
		Help:      "Briga",
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
		battleMutex.Lock()
		currentBattles[msg.Author.ID] = msg.Mentions[0].Username
		currentBattles[msg.Mentions[0].ID] = msg.Author.Username
		battleMutex.Unlock() 
		galoAuthor, _ := utils.GetGaloDB(msg.Author.ID)
		authorLevel := utils.CalcLevel(galoAuthor.Xp)
		galoAdv, _ := utils.GetGaloDB(msg.Mentions[0].ID)
		AdvLevel := utils.CalcLevel(galoAdv.Xp)
		if galoAuthor.Type == 0{
			galoType := rand.Intn(len(utils.Classes))
			galoAuthor.Type = galoType
			utils.SaveGaloDB(msg.Author.ID,galoAuthor)
		}    
		if galoAdv.Type == 0{
			galoType := rand.Intn(len(utils.Classes))
			galoAdv.Type = galoType
		}
		user := msg.Mentions[0]
		embed := &disgord.Embed{
			Title: "Briga de galo",
			Color: 16776960,
			Footer: &disgord.EmbedFooter{
				Text: "Use j!galo para ver informaçoes sobre seu galo",
			},
			Image: &disgord.EmbedImage{
				URL: "https://sports-images.vice.com/images/articles/meta/2015/03/11/on-the-edge-of-the-pit-cockfighting-in-america-1426077876.jpeg",
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

		if err == nil{
			battle := utils.CreateBattle(&galoAuthor,&galoAdv)
			var lastEffects string
			for {	
				effects := battle.Play()
				var text string 
				author := msg.Author
				affected := user
				turn := battle.GetTurn()
				if turn == 0{
					author = user
					affected = msg.Author
				} 
				for _, effect := range effects{
					if effect.Effect == utils.Damaged{
						text += fmt.Sprintf("%s **%s** Usou **%s** em **%s** causando **%d** de dano\n",rinhaEmojis[battle.GetReverseTurn()],author.Username,effect.Skill.Name,affected.Username,effect.Damage)
					}else if effect.Effect == utils.Effected{
						text += fmt.Sprintf("**%s** Causou sangramento em **%s**\n",author.Username,affected.Username)
					}else if effect.Effect == utils.SideEffected{
						if effect.Self{
							text += fmt.Sprintf("**%s** Tomou **%d** de dano de sangramento\n",author.Username,effect.Damage)
						}else{
							text += fmt.Sprintf("**%s** Tomou **%d** de dano de sangramento\n",affected.Username,effect.Damage)
						}
					}else if effect.Effect == utils.NotEffective{
						text += fmt.Sprintf("%s **%s** Usou **%s** em **%s** causando **%d** de dano **reduzido**\n",rinhaEmojis[battle.GetReverseTurn()],author.Username,effect.Skill.Name,affected.Username,effect.Damage)
					}
				}
				embed.Color = rinhaColors[battle.GetReverseTurn()]
				embed.Description = lastEffects + text
				embed.Fields = []*disgord.EmbedField{
					&disgord.EmbedField{
						Name:   fmt.Sprintf("%s Level %d", msg.Author.Username, authorLevel),
						Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, 100),
						Inline: true,
					},
					&disgord.EmbedField{
						Name:   fmt.Sprintf("%s Level %d", user.Username, AdvLevel),
						Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, 100),
						Inline: true,
					},
				}
				if 0 >= battle.Fighters[0].Life  || 0>= battle.Fighters[1].Life{
					winner := author
					winnerTurn := battle.GetTurn()
					if 0 >= battle.Fighters[turn].Life {
						winner = affected
						winnerTurn = turn
						turn = battle.GetReverseTurn()
					} 
					xpOb := (rand.Intn(10) + 5) - (2 * (utils.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp) - utils.CalcLevel(battle.Fighters[turn].Galo.Xp)))
							if 0 > xpOb {
								xpOb = 0
							}
							battle.Fighters[winnerTurn].Galo.Xp += xpOb
							utils.SaveGaloDB(winner.ID, *battle.Fighters[winnerTurn].Galo)
							if utils.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp) > utils.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp - xpOb) {
								nextLevel := utils.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp) 
								nextSkill := utils.GetNextSkill(*battle.Fighters[winnerTurn].Galo)
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
							embed.Description += fmt.Sprintf("\nO galo de **%s** venceu e ganhou %d de XP (%d/%d)", winner.Username, xpOb, battle.Fighters[winnerTurn].Galo.Xp, utils.CalcXP(utils.CalcLevel(battle.Fighters[winnerTurn].Galo.Xp) +1 ))
							edit(message, embed)
							battleMutex.Lock()
							delete(currentBattles, msg.Author.ID)
							delete(currentBattles, msg.Mentions[0].ID)
							battleMutex.Unlock()
							break
				} 				
				edit(message,embed)
				lastEffects = text
				time.Sleep(4 * time.Second)
			}
			
		}else{
			battleMutex.Lock()
			delete(currentBattles, msg.Author.ID)
			delete(currentBattles, msg.Mentions[0].ID)
			battleMutex.Unlock()
		}
	}else{
		msg.Reply(context.Background(),session,"aVoce precisa mencionar alguem")
	}
}
