package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

func main() {
	// Создаем дисплей
	rect := image.Rect(0, 0, 800, 480)
	display := image.NewRGBA(rect)

	// Создаем устройство ввода

	// Создаем гуй
	gui, err := sgui.New(display, nil)
	if err != nil {
		panic(err)
	}

	// Создаем виджеты
	button := widget.NewButton(widget.Button{
		Size:      image.Point{X: 120, Y: 40},
		Onclick:   nil,
		Label:     "I'm button",
		LabelSize: 20,
	})

	background := widget.NewRectangle(
		image.Point{800, 480},
		color.RGBA{50, 50, 50, 255},
	)
	ind := widget.NewIndicator(20)
	ind.AddState(color.RGBA{255, 0, 0, 255})

	// Добавляем виджеты
	gui.AddWidget(0, 0, background)
	gui.AddWidget(100, 100, button)
	gui.AddWidget(100, 150, ind)

	// Принудительно отрисовываем холст
	gui.Render()

	f, err := os.Create("example.png")
	if err != nil {
		return
	}
	defer f.Close()

	err = png.Encode(f, display)
	if err != nil {
		return
	}
}
