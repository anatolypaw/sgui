package sgui

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"
)

type Sgui struct {
	display draw.Image //

	objects     []Object // виджеты и их положение на дисплее
	inputDevice IInput   // Устройство ввода
}

type Object struct {
	Widget   IWidget
	Position image.Point
}

// Интерфейс устройства ввода
type IInput interface {
	GetEvent() IEvent
}

// -
type IWidget interface {
	Render() *image.RGBA // Отрисовывает виджет
	Size() image.Point
	Tap() // Обработка нажатия и отпускания
	Release()
}

func New(display draw.Image, input IInput) (Sgui, error) {
	return Sgui{
		display:     display,
		inputDevice: input,
	}, nil
}

// Возвращает размер дисплея
func (ui *Sgui) Size() image.Point {
	return image.Point{
		X: ui.display.Bounds().Max.X,
		Y: ui.display.Bounds().Max.Y,
	}
}

// Добавляет объект (widget) на холст
func (ui *Sgui) AddWidget(x int, y int, w IWidget) {
	obj := Object{
		Widget:   w,
		Position: image.Point{X: -x, Y: -y},
	}
	ui.objects = append(ui.objects, obj)
}

// Обрабатывает события ввода
// События обрабатываем в горутинах, что бы не пропустить
// новые приходящие события
func (ui *Sgui) StartInputEventHandler() {
	go func() {
		for {
			event := ui.inputDevice.GetEvent()
			switch event.(type) {
			case EventTap:
				go ui.TapHandler(event.Position())
			case EventRelease:
				go ui.ReleaseHandler()
			}
		}
	}()
}

// Обработка нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку нажатия
func (ui *Sgui) TapHandler(pos image.Point) {
	for _, o := range ui.objects {
		o.Widget.Tap()
	}
	ui.Render()
	fmt.Printf("Tap event. pos %#v\n", pos)
}

// Обработка отпускания нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку  отжатия
func (ui *Sgui) ReleaseHandler() {
	for _, o := range ui.objects {
		o.Widget.Release()
	}
	ui.Render()
	fmt.Println("Release")
}

// Отрисовывает объекты на дисплей
func (ui *Sgui) Render() {

	start := time.Now()

	// Отрисовка на дисплей объектов в порядке их добавления
	for _, o := range ui.objects {
		draw.Draw(
			ui.display,
			ui.display.Bounds(),
			o.Widget.Render(),
			image.Point{o.Position.X, o.Position.Y},
			draw.Src)
	}

	log.Printf("Rendering  %v\n", time.Since(start))

}
