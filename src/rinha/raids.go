package rinha

import "asura/src/utils"

const KEY_CHANCE = 6
const NEWBIE_ADD = 5
const VIP_ADD = 1

func DropKey(userXP int, vip bool) bool {
	rand := utils.RandInt(1001)

	if 650 > userXP {
		return rand < (KEY_CHANCE + NEWBIE_ADD)
	}

	if vip {
		return rand < (KEY_CHANCE + VIP_ADD)
	}

	return rand < KEY_CHANCE
}

func GetKeyRarity() Rarity {
	rand := utils.RandInt(1001)
	if 4 > rand {
		return Mythic
	} else if 40 > rand {
		return Legendary
	} else if 195 > rand {
		return Epic
	} else if 400 > rand {
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
		return 6, 5
	}

	if rarity == Mythic {
		return 11, 12
	}

	return 1, 1
}

// xp money
func CalcRaidBattleWinPrize(rarity Rarity) (int, int) {
	xp := 70
	money := 25
	multiplier, _ := GetMultipliers(rarity)
	return xp * multiplier, money * multiplier
}

func CalcMaxRaidBattles(rarity Rarity) int {
	return int(rarity) + 5
}
