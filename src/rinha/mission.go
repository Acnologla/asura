package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

func MissionsToString(id disgord.Snowflake, user *entities.User) (arr []string) {
	for _, mission := range user.Missions {
		text := ""
		switch mission.Type {
		case entities.Win:
			text += fmt.Sprintf("Vencer %d galos (%d/%d)\nMoney: %d\nXp:  %d", (mission.Level+1)*3, mission.Progress, (mission.Level+1)*3, 45+5*mission.Level, 55*(mission.Level+1))
		case entities.Fight:
			text += fmt.Sprintf("Batalhar contra %d galos (%d/%d)\nMoney: %d\nXp:  %d", (mission.Level+1)*6, mission.Progress, (mission.Level+1)*6, 45+5*mission.Level, 55*(mission.Level+1))
		case entities.WinGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Vencer contra 3 galos do tipo %s (%d/3)\nMoney: 110\nXp:  260", className, mission.Progress)
		case entities.FightGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Batalhar contra 6 galos do tipo %s (%d/6)\nMoney: 110\nXp:  260", className, mission.Progress)
		case entities.PlayTrial:
			text += "Jogue uma batalha na trial (0/1)\nMoney: 50\nXp:  100"
		case entities.WinRaid:
			text += "Vencer uma raid (0/1)\nMoney: 125\nXp:  285"
		case entities.FightTower:
			text += fmt.Sprintf("Batalhar contra %d galos na torre (%d/%d)\nMoney: %d\nXp:  %d", (mission.Level+1)*6, mission.Progress, (mission.Level+1)*6, 50+5*mission.Level, 60*(mission.Level+1))
		case entities.WinDungeon:
			text += fmt.Sprintf("Vencer %d galos na dungeon (%d/%d)\nMoney: %d\nXp:  %d", (mission.Level+1)*3, mission.Progress, (mission.Level+1)*3, 55+5*mission.Level, 60*(mission.Level+1))
		}

		arr = append(arr, text)
	}
	return
}

func CreateMission() entities.Mission {
	missionType := entities.MissionType(utils.RandInt(8))
	level := utils.RandInt(5)
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
	need := 4 - len(user.Missions)
	for i := 0; need > i && uint64(i) < uint64((uint64(time.Now().Unix())-user.LastMission)/60/60/12); i++ {
		mission := CreateMission()
		missions = append(missions, &mission)
	}
	return missions
}
