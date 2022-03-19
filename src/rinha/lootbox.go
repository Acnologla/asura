package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"fmt"

	"strconv"

	"github.com/google/uuid"
)

type Lootboxes struct {
	Common    int
	Rare      int
	Epic      int
	Legendary int
	Items     int
	Cosmetic  int
	Normal    int
}

var Prices = [...][]int{{100, 0}, {400, 0}, {800, 0}, {2200, 0}, {0, 2}, {0, 3}, {320, 0}}
var LootNames = [...]string{"comum", "normal", "rara", "epica", "lendaria", "items", "cosmetica"}

func GenerateLootPrices() (text string) {
	for i, name := range LootNames {
		price := Prices[i][0]
		priceText := strconv.Itoa(price)
		if price == 0 {
			priceText = fmt.Sprintf("%d AsuraCoins", Prices[i][1])
		}
		text += fmt.Sprintf("[%s] Lootbox %s\n", priceText, name)
	}
	return
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
	pitVal := int(CalcPity(pity) * 140)
	if 140+pitVal >= value {
		return GetRandByType(Epic), true
	} else if 600 >= value {
		return GetRandByType(Rare), false
	}
	return GetRandByType(Common), false
}

func OpenNormal(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 6)
	if 6+pitVal >= value {
		return GetRandByType(Epic), true
	} else if 31 >= value {
		return GetRandByType(Rare), false
	}
	return GetRandByType(Common), false
}

func OpenCommon(pity int) (int, bool) {
	value := utils.RandInt(101)
	pitVal := int(CalcPity(pity) * 6)
	if 6+pitVal >= value {
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

func GetPrice(lootType int) (int, int) {
	return Prices[lootType][0], Prices[lootType][1]
}

func OpenCosmetic(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 15)
	if 8+pitVal >= value {
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
