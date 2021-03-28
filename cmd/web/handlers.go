package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func getPolish(w http.ResponseWriter, r *http.Request) {
	var words []Word

	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	result, err := db.Query("SELECT id, english, polish from engpol where english=?", params["word"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var word Word
		err := result.Scan(&word.ID, &word.English, &word.Polish)
		if err != nil {
			panic(err.Error())
		}
		words = append(words, word)
	}

	if len(words) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(words)
	}

}

func getEnglish(w http.ResponseWriter, r *http.Request) {
	var words []Word

	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	result, err := db.Query("SELECT id, english, polish from engpol where polish=?", params["word"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var word Word
		err := result.Scan(&word.ID, &word.English, &word.Polish)
		if err != nil {
			panic(err.Error())
		}
		words = append(words, word)
	}
	if len(words) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(words)
	}
}

// dołączanie nowego tłumaczenia do bazy danych
func createWord(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO engpol(english, polish) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	english := keyVal["english"]
	polish := keyVal["polish"]
	_, err = stmt.Exec(english, polish)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// aktualizacja tłumaczenia w bazie danych
func updateWord(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("UPDATE engpol SET english=?, polish=? WHERE english=? AND polish=?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	english := keyVal["english"]
	polish := keyVal["polish"]
	englishNew := keyVal["englishNew"]
	polishNew := keyVal["polishNew"]
	_, err = stmt.Exec(englishNew, polishNew, english, polish)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// obsługa usuwania tłumaczenia z bazy danych
func deleteWord(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("DELETE FROM engpol WHERE english=? AND polish=?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	english := keyVal["english"]
	polish := keyVal["polish"]

	_, err = stmt.Exec(english, polish)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
