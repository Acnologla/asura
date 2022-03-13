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
	Type   CosmeticType `json:"type"`
	Name   string       `json:"name"`
	Value  string       `json:"value"`
	Rarity Rarity       `json:"rarity"`
}

func (cosmetic Cosmetic) TypeToString() string {
	return [...]string{"Background", "Badge", "Skin"}[cosmetic.Type]
}

func (cosmetic Cosmetic) String() string {
	cosmeticType := cosmetic.TypeToString()
	return fmt.Sprintf("(%s) - %s %s", cosmetic.Rarity.String(), cosmeticType, cosmetic.Name)
}

func GetCosmeticRandByType(rarity Rarity) int {
	cosmeticArr := []*Cosmetic{}
	for _, cosmetic := range Cosmetics {
		if cosmetic.Rarity == rarity {
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
