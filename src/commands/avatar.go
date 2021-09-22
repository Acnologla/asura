package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"avatar", "ava"},
		Run:       runAvatar,
		Available: true,
		Cooldown:  2,
		Options: handler.GetOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.USER,
			Name:        "usuario",
			Description: "usuario para ver o avatar",
			Required:    false,
		}),
		Usage: "j!avatar <usuario>",
		Help:  "Veja o avatar de alguem",
	})
}

func runAvatar(session disgord.Session, msg *disgord.Message, args []*disgord.ApplicationCommandDataOption) (*disgord.Message, func(disgord.Snowflake)) {
	user := utils.GetUser(msg, args, session)
	avatar, _ := user.AvatarURL(512, true)
	return handler.CreateMessage(&disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Color:       65535,
			Description: fmt.Sprintf("**%s**\n[**Download**](%s)", user.Username, avatar),
			Image: &disgord.EmbedImage{
				URL: avatar,
			},
		},
	}), nil
}
