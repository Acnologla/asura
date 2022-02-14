package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
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

func GetEquippedItem(user *entities.User) int {
	for _, item := range user.Items {
		if item.Type == entities.NormalType && item.Equip {
			return item.ItemID
		}
	}
	return -1
}
func GetItem(user *entities.User) *Item {
	itemID := GetEquippedItem(user)
	if itemID == -1 {
		return nil
	}
	return Items[itemID]
}

func SellItem(item Item) int {
	if item.Level == 0 {
		return 50
	}
	if item.Level == 1 {
		return 170
	}
	if item.Level == 2 {
		return 390
	}
	if item.Level == 3 {
		return 1000
	}
	if item.Level == 4 {
		return 2000
	}
	return 0
}
