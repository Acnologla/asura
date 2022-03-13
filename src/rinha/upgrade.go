package rinha

import (
	"asura/src/entities"
	"fmt"
	"math"
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

func CalcUserXp(galo *entities.User) int {
	n := len(galo.Upgrades)
	if n == 0 {
		return 300
	}
	return 300 * (6 * int(math.Pow(float64(n), 2)))
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
