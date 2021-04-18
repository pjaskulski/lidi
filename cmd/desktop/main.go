package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func searchWord() {
	fmt.Println("Search start")
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("Lidi desktop")
	myWindow.Resize(fyne.NewSize(600, 400))

	lista := []string{}

	//for i := 4; i <= 103; i++ {
	//	lista = append(lista, fmt.Sprintf("%d word = słowo", i))
	//}

	data := binding.BindStringList(&lista)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	searchText := widget.NewEntry()
	searchBtn := widget.NewButton("Search", searchWord)
	rowSearch := container.New(layout.NewHBoxLayout())
	rowSearch.Add(searchText)
	rowSearch.Add(searchBtn)

	add := widget.NewButton("Append", func() {
		val := fmt.Sprintf("Item %d", data.Length()+1)
		data.Append(val)
	})
	myWindow.SetContent(container.NewBorder(rowSearch, add, nil, nil, list))
	myWindow.ShowAndRun()
}
