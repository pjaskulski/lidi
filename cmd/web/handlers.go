package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// standardowy json z informacją o błędzie
func writeErrorMessage(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusConflict)
	temp := &errorMessage{Message: err.Error()}
	json.NewEncoder(w).Encode(temp)
}

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

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorMessage(w, err)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	english := keyVal["english"]
	polish := keyVal["polish"]

	result, id, err := lidiDB.recordAdd(english, polish)
	if !result {
		w.WriteHeader(http.StatusConflict)
		if err != nil {
			temp := &errorMessage{Message: err.Error()}
			json.NewEncoder(w).Encode(temp)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
		temp := &Word{ID: strconv.Itoa(int(id))}
		json.NewEncoder(w).Encode(temp)
	}

}

// aktualizacja tłumaczenia w bazie danych
func updateWord(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorMessage(w, err)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	id, err := strconv.Atoi(keyVal["id"])
	if err != nil {
		writeErrorMessage(w, err)
		return
	}

	english := keyVal["english"]
	polish := keyVal["polish"]

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

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorMessage(w, err)
		return
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	id, err := strconv.Atoi(keyVal["id"])
	if err != nil {
		writeErrorMessage(w, err)
		return
	}

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
