package main

import "github.com/anatolypaw/sgui"

func main() {
	// Создаем дисплей

	// Создаем устройство ввода

	// Создаем гуй
	gui, err := sgui.New(nil, nil)
	if err != nil {
		panic(err)
	}

	_ = gui
}
