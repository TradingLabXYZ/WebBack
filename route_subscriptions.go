package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func SelectSubscriptionMonthlyPrice(w http.ResponseWriter, r *http.Request) {
	_, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed selecting subscription, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
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
