package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func translateEnglish(word string, runSpeak bool) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	format := "application/json"
	url := fmt.Sprintf("%s/api/en/%s", cfg.addressFlag, word)
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
		log.Println("No translation found")
		return
	} else if r.StatusCode >= 400 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(body))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
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

func translatePolish(word string, runSpeak bool) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	format := "application/json"
	url := fmt.Sprintf("%s/api/pl/%s", cfg.addressFlag, word)
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
		log.Println("No translation found")
		return
	} else if r.StatusCode >= 400 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(body))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
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
