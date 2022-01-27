package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"strconv"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"skills"},
		Run:       runSkills,
		Available: true,
		Cooldown:  3,
		Usage:     "j!skills",
		Help:      "Informação sobre seu galo",
		Category:  2,
	})
}

func runSkills(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author

	galo, _ := rinha.GetGaloDB(user.ID)

	skills := rinha.GetSkills(galo)
	if len(skills) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Você ainda não batalhou nenhuma vez")
		return
	}
	if len(args) == 0 || (args[0] != "use" && args[0] != "remove") {

		desc := ""

		if len(galo.Equipped) != 0 {
			desc += "**Equipped** \n"

			for i := 0; i < len(galo.Equipped); i++ {
				skill := rinha.Skills[galo.Type-1][galo.Equipped[i]]
				desc += "**" + strconv.Itoa(i) + ". " + skill.Name + "**  " + rinha.SkillToStringFormated(skill, galo) + "\n"
			}
		}

		desc += "\n**Inventory**\n"

		for i := 0; i < len(skills); i++ {
			skill := rinha.Skills[galo.Type-1][skills[i]]
			desc += "**" + strconv.Itoa(i) + ". " + skill.Name + "**  " + rinha.SkillToStringFormated(skill, galo) + "\n"
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
			msg.Reply(context.Background(), session, disgord.CreateMessageParams{
				Content: "Use 'j!skills use [numero_da_skill]' por exemplo 'j!skills use 1' para equipar ou 'j!skills remove [skill]' para retirar uma!",
			})
			return
		}
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere sua rinha terminar antes de equipar ou remover habilidades")
			return
		}
		battleMutex.RUnlock()
		i, err := strconv.Atoi(args[1])
		if args[0] == "use" {
			if len(galo.Equipped) >= 5 {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Voce ja tem 5 ou mais habilidades ativas! use ``j!skills remove [skill]`` para desativar uma habilidade e assim conseguir ativar outra!",
				})
				return
			}

			if err != nil || i < 0 || i >= len(skills) {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Você não tem essa habilidade",
				})
				return
			}
			if rinha.IsIntInList(skills[i], galo.Equipped) {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Você já está com essa habilidade equipada!",
				})
				return
			}
			rinha.UpdateGaloDB(user.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				galo.Equipped = append(galo.Equipped, skills[i])
				return galo, nil
			})
			msg.Reply(context.Background(), session, disgord.CreateMessageParams{
				Content: "Você equipou essa habilidade",
			})
		} else if args[0] == "remove" {
			if err != nil || i < 0 || i >= len(galo.Equipped) {
				msg.Reply(context.Background(), session, disgord.CreateMessageParams{
					Content: "Você nao esta equipado com essa habilidade",
				})
				return
			}
			rinha.UpdateGaloDB(user.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				if len(galo.Equipped)-1 == i {
					galo.Equipped = galo.Equipped[:i]
				} else {
					galo.Equipped = append(galo.Equipped[:i], galo.Equipped[i+1:]...)
				}
				return galo, nil
			})
			msg.Reply(context.Background(), session, disgord.CreateMessageParams{
				Content: "Você retirou essa habilidade",
			})
		}
	}
}
