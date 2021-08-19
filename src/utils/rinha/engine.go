package rinha

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

type RinhaOptions struct {
	GaloAuthor  Galo
	GaloAdv     Galo
	AuthorLevel int
	AdvLevel    int
	AuthorName  string
	AdvName     string
	NoItems     bool
	IDs         [2]disgord.Snowflake
}

var RinhaColors = [2]int{65280, 16711680}

func EditRinhaMessage(ic *disgord.InteractionCreate, newMessage *disgord.CreateMessageParams, message *disgord.Message) {
	if ic == nil {
		EditRinhaComponents(message, newMessage)
		return
	}
	handler.Client.SendInteractionResponse(context.Background(), ic, &disgord.InteractionResponse{
		Type: disgord.UpdateMessage,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Components: newMessage.Components,
			Embeds:     []*disgord.Embed{newMessage.Embed},
		},
	})
}

func EditRinhaComponents(message *disgord.Message, newMessage *disgord.CreateMessageParams) {
	utils.Try(func() error {
		msgUpdater := handler.Client.Channel(message.ChannelID).Message(message.ID).UpdateBuilder()
		msgUpdater.SetEmbed(newMessage.Embed)
		msgUpdater.Set("components", newMessage.Components)
		_, err := msgUpdater.Execute()
		return err
	}, 5)
}

func EditRinhaEmbed(message *disgord.Message, embed *disgord.Embed) {
	utils.Try(func() error {
		msgUpdater := handler.Client.Channel(message.ChannelID).Message(message.ID).UpdateBuilder()
		msgUpdater.SetEmbed(embed)
		_, err := msgUpdater.Execute()
		return err
	}, 5)
}
func GetImageTile(first *Galo, sec *Galo, turn int) string {
	if turn == 0 {
		return Sprites[turn^1][sec.Type-1]
	}
	return Sprites[turn^1][first.Type-1]
}

var RinhaEmojis = [2]string{"<:sverde:744682222644363296>", "<:svermelha:744682249408217118>"}

func CheckDead(battle Battle) bool {
	return 0 >= battle.Fighters[0].Life || 0 >= battle.Fighters[1].Life
}

func EffectToStr(effect *Result, affected string, author string, battle *Battle) string {
	if effect.Dodge {
		return fmt.Sprintf("%s **%s** desviou do ataque\n", RinhaEmojis[battle.GetReverseTurn()], affected)
	}
	if effect.Effect == Damaged {
		if effect.Skill.Self {
			return fmt.Sprintf("%s **%s** Usou **%s** em si mesmo\n", RinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name)
		}
		if effect.Reflected {
			return fmt.Sprintf("%s **%s** Refletiu o ataque **%s** causando **%d** de dano\n", RinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name, effect.Damage)
		}
		return fmt.Sprintf("%s **%s** Usou **%s** causando **%d** de dano\n", RinhaEmojis[battle.GetReverseTurn()], author, effect.Skill.Name, effect.Damage)
	} else if effect.Effect == Effected {
		effectLiteral := Effects[effect.EffectID]
		if effect.Self {
			return fmt.Sprintf(effectLiteral.Phrase+"\n", author, effect.Damage)
		}
		return fmt.Sprintf(effectLiteral.Phrase+"\n", affected, effect.Damage)
	} else if effect.Effect == NotEffective {
		return "**reduzido**\n"
	}
	return ""
}
func getRandAvalaibleSkills(skills []*NewSkill, round int) int {
	newSkillArr := []*NewSkill{}
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

func GenerateRinhaButtons(round int, skills []*NewSkill, embed *disgord.CreateMessageParams, author *Fighter) *disgord.CreateMessageParams {
	embed.Components[0].Components = []*disgord.MessageComponent{}
	for i, skill := range skills {
		embed.Components[0].Components = append(embed.Components[0].Components, &disgord.MessageComponent{
			CustomID: strconv.Itoa(i),
			Disabled: round-skill.Cooldown <= 3 && skill.Cooldown != 0,
			Type:     disgord.MessageComponentButton,
			Style:    disgord.Primary,
			Label:    Skills[author.Galo.Type-1][skill.Skill].Name,
		})
	}
	return embed
}

func getSkill(battle *Battle, options *RinhaOptions, message *disgord.Message, embed *disgord.CreateMessageParams, round int) (int, *disgord.InteractionCreate) {
	skill := make(chan int)
	turn := battle.GetTurn()
	author := battle.Fighters[turn]
	skills := author.Equipped
	var interaction *disgord.InteractionCreate
	go handler.RegisterBHandler(message, func(ic *disgord.InteractionCreate) {
		if ic.Member.User.ID == author.ID {
			skillId, _ := strconv.Atoi(ic.Data.CustomID)
			interaction = ic
			skill <- skillId
		}

	}, 25)
	selected := 0
	select {
	case selected = <-skill:
	case <-time.After(time.Second * 25):
		selected = getRandAvalaibleSkills(skills, round)
	}

	return selected, interaction
}

func RinhaEngine(battle *Battle, options *RinhaOptions, message *disgord.Message, embed *disgord.Embed) (int, *Battle) {
	var lastEffects string
	round := 1
	newMsg := &disgord.CreateMessageParams{
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
			URL: GetImageTile(&options.GaloAuthor, &options.GaloAdv, battle.GetReverseTurn()),
		}
		if CheckDead(*battle) {
			winnerTurn := battle.GetReverseTurn()
			embed.Image = &disgord.EmbedImage{
				URL: GetImageTile(&options.GaloAuthor, &options.GaloAdv, battle.GetTurn()),
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