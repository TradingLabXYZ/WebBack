package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
		WHERE userwallet = $1;`,
		session.UserWallet).Scan(pq.Array(&tradeCodes))
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(tradeCodes)
}

func ListSubtrades(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed listing subtrades, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "api" {
		log.Error("Failed listing trades, origin not api")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tradecode := mux.Vars(r)["tradecode"]

	var subtradeCodes []string

	err = Db.QueryRow(`
		SELECT
			ARRAY_AGG(code)
		FROM subtrades
		WHERE userwallet = $1
		AND tradecode = $2;`,
		session.UserWallet,
		tradecode,
	).Scan(pq.Array(&subtradeCodes))
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(subtradeCodes)
}

func GetSnapshot(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed getting snapshot, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "api" {
		log.Error("Failed getting snapshot, origin not api")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := SelectUser("wallet", session.UserWallet)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg":  "Failed getting snapshot, wrong user",
			"userWallet": session.UserWallet,
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	snapshot := user.GetSnapshot()

	export_snapshot := struct {
		UserDetails       UserDetails
		Trades            []Trade
		CountTrades       int
		TotalReturnUsd    string
		TotalReturnBtc    string
		Roi               float64
		TotalPortfolioUsd string
	}{
		snapshot.UserDetails,
		snapshot.Trades,
		snapshot.CountTrades,
		snapshot.TotalReturnUsd,
		snapshot.TotalReturnBtc,
		snapshot.Roi,
		snapshot.TotalPortfolioUsd,
	}

	json.NewEncoder(w).Encode(export_snapshot)
}
