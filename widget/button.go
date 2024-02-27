package widget

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/anatolypaw/sgui/painter"
	"github.com/anatolypaw/sgui/text2img"
)

// Имеет неограниченное количество переключающихся состояний
// Состояние представляет собой круг заданного цвета
// 1) Сначала нужно создать состояние через AddState()
// 2) Для изменения состояния испольузется SetState()

type Button struct {
	param ButtonParam

	tapped   bool // Флаг, что кнопка нажата
	disabled bool // флаг, что виджет не воспринимает события

	// Флаг, что изображение изменилось.
	// Сбрасывется после рендеринга
	sizeUpdated         bool
	textUpdated         bool
	releasedBaseUpdated bool
	pressedBaseUpdated  bool
	stateUpdated        bool

	textRender         *image.RGBA
	releasedRender     *image.RGBA
	finalRelesedRender *image.RGBA
	pressedRender      *image.RGBA
	finalPressedRender *image.RGBA
	backgroundRender   *image.RGBA // используется, когда кнопка скрыта

	// Если фнкция была передана, то она будет выполняться
	// каждый раз перед рендерингом.
	// если какой то из новых полученных параметров будет отличаться от текущих,
	// то он будет применен
	ParamSource func() ButtonParam
}

type ButtonParam struct {
	Size             image.Point
	OnClick          func()
	Text             string
	TextSize         float64
	ReleaseFillColor color.Color
	PressFillColor   color.Color
	BackgroundColor  color.Color
	CornerRadius     float64
	StrokeWidth      float64
	StrokeColor      color.Color
	TextColor        color.Color
	Hidden           bool
}

func NewButton(p *ButtonParam, ps func() ButtonParam) *Button {
	button := Button{}

	// Если параметры не переданы, то возвращаем пустую структуру
	if p == nil {
		// Не указан источник параметров
		if ps == nil {
			panic("Для виджета не указано ни одного источника параметров")
		}

		return &button
	}

	button.ParamSource = ps
	button.SetParam(*p)
	button.Render()

	return &button
}

// Установка параметров виджета
func (w *Button) SetParam(p ButtonParam) {
	w.param.Hidden = p.Hidden
	w.param.OnClick = p.OnClick

	w.SetSize(p.Size)
	w.SetBackground(p.BackgroundColor)
	w.SetText(p.Text, p.TextSize, p.TextColor)
	w.SetReleaseStyle(p.ReleaseFillColor, p.CornerRadius, p.StrokeWidth, p.StrokeColor)
	w.SetPressedStyle(p.PressFillColor, p.CornerRadius, p.StrokeWidth, p.StrokeColor)

}

// Установить размер
func (w *Button) SetSize(size image.Point) {
	if w.param.Size == size {
		return
	}
	w.param.Size = size
	w.sizeUpdated = true
}

