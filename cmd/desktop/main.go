package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Word struct {
	ID      string `json:"id"`
	English string `json:"english"`
	Polish  string `json:"polish"`
}

var data binding.ExternalStringList
var addressFlag string

// translation of the word from the text field
func startSearch(word string) {
	data.Set(nil)
	words := translateWord(word)
	for _, item := range words {
		data.Append(item)
	}
}

func main() {
	// api server address (lidi-server), if no environment variable is defined
	// (DICTIONARY_SERVER), app takes the default value (http://localhost:8080)
	addressFlag = os.Getenv("DICTIONARY_SERVER")
	if addressFlag == "" {
		addressFlag = "http://localhost:8080"
	}

	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("Lidi desktop")
	myWindow.Resize(fyne.NewSize(720, 400))

	lista := []string{}

	data = binding.BindStringList(&lista)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	searchText := newEnterEntry()
	searchText.SetPlaceHolder("type a word to translate...")

	searchBtn := widget.NewButton("Search", func() {
		startSearch(searchText.Text)
	})

	rowSearch := container.New(layout.NewBorderLayout(nil, nil, nil, searchBtn))
	rowSearch.Add(searchText)
	rowSearch.Add(layout.NewSpacer())
	rowSearch.Add(searchBtn)

	myWindow.SetContent(container.NewBorder(rowSearch, nil, nil, nil, list))
	myWindow.Canvas().Focus(searchText)
	myWindow.ShowAndRun()
}
