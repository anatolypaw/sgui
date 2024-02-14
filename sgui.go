package sgui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"time"

	"github.com/anatolypaw/sgui/painter"
)

type Canvas struct {
	display    *image.RGBA //
	background *image.RGBA

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

func New(display *image.RGBA, input IInput) (Canvas, error) {
	return Canvas{
		display:     display,
		inputDevice: input,
	}, nil
}

// Возвращает размер дисплея
func (ths *Canvas) Size() image.Point {
	return image.Point{
		X: ths.display.Bounds().Max.X,
		Y: ths.display.Bounds().Max.Y,
	}
}

// Добавляет объект (widget) на холст
func (ui *Canvas) AddWidget(x int, y int, w IWidget) {
	obj := Object{
		Widget:   w,
		Position: image.Point{X: -x, Y: -y},
	}
	ui.objects = append(ui.objects, obj)
}

// Обрабатывает события ввода
// События обрабатываем в горутинах, что бы не пропустить
// новые приходящие события
func (ths *Canvas) StartInputEventHandler() {
	if ths.inputDevice == nil {
		return
	}
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

// Заливка заднего фона сплошным цветом
func (ths *Canvas) SetBackground(c color.Color) {
	ths.background = painter.DrawRectangle(
		painter.Rectangle{
			Size: image.Point{
				ths.display.Bounds().Dx(),
				ths.display.Bounds().Dy(),
			},
			FillColor: c,
		},
	)

}

// Обработка нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку нажатия
func (ths *Canvas) TapHandler(event IEvent) {
	for _, o := range ths.objects {

		// определяем положение виджета на холсте
		wpos := image.Rect(
			o.Position.X,
			o.Position.Y,
			o.Widget.Size().X,
			o.Widget.Size().Y,
		)

		// Если позиция тапа внутри виджета, то вызываем обработку тапа
		if event.Position().In(wpos) {
			o.Widget.Tap()
		}
	}

	ths.Render()
	fmt.Printf("Event: Tap, pos %#v\n", event.Position())
}

// Обработка отпускания нажатия
// Ищем какой объект попал в точку нажатия и вызываем на нем
// обработку  отжатия
func (ths *Canvas) ReleaseHandler() {
	for _, o := range ths.objects {

		o.Widget.Release()
	}
	ths.Render()
	fmt.Println("Event: Release")
}

// Отрисовывает объекты на дисплей
func (ths *Canvas) Render() {

	start := time.Now()

	// Сначала рисуем background
	if ths.background != nil {
		copy(ths.display.Pix, ths.background.Pix)
	}

	// Отрисовка на дисплей объектов в порядке их добавления
	for _, o := range ths.objects {
		draw.Draw(
			ths.display,
			ths.display.Bounds(),
			o.Widget.Render(),
			image.Point{o.Position.X, o.Position.Y},
			draw.Over)
	}

	log.Printf("Rendering  %v\n", time.Since(start))
}
