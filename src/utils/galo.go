package utils

import (
	"math"
)

func CalcLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp) / 30)))
}