package rinha

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

type MissionType int

const (
	Win MissionType = iota
	Fight
	WinGalo
	FightGalo
)

type Mission struct {
	Type     MissionType
	Level    int
	Progress int
	Adv      int
}

func MissionsToString(id disgord.Snowflake, galo Galo) string {
	missions := PopulateMissions(galo)
	if len(missions) > len(galo.Missions) {
		galo.Missions = missions
		UpdateGaloDB(id, func(galo Galo) (Galo, error) {
			galo.Missions = missions
			galo.LastMission = uint64(time.Now().Unix())
			return galo, nil
		})
	}
	text := ""
	for _, mission := range galo.Missions {
		switch mission.Type {
		case Win:
			text += fmt.Sprintf("Vencer %d galos (%d/%d)\nMoney: **%d**\nXp:  **%d**", (mission.Level+1)*3, mission.Progress, (mission.Level+1)*3, 30+5*mission.Level, 35*(mission.Level+1))
		case Fight:
			text += fmt.Sprintf("Batalhar contra %d galos (%d/%d)\nMoney: **%d**\nXp:  **%d**", (mission.Level+1)*6, mission.Progress, (mission.Level+1)*6, 30+5*mission.Level, 35*(mission.Level+1))
		case WinGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Vencer contra 4 galos do tipo %s (%d/4)\nMoney: **38**\nXp:  **80**", className, mission.Progress)
		case FightGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Batalhar contra 8 galos do tipo %s (%d/8)\nMoney: **38**\nXp:  **80**", className, mission.Progress)
		}
		text += "\n"
	}
	return text
}

func RemoveMission(missions []Mission, i int) []Mission {
	newMissions := []Mission{}
	for j, mission := range missions {
		if j != i {
			newMissions = append(newMissions, mission)
		}
	}
	return newMissions
}

func CompleteMission(id disgord.Snowflake, galo, galoAdv Galo, winner bool, msg *disgord.Message) {
	tempGalo, _ := GetGaloDB(id)
	galo.Missions = tempGalo.Missions
	galo.LastMission = tempGalo.LastMission
	if len(galo.Missions) == 3 {
		UpdateGaloDB(id, func(galo Galo) (Galo, error) {
			galo.LastMission = uint64(time.Now().Unix())
			return galo, nil
		})
	}
	missions := PopulateMissions(galo)
	if len(missions) > len(galo.Missions) {
		galo.Missions = missions
		UpdateGaloDB(id, func(galo Galo) (Galo, error) {
			galo.Missions = missions
			galo.LastMission = uint64(time.Now().Unix())
			return galo, nil
		})
	}
	xp := 0
	money := 0
	toRemove := []int{}
	for i, mission := range galo.Missions {
		old := mission.Progress
		done := false
		switch mission.Type {
		case Win:
			if winner {
				mission.Progress++
				if mission.Progress == (mission.Level+1)*3 {
					xp += 35 * (mission.Level + 1)
					money += 30 + (5 * mission.Level)
					done = true
				}
			}
		case Fight:
			mission.Progress++
			if mission.Progress == (mission.Level+1)*6 {
				xp += 35 * (mission.Level + 1)
				money += 30 + (5 * mission.Level)
				done = true
			}
		case WinGalo:
			if winner && galoAdv.Type == mission.Adv {
				mission.Progress++
				if mission.Progress == 4 {
					xp += 80
					money += 38
					done = true
				}
			}
		case FightGalo:
			if galoAdv.Type == mission.Adv {
				mission.Progress++
				if mission.Progress == 8 {
					xp += 80
					money += 38
					done = true
				}
			}
		}
		if mission.Progress != old {
			if done {
				toRemove = append(toRemove, i)
			} else {
				galo.Missions[i] = mission
			}
		}
	}
	if galo.GaloReset > 0 {
		xp = xp / (galo.GaloReset + 1)
	}
	if HasUpgrade(galo.Upgrades, 0, 1) {
		money += 3
	}
	if len(toRemove) > 0 {
		text := "missão"
		if len(toRemove) > 1 {
			text = "missoes"
		}
		msg.Reply(context.Background(), handler.Client, fmt.Sprintf("<@%d> voce completou **%d** %s e recebeu **%d** de money e **%d** de xp", id, len(toRemove), text, money, xp))
	}
	for i := len(toRemove) - 1; i >= 0; i-- {
		galo.Missions = RemoveMission(galo.Missions, toRemove[i])
	}
	MissionUpdate(id, galo, xp, money)
}

func MissionUpdate(id disgord.Snowflake, galo Galo, xp int, money int) {
	UpdateGaloDB(id, func(gal Galo) (Galo, error) {
		gal.Missions = galo.Missions
		if xp != 0 {
			gal.Xp += xp
		}
		return gal, nil
	})
	if xp != 0 {
		ChangeMoney(id, money, 0)
	}
}

func CreateMission(galo Galo) Mission {
	missionType := MissionType(utils.RandInt(4))
	level := utils.RandInt(3)
	galoAdv := 0
	if missionType == FightGalo || missionType == WinGalo {
		galoAdv = GetRand()
	}
	return Mission{
		Type:  missionType,
		Level: level,
		Adv:   galoAdv,
	}
}

func PopulateMissions(galo Galo) []Mission {
	missions := galo.Missions
	need := 3 - len(galo.Missions)
	for i := 0; need > i && uint64(i) < uint64((uint64(time.Now().Unix())-galo.LastMission)/60/60/24); i++ {
		mission := CreateMission(galo)
		missions = append(missions, mission)
	}
	return missions
}
