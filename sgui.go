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
func (ths *Sgui) Size() image.Point {
	return image.Point{
		X: ths.display.Bounds().Max.X,
		Y: ths.display.Bounds().Max.Y,
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
func (ths *Sgui) StartInputEventHandler() {
	go func() {
		for {
			event := ths.inputDevice.GetEvent()
			switch event.(type) {
			case EventTap:
				go ths.TapHandler(event)
			case EventRelease:
				go ths.ReleaseHandler()
			}
		}
	}()
}

// Обработка нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку нажатия
func (ths *Sgui) TapHandler(event IEvent) {
	for _, o := range ths.objects {
		o.Widget.Tap()
	}
	ths.Render()
	fmt.Printf("Event: Tap, pos %#v\n", event.Position())
}

// Обработка отпускания нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку  отжатия
func (ths *Sgui) ReleaseHandler() {
	for _, o := range ths.objects {
		o.Widget.Release()
	}
	ths.Render()
	fmt.Println("Event: Release")
}

// Отрисовывает объекты на дисплей
func (ths *Sgui) Render() {

	start := time.Now()

	// Отрисовка на дисплей объектов в порядке их добавления
	for _, o := range ths.objects {
		draw.Draw(
			ths.display,
			ths.display.Bounds(),
			o.Widget.Render(),
			image.Point{o.Position.X, o.Position.Y},
			draw.Src)
	}

	log.Printf("Rendering  %v\n", time.Since(start))

}
