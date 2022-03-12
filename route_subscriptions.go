package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func SelectSubscriptionMonthlyPrice(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed selecting subscription, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "web" {
		log.Error("Failed selectin subscription, origin not web")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wallet := mux.Vars(r)["wallet"]

	observed, err := SelectUser("wallet", wallet)
	if err != nil {
		log.WithFields(log.Fields{
			"urlPath":  r.URL.Path,
			"observed": observed,
		}).Warn("Failed selecting subscription ws, user not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var sender string
	var createdat string
	var monthly_price int
	_ = Db.QueryRow(`
			SELECT DISTINCT ON(sender)
				sender,
				createdat AS eventcreatedat,
				payload#>>'{Value}' AS monthly_fee
			FROM smartcontractevents
			WHERE name = 'ChangePlan'
			AND LOWER(sender) = LOWER($1)
			ORDER BY 1, 2 DESC;`, wallet).Scan(
		&sender,
		&createdat,
		&monthly_price)

	json.NewEncoder(w).Encode(monthly_price)
}

func ManageUnsubscriptions() {
	for {
		err := Db.QueryRow(`
		DELETE FROM subscribers WHERE (subscribefrom, subscribeto) IN (
			SELECT
				sm.payload#>>'{Sender}' AS subscribefrom,
				sm.payload#>>'{To}' AS subscribeto
			FROM smartcontractevents sm
			INNER JOIN subscribers s1 ON(sm.payload#>>'{Sender}' = s1.subscribefrom)
			INNER JOIN subscribers s2 ON(sm.payload#>>'{To}' = s2.subscribeto)
			WHERE name = 'Subscribe'
			AND now() > (sm.createdat + (sm.payload#>>'{Weeks}')::INT * interval '1 week'));`)
		if err.Err() != nil {
			log.WithFields(log.Fields{
				"err": err.Err(),
			}).Warn("Failed managing unsubscriptions")
		}
		time.Sleep(2 * time.Hour)
	}
}
