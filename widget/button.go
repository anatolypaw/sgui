package widget

import (
	"image"
	"image/draw"

	"github.com/anatolypaw/sgui/painter"
	"github.com/anatolypaw/sgui/text2img"
)

// Имеет неограниченное количество переключающихся состояний
// Состояние представляет собой круг заданного цвета
// 1) Сначала нужно создать состояние через AddState()
// 2) Для изменения состояния испольузется SetState()

type Button struct {
	size         image.Point
	onClick      func()
	currentState int  // Текущее состояние
	tapped       bool // Флаг, что кнопка нажата

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	updated        bool
	releasedRender *image.RGBA
	pressedRenser  *image.RGBA
}

type ButtonParam struct {
	Size      image.Point
	Onclick   func()
	Label     string
	LabelSize float64
	Theme     ColorTheme
}

func NewButton(p ButtonParam) *Button {
	if p.Size.X <= 0 {
		p.Size.X = 1
	}

	if p.Size.Y <= 0 {
		p.Size.Y = 1
	}

	// Состояния кнопки.
	// Нормальное состояние
	releasedRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         p.Size,
			FillColor:    p.Theme.MainColor,
			BackColor:    p.Theme.BackgroundColor,
			CornerRadius: p.Theme.CornerRadius,
			StrokeWidth:  p.Theme.StrokeWidth,
			StrokeColor:  p.Theme.StrokeColor,
		},
	)

	// Нажатое состояние
	pressedRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         p.Size,
			FillColor:    p.Theme.SecondColor,
			BackColor:    p.Theme.BackgroundColor,
			CornerRadius: p.Theme.CornerRadius,
			StrokeWidth:  p.Theme.StrokeWidth,
			StrokeColor:  p.Theme.StrokeColor,
		},
	)

	// Получаем изображение текста и вычисляем его расположение
	// для размещения в середине кнопки
	textimg := text2img.Text2img(p.Label, p.LabelSize, p.Theme.TextColor)
	textMidPos := image.Point{
		X: -(p.Size.X - textimg.Rect.Dx()) / 2,
		Y: -(p.Size.Y - textimg.Rect.Dy()) / 2,
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

	return &Button{
		size:           p.Size,
		onClick:        p.Onclick,
		releasedRender: releasedRender,
		pressedRenser:  pressedRender,
		updated:        true,
	}
}

func (w *Button) Render() *image.RGBA {
	w.updated = false
	if w.tapped {
		return w.pressedRenser
	} else {
		return w.releasedRender
	}
}

// Вызвать при нажатии на кнопку
func (w *Button) Tap() {
	if w.currentState == 0 {
		w.currentState = 1
		w.tapped = true
		w.updated = true
	}

}

// Вызвать при отпускании кнопки
func (w *Button) Release() {
	if w.tapped {
		w.currentState = 0
		w.tapped = false
		w.updated = true
		w.Click()
	}

}

// Вызвывается когда предварительно нажатая кнопка была отпущенна
func (w *Button) Click() {

	if w.onClick != nil {
		go w.onClick()
	}
}

func (w *Button) Size() image.Point {
	return w.size
}

func (w *Button) Updated() bool {
	return w.updated
}
