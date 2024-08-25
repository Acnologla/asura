package commands

import (
	"asura/src/handler"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "rinhahelp",
		Description: "Tutorial de rinha",
		Run:         runRinhaHelp,
		Cooldown:    3,
	})
}

var rinhaTexts = map[string]string{
	"Galos":                "Existem vários tipos de galos diferentes, com várias raridades e habilidades, que você vai conhecer ao longo do tempo! De cara, você pode ter notado que a cor da borda de alguns galos é diferente da dos outros, e é isso que você usa pra descobrir a raridade de seu lutador. Galos comuns tem a borda branca, galos raros tem a borda azul, galos épicos tem a borda roxa, galos lendários tem a borda laranja, e galos míticos tem a borda colorida. É claro que as raridades influenciam a força de um galo, mas muitas coisas também, como as habilidades escolhidas, as vantagens e desvantagens de tipagem, e as evoluções ou resets que ele tiver. O sistema de habilidades do bot é fácil de entender, e tem 5 slots para ataques. Assim, quando você tiver 6 habilidades ou mais, poderá escolher 5 delas para usar em suas lutas, com `/skills` (Lembre-se de sempre usar o comando ao ganhar um novo ataque, para ver se ele tem algum efeito, ou alterar seu skillset, já que se você não escolher, o galo sempre vai usar os últimos 5 ataques desbloqueados). Por fim, existe o sistema de resets, que podem mudar completamente o jeito que você joga. Ao chegar no nível 35, você consegue dar upgrade no seu galo, que volta ele para o nível 1, fazendo ele perder todos os seus ataques, mas em compensação, ganhando +15% de dano e efeitos pra sempre. Mas cuidado, enquanto o número de resets que você pode conquistar é infinito, cada reset, o seu galo passa a ganhar -25% de xp. Se você conseguir alcançar 5 resets, vai poder evoluir o seu galo, fazendo ele ganhar um novo ataque exclusivo no nível 30!",
	"Economia":             "Por padrão, se ganha 8 moedas a cada vitória em treinos, mas existem formas de aumentar essa quantidade, como upgrades utilitários ou clans de alto nível. Além disso, usar o `/daily`, tentar completar o passe de batalha (visto usando `/battlpass`), cumprir missões com `/missions`, e vender itens ou galos com `/sell` também ajudam a farmar dinheiro. A parte de gastar é muito mais fácil, já que você pode comprar lootboxes usando `/lootbox buy`, trocar o nome de seus galos usando `/changename` e comprar upgrades pro seu clan com o dinheiro depositado no banco usando `/clan banco`. Além disso, você pode enviar dinheiro pra seus amigos com `/givemoney`, e trocar galos, itens, skins e backgrounds com `/trade` (Devido à um bug, apenas quem usar o comando consegue enviar coisas na troca, então para realmente trocar, os dois usuários tem que usar o comando um depois do outro).",
	"Tipos de luta":        "Existem vários tipos de batalhas no bot, e cabe a você testar cada uma! Para começar, com `/rinha`, você consegue duelar contra amigos. Outra forma é a Arena, um sistema em que você luta contra qualquer outra pessoa que use o comando `/arena batalhar` ao mesmo tempo que você. Para entrar na arena, você deve comprar ingressos com `/arena ingresso`, que custam 200 moedas. Em troca, você pode lutar até 12 batalhas, e se perder 3, seu ingresso acaba. Se vencer as 12, recebe a recompensa máxima, ganhando muitas moedas e xp!",
	"Itens e Dungeons":     "Itens são ferramentas que você pode equipar para aumentar seu poder de luta, sua vida, seu xp, e todo tipo de coisa! Eles são permanentes (Não acabam ou são destruídos depois de um certo tempo), e afetam todos os seus galos ao mesmo tempo. Para ganhar itens ou galos novos, você deve entrar nas Dungeons, usando `/dungeon`! As Dungeons são um sistema em que você luta contra chefões cada vez mais poderosos que podem te dar recompensas valiosas. Para ver e equipar seus itens, use `/item`. Itens funcionam em duelos, treinos e na própria Dungeon.",
	"Atributos e Upgrades": "Atributos e Upgrades são sistemas que afetam todos os galos da sua conta, e mudam completamente o jeito que você luta e upa. Eles são escolhidos usando pontos de **userxp** (a cada vitória em treinos se ganha 1), visto usando `/money`. A cada 100 userxp, você pode alocar um ponto em um dos 4 atributos, Vida, Dano Bruto, Dano de Efeitos e Sorte, usando `/attributes`, sendo que o máximo em cada um é 300 pontos. Cada ponto em **Vida** te dá __+3__ de vida máxima, cada ponto em **Dano Bruto** te dá __+0,6__ de dano em todos os seus ataques, cada ponto em **Dano de Efeitos** te dá __+0,4__ de dano em todos os seus efeitos, e cada ponto em **Sorte** __multiplica o seu pity__. Já os Upgrades são 3 Árvores de passivas, e você pode escolher uma das 3 ao alcançar **150** userxp, usando `/upgrades` (Cuidado, o bot conta a partir de zero, como 0, 1, 2, então o terceiro upgrade tem o id 2 e não 3). Você não pode trocar seu Upgrade depois de escolher, então pense bem! Os 3 iniciais são **Vida** (__+20__ de hp máximo), **Esquiva** (__4%__ de chance de desviar de ataques), e **Utilitários** (__+3__ de xp por treino). Quando você conseguir __800__, __5.000__ e __15.000__ userxp, poderá escolher Upgrades mais avançados, mas só da Árvore que você pegou no começo!",
	"Clan":                 "Os clans são grupos de jogadores que upam, treinam e lutam juntos. Para criar um clã, você deve usar `/clan create`, que custa 1000 moedas. No nível 1, os clans dão +10% de xp bônus em treinos, e a cada nível, ele ganha mais buffs. Com `/clan view`, você pode ver os status do seu clan, além da lista de membros. Se quiser convidar seus amigos, use `/clan invite`, mas você deve ser o Dono ou um Administrador. Todo mês o clwn recebe uma nova missão, vista usando `/clan mission`. A missão é matar 20.000 galos, e todos os membros do clan ajudam a completar. No fim, todos ganham a recompensa, que é uma enorme quantidade de moedas e xp. Por fim, existe a raid do clan, em que até 7 membros lutam contra um boss, que fica mais fácil quanto mais pessoas entrar na batalha. Você pode entrar na raid do clan a cada 4 horas, usando `/clan battle`.",
	"Rinhas ranqueadas":    "Quer ganhar elo e subir de rank? Quer ganhar prêmios valiosos e aparecer no ranking geral de jogadores? Então é com as rankeadas que você vai conseguir! Para começar, use `/ranked battle`, que funciona do mesmo jeito que a Arena, em que outros jogadores tem que usar o comando ao mesmo tempo que você. A cada vitória, você ganha elo, e com pontos suficientes, sobe de rank, ganhando muitas recompensas!",
	"Cosméticos":           "Você pode ganhar skins e backgrounds ao longo da sua jornada, ao comprar lootboxes cosméticas, que custam 500 moedas. Skins alteram o visual dos seus galos, e backgrounds alteram a imagem que aparece nos comandos `/galo` e `/profile`. Ao comprar a lootbox cosmética, pode vir qualquer skin de qualquer galo, além de qualquer background, que são equipados com `/skins` e `/equipitem`.",
	"Como jogar":           "Para começar, use o comando `/galo` e descubra qual é seu galo inicial! Existem 5 raridades de galos, comuns, raros, épicos, lendários e míticos, e você pode começar apenas com comuns, raros ou épicos. Depois, use `/daily` para receber o bônus diário de moedas e xp, que aumenta caso você esteja em meu Servidor Suporte. Use `/train` para treinar seu galo, aumentando sua força ao ganhar batalhas. Para ver suas tarefas, use `/mission`, mas saiba que essas missões só podem ser cumpridas em treinos. Eu tenho muitos comandos divertidos, e espero que aproveite!",
}

