package rinha

func GetBackground(galo Galo) string {
	bgs, indexes := GetBackgrounds(galo.Cosmetics)
	if IsVip(galo) {
		if galo.VipBackground != "" {
			return galo.VipBackground
		}
	}
	if len(bgs) != 0 {
		if galo.Background <= indexes[len(indexes)-1] {
			return Cosmetics[galo.Cosmetics[galo.Background]].Value
		}
	}
	return "https://i.imgur.com/F64ybgg.jpg"
}

func GetBackgrounds(cosmetics []int) ([]Cosmetic, []int) {
	arr := []Cosmetic{}
	indexes := []int{}
	for i, cosmetic := range cosmetics {
		if Cosmetics[cosmetic].Type == Background {
			arr = append(arr, *Cosmetics[cosmetic])
			indexes = append(indexes, i)
		}
	}
	return arr, indexes
}
