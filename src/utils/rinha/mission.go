package rinha

import (
	"fmt"
	"math/rand"
	"time"
	"asura/src/handler"
	"github.com/andersfylling/disgord"
	"context"
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

func MissionsToString(id disgord.Snowflake,galo Galo) string {
	missions := PopulateMissions(galo)
	if len(missions) > len(galo.Missions) {
		galo.Missions = missions
		galo.LastMission = uint64(time.Now().Unix())
		UpdateGaloDB(id, map[string]interface{}{
			"lastMission": galo.LastMission,
			"missions":    galo.Missions,
		})
	}
	text := ""
	for _, mission := range galo.Missions {
		switch mission.Type {
		case Win:
			text += fmt.Sprintf("Vencer %d galos (%d/%d)\nMoney: **%d**\nXp:  **%d**", mission.Level+1, mission.Progress, mission.Level+1, 20+5*mission.Level, 30*(mission.Level+1))
		case Fight:
			text += fmt.Sprintf("Batalhar contra %d galos (%d/%d)\nMoney: **%d**\nXp:  **%d**", (mission.Level+1)*2, mission.Progress, (mission.Level+1)*2, 20+5*mission.Level, 30*(mission.Level+1))
		case WinGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Vencer contra galo do tipo %s (0/1)\nMoney: **30**\nXp:  **90**", className)
		case FightGalo:
			className := Classes[mission.Adv].Name
			text += fmt.Sprintf("Batalhar contra 3 galos do tipo %s (%d/3)\nMoney: **30**\nXp:  **90**", className, mission.Progress)
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
	if len(galo.Missions) == 3{
		galo.LastMission = uint64(time.Now().Unix())
		UpdateGaloDB(id, map[string]interface{}{
			"lastMission": galo.LastMission,
		})
	}
	missions := PopulateMissions(galo)
	if len(missions) > len(galo.Missions) {
		galo.Missions = missions
		galo.LastMission = uint64(time.Now().Unix())
		UpdateGaloDB(id, map[string]interface{}{
			"lastMission": galo.LastMission,
			"missions":    galo.Missions,
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
				if mission.Progress == mission.Level+1 {
					xp += 30 * (mission.Level + 1)
					money += 20 + (5 * mission.Level)
					done = true
				}
			}
		case Fight:
			mission.Progress++
			if mission.Progress == (mission.Level+1)*2 {
				xp += 30 * (mission.Level + 1)
				money += 20 + (5 * mission.Level)
				done = true
			}
		case WinGalo:
			if winner && galoAdv.Type == mission.Adv {
				mission.Progress++
				xp += 90
				money += 30
				done = true
			}
		case FightGalo:
			if galoAdv.Type == mission.Adv {
				mission.Progress++
				if mission.Progress == 3 {
					xp += 90
					money += 30
					done = true
				}
			}
		}
		if mission.Progress != old {
			if done {
				toRemove = append(toRemove, i)
				msg.Reply(context.Background(), handler.Client, fmt.Sprintf("<@%d> voce completou uma missão e recebeu **%d** de money e **%d** de xp", id, money, xp))
			} else {
				galo.Missions[i] = mission
			}
		}
	}
	for i := len(toRemove)-1;i >= 0;i--{
		galo.Missions = RemoveMission(galo.Missions, toRemove[i])
	}
	MissionUpdate(id, galo, xp, money)
}

func MissionUpdate(id disgord.Snowflake, galo Galo, xp int, money int) {
	if galo.Xp != xp {
		ChangeMoney(id, money, 0)
		UpdateGaloDB(id, map[string]interface{}{
			"xp":       galo.Xp + xp,
			"missions": galo.Missions,
		})
	}
}

func CreateMission(galo Galo) Mission {
	missionType := MissionType(rand.Intn(4))
	level := rand.Intn(3)
	galoAdv := 0
	if missionType == FightGalo || missionType == WinGalo{
		galoAdv = GetRand()
	}
	return Mission{
		Type:  missionType,
		Level: level,
		Adv: galoAdv,
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
