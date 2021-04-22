package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrorNotFound error = errors.New("no translation found")

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
