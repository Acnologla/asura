package rinha

func GetBackground(galo Galo) string {
	if galo.Background != 0 {
		return Cosmetics[galo.Background].Value
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