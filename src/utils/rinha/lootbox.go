package rinha

import (
	"asura/src/utils"
)

/* var lootChances = [][3]int{
	{50},
	{240, 50},
	{450, 100},
}

func _Open(lootType int) int {
	lootChance := lootChances[lootType]
	value := utils.RandInt(1000) + 1
	for i := len(lootChance) - 1; i >= 0; i-- {
		randValue := lootChance[i]
		if randValue >= value && value != 0 {
			return GetRandByType(Rarity(i + 1))
		}
	}
	return GetRandByType(Common)
}
*/

const PityMultiplier = 1

func CalcPity(pity int) float64 {
	return (float64(pity) * PityMultiplier) / 100
}

func OpenEpic(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 6)
	if 6+pitVal >= value {
		return GetRandByType(Legendary), true
	} else if 250 >= value {
		return GetRandByType(Epic), false
	}
	return GetRandByType(Rare), false
}

func OpenRare(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 130)
	if 130+pitVal >= value {
		return GetRandByType(Epic), true
	} else if 570 >= value {
		return GetRandByType(Rare), false
	}
	return GetRandByType(Common), false
}

func OpenNormal(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 6)
	if 6+pitVal >= value {
		return GetRandByType(Epic), true
	} else if 30 >= value {
		return GetRandByType(Rare), false
	}
	return GetRandByType(Common), false
}

func OpenCommon(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 5)
	if 5+pitVal >= value {
		return GetRandByType(Rare), true
	}
	return GetRandByType(Common), false
}

func OpenLegendary(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 4)
	if 4+pitVal >= value {
		return GetRandByType(Mythic), true
	} else if 50 >= value {
		return GetRandByType(Legendary), false
	}
	return GetRandByType(Epic), false
}

func OpenItems(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 10)
	if 10+pitVal >= value {
		return GetItemByLevel(4), true
	}
	return GetItemByLevel(3), false
}

func Open(lootType string, galo Galo) (int, int) {
	var gal int
	var isRarest bool
	if lootType == "comum" {
		gal, isRarest = OpenCommon(galo.Pity)
	}
	if lootType == "rara" {
		gal, isRarest = OpenRare(galo.Pity)
	}
	if lootType == "epica" {
		gal, isRarest = OpenEpic(galo.Pity)
	}
	if lootType == "cosmetica" {
		gal, isRarest = OpenCosmetic(galo.Pity)
	}
	if lootType == "items" {
		gal, isRarest = OpenItems(galo.Pity)
	}
	if lootType == "lendaria" {
		gal, isRarest = OpenLegendary(galo.Pity)
	}
	if lootType == "normal" {
		gal, isRarest = OpenNormal(galo.Pity)
	}
	if isRarest {
		return gal, 0
	}
	price, asuraCoins := GetPrice(lootType)
	money := price
	if asuraCoins > 0 {
		money = asuraCoins * 1800
	}
	newPity := money / 100
	return gal, newPity + galo.Pity
}

func GetPrice(lootType string) (int, int) {
	if lootType == "comum" {
		return 100, 0
	}
	if lootType == "rara" {
		return 800, 0
	}
	if lootType == "epica" {
		return 2200, 0
	}
	if lootType == "cosmetica" {
		return 300, 0
	}
	if lootType == "items" {
		return 0, 4
	}
	if lootType == "lendaria" {
		return 0, 2
	}
	return 400, 0
}

func HaveLootbox(galo Galo, lootbox string) bool {
	if lootbox == "comum" {
		return galo.CommonLootbox > 0
	}
	if lootbox == "rara" {
		return galo.RareLootbox > 0
	}
	if lootbox == "cosmetica" {
		return galo.CosmeticLootbox > 0
	}
	if lootbox == "epica" {
		return galo.EpicLootbox > 0
	}
	if lootbox == "lendaria" {
		return galo.LegendaryLootbox > 0
	}
	if lootbox == "items" {
		return galo.ItemsLootbox > 0
	}
	return galo.Lootbox > 0
}

func GetNewLb(lootbox string, galo Galo, add bool) Galo {
	if lootbox == "comum" {
		if add {
			galo.CommonLootbox++
		} else {
			galo.CommonLootbox--
		}
	} else if lootbox == "rara" {
		if add {
			galo.RareLootbox++
		} else {
			galo.RareLootbox--
		}
	} else if lootbox == "epica" {
		if add {
			galo.EpicLootbox++
		} else {
			galo.EpicLootbox--
		}
	} else if lootbox == "cosmetica" {
		if add {
			galo.CosmeticLootbox++
		} else {
			galo.CosmeticLootbox--
		}
	} else if lootbox == "lendaria" {
		if add {
			galo.LegendaryLootbox++
		} else {
			galo.LegendaryLootbox--
		}
	} else if lootbox == "items" {
		if add {
			galo.ItemsLootbox++
		} else {
			galo.ItemsLootbox--
		}
	} else {
		if add {
			galo.Lootbox++
		} else {
			galo.Lootbox--
		}
	}
	return galo
}

func Sell(rarity Rarity, xp int, reset int) (int, int) {
	level := float64(CalcLevel(xp)+(reset*30)) - 1
	price := float64(rarity.Price())
	if reset == 0 {
		return int(price * (level/5 + 1)), 0
	}
	asuraCoins := reset
	if rarity > Rare {
		asuraCoins++
	}
	if rarity == Legendary {
		asuraCoins += 2
	}
	return 0, asuraCoins
}
