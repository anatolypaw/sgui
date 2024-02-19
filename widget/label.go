package widget

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/anatolypaw/sgui/painter"
	"github.com/anatolypaw/sgui/text2img"
)

type Label struct {
	param LabelParam

	hidden  bool // флаг, что надпись скрыта
	useBase bool // флаг, что используется основа

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	textUpdated      bool        // текст был изменен
	baseUpdated      bool        // основа была изменена
	textRender       *image.RGBA // Рендер текста
	baseRender       *image.RGBA // Рендер основы надписи (заливка, рамка, скругление)
	finalRender      *image.RGBA // Рендер текста на основе
	backgroundRender *image.RGBA // используется, когда виджет скрыт
}

type LabelParam struct {
	Size            image.Point
	Text            string
	TextSize        float64
	TextColor       color.Color
	FillColor       color.Color
	BackgroundColor color.Color
	CornerRadius    float64
	StrokeWidth     float64
	StrokeColor     color.Color
}

func NewLabel(p LabelParam) *Label {
	if p.Size.X <= 0 {
		p.Size.X = 1
	}

	if p.Size.Y <= 0 {
		p.Size.Y = 1
	}

	// Создаем background для скрытого состояния
	backgroundRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:      p.Size,
			BackColor: p.BackgroundColor,
		},
	)

	// Основа не нужна, если не указаны цвета рамки и заливки
	useBase := p.FillColor != nil && p.StrokeColor != nil

	// Создаем рендер основы надписи
	baseRender := painter.DrawRectangle(
		painter.Rectangle{
			Size:         p.Size,
			FillColor:    p.FillColor,
			BackColor:    p.BackgroundColor,
			CornerRadius: p.CornerRadius,
			StrokeWidth:  p.StrokeWidth,
			StrokeColor:  p.StrokeColor,
		},
	)

	// Создаем рендер текста и вычисляем его расположение
	// для размещения в середине виджета
	textRender := text2img.Text2img(p.Text, p.TextSize, p.TextColor)
	textMidPos := image.Point{
		X: -(p.Size.X - textRender.Rect.Dx()) / 2,
		Y: -(p.Size.Y-textRender.Rect.Dy())/2 - textRender.Rect.Dy()/12,
	}

	rect := image.Rectangle{image.Point{0, 0}, p.Size}
	finalRender := image.NewRGBA(rect)

	if useBase {
		// Добавляем основу в финальный рендер
		draw.Draw(finalRender,
			finalRender.Bounds(),
			baseRender,
			image.Point{0, 0},
			draw.Src)
	}

	// Добавляем текст в финальный рендер
	draw.Draw(finalRender,
		finalRender.Bounds(),
		textRender,
		textMidPos,
		draw.Over)

	return &Label{
		param:            p,
		textUpdated:      false,
		baseUpdated:      false,
		useBase:          useBase,
		textRender:       textRender,
		baseRender:       baseRender,
		finalRender:      finalRender,
		backgroundRender: backgroundRender,
	}
}

// Установить новый текст
func (w *Label) SetText(s string) {
	w.param.Text = s
	w.textUpdated = true
}

// Отобразить виджет
func (w *Label) Show() {
	w.hidden = false
}

// Обработка нажатия на виджет
func (w *Label) Tap() {
	// Игнорируем нажатие
}

// Обработка отпускания виджета
func (w *Label) Release() {
	// Игнорируем отпускание
}

// Render implements sgui.IWidget.
func (w *Label) Render() *image.RGBA {

	// Была изменена основа
	if w.baseUpdated {
		// TODO Рендерим новую основу

	}

	// Был изменен текст
	if w.textUpdated {
		// Создаем рендер текста и вычисляем его расположение
		// для размещения в середине виджета
		w.textRender = text2img.Text2img(
			w.param.Text,
			w.param.TextSize,
			w.param.TextColor)
	}

	if w.textUpdated || w.baseUpdated {

		// Добавляем основу в финальный рендер
		if w.useBase {
			draw.Draw(w.finalRender,
				w.finalRender.Bounds(),
				w.baseRender,
				image.Point{0, 0},
				draw.Src)

		}

		textMidPos := image.Point{
			X: -(w.param.Size.X - w.textRender.Rect.Dx()) / 2,
			Y: -(w.param.Size.Y-w.textRender.Rect.Dy())/2 - w.textRender.Rect.Dy()/12,
		}

		// Добавляем текст в финальный рендер
		draw.Draw(w.finalRender,
			w.finalRender.Bounds(),
			w.textRender,
			textMidPos,
			draw.Over)

		w.textUpdated = false
		w.baseUpdated = false
	}

	return w.finalRender
}

// Возвращает размер виджета
func (w *Label) Size() image.Point {
	return w.param.Size
}

// Если изображение виджета обновилось, но не был вызван Render()
// то возвращается true, иначе false
func (w *Label) Updated() bool {
	return w.textUpdated || w.baseUpdated
}

// Скрыть виджет
func (w *Label) Hide() {
	w.hidden = true
}

// Вовзращаем, скрыт ли виджет
func (w *Label) Hidden() bool {
	return w.hidden
}

// Считаем что виджет отключен, когда скрыт
func (w *Label) Disabled() bool {
	return w.hidden
}
