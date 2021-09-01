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
func OpenEpic() int {
	value := utils.RandInt(1001)
	if 10 >= value {
		return GetRandByType(Legendary)
	} else if 243 >= value {
		return GetRandByType(Epic)
	}
	return GetRandByType(Rare)
}

func OpenRare() int {
	value := utils.RandInt(1001)
	if 120 >= value {
		return GetRandByType(Epic)
	} else if 500 >= value {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func OpenNormal() int {
	value := utils.RandInt(101)
	if 5 >= value {
		return GetRandByType(Epic)
	} else if 25 >= value {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func OpenCommon() int {
	value := utils.RandInt(101)
	if 4 >= value {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func OpenLegendary() int {
	value := utils.RandInt(101)
	if 4 >= value {
		return GetRandByType(Legendary)
	}
	return GetRandByType(Epic)
}

func OpenItems() int {
	value := utils.RandInt(101)
	if 10 >= value {
		return GetItemByLevel(4)
	}
	return GetItemByLevel(3)
}

func Open(lootType string) int {
	if lootType == "comum" {
		return OpenCommon()
	}
	if lootType == "rara" {
		return OpenRare()
	}
	if lootType == "epica" {
		return OpenEpic()
	}
	if lootType == "cosmetica" {
		return OpenCosmetic()
	}
	if lootType == "items" {
		return OpenItems()
	}
	if lootType == "lendaria" {
		return OpenLegendary()
	}
	return OpenNormal()
}

func GetPrice(lootType string) (int, int) {
	if lootType == "comum" {
		return 100, 0
	}
	if lootType == "rara" {
		return 800, 0
	}
	if lootType == "epica" {
		return 1750, 0
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
	return 0, asuraCoins
}
