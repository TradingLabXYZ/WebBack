package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func (new_trade *NewTrade) InsertSubTrades() (err error) {
	subtrade_sql := `
		INSERT INTO subtrades (
			code, tradecode, userwallet, 
			createdat, type, reason, 
			quantity, avgprice, total, updatedat)
		VALUES %s;`

	valueStrings := []string{}
	valueArgs := []interface{}{}
	timeNow := time.Now()
	for i, subtrade := range new_trade.Subtrades {
		rand_subtrade_code := RandStringBytes(12)
		str1 := "$" + strconv.Itoa(1+i*10) + ","
		str2 := "$" + strconv.Itoa(2+i*10) + ","
		str3 := "$" + strconv.Itoa(3+i*10) + ","
		str4 := "$" + strconv.Itoa(4+i*10) + ","
		str5 := "$" + strconv.Itoa(5+i*10) + ","
		str6 := "$" + strconv.Itoa(6+i*10) + ","
		str7 := "$" + strconv.Itoa(7+i*10) + ","
		str8 := "$" + strconv.Itoa(8+i*10) + ","
		str9 := "$" + strconv.Itoa(9+i*10) + ","
		str10 := "$" + strconv.Itoa(10+i*10)
		str_n := "(" + str1 + str2 + str3 + str4 + str5 + str6 + str7 + str8 + str9 + str10 + ")"
		valueStrings = append(valueStrings, str_n)
		valueArgs = append(valueArgs, rand_subtrade_code)
		valueArgs = append(valueArgs, new_trade.Code)
		valueArgs = append(valueArgs, new_trade.UserWallet)
		valueArgs = append(valueArgs, subtrade.CreatedAt)
		valueArgs = append(valueArgs, subtrade.Type)
		valueArgs = append(valueArgs, subtrade.Reason)
		valueArgs = append(valueArgs, subtrade.Quantity)
		valueArgs = append(valueArgs, subtrade.AvgPrice)
		valueArgs = append(valueArgs, subtrade.Total)
		valueArgs = append(valueArgs, timeNow)
	}

	smt := fmt.Sprintf(subtrade_sql, strings.Join(valueStrings, ","))

	_, err = Db.Exec(smt, valueArgs...)
	if err != nil {
		log.WithFields(log.Fields{
			"subtradepayload": new_trade.Subtrades,
			"custom_msg":      "Failed inserting new subtrade",
		}).Error(err)
		return
	}
	return
}

func CreateSubtrade(w http.ResponseWriter, r *http.Request) {

	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed creating subtrade, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tradecode := mux.Vars(r)["tradecode"]
	if tradecode == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed creating subtrade, empty tradecode",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user_timezone, _ := time.LoadLocation(session.Timezone)
	now := time.Now().In(user_timezone)

	subtrade_sql := `
		INSERT INTO subtrades (
			code, tradecode, userwallet, 
			createdat, type, reason, 
			quantity, avgprice, total, updatedat)
		VALUES (
			SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2,
			$3, 'BUY', '', 0.0001, 0.0001,
			0.0001, $3)
		RETURNING code;`
	var subtrade_code string
	err = Db.QueryRow(
		subtrade_sql,
		tradecode,
		session.UserWallet,
		now,
	).Scan(&subtrade_code)
	if err != nil || subtrade_code == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed creating subtrade, wrong query",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte(subtrade_code))
}

func UpdateSubtrade(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed updating subtrade, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	subtrade := struct {
		Code      string      `json:"Code"`
		TradeCode string      `json:"TradeCode"`
		CreatedAt string      `json:"CreatedAt"`
		Type      string      `json:"Type"`
		Reason    string      `json:"Reason"`
		Quantity  json.Number `json:"Quantity"`
		AvgPrice  json.Number `json:"AvgPrice"`
		Total     json.Number `json:"Total"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&subtrade)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed creating subtrade, empty payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var subtrade_code string
	err = Db.QueryRow(`
		UPDATE subtrades
		SET
			createdat = $1,
			type = $2,
			reason = $3,
			quantity = $4,
			avgprice = $5,
			total = $6
		WHERE code = $7
		RETURNING code;`,
		subtrade.CreatedAt,
		subtrade.Type,
		subtrade.Reason,
		subtrade.Quantity,
		subtrade.AvgPrice,
		subtrade.Total,
		subtrade.Code).Scan(&subtrade_code)
	if err != nil || subtrade_code == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed creating subtrade, wrong query",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteSubtrade(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed deleting subtrade, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	subtrade_code := mux.Vars(r)["subtradecode"]
	if subtrade_code == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed deleting subtrade, empty subtradecode",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Db.Exec(`
		DELETE FROM subtrades
		WHERE code = $1;
		`, subtrade_code)

	w.WriteHeader(http.StatusOK)
}
