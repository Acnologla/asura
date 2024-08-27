package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ranked",
		Description: "Jogue rinha ranqueada e ganhe vários premios",
		Run:         runRanked,
		Cooldown:    10,
		Category:    handler.Rinha,
		Aliases:     []string{"ranqueada"},
		Options: []*disgord.ApplicationCommandOption{
			{
				Name:        "view",
				Description: "Veja seu rank e recompensas",
				Type:        disgord.OptionTypeSubCommand,
			},
			{
				Name:        "battle",
				Description: "batalhe na ranqueada",
				Type:        disgord.OptionTypeSubCommand,
			},
			{
				Name:        "rank",
				Description: "Veja os melhores jogadores",
				Type:        disgord.OptionTypeSubCommand,
			},
		},
	})
}

func runRanked(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.User.ID, "Galos")
	command := itc.Data.Options[0].Name
	if command == "view" {
		userRank := rinha.GetRank(&user)
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Seu elo",
						Color:       65535,
						Description: fmt.Sprintf("Seu rank atual é: %s\nElo: **%d**", userRank.ToString(), user.Rank),
					},
				},
			},
		}

	} else if command == "rank" {
		handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Carregando...",
			},
		})
		users := database.User.SortUsers(ctx, DEFAULT_RANK_LIMIT, "rank")
		var description string
		for i, user := range users {
			u, err := handler.Client.User(user.ID).Get()
			if err == nil {
				tag := u.Username
				rank := rinha.GetRank(user)
				description += fmt.Sprintf("[**%d**] - %s\nRank: %s (**%d** Elo)\n", i+1, tag, rank.ToString(), user.Rank)
			}
		}
		embed := &disgord.Embed{
			Title:       "Rank",
			Color:       65535,
			Description: description,
		}
		handler.Client.EditInteractionResponse(ctx, itc, &disgord.UpdateMessage{
			Embeds:  &([]*disgord.Embed{embed}),
			Content: nil,
		})
		return nil

	}

	ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))
	authorRinha := isInRinha(ctx, itc.Member.User)
	if authorRinha != "" {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: rinhaMessage(itc.Member.User.Username, authorRinha).Data.Content,
			},
		}
	}
	lockEvent(ctx, itc.Member.User.ID, "Ranked")
	defer unlockEvent(ctx, itc.Member.User.ID)

	user = database.User.GetUser(ctx, itc.Member.User.ID, "Galos")

	userRank := rinha.GetRank(&user)
	itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "Voce entrou em uma fila de espera para a rinha rankeada",
		},
	})
	message, err := ch.CreateMessage(&disgord.CreateMessage{
		Embeds: []*disgord.Embed{
			{
				Title: "Procurando oponente...",
				Color: 65535,
			},
		},
	})
	if err != nil {
		return nil
	}

	c := engine.AddToMatchMaking(itc.Member.User, user.RankedLastFight, message, engine.ArenaRanked, user.Rank)
	result := <-c
	if result.Type == engine.TimeExceeded {
		handler.Client.Channel(message.ChannelID).Message(message.ID).Update(&disgord.UpdateMessage{
			Embeds: &([]*disgord.Embed{
				{
					Title: "Nao consegui achar um oponente para voce",
					Color: 65535,
				},
			}),
		})
		return nil
	}
	if result.Type == engine.ArenaTie {
		return nil
	}

	advUser := database.User.GetUser(ctx, result.Adv, "Galos")
	if result.Type == engine.ArenaWin {
		database.User.UpdateUser(ctx, itc.Member.User.ID, func(u entities.User) entities.User {
			rank := rinha.CalcRank(u.Rank, advUser.Rank)
			u.Rank += rank
			u.RankedLastFight = result.Adv
			newRank := rinha.GetRank(&u)
			if newRank.Name != userRank.Name && u.Rank > u.LastRank {
				money, xp, asuraCoins := rinha.CalcRankedPrize(newRank)
				u.Money += money
				u.AsuraCoin += asuraCoins
				u.LastRank = u.Rank
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					r.Xp += xp
					return r
				})
				go ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{
						{
							Title:       "Ranqueada",
							Color:       16776960,
							Description: fmt.Sprintf("Parabens %s, voce upou para o rank %s, recebeu:\nMoney: **%d**\nXp: **%d**\n\nAsuraCoins: **%d**", itc.Member.User.Username, newRank.ToString(), money, xp, asuraCoins),
						},
					},
				})
			}
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{
					{
						Title:       "Ranqueada",
						Color:       16776960,
						Description: fmt.Sprintf("Parabens %s, voce venceu\nGanhou: **%d** de elo", itc.Member.User.Username, rank),
					},
				},
			})

			return u
		}, "Galos")
	} else if result.Type == engine.ArenaLose {
		database.User.UpdateUser(ctx, itc.Member.User.ID, func(u entities.User) entities.User {
			rank := rinha.CalcRank(advUser.Rank, u.Rank)
			u.Rank -= rank
			if u.Rank < 0 {
				u.Rank = 0
			}
			u.RankedLastFight = result.Adv
			ch.CreateMessage(&disgord.CreateMessage{
				Embeds: []*disgord.Embed{
					{
						Title:       "Ranqueada",
						Color:       16776960,
						Description: fmt.Sprintf("Parabens %s, voce perdeu\nPerdeu: **%d** de elo", itc.Member.User.Username, rank),
					},
				},
			})
			return u
		}, "Galos")
	}
	return nil
}
