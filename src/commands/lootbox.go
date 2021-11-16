package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

var lootTypes = []string{"comum", "rara", "normal", "cosmetica", "epica", "lendaria", "items"}

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
			Description: fmt.Sprintf("Money: **%d**\nAsuraCoins: **%d**\nPity: **%d%%**\n\n[100] Lootbox comum: **%d**\n[400] Lootbox normal: **%d**\n[800] Lootbox rara: **%d**\n[1750] Lootbox epica: **%d**\n[2 AsuraCoins] Lootbox lendaria: **%d**\n[4 AsuraCoins] Lootbox items: **%d**\n[300] Lootbox cosmetica: **%d**\n\nUse `j!lootbox buy <tipo>` para comprar lootbox\nUse `j!lootbox open <tipo>` para abrir lootbox\n Use `j!changename` para trocar o nome do galo (precisa de 100 money)\n%s\n\n**[Comprar Moedas e XP](https://acnologla.github.io/asura-site/donate)**", galo.Money, galo.AsuraCoin, galo.Pity*rinha.PityMultiplier, galo.CommonLootbox, galo.Lootbox, galo.RareLootbox, galo.EpicLootbox, galo.LegendaryLootbox, galo.ItemsLootbox, galo.CosmeticLootbox, rinha.VipMessage(galo)),
		})
	}
	if len(args) == 0 {
		normal()
		return
	}
	if args[0] == "open" || args[0] == "abrir" {
		if 2 > len(args) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Você precisa decidir o tipo de lootbox para abrir\nj!lootbox open <tipo>")
			return
		}
		lootType := strings.ToLower(args[1])
		if !utils.Includes(lootTypes, lootType) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Tipo de caixa invalido, use j!lootbox para ver os tipos")
			return
		}
		if !rinha.HaveLootbox(galo, lootType) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Você não tem essa lootbox\nuse j!lootbox para ver as suas lootboxs")
			return
		}
		isCosmetic := lootType == "cosmetica"
		isItem := lootType == "items"
		if len(galo.Galos) >= 10 && !isCosmetic && !isItem {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Você atingiu o limite maximo de galos (10) use `j!equip` para remover um galo")
			return
		}
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere a sua rinha terminar antes de abrir lootboxs")
			return
		}
		battleMutex.RUnlock()
		oldPity := galo.Pity
		result, pity := rinha.Open(lootType, galo)
		avatar, _ := msg.Author.AvatarURL(512, true)
		extraMsg := ""
		sold := "no"
		embed := &disgord.Embed{
			Title: "Abrindo lootbox...",
		}
		if isCosmetic {
			newCosmetic := rinha.Cosmetics[result]
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				galo = rinha.GetNewLb(lootType, galo, false)
				if !rinha.IsIntInList(result, galo.Cosmetics) {
					galo.Cosmetics = append(galo.Cosmetics, result)
				} else {
					price := 150
					sold = "yes"
					galo.Money += price
					extraMsg = fmt.Sprintf("\nComo você ja tinha esse background voce ganhou **%d** de dinheiro", price)
				}
				galo.Pity = pity
				return galo, nil
			})
			message, err := msg.Reply(context.Background(), session, embed)
			if err == nil {
				embed.Description = "Selecionando cosmetico..."
				for i := 0; i < 6; i++ {
					rand := rinha.GetRandCosmetic()
					if i == 5 {
						embed.Title = "Lootbox open"
						rand = result
						embed.Description = "Voce abriu uma lootbox " + lootType + " e ganhou o cosmetico **" + newCosmetic.Name + "**\nRaridade: " + newCosmetic.Rarity.String() + extraMsg
						embed.Footer = &disgord.EmbedFooter{
							IconURL: avatar,
							Text:    rinha.CosmeticCommand(*newCosmetic),
						}
					}
					randCosmetic := rinha.Cosmetics[rand]
					embed.Color = randCosmetic.Rarity.Color()
					embed.Image = &disgord.EmbedImage{
						URL: randCosmetic.Value,
					}
					utils.Try(func() error {
						msgUpdater := handler.Client.Channel(message.ChannelID).Message(message.ID).UpdateBuilder()
						msgUpdater.SetEmbed(embed)
						_, err := msgUpdater.Execute()
						return err
					}, 3)
					time.Sleep(time.Millisecond * 3500)
				}
			}
			return
		}
		if isItem {
			newItem := rinha.Items[result]
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				galo = rinha.GetNewLb(lootType, galo, false)
				if !rinha.IsIntInList(result, galo.Items) {
					galo.Items = append(galo.Items, result)
				}
				galo.Pity = pity
				return galo, nil
			})
			embed.Title = "Lootbox open"
			embed.Description = fmt.Sprintf("Voce abriu uma caixa e items e ganhou o item: **%s**\nRaridade: **%s**", newItem.Name, rinha.LevelToString(newItem.Level))
			msg.Reply(context.Background(), session, embed)
			return
		}
		newGalo := rinha.Classes[result]
		tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			galo = rinha.GetNewLb(lootType, galo, false)
			if !rinha.HaveGalo(result, galo.Galos) && galo.Type != result {
				galo.Galos = append(galo.Galos, rinha.SubGalo{
					Type: result,
					Xp:   0,
				})
			} else {
				price, _ := rinha.Sell(newGalo.Rarity, 0, 0)
				sold = "yes"
				galo.Money += price
				extraMsg = fmt.Sprintf("\nComo você ja tinha esse galo voce ganhou **%d** de dinheiro", price)
			}
			galo.Pity = pity
			return galo, nil
		})
		telemetry.Debug(fmt.Sprintf("%s wins %s", tag, newGalo.Name), map[string]string{
			"galo":     newGalo.Name,
			"user":     strconv.FormatUint(uint64(msg.Author.ID), 10),
			"rarity":   newGalo.Rarity.String(),
			"lootType": lootType,
			"sold":     sold,
			"pity":     strconv.Itoa(oldPity),
		})
		message, err := msg.Reply(context.Background(), session, embed)
		if err == nil {
			embed.Description = "Selecionando galo..."
			for i := 0; i < 6; i++ {
				rand := rinha.GetRand()
				if i == 5 {
					embed.Title = "Lootbox open"
					rand = result
					embed.Description = "Voce abriu uma lootbox " + lootType + " e ganhou o galo **" + newGalo.Name + "**\nRaridade: " + newGalo.Rarity.String() + extraMsg
					embed.Footer = &disgord.EmbedFooter{
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
					msgUpdater := handler.Client.Channel(message.ChannelID).Message(message.ID).UpdateBuilder()
					msgUpdater.SetEmbed(embed)
					_, err := msgUpdater.Execute()
					return err
				}, 3)
				time.Sleep(time.Millisecond * 3500)
			}
		}
	} else if args[0] == "buy" || args[0] == "comprar" {
		if 2 > len(args) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Você precisa decidir o tipo de lootbox para comprar\nj!lootbox buy <tipo>")
			return
		}
		lootType := strings.ToLower(args[1])
		if !utils.Includes(lootTypes, lootType) {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", Tipo de caixa inválido, use j!lootbox para ver os tipos")
			return
		}
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere sua rinha terminar antes de comprar lootboxs")
			return
		}
		battleMutex.RUnlock()
		price, asuraCoins := rinha.GetPrice(lootType)
		if asuraCoins == 0 {
			err := rinha.ChangeMoney(msg.Author.ID, -price, price)
			if err != nil {
				msg.Reply(context.Background(), session, fmt.Sprintf("%s, Você precisa ter %d de dinheiro para comprar uma lootbox %s, use `j!lootbox` para ver seu dinheiro", msg.Author.Mention(), price, lootType))
				return
			}
		} else {
			e := false
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				if asuraCoins > galo.AsuraCoin {
					e = true
					return galo, nil
				}
				galo.AsuraCoin -= asuraCoins
				return galo, nil
			})
			if e {
				msg.Reply(context.Background(), session, fmt.Sprintf("%s, Você precisa ter %d de dinheiro para comprar uma lootbox %s, use `j!lootbox` para ver seu dinheiro", msg.Author.Mention(), asuraCoins, lootType))
				return
			}
		}
		rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
			return rinha.GetNewLb(lootType, galo, true), nil
		})
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Você comprou uma lootbox "+lootType+" use `j!lootbox open "+lootType+"` para abrir")
	} else {
		normal()
	}
}
