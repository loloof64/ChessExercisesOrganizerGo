package main

import (
	"chess_exercises_organizer/engine"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "ChessExercisesOrganizer",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})

	app.Bind(engine.NewUciEngine())
	app.Run()
}
