package rinha

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
