package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strings"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"clan", "guilda"},
		Run:       runClan,
		Available: true,
		Cooldown:  3,
		Usage:     "j!clan",
		Help:      "Informação sobre o cla",
		Category:  1,
	})
}

func runClan(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Clan == "" {
		if len(args) > 1 {
			if args[0] == "create" {
				text := rinha.Format(strings.Join(args[1:], " "))
				if len(text) >= 25 {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", O nome do clan pode ter no maximo 24 caracteres")
					return
				}
				clan := rinha.GetClan(text)
				if clan.CreatedAt != 0 {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Ja tem um clan com esse nome")
					return
				}
				err := rinha.ChangeMoney(msg.Author.ID, -1000, 1000)
				if err != nil {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce precisa ter 1000 de dinheiro para criar um clan")
					return
				}
				rinha.CreateClan(text, msg.Author.ID)
				rinha.UpdateGaloDB(msg.Author.ID, map[string]interface{}{
					"clan": text,
				})
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Clan **"+text+"** criado com sucesso\nUse **j!clan invite <usuario>** para convidar membros")
				return
			}
		}
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title:       "Clan",
			Color:       65535,
			Description: "Voce não esta em nenhum clan\nUse j!clan create <nome> para criar um (**1000** de dinheiro)",
		})
	} else {
		clan := rinha.GetClan(galo.Clan)
		if len(args) > 0 {
			if strings.ToLower(args[0]) == "leave" {
				role := rinha.GetMember(clan, msg.Author.ID)
				members := rinha.RemoveMember(clan, msg.Author.ID)
				if len(members) == 0 {
					rinha.DeleteClan(galo.Clan)
					msg.Reply(context.Background(), session, msg.Author.Mention()+", O clan **"+galo.Clan+"** foi deletado com sucesso")
					return
				}
				if role.Role == rinha.Owner {
					members[0].Role = rinha.Owner
				}
				rinha.UpdateClan(galo.Clan, map[string]interface{}{
					"members": members,
				})
				rinha.UpdateGaloDB(msg.Author.ID, map[string]interface{}{
					"clan": "",
				})
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce saiu do clan **"+galo.Clan+"** com sucesso")
				return
			}
		}
		if len(args) > 1 {
			role := rinha.GetMember(clan, msg.Author.ID)
			if strings.ToLower(args[0]) == "invite" {
				if len(msg.Mentions) == 0 {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Use j!clan invite @user")
					return
				}
				user := msg.Mentions[0]
				if user.ID == msg.Author.ID || user.Bot {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Usuario invalido")
					return
				}
				if role.Role == rinha.Member {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Apenas donos e administradores podem convidar pessoas ao clan")
					return
				}
				confirmMsg, confirmErr := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
					Content: msg.Mentions[0].Mention(),
					Embed: &disgord.Embed{
						Color:       65535,
						Description: fmt.Sprintf("**%s** clique na reação abaixo para aceitar o convite", msg.Mentions[0].Username),
					},
				})
				if confirmErr == nil {
					utils.Confirm(confirmMsg, msg.Mentions[0].ID, func() {
						clan = rinha.GetClan(galo.Clan)
						level := rinha.ClanXpToLevel(clan.Xp)
						maxMembers := 15
						if level >= 3{
							maxMembers = 21
						}
						if rinha.IsInClan(clan, user.ID) {
							msg.Reply(context.Background(), session, msg.Author.Mention()+", Este usuario ja esta no clan")
							return
						}
						if len(clan.Members) >=  maxMembers {
							msg.Reply(context.Background(), session, msg.Author.Mention()+", Este clan ja esta cheio\nRemova algum usuario usando j!clan remove @user")
							return
						}
						invited, _ := rinha.GetGaloDB(user.ID)
						if invited.Clan != "" {
							msg.Reply(context.Background(), session, msg.Author.Mention()+", Este usuario esta em outro clan ("+invited.Clan+")")
							return
						}
						clan.Members = append(clan.Members, rinha.ClanMember{
							ID:   uint64(user.ID),
							Role: rinha.Member,
						})
						rinha.UpdateClan(galo.Clan, map[string]interface{}{
							"members": clan.Members,
						})
						rinha.UpdateGaloDB(user.ID, map[string]interface{}{
							"clan": galo.Clan,
						})
						msg.Reply(context.Background(), session, user.Mention()+", Voce entrou para o clan **"+galo.Clan+"** com sucesso")
					})
				}
				return
			}
			if strings.ToLower(args[0]) == "remove" {
				user := utils.GetUser(msg, args[1:], session)
				if user.ID == msg.Author.ID || user.Bot {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Usuario invalido")
					return
				}
				if role.Role == rinha.Member {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Apenas donos e administradores podem remover pessoas do clan")
					return
				}
				if !rinha.IsInClan(clan, user.ID) {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Este usuario nao esta no clan")
					return
				}
				member := rinha.GetMember(clan, user.ID)
				if member.Role >= role.Role {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce não tem permissoes para remover esse usuario")
					return
				}
				members := rinha.RemoveMember(clan, user.ID)
				rinha.UpdateClan(galo.Clan, map[string]interface{}{
					"members": members,
				})
				rinha.UpdateGaloDB(user.ID, map[string]interface{}{
					"clan": "",
				})
				msg.Reply(context.Background(), session, msg.Author.Mention()+", O usuario **"+user.Username+"** foi removido com sucesso")
				return
			}
		}
		level := rinha.ClanXpToLevel(clan.Xp)
		memberMsg := ""
		for _, member := range clan.Members {
			user, err := handler.Client.User(disgord.Snowflake(member.ID)).Get()
			if err == nil {
				memberMsg += fmt.Sprintf("[**%s**] %s\n", member.Role.ToString(), user.Username)
			}
		}
		benefits := rinha.GetBenefits(clan.Xp)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title: galo.Clan,
			Color: 65535,
			Footer: &disgord.EmbedFooter{
				Text: "Use j!clan invite <usuario> para convidar alguem | Use j!clan remove <usuario> para remover alguem",
			},
			Description: fmt.Sprintf("Level: **%d** (%d/%d)\nVantagens do clan:\n %s\nMembros (%d/15):\n %s", level, clan.Xp, rinha.ClanLevelToXp(level), benefits, len(clan.Members), memberMsg),
		})
	}
}
