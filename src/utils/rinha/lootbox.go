package rinha

import (
	"math/rand"
)

func OpenRare() int {
	value := rand.Intn(1001)
	if 99 >= value {
		return GetRandByType(Epic)
	} else if 450 >= value {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func OpenNormal() int {
	value := rand.Intn(101)
	if 4 >= value {
		return GetRandByType(Epic)
	} else if 24 >= value {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func OpenCommon() int {
	value := rand.Intn(101)
	if 4 >= value {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func Open(lootType string) int{
	if lootType == "comum"{
		return OpenCommon()
	}
	if lootType == "rara"{
		return OpenRare()
	}
	return OpenNormal()
}


func GetPrice(lootType string) int{
	if lootType == "comum"{
		return 100
	}
	if lootType == "rara"{
		return 800
	}
	return 400
}

func HaveLootbox(galo Galo, lootbox string) bool{
	if lootbox == "comum"{
		return galo.CommonLootbox > 0 
	}
	if lootbox == "rara"{
		return galo.RareLootbox > 0
	}
	return galo.Lootbox > 0
}

func GetNewLb(lootbox string, galo Galo, add bool) (string, int){
	if lootbox == "comum"{
		if add{
			galo.CommonLootbox++
		}else{
			galo.CommonLootbox--
		}
		return "commonLootbox", galo.CommonLootbox 
	}
	if lootbox == "rara"{
		if add{
			galo.RareLootbox++
		}else{
			galo.RareLootbox--
		}
		return "rareLootbox", galo.RareLootbox 
	}
	if add{
		galo.Lootbox++
	}else{
		galo.Lootbox--
	}
	return "lootbox", galo.Lootbox
}


func Sell(rarity Rarity, xp int) int{
	level := CalcLevel(xp)
	return rarity.Price() * (level / 5 + 1)
}