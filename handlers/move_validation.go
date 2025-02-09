package handlers

import "fmt"

type CastlingRights struct {
    WhiteKingSide  bool
    WhiteQueenSide bool
    BlackKingSide  bool
    BlackQueenSide bool
}

var initialPositions = map[string]bool{
    "e1": true, // White King
    "e8": true, // Black King
    "a1": true, // White Queen Rook
    "h1": true, // White King Rook
    "a8": true, // Black Queen Rook
    "h8": true, // Black King Rook
}

func IsSquareUnderAttack(board [8][8]rune, row, col int, isWhitePiece bool) bool {
    for i := 0; i < 8; i++ {
        for j := 0; j < 8; j++ {
            piece := board[i][j]
            if piece == 0 {
                continue
            }
            
            if isWhite(piece) == isWhitePiece {
                continue
            }
            
            if IsValidMove(board, piece, i, j, row, col, nil) {
                return true
            }
        }
    }
    return false
}

func IsInCheck(board [8][8]rune, isWhiteKing bool, kingRow, kingCol int) bool {
    return IsSquareUnderAttack(board, kingRow, kingCol, isWhiteKing)
}

func IsCastleable(board [8][8]rune, fromRow, fromCol, toRow, toCol int) bool {
    piece := board[fromRow][fromCol]
    
    if piece != 'K' && piece != 'k' {
        return false
    }
    
    if (piece == 'K' && (fromRow != 7 || fromCol != 4)) || 
       (piece == 'k' && (fromRow != 0 || fromCol != 4)) {
        return false
    }
    
    if abs(fromCol-toCol) != 2 || fromRow != toRow {
        return false
    }
    
    isKingSide := toCol > fromCol
    row := fromRow
    
    if IsInCheck(board, piece == 'K', row, fromCol) {
        return false
    }
    
    if isKingSide {
        if (piece == 'K' && board[7][7] != 'R') || (piece == 'k' && board[0][7] != 'r') {
            return false
        }
        
        for col := fromCol + 1; col < 7; col++ {
            if board[row][col] != 0 {
                return false
            }
            if IsSquareUnderAttack(board, row, col, piece == 'K') {
                return false
            }
        }
    } else {
        if (piece == 'K' && board[7][0] != 'R') || (piece == 'k' && board[0][0] != 'r') {
            return false
        }
        
        for col := fromCol - 1; col > 0; col-- {
            if board[row][col] != 0 {
                return false
            }
            if IsSquareUnderAttack(board, row, col, piece == 'K') {
                return false
            }
        }
    }
    
    return true
}

func IsValidMove(board [8][8]rune, piece rune, fromRow, fromCol, toRow, toCol int, promotionPiece *rune) bool {
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
            if toRow == 0 {
                fmt.Println("Promote pawn to (Q, R, B, N): ")
                fmt.Scanf("%c", promotionPiece)
                if *promotionPiece == 'Q' || *promotionPiece == 'R' || *promotionPiece == 'B' || *promotionPiece == 'N' {
                    return true
                } else {
                    fmt.Println("Invalid promotion piece")
                    return false
                }
            }
            return true 
        }
    case 'p': 
        if fromCol == toCol && (toRow == fromRow+1 || (fromRow == 1 && toRow == 3)) { 
            if toRow == 7 {
                fmt.Println("Promote pawn to (q, r, b, n): ")
                fmt.Scanf("%c", promotionPiece)
                if *promotionPiece == 'q' || *promotionPiece == 'r' || *promotionPiece == 'b' || *promotionPiece == 'n' {
                    return true
                } else {
                    fmt.Println("Invalid promotion piece")
                    return false
                }
            }
            return true 
        }
    case 'R', 'r':
        if fromRow == toRow || fromCol == toCol { 
            return true 
        }
    case 'N', 'n': 
        rowDiff, colDiff := abs(fromRow-toRow), abs(fromCol-toCol)
        if (rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2) {
            return true 
        }
    case 'B', 'b':
        if abs(fromRow-toRow) == abs(fromCol-toCol) {
            return true 
        }
    case 'Q', 'q': 
        if fromRow == toRow || fromCol == toCol || abs(fromRow-toRow) == abs(fromCol-toCol) {
            return true
        }
    case 'K', 'k':
        if abs(fromRow-toRow) <= 1 && abs(fromCol-toCol) <= 1 {
            return true
        }
        // Check for castling
        if IsCastleable(board, fromRow, fromCol, toRow, toCol) {
            return true
        }
    }

    fmt.Println("Invalid move for", string(piece))
    return false
}

func isWhite(piece rune) bool {
    return piece >= 'A' && piece <= 'Z'
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}