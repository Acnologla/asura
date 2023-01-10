package engine

import (
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

type RinhaOptions struct {
	GaloAuthor  *entities.User
	GaloAdv     *entities.User
	MessageID   [2]*disgord.Message
	AuthorLevel int
	AdvLevel    int
	AuthorName  string
	AdvName     string
	NoItems     bool
	IDs         [2]disgord.Snowflake
	Waiting     []*entities.User
	Usernames   []string
	Images      [2]string
}

var RinhaColors = [2]int{65280, 16711680}

func EditRinhaMessage(ic *disgord.InteractionCreate, newMessage *disgord.CreateMessage, message *disgord.Message) {
	if ic == nil {
		EditRinhaComponents(message, newMessage)
		return
	}
	handler.Client.SendInteractionResponse(context.Background(), ic, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackUpdateMessage,
		Data: &disgord.CreateInteractionResponseData{
			Components: newMessage.Components,
			Embeds:     []*disgord.Embed{newMessage.Embed},
		},
	})
}

func EditRinhaComponents(message *disgord.Message, newMessage *disgord.CreateMessage) {
	utils.Try(func() error {
		_, err := handler.Client.Channel(message.ChannelID).Message(message.ID).Update(&disgord.UpdateMessage{
			Embeds:     &([]*disgord.Embed{newMessage.Embed}),
			Components: &newMessage.Components,
		})

		return err
	}, 5)
}

func editRinhaEmbed(msg *disgord.Message, embed *disgord.Embed) {
	utils.Try(func() error {
		_, err := handler.Client.Channel(msg.ChannelID).Message(msg.ID).Update(&disgord.UpdateMessage{
			Embeds: &([]*disgord.Embed{embed}),
		})
		return err
	}, 5)
}

func EditRinhaEmbed(msg *disgord.Message, embed *disgord.Embed, msgs [2]*disgord.Message) {
	if msgs[0] != nil {
		for _, msg := range msgs {
			editRinhaEmbed(msg, embed)
		}
	} else {
		editRinhaEmbed(msg, embed)
	}

}
func GetImageTile(options *RinhaOptions, turn int) string {
	if turn == 0 {
		return options.Images[1]
	}
	return options.Images[0]
}

var RinhaEmojis = [2]string{"<:sverde:744682222644363296>", "<:svermelha:744682249408217118>"}

func CheckDead(battle rinha.Battle) bool {
	return 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life
}

func EffectToStr(effect *rinha.Result, affected string, author string, battle *rinha.Battle) string {
	if effect.ExtraMsg != "" {
		return fmt.Sprintf("%s **%s** %s\n", RinhaEmojis[battle.GetReverseTurn()], affected, effect.ExtraMsg)
	}
	if effect.Reset {
		return "O tempo volta e a batalha recomeça"
	}
	if effect.Dodge {
		return fmt.Sprintf("%s **%s** desviou do ataque\n", RinhaEmojis[battle.GetReverseTurn()], affected)
	}
	if effect.Effect == rinha.Damaged {
		if effect.Skill.Self {
			return fmt.Sprintf("%s **%s** Usou **%s** em si mesmo\n", RinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name)
		}
		if effect.Reflected {
			return fmt.Sprintf("%s **%s** Refletiu o ataque **%s** causando **%d** de dano\n", RinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name, effect.Damage)
		}
		return fmt.Sprintf("%s **%s** Usou **%s** causando **%d** de dano\n", RinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name, effect.Damage)
	} else if effect.Effect == rinha.Effected {
		effectLiteral := rinha.Effects[effect.EffectID]
		if effect.Self {
			return fmt.Sprintf(effectLiteral.Phrase+"\n", author, effect.Damage)
		}
		return fmt.Sprintf(effectLiteral.Phrase+"\n", affected, effect.Damage)
	} else if effect.Effect == rinha.NotEffective {
		return "**reduzido**\n"
	}
	return ""
}
func getRandAvalaibleSkills(skills []*rinha.NewSkill, round int) int {
	newSkillArr := []*rinha.NewSkill{}
	for _, skill := range skills {
		if round-skill.Cooldown > 3 || skill.Cooldown == 0 {
			newSkillArr = append(newSkillArr, skill)
		}
	}
	selected := newSkillArr[utils.RandInt(len(newSkillArr))]
	for i, skill := range skills {
		if skill.Skill == selected.Skill {
			return i
		}
	}
	return -1
}

