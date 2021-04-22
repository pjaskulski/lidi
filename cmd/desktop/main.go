package main

import (
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

func startSearch(word string) {
	data.Set(nil)
	words := translateEnglish(word)
	for _, item := range words {
		data.Append(item)
	}
}

func main() {
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
