package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"fmt"
)

func GetShopRooster() int {
	classTypeArr := []*Class{}
	for _, class := range Classes {
		if class.Rarity < Special && class.Rarity != -1 {
			classTypeArr = append(classTypeArr, class)
		}
	}
	selected := classTypeArr[utils.RandInt(len(classTypeArr))]
	for i, class := range Classes {
		if class.Name == selected.Name {
			return i
		}
	}
	return -1
}

func GetShopRarity(t entities.ShopItemType, value int) int {
	switch t {
	case entities.Shards:
		return value
	case entities.Cosmetics:
		return int(Cosmetics[value].Rarity)
	case entities.Roosters:
		return int(Classes[value].Rarity)
	case entities.Items:
		return Items[value].Level
	}

	return 1
}

func GetShopItemRarity(item *entities.ShopItem) string {
	switch item.Type {
	case entities.Shards:
		return Rarity(item.Value).String()
	case entities.Cosmetics:
		return Cosmetics[item.Value].Rarity.String()
	case entities.Roosters:
		return Classes[item.Value].Rarity.String()
	case entities.Items:
		return LevelToString(Items[item.Value].Level)
	case entities.AsuraCoin:
		return "Lendario"
	case entities.Xp:
		return "Comum"
	}

	return ""
}

func GetShopItemName(item *entities.ShopItem) string {
	switch item.Type {
	case entities.Shards:
		return fmt.Sprintf("Shard %s", Rarity(item.Value).String())
	case entities.Cosmetics:
		return fmt.Sprintf("Cosmético %s", Cosmetics[item.Value].Name)
	case entities.Roosters:
		return fmt.Sprintf("Galo %s", Classes[item.Value].Name)
	case entities.Items:
		return fmt.Sprintf("Item %s", Items[item.Value].Name)
	case entities.AsuraCoin:
		return fmt.Sprintf("%d AsuraCoins", item.Value)
	case entities.Xp:
		return fmt.Sprintf("%d XP", item.Value)
	}

	return ""
}

func GetShopItem() int {
	itemTypearr := []*Item{}
	for _, item := range Items {
		if item.Level < 5 {
			itemTypearr = append(itemTypearr, item)
		}
	}
	selected := itemTypearr[utils.RandInt(len(itemTypearr))]
	for i, item := range Items {
		if item.Name == selected.Name {
			return i
		}
	}
	return -1
}

func GetShopCosmetic() int {
	cosmeticTypeArr := []*Cosmetic{}
	for _, cosmetic := range Cosmetics {
		if cosmetic.Rarity < Special && (cosmetic.Type == Background || cosmetic.Type == Skin) {
			cosmeticTypeArr = append(cosmeticTypeArr, cosmetic)
		}
	}
	selected := cosmeticTypeArr[utils.RandInt(len(cosmeticTypeArr))]
	for i, cosmetic := range Cosmetics {
		if cosmetic.Name == selected.Name {
			return i
		}
	}
	return -1
}

func GetShopShard() int {
	rand := utils.RandInt(100)
	if 45 > rand {
		return 2
	} else if 20 > rand {
		return 3
	} else if 5 > rand {
		return 5
	}

	return 1
}

var shopItemValueGenerator = map[entities.ShopItemType]func() int{
	entities.Shards:    GetShopShard,
	entities.Cosmetics: GetShopCosmetic,
	entities.Roosters:  GetShopRooster,
	entities.Items:     GetShopItem,
	entities.AsuraCoin: func() int { return utils.RandInt(3) + 1 },
	entities.Xp:        func() int { return utils.RandInt(901) + 100 },
}

func GenerateShopItem() *entities.ShopItem {
	rand := utils.RandInt(5)
	randDiscount := utils.RandInt(100)
	discount := 1.0
	if randDiscount < 9 {
		discount = 0.8
	}
	t := entities.ShopItemType(rand)
	value := shopItemValueGenerator[t]()
	return &entities.ShopItem{
		Type:     t,
		Discount: float64(discount),
		Value:    value,
		Rarity:   GetShopRarity(t, value),
	}
}

func GenerateShop() *entities.Shop {
	return &entities.Shop{
		GenerateShopItem(),
		GenerateShopItem(),
		GenerateShopItem(),
		GenerateShopItem(),
		GenerateShopItem(),
	}
}
