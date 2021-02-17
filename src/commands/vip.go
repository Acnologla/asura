package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
	"time"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"vip", "setvip"},
		Run:       runVip,
		Available: false,
		Cooldown:  1,
		Usage:     "j!vip @user",
		Help:      "vip",
	})
}

func runVip(session disgord.Session, msg *disgord.Message, args []string) {
	if msg.Author.ID != utils.StringToID("365948625676795904") {
		return
	}
	user := utils.GetUser(msg, args, session)
	if user.ID == msg.Author.ID {
		return
	}
	if len(args) == 0 {
		return
	}
	months, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}
	var vipTime uint64 = 0
	rinha.UpdateGaloDB(user.ID, func(galo rinha.Galo) (rinha.Galo, error) {
		vipTime = uint64(months * 30 * 24 * 60 * 60)
		galo.Vip = uint64(time.Now().Unix()) + vipTime
		if !rinha.IsIntInList(1, galo.Cosmetics) {
			galo.Cosmetics = append(galo.Cosmetics, 1)
		}
		return galo, nil
	})
	msg.Reply(context.Background(), session, fmt.Sprintf("**%s** agora é vip por **%d** meses", user.Username, vipTime/30/24/60/60))
}
