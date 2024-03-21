package widget

import (
	"image"
)

// Если добавить этот виджет на экран, то он будет перехватывать все нажатия

type TapHooker struct {
	disabled bool // флаг, что виджет не воспринимает события
	OnClick  func(image.Point)
}

func NewTapHooker(f func(image.Point)) *TapHooker {
	return &TapHooker{
		OnClick: f,
	}
}

// Вызвать при нажатии на кнопку
func (w *TapHooker) Tap(pos image.Point) {
	if w.OnClick != nil {
		go w.OnClick(pos)
	}
}

// Вызвать при отпускании кнопки
func (w *TapHooker) Release(pos image.Point) {
}

func (w *TapHooker) Update() {
}

func (w *TapHooker) Render() *image.RGBA {
	return nil

}

func (w *TapHooker) Size() image.Point {
	return image.Point{}
}

// Используется для определения, нужно ли вызывать функцию рендеринга
func (w *TapHooker) Updated() bool {
	return false
}

func (w *TapHooker) Hide() {
}

func (w *TapHooker) Show() {
}

func (w *TapHooker) Disabled() bool {
	return w.disabled
}

func (w *TapHooker) Hidden() bool {
	return false
}

// Ловить все события
func (w *TapHooker) IsHookAllEvent() bool {
	return true
}
