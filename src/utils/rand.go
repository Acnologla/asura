package utils

import (
	"crypto/rand"
	"math/big"
	deterministic "math/rand"
)

func RandInt(max int) int{
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil{
		return deterministic.Intn(max)
	}
	return int(value.Int64())
}
