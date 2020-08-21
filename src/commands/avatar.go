package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"avatar", "ava"},
		Run:       runAvatar,
		Available: true,
		Cooldown:  2,
		Usage:     "j!avatar <usuario>",
		Help:      "Veja o avatar de alguem",
	})
}

func runAvatar(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg, args,session)
	avatar, _ := user.AvatarURL(512, false)
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Color:       65535,
			Description: fmt.Sprintf("**%s**\n[**Download**](%s)", user.Username, avatar),
			Image: &disgord.EmbedImage{
				URL: avatar,
			},
		},
	})
}
