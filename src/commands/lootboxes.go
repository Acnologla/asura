package commands

import (
	"asura/src/utils"
	"asura/src/rinha"
	"asura/src/handler"
	"strings"
	"asura/src/entities"
	"asura/src/database"
	"fmt"
	"context"
	
	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "lootboxes",
		Description: "compre e/ou abra diversas lootboxes rapidamente!",
		Run:         runLootboxs,
		Cooldown:    5,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type: 			disgord.OptionTypeSubCommand,
			Name: 			"buy",
			Description:	"compre lootboxes",
			Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
				Type: 		disgord.OptionTypeString,
				Name: 		"type",
				Required: 	true,
				Description:"tipo da lootbox (lendaria, items, epica, rara, cosmetico, normal, comum)",
			},
			&disgord.ApplicationCommandOption{
				Type: 		disgord.OptionTypeNumber,
				Name: 		"quantity",
				Required: 	true,
				Description:"quantidade para comprar",
				MinValue:   1,
				MaxValue:	100,
			}),
		},
		&disgord.ApplicationCommandOption{
			Type: 			disgord.OptionTypeSubCommand,
			Name: 			"open",
			Description:	"abra as lootboxes",
			Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
				Type: 		disgord.OptionTypeString,
				Name: 		"type",
				Required: 	true,
				Description:"tipo da lootbox (lendaria, items, epica, rara, cosmetico, normal, comum)",
			},
			&disgord.ApplicationCommandOption{
				Type: 		disgord.OptionTypeNumber,
				Name: 		"quantity",
				Required: 	true,
				Description:"quantidade para abrir",
				MinValue:   1,
				MaxValue:	10,
			}),
		}),
	})
}

