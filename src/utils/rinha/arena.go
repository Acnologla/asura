package rinha

func CalcArena(gal Galo) (Galo, int, int) {
	wins := gal.Arena.Win
	money := 80 * wins
	if money == 0 {
		money = 40
	}
	xp := 20 * wins
	gal.Xp += xp
	gal.Money += money
	gal.Arena.Win = 0
	gal.Arena.Lose = 0
	gal.Arena.Active = false
	return gal, xp, money
}
