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
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

func (new_trade *NewTrade) InsertSubTrades() (err error) {
	subtrade_sql := `
		INSERT INTO subtrades (
			code, tradecode, usercode, 
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
		valueArgs = append(valueArgs, new_trade.Usercode)
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
			"custom_msg": "Failed inserting new subtrade",
		}).Error(err)
		return
	}
	return
}

func CreateSubtrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting CreateUpdate..."))

	tradecode := mux.Vars(r)["tradecode"]

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

	subtrade_sql := `
		INSERT INTO subtrades (
			code, tradecode, usercode, 
			createdat, type, reason, 
			quantity, avgprice, total, updatedat)
		VALUES (
			SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2,
			current_timestamp, 'BUY', '', 0, 0,
			0, current_timestamp)`
	Db.Exec(
		subtrade_sql,
		tradecode,
		user.Code,
	)
	if err != nil {
		log.Error(err)
	}
	return
}

func UpdateSubtrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateSubtrade..."))

	_, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	subtrade := struct {
		Code      string      `json:"Code"`
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
		log.Error(err)
	}

	Db.Exec(`
		UPDATE subtrades
		SET
			createdat = $1,
			type = $2,
			reason = $3,
			quantity = $4,
			avgprice = $5,
			total = $6
		WHERE code = $7;`,
		subtrade.CreatedAt,
		subtrade.Type,
		subtrade.Reason,
		subtrade.Quantity,
		subtrade.AvgPrice,
		subtrade.Total,
		subtrade.Code,
	)

	json.NewEncoder(w).Encode("OK")
}

func DeleteSubtrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting DeleteSubTrade..."))

	_, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	subtrade_code := mux.Vars(r)["subtradecode"]

	Db.Exec(`
		DELETE FROM subtrades
		WHERE code = $1;
		`, subtrade_code)

	json.NewEncoder(w).Encode("OK")
}
