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

type NewTrade struct {
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
		Usercode  string
	} `json:"subtrades"`
	Usercode string
	Code     string
}

func CreateTrade(w http.ResponseWriter, r *http.Request) {
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

	new_trade := NewTrade{
		Usercode: user.Code,
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&new_trade)
	if err != nil {
		log.Error(err)
	}

	new_trade.InsertTrade()
	new_trade.InsertSubTrades()

}

func (new_trade *NewTrade) InsertTrade() {
	trade_sql := `
		INSERT INTO trades (code, usercode, exchange, firstpair, secondpair, createdat, updatedat, isopen)
		VALUES (SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4, current_timestamp, current_timestamp, true)
		RETURNING code;`
	err := Db.QueryRow(
		trade_sql,
		new_trade.Usercode,
		new_trade.Exchange,
		new_trade.FirstPairId,
		new_trade.SecondPairId,
	).Scan(&new_trade.Code)
	if err != nil {
		log.Error(err)
	}
	return
}

func (new_trade *NewTrade) InsertSubTrades() {
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
		valueArgs = append(valueArgs, RandStringBytes(12))
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
	_, err := Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
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
		WHERE code = $1;
		`, tradecode)

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
		WHERE code = $1;
		`, tradecode)

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
