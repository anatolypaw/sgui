package widget

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/anatolypaw/sgui/painter"
	"github.com/anatolypaw/sgui/text2img"
)

// Имеет неограниченное количество переключающихся состояний
// Состояние представляет собой круг заданного цвета
// 1) Сначала нужно создать состояние через AddState()
// 2) Для изменения состояния испольузется SetState()

type button struct {
	size         image.Point
	onClick      func()
	currentState int  // Текущее состояние
	tapped       bool // Флаг, что кнопка нажата

	releasedRender *image.RGBA
	pressedRenser  *image.RGBA
}

type Button struct {
	Size      image.Point
	Onclick   func()
	Label     string
	LabelSize float64
}

func NewButton(b Button) *button {
	if b.Size.X <= 0 {
		b.Size.X = 1
	}

	if b.Size.Y <= 0 {
		b.Size.Y = 1
	}

	// Состояния кнопки.
	// 0 - кнопка отжата
	// 1 - кнопка нажата
	releasedRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         b.Size,
			FillColor:    color.RGBA{94, 94, 94, 255},
			CornerRadius: 8,
			StrokeWidth:  1,
			StrokeColor:  color.RGBA{34, 34, 34, 255},
		},
	)

	pressedRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         b.Size,
			FillColor:    color.RGBA{118, 118, 118, 255},
			CornerRadius: 8,
			StrokeWidth:  1,
			StrokeColor:  color.RGBA{34, 34, 34, 255},
		},
	)

	// Получаем изображение текста и вычисляем его расположение
	// для размещения в середине кнопки
	textimg := text2img.Text2img(b.Label, b.LabelSize)
	textMidPos := image.Point{
		X: -(b.Size.X - textimg.Rect.Dx()) / 2,
		Y: -(b.Size.Y - textimg.Rect.Dy()) / 2,
	}

	// Наносим текст на оба состояния кнопки
	draw.Draw(releasedRender,
		releasedRender.Bounds(),
		textimg,
		textMidPos,
		draw.Over)

	draw.Draw(pressedRender,
		pressedRender.Bounds(),
		textimg,
		textMidPos,
		draw.Over)

	return &button{
		size:           b.Size,
		onClick:        b.Onclick,
		releasedRender: releasedRender,
		pressedRenser:  pressedRender,
	}
}

func (w *button) Render() *image.RGBA {
	if w.tapped {
		return w.pressedRenser
	} else {
		return w.releasedRender
	}
}

// Вызвать при нажатии на кнопку
func (w *button) Tap() {
	w.currentState = 1
	w.tapped = true
}

// Вызвать при отпускании кнопки
func (w *button) Release() {
	if w.tapped {
		w.Click()
	}
	w.currentState = 0
	w.tapped = false
}

// Вызвывается когда предварительно нажатая кнопка была отпущенна
func (w *button) Click() {

	if w.onClick != nil {
		go w.onClick()
	}
}

func (w *button) Size() image.Point {
	return w.size
}
