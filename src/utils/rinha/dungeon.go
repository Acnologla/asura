package rinha

import (
	"math/rand"
)

type Room struct {
	Boss  Galo `json:"boss"`
	Level int  `json:"level"`
}

func Diff(galo Galo, diffGalo Galo) map[string]interface{} {
	diff := map[string] interface{}{}
	if galo.CommonLootbox != diffGalo.CommonLootbox{
		diff["commonLootbox"] = diffGalo.CommonLootbox
	}
	if galo.Lootbox != diffGalo.Lootbox{
		diff["lootbox"] = diffGalo.Lootbox
	}
	if galo.RareLootbox != diffGalo.RareLootbox{
		diff["rareLootbox"] = diffGalo.RareLootbox
	}
	if len(galo.Items) != len(diffGalo.Items){
		diff["items"] = diffGalo.Items
	}
	return diff
}

func DungeonWin(level int, galo Galo) (Galo, string){
	value := rand.Intn(101)
	msg := ""
	if level == 0{
		if 3 >= value {
			galo.Items, msg = AddItem(1, galo.Items)
		}else if value >= 52{
			galo.Items, msg = AddItem(0, galo.Items)
		}else{
			galo.CommonLootbox++
			msg = "Ganhou uma caixa comum"
		}
	}else if level == 1{
		if 11 >= value{
			galo.Items, msg = AddItem(1, galo.Items)
		}else if 55 >= value{
			galo.Items, msg = AddItem(0, galo.Items)
		}else{
			galo.CommonLootbox++
			msg = "Ganhou uma caixa comum"
		}
	}else if level == 2{	
		if 3 >= value {
			galo.Items, msg = AddItem(2, galo.Items)
		}else if value >= 52{
			galo.Items, msg = AddItem(1, galo.Items)
		}else{
			galo.Lootbox++
			msg = "Ganhou uma caixa Normal"
		}
	}else if level == 3 {
		if 11 >= value{
			galo.Items, msg = AddItem(2, galo.Items)
		}else if 55 >= value{
			galo.Items, msg = AddItem(1, galo.Items)
		}else{
			galo.Lootbox++
			msg = "Ganhou uma caixa normal"
		}
	}else if level == 4{
		if 49 >= value{
			galo.Items, msg = AddItem(2, galo.Items)
		}else{
			galo.RareLootbox++
			msg = "Ganhou uma caixa rara"
		}
	}
	return galo, msg
}