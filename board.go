package main

import (
	"image/color"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

func main() {
	chessApp := app.New()
	window := chessApp.NewWindow("Chess Game")
	window.Resize(fyne.NewSize(600, 600))

	board := container.NewGridWithColumns(boardSize)
	parsedBoard := parseFEN(strings.Split(startFenNotation, " ")[0])

	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			squareColor := color.White
			if (row+col)%2 == 1 {
				squareColor = color.Black
			}
			square := canvas.NewRectangle(squareColor)
			square.SetMinSize(fyne.NewSize(75, 75))

			blockHasPiece := parsedBoard[row][col]
			if blockHasPiece != 0 {
				imagePath := filepath.Join(pieceDir, mpPieceToImage[blockHasPiece])
				pieceImage := canvas.NewImageFromFile(imagePath)
				pieceImage.FillMode = canvas.ImageFillContain
				pieceImage.Resize(fyne.NewSize(75, 75))

				cell := container.NewStack(square, pieceImage)
				board.Add(cell)
			} else {
				board.Add(square)
			}
		}
	}

	content := container.NewVBox(
		widget.NewLabel("Chess Game"),
		board,
	)

	window.SetContent(content)
	window.ShowAndRun()
}