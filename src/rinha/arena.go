package rinha

import "asura/src/entities"

func CalcArena(user *entities.User) (int, int) {
	wins := user.ArenaWin
	money := 60 * wins
	xp := 25 * wins
	user.Money += money
	user.ArenaWin = 0
	user.ArenaLose = 0
	user.ArenaActive = false
	user.ArenaLastFight = 0
	return xp, money
}
