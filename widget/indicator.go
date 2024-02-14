package widget

import (
	"image"
	"image/color"
	"log/slog"

	"github.com/anatolypaw/sgui/painter"
)

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
}

func NewIndicator(size int) *BitIndicator {
	if size <= 0 {
		size = 1
	}
	return &BitIndicator{size: size}
}

func (w *BitIndicator) AddState(c color.Color) {
	circle := painter.Circle{
		Radius:      w.size / 2,
		FillColor:   c,
		StrokeWidth: 1,
		StrokeColor: color.RGBA{34, 34, 34, 255},
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

	w.currentState = s

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
	return w.states[w.currentState].img
}

func (w *BitIndicator) Tap() {
}

func (w *BitIndicator) Release() {
}

func (w *BitIndicator) Size() image.Point {
	return image.Point{
		w.size,
		w.size,
	}
}
