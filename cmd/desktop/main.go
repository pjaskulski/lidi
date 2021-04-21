package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
var ErrorNotFound error = errors.New("no translation found")

type enterEntry struct {
	widget.Entry
}

func (e *enterEntry) onEsc() {
	e.Entry.SetText("")
}

func (e *enterEntry) onEnter() {
	startSearch(e.Entry.Text)
}

func newEnterEntry() *enterEntry {
	entry := &enterEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *enterEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		e.onEsc()
	case fyne.KeyReturn:
		e.onEnter()
	default:
		e.Entry.TypedKey(key)
	}
}

// funkcja pobiera tłumaczenie poprzez REST API z serwera lidi-server
// zwracana jest odpowiedź w formie tablicy bajów i błąd (lub nil)
func getTranslation(word string, lang string) ([]byte, error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	addressFlag := "http://localhost:8080"
	format := "application/json"
	url := fmt.Sprintf("%s/api/%s/%s", addressFlag, lang, word)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", format)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		return nil, ErrorNotFound
	} else if r.StatusCode >= 400 {
		txtError := fmt.Sprintf("error: %d", r.StatusCode)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			txtError += ", " + err.Error()
		}
		txtError += ", " + string(body)
		return nil, errors.New(txtError)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// tłumaczennie z angielskiego na polski
func translateEnglish(word string) []string {
	var result []string
	var answer []byte
	var err error
	var from string = "eng"

	answer, err = getTranslation(word, "en")

	// if eng->pl translation not found try pl->eng
	if err != nil && err == ErrorNotFound {
		from = "pl"
		answer, err = getTranslation(word, "pl")
	}

	if err != nil {
		result = append(result, err.Error())
		return result
	}

	var translateWords []Word

	err = json.Unmarshal(answer, &translateWords)
	if err != nil {
		result = append(result, err.Error())
		return result
	}

	for _, item := range translateWords {
		var line string
		if from == "eng" {
			line = item.English + " = " + item.Polish
		} else {
			line = item.Polish + " = " + item.English
		}

		result = append(result, line)
	}

	return result
}

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
