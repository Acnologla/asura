package rinha

const BASE_TOWER_XP = 60
const BASE_TOWER_MONEY = 30
const FLOOR_XP_BONUS = 18
const FLOOR_MONEY_BONUS = 8
const MAXIMUM_FLOOR = 150
const REWARD_FLOOR_DIVIDER = 5

func CalcTowerReward(floor int) (int, int) {
	if floor%REWARD_FLOOR_DIVIDER != 0 || floor == 0 {
		return 0, 0
	}

	f := floor / REWARD_FLOOR_DIVIDER
	xp := BASE_TOWER_XP + (f * FLOOR_XP_BONUS)
	money := BASE_TOWER_MONEY + (f * FLOOR_MONEY_BONUS)

	return xp, money
}

var rewards = map[int]int{
	25:  2,
	50:  2,
	75:  3,
	100: 3,
	150: 4,
}

func CalcFloorReward(floor int) int {
	if reward, exists := rewards[floor]; exists {
		return reward
	}
	return -1
}

func GetFloorRarity(floor int) Rarity {
	switch {
	case floor >= 100:
		return Mythic
	case floor >= 75:
		return Legendary
	case floor >= 50:
		return Epic
	case floor >= 25:
		return Rare
	default:
		return Common
	}
}
