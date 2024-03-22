package widget

import (
	"image"
	"image/color"

	"github.com/anatolypaw/sgui/painter"
)

// Прямоугольник без действий залитый сплошным цветом
type rectangle struct {
	size       image.Point
	render     *image.RGBA
	background *image.RGBA
	hidden     bool
	updated    bool
}

// Disabled implements sgui.IWidget.
func (w *rectangle) Disabled() bool {
	return false
}

// Hidden implements sgui.IWidget.
func (w *rectangle) Hidden() bool {

	return w.hidden
}

// Hide implements sgui.IWidget.
func (w *rectangle) Hide() {
	w.updated = true
	w.hidden = true
}

// Show implements sgui.IWidget.
func (w *rectangle) Show() {
	w.updated = true
	w.hidden = false
}

// Update implements sgui.IWidget.
func (w *rectangle) Update() {
}

// Updated implements sgui.IWidget.
func (w *rectangle) Updated() bool {
	return w.updated
}

func NewRectangle(size image.Point, color color.Color, background color.Color) *rectangle {
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

	back := painter.DrawRectangle(
		painter.Rectangle{
			Size:      size,
			FillColor: background,
		},
	)

	return &rectangle{
		size:       size,
		render:     img,
		background: back,
	}
}

func (w *rectangle) Render() *image.RGBA {
	w.updated = false
	if w.hidden {
		return w.background
	}

	return w.render
}

// Вызвать при нажатии
func (w *rectangle) Tap(pos image.Point) {
}

// Вызвать при отпускании
func (w *rectangle) Release(pos image.Point) {
}

// Вызвывается когда предварительно нажатая кнопка была отпущенна
func (w *rectangle) Click() {
}

func (w *rectangle) Size() image.Point {
	return w.size
}
