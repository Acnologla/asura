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
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title:       "Conquistas",
					Color:       65535,
					Description: description,
				},
			},
		},
	}
}
