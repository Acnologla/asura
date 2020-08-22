package commands

import (
	"asura/src/database"
	"asura/src/utils"
	"asura/src/handler"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"math/rand"
	"sync"
	"time"
)

const (
	rinhaUserTitle  = "%s lvl %d"
	rinhaDescAttack = "%s **%s** usou a **Galada** dando **%d** De dano!"
)

type BattleUser struct {
	User *disgord.User
	Galo *database.Galo
	Life int
	Level int
}

var (
	currentBattles = map[disgord.Snowflake]string{}
	battleMutex    *sync.RWMutex = &sync.RWMutex{}
	rinhaEmojis = [2]string{"<:sverde:744682222644363296>","<:svermelha:744682249408217118>"}
	rinhaColors = [2]int{65280,16711680}
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"rinha", "brigadegalo", "briga","rinhateste"},
		Run:       runRinha,
		Available: true,
		Cooldown:  10,
		Usage:     "j!rinha",
		Help:      "Briga",
	})
}


func edit(message *disgord.Message,embed *disgord.Embed){
	utils.Try(func()error{
		msgUpdater := handler.Client.UpdateMessage(context.Background(),message.ChannelID,message.ID)
		msgUpdater.SetEmbed(embed)
		_,err := msgUpdater.Execute()
		return err
	},5)
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
		if currentBattles[msg.Author.ID] != ""{
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Voce ja esta lutando com o " + currentBattles[msg.Author.ID])
			return
		}
		if currentBattles[msg.Mentions[0].ID] != ""{
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Este usuario ja esta lutando com o " + currentBattles[msg.Mentions[0].ID])
			return
		}
		battleMutex.RUnlock()
		battleMutex.Lock()
		currentBattles[msg.Author.ID] = msg.Mentions[0].Username
		currentBattles[msg.Mentions[0].ID] = msg.Author.Username
		battleMutex.Unlock()
		galoAuthor, _ := database.GetGaloDB(msg.Author.ID)
		authorLevel := utils.CalcLevel(galoAuthor.Xp)
		galoAdv, _ := database.GetGaloDB(msg.Mentions[0].ID)
		AdvLevel := utils.CalcLevel(galoAdv.Xp)
		advLife := 100
		authorLife := 100
		if AdvLevel >= 10{
			advLife = 150
		}
		if authorLevel >= 10{
			authorLife = 150
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
					Name: fmt.Sprintf("%s Level %d",msg.Author.Username,authorLevel+1),
					Value: fmt.Sprintf("%d/%d",authorLife,authorLife),
					Inline: true,
				},
				&disgord.EmbedField{
					Name: fmt.Sprintf("%s Level %d",user.Username,AdvLevel+1),
					Value: fmt.Sprintf("%d/%d",advLife,advLife),
					Inline: true,
				},
			},
		}
		message,err := msg.Reply(context.Background(),session,&disgord.CreateMessageParams{
			Content: msg.Author.Mention(),
			Embed: embed,
		})

		if err == nil{
			turn := 0 
			firstPlayer := &BattleUser{
				User: msg.Author,
				Galo: &galoAuthor,
				Life: authorLife,
				Level: authorLevel,
			}
			secondPlayer := &BattleUser{
				User: user,
				Galo: &galoAdv,
				Life: advLife,
				Level: AdvLevel,
			}
			rinhaPlayers := [2]*BattleUser{firstPlayer,secondPlayer}
			var lastAtack string
			for{
				var currentAtack string
				var num int
				var s int
				if rinhaPlayers[turn].Level > 10{
					num := rinhaPlayers[turn].Level - 10
					m := len(skills)-10
					if rinhaPlayers[turn].Level - 10 > m{
						num = m
					}
					s = num
				}
				if rinhaPlayers[turn].Level > len(skills[s:])-1{
					num = rand.Intn(len(skills[s:]))
				}else if rinhaPlayers[turn].Level > 0{
					num = rand.Intn(rinhaPlayers[turn].Level+1)
				}
				atack := skills[s:][num]
				damage := rand.Intn(atack.Damage[1]-1 - atack.Damage[0]) + atack.Damage[0]
				if lastAtack == ""{
					damage = int(damage/2)
				}
				var currentTurn = turn
				if turn == 0{
					turn = 1
				}else{
					turn =0
				}
				rinhaPlayers[turn].Life -= damage
				currentAtack = fmt.Sprintf("%s **%s** Usou **%s** em **%s** causando **%d** de dano!",rinhaEmojis[currentTurn],rinhaPlayers[currentTurn].User.Username,atack.Name,rinhaPlayers[turn].User.Username,damage) 
				embed.Color = rinhaColors[currentTurn]
				embed.Description = lastAtack + currentAtack
				embed.Fields = []*disgord.EmbedField{
					&disgord.EmbedField{
						Name: fmt.Sprintf("%s Level %d",msg.Author.Username,authorLevel+1),
						Value: fmt.Sprintf("%d/%d",rinhaPlayers[0].Life,authorLife),
						Inline: true,
					},
					&disgord.EmbedField{
						Name: fmt.Sprintf("%s Level %d",user.Username,AdvLevel+1),
						Value: fmt.Sprintf("%d/%d",rinhaPlayers[1].Life,advLife),
						Inline: true,
					},
				}
				lastAtack = currentAtack + "\n"
				if 0>= rinhaPlayers[turn].Life{
					xpOb := (rand.Intn(10) + 5) - (2 * (rinhaPlayers[currentTurn].Level - rinhaPlayers[turn].Level))
					if 0 > xpOb{
						xpOb = 0
					}
					rinhaPlayers[currentTurn].Galo.Xp += xpOb
					database.Database.NewRef(fmt.Sprintf("galo/%d",rinhaPlayers[currentTurn].User.ID)).Set(context.Background(), &rinhaPlayers[currentTurn].Galo)
					if (utils.CalcLevel(rinhaPlayers[currentTurn].Galo.Xp) > rinhaPlayers[currentTurn].Level){
						nextLevel := utils.CalcLevel(rinhaPlayers[currentTurn].Galo.Xp) 
						var nextSkill string
						if nextLevel <= len(skills)-1{
							nextSkill = fmt.Sprintf("e desbloqueando a habilidade %s",skills[nextLevel].Name)
						}
						msg.Reply(context.Background(),session,&disgord.CreateMessageParams{
							Embed: &disgord.Embed{
								Title: "Galo upou de nivel",
								Color: 65535,
								Description: fmt.Sprintf("O galo de %s upou para o nivel %d %s",rinhaPlayers[currentTurn].User.Username,nextLevel+1,nextSkill),
							},
						})
					}
					winner := rinhaPlayers[currentTurn]
					embed.Description += fmt.Sprintf("\n\nO galo de **%s** venceu e ganhou %d de XP (%d/%d)",winner.User.Username, xpOb,winner.Galo.Xp,utils.CalcXP(winner.Level+1))
					edit(message,embed)
					battleMutex.Lock()
					delete(currentBattles,msg.Author.ID)
					delete(currentBattles,msg.Mentions[0].ID)
					battleMutex.Unlock()
					break
				}
				edit(message,embed)
				time.Sleep(4 * time.Second)
			}
		}
	} else {
		msg.Reply(context.Background(), session, "Você tem que mencionar alguem!")
	}

}
