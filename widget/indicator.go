package widget

import (
	"image"
	"image/color"
	"log/slog"

	"github.com/anatolypaw/sgui/painter"
)

// Простой круглый индикатор
// Имеет неограниченное количество переключающихся состояний
// Состояние представляет собой круг заданного цвета
// 1) Сначала нужно создать состояние через AddState()
// 2) Для изменения состояния испольузется SetState()

type bitIndicatorState struct {
	img *image.RGBA
}

type BitIndicator struct {
	size         int
	currentState int // Текущее состояние
	states       []bitIndicatorState
	theme        ColorTheme

	// Если эта функция указана, то она выполняется перед рендерингом
	// Предназначена для получения состояния индикатора
	stateLoader func() int

	hidden   bool
	disabled bool

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	updated bool
}

func NewIndicator(size int, stateLoader func() int, theme ColorTheme) *BitIndicator {
	if size <= 0 {
		size = 1
	}
	return &BitIndicator{
		size:        size,
		stateLoader: stateLoader,
		updated:     true,
		theme:       theme,
	}
}

func (w *BitIndicator) AddState(c color.Color) {
	circle := painter.Circle{
		Radius:      w.size / 2,
		FillColor:   c,
		BackColor:   w.theme.BackgroundColor,
		StrokeWidth: w.theme.StrokeWidth,
		StrokeColor: w.theme.StrokeColor,
	}
	img := painter.DrawCircle(circle)
	w.states = append(w.states, bitIndicatorState{img: img})
}

func (w *BitIndicator) SetState(s int) {
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

func (w *BitIndicator) GetState() int {
	return w.currentState
}

func (w *BitIndicator) States() int {
	return len(w.states)
}

func (w *BitIndicator) Render() *image.RGBA {
	if w.states == nil {
		w.AddState(color.RGBA{0, 0, 0, 0})
		slog.Error("No states for BitIndicator. Created empty state")
	}

	w.updated = false
	return w.states[w.currentState].img
}

func (w *BitIndicator) Tap(pos image.Point) {
}

func (w *BitIndicator) Release(pos image.Point) {
}

func (w *BitIndicator) Size() image.Point {
	return image.Point{
		w.size,
		w.size,
	}
}

func (w *BitIndicator) Updated() bool {
	return w.updated
}

func (w *BitIndicator) Hide() {
	w.hidden = true
}

func (w *BitIndicator) Show() {
	w.hidden = false

}

func (w *BitIndicator) Disabled() bool {
	return w.disabled
}

func (w *BitIndicator) Hidden() bool {
	return w.hidden
}

func (w *BitIndicator) Update() {
	// Получаем статус индикатора с внешней функции
	if w.stateLoader != nil {
		w.SetState(w.stateLoader())
	}
}
