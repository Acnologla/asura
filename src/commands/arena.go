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
		Name:        "arena",
		Description: "Arena de galos",
		Run:         runArena,
		Cooldown:    10,
		Category:    handler.Rinha,
		Options: []*disgord.ApplicationCommandOption{
			{
				Name:        "ingresso",
				Description: "Comprar um ingresso para a arena",
				Type:        disgord.OptionTypeSubCommand,
			},
			{
				Name:        "status",
				Description: "status da arena",
				Type:        disgord.OptionTypeSubCommand,
			},
			{
				Name:        "batalhar",
				Description: "batalha na arena",
				Type:        disgord.OptionTypeSubCommand,
			},
		},
	})
}

const arenaPrice = 200

func runArena(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.User.ID, "Galos")
	command := itc.Data.Options[0].Name
	if command == "ingresso" {
		database.User.UpdateUser(ctx, itc.Member.User.ID, func(u entities.User) entities.User {
			if u.ArenaActive {
				itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: "Use /arena status para ver o status da arena",
					},
				})
				return u
			}
			if arenaPrice > user.Money {
				itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: fmt.Sprintf("Você precisa ter **%d** de dinheiro para comprar um ingresso na arena", arenaPrice),
					},
				})
				return u
			}
			u.Money -= arenaPrice
			u.ArenaActive = true
			u.ArenaWin = 0
			u.ArenaLose = 0
			u.ArenaLastFight = 0
			itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Você comprou um ingresso para a arena use /arena batalhar para batalhar",
				},
			})
			return u
		})
		return nil
	} else if command == "status" {
		if !user.ArenaActive {
			return entities.CreateMsg().Content("Você precisa comprar um ingresso para a arena /arena ingresso").Res()
		}
		return entities.CreateMsg().Content(fmt.Sprintf("Você esta com **%d/12** vitorias na arena e  **%d** derrotas, com 3 voce é eliminado", user.ArenaWin, user.ArenaLose)).Res()
	}

	if !user.ArenaActive {
		return entities.CreateMsg().Content("Você precisa comprar um ingresso para a arena").Res()
	}
	ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))
	authorRinha := isInRinha(ctx, itc.Member.User)
	if authorRinha != "" {
		ch.CreateMessage(&disgord.CreateMessage{
			Content: rinhaMessage(itc.Member.User.Username, authorRinha).Data.Content,
		})
		return nil
	}
	lockEvent(ctx, itc.Member.User.ID, "Arena")
	defer unlockEvent(ctx, itc.Member.User.ID)
	itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "Voce entrou em uma fila de espera para a arena",
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

	c := engine.AddToMatchMaking(itc.Member.User, user.ArenaLastFight, message)
	result := <-c
	if result == engine.TimeExceeded {
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
	if result == engine.ArenaTie {
		return nil
	}
	if result == engine.ArenaWin {
		database.User.UpdateUser(ctx, itc.Member.User.ID, func(u entities.User) entities.User {
			u.ArenaWin++
			if u.ArenaWin >= 12 {
				xp, money := rinha.CalcArena(&u)
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					r.Xp += xp
					return r
				})
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{
						{
							Title:       "Arena",
							Color:       16776960,
							Description: fmt.Sprintf("Parabens %s, voce atingiu o limite de vitorias na arena\nPremios:\nXp: **%d**\nMoney: **%d**", itc.Member.User.Username, xp, money),
						},
					},
				})
			} else {
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{
						{
							Title:       "Arena",
							Color:       16776960,
							Description: fmt.Sprintf("Parabens %s, voce venceu\n %d/12 Vitorias", itc.Member.User.Username, u.ArenaWin),
						},
					},
				})
			}
			return u
		}, "Galos")
	} else if result == engine.ArenaLose {
		database.User.UpdateUser(ctx, itc.Member.User.ID, func(u entities.User) entities.User {
			u.ArenaLose++
			if u.ArenaLose >= 3 {
				xp, money := rinha.CalcArena(&u)
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					r.Xp += xp
					return r
				})
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{
						{
							Color:       16711680,
							Title:       "Arena",
							Description: fmt.Sprintf("Parabens %s, voce atingiu o limite de derrotas na arena\nPremios:\nXp: **%d**\nMoney: **%d**", itc.Member.User.Username, xp, money),
						},
					},
				})
			} else {
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{
						{
							Color:       16711680,
							Title:       "Arena",
							Description: fmt.Sprintf("Parabens %s, voce perdeu. %d/3 Derrotas", itc.Member.User.Username, u.ArenaLose),
						},
					},
				})
			}
			return u
		}, "Galos")
	}
	return nil
}
