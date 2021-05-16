package rinha

func CalcArena(gal Galo) (Galo, int, int ){
	wins := gal.Arena.Win
	money := 75 * wins
	xp := 22 * wins
	gal.Xp += xp
	gal.Money += money
	gal.Arena.Win = 0
	gal.Arena.Lose = 0
	gal.Arena.Active = false
	return gal,xp, money
}