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
}

type textIndicator struct {
	param TextIndicatorParam

	currentState int
	states       []*image.RGBA
	hidden       bool
	disabled     bool

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	updated bool
}

func NewTextIndicator(p TextIndicatorParam) *textIndicator {
	if p.Size.X <= 0 {
		p.Size.X = 1
	}

	if p.Size.Y <= 0 {
		p.Size.Y = 1
	}

	return &textIndicator{
		param: p,
	}
}

func (w *textIndicator) AddState(
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

func (w *textIndicator) SetState(s int) {
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

func (w *textIndicator) GetState() int {
	return w.currentState
}

func (w *textIndicator) States() int {
	return len(w.states)
}

func (w *textIndicator) Render() *image.RGBA {
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
	w.updated = false
	return w.states[w.currentState]
}

func (w *textIndicator) Tap() {
}

func (w *textIndicator) Release() {
}

func (w *textIndicator) Size() image.Point {
	return image.Point{
		w.param.Size.X,
		w.param.Size.Y,
	}
}

func (w *textIndicator) Updated() bool {
	return w.updated
}

func (w *textIndicator) Hide() {
	w.hidden = true
}

func (w *textIndicator) Show() {
	w.hidden = false

}

func (w *textIndicator) Disabled() bool {
	return w.disabled
}

func (w *textIndicator) Hidden() bool {
	return w.hidden
}
