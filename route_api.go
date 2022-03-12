package main

import (
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func ListTrades(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed listing trades, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "api" {
		log.Error("Failed listing trades, origin not api")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tradeCodes []string

	err = Db.QueryRow(`
		SELECT
			ARRAY_AGG(code)
		FROM trades
		WHERE userwallet = $1;`, session.UserWallet).Scan(pq.Array(&tradeCodes))
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(tradeCodes)
}
