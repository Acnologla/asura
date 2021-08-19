package rinha

import (
	"asura/src/utils"
)

type Room struct {
	Boss  Galo `json:"boss"`
	Level int  `json:"level"`
}

func DungeonWin(level int, galo Galo) (Galo, string) {
	value := utils.RandInt(101)
	msg := ""
	if level == 0 {
		if 3 >= value {
			galo.Items, msg = AddItem(1, galo.Items)
		} else if value >= 52 {
			galo.Items, msg = AddItem(0, galo.Items)
		} else {
			galo.CommonLootbox++
			msg = "Ganhou uma caixa comum"
		}
	} else if level == 1 {
		if 11 >= value {
			galo.Items, msg = AddItem(1, galo.Items)
		} else if 55 >= value {
			galo.Items, msg = AddItem(0, galo.Items)
		} else {
			galo.CommonLootbox++
			msg = "Ganhou uma caixa comum"
		}
	} else if level == 2 {
		if 3 >= value {
			galo.Items, msg = AddItem(2, galo.Items)
		} else if value >= 52 {
			galo.Items, msg = AddItem(1, galo.Items)
		} else {
			galo.Lootbox++
			msg = "Ganhou uma caixa Normal"
		}
	} else if level == 3 {
		if 11 >= value {
			galo.Items, msg = AddItem(2, galo.Items)
		} else if 55 >= value {
			galo.Items, msg = AddItem(1, galo.Items)
		} else {
			galo.Lootbox++
			msg = "Ganhou uma caixa normal"
		}
	} else if level == 4 {
		if 49 >= value {
			galo.Items, msg = AddItem(2, galo.Items)
		} else {
			galo.RareLootbox++
			msg = "Ganhou uma caixa rara"
		}
	} else if level == 5 {
		if 4 >= value {
			galo.Items, msg = AddItem(3, galo.Items)
		} else if 49 >= value {
			galo.Items, msg = AddItem(2, galo.Items)
		} else {
			galo.RareLootbox++
			msg = "Ganhou uma caixa rara"
		}
	} else if level == 6 {
		if 49 >= value || galo.DungeonReset != 0 {
			galo.Items, msg = AddItem(3, galo.Items)
		} else {
			galo.EpicLootbox++
			msg = "Ganhou uma caixa epica"
		}
	}
	return galo, msg
}
