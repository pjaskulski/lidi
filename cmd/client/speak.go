package main

import (
	"log"
	"os"

	htgotts "github.com/hegedustibor/htgo-tts"
)

func speak(word string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	speech := htgotts.Speech{Folder: home + "/lidi-audio", Language: "en"}
	speech.Speak(word)
}
