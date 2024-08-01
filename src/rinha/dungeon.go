package rinha

import (
	"asura/src/entities"
)

type Room struct {
	Boss  entities.Rooster `json:"boss"`
	Level int              `json:"level"`
}

type DungeonWin struct {
	Percentage  int
	PrizeType   entities.ItemType
	PrizeRarity int
}

var DungeonsPercentages = [][]DungeonWin{
	{
		{Percentage: 7, PrizeType: entities.NormalType, PrizeRarity: 1},
		{Percentage: 52, PrizeType: entities.NormalType, PrizeRarity: 0},
		{PrizeType: entities.LootboxType, PrizeRarity: 0},
	},
	{
		{Percentage: 11, PrizeType: entities.NormalType, PrizeRarity: 1},
		{Percentage: 55, PrizeType: entities.NormalType, PrizeRarity: 0},
		{PrizeType: entities.LootboxType, PrizeRarity: 0},
	},
	{
		{Percentage: 8, PrizeType: entities.NormalType, PrizeRarity: 2},
		{Percentage: 52, PrizeType: entities.NormalType, PrizeRarity: 1},
		{PrizeType: entities.LootboxType, PrizeRarity: 1},
	},
	{
		{Percentage: 12, PrizeType: entities.NormalType, PrizeRarity: 2},
		{Percentage: 55, PrizeType: entities.NormalType, PrizeRarity: 1},
		{PrizeType: entities.LootboxType, PrizeRarity: 1},
	},
	{
		{Percentage: 49, PrizeType: entities.NormalType, PrizeRarity: 2},
		{PrizeType: entities.LootboxType, PrizeRarity: 2},
	},
	{
		{Percentage: 8, PrizeType: entities.NormalType, PrizeRarity: 3},
		{Percentage: 49, PrizeType: entities.NormalType, PrizeRarity: 2},
		{PrizeType: entities.LootboxType, PrizeRarity: 2},
	},
	{
		{Percentage: 49, PrizeType: entities.NormalType, PrizeRarity: 3},
		{PrizeType: entities.LootboxType, PrizeRarity: 3},
	},
	{
		{Percentage: 1, PrizeType: entities.LootboxType, PrizeRarity: 8},
		{Percentage: 4, PrizeType: entities.LootboxType, PrizeRarity: 7},
		{Percentage: 49, PrizeType: entities.NormalType, PrizeRarity: 4},
		{PrizeType: entities.LootboxType, PrizeRarity: 4},
	},
}