func GenerateRinhaButtons(round int, skills []*rinha.NewSkill, embed *disgord.CreateMessage, author *rinha.Fighter) *disgord.CreateMessage {
	embed.Components[0].Components = []*disgord.MessageComponent{}
	for i, skill := range skills {
		embed.Components[0].Components = append(embed.Components[0].Components, &disgord.MessageComponent{
			CustomID: strconv.Itoa(i),
			Disabled: round-skill.Cooldown <= 3 && skill.Cooldown != 0,
			Type:     disgord.MessageComponentButton,
			Style:    disgord.Primary,
			Label:    rinha.Skills[author.Galo.Type-1][skill.Skill].Name,
		})
	}
	return embed
}

func getSkill(battle *rinha.Battle, options *RinhaOptions, message *disgord.Message, embed *disgord.CreateMessage, round int) (int, *disgord.InteractionCreate) {
	skill := make(chan int)
	turn := battle.GetTurn()
	author := battle.Fighters[turn]
	skills := author.Equipped
	var interaction *disgord.InteractionCreate
	/*	go handler.RegisterBHandler(message, func(ic *disgord.InteractionCreate) {
			if ic.Member.User.ID == author.ID {
				skillId, _ := strconv.Atoi(ic.Data.CustomID)
				interaction = ic
				skill <- skillId
			}

		}, 25)
	*/
	selected := 0
	select {
	case selected = <-skill:
	case <-time.After(time.Second * 25):
		selected = getRandAvalaibleSkills(skills, round)
	}

	return selected, interaction
}

func RinhaEngineNew(battle *rinha.Battle, options *RinhaOptions, message *disgord.Message, embed *disgord.Embed) (int, *rinha.Battle) {
	var lastEffects string
	round := 1
	newMsg := &disgord.CreateMessage{
		Embed: embed,
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
			},
		},
	}
	newMsg = GenerateRinhaButtons(round, battle.Fighters[0].Equipped, newMsg, battle.Fighters[0])
	EditRinhaComponents(message, newMsg)
	for {
		var skill int
		var ic *disgord.InteractionCreate
		if battle.ReflexType == 1 || battle.ReflexType == 2 {
			skill = -1
			round--
			time.Sleep(time.Second * 2)
		} else {
			skill, ic = getSkill(battle, options, message, newMsg, round)
			battle.Fighters[battle.GetTurn()].Equipped[skill].Cooldown = round
		}
		effects := battle.Play(skill)
		var text string
		authorName := options.AuthorName
		affectedName := options.AdvName
		turn := battle.GetTurn()

		if turn == 0 {
			authorName = options.AdvName
			affectedName = options.AuthorName
		}
		for _, effect := range effects {
			text += EffectToStr(effect, affectedName, authorName, battle)
		}
		if round >= 70 {
			if battle.Fighters[1].Life >= battle.Fighters[0].Life {
				text += "\n" + options.AuthorName + " Foi executado"
				battle.Fighters[0].Life = 0
				battle.Turn = false
			} else {
				battle.Fighters[1].Life = 0
				text += "\n" + options.AdvName + " Foi executado"
				battle.Turn = true
			}
		}
		embed.Color = RinhaColors[battle.GetReverseTurn()]
		embed.Description = lastEffects + "\n" + text
		if (battle.Stun || battle.ReflexType == 3) && !CheckDead(*battle) {
			battle.Turn = !battle.Turn
			if battle.Stun {
				battle.Stun = false
			}
			if battle.ReflexType == 3 {
				battle.ReflexType = 0
			}
			round--
			turn = battle.GetTurn()
		}
		if battle.ReflexType != 1 && battle.ReflexType != 2 {
			newMsg = GenerateRinhaButtons(round, battle.Fighters[turn].Equipped, newMsg, battle.Fighters[turn])
		}

		u, _ := handler.Client.User(options.IDs[turn]).Get()
		avatar, _ := u.AvatarURL(128, true)
		embed.Fields = []*disgord.EmbedField{
			{
				Name:   fmt.Sprintf("%s Level %d", options.AuthorName, options.AuthorLevel),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, battle.Fighters[0].MaxLife),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s Level %d", options.AdvName, options.AdvLevel),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, battle.Fighters[1].MaxLife),
				Inline: true,
			},
		}
		embed.Footer = &disgord.EmbedFooter{
			Text:    u.Username,
			IconURL: avatar,
		}
		embed.Image = &disgord.EmbedImage{
			URL: GetImageTile(options, battle.GetReverseTurn()),
		}
		if CheckDead(*battle) {
			winnerTurn := battle.GetReverseTurn()
			embed.Image = &disgord.EmbedImage{
				URL: GetImageTile(options, battle.GetTurn()),
			}
			if winnerTurn == 1 {
				embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.AdvName)
			} else {
				embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.AuthorName)
			}
			newMsg.Embed = embed
			newMsg.Components = []*disgord.MessageComponent{}
			EditRinhaMessage(ic, newMsg, message)
			return winnerTurn, battle
		}
		newMsg.Embed = embed
		EditRinhaMessage(ic, newMsg, message)
		lastEffects = text
		round++
	}
}

