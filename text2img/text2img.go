// Преобразует текст в изображение

package text2img

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func Text2img(label string, size float64) *image.RGBA {
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

	img := image.NewRGBA(image.Rect(0, 0, meas.Round(), int(metrics.Height/70)))
	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(0, int(size*0.80)),
	}

	drawer.DrawString(label)

	return img

}

func foo() {

	text := "Привет Как ДКЕЛАр"
	size := float64(50)

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
	meas := font.MeasureString(face, text)

	fmt.Printf("metrics: %#v\nmeasure %#v\n", metrics, meas)

	img := image.NewRGBA(image.Rect(0, 0, meas.Round(), int(metrics.Height/70)))

	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(0, int(size*0.80)),
	}

	drawer.DrawString(text)

	fname := "./test_result/text.png"
	file, err := os.Create(fname)
	if err != nil {
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return
	}

}
