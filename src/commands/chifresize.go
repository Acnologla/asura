package commands

import (
	"asura/src/handler"
	"context"
	"asura/src/utils"
	"github.com/andersfylling/disgord"
	"strconv"
	"fmt"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"chifresize","cs","cornosize"},
		Run:       runChifresize,
		Available: true,
		Cooldown:  1,
		Usage:     "j!chifresize <usuario>",
		Help:      "Veja o tamanho do seu chifre",
	})
}

func runChifresize(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg,args)
	idString :=  strconv.FormatUint(uint64(user.ID), 10)
	result, _ := strconv.Atoi(string(idString[3:4]))
	random,_ := strconv.Atoi(string(idString[5]))
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: msg.Author.Mention(),
		Embed: &disgord.Embed{
			Description: fmt.Sprintf("%s Seu chifre tem **%d** cms de altura e **%d** cms de circunferencia",user.Mention(),result * 3,result + random),
			Color:       65535,
			Title:       ":ox: Tamanho do chifre do " + user.Username,
		},})
}
