package rinha

type ArenaResult = int

const (
	TimeExceeded ArenaResult = iota
	ArenaWin
	ArenaLose
	ArenaTie
)

func CalcArena(gal Galo) (Galo, int, int) {
	wins := gal.Arena.Win
	money := 65 * wins
	xp := 20 * wins
	gal.Xp += xp
	gal.Money += money
	gal.Arena.Win = 0
	gal.Arena.Lose = 0
	gal.Arena.Active = false
	gal.Arena.LastFight = 0
	return gal, xp, money
}
