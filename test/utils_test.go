package test

import (
	"asura/src/utils"
	"testing"
)

func TestArray(t *testing.T) {
	t.Run("TestIncludes", func(t *testing.T) {
		arr := []string{"1", "2"}
		if !utils.Includes(arr, "1") {
			t.Errorf("'1' is in the array")
		}
		if utils.Includes(arr, "3") {
			t.Errorf("'3' is not in the array")

		}
	})
	t.Run("TestIndexOf", func(t *testing.T) {
		arr := []string{"1", "2"}
		if utils.IndexOf(arr, "1") != 0 {
			t.Errorf("'1' is in the 0 position of array")
		}
		if utils.IndexOf(arr, "3") != -1 {
			t.Errorf("'3' is not in the array")
		}
	})
	t.Run("TestSplice", func(t *testing.T) {
		arr := []string{"1", "2"}
		if len(utils.Splice(arr, 1)) != 1 {
			t.Errorf("Array length must be one")
		}
	})
}

func TestIsNumber(t *testing.T) {
	if utils.IsNumber("acnologia") {
		t.Errorf("This must return false")
	}
	if !utils.IsNumber("1") {
		t.Errorf("This must return true")
	}
}

func TestBoard(t *testing.T) {
	board := utils.MakeBoard(5, 5)
	if len(board[0]) != 5 {
		t.Errorf("Board size must be 5")
	}
	board2 := utils.MakeBoard(5, 5)
	board2[0][0] = 10
	if utils.IsEqual(board, board2) {
		t.Errorf("This must be false")
	}
}

func TestImage(t *testing.T) {
	if !utils.CheckImage("https://i.imgur.com/qrBSYBQ.png") {
		t.Errorf("This must return true")
	}
	if utils.CheckImage("https://www.youtube.com/") {
		t.Errorf("This must return false")
	}
}
