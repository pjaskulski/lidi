package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// pobranie tłumaczenia na polski
func getPolish(w http.ResponseWriter, r *http.Request) {
	var words []Word

	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	words = lidiDB.recordFind("english", params["word"])

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

	words = lidiDB.recordFind("polish", params["word"])

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

	result, err := lidiDB.recordAdd(english, polish)
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
	id, err := strconv.Atoi(keyVal["id"])
	if err != nil {
		log.Fatal(err)
	}
	english := keyVal["english"]
	polish := keyVal["polish"]

	w.Header().Set("Content-Type", "application/json")

	result, err := lidiDB.recordUpdate(id, english, polish)
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
	id, err := strconv.Atoi(keyVal["id"])
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")

	result, err := lidiDB.recordDelete(id)

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