func GetUsername(name, realName string) string {
	if name == "" {
		return realName
	}
	return name
}

func RinhaEngine(battle *rinha.Battle, options *RinhaOptions, msg *disgord.Message, embed *disgord.Embed) (int, *rinha.Battle) {
	var lastEffects string
	round := 0
	for {
		effects := battle.Play(-1)
		var text string

		authorName := GetUsername(battle.Fighters[0].Username, options.AuthorName)
		affectedName := options.AdvName

		turn := battle.GetTurn()

		if turn == 0 {
			authorName = options.AdvName
			affectedName = GetUsername(battle.Fighters[0].Username, options.AuthorName)
		}

		for _, effect := range effects {
			text += EffectToStr(effect, affectedName, authorName, battle)
		}
		if round >= 35 {
			if battle.Fighters[1].Life >= battle.Fighters[0].Life {
				text += "\n" + options.AuthorName + " Foi executado"
				battle.Fighters[0].Life = 0
				battle.Turn = false
			} else {
				battle.Fighters[1].Life = 0
				text += "\n" + options.AdvName + " Foi executado"
				battle.Turn = true
			}
		}
		embed.Color = RinhaColors[battle.GetReverseTurn()]
		embed.Description = lastEffects + "\n" + text
		level := rinha.CalcLevel(battle.Fighters[0].Galo.Xp)
		name := GetUsername(battle.Fighters[0].Username, options.AuthorName)
		embed.Fields = []*disgord.EmbedField{
			{
				Name:   fmt.Sprintf("%s Level %d", name, level),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, battle.Fighters[0].MaxLife),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s Level %d", options.AdvName, options.AdvLevel),
				Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, battle.Fighters[1].MaxLife),
				Inline: true,
			},
		}

		embed.Image = &disgord.EmbedImage{
			URL: GetImageTile(options, turn),
		}
		if 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life {
			if len(battle.Waiting) > 1 && 0 >= battle.Fighters[0].Life {
				battle.Waiting = SpliceWaiting(battle.Waiting, battle.WaitingN)
			} else {
				winnerTurn := battle.GetReverseTurn()

				if winnerTurn == 1 {
					embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.AdvName)
				} else {
					embed.Description += fmt.Sprintf("\n**%s** venceu a batalha!", options.AuthorName)
				}
				EditRinhaEmbed(msg, embed, options.MessageID)
				return winnerTurn, battle
			}
		}

		EditRinhaEmbed(msg, embed, options.MessageID)
		lastEffects = text
		round++
		time.Sleep(4 * time.Second)
	}
}
func SpliceWaiting(slice []*rinha.Fighter, s int) []*rinha.Fighter {
	newArr := []*rinha.Fighter{}
	for i, element := range slice {
		if i != s {
			newArr = append(newArr, element)
		}
	}
	return newArr
}
func ExecuteRinha(itc *disgord.InteractionCreate, session disgord.Session, options RinhaOptions, newEngine bool) (int, *rinha.Battle) {
	authorGal := rinha.GetEquippedGalo(options.GaloAuthor)
	advGal := rinha.GetEquippedGalo(options.GaloAdv)
	authorURL := rinha.GetGaloImage(authorGal, options.GaloAuthor.Items)
	advURL := rinha.GetGaloImage(advGal, options.GaloAdv.Items, "reverse")
	options.Images = [2]string{authorURL, advURL}
	if rinha.HasUpgrade(options.GaloAuthor.Upgrades, 2, 1, 0) {
		advGal.Xp = rinha.CalcXP(rinha.CalcLevel(advGal.Xp)) - 1
		if rinha.HasUpgrade(options.GaloAuthor.Upgrades, 2, 1, 0, 0) {
			advGal.Xp = rinha.CalcXP(rinha.CalcLevel(advGal.Xp)) - 1
		}
	}
	if rinha.HasUpgrade(options.GaloAdv.Upgrades, 2, 1, 0) {
		authorGal.Xp = rinha.CalcXP(rinha.CalcLevel(authorGal.Xp)) - 1
		if rinha.HasUpgrade(options.GaloAdv.Upgrades, 2, 1, 0, 0) {
			authorGal.Xp = rinha.CalcXP(rinha.CalcLevel(authorGal.Xp)) - 1
		}
	}
	if authorGal.Xp < 0 {
		authorGal.Xp = 0
	}
	if advGal.Xp < 0 {
		advGal.Xp = 0
	}
	options.AdvLevel = rinha.CalcLevel(advGal.Xp)
	options.AuthorLevel = rinha.CalcLevel(authorGal.Xp)
	u, _ := handler.Client.User(options.IDs[0]).Get()
	var avatar string
	if u != nil {
		avatar, _ = u.AvatarURL(128, true)
	}
	embed := &disgord.Embed{
		Title: "Briga de galo",
		Color: RinhaColors[0],
		Footer: &disgord.EmbedFooter{
			Text:    u.Username,
			IconURL: avatar,
		},
		Image: &disgord.EmbedImage{
			URL: GetImageTile(&options, 1),
		},
	}
	if newEngine {
		if 3 > options.AdvLevel || 3 > options.AuthorLevel {
			return -1, nil
		}
	}
	var msg *disgord.Message
	var err error
	if options.MessageID[0] == nil {
		msg, err = session.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
			Content: itc.Member.User.Mention(),
			Embeds:  []*disgord.Embed{embed},
		})
		if err != nil {
			return -1, nil
		}
	}

	battle := rinha.CreateBattle(options.GaloAuthor, options.GaloAdv, options.NoItems, options.IDs[0], options.IDs[1], options.Waiting, options.Usernames)
	embed.Fields = []*disgord.EmbedField{
		{
			Name:   fmt.Sprintf("%s Level %d", options.AuthorName, options.AuthorLevel),
			Value:  fmt.Sprintf("%d/%d", battle.Fighters[0].Life, battle.Fighters[0].MaxLife),
			Inline: true,
		},
		{
			Name:   fmt.Sprintf("%s Level %d", options.AdvName, options.AdvLevel),
			Value:  fmt.Sprintf("%d/%d", battle.Fighters[1].Life, battle.Fighters[1].MaxLife),
			Inline: true,
		},
	}
	EditRinhaEmbed(msg, embed, options.MessageID)
	/*	if newEngine {
			return RinhaEngineNew(&battle, &options, message, embed)
		}
	*/
	return RinhaEngine(&battle, &options, msg, embed)
}
