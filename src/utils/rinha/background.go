package rinha

func GetBackground(galo Galo) string {
	if IsVip(galo) {
		return "./resources/wallVip.jpg"
	}
	return "./resources/wall.jpg"
}
