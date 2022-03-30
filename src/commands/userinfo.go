package commands

import (
	"asura/src/firebase"
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "userinfo",
		Description: translation.T("UserinfoHelp", "pt"),
		Run:         runUserinfo,
		Cooldown:    3,
		Cache:       60,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user info",
			Required:    false,
		}, &disgord.ApplicationCommandOption{
			Name:        "id",
			Type:        disgord.OptionTypeString,
			Description: "user id",
			Required:    false,
		}),
	})
}

func runUserinfo(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	if len(itc.Data.Options) > 0 {
		user = utils.GetUser(itc, 0)
	}
	var userinfo firebase.User
	var private bool
	avatar, _ := user.AvatarURL(512, true)
	authorAvatar, _ := itc.Member.User.AvatarURL(512, true)
	id := strconv.FormatUint(uint64(user.ID), 10)
	firebase.Database.NewRef("users/"+id).Get(ctx, &userinfo)
	firebase.Database.NewRef("private/"+id).Get(ctx, &private)
	oldAvatars := ""
	oldUsernames := ""
	date := ((uint64(user.ID) >> 22) + 1420070400000) / 1000
	if len(userinfo.Usernames) > 0 && !private {
		if len(userinfo.Usernames) > 20 {
			oldUsernames = strings.Join(userinfo.Usernames[:20], "** | **")
		} else {
			oldUsernames = strings.Join(userinfo.Usernames, "** | **")
		}
	} else if private {
		oldUsernames = translation.T("Private", translation.GetLocale(itc))
	} else {
		oldUsernames = translation.T("NoUsernames", translation.GetLocale(itc))
	}
	if len(userinfo.Avatars) > 0 && !private {
		avats := userinfo.Avatars
		if len(avats) > 10 {
			avats = avats[:10]
		}
		for i, avatar := range avats {
			oldAvatars += fmt.Sprintf("[**Link**](%s)", avatar)
			if i != len(avats) {
				oldAvatars += "** | **"
			}
		}
	} else if private {
		oldAvatars = translation.T("Private", translation.GetLocale(itc))
	} else {
		oldAvatars = translation.T("NoAvatars", translation.GetLocale(itc))
	}
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Color: 65535,
					Title: fmt.Sprintf("%s(%s)", user.Username, user.ID),
					Thumbnail: &disgord.EmbedThumbnail{
						URL: avatar,
					},
					Description: translation.T("AccountCreated", translation.GetLocale(itc), int((uint64(time.Now().Unix())-date)/60/60/24)),
					Footer: &disgord.EmbedFooter{
						Text:    itc.Member.User.Username,
						IconURL: authorAvatar,
					},
					Fields: []*disgord.EmbedField{
						{
							Name:  translation.T("OldUsernames", translation.GetLocale(itc)),
							Value: oldUsernames,
						},
						{
							Name:  translation.T("OldAvatars", translation.GetLocale(itc)),
							Value: oldAvatars,
						},
					},
				},
			},
		},
	}
}
