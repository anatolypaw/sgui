package widget

import "image/color"

type ColorTheme struct {
	BackgroundColor color.Color
	MainColor       color.Color
	SecondColor     color.Color
	TextColor       color.Color
	StrokeColor     color.Color
	StrokeWidth     float64
	CornerRadius    float64
}
