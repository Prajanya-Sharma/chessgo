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

func movePiece(fromRow, fromCol, toRow, toCol int) {
	if fromRow == toRow && fromCol == toCol {
		fmt.Println("Invalid move: Same position")
		return
	}

	piece := parsedBoard[fromRow][fromCol]

	if !handlers.IsValidMove(parsedBoard, piece, fromRow, fromCol, toRow, toCol) {
		fmt.Println("Invalid move for piece:", string(piece))
		return
	}

	parsedBoard[toRow][toCol] = piece
	parsedBoard[fromRow][fromCol] = 0 
	fmt.Printf("Moved %c from (%d, %d) to (%d, %d)\n", piece, fromRow, fromCol, toRow, toCol)

	pieceSelected = false
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
						// First click: Select piece
						selectedRow, selectedCol = rowCopy, colCopy
						pieceSelected = true
						fmt.Println("Piece selected at:", selectedRow, selectedCol)
					} else {
						// Second click: Move piece
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

	content := container.NewVBox(
		widget.NewLabel("Chess Game"),
		generateChessBoard(),
	)

	window.SetContent(content)
	window.ShowAndRun()
}
