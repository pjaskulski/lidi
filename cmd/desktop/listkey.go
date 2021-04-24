package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	htgotts "github.com/hegedustibor/htgo-tts"
)

var currentWord string = ""

// widget with a text field with support for Esc and Enter keys

type keyList struct {
	*widget.List
	current widget.ListItemID
}

/* function converts text to speech, uses google api thanks
   htgo-tts library, play downloaded mp3 file via mplayer
   the downloaded files for English words are stored in a subdirectory
   'lidi-audio' in the user's home directory therefore need not be
   retrieved again the next time the word is pronounced. */
func speak(word string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	speakLang := "en"
	if from == "en" {
		speakLang = "pl"
	}
	speech := htgotts.Speech{Folder: home + "/lidi-audio", Language: speakLang}
	speech.Speak(word)
}

func (l *keyList) onUp() {
	fmt.Println("Up")
}

func (l *keyList) onDown() {
	fmt.Println("Down")
}

func newKeyList(data binding.ExternalStringList) *keyList {

	list := &keyList{widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(2,
				widget.NewLabel("template"),
			)
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			item := i.(binding.String)
			textLabel := o.(*fyne.Container).Objects[0].(*widget.Label)
			textLabel.Bind(item)
		}), -1}

	list.ExtendBaseWidget(list)
	return list
}

func (l *keyList) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		l.onUp()
	case fyne.KeyReturn:
		l.onDown()
	}
}
