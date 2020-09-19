package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"github.com/andersfylling/disgord"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"skills"},
		Run:       runSkills,
		Available: true,
		Cooldown:  3,
		Usage:     "j!skills",
		Help:      "Informação sobre seu galo",
	})
}

func runSkills(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author

	galo, _ := utils.GetGaloDB(user.ID)

	
	if len(args) == 0 || (args[0] != "use" && args[0] != "remove") {
	
		desc := ""

		if len(galo.Equipped) != 0 {
			desc += "**Equipped** \n"
	
			for i := 0; i < len(galo.Equipped); i++ {
				skill := utils.Skills[galo.Equipped[i]]
				desc += "**" + strconv.Itoa(i) + "**. [" + strconv.Itoa(skill.Damage[0]) + " - " + strconv.Itoa(skill.Damage[1]) + "] " +  skill.Name + "\n"
			}
		}
	
		desc += "\n**Inventory**\n"

		for i := 0; i < len(galo.Skills); i++ {
			skill := utils.Skills[galo.Skills[i]]
			desc += "**" + strconv.Itoa(i) + "**. [" + strconv.Itoa(skill.Damage[0]) + " - " + strconv.Itoa(skill.Damage[1]) + "] " +  skill.Name + "\n"
		}
		
		msg.Reply(context.Background(), session, disgord.CreateMessageParams{
			Content: msg.Author.Mention(),
			Embed: &disgord.Embed{
				Title: "Skills do seu galo",
				Color: 65535,
				Footer: &disgord.EmbedFooter{
					Text: "Use 'j!skills use [numero_da_skill]' por exemplo 'j!skills use 1' para equipar ou 'j!skills remove [skill]' para retirar uma!",
				},
				Description: desc,
			},
		})
	} else {
		if len(args) == 1 {
			msg.Reply(context.Background(),session, disgord.CreateMessageParams{
				Content: "Voce esta usando errado bob lolo",
			})
			return
		}
		i, err := strconv.Atoi(args[1])
		if args[0] == "use"{
			if len(galo.Equipped) >= 5 {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Voce ja tem 5 ou mais habilidades ativas! use ``j!skills unequip [skill]`` para desativar uma habilidade e assim conseguir ativar outra!",
				})
				return
			}

			if err != nil || i < 0 || i > len(galo.Skills) {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Voce esta usando errado bobo",
				})
				return
			} 
			
			if utils.IdInSkills(galo.Skills[i], galo.Equipped) {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Voce já está com essa habilidade equipada!",
				})
				return
			}
			galo.Equipped = append(galo.Equipped, galo.Skills[i])
			msg.Reply(context.Background(), session, disgord.CreateMessageParams{
				Content: "Voce equipou essa habilidade",
			})
			utils.SaveGaloDB(user.ID, galo)
		} else if args[0] == "remove" {
			if err != nil || i < 0 || i > len(galo.Equipped) {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Voce esta usando errado bobo",
				})
				return
			} 
			msg.Reply(context.Background(), session, disgord.CreateMessageParams{
				Content: "Voce retirou essa habilidade",
			})
			galo.Equipped = append(galo.Equipped[:i], galo.Equipped[i+1:]...)
			utils.SaveGaloDB(user.ID, galo)
		}
	}
}
