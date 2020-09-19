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
	
	galo, _ := utils.GetGaloDB(user.ID)
	level := utils.CalcLevel(galo.Xp)
	nextLevelXP := utils.CalcXP(level+1)
	curLevelXP := utils.CalcXP(level)

	var fields []*disgord.EmbedField

	galo.Skills = []int{0,1,5}

	for i := 0; i < len(galo.Equipped); i++ {
		skill := utils.Skills[galo.Equipped[i]]
		fields = append(fields, &disgord.EmbedField{
			Name:   skill.Name,
			Value:  fmt.Sprintf("Dano: %d - %d", skill.Damage[0], skill.Damage[1]-1),
			Inline: true,
		})
	}

	msg.Reply(context.Background(), session, disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Title: "Galo do " + user.Username,
			Color: 65535,
			Thumbnail: &disgord.EmbedThumbnail{
				URL: "https://blogs.uai.com.br/cantodogalo/wp-content/uploads/sites/32/2017/09/galo-imagem.jpg",
			},
			Footer: &disgord.EmbedFooter{
				Text: "Use j!skills para ver os skills e equipa-las",
			},
			Description: fmt.Sprintf("Level **%d**\nXP: **%d/%d**", level, galo.Xp - curLevelXP, nextLevelXP - curLevelXP),
			Fields:      fields,
		},
	})
}
