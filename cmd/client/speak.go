package main

import (
	"log"
	"os"

	htgotts "github.com/hegedustibor/htgo-tts"
)

/* funkcja przetwarza teskt na mowę, wykorzystuje api google dzięki
   bibliotece htgo-tts, odtwarza pobrany plik mp3 poprzez mplayer-a
   pobrane pliki dla angielskich słów są przechowywane w podkatalogu
   'lidi-audio' w katalogu domowym użytkownika dlatego nie muszą być
   pobierane powtórnie przy kolejnym odtworzeniu wymowy słowa. */
func speak(word string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	speech := htgotts.Speech{Folder: home + "/lidi-audio", Language: "en"}
	speech.Speak(word)
}
