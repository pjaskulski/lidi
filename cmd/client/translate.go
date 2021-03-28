package main

import (
	"bytes"
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
	}

	req.Header.Set("Accept", format)

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
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
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	todo := &Word{
		English: words[0],
		Polish:  words[1],
	}

	post, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	format := "application/json"
	url := fmt.Sprintf("%s/api/add", cfg.addressFlag)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(post))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", format)

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 201 {
		log.Fatal("Failed to add a translation to the dictionary")
	}

	fmt.Println("New translation accepted")
}

func updateTranslation(translation string) {
	oldAndNew := strings.Split(translation, ":")
	if len(oldAndNew) != 2 {
		log.Fatal("error: update translation in form OldEnglish=OldPolish:NewEnglish=NewPolish was expected ex. house=dom:home=dom")
	}

	oldTranslation := strings.Split(oldAndNew[0], "=")
	if len(oldTranslation) != 2 {
		log.Fatal("error: update translation in form OldEnglish=OldPolish:NewEnglish=NewPolish was expected ex. house=dom:home=dom")
	}
	newTranslation := strings.Split(oldAndNew[1], "=")
	if len(newTranslation) != 2 {
		log.Fatal("error: update translation in form OldEnglish=OldPolish:NewEnglish=NewPolish was expected ex. house=dom:home=dom")
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	todo := &WordUpdate{
		English:    oldTranslation[0],
		Polish:     oldTranslation[1],
		EnglishNew: newTranslation[0],
		PolishNew:  newTranslation[1],
	}

	post, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	format := "application/json"
	url := fmt.Sprintf("%s/api/update", cfg.addressFlag)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(post))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", format)

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		log.Fatal("Failed to update a translation")
	}

	fmt.Println("Update accepted")
}

func deleteTranslation(translation string) {
	words := strings.Split(translation, "=")
	if len(words) != 2 {
		log.Fatal("error: delete translation in form English=Polish was expected ex. house=dom")
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	todo := &Word{
		English: words[0],
		Polish:  words[1],
	}

	post, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	format := "application/json"
	url := fmt.Sprintf("%s/api/delete", cfg.addressFlag)
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(post))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", format)

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		log.Fatal("Failed to delete a translation")
	}

	fmt.Println("Translation deleted")
}
