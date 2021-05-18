package rinha

import (
	"asura/src/utils"
)

var lootChances = [][3]int{
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

func OpenEpic() int {
	value := utils.RandInt(1001)
	if 14 >= value {
		return GetRandByType(Legendary)
	} else if 241 >= value {
		return GetRandByType(Epic)
	}
	return GetRandByType(Rare)
}

func OpenRare() int {
	value := utils.RandInt(1001)
	if 6 >= value {
		return GetRandByType(Legendary)
	} else if 120 >= value {
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
	return OpenNormal()
}

func GetPrice(lootType string) int {
	if lootType == "comum" {
		return 100
	}
	if lootType == "rara" {
		return 800
	}
	if lootType == "epica" {
		return 1800
	}
	if lootType == "cosmetica" {
		return 320
	}
	return 400
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
	} else {
		if add {
			galo.Lootbox++
		} else {
			galo.Lootbox--
		}
	}
	return galo
}

func Sell(rarity Rarity, xp int) int {
	level := CalcLevel(xp)
	return rarity.Price() * (level/5 + 1)
}
