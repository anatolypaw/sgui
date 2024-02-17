package sgui

import (
	"image"
	"image/color"

	"github.com/anatolypaw/sgui/painter"
)

// На экране размещаются различные виджеты
type Screen struct {
	Background       *image.RGBA // Изображение бэкграунда
	Objects          []Object    // виджеты и их положение на дисплее
	Size             image.Rectangle
	BackgroundRefill bool
}

type Object struct {
	Widget   IWidget
	Position image.Point
}

// -
type IWidget interface {
	Render() *image.RGBA // Отрисовывает виджет, сбрасывает флаг updated
	Size() image.Point
	Updated() bool //Возвращает флаг, изменилось ли изображение виджета
	Tap()          // Обработка нажатия и отпускания
	Release()
	Hide()          // Включает флаг скрытия виджета
	Show()          // Отключает флаг скрытия виджета
	Hidden() bool   // Возвращает флаг, скрыт ли виджет
	Disabled() bool // Возвращает флаг, воспринимает ли виджет события
}

// Создает экран
func NewScreen(size image.Rectangle) Screen {
	return Screen{
		Size:             size,
		BackgroundRefill: true,
	}
}

// Добавляет объект (widget) на экран
func (ui *Screen) AddWidget(x int, y int, w IWidget) {
	obj := Object{
		Widget:   w,
		Position: image.Point{X: x, Y: y},
	}
	ui.Objects = append(ui.Objects, obj)
}

// Заливка заднего фона сплошным цветом
func (ths *Screen) SetBackground(c color.Color) {
	ths.Background = painter.DrawRectangle(
		painter.Rectangle{
			Size: image.Point{
				ths.Size.Dx(),
				ths.Size.Dy(),
			},
			FillColor: c,
		},
	)

}
