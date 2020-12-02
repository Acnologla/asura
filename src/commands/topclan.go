package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"topclan", "melhoresclans"},
		Run:       runTopClans,
		Available: true,
		Cooldown:  15,
		Usage:     "j!topclans",
		Help:      "Veja os melhores clans",
		Category:  1,
	})
}

func runTopClans(session disgord.Session, msg *disgord.Message, args []string) {
	q := database.Database.NewRef("clan").OrderByChild("xp")
	result, err := q.GetOrdered(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var text string
	for i := len(result) - 1; 0 <= i; i-- {
		if i > len(result)-11 {
			var clan rinha.Clan
			if err := result[i].Unmarshal(&clan); err != nil {
				continue
			}
			name := result[i].Key()
			text += fmt.Sprintf("[%d] - %s (%d/15)\nLevel: %d (%d XP)\n", len(result)-i, name, len(clan.Members) ,rinha.ClanXpToLevel(clan.Xp), clan.Xp)
		} 
	}
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Description: text,
			Color: 65535,
			Title: "Topclans",
		},
	})
}
