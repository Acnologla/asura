package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"asura/src/telemetry"
	"context"
	"fmt"
	"sync"

	"github.com/andersfylling/disgord"
)

//this command is not executed in the normal Way

func SendLootbox(msg *disgord.Message) {
	ctx := context.Background()
	rarity, _ := rinha.MessageRandomLootbox()
	users := []*disgord.User{}
	rand := rinha.GetRandByType(rarity)
	embeds := []*disgord.Embed{
		{
			Title:       fmt.Sprintf("Um galo enfurecido **%s** caiu no servidor", rinha.Classes[rand].Name),
			Description: "Clique no botão abaixo para participar desta batalha!",
			Color:       65535,
		},
	}
	embeds[0].Footer = &disgord.EmbedFooter{
		Text: "Para desativar ou mudar de canal essas mensagens use o comando /config (precisa ter a permissão de gerenciar servidor)",
	}
	message := disgord.CreateMessage{
		Embeds: embeds,
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
	}
	newMessage, err := msg.Reply(context.Background(), handler.Client, message)
	guild, _ := handler.Client.Cache().GetGuild(msg.GuildID)
	mutex := sync.Mutex{}
	if err == nil {
		telemetry.Debug(fmt.Sprintf("%s drop lootbox", guild.Name), map[string]string{
			"id": msg.GuildID.String(),
		})
		var itc *disgord.InteractionCreate
		handler.RegisterHandler(newMessage.ID, func(ic *disgord.InteractionCreate) {
			itc = ic
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
			if len(users) >= 10 {
				handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: "A batalha ja chegou ao maximo (5)",
					},
				})
				return
			}
			users = append(users, ic.Member.User)
			embed := getBattleEmbed(users, "Batalha contra o chefe", 1)
			embed.Footer.Text += " (para desativar ou mudar de canal essas mensagens use o comando /config)"
			handler.Client.SendInteractionResponse(context.Background(), ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.CreateInteractionResponseData{
					Embeds:     []*disgord.Embed{embed},
					Components: message.Components,
				},
			})

		}, 60)

		if len(users) == 0 {
			return
		}
		var usersDb []*entities.User
		for _, user := range users {
			u := database.User.GetUser(ctx, user.ID, "Galos", "Trials", "Items")
			usersDb = append(usersDb, &u)
		}
		sumOfLevels := 0
		for _, user := range usersDb {
			galo := rinha.GetEquippedGalo(user)
			sumOfLevels += rinha.CalcLevel(galo.Xp)
		}
		bossLevel := sumOfLevels / len(users)
		galoAdv := entities.Rooster{
			Xp:      rinha.CalcXP(bossLevel) + 1,
			Type:    rand,
			Equip:   true,
			Evolved: true,
			Resets:  len(users) / 2,
		}
		userAdv := entities.User{
			Galos:      []*entities.Rooster{&galoAdv},
			Attributes: [6]int{0, 0, 0, 0, 0, 150},
		}
		usernames := make([]string, len(usersDb))
		for i, user := range users {
			usernames[i] = user.Username
		}
		galo := rinha.GetEquippedGalo(usersDb[0])
		winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
			GaloAuthor:  usersDb[0],
			GaloAdv:     &userAdv,
			IDs:         [2]disgord.Snowflake{users[0].ID},
			AuthorName:  rinha.GetName(users[0].Username, *galo),
			AdvName:     "Chefe",
			AuthorLevel: rinha.CalcLevel(galo.Xp),
			AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
			Waiting:     usersDb,
			Usernames:   usernames,
		}, false)
		if winner == 0 {
			for _, user := range users {
				database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
					u.Money += 300
					database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
						r.Xp += 150 / (r.Resets + 1)
						return r
					})
					return u
				}, "Galos")
			}
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: "O boss foi derrotado\nRecompensas:\nDinheiro: **300**\nXp: **150**",
			})
		} else {
			handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
				Content: "O boss venceu",
			})
		}
	}

}
