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

type Button struct {
	size         image.Point
	onClick      func()
	currentState int  // Текущее состояние
	tapped       bool // Флаг, что кнопка нажата
	hidden       bool // флаг, что кнопка скрыта
	disabled     bool // флаг, что виджет не воспринимает события

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	updated          bool
	releasedRender   *image.RGBA
	pressedRenser    *image.RGBA
	backgroundRender *image.RGBA // используется, когда кнопка скрыта
}

type ButtonParam struct {
	Size            image.Point
	Onclick         func()
	Label           string
	LabelSize       float64
	ReleaseColor    color.Color
	PressColor      color.Color
	BackgroundColor color.Color
	CornerRadius    float64
	StrokeWidth     float64
	StrokeColor     color.Color
	TextColor       color.Color
}

func NewButton(p ButtonParam) *Button {
	if p.Size.X <= 0 {
		p.Size.X = 1
	}

	if p.Size.Y <= 0 {
		p.Size.Y = 1
	}

	// Создаем background для скрытого состояния
	backgroundRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:      p.Size,
			BackColor: p.BackgroundColor,
		},
	)

	// Состояния кнопки.
	// Нормальное состояние
	releasedRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         p.Size,
			FillColor:    p.ReleaseColor,
			BackColor:    p.BackgroundColor,
			CornerRadius: p.CornerRadius,
			StrokeWidth:  p.StrokeWidth,
			StrokeColor:  p.StrokeColor,
		},
	)

	// Нажатое состояние
	pressedRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         p.Size,
			FillColor:    p.PressColor,
			BackColor:    p.BackgroundColor,
			CornerRadius: p.CornerRadius,
			StrokeWidth:  p.StrokeWidth,
			StrokeColor:  p.StrokeColor,
		},
	)

	// Получаем изображение текста и вычисляем его расположение
	// для размещения в середине кнопки
	textRender := text2img.Text2img(p.Label, p.LabelSize, p.TextColor)
	textMidPos := image.Point{
		X: -(p.Size.X - textRender.Rect.Dx()) / 2,
		Y: -(p.Size.Y - textRender.Rect.Dy()) / 2,
	}

	// Наносим текст на оба состояния кнопки
	draw.Draw(releasedRender,
		releasedRender.Bounds(),
		textRender,
		textMidPos,
		draw.Over)

	draw.Draw(pressedRender,
		pressedRender.Bounds(),
		textRender,
		textMidPos,
		draw.Over)

	return &Button{
		size:             p.Size,
		onClick:          p.Onclick,
		releasedRender:   releasedRender,
		pressedRenser:    pressedRender,
		backgroundRender: backgroundRender,
		updated:          true,
	}
}

func (w *Button) Render() *image.RGBA {
	w.updated = false
	if w.hidden {
		return w.backgroundRender
	}

	if w.tapped {
		return w.pressedRenser
	}
	return w.releasedRender
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

// Скрывает кнопку
func (w *Button) Hide() {
	if !w.hidden {
		w.hidden = true
		w.updated = true
	}
}

func (w *Button) Show() {
	if w.hidden {
		w.hidden = false
		w.updated = true
	}

}

func (w *Button) Disabled() bool {
	return w.disabled
}

func (w *Button) Hidden() bool {
	return w.hidden
}
