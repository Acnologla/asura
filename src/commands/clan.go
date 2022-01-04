package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"clan", "guilda"},
		Run:       runClan,
		Available: true,
		Cooldown:  10,
		Usage:     "j!clan",
		Help:      "Informação sobre o cla",
		Category:  2,
	})
}

func runClan(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Clan == "" {
		if len(args) > 1 {
			if args[0] == "create" {
				text := rinha.Format(strings.Join(args[1:], " "))
				if text == "_test" {
					return
				}
				if len(text) >= 25 {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", O nome do clan pode ter no máximo 24 caracteres.")
					return
				}
				if len(text) <= 4 {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", O nome do clan pode ter no mínimo 5 caracteres.")
					return
				}
				clan := rinha.GetClan(text)
				if clan.CreatedAt != 0 {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Infelizmente já possui um clan com esse nome!")
					return
				}
				err := rinha.ChangeMoney(msg.Author.ID, -1000, 1000)
				if err != nil {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Você precisa ter 1000 de dinheiro para criar um clan!")
					return
				}
				rinha.CreateClan(text, msg.Author.ID)
				rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
					gal.Clan = text
					return gal, nil
				})
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Clan **"+text+"** criado com sucesso!\nUse **j!clan invite <usuario>** para convidar membros")
				return
			}
		}
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title:       "Clan",
			Color:       65535,
			Description: "Você não esta em nenhum clan\nUse j!clan create <nome> para criar um (**1000** de dinheiro)",
		})
	} else {
		clan := rinha.GetClan(galo.Clan)
		if len(args) > 0 {
			if strings.ToLower(args[0]) == "leave" {
				role := rinha.GetMember(clan, msg.Author.ID)
				members := rinha.RemoveMember(clan, msg.Author.ID)
				if len(members) == 0 {
					rinha.DeleteClan(galo.Clan)
					rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
						gal.Clan = ""
						return gal, nil
					})
					msg.Reply(context.Background(), session, msg.Author.Mention()+", O clan **"+galo.Clan+"** foi deletado com sucesso")
					return
				}
				if role.Role == rinha.Owner {
					members[0].Role = rinha.Owner
				}
				rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
					clan.Members = members
					return clan, nil
				})
				rinha.UpdateGaloDB(msg.Author.ID, func(gal rinha.Galo) (rinha.Galo, error) {
					gal.Clan = ""
					return gal, nil
				})
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce saiu do clan **"+galo.Clan+"** com sucesso")
				return
			}
			if strings.ToLower(args[0]) == "background" {
				role := rinha.GetMember(clan, msg.Author.ID)
				if role.Role < rinha.Admin {
					return
				}
				if !utils.CheckImage(args[1]) {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Imagem invalida!")
					return
				}
				rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
					if clan.Money < 10000 {
						msg.Reply(context.Background(), session, msg.Author.Mention()+", O clan precisa ter 10.000 de dinheiro para trocar o background")
						return clan, nil
					}
					clan.Money -= 10000
					clan.Background = args[1]
					msg.Reply(context.Background(), session, "O background do clan foi alterado com sucesso")
					return clan, nil
				})
				return
			}
			if strings.ToLower(args[0]) == "banco" {
				msg.Reply(context.Background(), session, &disgord.Embed{
					Title:       "Banco do clan",
					Color:       65535,
					Description: fmt.Sprintf("Dinheiro: **%d**\n\nUpgrades:\n[**Membros**] - Membro adicional para clan (%d)\n[**Artilharia**] - Artilharia (%d)\n[**Missao**] - Bonus para missao do clan (%d)\n\nUse j!clan depositar <dinheiro> para depositar dinheiro\nUse j!clan upgrade <upgrade> para dar upgrade em algo\nUse j!clan background <background> para setar um background para seu clan (custa 10000)", clan.Money, rinha.CalcClanUpgrade(clan.Upgrades.Members, 1), rinha.CalcClanUpgrade(clan.Upgrades.Artillery, 2), rinha.CalcClanUpgrade(clan.Upgrades.Mission, 1)),
				})
				return
			}
			if strings.ToLower(args[0]) == "missao" || strings.ToLower(args[0]) == "missão" {
				clan = rinha.PopulateClanMissions(clan, galo.Clan, true)
				msg.Reply(context.Background(), session, &disgord.Embed{
					Title:       "Clan missão",
					Color:       65535,
					Description: rinha.MissionToString(clan),
				})
				return
			}
		}
		if len(args) > 1 {
			role := rinha.GetMember(clan, msg.Author.ID)
			if strings.ToLower(args[0]) == "upgrade" {
				if role.Role >= rinha.Admin {
					txt := strings.ToLower(args[1])
					if txt != "membros" && txt != "artilharia" && txt != "missao" {
						msg.Reply(context.Background(), session, msg.Author.Mention()+", Upgrade invalido")
						return
					}
					done := false
					if txt == "missao" {
						price := rinha.CalcClanUpgrade(clan.Upgrades.Mission, 1)
						rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
							if clan.Money >= price {
								clan.Upgrades.Mission++
								clan.Money -= price
								done = true
							}
							return clan, nil
						})
					}
					if txt == "membros" {
						price := rinha.CalcClanUpgrade(clan.Upgrades.Members, 1)
						rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
							if clan.Money >= price {
								clan.Upgrades.Members++
								clan.Money -= price
								done = true
							}
							return clan, nil
						})
					} else if txt == "artilharia" {
						price := rinha.CalcClanUpgrade(clan.Upgrades.Artillery, 2)
						rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
							if clan.Money >= price {
								clan.Upgrades.Artillery++
								clan.Money -= price
								done = true
							}
							return clan, nil
						})
					}
					if done {
						msg.Reply(context.Background(), session, msg.Author.Mention()+", Upgrade comprado com sucesso")
					} else {
						msg.Reply(context.Background(), session, msg.Author.Mention()+", Seu clan nao tem dinheiro para isso")
					}
				}
				return
			}
			if strings.ToLower(args[0]) == "depositar" {
				if clan.Money >= rinha.MaxMoney {
					return
				}
				money, err := strconv.Atoi(args[1])
				if err != nil {
					return
				}
				err = rinha.ChangeMoney(msg.Author.ID, -money, money)
				if err != nil {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem esse dinheiro")
					return
				}
				rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
					clan.Money += money
					return clan, nil
				})
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce depositou dinheiro no clan com sucesso")
				return
			}
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
				text := fmt.Sprintf("%s Voce foi convidado para o clan %s", msg.Mentions[0].Username, galo.Clan)
				utils.Confirm(text, msg.ChannelID, msg.Mentions[0].ID, func() {
					clan = rinha.GetClan(galo.Clan)
					maxMembers := rinha.GetMaxMembers(clan)
					if rinha.IsInClan(clan, user.ID) {
						msg.Reply(context.Background(), session, msg.Author.Mention()+", Este usuario ja esta no clan")
						return
					}
					if len(clan.Members) >= maxMembers {
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
					rinha.UpdateClan(galo.Clan, func(clanUpdate rinha.Clan) (rinha.Clan, error) {
						clanUpdate.Members = clan.Members
						return clanUpdate, nil
					})
					rinha.UpdateGaloDB(user.ID, func(gal rinha.Galo) (rinha.Galo, error) {
						gal.Clan = galo.Clan
						return gal, nil
					})
					msg.Reply(context.Background(), session, user.Mention()+", Voce entrou para o clan **"+galo.Clan+"** com sucesso")
				})

				return
			}
			if strings.ToLower(args[0]) == "admin" {
				user := utils.GetUser(msg, args[1:], session)
				if user.ID == msg.Author.ID || user.Bot {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Usuario invalido")
					return
				}
				if role.Role != rinha.Owner {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Apenas donos podem promover as pessoas a administrador")
					return
				}
				if !rinha.IsInClan(clan, user.ID) {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Este usuario nao esta no clan")
					return
				}
				members := rinha.PromoteMember(clan, user.ID)
				rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
					clan.Members = members
					return clan, nil
				})
				msg.Reply(context.Background(), session, msg.Author.Mention()+", O usuario **"+user.Username+"** foi promovido com sucesso")
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
				rinha.UpdateClan(galo.Clan, func(clan rinha.Clan) (rinha.Clan, error) {
					clan.Members = members
					return clan, nil
				})
				rinha.UpdateGaloDB(user.ID, func(gal rinha.Galo) (rinha.Galo, error) {
					gal.Clan = ""
					return gal, nil
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
				memberMsg += fmt.Sprintf("[**%s**] %s (**%d** XP)\n", member.Role.ToString(), user.Username, member.Xp)
			}
		}
		benefits := rinha.GetBenefits(clan.Xp)
		maxMembers := rinha.GetMaxMembers(clan)
		bg := clan.Background
		embed := &disgord.Embed{
			Title: galo.Clan,
			Color: 65535,
			Footer: &disgord.EmbedFooter{
				Text: "Use j!clan invite <user> para convidar | j!clan remove <user> para remover | j!clan admin <user> para promover membros",
			},
			Description: fmt.Sprintf("Level: **%d** (%d/%d)\nVantagens do clan:\n %s\nMembros (%d/%d):\n %s\nUse **j!clan missao** para ver a missão mensal do clan\nUse **j!clan banco** para ver o banco do clan", level, clan.Xp, rinha.ClanLevelToXp(level), benefits, len(clan.Members), maxMembers, memberMsg),
		}
		if bg != "" {
			embed.Image = &disgord.EmbedImage{
				URL: bg,
			}
		}
		msg.Reply(context.Background(), session, embed)
	}
}
