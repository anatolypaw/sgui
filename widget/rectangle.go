package widget

import (
	"image"
	"image/color"

	"github.com/anatolypaw/sgui/painter"
)

// Прямоугольник без действий залитый сплошным цветом
type rectangle struct {
	size   image.Point
	render *image.RGBA
}

func NewRectangle(size image.Point, color color.Color) *rectangle {
	if size.X <= 0 {
		size.X = 1
	}

	if size.Y <= 0 {
		size.Y = 1
	}

	img := painter.DrawRectangle(
		painter.Rectangle{
			Size:      size,
			FillColor: color,
		},
	)

	return &rectangle{
		size:   size,
		render: img,
	}
}

func (w *rectangle) Render() *image.RGBA {
	return w.render
}

// Вызвать при нажатии
func (w *rectangle) Tap() {
}

// Вызвать при отпускании
func (w *rectangle) Release() {
}

// Вызвывается когда предварительно нажатая кнопка была отпущенна
func (w *rectangle) Click() {
}

func (w *rectangle) Size() image.Point {
	return w.size
}
