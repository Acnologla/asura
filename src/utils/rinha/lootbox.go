package rinha

import (
	"math/rand"
)

func Open() int {
	value := rand.Intn(101)
	if 4 >= value {
		return GetRandByType(Epic)
	} else if 25 >= value {
		return GetRandByType(Rare)
	} else {
		return GetRandByType(Common)
	}
}