package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"math/rand"
	"fmt"
	"github.com/andersfylling/disgord"
)

var galoColors = []int{65535,16777214,8421504,16711680}

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"galo", "galolevel", "meugalo"},
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
	
	galo, _ := rinha.GetGaloDB(user.ID)
	level := rinha.CalcLevel(galo.Xp)
	nextLevelXP := rinha.CalcXP(level+1)
	curLevelXP := rinha.CalcXP(level)
	if galo.Type == 0 {
		galoType :=  rand.Intn(len(rinha.Classes)-1)+1
		galo.Type = galoType
		rinha.SaveGaloDB(user.ID,galo)
	}
	
	
	var fields []*disgord.EmbedField
	skills := rinha.GetEquipedSkills(&galo)
	for i := 0; i < len(skills); i++ {
		skill := rinha.Skills[galo.Type-1][skills[i]]
		fields = append(fields, &disgord.EmbedField{
			Name:   skill.Name,
			Value:  rinha.SkillToString(skill),
			Inline: false,
		})
	}

	msg.Reply(context.Background(), session, disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Title: "Galo do " + user.Username,
			Color: galoColors[galo.Type-1],
			Thumbnail: &disgord.EmbedThumbnail{
				URL: "https://blogs.uai.com.br/cantodogalo/wp-content/uploads/sites/32/2017/09/galo-imagem.jpg",
			},
			Footer: &disgord.EmbedFooter{
				Text: "Use j!skills para ver os skills e equipa-las",
			},
			Description: fmt.Sprintf("Level: **%d\n**XP: **%d/%d**\nTipo: **%s**\nHabilidades:", level, galo.Xp - curLevelXP, nextLevelXP - curLevelXP,rinha.Classes[galo.Type].Name),
			Fields:      fields,
		},
	})
}
