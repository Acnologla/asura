package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"upgrades", "melhorias", "up"},
		Run:       runUpgrades,
		Available: true,
		Cooldown:  5,
		Usage:     "j!upgrades",
		Help:      "Informação sobre seus upgrades",
		Category:  2,
	})
}

func runUpgrades(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author
	galo, _ := rinha.GetGaloDB(user.ID)
	if len(args) == 0 {
		desc := fmt.Sprintf("Upgrade Xp: **%d/%d**", galo.UserXp, rinha.CalcUserXp(galo))
		if len(galo.Upgrades) > 0 {
			desc += "\n\nUpgrades:\n"
			upgrades := rinha.Upgrade{
				Childs: rinha.Upgrades,
			}
			for i, upgrade := range galo.Upgrades {
				upgrades = upgrades.Childs[upgrade]
				v := strings.Repeat("-", i*5)
				desc += fmt.Sprintf("%s%s\n%s%s\n", v, upgrades.Name, v, upgrades.Value)
			}
		}
		if rinha.HavePoint(galo) {
			desc += "\nVoce tem um ponto de upgrade disponivel.\nUse **j!upgrade <id do upgrade>** para dar updgrade.\n\nUpgrades:\n\n" + rinha.UpgradesToString(galo)
		}
		msg.Reply(context.Background(), session, disgord.CreateMessageParams{
			Content: msg.Author.Mention(),
			Embed: &disgord.Embed{
				Title:       "Upgrades",
				Color:       65535,
				Description: desc,
			},
		})
	} else {
		i, err := strconv.Atoi(args[0])
		if !rinha.HavePoint(galo) {
			msg.Reply(context.Background(), session, "Você não tem pontos de upgrades, use j!upgrades para ver seus pontos de upgrades")
			return
		}
		if err != nil {
			msg.Reply(context.Background(), session, "Id do upgrade invalido use j!upgrade para ver os upgrades disponiveis")
			return
		}
		upgrades := rinha.GetCurrentUpgrade(galo)
		if 0 > i || i >= len(upgrades.Childs) {
			msg.Reply(context.Background(), session, "Id do upgrade inválido use j!upgrade para ver os upgrades disponiveis")
			return
		}
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			galo.Upgrades = append(galo.Upgrades, i)
			return galo, nil
		})
		msg.Reply(context.Background(), session, fmt.Sprintf("Upgrade %s adquirido com sucesso", upgrades.Childs[i].Name))
	}
}
