package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
)

func HasEgg(u *entities.User) bool {
	return u.Egg >= 0
}

func ShardToPrice(rarity Rarity) int {
	switch rarity {
	case Rare:
		return 250
	case Epic:
		return 600
	case Legendary:
		return 1500
	case Mythic:
		return 3400
	}
	return 0
}

type EggChance struct {
	level     int
	mythic    int
	legendary int
	epic      int
	rare      int
}

var eggChances = []EggChance{
	{10, 0, 0, 0, 1000},
	{20, 0, 0, 0, 2500},
	{30, 0, 0, 1500, 8500},
	{40, 0, 700, 9300, 0},
	{50, 0, 1500, 8500, 0},
	{60, 0, 2000, 8000, 0},
	{70, 0, 2500, 7500, 0},
	{80, 0, 5000, 5000, 0},
	{90, 600, 9200, 0, 0},
	{100, 1200, 8800, 0, 0},
	{110, 1500, 8500, 0, 0},
	{120, 2500, 7500, 0, 0},
	{130, 3000, 7000, 0, 0},
	{140, 4000, 6000, 0, 0},
	{150, 5000, 5000, 0, 0},
	{160, 7000, 3000, 0, 0},
	{170, 10000, 0, 0, 0},
}

func GetRoosterFromEgg(eggXp int) Rarity {
	level := CalcLevel(eggXp)
	rand := utils.RandInt(10000)

	if rand == 0 {
		rand = 1
	}

	for _, chance := range eggChances {
		if level <= chance.level {
			if rand <= chance.mythic {
				return Mythic
			} else if rand <= chance.mythic+chance.legendary {
				return Legendary
			} else if rand <= chance.mythic+chance.legendary+chance.epic {
				return Epic
			} else if rand <= chance.mythic+chance.legendary+chance.epic+chance.rare {
				return Rare
			}
			return Common
		}
	}
	return Common
}
