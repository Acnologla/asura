package rinha

import (
	"asura/src/utils"
	"fmt"
)

type CosmeticType int

const (
	Background CosmeticType = iota
	Badge
	Skin
)

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

func GetCosmeticRandBgByType(rarity Rarity) int {
	cosmeticArr := []*Cosmetic{}
	for _, cosmetic := range Cosmetics {
		if cosmetic.Rarity == rarity && cosmetic.Type == Background {
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

func GetBadges(galo Galo) []*Cosmetic {
	badges := []*Cosmetic{}
	for _, _cosmetic := range galo.Cosmetics {
		cosmetic := Cosmetics[_cosmetic]
		if cosmetic.Type == Badge {
			badges = append(badges, cosmetic)
		}
	}
	return badges
}

func OpenCosmeticBg() int {
	value := utils.RandInt(1001)
	if 15 >= value {
		return GetCosmeticRandBgByType(Legendary)
	} else if 120 >= value {
		return GetCosmeticRandBgByType(Epic)
	} else if 460 >= value {
		return GetCosmeticRandBgByType(Rare)
	}
	return GetCosmeticRandBgByType(Common)
}

func OpenCosmetic(pity int) (int, bool) {
	value := utils.RandInt(1001)
	pitVal := int(CalcPity(pity) * 15)
	if 15+pitVal >= value {
		return GetCosmeticRandByType(Legendary), true
	} else if 120 >= value {
		return GetCosmeticRandByType(Epic), false
	} else if 460 >= value {
		return GetCosmeticRandByType(Rare), false
	}
	return GetCosmeticRandByType(Common), false
}

func GetRandCosmetic() int {
	return utils.RandInt(len(Cosmetics)-1) + 1
}

func CosmeticCommand(cosmetic Cosmetic) string {
	if cosmetic.Type == Background {
		return "Use j!background para equipar o background"
	}
	return ""
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
