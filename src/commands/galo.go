package commands

import (
	"asura/src/handler"
	"context"
	"io/ioutil"
	"github.com/andersfylling/disgord"
	"encoding/json"
	"asura/src/database"
	"strconv"
	"asura/src/utils"
	"fmt"
)

type skill struct{
	Name   string `json:"name"`
    Damage [2]int `json:"damage"`
}
var skills []skill
func init() {
	byteValue, _ := ioutil.ReadFile("./resources/atacks.json")
	json.Unmarshal([]byte(byteValue), &skills)
	handler.Register(handler.Command{
		Aliases:   []string{"galo","galolevel","meugalo"},
		Run:       runGalo,
		Available: true,
		Cooldown:  3,
		Usage:     "j!galo",
		Help:      "Informação sobre seu galo",
	})
}

func runGalo(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author
	if len(msg.Mentions) > 0 {
		user = msg.Mentions[0]
	}
	var galo database.Galo
	id := strconv.FormatUint(uint64(user.ID),10)
	database.Database.NewRef("galo/"+id).Get(context.Background(), &galo)
	level := utils.CalcLevel(galo.Xp) +1
	nextLevel := utils.CalcXP(level)
	var nextSkill string
	if level <= len(skills)-1{
		nextSkill = skills[level].Name
	} 
	var fields []*disgord.EmbedField
	for i:= 0; i < level;i++{
		if len(skills)-1 >= i {
			fields = append(fields,&disgord.EmbedField{
				Name: skills[i].Name,
				Value: fmt.Sprintf("Dano: %d - %d",skills[i].Damage[0], skills[i].Damage[1]-1),
				Inline: true,
			})
		}
	}
	msg.Reply(context.Background(),session,disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Title: "Galo do " + user.Username,
			Color: 65535,
			Thumbnail: &disgord.EmbedThumbnail{
				URL: "https://blogs.uai.com.br/cantodogalo/wp-content/uploads/sites/32/2017/09/galo-imagem.jpg",
			},
			Description: fmt.Sprintf("Level **%d**\nXP: **%d/%d**\nProxima habilidade: **%s**\n\nHabilidades Atuais:",level,galo.Xp,nextLevel,nextSkill),
			Fields: fields,
		},
	})	
}