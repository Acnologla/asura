package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"context"
	"fmt"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "attributes",
		Description: translation.T("AtributtesHelp", "pt"),
		Run:         runAttributes,
		Cooldown:    12,
		Category:    handler.Profile,
		Aliases:     []string{"atributos"},
	})
}

const MAX_POINTS = 250

func getPointOptions(selectPoints int) []*disgord.SelectMenuOption {
	return []*disgord.SelectMenuOption{
		{
			Label:   "1 Ponto",
			Value:   "1",
			Default: selectPoints == 1,
		},
		{
			Label:   "10 Pontos",
			Value:   "10",
			Default: selectPoints == 10,
		},
		{
			Label:   "100 Pontos",
			Value:   "100",
			Default: selectPoints == 100,
		},
	}
}

var attributeButtons = []*disgord.MessageComponent{
	{
		Type:     disgord.MessageComponentButton,
		Label:    "Vida",
		CustomID: "1",
		Style:    disgord.Primary,
	},
	{
		Type:     disgord.MessageComponentButton,
		Label:    "Dano",
		CustomID: "2",
		Style:    disgord.Primary,
	},
	{
		Type:     disgord.MessageComponentButton,
		Label:    "Sorte",
		CustomID: "3",
		Style:    disgord.Primary,
	},
}

func calcPoints(userDb entities.User) int {
	return int(userDb.UserXp / 100)
}

func calcAvailPoints(userDb entities.User) int {
	return calcPoints(userDb) - (userDb.Attributes[0] + userDb.Attributes[1] + userDb.Attributes[2])
}

func generateDesc(userDb entities.User, selectPoints int) (desc string) {
	points := calcAvailPoints(userDb)
	healthPoints := userDb.Attributes[0]
	attackPoints := userDb.Attributes[1]
	luckyPoints := userDb.Attributes[2]
	desc += fmt.Sprintf("Vida extra: **%d**\nDano extra: **%d**\nSorte: **%d**", healthPoints, attackPoints, luckyPoints)
	desc += "\n\nVoce ganha um ponto a cada 100 rinhas\nUm ponto em vida te da 3 de vida\nUm ponto em dano te da 1 de dano"
	desc += fmt.Sprintf("\n\nVoce tem **%d** pontos para gastar\nVoce esta colocando atualmente **%d** pontos", points, selectPoints)
	return
}

func runAttributes(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	userDb := database.User.GetUser(ctx, user.ID)
	selectPoints := 1
	desc := generateDesc(userDb, selectPoints)
	response := &disgord.CreateInteractionResponseData{
		Embeds: []*disgord.Embed{{
			Title:       "Atributos",
			Color:       65535,
			Description: desc,
		}},
	}

	if calcAvailPoints(userDb) > 0 {
		response.Components = append(response.Components, &disgord.MessageComponent{

			Type: disgord.MessageComponentActionRow,
			Components: []*disgord.MessageComponent{
				{
					Type:        disgord.MessageComponentSelectMenu,
					Options:     getPointOptions(selectPoints),
					Placeholder: "Selecione a quantidade de pontos para botar",
					CustomID:    "points-select",
					MaxValues:   1,
				},
			},
		},
			&disgord.MessageComponent{
				Type:       disgord.MessageComponentActionRow,
				Components: attributeButtons,
			})
		handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: response,
		})
		updateMessage := func(itc *disgord.InteractionCreate) {
			r := &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{{
					Title:       "Atributos",
					Color:       65535,
					Description: generateDesc(userDb, selectPoints),
				}},
			}
			r.Components = response.Components
			r.Components[0].Components[0].Options = getPointOptions(selectPoints)
			handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: r,
			})
		}

		handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
			if ic.Member.UserID == user.ID {
				if ic.Data.CustomID == "points-select" {
					if len(ic.Data.Values) == 0 {
						return
					}
					val := ic.Data.Values[0]
					newPoints, err := strconv.Atoi(val)
					if err != nil {
						return
					}
					selectPoints = newPoints
				}
				if ic.Data.CustomID == "1" || ic.Data.CustomID == "2" || ic.Data.CustomID == "3" {
					n, err := strconv.Atoi(ic.Data.CustomID)
					if err != nil {
						return
					}
					database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
						if selectPoints > calcAvailPoints(u) {
							return u
						}
						if u.Attributes[n-1] >= MAX_POINTS {
							return u
						}
						u.Attributes[n-1] += selectPoints
						userDb = u
						return u
					})
				}
				updateMessage(ic)
			}
		}, 100)
		return nil
	}

	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: response,
	}

}
