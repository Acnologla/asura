package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"fmt"
)

type CosmeticType int

const (
	Background CosmeticType = iota
	Badge
	Skin
)

func GetBackground(galo *entities.User) string {
	var equipBG *Cosmetic
	if IsVip(galo) {
		if galo.VipBackground != "" {
			return galo.VipBackground
		}
	}
	for _, item := range galo.Items {
		if item.Equip && item.Type == entities.CosmeticType {
			c := Cosmetics[item.ItemID]
			if c.Type == Background {
				equipBG = c
				break
			}
		}
	}
	if equipBG != nil {
		return equipBG.Value
	}
	return "https://i.imgur.com/F64ybgg.jpg"
}

type Cosmetic struct {
	Type         CosmeticType `json:"type"`
	Name         string       `json:"name"`
	Value        string       `json:"value"`
	ReverseValue string       `json:"reverseValue"`
	Rarity       Rarity       `json:"rarity"`
	Extra        int          `json:"extra"`
}

func (cosmetic Cosmetic) TypeToString() string {
	return [...]string{"Background", "Badge", "Skin"}[cosmetic.Type]
}

func (cosmetic Cosmetic) String() string {
	cosmeticType := cosmetic.TypeToString()
	return fmt.Sprintf("(%s) - %s %s", cosmetic.Rarity.String(), cosmeticType, cosmetic.Name)
}

func GetRandCosmeticType() CosmeticType {
	val := utils.RandInt(101)
	if val > 70 {
		return Background
	}
	return Skin
}

func GetCosmeticRandByType(rarity Rarity) int {
	t := GetRandCosmeticType()
	cosmeticArr := []*Cosmetic{}
	for _, cosmetic := range Cosmetics {
		if cosmetic.Rarity == rarity && (cosmetic.Type == t) {
			cosmeticArr = append(cosmeticArr, cosmetic)
		}
	}
	selected := cosmeticArr[utils.RandInt(len(cosmeticArr))]
	for i, cosmetic := range Cosmetics {
		if cosmetic.Name == selected.Name {
			return i
		}
	}
	return -1

}

func GetGaloImage(galo *entities.Rooster, items []*entities.Item, def ...string) string {
	for _, item := range items {
		if item.Type == entities.CosmeticType && item.Equip {
			cosmetic := Cosmetics[item.ItemID]
			if cosmetic.Type == Skin && cosmetic.Extra == galo.Type {
				if len(def) > 0 {
					return cosmetic.ReverseValue
				}
				return cosmetic.Value
			}
		}
	}
	if len(def) > 0 {
		return Sprites[1][galo.Type-1]
	}
	return Sprites[0][galo.Type-1]
}

func SellCosmetic(cosmetic Cosmetic) int {
	if cosmetic.Rarity == Common {
		return 15
	}
	if cosmetic.Rarity == Rare {
		return 90
	}
	if cosmetic.Rarity == Epic {
		return 235
	}
	if cosmetic.Rarity == Legendary {
		return 500
	}
	return 0
}

func GetCosmeticsByTypes(Items []*entities.Item, cType CosmeticType) ([]*Cosmetic, []*entities.Item) {
	cosmetics := []*Cosmetic{}
	items := []*entities.Item{}
	for _, item := range Items {
		if item.Type == entities.CosmeticType {
			cosmetic := Cosmetics[item.ItemID]
			if cosmetic.Type == cType {
				cosmetics = append(cosmetics, cosmetic)
				items = append(items, item)
			}
		}
	}
	return cosmetics, items
}

func GetBadges(galo *entities.User) []*Cosmetic {
	badges := []*Cosmetic{}
	for _, item := range galo.Items {
		if item.Type == entities.CosmeticType {
			cosmetic := Cosmetics[item.ItemID]
			if cosmetic.Type == Badge {
				badges = append(badges, cosmetic)
			}
		}
	}
	return badges
}