var textsOrder = []string{
	"Galos",
	"Economia",
	"Tipos de luta",
	"Itens e Dungeons",
	"Atributos e Upgrades",
	"Clan",
	"Rinhas ranqueadas",
	"Cosméticos",
	"Como jogar",
}

func generateComponents() (buttons []*disgord.MessageComponent) {
	opts := []*disgord.SelectMenuOption{}
	for _, name := range textsOrder {
		opts = append(opts, &disgord.SelectMenuOption{
			Label: name,
			Value: name,
		})
	}
	buttons = []*disgord.MessageComponent{
		{
			Type:        disgord.MessageComponentButton + 1,
			Options:     opts,
			Placeholder: "Selecione a categoria",
			MaxValues:   1,
			CustomID:    "1",
		},
	}
	return
}

func runRinhaHelp(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	embed := disgord.Embed{
		Title: "Tutorial de rinha",
		Footer: &disgord.EmbedFooter{
			Text: "Selecione a baixo a categoria do tutorial",
		},
		Description: rinhaTexts["Como jogar"],
		Color:       65535,
	}
	components := []*disgord.MessageComponent{
		{
			Type:       disgord.MessageComponentActionRow,
			Components: generateComponents(),
		},
	}
	err := handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds:     []*disgord.Embed{&embed},
			Components: components,
		},
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	user := itc.Member.User

	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		if ic.Member.UserID == user.ID {
			if len(ic.Data.Values) == 0 {
				return
			}
			key := ic.Data.Values[0]
			embed.Title = key
			embed.Description = rinhaTexts[key]
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.CreateInteractionResponseData{
					Embeds:     []*disgord.Embed{&embed},
					Components: components,
				},
			})
			handler.Client.EditInteractionResponse(ctx, itc, &disgord.UpdateMessage{
				Embeds: &([]*disgord.Embed{&embed}),
			})
		}
	}, 500)
	return nil

}
