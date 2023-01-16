package rinha

import (
	"asura/src/entities"
	"fmt"
	"math"
)

type BattlePassType int

const (
	BattlePassItem BattlePassType = iota
	BattlePassMoney
	BattlePassXp
	BattlePassCoins
)

type BattlePassLevel struct {
	Type     BattlePassType    `json:"type"`
	ItemType entities.ItemType `json:"itemType"`
	Value    int               `json:"value"`
}

func CalcBPLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp)/50))) + 1
}

func CalcBPXP(level int) int {
	return int(math.Pow(float64(level-1), 2)) * 50
}

func BPLevelToString(level int) string {
	currentLevel := BattlePass[level]
	switch currentLevel.Type {
	case BattlePassItem:
		if currentLevel.ItemType == entities.CosmeticType {
			return fmt.Sprintf("Cosmetico %s", Cosmetics[currentLevel.Value].Name)
		} else {
			return fmt.Sprintf("Lootbox %s", LootNames[currentLevel.Value])
		}
	case BattlePassMoney:
		return fmt.Sprintf("%d Moedas", currentLevel.Value)
	case BattlePassXp:
		return fmt.Sprintf("%d XP", currentLevel.Value)
	case BattlePassCoins:
		return fmt.Sprintf("%d AsuraCoins", currentLevel.Value)
	}
	return ""
}
