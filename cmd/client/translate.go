package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type errorMessage struct {
	Message string
}

var ErrorNotFound error = errors.New("no translation found")

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
		log.Fatal(err)
	}
	return data, nil
}

// tłumaczennie z angielskiego na polski
func translateEnglish(word string, runSpeak bool, showID bool) {

	data, err := getTranslation(word, "en")
	if err != nil {
		if err != ErrorNotFound {
			log.Fatal(err)
		}
		fmt.Println(err)
		os.Exit(1)
	}

	var translateWords []Word

	err = json.Unmarshal(data, &translateWords)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range translateWords {
		if showID {
			fmt.Println("ID:", item.ID, item.English, " = ", item.Polish)
		} else {
			fmt.Println(item.English, " = ", item.Polish)
		}
		if runSpeak {
			speak(item.English)
		}
	}
}

// tłumaczenie z polskiego na angielski
func translatePolish(word string, runSpeak bool, showID bool) {

	data, err := getTranslation(word, "pl")
	if err != nil {
		if err != ErrorNotFound {
			log.Fatal(err)
		}
		fmt.Println(err)
		os.Exit(1)
	}

	var translateWords []Word

	err = json.Unmarshal(data, &translateWords)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range translateWords {
		if showID {
			fmt.Println("ID:", item.ID, item.Polish, " = ", item.English)
		} else {
			fmt.Println(item.Polish, " = ", item.English)
		}
		if runSpeak {
			speak(item.English)
		}
	}
}

// dołączanie nowego tłumaczenia
func addTranslation(translation string) {

	words := strings.Split(translation, "=")
	if len(words) != 2 {
		log.Fatal("error: new translation in form English=Polish was expected ex. house=dom")
	}
	if words[0] == "" || words[1] == "" {
		log.Fatal("error: empty value not accepted, English=Polish was expected")
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

	url := fmt.Sprintf("%s/api/add", cfg.addressFlag)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(post))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 201 {
		var msg errorMessage

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Failed to add a translation to the dictionary: ", err.Error())
		}

		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Fatal("Failed to add a translation to the dictionary: ", err.Error())
		}
		log.Fatal("Failed to add a translation to the dictionary: ", msg.Message)
	}

	var idJson Word

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Failed to read json: ", err.Error())
	}

	err = json.Unmarshal(body, &idJson)
	if err != nil {
		log.Fatal("Failed to unmarshal json: ", err.Error())
	}

	fmt.Println("New translation accepted", "ID:", idJson.ID)
}

// aktualizacja istniejącego tłumaczenia
func updateTranslation(recID string, word string) {

	if recID == "" {
		log.Fatal("error: valid record id was expected")
	}

	_, err := strconv.Atoi(recID)
	if err != nil {
		log.Fatal("error: valid record id was expected")
	}

	translation := strings.Split(word, "=")
	if len(translation) != 2 {
		log.Fatal("error: updated translation in form English=Polish was expected ex. home=dom")
	}
	if translation[0] == "" || translation[1] == "" {
		log.Fatal("error: empty value not accepted, English=Polish was expected")
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	todo := &Word{
		ID:      recID,
		English: translation[0],
		Polish:  translation[1],
	}

	post, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/api/update", cfg.addressFlag)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(post))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		var msg errorMessage

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Failed to update a translation: ", err.Error())
		}

		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Fatal("Failed to update a translation: ", err.Error())
		}
		log.Fatal("Failed to update a translation: ", msg.Message)
	}

	fmt.Println("Update accepted")
}

// usunięcie tłumaczenia
func deleteTranslation(recID string) {

	if recID == "" {
		log.Fatal("error: valid record id was expected")
	}

	_, err := strconv.Atoi(recID)
	if err != nil {
		log.Fatal("error: valid record id was expected")
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	todo := &Word{
		ID: recID,
	}

	post, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/api/delete", cfg.addressFlag)
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(post))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		var msg errorMessage

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Failed to delete a translation: ", err.Error())
		}

		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Fatal("Failed to delete a translation: ", err.Error())
		}
		log.Fatal("Failed to delete a translation: ", msg.Message)
	}

	fmt.Println("Translation deleted")
}
