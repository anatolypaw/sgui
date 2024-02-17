/*
Демонстрация работы интерфейса.
Веб использован просто как удобный
с точки зрения реализации способ взаимодействия.
С практической точки зрения он бесполезен, так как идет постоянная загрузка
новый кадров в png, что съедает много трафика.
Предназначен для отладки и разработки sgui
*/

package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"log/slog"
	"net/http"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

func main() {
	// Создаем дисплей
	rect := image.Rect(0, 0, 800, 480)
	display := image.NewRGBA(rect)

	// Создаем гуй
	gui, _ := sgui.New(display, nil)

	// Создаем тему
	theme := widget.ColorTheme{
		BackgroundColor: color.RGBA{240, 240, 240, 255},
		MainColor:       color.RGBA{200, 200, 200, 255},
		SecondColor:     color.RGBA{180, 180, 180, 255},
		StrokeColor:     color.RGBA{60, 60, 60, 255},
		StrokeWidth:     2,
		CornerRadius:    10,
	}

	// Создаем экраны
	mainScreen := sgui.NewScreen(gui.SizeDisplay())
	mainScreen.SetBackground(theme.BackgroundColor)

	secondScreen := sgui.NewScreen(gui.SizeDisplay())
	secondScreen.SetBackground(theme.BackgroundColor)

	// Создаем виджеты на основной экран
	ind := widget.NewIndicator(20, theme)
	ind.AddState(color.RGBA{255, 0, 0, 255})
	ind.AddState(color.RGBA{0, 255, 0, 255})

	button2 := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				if ind.GetState() == 0 {
					ind.SetState(1)
				} else {
					ind.SetState(0)
				}
			},
			Label:           "Button 2",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	button1 := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				if button2.Hidden() {
					button2.Show()
				} else {
					button2.Hide()
				}
			},
			Label:           "Hide",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	buttonSetSecondScreen := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				gui.SetScreen(&secondScreen)
			},
			Label:           "2 экран",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	buttonSetMainScreen := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				gui.SetScreen(&mainScreen)
			},
			Label:           "1 экран",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	// Добавляем виджеты на холст
	mainScreen.AddWidget(10, 10, button1)
	mainScreen.AddWidget(10, 60, button2)
	mainScreen.AddWidget(130, 70, ind)
	mainScreen.AddWidget(10, 200, buttonSetSecondScreen)

	secondScreen.AddWidget(10, 10, buttonSetMainScreen)

	// Устанавливаем активный экран
	gui.SetScreen(&mainScreen)

	http.Handle("/", http.HandlerFunc(MainPage))
	http.Handle("/render.png", http.HandlerFunc(GetRender(&gui, display)))
	http.Handle("/event", http.HandlerFunc(Event(&gui)))

	//
	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", nil)
}

// Возвращает тело основной страницы
func MainPage(w http.ResponseWriter, r *http.Request) {
	page := `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Обновляемое PNG изображение</title>
    <style>
        body {
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }

        img {
            max-width: 100%;
            max-height: 100%;
            cursor: pointer; /* Добавляем стиль указателя мыши для изображения */
        }
    </style>
</head>
<body>
    <img id="refreshingImage" src="render.png" alt="PNG Image">
    
    <script>
        document.getElementById("refreshingImage").addEventListener("mousedown", handleMouseDown);
        document.getElementById("refreshingImage").addEventListener("mouseup", handleMouseUp);

        function handleMouseDown(event) {
            sendEvent("mousedown", event);
        }

        function handleMouseUp(event) {
            sendEvent("mouseup", event);
        }

        function sendEvent(eventName, event) {
            var img = document.getElementById("refreshingImage");
            var imgRect = img.getBoundingClientRect();

            var relativeX = event.clientX - imgRect.left;
            var relativeY = event.clientY - imgRect.top;

            var coordinates = {
                x: relativeX,
                y: relativeY
            };

            var eventData = {
                event: eventName,
                coordinates: coordinates
            };

            fetch('/event', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(eventData),
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                console.log(data);
            })
            .catch(error => {
                console.error('There was an error!', error);
            });
        }

        // Функция для обновления изображения
        function updateImage() {
            var img = document.getElementById("refreshingImage");
            var currentSrc = img.src;
            // Добавляем случайное значение в URL для предотвращения кеширования
            img.src = currentSrc.split("?")[0] + "?" + new Date().getTime();
        }

        // Обновляем изображение каждые 100 миллисекунд (10 раз в секунду)
        setInterval(updateImage, 100);
    </script>
</body>
</html>

	`
	fmt.Fprint(w, page)
}

// Возвращает рендер
func GetRender(gui *sgui.Sgui, display *image.RGBA) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "image/png")

		gui.Render()
		err := png.Encode(w, display)
		if err != nil {
			log.Println(err)
		}
	}
}

// Передает событие в гуй
func Event(gui *sgui.Sgui) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type EventData struct {
			Event       string `json:"event"`
			Coordinates struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"coordinates"`
		}

		var event EventData

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&event)
		if err != nil {
			slog.Error("Json decoder", "error", err)
			return
		}

		fmt.Printf("%v\n", event)

		if event.Event == "mousedown" {
			tap := sgui.EventTap{
				Pos: image.Point{event.Coordinates.X, event.Coordinates.Y},
			}
			gui.Event(tap)
		}

		if event.Event == "mouseup" {
			release := sgui.EventRelease{
				Pos: image.Point{event.Coordinates.X, event.Coordinates.Y},
			}
			gui.Event(release)
		}

	}
}
