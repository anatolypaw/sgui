package painter

import (
	"image"
	"image/color"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const quarterCircleControl = 1 - 0.55228

type Circle struct {
	Radius      int
	FillColor   color.Color
	StrokeWidth float64
	StrokeColor color.Color
}

type Rectangle struct {
	Width        int
	Height       int
	FillColor    color.Color
	CornerRadius float64
	StrokeWidth  float64
	StrokeColor  color.Color
}

// Круг с обводкой
func DrawCircle(c Circle) *image.RGBA {
	size := c.Radius * 2
	rect := image.Rect(0, 0, size, size)
	img := image.NewRGBA(rect)

	mid := float64(c.Radius)
	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())

	if c.FillColor != nil {
		// Рисуем основу
		filler := rasterx.NewFiller(size, size, scanner)
		filler.SetColor(c.FillColor)
		rasterx.AddCircle(mid, mid, float64(c.Radius), filler)
		filler.Draw()
	}

	if c.StrokeColor != nil && c.StrokeWidth > 0 {
		// Рисуем обводку
		dasher := rasterx.NewDasher(size, size, scanner)
		dasher.SetColor(c.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(c.StrokeWidth*64), 0, nil, nil, nil, 0, nil, 0)
		rasterx.AddCircle(mid, mid, float64(c.Radius)-c.StrokeWidth/2, dasher)
		dasher.Draw()

	}

	return img
}

// Скругленный рямоугольник с обводкой
func DrawRectangle(r Rectangle) *image.RGBA {
	rect := image.Rect(0, 0, r.Width, r.Height)
	img := image.NewRGBA(rect)

	scanner := rasterx.NewScannerGV(r.Width, r.Height, img, img.Bounds())

	if r.FillColor != nil {
		// Рисуем основу
		filler := rasterx.NewFiller(r.Width, r.Height, scanner)
		filler.SetColor(r.FillColor)
		if r.CornerRadius <= 0 {
			rasterx.AddRect(0, 0, float64(r.Width), float64(r.Height), 0, filler)
		} else {
			rasterx.AddRoundRect(0, 0, float64(r.Width), float64(r.Height),
				r.CornerRadius, r.CornerRadius,
				0, rasterx.RoundGap, filler)
		}
		filler.Draw()
	}

	// Рисуем обводку
	if r.StrokeColor != nil && r.StrokeWidth > 0 {
		stk := float64(r.StrokeWidth / 2.1)
		p1x, p1y := stk, stk
		p2x, p2y := float64(r.Width)-stk, stk
		p3x, p3y := float64(r.Width)-stk, float64(r.Height)-stk
		p4x, p4y := stk, float64(r.Height)-stk
		rad := float64(r.CornerRadius) - r.StrokeWidth

		c := quarterCircleControl * rad
		dasher := rasterx.NewDasher(r.Width, r.Height, scanner)
		dasher.SetColor(r.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(r.StrokeWidth*64), 0, nil, nil, nil, 0, nil, 0)
		if c > 0 {
			dasher.Start(rasterx.ToFixedP(p1x, p1y+rad))
			dasher.CubeBezier(rasterx.ToFixedP(p1x, p1y+c),
				rasterx.ToFixedP(p1x+c, p1y),
				rasterx.ToFixedP(p1x+rad, p2y))
		} else {
			dasher.Start(rasterx.ToFixedP(p1x, p1y))
		}
		dasher.Line(rasterx.ToFixedP(p2x-rad, p2y))
		if c > 0 {
			dasher.CubeBezier(rasterx.ToFixedP(p2x-c, p2y),
				rasterx.ToFixedP(p2x, p2y+c),
				rasterx.ToFixedP(p2x, p2y+rad))
		}
		dasher.Line(rasterx.ToFixedP(p3x, p3y-rad))
		if c > 0 {
			dasher.CubeBezier(rasterx.ToFixedP(p3x, p3y-c),
				rasterx.ToFixedP(p3x-c, p3y),
				rasterx.ToFixedP(p3x-rad, p3y))
		}
		dasher.Line(rasterx.ToFixedP(p4x+rad, p4y))
		if c > 0 {
			dasher.CubeBezier(rasterx.ToFixedP(p4x+c, p4y),
				rasterx.ToFixedP(p4x, p4y-c),
				rasterx.ToFixedP(p4x, p4y-rad))
		}
		dasher.Stop(true)
		dasher.Draw()
	}

	return img

}

// Добавляет текст на изображение
func AddLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
