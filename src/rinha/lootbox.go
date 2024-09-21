package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"fmt"

	"strconv"

	"github.com/google/uuid"
)

type Lootboxes struct {
	Common      int
	Rare        int
	Epic        int
	Legendary   int
	Items       int
	Cosmetic    int
	Normal      int
	Mythic      int
	ItemsMythic int
}

var Prices = [...][]int{{100, 0}, {400, 0}, {800, 0}, {2200, 0}, {0, 3}, {0, 2}, {500, 0}, {-1, 10}, {-1, 20}}
var LootNames = [...]string{"comum", "normal", "rara", "epica", "lendaria", "items", "cosmetica", "mistica", "items mistica"}

func GenerateLootPrices() (text string) {
	for i, name := range LootNames {
		price := Prices[i][0]
		if price != -1 {
			priceText := strconv.Itoa(price)
			if price == 0 {
				priceText = fmt.Sprintf("%d AsuraCoins", Prices[i][1])
			}
			text += fmt.Sprintf("[%s] Lootbox %s\n", priceText, name)
		}
	}
	return
}

func MessageRandomLootbox() (Rarity, int) {
	rand := utils.RandInt(100)
	if rand < 22 {
		return Epic, 2
	}
	if rand < 48 {
		return Rare, 1
	}
	return Common, 0
}

func GetUserLootboxes(user *entities.User) (loot []string) {
	for _, item := range user.Items {
		if item.Type == entities.LootboxType {
			loot = append(loot, LootNames[item.ItemID])
		}
	}
	return
}

func GetLootboxes(user *entities.User) (lootboxes Lootboxes) {
	for _, item := range user.Items {
		if item.Type == entities.LootboxType {
			switch item.ItemID {
			case 0:
				lootboxes.Common += item.Quantity
			case 1:
				lootboxes.Normal += item.Quantity
			case 2:
				lootboxes.Rare += item.Quantity
			case 3:
				lootboxes.Epic += item.Quantity
			case 4:
				lootboxes.Legendary += item.Quantity
			case 5:
				lootboxes.Items += item.Quantity
			case 6:
				lootboxes.Cosmetic += item.Quantity
			case 7:
				lootboxes.Mythic += item.Quantity
			case 8:
				lootboxes.ItemsMythic += item.Quantity

			}
		}
	}
	return
}

func GetLbIndex(name string) int {
	for i, lbName := range LootNames {
		if name == lbName {
			return i
		}
	}
	return -1
}

func CalcPity(pity int) float64 {
	return (float64(pity) * PityMultiplier) / 100
}

func OpenEpic(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 7)
	if 7+pitVal >= value {
		return GetRandByType(Legendary), true
	} else if 300 >= value {
		return GetRandByType(Epic), false
	}
	return GetRandByType(Rare), false
}

func OpenRare(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 175)
	if 175+pitVal >= value {
		return GetRandByType(Epic), true
	} else if 640 >= value {
		return GetRandByType(Rare), false
	}
	return GetRandByType(Common), false
}

func OpenNormal(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 8)
	if 8+pitVal >= value {
		return GetRandByType(Epic), true
	} else if 35 >= value {
		return GetRandByType(Rare), false
	}
	return GetRandByType(Common), false
}

func OpenCommon(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 10)
	if 10+pitVal >= value {
		return GetRandByType(Rare), true
	}
	return GetRandByType(Common), false
}

func OpenLegendary(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 1.5)
	if 2+pitVal >= value {
		return GetRandByType(Mythic), true
	} else if 90 >= value {
		return GetRandByType(Legendary), false
	}
	return GetRandByType(Epic), false
}

func OpenMythic(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 30)
	if 30+pitVal >= value {
		return GetRandByType(Mythic), true
	}
	return GetRandByType(Legendary), false
}

func OpenItems(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 10)
	if 10+pitVal >= value {
		return GetItemByLevel(4), true
	}
	return GetItemByLevel(3), false
}

func OpenItemsMythic(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 1)
	if 1+pitVal >= value {
		return GetItemByLevel(6), true
	}
	return GetItemByLevel(4), false
}

func GetPrice(lootType int) (int, int) {
	return Prices[lootType][0], Prices[lootType][1]
}

func OpenCosmetic(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 7)
	if 7+pitVal >= value {
		return GetCosmeticRandByType(Legendary), true
	} else if 70 >= value {
		return GetCosmeticRandByType(Epic), false
	} else if 330 >= value {
		return GetCosmeticRandByType(Rare), false
	}
	return GetCosmeticRandByType(Common), false
}

func Open(lootType int, user *entities.User) (int, int) {
	var gal int
	var isRarest bool
	switch lootType {
	case 0:
		gal, isRarest = OpenCommon(user.Pity)
	case 1:
		gal, isRarest = OpenNormal(user.Pity)
	case 2:
		gal, isRarest = OpenRare(user.Pity)
	case 3:
		gal, isRarest = OpenEpic(user.Pity)
	case 4:
		gal, isRarest = OpenLegendary(user.Pity)
	case 5:
		gal, isRarest = OpenItems(user.Pity)
	case 6:
		gal, isRarest = OpenCosmetic(user.Pity)
	case 7:
		gal, isRarest = OpenMythic(user.Pity)
	case 8:
		gal, isRarest = OpenItemsMythic(user.Pity)
	}

	if isRarest {
		return gal, 0
	}
	price, asuraCoins := GetPrice(lootType)
	money := price
	if asuraCoins > 0 {
		money = asuraCoins * 1900
	}
	newPity := money / 80
	return gal, newPity + user.Pity
}

func GetLbID(items []*entities.Item, lootType int) (id uuid.UUID, exists bool) {
	for _, item := range items {
		if item.Type == entities.LootboxType && item.ItemID == lootType {
			id = item.ID
			exists = true
		}
	}
	return
}

func GetTrialLootbox(rarity Rarity) int {

	if rarity == Epic {
		rand := utils.RandInt(100)
		if 33 >= rand {
			return 3
		}
		return 2
	}

	if rarity == Legendary {
		rand := utils.RandInt(100)
		if 12 >= rand {
			return 4
		}
		return 3
	}

	if rarity == Mythic {
		rand := utils.RandInt(100)
		if rand == 0 {
			return 8
		} else if 2 >= rand {
			return 7
		}
		return 4
	}

	rand := utils.RandInt(100)
	if 38 >= rand {
		return 2
	}
	return 1

}
