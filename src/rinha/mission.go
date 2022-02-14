package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

func MissionsToString(id disgord.Snowflake, user *entities.User) string {
	text := ""
	for _, mission := range user.Missions {
		switch mission.Type {
		case entities.Win:
			text += fmt.Sprintf("Vencer %d galos (%d/%d)\nMoney: **%d**\nXp:  **%d**", (mission.Level+1)*3, mission.Progress, (mission.Level+1)*3, 35+5*mission.Level, 40*(mission.Level+1))
		case entities.Fight:
			text += fmt.Sprintf("Batalhar contra %d galos (%d/%d)\nMoney: **%d**\nXp:  **%d**", (mission.Level+1)*6, mission.Progress, (mission.Level+1)*6, 35+5*mission.Level, 40*(mission.Level+1))
		case entities.WinGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Vencer contra 4 galos do tipo %s (%d/4)\nMoney: **60**\nXp:  **130**", className, mission.Progress)
		case entities.FightGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Batalhar contra 8 galos do tipo %s (%d/8)\nMoney: **60**\nXp:  **130**", className, mission.Progress)
		}
		text += "\n"
	}
	return text
}

func CreateMission() entities.Mission {
	missionType := entities.MissionType(utils.RandInt(4))
	level := utils.RandInt(3)
	galoAdv := 0
	if missionType == entities.FightGalo || missionType == entities.WinGalo {
		galoAdv = GetRand()
	}
	return entities.Mission{
		Type:  missionType,
		Level: level,
		Adv:   galoAdv,
	}
}

func PopulateMissions(user *entities.User) []*entities.Mission {
	missions := []*entities.Mission{}
	need := 3 - len(user.Missions)
	for i := 0; need > i && uint64(i) < uint64((uint64(time.Now().Unix())-user.LastMission)/60/60/24); i++ {
		mission := CreateMission()
		missions = append(missions, &mission)
	}
	return missions
}
