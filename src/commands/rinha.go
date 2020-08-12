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

type Battle struct {
	Person *disgord.User
	Other  *disgord.User
	Life   [2]int
	Round  bool
}

var (
	currentBattles []Battle      = []Battle{}
	battleMutex    *sync.RWMutex = &sync.RWMutex{}
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"rinha", "brigadegalo", "briga"},
		Run:       runRinha,
		Available: true,
		Cooldown:  5,
		Usage:     "j!rinha",
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

		galoAuthor, _ := database.GetGaloDB(msg.Author.ID)
		galoAdv, _ := database.GetGaloDB(msg.Mentions[0].ID)

		embed := disgord.Embed{
			Color:       65535,
			Title:       "Briga de galo",
			Fields: []*disgord.EmbedField{
				&disgord.EmbedField{
					Name:   fmt.Sprintf(rinhaUserTitle, msg.Author.Username, utils.CalcLevel(galoAuthor.Xp)+1),
					Value:  "Normal \n**(100/100)**",
					Inline: true,
				},
				&disgord.EmbedField{
					Name:   fmt.Sprintf(rinhaUserTitle, msg.Mentions[0].Username, utils.CalcLevel(galoAdv.Xp)+1),
					Value:  "Normal \n**(100/100)**",
					Inline: true,
				},
			},
			Image: &disgord.EmbedImage{
				URL: "https://i.imgur.com/T1IOzQo.png",
			},
		}

		ctx := context.Background()

		battle := Battle{
			msg.Author,
			msg.Mentions[0],
			[2]int{100, 100},
			false,
		}
		
		battleMutex.Lock()
		currentBattles = append(currentBattles, battle)
		battleMutex.Unlock()

		var mes *disgord.Message

			for {

				if battle.Life[0] <= 0 || battle.Life[1] <= 0 {
					index := 0
					for index = 0; index < len(currentBattles); index++ {
						if currentBattles[index] == battle {
							break
						}
					}
					battleMutex.Lock()	
					if index == len(currentBattles) {
						currentBattles = currentBattles[0:len(currentBattles)-1]
					}else{	
						currentBattles = append(currentBattles[0:index],currentBattles[(index+1):]...)
					}
					battleMutex.Unlock()

					builder := handler.Client.UpdateMessage(ctx, mes.ChannelID, mes.ID)
					builder.SetEmbed(&disgord.Embed{
						Color:       65535,
						Description:       "Alguem venceu e ganhou X levels",
					})
					time.Sleep(4 * time.Second)
					builder.Execute()

					return
				} 


				damage := rand.Intn(20) + 10
			
				attacker := battle.Person
				defended := rand.Intn(20) < 4

				if battle.Round {
					attacker = battle.Other
				}

				if !defended {
					if battle.Round {
						battle.Life[0] -= damage
					} else {
						battle.Life[1] -= damage
					}
				}

				embed.Description = fmt.Sprintf(rinhaDescAttack, ":crossed_swords:", attacker.Username, damage)

				if defended {
					if battle.Round {
						embed.Description = embed.Description + fmt.Sprintf("\n:shield:%s Defendeu!\n", battle.Person.Username)
					} else {
						embed.Description = embed.Description + fmt.Sprintf("\n:shield:%s Defendeu!\n", battle.Other.Username)
					}
				}

				embed.Fields[0].Value = fmt.Sprintf("Normal \n**(%d/100)**", battle.Life[0])
				embed.Fields[1].Value = fmt.Sprintf("Normal \n**(%d/100)**", battle.Life[1])

				battle.Round = !battle.Round
				

				if mes == nil {
					msg, err := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
						Embed: &embed,
					})
					mes = msg
					if err != nil {
						return
					}
				} else {
					builder := handler.Client.UpdateMessage(ctx, mes.ChannelID, mes.ID)
					builder.SetEmbed(&embed)
					builder.Execute()
				}

				time.Sleep(2 * time.Second)
			}

	} else {
		msg.Reply(context.Background(), session, "Você tem que mencionar alguem!")
	}

}
