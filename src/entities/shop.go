package entities

import (
	"asura/src/utils"
	"math"
	"sync"

	"github.com/andersfylling/disgord"
)

type ShopItemType int

const (
	Shards ShopItemType = iota
	Cosmetics
	Roosters
	Items
	AsuraCoin
	Xp
)

type ShopItem struct {
	Type     ShopItemType
	Value    int
	Rarity   int
	Discount float64
	Users    []disgord.Snowflake
	sync.Mutex
}

func (item *ShopItem) price() (int, int) {
	switch item.Type {
	case Shards:
		if item.Rarity == 5 {
			return 0, 2
		} else if item.Rarity == 3 {
			return 14500, 0
		} else if item.Rarity == 2 {
			return 4600, 0
		}
		return 950, 0
	case Cosmetics:
		if item.Rarity == 3 {
			return 0, 2
		}
		return (item.Rarity + 1) * 4200, 0
	case Roosters:
		if item.Rarity == 3 {
			return 0, 3
		}
		if item.Rarity == 2 {
			return 5900, 0
		}
		if item.Rarity == 1 {
			return 800, 0
		}
		return 200, 0
	case Items:
		if item.Rarity == 4 {
			return 0, 4
		}
		if item.Rarity == 3 {
			return 10500, 0
		}
		return ((item.Rarity * item.Rarity * 2) + 1) * 420, 0
	case AsuraCoin:
		return 18000 * item.Value, 0
	case Xp:
		return item.Value * 2, 0
	}

	return 0, 0

}

func (item *ShopItem) OriginalPrice() (int, int) {
	return item.price()
}

func (item *ShopItem) Price() (int, int) {
	moneyPrice, asuraCoinPrice := item.price()
	return int(math.Round(float64(moneyPrice) * item.Discount)), int(math.Round((float64(asuraCoinPrice) * item.Discount)))
}

func (item *ShopItem) CanBuy(user disgord.Snowflake) bool {
	item.Lock()
	defer item.Unlock()
	return !utils.Has(item.Users, user)
}

func (item *ShopItem) Buy(user disgord.Snowflake) {
	item.Lock()
	defer item.Unlock()
	item.Users = append(item.Users, user)
}

type Shop = [5]*ShopItem
