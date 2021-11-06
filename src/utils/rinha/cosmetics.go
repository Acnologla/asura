package rinha

import (
	"asura/src/utils"
	"fmt"
)

type CosmeticType int

const (
	Background CosmeticType = iota
)

type Cosmetic struct {
	Type   CosmeticType `json:"type"`
	Name   string       `json:"name"`
	Value  string       `json:"value"`
	Rarity Rarity       `json:"rarity"`
}

func (cosmetic Cosmetic) TypeToString() string {
	return [...]string{"Background"}[cosmetic.Type]
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
	pitVal := CalcPity(pity) * 15
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
