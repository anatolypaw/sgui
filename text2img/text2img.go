// Преобразует текст в изображение

package text2img

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func Text2img(label string, size float64, color color.Color) *image.RGBA {
	fnt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingNone,
	})

	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	metrics := face.Metrics()
	meas := font.MeasureString(face, label)

	var uniform *image.Uniform

	if color != nil {
		uniform = image.NewUniform(color)
	} else {
		uniform = image.White
	}

	img := image.NewRGBA(image.Rect(0, 0, meas.Round(), int(metrics.Height/70)))
	drawer := font.Drawer{
		Dst:  img,
		Src:  uniform,
		Face: face,
		Dot:  fixed.P(0, int(size*0.80)),
	}

	drawer.DrawString(label)

	return img

}
