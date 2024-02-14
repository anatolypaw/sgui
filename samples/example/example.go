package main

import (
	"fmt"
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
	gui.SetBackground(color.RGBA{50, 50, 50, 255})

	for i := 0; i < 5; i++ {
		for n := 0; n < 10; n++ {
			// Создаем виджеты

			ind := widget.NewIndicator(20)
			ind.AddState(color.RGBA{255, 0, 0, 255})
			ind.AddState(color.RGBA{0, 255, 0, 255})

			button := widget.NewButton(widget.ButtonParam{
				Size: image.Point{X: 110, Y: 40},
				Onclick: func() {
					if ind.GetState() == 0 {
						ind.SetState(1)
					} else {
						ind.SetState(0)
					}
				},
				Label:     fmt.Sprintf("Button %v", n+(i*10)),
				LabelSize: 20,
			})

			// Добавляем виджеты
			gui.AddWidget(10+i*160, 10+(n*47), button)
			gui.AddWidget(130+i*160, 20+(n*47), ind)

		}
	}

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
