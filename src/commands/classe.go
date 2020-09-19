package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"classe"},
		Run:       runClass,
		Available: true,
		Cooldown:  2,
		Usage:     "j!classe ou | j!classe escolher <nome>",
		Help:      "Escolher",
	})
}

func displayClasses(session disgord.Session, msg *disgord.Message, args []string) {
	var fields []*disgord.EmbedField

	for i := 1; i < len(utils.Classes); i++ {
		fields = append(fields, &disgord.EmbedField{
			Name:   utils.Classes[i].Name,
			Value:  utils.Classes[i].Desc,
			Inline: false,
		})
	}

	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Title:       "Classes",
			Color:       65535,
			Description: "Classes disponiveis para o level 5:",
			Footer: &disgord.EmbedFooter{
				Text: "Use 'j!classe escolher [nome]' para escolher uma classe!",
			},
			Fields: fields,
		},
	})
}

func runClass(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author

	galo, _ := utils.GetGaloDB(user.ID)
	level := utils.CalcLevel(galo.Xp)

	if level >= 0 {
		if len(args) == 2 && args[0] == "escolher" {
			if galo.Type != 0 {
				msg.Reply(context.Background(), session, "Você já escolheu uma classe!")
				return
			}
			var class *utils.Class
			class_id := 0
			for i := 1; i < len(utils.Classes); i++ {
				if utils.Classes[i].Name == args[1] {
					class = &utils.Classes[i]
					class_id = i
				}
			}
			if class == nil {
				msg.Reply(context.Background(), session, "Essa classe não existe ou não pode ser usada por voce!")
			} else {
				galo.Type = class_id
				utils.SaveGaloDB(user.ID, galo)
				msg.Reply(context.Background(), session, "Parabens agora voce pode adquirir skills da classe '"+class.Name+"'!")
			}
		} else {
			displayClasses(session, msg, args)
		}
	} else {
		msg.Reply(context.Background(), session, "Você tem que ser level 5 ou maior para escolher uma classe!")
	}

}