func runLootboxs(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	author := itc.Member.User
	command := itc.Data.Options[0].Name

	authorRinha := isInRinha(ctx, author)
	if authorRinha != "" {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: rinhaMessage(author.Username, authorRinha).Data.Content,
			},
		}
	}

	lockEvent(ctx, author.ID, "Comprando ou abrindo lootboxes.")
	defer unlockEvent(ctx, author.ID)

	switch command {
	case "buy":
		lootType := itc.Data.Options[0].Options[0].Value.(string)
		lootboxes := rinha.LootNames
		found := false

		for _, name := range lootboxes {
			if lootType == name {
				found = true
				break
			}
		}
		if !found {
			lootNamesMapped := strings.Join(lootboxes[:], ", ")
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: fmt.Sprintf("lootbox inválida\nopções: %s", lootNamesMapped),
				},
			}
		}

		quantity := int(itc.Data.Options[0].Options[1].Value.(float64))
		index := rinha.GetLbIndex(lootType)
		price := rinha.Prices[index] //indice 0 = money, indice 1 = asuracoin
		totalPrice := 0
		done := false
		moneyType := "dinheiro"

		database.User.UpdateUser(ctx, author.ID, func (u entities.User) entities.User {
			if price[0] == 0 {
				totalPrice = quantity * price[1]
				moneyType = "asura coins"
				if u.AsuraCoin >= totalPrice {
					done = true
					u.AsuraCoin -= totalPrice
				}
			} else {
				totalPrice = quantity * price[0]
				if u.Money >= totalPrice {
					done = true
					u.Money -= totalPrice
				}
			}
			if done {
				database.User.InsertManyItem(ctx, itc.Member.UserID, u.Items, index, entities.LootboxType, quantity)
				itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: fmt.Sprintf("Você comprou **%d** lootboxes **%s** por **%d %s**", quantity, lootType, totalPrice, moneyType),
					},
				})
			} else {
				itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: fmt.Sprintf("Você não tem **%d %s** para comprar **%d** lootboxes **%s**", totalPrice, moneyType, quantity, lootType),
					},
				})
			}
			return u
		}, "Items")
	case "open":
		user := database.User.GetUser(ctx, itc.Member.UserID, "Items")
		lootbox := rinha.GetLootboxes(&user)
		quantity := int(itc.Data.Options[0].Options[1].Value.(float64))
		lootType := itc.Data.Options[0].Options[0].Value.(string)

		lootboxesNames := rinha.LootNames
		found := false

		for _, name := range lootboxesNames {
			if lootType == name {
				found = true
				break
			}
		}
		if !found {
			lootNamesMapped := strings.Join(lootboxesNames[:], ", ")
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: fmt.Sprintf("lootbox inválida\nopções: %s", lootNamesMapped),
				},
			}
		}

		lootIndex := rinha.GetLbIndex(lootType)
		lootTotal := rinha.GetLootboxQuantity(lootbox, lootType)
		sellMsg := "\n"
		msg := "\n"
		image := ""
		name := ""
		newVal := -1
		winType := ""
		stopMsg := ""
		pity := 0
		var rarity rinha.Rarity

		if quantity > lootTotal {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: fmt.Sprintf("Você não tem essa quantia de lootboxes **%s**", lootType),
				},
			}
		}

		itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Abrindo lootboxes",
			},
		})

		database.User.UpdateUser(ctx, author.ID, func (u entities.User) entities.User {
			if lootType != "cosmetica" && lootType != "items" {
				if len(u.Galos) >= 10 {
					itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
						Type: disgord.InteractionCallbackChannelMessageWithSource,
						Data: &disgord.CreateInteractionResponseData{
					 		Content: "Voce já chegou no limite de galos (10)",
					 	},
					 })
					return u
				}
			}

			if lootType == "items" {
				for i := 1; i <= quantity; i++ {
					userDatabase := database.User.GetUser(ctx, author.ID, "Items", "Galos")

					lbID, _ := rinha.GetLbID(userDatabase.Items, lootIndex)
					database.User.RemoveItem(ctx, userDatabase.Items, lbID)
					newVal, pity = rinha.Open(lootIndex, &u)
					u.Pity = pity
					database.User.InsertItem(ctx, author.ID, userDatabase.Items, newVal, entities.CosmeticType)
					item := rinha.Items[newVal]
					if item.Level == 4 {
						rarity = rinha.Mythic
					} else {
						rarity = rinha.Legendary
					}
					winType = "item"
					name = item.Name
					msg += fmt.Sprintf("\nVocê ganhou o %s **%s** (%s)", winType, name, rarity)
				}
			} else if lootType == "cosmetica" {
				for i := 1; i <= quantity; i++ {
					userDatabase := database.User.GetUser(ctx, author.ID, "Items", "Galos")

					lbID, _ := rinha.GetLbID(userDatabase.Items, lootIndex)
					database.User.RemoveItem(ctx, userDatabase.Items, lbID)
					newVal, pity = rinha.Open(lootIndex, &u)
					u.Pity = pity
					database.User.InsertItem(ctx, author.ID, userDatabase.Items, newVal, entities.CosmeticType)
					cosmetic := rinha.Cosmetics[newVal]
					image = cosmetic.Value
					winType = "cosmetico"
					name = cosmetic.Name
					rarity = cosmetic.Rarity
					msg += fmt.Sprintf("\nVocê ganhou o %s [%s](%s) (%s)", winType, name, image, rarity)
				}
			} else {
				for i := 1; i <= quantity; i++ {
					userDatabase := database.User.GetUser(ctx, author.ID, "Items", "Galos")
					if len(userDatabase.Galos) >= 10 {
						stopMsg = fmt.Sprintf("\na abertura do restante das lootboxes foi cancelada pois você atingiu o limite de 10 galos.")
						//fmt.Println(fmt.Sprintf("user %s stoped open many lootboxes because already have 10 rooster", author.Username))
						break
						return u
					}

					lbID, _ := rinha.GetLbID(userDatabase.Items, lootIndex)
					database.User.RemoveItem(ctx, userDatabase.Items, lbID)
					newVal, pity = rinha.Open(lootIndex, &u)
					u.Pity = pity
					winType = "galo"
					galo := rinha.Classes[newVal]
					image = rinha.Sprites[0][newVal-1]
					name = galo.Name
					rarity = galo.Rarity
					if !rinha.HaveRooster(userDatabase.Galos, newVal) {
						database.User.InsertRooster(ctx, &entities.Rooster{
							UserID: author.ID,
							Type:   newVal,
						})
						msg += fmt.Sprintf("\nvocê ganhou o %s [%s](%s) (%s)", winType, name, image, rarity)
					} else {
						money, _ := rinha.Sell(galo.Rarity, 0, 0)
						u.Money += money
						sellMsg += fmt.Sprintf("\no %s [%s](%s) foi vendido automaticamente por **%d** moedas", winType, galo.Name, image, money)
					}
				}
			}
			return u
		}, "Galos", "Items")

		ch := handler.Client.Channel(itc.ChannelID)
		if len(msg) <= 1 {
			msg = "Voce já chegou no limite de galos (10)"
		}
		ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{
				{
					Color:       16776960,
					Title:       fmt.Sprintf("**Abertura de lootboxes - %d %s**", quantity, lootType),
					Description: fmt.Sprintf("%s%s%s", msg, sellMsg, stopMsg),
				},
			},
		})
	}
	return nil
}