package utils

import (
	"asura/src/handler"
	"context"
	"strconv"

	"github.com/andersfylling/disgord"
)

func GenerateOptions(options ...*disgord.ApplicationCommandOption) []*disgord.ApplicationCommandOption {
	return options
}

func GetOptionsUser(options []*disgord.ApplicationCommandDataOption, itc *disgord.InteractionCreate, i int) *disgord.User {
	opt := options[i]
	idStr := opt.Value.(string)
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if itc.Data.Resolved == nil || opt.Type == disgord.OptionTypeString {
		u, _ := handler.Client.User(disgord.Snowflake(id)).Get()
		return u
	}
	return itc.Data.Resolved.Users[disgord.Snowflake(id)]
}

func GetUser(itc *disgord.InteractionCreate, i int) *disgord.User {
	opt := itc.Data.Options[i]
	idStr := opt.Value.(string)
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if itc.Data.Resolved == nil || opt.Type == disgord.OptionTypeString {
		u, err := handler.Client.User(disgord.Snowflake(id)).Get()
		if err != nil {
			return itc.Member.User
		}
		return u
	}
	return itc.Data.Resolved.Users[disgord.Snowflake(id)]
}

func GetUserOrAuthor(itc *disgord.InteractionCreate, i int) *disgord.User {
	if i < len(itc.Data.Options) {
		return GetUser(itc, i)
	}
	return itc.Member.User
}

func createConfirmMessage(title string) *disgord.CreateInteractionResponseData {
	return &disgord.CreateInteractionResponseData{
		Embeds: []*disgord.Embed{
			{
				Title: title,
				Color: 65535,
			},
		},
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Aceitar",
						Style:    disgord.Success,
						CustomID: "yes",
					},
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Recusar",
						Style:    disgord.Danger,
						CustomID: "no",
					},
				},
			},
		},
	}
}

func confirmHandler(userID, id disgord.Snowflake, callback func()) {
	done := false
	handler.RegisterHandler(id, func(interaction *disgord.InteractionCreate) {
		if userID == interaction.Member.User.ID && !done {
			done = true
			go handler.Client.Channel(interaction.ChannelID).Message(interaction.Message.ID).Delete()

			if interaction.Data.CustomID == "yes" {
				callback()
			}
			handler.DeleteHandler(id)
		}
	}, 120)
}

func ConfirmMessage(ctx context.Context, title string, itc *disgord.InteractionCreate, id disgord.Snowflake, callback func()) {
	msg := createConfirmMessage(title)
	message, err := handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
		Components: msg.Components,
		Embeds:     msg.Embeds,
	})
	if err != nil {
		return
	}
	confirmHandler(id, message.ID, callback)
}

func Confirm(ctx context.Context, title string, itc *disgord.InteractionCreate, id disgord.Snowflake, callback func()) {
	msgID, _ := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: createConfirmMessage(title),
	})
	confirmHandler(id, msgID, callback)
}
