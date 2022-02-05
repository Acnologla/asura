package rinha

import (
	"asura/src/entities"
	"fmt"

	"strconv"
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
