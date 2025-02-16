package rinha

import (
	"asura/src/entities"
	"fmt"
)

type Upgrade struct {
	Name   string    `json:"name"`
	Value  string    `json:"value"`
	Childs []Upgrade `json:"childs"`
}

func HasUpgrade(upgrades []int, upgradeList ...int) bool {
	if len(upgradeList) > len(upgrades) {
		return false
	}
	isFalse := true
	for i, upgrade := range upgradeList {
		if upgrades[i] != upgrade {
			isFalse = false
			break
		}
	}
	return isFalse
}

func CalcUserXp(user *entities.User) int {
	n := len(user.Upgrades)
	switch n {
	case 0:
		return 200
	case 1:
		return 1500
	case 2:
		return 8000
	case 3:
		return 20000
	case 4:
		return 50000
	default:
		return 100000
	}
}

func GetCurrentUpgrade(galo *entities.User) Upgrade {
	selected := Upgrade{
		Childs: Upgrades,
	}
	for _, upgrade := range galo.Upgrades {
		if selected.Name != "" {
			selected = selected.Childs[upgrade]
		} else {
			selected = Upgrades[upgrade]
		}
	}
	return selected
}

func UpgradesToString(galo *entities.User) (text string) {
	selected := GetCurrentUpgrade(galo)
	for i, upgrade := range selected.Childs {
		text += fmt.Sprintf("[%d] - %s\n%s\n", i, upgrade.Name, upgrade.Value)
	}
	return
}

func HavePoint(galo *entities.User) bool {
	return galo.UserXp >= CalcUserXp(galo)
}
