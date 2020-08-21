package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"strings"
	"fmt"
	"time"
	"asura/src/database"
	"strconv"
	"context"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"userinfo", "usuario","uinfo"},
		Run:       runUserinfo,
		Available: true,
		Cooldown:  4,
		Usage:     "j!userinfo <usuario.",
		Help:      "Veja as informaçoes de um usuario",
	})
}

func runUserinfo(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg,args,session)
	var userinfo database.User
	var private bool
	avatar,_ := user.AvatarURL(512,false)
	authorAvatar,_ := msg.Author.AvatarURL(512,false)
	id := strconv.FormatUint(uint64(user.ID),10)
	database.Database.NewRef("users/"+id).Get(context.Background(), &userinfo)
	database.Database.NewRef("private/"+id).Get(context.Background(), &private)
	guilds := session.GetConnectedGuilds()
	oldAvatars := ""
	oldUsernames := ""
	filteredGuilds := ""
	count := 0
	date := ((uint64(user.ID) >> 22) +1420070400000) / 1000
	for i,guild := range guilds{
		_, is :=session.GetMember(context.Background(),guild, user.ID)
		if count >= 12{
			break
		}
		if is == nil {
			name,err := session.GetGuild(context.Background(),guild)
			if err == nil {
				filteredGuilds += name.Name
				count ++
				if i != len(guilds){
					filteredGuilds+= "** | **"
				}
			}
		}
	}
	if len(userinfo.Usernames) >0 {
		if len(userinfo.Usernames) > 12{
			oldUsernames = strings.Join(userinfo.Usernames[:12], "** | **")
		}else{
			oldUsernames = strings.Join(userinfo.Usernames, "** | **")
		}
	}else {
		oldUsernames = "Nenhum nome antigo registrado"
	}
	if len(userinfo.Avatars) >0  && !private {
		avats := userinfo.Avatars
		if len(avats) > 12{
			avats = avats[:12]
		}
		for i, avatar := range avats{
			oldAvatars+= fmt.Sprintf("[**Link**](%s)",avatar)
			if i != len(avats){
				oldAvatars+= "** | **"
			}
		}
	}else if private{
		oldAvatars = "O historico desse usuario é privado, use j!private para deixar publico"
	}else{
		oldAvatars = "Nenhum avatar antigo registrado"
	}
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Color:       65535,
			Title:      fmt.Sprintf("%s(%s)",user.Username,user.ID),
			Thumbnail: &disgord.EmbedThumbnail{
				URL: avatar,
			},
			Description: fmt.Sprintf("Conta criada a **%d** Dias", int(( uint64(time.Now().Unix()) - date) / 60 / 60 / 24)),
			Footer: &disgord.EmbedFooter{
				Text: msg.Author.Username,
				IconURL: authorAvatar,
			},
			Fields: []*disgord.EmbedField {
				&disgord.EmbedField{
				Name: "Nomes Antigos",
				Value:	oldUsernames,
			},
			&disgord.EmbedField{
				Name: "Avatares Antigos",
				Value:	oldAvatars,
			},
			&disgord.EmbedField{
				Name: fmt.Sprintf("Servidores compartilhados (%d)",count),
				Value:	filteredGuilds,
			},
			&disgord.EmbedField{
				Name: "Mais informaçoes",
				Value:	fmt.Sprintf("[**Clique aqui**](https://asura-site.glitch.me/user/%s)",id),
			},
		},
		}})
}