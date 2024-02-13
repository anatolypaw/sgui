package painter

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestDrawCircle(t *testing.T) {
	tests := []struct {
		name string
		c    Circle
	}{
		{"red", Circle{
			Radius:      50,
			FillColor:   color.RGBA{0, 255, 0, 255},
			StrokeWidth: 3,
			StrokeColor: color.RGBA{255, 0, 0, 200},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fname := fmt.Sprintf("./test_result/circle_%s.png", tt.name)
			f, err := os.Create(fname)
			if err != nil {
				return
			}
			defer f.Close()

			img := DrawCircle(tt.c)

			err = png.Encode(f, img)
			if err != nil {
				return
			}

		})
	}
}

func TestDrawRectangle(t *testing.T) {
	tests := []struct {
		name string
		r    Rectangle
	}{
		{"red", Rectangle{
			Width:        50,
			Height:       15,
			FillColor:    color.RGBA{94, 94, 94, 255},
			CornerRadius: 5,
			StrokeWidth:  1,
			StrokeColor:  color.RGBA{34, 34, 34, 255},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fname := fmt.Sprintf("./test_result/rectangle_%s.png", tt.name)
			f, err := os.Create(fname)
			if err != nil {
				return
			}
			defer f.Close()

			img := DrawRectangle(tt.r)

			err = png.Encode(f, img)
			if err != nil {
				return
			}

		})
	}
}
