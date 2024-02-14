package sgui

import "image"

type IEvent interface {
	Position() image.Point
}

// Нажатие
type EventTap struct {
	Pos image.Point
}

func (e EventTap) Position() image.Point {
	return e.Pos
}

// Отпускание
type EventRelease struct {
	Pos image.Point
}

func (e EventRelease) Position() image.Point {
	return e.Pos
}
