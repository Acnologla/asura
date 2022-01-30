package utils

func MakeBoard(x, y int) []([]int) {
	board := []([]int){}
	for i := 0; i < x; i++ {
		row := []int{}
		for j := 0; j < y; j++ {
			row = append(row, 0)
		}
		board = append(board, row)
	}
	return board
}

func IsEqual(oldBoard []([]int), board []([]int)) bool {
	equal := true
	for i, row := range oldBoard {
		for j, tile := range row {
			if tile != board[i][j] {
				equal = false
				break
			}
		}
	}
	return equal
}

func DeepClone(board []([]int)) []([]int) {
	arr := []([]int){}
	for _, row := range board {
		newRow := []int{}
		newRow = append(newRow, row...)
		arr = append(arr, newRow)
	}
	return arr
}
