package rinha

import "fmt"

type CosmesticType int

const (
	Background CosmesticType = iota
)

type Cosmestic struct {
	Type   CosmesticType
	Name   string
	Value  string
	Rarity Rarity
}

func (cosmetic Cosmestic) TypeToString() string {
	return [...]string{"Background"}[cosmetic.Type]
}

func (cosmetic Cosmestic) String() string {
	cosmeticType := cosmetic.TypeToString()
	return fmt.Sprintf("(%s) - %s %s", cosmetic.Rarity.String(), cosmeticType, cosmetic.Name)
}
