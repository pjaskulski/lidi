package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// pobranie tłumaczenia na polski
func getPolish(w http.ResponseWriter, r *http.Request) {
	var words []Word

	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	words = recordFind("english", params["word"])

	if len(words) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(words)
	}

}

// pobranie tłumaczenia na anglielski
func getEnglish(w http.ResponseWriter, r *http.Request) {
	var words []Word

	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	words = recordFind("polish", params["word"])

	if len(words) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(words)
	}
}

// dołączanie nowego tłumaczenia do bazy danych
func createWord(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	english := keyVal["english"]
	polish := keyVal["polish"]

	w.Header().Set("Content-Type", "application/json")

	result, err := recordAdd(english, polish)
	if !result {
		w.WriteHeader(http.StatusConflict)
		if err != nil {
			temp := &errorMessage{Message: err.Error()}
			json.NewEncoder(w).Encode(temp)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}

}

// aktualizacja tłumaczenia w bazie danych
func updateWord(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	english := keyVal["english"]
	polish := keyVal["polish"]
	englishNew := keyVal["englishNew"]
	polishNew := keyVal["polishNew"]

	w.Header().Set("Content-Type", "application/json")

	result, err := recordUpdate(englishNew, polishNew, english, polish)
	if !result {
		w.WriteHeader(http.StatusConflict)
		if err != nil {
			temp := &errorMessage{Message: err.Error()}
			json.NewEncoder(w).Encode(temp)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// obsługa usuwania tłumaczenia z bazy danych
func deleteWord(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	english := keyVal["english"]
	polish := keyVal["polish"]

	w.Header().Set("Content-Type", "application/json")

	result, err := recordDelete(english, polish)

	if !result {
		w.WriteHeader(http.StatusConflict)
		if err != nil {
			temp := &errorMessage{Message: err.Error()}
			json.NewEncoder(w).Encode(temp)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

}
