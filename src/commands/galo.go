package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
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
		Category:  1,
	})
}

func runGalo(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg, args, session)
	galo, _ := rinha.GetGaloDB(user.ID)
	level := rinha.CalcLevel(galo.Xp)
	nextLevelXP := rinha.CalcXP(level + 1)
	curLevelXP := rinha.CalcXP(level)
	if galo.Type == 0 {
		galoType := rinha.GetRandByType(rinha.Common)
		galo.Type = galoType
		rinha.SaveGaloDB(user.ID, galo)
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
	name := "Galo do " + user.Username
	if galo.Name != ""{
		name = galo.Name
	}
	msg.Reply(context.Background(), session, disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Title: name,
			Color: rinha.Classes[galo.Type].Rarity.Color(),
			Thumbnail: &disgord.EmbedThumbnail{
				URL: rinha.Sprites[0][galo.Type-1],
			},
			Footer: &disgord.EmbedFooter{
				Text: "Use j!skills para ver os skills e equipa-las",
			},
			Description: fmt.Sprintf("Level: **%d** (**%d/%d**)\nTipo: **%s** (**%s**)\n **%d** Vitorias | **%d** Derrotas\nHabilidades:", level, galo.Xp-curLevelXP, nextLevelXP-curLevelXP, rinha.Classes[galo.Type].Name, rinha.Classes[galo.Type].Rarity.String(), galo.Win, galo.Lose),
			Fields:      fields,
		},
	})
}
