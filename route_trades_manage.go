package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func CreateTrade(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed creating trade, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var new_trade NewTrade
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&new_trade)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed decoding new trade payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	new_trade.UserWallet = session.UserWallet

	if len(new_trade.Subtrades) == 0 {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed creating trade, missing subtrades",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = new_trade.InsertTrade()
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = new_trade.InsertSubTrades()
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
		}).Error(err)
		Db.Exec(`DELETE FROM trades WHERE code = $1`, new_trade.Code)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (new_trade *NewTrade) InsertTrade() (err error) {
	trade_sql := `
		INSERT INTO trades (
			code, userwallet, exchange, firstpair,
			secondpair, createdat, updatedat)
		VALUES (
			SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4,
			current_timestamp, current_timestamp)
		RETURNING code;`
	err = Db.QueryRow(
		trade_sql,
		new_trade.UserWallet,
		new_trade.Exchange,
		new_trade.FirstPairId,
		new_trade.SecondPairId,
	).Scan(&new_trade.Code)
	if err != nil {
		log.WithFields(log.Fields{
			"tradePayload": new_trade,
			"customMsg":    "Failed inserting new trade",
		}).Error(err)
		return
	}
	return
}

func DeleteTrade(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed deleting trade, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tradecode := mux.Vars(r)["tradecode"]
	if tradecode == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed deleting trade, empty tradecode",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Db.Exec(`
		DELETE FROM trades
		WHERE code = $1
		RETURNING userwallet;`, tradecode)

	w.WriteHeader(http.StatusOK)
}
