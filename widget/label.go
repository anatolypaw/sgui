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

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	sizeUpdated    bool // размер виджета был изменен
	textUpdated    bool // текст был изменен
	baseUpdated    bool // основа была изменена
	visibleUpdated bool // Изменена видимость виджета

	textRender       *image.RGBA // Рендер текста
	baseRender       *image.RGBA // Рендер основы надписи (заливка, рамка, скругление)
	finalRender      *image.RGBA // Рендер текста на основе
	backgroundRender *image.RGBA // используется, когда виджет скрыт

	// Если фнкция была передана, то она будет выполняться
	// каждый раз перед рендерингом.
	// если какой то из новых полученных параметров будет отличаться от текущих,
	// то он будет применен
	ParamSource func() LabelParam
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
	Hidden          bool
}

// Должен быть передан хотя бы один параметр
func NewLabel(p *LabelParam, ps func() LabelParam) *Label {
	label := Label{}

	// Если параметры не переданы, то возвращаем пустую структуру
	if p == nil {

		// Не указан источник параметров
		if ps == nil {
			panic("Для виджета label не указано ни одного источника параметров")
		}

		label.ParamSource = ps
		return &label
	}

	label.SetParam(*p)
	label.Render()

	return &label
}

// Установка параметров виджета
func (w *Label) SetParam(p LabelParam) {
	if w.param.Hidden != p.Hidden {
		w.param.Hidden = p.Hidden
		w.visibleUpdated = true
	}

	w.SetSize(p.Size)
	w.SetBackground(p.BackgroundColor)
	w.SetBase(p.FillColor, p.CornerRadius, p.StrokeWidth, p.StrokeColor)
	w.SetText(p.Text, p.TextSize, p.TextColor)
}

// Установить размер
func (w *Label) SetSize(size image.Point) {
	if w.param.Size == size {
		return
	}

	w.param.Size = size
	w.sizeUpdated = true
}

// Установить задний фон
func (w *Label) SetBackground(c color.Color) {
	// Создаем background для скрытого состояния

	// Цвет и размер не изменился, пропускаем
	if w.param.BackgroundColor == c &&
		!w.sizeUpdated {
		return
	}

	// Обновляем параметры
	w.param.BackgroundColor = c

	w.backgroundRender = painter.DrawRectangle(
		painter.Rectangle{
			Size:      w.param.Size,
			BackColor: w.param.BackgroundColor,
		},
	)

}

// Установить основу надписи (подложку)
func (w *Label) SetBase(
	fillColor color.Color,
	cornerRadius float64,
	strokeWidth float64,
	strokeColor color.Color,
) {
	// Проверяем, отличаются ли новые параметры от сущесвубщих
	// Если не отличаются, то выходим
	if w.param.FillColor == fillColor &&
		w.param.CornerRadius == cornerRadius &&
		w.param.StrokeWidth == strokeWidth &&
		w.param.StrokeColor == strokeColor &&
		!w.sizeUpdated {
		return
	}

	//Обновляем параметры
	w.param.FillColor = fillColor
	w.param.CornerRadius = cornerRadius
	w.param.StrokeWidth = strokeWidth
	w.param.StrokeColor = strokeColor

	w.baseUpdated = true

	// Создаем рендер основы надписи
	w.baseRender = painter.DrawRectangle(
		painter.Rectangle{
			Size:         w.param.Size,
			FillColor:    w.param.FillColor,
			BackColor:    w.param.BackgroundColor,
			CornerRadius: w.param.CornerRadius,
			StrokeWidth:  w.param.StrokeWidth,
			StrokeColor:  w.param.StrokeColor,
		},
	)

}

// Установить новый текст
func (w *Label) SetText(text string, size float64, color color.Color) {
	// Проверяем, отличаются ли новые параметры от сущесвубщих
	// Если не отличаются, то выходим
	if w.param.Text == text &&
		w.param.TextSize == size &&
		w.param.TextColor == color &&
		!w.sizeUpdated {
		return
	}

	// Обновляем параметры
	w.param.Text = text
	w.param.TextSize = size
	w.param.TextColor = color

	// Создаем рендер текста и вычисляем его расположение
	// для размещения в середине виджета
	w.textRender = text2img.Text2img(
		w.param.Text,
		w.param.TextSize,
		w.param.TextColor,
	)
	w.textUpdated = true
}

// Отобразить виджет
func (w *Label) Show() {
	w.param.Hidden = false
}

// Обработка нажатия на виджет
func (w *Label) Tap() {
	// Игнорируем нажатие
}

// Обработка отпускания виджета
func (w *Label) Release() {
	// Игнорируем отпускание
}

// Обновление внутреннего состояния виджета
func (w *Label) Update() {
	// Обновляем параметры виджета
	if w.ParamSource != nil {
		param := w.ParamSource()
		w.SetParam(param)
	}
}

// Render implements sgui.IWidget.
func (w *Label) Render() *image.RGBA {

	// Подгонка размера финального рендера
	if w.sizeUpdated {
		rect := image.Rect(0, 0, w.param.Size.X, w.param.Size.Y)
		w.finalRender = image.NewRGBA(rect)

		w.sizeUpdated = false
	}

	// Композиция слоев
	// Был изменен текст или основа
	if w.textUpdated || w.baseUpdated {
		// Рисуем основу
		draw.Draw(w.finalRender,
			w.finalRender.Bounds(),
			w.baseRender,
			image.Point{0, 0},
			draw.Src)

		// Рисуем текст
		textMidPos := image.Point{
			X: -(w.param.Size.X - w.textRender.Rect.Dx()) / 2,
			Y: -(w.param.Size.Y-w.textRender.Rect.Dy())/2 - w.textRender.Rect.Dy()/12,
		}

		draw.Draw(w.finalRender,
			w.finalRender.Bounds(),
			w.textRender,
			textMidPos,
			draw.Over)

		w.textUpdated = false
	}

	w.baseUpdated = false
	w.textUpdated = false

	// Виджет скрыт
	if w.param.Hidden {
		return w.backgroundRender
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
	return w.textUpdated || w.baseUpdated || w.sizeUpdated || w.visibleUpdated
}

// Скрыть виджет
func (w *Label) Hide() {
	w.param.Hidden = true
}

// Вовзращаем, скрыт ли виджет
func (w *Label) Hidden() bool {
	return w.param.Hidden
}

// Считаем что виджет отключен, когда скрыт
func (w *Label) Disabled() bool {
	return w.param.Hidden
}
