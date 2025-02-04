package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const boardSize = 8

func main() {
	chessApp := app.New()
	window := chessApp.NewWindow("Chess Game")
	window.Resize(fyne.NewSize(600, 600))

	board := container.NewGridWithColumns(boardSize)

	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			squareColor := color.White
			if (row+col)%2 == 1 {
				squareColor = color.Black
			}
			square := canvas.NewRectangle(squareColor)
			square.SetMinSize(fyne.NewSize(75, 75))
			board.Add(square)
		}
	}

	content := container.NewVBox(
		widget.NewLabel("Chess Game"),
		board,
	)

	window.SetContent(content)
	window.ShowAndRun()
}