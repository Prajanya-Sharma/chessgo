package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"chess-engine/handlers"
)

const boardSize = 8
const pieceDir = "chess-gui/peices"

const startFenNotation = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var mpPieceToImage = map[rune]string{
	'P': "whitePawn.svg", 'N': "whiteKnight.svg", 'B': "whiteBishop.svg", 'R': "whiteRook.svg",
	'Q': "whiteQueen.svg", 'K': "whiteKing.svg",
	'p': "blackPawn.svg", 'n': "blackKnight.svg", 'b': "blackBishop.svg", 'r': "blackRook.svg",
	'q': "blackQueen.svg", 'k': "blackKing.svg",
}

var selectedRow, selectedCol int
var pieceSelected bool
var parsedBoard [8][8]rune
var whiteTurn = true

var boardContainer *fyne.Container
var boardCells [8][8]*fyne.Container

func parseFEN(fen string) [8][8]rune {
	var board [8][8]rune
	rows := strings.Split(fen, "/")
	for rowIdx, row := range rows {
		colIdx := 0
		for _, char := range row {
			if char >= '1' && char <= '8' {
				colIdx += int(char - '0')
			} else {
				board[rowIdx][colIdx] = char
				colIdx++
			}
		}
	}
	return board
}

func isPathClear(fromRow, fromCol, toRow, toCol int) bool {
	rowStep, colStep := 0, 0
	piece := parsedBoard[fromRow][fromCol]
	if piece == 'N' || piece == 'n' {
		return true
	}
	if fromRow < toRow {
		rowStep = 1
	} else if fromRow > toRow {
		rowStep = -1
	}

	if fromCol < toCol {
		colStep = 1
	} else if fromCol > toCol {
		colStep = -1
	}

	r, c := fromRow+rowStep, fromCol+colStep
	for r != toRow || c != toCol {
		if parsedBoard[r][c] != 0 {
			return false // Path is blocked
		}
		r += rowStep
		c += colStep
	}

	return true // Path is clear
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func movePiece(fromRow, fromCol, toRow, toCol int) {
	if fromRow == toRow && fromCol == toCol {
		pieceSelected = false
		return
	}

	piece := parsedBoard[fromRow][fromCol]
	isWhitePiece := piece >= 'A' && piece <= 'Z'
	targetPiece := parsedBoard[toRow][toCol]
	isTargetWhitePiece := targetPiece >= 'A' && targetPiece <= 'Z'

	// Check if it's the right player's turn
	if (whiteTurn && !isWhitePiece) || (!whiteTurn && isWhitePiece) {
		fmt.Println("Not your turn!")
		return
	}

	// ye hanlers se castling move detection
	if (piece == 'K' || piece == 'k') && abs(fromCol-toCol) == 2 {
		if handlers.IsCastleable(parsedBoard, fromRow, fromCol, toRow, toCol) {
			isKingSide := toCol > fromCol
			rookFromCol := 0
			rookToCol := 3
			if isKingSide {
				rookFromCol = 7
				rookToCol = 5
			}

			// ye hai for Move king
			parsedBoard[toRow][toCol] = piece
			parsedBoard[fromRow][fromCol] = 0

			// ye hai for Move rook
			rook := 'R'
			if !isWhitePiece {
				rook = 'r'
			}
			parsedBoard[toRow][rookToCol] = rook
			parsedBoard[toRow][rookFromCol] = 0

			// Update UI for both king and rook moves individually kyunki deletion bhi krni hai
			updateBoardUI(fromRow, fromCol, toRow, toCol)
			updateBoardUI(toRow, rookFromCol, toRow, rookToCol)

			whiteTurn = !whiteTurn
			pieceSelected = false
			return
		}
		pieceSelected = false
		return
	}

	if targetPiece != 0 && ((whiteTurn && isTargetWhitePiece) || (!whiteTurn && !isTargetWhitePiece)) {
		fmt.Println("Can't capture own piece")
		return
	}

	if !isPathClear(fromRow, fromCol, toRow, toCol) {
		fmt.Println("Path is blocked for piece:", string(piece))
		return
	}

	var promotionPiece rune
	if !handlers.IsValidMove(parsedBoard, piece, fromRow, fromCol, toRow, toCol, &promotionPiece) {
		fmt.Println("Invalid move for piece:", string(piece))
		return
	}

	if promotionPiece != 0 {
		piece = promotionPiece
	}

	parsedBoard[toRow][toCol] = piece
	parsedBoard[fromRow][fromCol] = 0
	fmt.Printf("Moved %c from (%d, %d) to (%d, %d)\n", piece, fromRow, fromCol, toRow, toCol)

	whiteTurn = !whiteTurn
	pieceSelected = false

	updateBoardUI(fromRow, fromCol, toRow, toCol)
}

func updateBoardUI(fromRow, fromCol, toRow, toCol int) {
	fromCell := boardCells[fromRow][fromCol]
	toCell := boardCells[toRow][toCol]

	fromCell.Objects = fromCell.Objects[:1] 
	toCell.Objects = toCell.Objects[:1]     

	if parsedBoard[fromRow][fromCol] != 0 {
		imagePath := filepath.Join(pieceDir, mpPieceToImage[parsedBoard[fromRow][fromCol]])
		pieceImage := canvas.NewImageFromFile(imagePath)
		pieceImage.FillMode = canvas.ImageFillContain
		pieceImage.Resize(fyne.NewSize(75, 75))

		rowCopy, colCopy := fromRow, fromCol
		tapButton := widget.NewButton(" ", func() {
			if !pieceSelected {
				selectedRow, selectedCol = rowCopy, colCopy
				pieceSelected = true
				fmt.Println("Piece selected at:", selectedRow, selectedCol)
			} else {
				movePiece(selectedRow, selectedCol, rowCopy, colCopy)
			}
		})
		tapButton.Importance = widget.LowImportance
		tapButton.Resize(fyne.NewSize(75, 75))

		fromCell.Add(pieceImage)
		fromCell.Add(tapButton)
	} else {
		rowCopy, colCopy := fromRow, fromCol
		tapButton := widget.NewButton(" ", func() {
			if pieceSelected {
				movePiece(selectedRow, selectedCol, rowCopy, colCopy)
			}
		})
		tapButton.Importance = widget.LowImportance
		tapButton.Resize(fyne.NewSize(75, 75))
		fromCell.Add(tapButton)
	}

	if parsedBoard[toRow][toCol] != 0 {
		imagePath := filepath.Join(pieceDir, mpPieceToImage[parsedBoard[toRow][toCol]])
		pieceImage := canvas.NewImageFromFile(imagePath)
		pieceImage.FillMode = canvas.ImageFillContain
		pieceImage.Resize(fyne.NewSize(75, 75))

		rowCopy, colCopy := toRow, toCol
		tapButton := widget.NewButton(" ", func() {
			if !pieceSelected {
				selectedRow, selectedCol = rowCopy, colCopy
				pieceSelected = true
				fmt.Println("Piece selected at:", selectedRow, selectedCol)
			} else {
				movePiece(selectedRow, selectedCol, rowCopy, colCopy)
			}
		})
		tapButton.Importance = widget.LowImportance
		tapButton.Resize(fyne.NewSize(75, 75))

		toCell.Add(pieceImage)
		toCell.Add(tapButton)
	} else {
		rowCopy, colCopy := toRow, toCol
		tapButton := widget.NewButton(" ", func() {
			if pieceSelected {
				movePiece(selectedRow, selectedCol, rowCopy, colCopy)
			}
		})
		tapButton.Importance = widget.LowImportance
		tapButton.Resize(fyne.NewSize(75, 75))
		toCell.Add(tapButton)
	}

	fromCell.Refresh()
	toCell.Refresh()
}

func generateChessBoard() *fyne.Container {
	board := container.NewGridWithColumns(boardSize)

	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			squareColor := color.White
			if (row+col)%2 == 1 {
				squareColor = color.Black
			}
			square := canvas.NewRectangle(squareColor)
			square.SetMinSize(fyne.NewSize(75, 75))

			blockHasPiece := parsedBoard[row][col]
			cell := container.NewStack(square)

			if blockHasPiece != 0 {
				imagePath := filepath.Join(pieceDir, mpPieceToImage[blockHasPiece])
				pieceImage := canvas.NewImageFromFile(imagePath)
				pieceImage.FillMode = canvas.ImageFillContain
				pieceImage.Resize(fyne.NewSize(75, 75))

				rowCopy, colCopy := row, col
				tapButton := widget.NewButton(" ", func() {
					if !pieceSelected {
						selectedRow, selectedCol = rowCopy, colCopy
						pieceSelected = true
						fmt.Println("Piece selected at:", selectedRow, selectedCol)
					} else {
						movePiece(selectedRow, selectedCol, rowCopy, colCopy)
					}
				})
				tapButton.Importance = widget.LowImportance
				tapButton.Resize(fyne.NewSize(75, 75))

				cell = container.NewStack(square, pieceImage, tapButton)
			} else {
				rowCopy, colCopy := row, col
				tapButton := widget.NewButton(" ", func() {
					if pieceSelected {
						movePiece(selectedRow, selectedCol, rowCopy, colCopy)
					}
				})
				tapButton.Importance = widget.LowImportance
				tapButton.Resize(fyne.NewSize(75, 75))
				cell = container.NewStack(square, tapButton)
			}

			boardCells[row][col] = cell
			board.Add(cell)
		}
	}

	return board
}

func main() {
	chessApp := app.New()
	window := chessApp.NewWindow("Chess Game")
	window.Resize(fyne.NewSize(600, 600))

	parsedBoard = parseFEN(strings.Split(startFenNotation, " ")[0])

	boardContainer = container.NewVBox(
		widget.NewLabel("Chess Game"),
		generateChessBoard(),
	)

	window.SetContent(boardContainer)
	window.ShowAndRun()
}