package commands

import (
	"asura/src/firebase"
	"asura/src/handler"
	"context"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "private",
		Description: translation.T("PrivateHelp", "pt"),
		Run:         runPrivate,
		Cooldown:    5,
	})
}

func runPrivate(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	var acc bool
	id := strconv.FormatUint(uint64(itc.Member.UserID), 10)
	ctx := context.Background()
	if err := firebase.Database.NewRef("private/"+id).Get(ctx, &acc); err != nil {
		return nil
	}
	newVal := !acc
	if err := firebase.Database.NewRef("private/"+id).Set(ctx, &(newVal)); err != nil {
		return nil
	}
	if acc {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("PrivateDisabled", translation.GetLocale(itc)),
			},
		}
	} else {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("PrivateEnabled", translation.GetLocale(itc)),
			},
		}
	}
}
