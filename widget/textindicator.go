package widget

import (
	"image"
	"image/color"
	"image/draw"
	"log/slog"

	"github.com/anatolypaw/sgui/painter"
	"github.com/anatolypaw/sgui/text2img"
)

// Представляет собой множество переключаемых надписей

type TextIndicatorParam struct {
	Size            image.Point
	BackgroundColor color.Color
	CornerRadius    float64
	StrokeWidth     float64

	// Если эта функция указана, то данные берутся из нее
	// Предназначена для получения состояния индикатора
	StateSource func() int
}

type TextIndicator struct {
	param TextIndicatorParam

	currentState int
	states       []*image.RGBA
	hidden       bool
	disabled     bool

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	updated bool
}

func NewTextIndicator(p TextIndicatorParam) *TextIndicator {
	if p.Size.X <= 0 {
		p.Size.X = 1
	}

	if p.Size.Y <= 0 {
		p.Size.Y = 1
	}

	return &TextIndicator{
		param: p,
	}
}

func (w *TextIndicator) AddState(
	text string,
	textSize float64,
	textColor color.Color,
	fillColor color.Color,
	strokeColor color.Color) {

	// Создаем рендер основы надписи
	baseRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         w.param.Size,
			FillColor:    fillColor,
			BackColor:    w.param.BackgroundColor,
			CornerRadius: w.param.CornerRadius,
			StrokeWidth:  w.param.StrokeWidth,
			StrokeColor:  strokeColor,
		},
	)

	// Создаем рендер текста и вычисляем его расположение
	// для размещения в середине виджета
	textRender := text2img.Text2img(text, textSize, textColor)
	textMidPos := image.Point{
		X: -(w.param.Size.X - textRender.Rect.Dx()) / 2,
		Y: -(w.param.Size.Y-textRender.Rect.Dy())/2 - textRender.Rect.Dy()/12,
	}

	// Добавляем текст в базовый ренден
	draw.Draw(baseRender,
		baseRender.Bounds(),
		textRender,
		textMidPos,
		draw.Over)

	// Добавляем рендер в состояния
	w.states = append(w.states, baseRender)
}

func (w *TextIndicator) SetState(s int) {
	if s < 0 {
		w.currentState = 0
		return
	}

	if s > len(w.states)-1 {
		w.currentState = len(w.states) - 1
		return
	}

	if w.currentState != s {
		w.currentState = s
		w.updated = true
	}

}

func (w *TextIndicator) GetState() int {
	return w.currentState
}

func (w *TextIndicator) States() int {
	return len(w.states)
}

func (w *TextIndicator) Render() *image.RGBA {
	if w.states == nil {
		w.AddState(
			"NO STATE",
			10,
			color.White,
			color.Black,
			nil,
		)
		slog.Error("No states for BitIndicator. Created empty state")
	}

	// Получаем статус индикатора с внешней функции
	if w.param.StateSource != nil {
		w.SetState(w.param.StateSource())
	}

	w.updated = false
	return w.states[w.currentState]
}

func (w *TextIndicator) Tap(pos image.Point) {
}

func (w *TextIndicator) Release(pos image.Point) {
}

func (w *TextIndicator) Size() image.Point {
	return image.Point{
		w.param.Size.X,
		w.param.Size.Y,
	}
}

func (w *TextIndicator) Updated() bool {
	return w.updated
}

func (w *TextIndicator) Hide() {
	w.hidden = true
}

func (w *TextIndicator) Show() {
	w.hidden = false

}

func (w *TextIndicator) Disabled() bool {
	return w.disabled
}

func (w *TextIndicator) Hidden() bool {
	return w.hidden
}

func (w *TextIndicator) Update() {}
