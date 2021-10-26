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

func InsertTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting InsertTrade..."))

	session, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user, err := SelectUser("email", session.Email)
	if err != nil {
		log.Warn("User not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	trade := struct {
		Exchange     string `json:"Exchange"`
		FirstPairId  int    `json:"FirstPair"`
		SecondPairId int    `json:"SecondPair"`
		Subtrades    []struct {
			CreatedAt string      `json:"CreatedAt"`
			Type      string      `json:"Type"`
			Reason    string      `json:"Reason"`
			Quantity  json.Number `json:"Quantity"`
			AvgPrice  json.Number `json:"AvgPrice"`
			Total     json.Number `json:"Total"`
		} `json:"subtrades"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&trade)
	if err != nil {
		log.Error(err)
	}

	var trade_code string
	trade_sql := `
		INSERT INTO trades (code, usercode, exchange, firstpair, secondpair, createdat, updatedat, isopen)
		VALUES (SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4, current_timestamp, current_timestamp, true)
		RETURNING code;`
	err = Db.QueryRow(
		trade_sql,
		user.Code,
		trade.Exchange,
		trade.FirstPairId,
		trade.SecondPairId,
	).Scan(&trade_code)
	if err != nil {
		log.Error(err)
	}

	for _, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (code, tradecode, usercode, createdat, type, reason, quantity, avgprice, total, updatedat)
		VALUES (SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4, $5, $6, $7, $8, current_timestamp)`
		Db.Exec(
			subtrade_sql,
			trade_code,
			user.Code,
			subtrade.CreatedAt,
			subtrade.Type,
			subtrade.Reason,
			subtrade.Quantity,
			subtrade.AvgPrice,
			subtrade.Total,
		)
	}
	if err != nil {
		log.Error(err)
	}

	json.NewEncoder(w).Encode("OK")
}

func UpdateTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateTrade..."))

	trade := struct {
		Code         string `json:"Code"`
		Usercode     string `json:"Usercode"`
		Exchange     string `json:"Exchange"`
		FirstPairId  int    `json:"FirstPairId"`
		SecondPairId int    `json:"SecondPairId"`
		Subtrades    []struct {
			CreatedAt string      `json:"CreatedAt"`
			Type      string      `json:"Type"`
			Reason    string      `json:"Reason"`
			Quantity  json.Number `json:"Quantity"`
			AvgPrice  json.Number `json:"AvgPrice"`
			Total     json.Number `json:"Total"`
		} `json:"Subtrades"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&trade)
	if err != nil {
		log.Error(err)
	}

	Db.Exec(`
		DELETE FROM subtrades
		WHERE tradecode = $1;
	`, trade.Code)

	for _, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (code, tradecode, usercode, createdat, type, reason, quantity, avgprice, total, updatedat)
		VALUES (SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4, $5, $6, $7, $8, current_timestamp)`
		Db.Exec(
			subtrade_sql,
			trade.Code,
			trade.Usercode,
			subtrade.CreatedAt,
			subtrade.Type,
			subtrade.Reason,
			subtrade.Quantity,
			subtrade.AvgPrice,
			subtrade.Total,
		)
	}

	json.NewEncoder(w).Encode("OK")
}

func CloseTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting CloseTrade..."))

	tradecode := mux.Vars(r)["tradecode"]

	Db.Exec(`
		UPDATE trades
		SET
			isopen = False,
			updatedat = current_timestamp
		WHERE code = $1;
		`, tradecode)

	json.NewEncoder(w).Encode("OK")
}

func OpenTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting OpenTrade..."))

	tradecode := mux.Vars(r)["tradecode"]

	Db.Exec(`
		UPDATE trades
		SET
			isopen = True,
			updatedat = current_timestamp
		WHERE code = $1;
		`, tradecode)

	json.NewEncoder(w).Encode("OK")
}

func DeleteTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting DeleteTrade..."))

	tradecode := mux.Vars(r)["tradecode"]

	Db.Exec(`
		DELETE FROM subtrades
		WHERE tradecode IN (
			SELECT code
			FROM trades
			WHERE code = $1);
		`, tradecode)

	Db.Exec(`
		DELETE FROM trades
		WHERE code = $1;
		`, tradecode)

	json.NewEncoder(w).Encode("OK")
}
