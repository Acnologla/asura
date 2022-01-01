package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"backgrounds", "bg"},
		Run:       runBackground,
		Available: true,
		Cooldown:  4,
		Usage:     "j!backgrounds",
		Help:      "Troque seu background",
		Category:  2,
	})
}

func runBackground(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if len(args) == 0 {
		text := ""
		bgs, indexes := rinha.GetBackgrounds(galo.Cosmetics)
		if len(bgs) == 0 {
			text = "Background equipado: **padrão**\n\n"
		} else {
			equippedBg := rinha.Cosmetics[galo.Cosmetics[galo.Background]]
			text = fmt.Sprintf("Background equipado: **%s** (%s)\n\n", equippedBg.Name, equippedBg.Rarity.String())
		}
		for i, background := range bgs {
			text += fmt.Sprintf("[%d] - %s (Raridade: **%s**) \n", indexes[i], background.Name, background.Rarity.String())
		}
		if len(bgs) == 0 {
			text = "Você não tem nenhum background, compre uma caixa cosmetica usando j!lootbox, para obter um"
		}
		avatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       65535,
			Title:       "Backgrounds",
			Description: text,
			Footer: &disgord.EmbedFooter{
				IconURL: avatar,
				Text:    "Use j!background <numero do background> para equipar um background",
			},
		})
	} else {
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere sua rinha terminar para equipar backgrounds")
			return
		}
		battleMutex.RUnlock()
		if len(args) >= 2 {
			if args[0] == "vip" {
				if rinha.IsVip(galo) {
					if utils.CheckImage(args[1]) {
						rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
							galo.VipBackground = args[1]
							return galo, nil
						})
						msg.Reply(context.Background(), session, "Background VIP alterado com sucesso")
						return
					} else {
						msg.Reply(context.Background(), session, "Imagem invalida")
						return
					}
				}
			}
			if args[1] == "remove" || args[1] == "vender" {
				value, err := strconv.Atoi(args[0])
				if err != nil {
					msg.Reply(context.Background(), session, "Use j!background <número do background> remove para vender um background")
					return
				}
				if value < 0 || len(galo.Cosmetics) <= value {
					msg.Reply(context.Background(), session, "Background inválido")
					return
				}
				sellBg := galo.Cosmetics[value]
				cosmetic := rinha.Cosmetics[sellBg]
				priceToSell := rinha.SellCosmetic(*cosmetic)
				utils.Confirm(fmt.Sprintf("Voce deseja vender o Background **%s** por **%d** de dinheiro", cosmetic.Name, priceToSell), msg.ChannelID, msg.Author.ID, func() {
					var price int
					rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
						sellBg = galo.Cosmetics[value]
						cosmetic = rinha.Cosmetics[sellBg]
						for i := value; i < len(galo.Cosmetics)-1; i++ {
							galo.Cosmetics[i] = galo.Cosmetics[i+1]
						}
						price = rinha.SellCosmetic(*cosmetic)
						galo.Money += price
						galo.Background = 0
						galo.Cosmetics = galo.Cosmetics[0 : len(galo.Cosmetics)-1]
						return galo, nil
					})
					tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
					telemetry.Debug(fmt.Sprintf("%s Sell %s", tag, cosmetic.Name), map[string]string{
						"cosmetic": cosmetic.Name,
						"user":     strconv.FormatUint(uint64(msg.Author.ID), 10),
						"rarity":   cosmetic.Rarity.String(),
					})
					msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce vendeu o background **%s** por **%d** de dinheiro com sucesso", msg.Author.Mention(), cosmetic.Name, price))
				})
				return
			}
		}
		value, err := strconv.Atoi(args[0])
		if err != nil {
			msg.Reply(context.Background(), session, "Use j!background <número do background> para equipar um background")
			return
		}
		_, indexes := rinha.GetBackgrounds(galo.Cosmetics)
		if rinha.IsIntInList(value, indexes) {
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				galo.Background = value
				msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce equipou o background **%s**", msg.Author.Mention(), rinha.Cosmetics[galo.Cosmetics[galo.Background]].Name))
				return galo, nil
			})
		} else {
			msg.Reply(context.Background(), session, "Número inválido")
		}
	}
}
