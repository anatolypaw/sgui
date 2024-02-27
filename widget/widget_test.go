package widget

import (
	"fmt"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"testing"
)

func TestDrawCircle(t *testing.T) {
	tests := []struct {
		name   string
		radius int
	}{
		{"indicator", rand.Int() % 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			indicator := NewIndicator(tt.radius, nil, ColorTheme{
				BackgroundColor: color.White,
				StrokeColor:     color.Black,
				StrokeWidth:     2,
			})
			img := indicator.Render()

			fname := fmt.Sprintf("%s.png", tt.name)
			f, err := os.Create(fname)
			if err != nil {
				return
			}
			defer f.Close()

			err = png.Encode(f, img)
			if err != nil {
				return
			}

		})
	}
}
