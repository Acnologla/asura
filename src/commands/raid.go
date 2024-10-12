package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"context"
	"fmt"
	"sync"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "raid",
		Description: translation.T("RunRaid", "pt"),
		Run:         runRaid,
		Cooldown:    8,
		Category:    handler.Rinha,
	})
}

func isInRarityArray(arr []rinha.Rarity, rarity rinha.Rarity) bool {
	for _, r := range arr {
		if r == rarity {
			return true
		}
	}
	return false
}

func generateKeyOpts(user *entities.User) (opts []*disgord.SelectMenuOption) {
	raritys := []rinha.Rarity{}
	for _, item := range user.Items {
		if item.Type == entities.KeyType && !isInRarityArray(raritys, rinha.Rarity(item.ItemID)) {
			rarity := rinha.Rarity(item.ItemID)
			opts = append(opts, &disgord.SelectMenuOption{
				Label: fmt.Sprintf("Chave raid %s (%d/%d)", rarity.String(), item.Extra, rinha.CalcMaxRaidBattles(rarity)),
				Value: item.ID.String(),
			})
			raritys = append(raritys, rarity)
		}
	}
	return
}

func createFight(ctx context.Context, itc *disgord.InteractionCreate, user *entities.User, key *entities.Item) {
	users := []*disgord.User{
		itc.Member.User,
	}
	keyRarity := rinha.Rarity(key.ItemID)
	raidTitle := fmt.Sprintf("Raid %s", keyRarity.String())
	/*
		handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: fmt.Sprintf("A raid %s foi iniciada, reaja com o botão para entrar na batalha", keyRarity.String()),
			},
		})
		/*
			err := handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Embeds: []*disgord.Embed{getBattleEmbed(users, raidTitle)},
					Components: []*disgord.MessageComponent{
						{
							Type: disgord.MessageComponentActionRow,
							Components: []*disgord.MessageComponent{
								{
									Type:     disgord.MessageComponentButton,
									Label:    "Entrar",
									CustomID: "joinBattle",
									Style:    disgord.Primary,
								},
							},
						},
					},
				},
			})
	*/
	msg, err := handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
		Embeds: []*disgord.Embed{getBattleEmbed(users, raidTitle)},
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Entrar",
						CustomID: "joinBattle",
						Style:    disgord.Primary,
					},
				},
			},
		}})

	if err != nil {
		return
	}
	mutex := sync.Mutex{}

	handler.RegisterHandler(msg.ID, func(ic *disgord.InteractionCreate) {
		mutex.Lock()
		defer mutex.Unlock()
		if isInUsers(users, ic.Member.User) {
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Voce ja esta na batalha",
				},
			})
			return
		}
		if len(users) >= 5 {
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "A batalha ja chegou ao maximo (5)",
				},
			})
			return
		}
		users = append(users, ic.Member.User)
		handler.Client.EditInteractionResponse(ctx, itc, &disgord.UpdateMessage{
			Embeds: &([]*disgord.Embed{getBattleEmbed(users, raidTitle)}),
		})
		handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: fmt.Sprintf("%s Entrou na batalha", ic.Member.User.Mention()),
			},
		})
	}, 120)
	var usersDb []*entities.User
	for _, user := range users {
		u := database.User.GetUser(ctx, user.ID, "Galos", "Trials")
		usersDb = append(usersDb, &u)
	}

	usernames := make([]string, len(usersDb))
	for i, user := range users {
		usernames[i] = user.Username
	}
	galo := rinha.GetEquippedGalo(usersDb[0])
	xpMultiplier, atributeMultiplier := rinha.GetMultipliers(keyRarity)
	galoXp := (rinha.CalcXP(15 + key.Extra*5)) * xpMultiplier
	galoAttrs := 140 * atributeMultiplier
	galoOthersAttr := 20 * atributeMultiplier
	galoAdv := &entities.Rooster{
		Type:    rinha.GetRandByType(keyRarity),
		Xp:      galoXp,
		Evolved: key.Extra > 5,
		Resets:  key.Extra,
	}
	userAdv := entities.User{
		Galos:      []*entities.Rooster{galoAdv},
		Attributes: [5]int{galoAttrs, galoOthersAttr, galoOthersAttr, galoOthersAttr, galoOthersAttr},
	}

	winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
		GaloAuthor:  usersDb[0],
		GaloAdv:     &userAdv,
		IDs:         [2]disgord.Snowflake{users[0].ID},
		AuthorName:  rinha.GetName(users[0].Username, *galo),
		AdvName:     "Chefe da raid",
		AuthorLevel: rinha.CalcLevel(galo.Xp),
		AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
		Waiting:     usersDb,
		Usernames:   usernames,
	}, false)
	lb := -1
	if winner == 0 {
		xp, money := rinha.CalcRaidBattleWinPrize(keyRarity)
		for _, user := range users {
			database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
				u.Money += money / len(users)
				if user.ID == users[0].ID {
					if key.Extra+1 == rinha.CalcMaxRaidBattles(keyRarity) {
						database.User.RemoveItem(ctx, u.Items, key.ID)
						lb = rinha.RaidLootbox(keyRarity)
						database.User.InsertItem(ctx, user.ID, u.Items, lb, entities.LootboxType)
					} else {
						database.User.UpdateItem(ctx, &u, key.ID, func(i entities.Item) entities.Item {
							i.Extra++
							return i
						})
					}
				}
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					r.Xp += (xp / len(users)) / (r.Resets + 1)
					return r
				})
				return u
			}, "Galos", "Items")
		}
		handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{{
				Color:       16776960,
				Title:       "Raid battle",
				Description: fmt.Sprintf("A batalha da raid raid foi completada, os participantes receberam **%d** de dinheiro e **%d** de xp divididos entre si\nStatus: **%d/%d**", money, xp, key.Extra+1, rinha.CalcMaxRaidBattles(keyRarity)),
			}},
		})
		if lb != -1 {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{{
					Color:       65535,
					Title:       "Raid win",
					Description: fmt.Sprintf("A raid foi completada, o dono da chave recebeu uma caixa **%s**", rinha.LootNames[lb]),
				}},
			})
		}
	}
	if winner == 1 {
		handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
			Content: "Parabéns vocês perderam",
		})
	}
}

