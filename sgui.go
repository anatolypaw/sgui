package sgui

import (
	_ "fmt"
	"image"
	"image/draw"
	"log"
	_ "log"
)

// Основа. Отображает экраны.
// Одновременно активен может быть только один экран
type Sgui struct {
	Display      *image.RGBA //
	InputDevice  IInput      // Устройство ввода
	ActiveScreen *Screen     // Активный экран, который будет обрабатываться
	Overlay      *Overlay    // Отрисовывается поверх всех экранов, не
}

// Интерфейс устройства ввода
type IInput interface {
	GetEvent() IEvent
}

func New(display *image.RGBA, input IInput) (Sgui, error) {
	return Sgui{
		Display:     display,
		InputDevice: input,
	}, nil
}

// Устанавливает активный экран
func (ths *Sgui) SetScreen(screen *Screen) {
	screen.BackgroundRefill = true

	// Ждем, когда завершится обработка действующего экрана
	if ths.ActiveScreen != nil {
		ths.ActiveScreen.mu.Lock()
		defer ths.ActiveScreen.mu.Unlock()
	}

	ths.ActiveScreen = screen

	// Запускаем функцию, которая отрабатывает при ключении этого экрана
	if screen.RunOnce != nil {
		screen.RunOnce()
	}

}

// Устанавливает оверлей
func (ths *Sgui) SetOverlay(overlay *Overlay) {
	// Ждем, когда завершится обработка действующего экрана
	if ths.Overlay != nil {
		ths.Overlay.mu.Lock()
		defer ths.Overlay.mu.Unlock()
	}

	ths.Overlay = overlay
}

// Возвращает размер дисплея
func (ths *Sgui) SizeDisplay() image.Rectangle {
	return ths.Display.Bounds()
}

// Обрабатывает события ввода
// События обрабатываем в горутинах, что бы не пропустить
// новые приходящие события
func (ths *Sgui) StartInputEventHandler() {
	if ths.InputDevice == nil {
		return
	}
	go func() {
		for {
			event := ths.InputDevice.GetEvent()
			ths.Event(event)
		}
	}()
}

// Обрабатывает соыбытие ввода
func (ths *Sgui) Event(event IEvent) {
	// Передача нажатия в хук экрана
	if ths.ActiveScreen.TapHooker != nil {
		switch event.(type) {
		case EventTap:
			ths.ActiveScreen.TapHooker(event.Position())
		}
	}

	// Поиск виджетов в зоне нажатия и передача ему события
	for _, o := range ths.ActiveScreen.Objects {
		// если виджет отключен или скрыт, не передаем ему событие
		if o.Widget.Disabled() || o.Widget.Hidden() {
			continue
		}

		// определяем положение виджета на холсте
		wpos := image.Rect(
			o.Position.X,
			o.Position.Y,
			o.Widget.Size().X+o.Position.X,
			o.Widget.Size().Y+o.Position.Y,
		)

		// Если позиция тапа внутри виджета, то вызываем обработку тапа
		switch event.(type) {
		case EventTap:
			if event.Position().In(wpos) {
				go o.Widget.Tap(event.Position())
			}
		case EventRelease:
			go o.Widget.Release(event.Position())
		}
	}

}

// Отрисовывает объекты на дисплей
func (ths *Sgui) Render() {
	// Проверяем, установлен ли экран
	if ths.ActiveScreen == nil {
		return
	}
	ths.ActiveScreen.mu.Lock()
	defer ths.ActiveScreen.mu.Unlock()

	// Сначала рисуем background
	if ths.ActiveScreen.Background != nil && ths.ActiveScreen.BackgroundRefill {
		copy(ths.Display.Pix, ths.ActiveScreen.Background.Pix)
	}

	// Отрисовка на дисплей объектов с экрана, в порядке их добавления
	for _, o := range ths.ActiveScreen.Objects {
		ths.DrawObject(&o)
	}

	if ths.ActiveScreen.BackgroundRefill {
		ths.ActiveScreen.BackgroundRefill = false
	}

	// Отрисовываем оверлей
	if ths.Overlay == nil {
		return
	}
	ths.Overlay.mu.Lock()
	defer ths.Overlay.mu.Unlock()

	// Отрисовка на дисплей объектов с экрана, в порядке их добавления
	for _, o := range ths.Overlay.Objects {
		ths.DrawObject(&o)
	}

}

func (ths *Sgui) DrawObject(o *Object) {
	// Обновление внутреннего состояния виджета
	o.Widget.Update()

	// Если изображение виджета не менялось,
	// то и перерисовывать его не нужно. Пропускаем этот виджет
	// Если была отрисовка бэкграунда, то виджет нужно снова отрисовать
	if !o.Widget.Updated() && !ths.ActiveScreen.BackgroundRefill {
		return
	}

	wr := o.Widget.Render()

	if wr == nil {
		log.Println("SGUI: Widget render error - no render")
		return
	}

	draw.Draw(
		ths.Display,
		ths.Display.Bounds(),
		wr,
		image.Point{-o.Position.X, -o.Position.Y},
		draw.Src)
}
