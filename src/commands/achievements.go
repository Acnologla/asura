package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "conquistas",
		Description: "Veja suas conquistas",
		Run:         runAchievements,
		Cooldown:    5,
		Category:    handler.Profile,
	})
}

func runAchievements(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items")
	description := ""
	for _, item := range user.Items {
		if item.Type == entities.AchievementType {
			ach := rinha.Achievements[item.ItemID]
			description += fmt.Sprintf("**%s** - %s\n", ach.Name, ach.Description)
		}
	}

	description2 := ""

	for i, ach := range rinha.Achievements {
		if !rinha.HasAchievement(&user, i) {
			description2 += fmt.Sprintf("**%s** - %s\n", ach.Name, ach.Description)
		}
	}

	r := &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title:       "Conquistas",
					Color:       65535,
					Description: description,
				},
			},
			Components: []*disgord.MessageComponent{
				{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{

						{
							Type:  disgord.MessageComponentButton,
							Style: disgord.Success,

							Label:    "Obtidas",
							CustomID: "obtained",
						},
						{
							Type:  disgord.MessageComponentButton,
							Style: disgord.Danger,

							Label:    "Não obtidas",
							CustomID: "notObtained",
						},
					},
				},
			},
		},
	}

	err := handler.Client.SendInteractionResponse(ctx, itc, r)
	if err != nil {
		return nil
	}
	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		if ic.Member.UserID == user.ID {
			customID := ic.Data.CustomID
			if customID == "obtained" {
				r.Data.Embeds[0].Description = description
			} else {
				r.Data.Embeds[0].Description = description2
			}
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: r.Data,
			})

		}
	}, 200)
	return nil
}
