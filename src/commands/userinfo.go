package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func init() {

	handler.Register(handler.Command{
		Aliases:   []string{"userinfo", "usuario", "uinfo"},
		Run:       runUserinfo,
		Available: true,
		Cooldown:  4,
		Usage:     "j!userinfo <usuario.",
		Help:      "Veja as informaçoes de um usuario",
	})
}

func runUserinfo(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg, args, session)
	ctx := context.Background()
	var userinfo database.User
	var private bool
	avatar, _ := user.AvatarURL(512, true)
	authorAvatar, _ := msg.Author.AvatarURL(512, true)
	id := strconv.FormatUint(uint64(user.ID), 10)
	database.Database.NewRef("users/"+id).Get(ctx, &userinfo)
	database.Database.NewRef("private/"+id).Get(ctx, &private)
	oldAvatars := ""
	oldUsernames := ""
	date := ((uint64(user.ID) >> 22) + 1420070400000) / 1000
	if len(userinfo.Usernames) > 0 && !private {
		if len(userinfo.Usernames) > 10 {
			oldUsernames = strings.Join(userinfo.Usernames[:10], "** | **")
		} else {
			oldUsernames = strings.Join(userinfo.Usernames, "** | **")
		}
	} else if private {
		oldUsernames = "O histórico desse usuario é privado, use j!private para deixar publico"
	} else {
		oldUsernames = "Nenhum nome antigo registrado"
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
		oldAvatars = "O histórico desse usuario é privado, use j!private para deixar público!"
	} else {
		oldAvatars = "Nenhum avatar antigo registrado"
	}
	msg.Reply(ctx, session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Color: 65535,
			Title: fmt.Sprintf("%s(%s)", user.Username, user.ID),
			Thumbnail: &disgord.EmbedThumbnail{
				URL: avatar,
			},
			Description: fmt.Sprintf("Conta criada a **%d** Dias", int((uint64(time.Now().Unix())-date)/60/60/24)),
			Footer: &disgord.EmbedFooter{
				Text:    msg.Author.Username,
				IconURL: authorAvatar,
			},
			Fields: []*disgord.EmbedField{
				{
					Name:  "Nomes Antigos",
					Value: oldUsernames,
				},
				{
					Name:  "Avatares Antigos",
					Value: oldAvatars,
				},
			},
		}})
}
