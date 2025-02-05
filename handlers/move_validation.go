package handlers

import "fmt"

func IsValidMove(board [8][8]rune, piece rune, fromRow, fromCol, toRow, toCol int) bool {

	if toRow < 0 || toRow >= 8 || toCol < 0 || toCol >= 8 {
		fmt.Println("Move out of bounds")
		return false
	}

	if board[toRow][toCol] != 0 {
		if (isWhite(piece) && isWhite(board[toRow][toCol])) || 
		   (!isWhite(piece) && !isWhite(board[toRow][toCol])) {
			fmt.Println("Can't capture own piece")
			return false
		}
	}

	switch piece {
	case 'P':
		if fromCol == toCol && (toRow == fromRow-1 || (fromRow == 6 && toRow == 4)) { 
			return true // Normal move or double-step from starting position
		}
	case 'p': 
		if fromCol == toCol && (toRow == fromRow+1 || (fromRow == 1 && toRow == 3)) { 
			return true // Normal move or double-step from starting position
		}
	case 'R', 'r':
		if fromRow == toRow || fromCol == toCol { 
			return true // Can move in straight lines
		}
	case 'N', 'n': 
		rowDiff, colDiff := abs(fromRow-toRow), abs(fromCol-toCol)
		if (rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2) {
			return true // L-shape movement
		}
	case 'B', 'b':
		if abs(fromRow-toRow) == abs(fromCol-toCol) {
			return true // Moves diagonally
		}
	case 'Q', 'q': 
		if fromRow == toRow || fromCol == toCol || abs(fromRow-toRow) == abs(fromCol-toCol) {
			return true // Moves like a rook or bishop
		}
	case 'K', 'k':
		if abs(fromRow-toRow) <= 1 && abs(fromCol-toCol) <= 1 {
			return true // Moves one square in any direction
		}
	}

	fmt.Println("Invalid move for", string(piece))
	return false
}

func isWhite(piece rune) bool {
	return piece == 'P' || piece == 'N' || piece == 'B' || piece == 'R' || piece == 'Q' || piece == 'K'
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
