package rinha

import (
	"asura/src/entities"
	"fmt"
)

type Rank struct {
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Emoji string `json:"emoji"`
}

func (rank Rank) ToString() string {
	return fmt.Sprintf("[%s] **%s**", rank.Emoji, rank.Name)
}

var unkranked = Rank{
	Name:  "Unranked",
	Rank:  0,
	Emoji: "",
}

func GetRank(u *entities.User) *Rank {
	i := 0
	for j, rank := range Ranks {
		if rank.Rank > u.Rank {
			i = j - 1
			break
		}
		if j == len(Ranks)-1 {
			i = j
		}
	}

	if i == -1 {
		return &unkranked
	}

	return Ranks[i]
}

const BASE_RANK = 15
const MIN_RANK = 5
const MAX_RANK = 25

func CalcRank(winnerRank, loserRank int) int {
	diff := float64(winnerRank - loserRank)
	v := int(BASE_RANK * (1.0 - diff/400))
	if v < MIN_RANK {
		return MIN_RANK
	}
	if v > MAX_RANK {
		return MAX_RANK
	}
	return v
}

func CalcRankedPrize(rank *Rank) (money int, xp int, asuraCoins int) {
	money = int(float64(rank.Rank) * 0.85)
	xp = int(float64(rank.Rank) * 1.8)

	if rank.Rank > 2500 {
		money = 0
		asuraCoins = rank.Rank / 1000
	}

	return
}
