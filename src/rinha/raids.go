package rinha

import "asura/src/utils"

const KEY_CHANCE = 8
const NEWBIE_ADD = 6
const VIP_ADD = 1

func DropKey(userXP int, vip bool, clanLevel int) bool {
	rand := utils.RandInt(1001)
	add := 0
	if 900 > userXP {
		add += NEWBIE_ADD
	}

	if vip {
		add += VIP_ADD
	}

	if clanLevel >= 15 {
		add++
	}

	return rand < (KEY_CHANCE + add)
}

func GetKeyRarity() Rarity {
	rand := utils.RandInt(1001)
	if rand == 0 {
		return God
	} else if 5 > rand {
		return Mythic
	} else if 42 > rand {
		return Legendary
	} else if 197 > rand {
		return Epic
	} else if 445 > rand {
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

	if rarity == God {
		return 23, 22
	}

	return 1, 1
}

// xp money
func CalcRaidBattleWinPrize(rarity Rarity) (int, int) {
	xp := 75
	money := 25
	multiplier, _ := GetMultipliers(rarity)
	return xp * multiplier, money * multiplier
}

func CalcMaxRaidBattles(rarity Rarity) int {
	return int(rarity) + 5
}
