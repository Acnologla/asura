package commands

import (
	"asura/src/cache"
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"asura/src/utils"
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "natal",
		Description: "Comando do evento de natal",
		Run:         runNatal,
		Cooldown:    5,
		Category:    handler.Rinha,
	})
}

const bossN = 52
const bossChance = 1

const (
	bossXP    = 400
	bossMoney = 500
)

func runNatal(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	discordUser := itc.Member.User
	userID := discordUser.ID
	redisKey := fmt.Sprintf("natalBattle:%s", userID.String())
	lastTime := cache.Client.Get(ctx, redisKey)
	if lastTime.Val() != "" {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "O chefe só estara disponivel daqui 3 horas",
			},
		}
	}

	authorRinha := isInRinha(ctx, discordUser)
	if authorRinha != "" {

		return rinhaMessage(discordUser.Username, authorRinha)
	}
	lockEvent(ctx, userID, "Arena")
	defer unlockEvent(ctx, userID)

	u := database.User.GetUser(ctx, userID, "Items", "Galos", "Trials")
	rooster := rinha.GetEquippedGalo(&u)
	duration := time.Hour * 3
	cache.Client.Set(ctx, redisKey, "1", duration)
	rarity := rinha.GetRarity(rooster)
	gAdv := &entities.Rooster{
		Type:    bossN,
		Xp:      rinha.CalcXP(65),
		Evolved: rarity > rinha.Epic,
		Resets:  1 + rooster.Resets,
	}

	uAdv := entities.User{
		Galos:      []*entities.Rooster{gAdv},
		Attributes: u.Attributes,
		Items:      u.Items,
		Upgrades:   u.Upgrades,
	}

	if rarity == rinha.Mythic {
		uAdv.Attributes[0] += uAdv.Attributes[0]
	}

	if rarity >= rinha.Legendary {
		uAdv.Attributes[0] += 150
	}
	itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "A batalha esta iniciando",
		},
	})

	winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
		GaloAuthor: &u,
		GaloAdv:    &uAdv,
		IDs:        [2]disgord.Snowflake{discordUser.ID},

		AuthorName:  rinha.GetName(discordUser.Username, *rooster),
		AdvName:     "Boss " + rinha.Classes[bossN].Name,
		AuthorLevel: rinha.CalcLevel(rooster.Xp),
		AdvLevel:    rinha.CalcLevel(gAdv.Xp),
		NoItems:     false,
	}, false)

	if winner == -1 {
		return nil
	}

	ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))

	if winner == 0 {
		chance := utils.RandInt(101)
		winG := bossChance >= chance
		if u.EventStatus > 0 {
			winG = false
		}
		database.User.UpdateUser(ctx, userID, func(u entities.User) entities.User {
			database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
				r.Xp += bossXP
				return r
			})
			u.Money += bossMoney
			if winG {
				u.EventStatus = 1
				database.User.InsertRooster(ctx, &entities.Rooster{
					UserID: u.ID,
					Type:   bossN,
				})
			}
			return u
		}, "Galos")
		description := fmt.Sprintf("Voce derrotou o chefe de natal e ganhou\nXP: **%d**\nDinheiro: **%d**", bossXP, bossMoney)
		if winG {
			description += "\n\n**Voce ganhou um presente especial**"
		}
		ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{{
				Color:       16776960,
				Title:       "Natal",
				Description: description,
			}},
		})
	} else {
		ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{{
				Color:       16711680,
				Title:       "Natal",
				Description: "Parabens, voce perdeu para o chefe de natal",
			}},
		})
	}
	return nil
}
