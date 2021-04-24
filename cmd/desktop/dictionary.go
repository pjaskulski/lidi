package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	htgotts "github.com/hegedustibor/htgo-tts"
)

var ErrorNotFound error = errors.New("no translation found")
var from string

// function gets translation via REST API from lidi-server,
// response is returned as a byte array and error (or nil)
func getTranslation(word string, lang string) ([]byte, error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

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

// translation from English to Polish or from Polish to English
func translateWord(word string) []string {
	var result []string
	var answer []byte
	var err error

	from = "en"
	answer, err = getTranslation(word, from)

	// if eng->pl translation not found try pl->eng
	if err != nil && err == ErrorNotFound {
		from = "pl"
		answer, err = getTranslation(word, from)
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
		if from == "en" {
			line = item.Polish
		} else {
			line = item.English
		}

		result = append(result, line)
	}

	return result
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
