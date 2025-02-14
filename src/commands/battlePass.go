package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "battlepass",
		Description: translation.T("BattlePassHelp", "pt"),
		Run:         runBattlePass,
		Cooldown:    5,
		Category:    handler.Rinha,
		Aliases:     []string{"passe"},
		AliasesMsg:  []string{"bp"},
	})
}

func runBattlePass(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	rinhaUser := database.User.GetUser(ctx, itc.Member.User.ID)
	bpLevel := rinha.CalcBPLevel(rinhaUser.BattlePass)
	nextLevels := ""
	for i := bpLevel + 1; i < bpLevel+10; i++ {
		if i >= len(rinha.BattlePass) {
			break
		}
		nextLevels += fmt.Sprintf("[ **%d** ] - %s\n", i, rinha.BPLevelToString(i))
	}

	extraMsg := ""
	isVip := rinha.IsVip(&rinhaUser)

	if !isVip && bpLevel >= len(rinha.BattlePass)/2 {
		extraMsg = "\n\n**Você pode comprar o passe de batalha vip para ter acesso a todos os niveis**"
	}
	curLevelXP := rinha.CalcBPXP(bpLevel)
	return entities.CreateMsg().Embed(&disgord.Embed{
		Title:       "Passe de batalha",
		Color:       65535,
		Description: fmt.Sprintf("Seu nivel atual é **%d** (%d/%d)%s\n\nRecompensas para os proximos 10 niveis:\n\n%s", bpLevel, rinhaUser.BattlePass-curLevelXP, rinha.CalcBPXP(bpLevel+1)-curLevelXP, extraMsg, nextLevels),
	}).Res()
}
