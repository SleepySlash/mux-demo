package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	serverApi := router.PathPrefix("/api").Subrouter()
	serverApi.HandleFunc("/getUserProfile/{id}", getUserProfile).Methods("GET")

	serverApi.HandleFunc("/createProfile", createProfile).Methods("POST")
	serverApi.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	serverApi.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	serverApi.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", serverApi))
}
