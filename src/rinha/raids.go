package rinha

import "asura/src/utils"

const KEY_CHANCE = 9

func DropKey() bool {
	rand := utils.RandInt(1000)
	return rand < KEY_CHANCE
}

func GetKeyRarity() Rarity {
	rand := utils.RandInt(1001)
	if 10 > rand {
		return Mythic
	} else if 55 > rand {
		return Legendary
	} else if 190 > rand {
		return Epic
	} else if 390 > rand {
		return Rare
	}
	return Common

}

// xp multiplier atribbutes multiplier
func GetMultipliers(rarity Rarity) (int, int) {
	if rarity == Epic {
		return 3, 2
	}
	if rarity == Legendary {
		return 5, 4
	}
	if rarity == Mythic {
		return 10, 7
	}
	return 1, 1
}

// xp money
func CalcRaidBattleWinPrize(rarity Rarity) (int, int) {
	xp := 60
	money := 30
	multiplier, _ := GetMultipliers(rarity)
	return xp * multiplier, money * multiplier
}

func CalcMaxRaidBattles(rarity Rarity) int {
	return int(rarity) + 5
}
