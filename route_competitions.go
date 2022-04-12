package main

import(
	"encoding/json"
	"net/http"
)


func GetCountPartecipants(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("25")
}

func GetPartecipants(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("TEMP")
}

func InsertPrediction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
