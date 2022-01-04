package rinha

import (
	"asura/src/utils"
	"fmt"
	"math"
	"time"
)

type Item struct {
	Name    string  `json:"name"`
	Level   int     `json:"level"`
	Payload float64 `json:"payload"`
	Effect  int     `json:"effect"`
}

func GetItemByLevel(level int) int {
	itemTypearr := []*Item{}
	for _, item := range Items {
		if item.Level == level {
			itemTypearr = append(itemTypearr, item)
		}
	}
	selected := itemTypearr[utils.RandInt(len(itemTypearr))]
	for i, item := range Items {
		if item.Name == selected.Name {
			return i
		}
	}
	return -1
}

func LevelToString(level int) string {
	return map[int]string{
		0: "comum",
		1: "raro",
		2: "epico",
		3: "lendario",
		4: "especial",
		5: "evento",
	}[level]
}

func AddItem(level int, items []int) ([]int, string) {
	item := GetItemByLevel(level)
	if !IsIntInList(item, items) {
		items = append(items, item)
	}
	msg := fmt.Sprintf("Ganhou um item %s\n**%s** (use **j!items** para equipar seu item)", LevelToString(level), Items[item].Name)
	return items, msg
}

func CanTrade(galo Galo) bool {
	now := uint64(time.Now().Unix())
	return (now-galo.Cooldowns.TradeItem)/60/60/24 >= 3
}

func ItemToString(item *Item) string {
	if item.Effect == 1 {
		return fmt.Sprintf("Aumenta o dano dos seus ataques em **%d%%** ", int(math.Round((item.Payload-1)*100)))
	}
	if item.Effect == 2 {
		return fmt.Sprintf("Tem chance de causar **%s**", Effects[int(item.Payload)].Name)
	}
	if item.Effect == 3 {
		return fmt.Sprintf("Aumenta a probabilidade dos seus efeitos em **%d%%**", int(math.Round((item.Payload)*100)))
	}
	if item.Effect == 4 {
		return fmt.Sprintf("Aumenta sua vida maxima em **%d%%**", int(math.Round((item.Payload-1)*100)))
	}
	if item.Effect == 5 {
		return fmt.Sprintf("Bloqueia **%d%%** de dano", int(math.Round((1-item.Payload)*100)))
	}
	if item.Effect == 6 {
		return fmt.Sprintf("Aumenta o dano dos seus efeitos em **%d%%** ", int(math.Round((item.Payload-1)*100)))
	}
	if item.Effect == 7 {
		return fmt.Sprintf("Voce reflete todo o dano levado de ataques em **%d%%** ", int(item.Payload*100))
	}
	if item.Effect == 8 {
		return fmt.Sprintf("Voce ganha **%d**%% de xp adicional por train", int(item.Payload*100))
	}
	if item.Effect == 9 {
		return "Voce ganha 1 de ouro extra e 3 de xp adicional por train"
	}
	if item.Effect == 10 {
		return "Tem chance de ganhar 1 de xp extra para seu clan nos treinos"
	}
	return ""
}

func GetRandItem() int {
	return utils.RandInt(len(Items))
}

func GetItem(galo Galo) *Item {
	return Items[galo.Items[0]]
}

func LevelToPrice(item Item) int {
	if item.Level == 0 {
		return 80
	}
	if item.Level == 1 {
		return 200
	}
	if item.Level == 2 {
		return 500
	}
	if item.Level == 3 {
		return 1000
	}
	if item.Level == 4 {
		return 2000
	}
	return 0
}
