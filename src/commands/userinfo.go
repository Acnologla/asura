package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
	"strings"
	"time"
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
	filteredGuilds := ""
	count := 0
	date := ((uint64(user.ID) >> 22) + 1420070400000) / 1000
	cGuilds := session.GetConnectedGuilds()
	for i, guild := range cGuilds {
		guil, _ := handler.Client.Guild(guild).Get()
		if guil == nil {
			continue
		}
		var is bool
		if count >= 10 {
			break
		}
		for _, member := range guil.Members {
			if member.User != nil {
				if member.User.ID == user.ID {
					is = true
					break
				}
			}
		}
		if is {
			filteredGuilds += guil.Name
			count++
			if i != len(cGuilds) {
				filteredGuilds += "** | **"
			}

		}
	}
	if filteredGuilds == "" {
		filteredGuilds = "Nenhum servidor compartilhado"
	}
	if len(userinfo.Usernames) > 0 {
		if len(userinfo.Usernames) > 12 {
			oldUsernames = strings.Join(userinfo.Usernames[:12], "** | **")
		} else {
			oldUsernames = strings.Join(userinfo.Usernames, "** | **")
		}
	} else {
		oldUsernames = "Nenhum nome antigo registrado"
	}
	if len(userinfo.Avatars) > 0 && !private {
		avats := userinfo.Avatars
		if len(avats) > 12 {
			avats = avats[:12]
		}
		for i, avatar := range avats {
			oldAvatars += fmt.Sprintf("[**Link**](%s)", avatar)
			if i != len(avats) {
				oldAvatars += "** | **"
			}
		}
	} else if private {
		oldAvatars = "O historico desse usuario é privado, use j!private para deixar publico"
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
				&disgord.EmbedField{
					Name:  "Nomes Antigos",
					Value: oldUsernames,
				},
				&disgord.EmbedField{
					Name:  "Avatares Antigos",
					Value: oldAvatars,
				},
				&disgord.EmbedField{
					Name:  fmt.Sprintf("Servidores compartilhados (%d)", count),
					Value: filteredGuilds,
				},
				&disgord.EmbedField{
					Name:  "Mais informaçoes",
					Value: fmt.Sprintf("[**Clique aqui**](https://asura-site.glitch.me/user/%s)", id),
				},
			},
		}})
}
