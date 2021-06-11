package test

import (
	"asura/src/utils/rinha"
	"testing"
)

func TestGetBackground(t *testing.T) {
	val := rinha.GetBackground(rinha.Galo{})
	if val != "https://i.imgur.com/F64ybgg.jpg" {
		t.Errorf("This should return the default background")
	}
}

func TestGetName(t *testing.T) {
	name := rinha.GetName("acno", rinha.Galo{
		Name: "acno o galo",
	})
	if name != "acno o galo" {
		t.Errorf("This should be 'acno o galo'")
	}
}

func TestIsIntInList(t *testing.T) {
	arr := []int{1, 2}
	if !rinha.IsIntInList(1, arr) {
		t.Errorf("1 is in the intList")
	}
	if rinha.IsIntInList(10, arr) {
		t.Errorf("10 isnt in the list")
	}
}

func TestLootbox(t *testing.T) {
	galo := rinha.Galo{
		CommonLootbox: 1,
	}
	if !rinha.HaveLootbox(galo, "comum") {
		t.Errorf("this should return true")
	}
	newGalo := rinha.GetNewLb("comum", galo, false)
	if newGalo.CommonLootbox != 0 {
		t.Errorf("this sould have 0 lootboxs")
	}
}
