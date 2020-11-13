package commands

import (
	"asura/src/handler"
	"context"
	"fmt"
	"asura/src/utils/rinha"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"lootbox","lb"},
		Run:       runLootbox,
		Available: true,
		Cooldown:  3,
		Usage:     "j!lootbox",
		Help:      "Abra lootboxs",
		Category:  1,
	})
}

func runLootbox(session disgord.Session, msg *disgord.Message, args []string) {
	galo,_ := rinha.GetGaloDB(msg.Author.ID)
	normal := func(){
		msg.Reply(context.Background(),session,&disgord.Embed{
			Title: "Lootbox",
			Color: 65535,
			Description: fmt.Sprintf("Lootbox: **%d**\nMoney: **%d**\n\nUse `j!lootbox buy` para comprar lootbox\nUse `j!lootbox open` para abrir lootbox",galo.Lootbox, galo.Money),
		})
	}
	if len(args) == 0 {
		normal()
		return
	}
	if args[0] == "open" || args[0] == "abrir"{
		if galo.Lootbox == 0 {
			msg.Reply(context.Background(),session,"Voce tem 0 lootboxs, use `j!lootbox buy` para comprar uma")
			return
		}
		//TODO open lootbox
	}else if args[0] == "buy" || args[0] == "comprar" {
		//TODO buy lootbox
	}else{
		normal()
	}
}
