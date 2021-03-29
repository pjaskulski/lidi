package main

import "github.com/gorilla/mux"

func RegisterRoutes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/en/{word}", getPolish).Methods("GET")
	router.HandleFunc("/api/add", createWord).Methods("POST")
	router.HandleFunc("/api/pl/{word}", getEnglish).Methods("GET")
	router.HandleFunc("/api/update", updateWord).Methods("PUT")
	router.HandleFunc("/api/delete", deleteWord).Methods("DELETE")

	return router
}
