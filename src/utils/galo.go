package utils

import (
	"math"
)

func CalcLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp) / 30)))
}

func CalcXP(level int) int{
	return int(math.Pow(float64(level), 2)) * 30
}