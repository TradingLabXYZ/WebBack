package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
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
	Subtrades    []NewSubtrade `json:"subtrades"`
	Usercode     string
	Code         string
}

func CreateTrade(w http.ResponseWriter, r *http.Request) {
	_, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var new_trade NewTrade
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&new_trade)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed decoding new trade payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = new_trade.InsertTrade()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = new_trade.InsertSubTrades()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
			"tradepayload": new_trade,
			"custom_msg":   "Failed inserting new trade",
		}).Error(err)
		return
	}
	return
}

func CloseTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting CloseTrade..."))

	_, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tradecode := mux.Vars(r)["tradecode"]

	Db.Exec(`
		UPDATE trades
		SET
			isopen = False,
			updatedat = current_timestamp
		WHERE code = $1;`, tradecode)

	Db.Exec(`
		UPDATE subtrades
		SET
			updatedat = current_timestamp
		WHERE tradecode = $1;`, tradecode)

	json.NewEncoder(w).Encode("OK")
}

func OpenTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting OpenTrade..."))

	_, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tradecode := mux.Vars(r)["tradecode"]

	Db.Exec(`
		UPDATE trades
		SET
			isopen = True,
			updatedat = current_timestamp
		WHERE code = $1;`, tradecode)

	Db.Exec(`
		UPDATE subtrades
		SET
			updatedat = current_timestamp
		WHERE tradecode = $1;`, tradecode)

	json.NewEncoder(w).Encode("OK")
}

func DeleteTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting DeleteTrade..."))

	_, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tradecode := mux.Vars(r)["tradecode"]

	Db.Exec(`
		DELETE FROM trades
		WHERE code = $1;
		`, tradecode)

	json.NewEncoder(w).Encode("OK")
}
