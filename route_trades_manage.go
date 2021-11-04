package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type NewSubtrade struct {
	CreatedAt string      `json:"CreatedAt"`
	Type      string      `json:"Type"`
	Reason    string      `json:"Reason"`
	Quantity  json.Number `json:"Quantity"`
	AvgPrice  json.Number `json:"AvgPrice"`
	Total     json.Number `json:"Total"`
	Usercode  string
}

type NewTrade struct {
	Exchange     string        `json:"Exchange"`
	FirstPairId  int           `json:"FirstPair"`
	SecondPairId int           `json:"SecondPair"`
	Subtrades    []NewSubtrade `json:"Subtrades"`
	Usercode     string
	Code         string
}

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

	new_trade.Usercode = session.UserCode

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
			code, usercode, exchange, firstpair,
			secondpair, createdat, updatedat, isopen)
		VALUES (
			SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4,
			current_timestamp, current_timestamp, true)
		RETURNING code;`
	err = Db.QueryRow(
		trade_sql,
		new_trade.Usercode,
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

func ChangeTradeStatus(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed changing trade status, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tradecode := mux.Vars(r)["tradecode"]
	if tradecode == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed changing trade status, empty tradecode",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	to_status_string := mux.Vars(r)["tostatus"]
	to_status, err := strconv.ParseBool(to_status_string)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed changing trade status, wrong bool",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var sentinel_1 string
	err = Db.QueryRow(`
		UPDATE trades
		SET
			isopen = $1,
			updatedat = current_timestamp
		WHERE code = $2
		RETURNING usercode;`, to_status, tradecode).Scan(&sentinel_1)
	if err != nil || sentinel_1 == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"tradeCode":   tradecode,
			"customMsg":   "Failed changing trade status, UPDATE trades",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var sentinel_2 string
	err = Db.QueryRow(`
		UPDATE subtrades
		SET
			updatedat = current_timestamp
		WHERE tradecode = $1
		RETURNING usercode;`, tradecode).Scan(&sentinel_2)
	if err != nil || sentinel_2 == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"tradeCode":   tradecode,
			"customMsg":   "Failed changing trade status, UPDATE subtrades",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		RETURNING usercode;`, tradecode)

	w.WriteHeader(http.StatusOK)
}
