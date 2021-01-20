package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"time"
	"github.com/andersfylling/disgord"
	"strconv"
	"strings"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"lootbox", "lb", "money", "dinheiro", "bal", "balance"},
		Run:       runLootbox,
		Available: true,
		Cooldown:  2,
		Usage:     "j!lootbox",
		Help:      "Abra lootboxs",
		Category:  1,
	})
}

func runLootbox(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	normal := func() {
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title:       "Lootbox",
			Color:       65535,
			Description: fmt.Sprintf("Money: **%d**\n\n[100] Lootbox comum: **%d**\n[400] Lootbox normal: **%d**\n[800] Lootbox rara: **%d**\n\nUse `j!lootbox buy <tipo>` para comprar lootbox\nUse `j!lootbox open <tipo>` para abrir lootbox\n Use `j!changename` para trocar o nome do galo (precisa de 100 Gold)", galo.Money, galo.CommonLootbox, galo.Lootbox, galo.RareLootbox),
		})
	}
	if len(args) == 0 {
		normal()
		return
	}
	if args[0] == "open" || args[0] == "abrir" {
		if 2 > len(args) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce precisa decidir o tipo de lootbox para abrir\nj!lootbox open <tipo>")
			return
		}
		lootType := strings.ToLower(args[1])
		if lootType != "comum" && lootType != "rara" && lootType != "normal" {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Tipo de caixa invalido, use j!lootbox para ver os tipos")
			return
		}
		if !rinha.HaveLootbox(galo, lootType) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem essa lootbox\nuse j!lootbox para ver suas lootbox")
			return
		}
		if len(galo.Galos) >= 7 {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce atingiu o limite maximo de galos (7) use `j!equip` para remover um galo")
			return
		}
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere sua rinha terminar antes de abrir lootboxs")
			return
		}
		battleMutex.RUnlock()
		result := rinha.Open(lootType)
		newGalo := rinha.Classes[result]
		tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
		extraMsg := ""
		sold := "no"
		if !rinha.HaveGalo(result, galo.Galos) && galo.Type != result {
			galo.Galos = append(galo.Galos, rinha.SubGalo{
				Type: result,
				Xp:   0,
			})
		} else {
			price := rinha.Sell(newGalo.Rarity, 0)
			sold = "yes"
			rinha.ChangeMoney(msg.Author.ID, price, 0)
			extraMsg = fmt.Sprintf("\nComo voce ja tinha esse galo voce ganhou **%d** de dinheiro", price)
		}
		telemetry.Debug(fmt.Sprintf("%s wins %s", tag, newGalo.Name), map[string]string{
			"galo":     newGalo.Name,
			"user":     strconv.FormatUint(uint64(msg.Author.ID), 10),
			"rarity":   newGalo.Rarity.String(),
			"lootType": lootType,
			"sold":     sold,
		})
		update := map[string]interface{}{
			"galos": galo.Galos,
		}
		avatar, _ := msg.Author.AvatarURL(512, true)
		newLb, value := rinha.GetNewLb(lootType, galo, false)
		update[newLb] = value
		rinha.UpdateGaloDB(msg.Author.ID, update)
		embed := &disgord.Embed{
			Title: "Abrindo lootbox",
		}
		message, err :=  msg.Reply(context.Background(), session, embed)
		if err == nil {
			embed.Description = "Selecionando galo..."
			for i:=0; i < 6;i++{
				rand := rinha.GetRand()
				if i == 5{
					embed.Title = "Lootbox open"
					rand = result
					embed.Description =  "Voce abriu uma lootbox " + lootType + " e ganhou o galo **" + newGalo.Name + "**\nRaridade: " + newGalo.Rarity.String() + extraMsg
					embed.Footer =  &disgord.EmbedFooter{
						IconURL: avatar,
						Text:    "Use j!equipar para equipar ou vender esse galo",
					}
				}
				randClass := rinha.Classes[rand]
				embed.Color = randClass.Rarity.Color()
				embed.Image = &disgord.EmbedImage{
					URL: rinha.Sprites[0][rand-1],
				}
				utils.Try(func() error {
					msgUpdater := handler.Client.Channel(message.ChannelID).Message(message.ID).Update()
					msgUpdater.SetEmbed(embed)
					_, err := msgUpdater.Execute()
					return err
				}, 3)
				time.Sleep(time.Millisecond * 2500)
			}
		}
	} else if args[0] == "buy" || args[0] == "comprar" {
		if 2 > len(args) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce precisa decidir o tipo de lootbox para comprar\nj!lootbox buy <tipo>")
			return
		}
		lootType := strings.ToLower(args[1])
		if lootType != "comum" && lootType != "rara" && lootType != "normal" {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Tipo de caixa invalido, use j!lootbox para ver os tipos")
			return
		}
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere sua rinha terminar antes de comprar lootboxs")
			return
		}
		battleMutex.RUnlock()
		price := rinha.GetPrice(lootType)
		err := rinha.ChangeMoney(msg.Author.ID, -price, price)
		if err != nil {
			msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce precisa ter %d de dinheiro para comprar uma lootbox %s, use `j!lootbox` para ver seu dinheiro", msg.Author.Mention(), price, lootType))
			return
		}
		update := map[string]interface{}{
			"galos": galo.Galos,
		}
		newLb, value := rinha.GetNewLb(lootType, galo, true)
		update[newLb] = value
		rinha.UpdateGaloDB(msg.Author.ID, update)
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce comprou uma lootbox "+lootType+" use `j!lootbox open "+lootType+"` para abrir")
	} else {
		normal()
	}
}
