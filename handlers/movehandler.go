package handlers

import (
	"fmt"
	pieces "chess-engine/peice_move_logic"
)

func HandlePieceClick(board [8][8]rune, piece rune, x, y int) {
	var possibleMoves []pieces.Move 

	switch piece {
	case 'R', 'r': 
		possibleMoves = pieces.GetRookMoves(board, x, y)
	default:
		fmt.Println("Move logic for this piece is not implemented yet.")
		return
	}

	fmt.Println("Possible Moves for", string(piece), "at", x, y, "->", possibleMoves)
}