// Установить задний фон
func (w *Button) SetBackground(c color.Color) {
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

// Установить основу для отжатого состояния(подложку)
func (w *Button) SetReleaseStyle(
	fillColor color.Color,
	cornerRadius float64,
	strokeWidth float64,
	strokeColor color.Color,
) {
	// Проверяем, отличаются ли новые параметры от сущесвубщих
	// Если не отличаются, то выходим
	if w.param.ReleaseFillColor == fillColor &&
		w.param.CornerRadius == cornerRadius &&
		w.param.StrokeWidth == strokeWidth &&
		w.param.StrokeColor == strokeColor &&
		!w.sizeUpdated &&
		!w.textUpdated {
		return
	}

	//Обновляем параметры
	w.param.ReleaseFillColor = fillColor
	w.param.CornerRadius = cornerRadius
	w.param.StrokeWidth = strokeWidth
	w.param.StrokeColor = strokeColor
	w.releasedBaseUpdated = true

	// Создаем рендер основы
	w.releasedRender = painter.DrawRectangle(
		painter.Rectangle{
			Size:         w.param.Size,
			FillColor:    w.param.ReleaseFillColor,
			BackColor:    w.param.BackgroundColor,
			CornerRadius: w.param.CornerRadius,
			StrokeWidth:  w.param.StrokeWidth,
			StrokeColor:  w.param.StrokeColor,
		},
	)

}

// Установить основу для нажатого состояния(подложку)
func (w *Button) SetPressedStyle(
	fillColor color.Color,
	cornerRadius float64,
	strokeWidth float64,
	strokeColor color.Color,
) {
	// Проверяем, отличаются ли новые параметры от сущесвубщих
	// Если не отличаются, то выходим
	if w.param.PressFillColor == fillColor &&
		w.param.CornerRadius == cornerRadius &&
		w.param.StrokeWidth == strokeWidth &&
		w.param.StrokeColor == strokeColor &&
		!w.sizeUpdated {
		return
	}
	//Обновляем параметры
	w.param.PressFillColor = fillColor
	w.param.CornerRadius = cornerRadius
	w.param.StrokeWidth = strokeWidth
	w.param.StrokeColor = strokeColor
	w.pressedBaseUpdated = true

	// Создаем рендер основы
	w.pressedRender = painter.DrawRectangle(
		painter.Rectangle{
			Size:         w.param.Size,
			FillColor:    w.param.PressFillColor,
			BackColor:    w.param.BackgroundColor,
			CornerRadius: w.param.CornerRadius,
			StrokeWidth:  w.param.StrokeWidth,
			StrokeColor:  w.param.StrokeColor,
		},
	)
}

// Установить новый текст
func (w *Button) SetText(text string, size float64, color color.Color) {
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

// Вызвать при нажатии на кнопку
func (w *Button) Tap() {
	if w.tapped {
		return
	}

	w.tapped = true
	w.stateUpdated = true

}

// Вызвать при отпускании кнопки
func (w *Button) Release() {
	if !w.tapped {
		return
	}
	w.tapped = false
	w.stateUpdated = true
	w.Click()
}

// Вызвывается когда предварительно нажатая кнопка была отпущенна
func (w *Button) Click() {
	if w.param.OnClick != nil {
		go w.param.OnClick()
	}
}

// Render implements sgui.IWidget.
func (w *Button) Render() *image.RGBA {
	// Подгонка размера финальных рендеров
	if w.sizeUpdated {
		rect := image.Rect(0, 0, w.param.Size.X, w.param.Size.Y)
		w.finalRelesedRender = image.NewRGBA(rect)
		w.finalPressedRender = image.NewRGBA(rect)
		w.sizeUpdated = false
	}

	// Композиция слоев

	// Была изменена основа отжатого состояния
	if w.releasedBaseUpdated {
		draw.Draw(w.finalRelesedRender,
			w.finalRelesedRender.Bounds(),
			w.releasedRender,
			image.Point{0, 0},
			draw.Src)

		// Текст нужно повторно отрисовать на основе
		w.textUpdated = true
	}

	// Была изменена основа нажатого состояния
	if w.pressedBaseUpdated {
		draw.Draw(w.finalPressedRender,
			w.finalPressedRender.Bounds(),
			w.pressedRender,
			image.Point{0, 0},
			draw.Src)
		// Текст нужно повторно отрисовать на основе
		w.textUpdated = true
	}

	// Был изменен текст
	if w.textUpdated {
		textMidPos := image.Point{
			X: -(w.param.Size.X - w.textRender.Rect.Dx()) / 2,
			Y: -(w.param.Size.Y-w.textRender.Rect.Dy())/2 - w.textRender.Rect.Dy()/12,
		}

		if w.releasedBaseUpdated {
			draw.Draw(w.finalRelesedRender,
				w.finalRelesedRender.Bounds(),
				w.textRender,
				textMidPos,
				draw.Over)
		}

		if w.pressedBaseUpdated {
			draw.Draw(w.finalPressedRender,
				w.finalPressedRender.Bounds(),
				w.textRender,
				textMidPos,
				draw.Over)
		}

		w.textUpdated = false
	}

	w.releasedBaseUpdated = false
	w.pressedBaseUpdated = false
	w.stateUpdated = false

	// Виджет скрыт
	if w.param.Hidden {
		return w.backgroundRender
	}

	if w.tapped {
		return w.finalPressedRender
	}
	return w.finalRelesedRender

}

func (w *Button) Size() image.Point {
	return w.param.Size
}

// Используется для определения, нужно ли вызывать функцию рендеринга
func (w *Button) Updated() bool {
	updated := w.sizeUpdated ||
		w.textUpdated ||
		w.releasedBaseUpdated ||
		w.pressedBaseUpdated ||
		w.stateUpdated

	return updated
}

// Скрывает кнопку
func (w *Button) Hide() {
	if w.param.Hidden {
		return
	}

	w.param.Hidden = true
}

func (w *Button) Show() {
	if !w.param.Hidden {
		return
	}
	w.param.Hidden = false

}

func (w *Button) Disabled() bool {
	return w.disabled
}

func (w *Button) Hidden() bool {
	return w.param.Hidden
}

func (w *Button) Update() {
	if w.param.Hidden {
		return
	}

	// Обновляем параметры виджета
	if w.ParamSource != nil {
		param := w.ParamSource()
		w.SetParam(param)
	}
}
