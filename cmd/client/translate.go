package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// funkcja pobiera tłumaczenie poprzez REST API z serwera lidi-server
// zwracana jest odpowiedź w formie tablicy bajów i błąd (lub nil)
func getTranslation(word string, lang string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	format := "application/json"
	url := fmt.Sprintf("%s/api/%s/%s", cfg.addressFlag, lang, word)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req.Header.Set("Accept", format)

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		return nil, errors.New("no translation found")
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
		log.Fatal(err)
	}
	return data, nil
}

// tłumaczennie z angielskiego na polski
func translateEnglish(word string, runSpeak bool) {
	data, err := getTranslation(word, "en")
	if err != nil {
		log.Fatal(err)
	}

	var translateWords []Word

	err = json.Unmarshal(data, &translateWords)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range translateWords {
		fmt.Println(item.English, " = ", item.Polish)
		if runSpeak {
			speak(item.English)
		}
	}
}

// tłumaczenie z polskiego na angielski
func translatePolish(word string, runSpeak bool) {
	data, err := getTranslation(word, "pl")
	if err != nil {
		log.Fatal(err)
	}

	var translateWords []Word

	err = json.Unmarshal(data, &translateWords)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range translateWords {
		fmt.Println(item.Polish, " = ", item.English)
		if runSpeak {
			speak(item.English)
		}
	}
}

func addTranslation(translation string) {
	words := strings.Split(translation, "=")
	if len(words) != 2 {
		log.Fatal("error: new translation in form English=Polish was expected ex. house=dom")
		return
	}

}
