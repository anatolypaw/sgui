package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

func main() {
	// Создаем дисплей
	rect := image.Rect(0, 0, 800, 480)
	display := image.NewRGBA(rect)

	// Создаем гуй
	gui, err := sgui.New(display, nil)
	if err != nil {
		panic(err)
	}
	backcolor := color.RGBA{50, 50, 50, 255}
	gui.SetBackground(backcolor)

	for i := 0; i < 5; i++ {
		for n := 0; n < 10; n++ {
			// Создаем виджеты
			ind := widget.NewIndicator(20, backcolor)
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
			},
				backcolor)

			// Добавляем виджеты на холст
			gui.AddWidget(10+i*160, 10+(n*47), button)
			gui.AddWidget(130+i*160, 20+(n*47), ind)

		}
	}

	// Принудительно отрисовываем холст
	start := time.Now()
	gui.Render()
	log.Printf("Rendering  %v\n", time.Since(start))

	// Сохраняем рендер в файл
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