func runRaid(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items", "Trials")
	opts := generateKeyOpts(&galo)
	if len(opts) == 0 {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Voce ainda não tem nenhuma chave, consiga chaves ganhando treinos!",
			},
		}
	}
	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title:       "Raids",
					Color:       65535,
					Description: "Selecione abaixo uma chave para abrir a raid, quanto maior a raridade mais dificil a raid (e tambem maior a recompensa).\n\nO premio sera dividido entre os participantes.\n\nAo final de cada raid o dono da chave ganhara uma lootbox de acordo com a dificuldade da raid",
				},
			},
			Components: []*disgord.MessageComponent{
				{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{
						{
							Type:        disgord.MessageComponentButton + 1,
							Style:       disgord.Primary,
							Placeholder: "Selecione a chave para abrir a raid",
							CustomID:    "selectKey",
							Options:     opts,
							MaxValues:   1,
						},
					},
				},
			},
		},
	})
	done := false
	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		if done {
			return
		}
		userIC := ic.Member.User
		if userIC.ID != itc.Member.UserID {
			return
		}
		if len(ic.Data.Values) == 0 {
			return
		}
		val := ic.Data.Values[0]
		var key *entities.Item
		database.User.UpdateUser(ctx, userIC.ID, func(u entities.User) entities.User {
			for _, item := range u.Items {
				if item.ID.String() == val {
					key = item
					break
				}
			}
			return u
		}, "Items")
		if key != nil {
			done = true
			handler.Client.Channel(ic.ChannelID).Message(ic.Message.ID).Delete()
			createFight(ctx, ic, &galo, key)
		}
	}, 60)
	return nil
}
