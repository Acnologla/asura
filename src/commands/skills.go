package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "skills",
		Description: translation.T("SkillsHelp", "pt"),
		Run:         runSkills,
		Cooldown:    15,
		Category:    handler.Profile,
	})
}

func skillsToIntArray(skills []*rinha.NewSkill) []int {
	var arr []int
	for _, skill := range skills {
		arr = append(arr, skill.Skill)
	}
	return arr
}

func skillsToString(skills []*rinha.NewSkill, galo *entities.Rooster) string {
	var str string
	for i, _skill := range skills {
		skill := rinha.Skills[galo.Type-1][_skill.Skill]
		str += fmt.Sprintf("[**%d**]- %s %s\n", i, skill.Name, rinha.SkillToStringFormated(skill, galo))
	}
	return str
}

func generateSkillsOptions(skills []int, galo *entities.Rooster) (opts []*disgord.SelectMenuOption) {
	for i, skill := range skills {
		skill := rinha.Skills[galo.Type-1][skill]
		opts = append(opts, &disgord.SelectMenuOption{
			Value:       strconv.Itoa(i),
			Label:       skill.Name,
			Description: "Equip " + skill.Name,
		})
	}
	return
}

func generateNumberOptions(i int) (opts []*disgord.SelectMenuOption) {
	for j := 0; j < i; j++ {
		toStr := strconv.Itoa(j)
		opts = append(opts, &disgord.SelectMenuOption{
			Value:       toStr,
			Label:       toStr,
			Description: fmt.Sprintf("Skill numero %s", toStr),
		})
	}
	return
}

func runSkills(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(itc.Member.User.ID, "Galos")
	galo := rinha.GetEquippedGalo(&user)
	equipedSkills := rinha.GetEquipedSkills(galo)
	skills := rinha.GetSkills(*galo)
	skillNumber := 0
	genEmbed := func() []*disgord.Embed {
		return []*disgord.Embed{{
			Title: "Skills",
			Color: 65535,
			Description: translation.T("ChangingSkill", translation.GetLocale(itc), map[string]interface{}{
				"skills": skillsToString(equipedSkills, galo),
				"number": skillNumber,
			}),
		}}
	}
	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: genEmbed(),
			Components: []*disgord.MessageComponent{
				{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{
						{
							Type:        disgord.MessageComponentButton + 1,
							MaxValues:   1,
							Placeholder: translation.T("SkillNumber", translation.GetLocale(itc)),
							CustomID:    "number",
							Options:     generateNumberOptions(len(equipedSkills)),
						},
					},
				},
				{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{
						{
							Type:        disgord.MessageComponentButton + 1,
							MaxValues:   1,
							Placeholder: translation.T("SkillSelect", translation.GetLocale(itc)),
							CustomID:    "skill",
							Options:     generateSkillsOptions(skills, galo),
						},
					},
				},
			},
		},
	})
	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		if ic.Member.UserID != itc.Member.UserID {
			return
		}
		name := ic.Data.CustomID
		if len(ic.Data.Values) == 0 {
			return
		}
		number := ic.Data.Values[0]
		data, _ := strconv.Atoi(number)
		if name == "number" {
			skillNumber = data
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.CreateInteractionResponseData{
					Embeds: genEmbed(),
				},
			})
		} else {
			sucess := false
			arrSkill := skillsToIntArray(equipedSkills)
			if rinha.IsIntInList(data, arrSkill) {
				return
			}
			database.User.UpdateUser(itc.Member.UserID, func(u entities.User) entities.User {
				database.User.UpdateEquippedRooster(u, func(r entities.Rooster) entities.Rooster {
					skills := rinha.GetSkills(r)
					if r.ID != galo.ID || r.Resets != galo.Resets || data >= len(skills) {
						return r
					}
					sucess = true
					equipedSkills[skillNumber] = &rinha.NewSkill{
						Skill: data,
					}
					r.Equipped = skillsToIntArray(equipedSkills)
					return r
				})
				return u
			}, "Galos")
			if sucess {
				handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackUpdateMessage,
					Data: &disgord.CreateInteractionResponseData{
						Embeds: genEmbed(),
					},
				})
			}
		}

	}, 160)
	return nil
}
