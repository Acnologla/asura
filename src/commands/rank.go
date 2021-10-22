package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"sort"

	"firebase.google.com/go/db"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"rank", "top", "tops"},
		Run:       runRank,
		Available: true,
		Cooldown:  20,
		Usage:     "j!rank <rank>",
		Help:      "Veja os ranks",
		Category:  1,
	})
}

func childToQuery(str string) string {
	return map[string]string{
		"players":      "dungeonreset",
		"money":        "money",
		"level":        "xp",
		"clan":         "xp",
		"clandinheiro": "money",
		"vitorias":     "win",
		"derrotas":     "lose",
	}[str]
}

func sortDungeon(values []db.QueryNode) []db.QueryNode {
	sort.Slice(values, func(i, j int) bool {
		var val1 rinha.Galo
		var val2 rinha.Galo
		values[i].Unmarshal(&val1)
		values[j].Unmarshal(&val2)
		return val1.Dungeon*(val1.DungeonReset+1) < val2.Dungeon*(val2.DungeonReset+1)
	})
	return values
}

func top(topType string, session disgord.Session) (text string) {
	query := "galo"
	if topType == "clan" || topType == "clandinheiro" {
		query = "clan"
	}
	child := childToQuery(topType)

	q := database.Database.NewRef(query).OrderByChild(child)
	if topType == "players" {
		q = q.StartAt(1)
	}
	q = q.LimitToLast(15)
	result, err := q.GetOrdered(context.Background())
	if topType == "players" {
		result = sortDungeon(result)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := len(result) - 1; 0 <= i; i-- {
		if query == "clan" {
			var clan rinha.Clan
			if err := result[i].Unmarshal(&clan); err != nil {
				continue
			}
			name := result[i].Key()
			maxMembers := rinha.GetMaxMembers(clan)
			if topType == "clandinheiro" {
				money := clan.Money
				text += fmt.Sprintf("[%d] - %s (%d/%d)\nDinheiro: %d\n", len(result)-i, name, len(clan.Members), maxMembers, money)
			} else {
				level := rinha.ClanXpToLevel(clan.Xp)
				text += fmt.Sprintf("[%d] - %s (%d/%d)\nLevel: %d (%d XP)\n", len(result)-i, name, len(clan.Members), maxMembers, level, clan.Xp)
			}

		} else {
			var gal rinha.Galo
			if err := result[i].Unmarshal(&gal); err != nil {
				continue
			}
			converted := utils.StringToID(result[i].Key())
			user, err := session.User(converted).Get()
			var username string
			name := child
			if err != nil {
				username = "Anonimo"
			} else {
				username = user.Username + "#" + user.Discriminator.String()
			}
			value := 0
			if child == "money" {
				value = gal.Money
			}
			if child == "dungeonreset" {
				name = "Andar da dungeon"
				value = 18*(gal.DungeonReset) + gal.Dungeon
			}
			if child == "xp" {
				name = "Level"
				value = rinha.CalcLevel(gal.Xp)
			}
			if child == "win" {
				name = "Vitorias"
				value = gal.Win
			}
			if child == "lose" {
				name = "Derrotas"
				value = gal.Lose
			}
			text += fmt.Sprintf("[%d] - %s\n%s: %d\n", len(result)-i, username, name, value)
		}
	}
	return
}

func runRank(session disgord.Session, msg *disgord.Message, args []string) {
	if len(args) == 0 {
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Content: msg.Author.Mention(),
			Embed: &disgord.Embed{
				Description: "Use `j!rank players` para ver os melhores jogadores\nUse `j!rank money` para ver os jogadores com mais dinheiro\nUse `j!rank clan` para ver os melhores clan\nUse `j!rank level` para ver os galos com o maior nivel\nUse `j!rank vitorias` para ver os jogadores com mais vitorias\nUse `j!rank derrotas` para ver os jogadores com mais derrotas\nUse `j!rank clandinheiro` para ver os clans com o maior dinheiro",
				Title:       "Ranks",
				Color:       65535,
			},
		})
	} else {
		text := args[0]
		if childToQuery(args[0]) == "" {
			msg.Reply(context.Background(), session, msg.Author.Mention()+", rank invalido\nuse `j!ranks` para ver os ranks disponiveis")
			return
		}
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Content: msg.Author.Mention(),
			Embed: &disgord.Embed{
				Description: top(text, session),
				Color:       65535,
				Title:       "Rank " + text,
			},
		})
	}

}
