package sgui

import (
	"image"
	"sync"
)

// На экране размещаются различные виджеты
type Overlay struct {
	Objects []Object // виджеты и их положение на дисплее
	Size    image.Rectangle
	mu      sync.Mutex // Блокировка, когда идет работа с экраном
}

// Создает оверлей
func NewOverlay(size image.Rectangle) Overlay {
	return Overlay{
		Size: size,
	}
}

// Добавляет объект (widget) на экран
func (ui *Overlay) AddWidget(x int, y int, w IWidget) {
	obj := Object{
		Widget:   w,
		Position: image.Point{X: x, Y: y},
	}
	ui.Objects = append(ui.Objects, obj)
}
