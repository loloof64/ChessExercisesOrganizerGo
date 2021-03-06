package main

import (
	"chess_exercises_organizer/appsettings"
	"chess_exercises_organizer/engine"
	"chess_exercises_organizer/text"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1100,
		Height: 800,
		Title:  "ChessExercisesOrganizer",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})

	app.Bind(engine.NewUciEngine())
	app.Bind(text.NewTextFileManager())
	app.Bind(appsettings.NewSettings())
	app.Run()
}
