package sgui

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"
)

type sgui struct {
	display draw.Image //

	objects     []Object // виджеты и их положение на дисплее
	inputDevice IInput   // Устройство ввода
}

type Size struct {
	Width  int
	Height int
}

type Position struct {
	X int
	Y int
}

// Интерфейсы ввода
type IInput interface {
	GetEvent() IEvent
}

type IEvent interface {
	Position() Position
}

type Tap struct {
	Pos Position
}

func (t Tap) Position() Position {
	return t.Pos
}

type Release struct {
	Pos Position
}

func (t Release) Position() Position {
	return t.Pos
}

// -
type IWidget interface {
	Render() *image.RGBA // Отрисовывает виджет
	Size() Size
	Tap() // Обработка нажатия и отпускания
	Release()
}

type Object struct {
	Widget   IWidget
	Position Position
}

func New(display draw.Image, input IInput) (sgui, error) {
	return sgui{
		display:     display,
		inputDevice: input,
	}, nil
}

// Возвращает размер дисплея
func (ui *sgui) Size() Size {
	return Size{
		Width:  ui.display.Bounds().Max.X,
		Height: ui.display.Bounds().Max.Y,
	}
}

// Добавляет объект (widget) на холст
func (ui *sgui) AddWidget(x int, y int, w IWidget) {
	obj := Object{
		Widget:   w,
		Position: Position{X: -x, Y: -y},
	}
	ui.objects = append(ui.objects, obj)
}

// Обрабатывает события ввода
// События обрабатываем в горутинах, что бы не пропустить
// новые приходящие события
func (ui *sgui) StartInputWorker() {
	go func() {
		for {
			event := ui.inputDevice.GetEvent()
			switch event.(type) {
			case Tap:
				go ui.Tap(event.Position())
			case Release:
				go ui.Release()
			}
		}
	}()
}

// Обработка нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку нажатия
func (ui *sgui) Tap(pos Position) {
	for _, o := range ui.objects {
		o.Widget.Tap()
	}
	ui.Render()
	fmt.Printf("Tap event. pos %#v\n", pos)
}

// Обработка отпускания нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку  отжатия
func (ui *sgui) Release() {
	for _, o := range ui.objects {
		o.Widget.Release()
	}
	ui.Render()
	fmt.Println("Release")
}

// Отрисовывает объекты на дисплей
func (ui *sgui) Render() {

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
