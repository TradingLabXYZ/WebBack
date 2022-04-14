package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func InsertPrediction(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed inserting prediction, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "web" {
		log.Error("Failed inserting prediction, origin not web")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitionname := mux.Vars(r)["competition"]
	prediction := mux.Vars(r)["prediction"]
	s_prediction := map[string]string{"prediction": prediction}

	s_payload, err := json.Marshal(s_prediction)
	if err != nil {
		return
	}

	statement := `
		INSERT INTO submissions (
			updatedat, competitionname, userwallet, payload)
		VALUES (current_timestamp, $1, $2, $3)
		ON CONFLICT (userwallet) DO UPDATE
		SET updatedat = now(), payload = EXCLUDED.payload;`
	_, err = Db.Exec(
		statement,
		competitionname,
		session.UserWallet,
		s_payload)
	if err != nil {
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SelectPrediction(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed selecting prediction, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "web" {
		log.Error("Failed selecting prediction, origin not web")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitionname := mux.Vars(r)["competition"]

	var prediction string
	err = Db.QueryRow(`
			SELECT payload#>>'{prediction}'
			FROM submissions
			WHERE competitionname = $1
			AND userwallet = $2;`,
		competitionname,
		session.UserWallet).Scan(&prediction)
	if err != nil {
		log.Error(err)
		return
	}

	json.NewEncoder(w).Encode(prediction)
}

func UpdatePrediction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func DeletePrediction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func GetCountPartecipants(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("25")
}

func GetPartecipants(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("TEMP")
}
