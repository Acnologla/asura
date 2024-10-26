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
	if IsVip(galo) {
		if galo.VipBackground != "" {
			return galo.VipBackground
		}
	}
	bgs := []*Cosmetic{}
	for _, item := range galo.Items {
		if item.Equip && item.Type == entities.CosmeticType {
			c := Cosmetics[item.ItemID]
			if c.Type == Background {
				bgs = append(bgs, c)
			}
		}
	}

	if len(bgs) > 0 {
		bg := bgs[utils.RandInt(len(bgs))]
		return bg.Value
	}
	return "https://iili.io/JWj59Nj.png"
	//return "https://i.imgur.com/F64ybgg.jpg"
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
	if val > 65 {
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
	skins := []*Cosmetic{}
	for _, item := range items {
		if item.Type == entities.CosmeticType && item.Equip {
			cosmetic := Cosmetics[item.ItemID]
			if cosmetic.Type == Skin && cosmetic.Extra == galo.Type {
				skins = append(skins, cosmetic)

			}
		}
	}

	if len(skins) > 0 {
		cosmetic := skins[utils.RandInt(len(skins))]
		if len(def) > 0 {
			return cosmetic.ReverseValue
		}
		return cosmetic.Value
	}

	i := galo.Type - 1
	if i == -1 {
		i = 0
		fmt.Println(galo.UserID)
	}

	if len(def) > 0 {
		return Sprites[1][i]
	}

	return Sprites[0][i]
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

func FindRandomSkin(galo *entities.Rooster) entities.Item {
	skinArr := utils.Filter(Cosmetics, func(c *Cosmetic) bool {
		return c.Type == Skin && c.Extra == galo.Type
	})

	selected := skinArr[utils.RandInt(len(skinArr))]
	i := utils.GetIndex(Cosmetics, func(c *Cosmetic) bool {
		return c.Name == selected.Name
	})

	return entities.Item{
		ItemID: i,
		Type:   entities.CosmeticType,
		Equip:  true,
	}
}
