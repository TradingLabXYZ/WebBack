package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

func InsertTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting InsertTrade..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	trade := struct {
		Exchange   string `json:"Exchange"`
		FirstPair  string `json:"FirstPair"`
		SecondPair string `json:"SecondPair"`
		Subtrades  []struct {
			Timestamp string      `json:"Timestamp"`
			Type      string      `json:"Type"`
			Reason    string      `json:"Reason"`
			Quantity  json.Number `json:"Quantity"`
			AvgPrice  json.Number `json:"AvgPrice"`
			Total     json.Number `json:"Total"`
		} `json:"subtrades"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&trade)
	if err != nil {
		panic(err)
	}

	var next_user_trade int
	next_trade_sql := `
		SELECT
			CASE
				WHEN MAX(usertrade) + 1 IS NULL THEN 1
				ELSE MAX(usertrade) + 1
			END
		FROM trades
		WHERE userid = $1;`

	err = DbWebApp.QueryRow(next_trade_sql, user.Id).Scan(&next_user_trade)
	if err != nil {
		panic(err.Error())
	}

	var trade_id int
	trade_sql := `
		INSERT INTO trades (userid, usertrade, exchange, firstpair, secondpair, createdat, updatedat, isopen)
		VALUES ($1, $2, $3, $4, $5, current_timestamp, current_timestamp, true)
		RETURNING id;`
	err = DbWebApp.QueryRow(
		trade_sql,
		user.Id,
		next_user_trade,
		trade.Exchange,
		trade.FirstPair,
		trade.SecondPair,
	).Scan(&trade_id)
	if err != nil {
		panic(err.Error())
	}

	for i, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (tradeid, subtradeid, tradetimestamp, type, reason, quantity, avgprice, total, createdat, updatedat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, current_timestamp, current_timestamp)`
		_, err = DbWebApp.Exec(
			subtrade_sql,
			trade_id,
			i+1,
			subtrade.Timestamp,
			subtrade.Type,
			subtrade.Reason,
			subtrade.Quantity,
			subtrade.AvgPrice,
			subtrade.Total,
		)
	}
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode("OK")
}

func SelectTrades(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectTrades..."))

	isopen := mux.Vars(r)["isopen"]
	username := mux.Vars(r)["username"]

	_ = SelectSession(r)

	user := UserByUsername(username)

	type Subtrade struct {
		SubtradeId int
		Timestamp  string
		Type       string
		Reason     string
		Quantity   float64
		AvgPrice   float64
		Total      float64
	}

	type Trade struct {
		Id              int
		Usertrade       int
		Exchange        string
		FirstPair       string
		SecondPair      string
		FirstPairPrice  float64
		SecondPairPrice float64
		QtyBuys         float64
		QtySells        float64
		TotalBuys       float64
		TotalSells      float64
		QtyAvailable    float64
		CurrentPrice    float64
		ActualReturn    float64
		FutureReturn    float64
		TotalReturn     float64
		ReturnBtc       float64
		ReturnUsd       float64
		Roi             float64
		Subtrades       []Subtrade
		BtcPrice        float64
	}

	trades := []Trade{}

	trades_sql := `
		WITH
			CURRENT_PRICE AS (
				SELECT
					symbol,
					price
				FROM coinmarketcap
				WHERE createdat = (SELECT MAX(createdat) FROM coinmarketcap)),
			TRADES_MACRO AS (
				SELECT
					t.id,
					t.usertrade,
					t.exchange,
					t.firstpair,
					t.secondpair,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'BUY' THEN s.quantity END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'BUY' THEN s.quantity END) 
					END AS qtybuys,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'SELL' THEN s.quantity END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'SELL' THEN s.quantity END) 
					END AS qtysells,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'BUY' THEN s.total END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'BUY' THEN s.total END) 
					END AS totalbuys,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'SELL' THEN s.total END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'SELL' THEN s.total END) 
					END AS totalsells
				FROM trades t 
				LEFT JOIN subtrades s ON(t.usertrade  = s.tradeid)
				WHERE t.userid = $1
				AND t.isopen = $2
				GROUP BY 1, 2, 3, 4, 5),
			TRADES_MICRO AS (
				SELECT
					t.id,
					t.usertrade,
					t.exchange,
					t.firstpair,
					c1.price AS firstpairprice,
					t.secondpair,
					c2.price AS secondpairprice,
					t.qtybuys,
					t.qtysells,
					t.totalbuys,
					t.totalsells,
					t.qtybuys - t.qtysells AS qtyavailable,
					(c2.price / c1.price) AS currentprice,
					t.totalsells - t.totalbuys AS actualreturn,
					(t.qtybuys - t.qtysells) * (c2.price / c1.price) AS futurereturn,
					t.totalsells - t.totalbuys + (t.qtybuys - t.qtysells) * (c2.price / c1.price) AS totalreturn,
					CASE
						WHEN t.totalbuys = 0 THEN 0
						ELSE (((t.qtybuys - t.qtysells) * (c2.price / c1.price) + t.totalsells) / t.totalbuys - 1) * 100
					END AS roi
				FROM TRADES_MACRO t
				LEFT JOIN CURRENT_PRICE c1 ON(t.firstpair = c1.symbol)
				LEFT JOIN CURRENT_PRICE c2 ON(t.secondpair = c2.symbol))
		SELECT
			t.id,
			t.usertrade,
			t.exchange,
			t.firstpair,
			t.firstpairprice,
			t.secondpair,
			t.secondpairprice,
			t.qtybuys,
			t.qtysells,
			t.totalbuys,
			t.totalsells,
			t.qtyavailable,
			t.currentprice,
			t.actualreturn,
			t.futurereturn,
			t.totalreturn,
			t.totalreturn * t.firstpairprice / c3.price as returnbtc,
			t.totalreturn * t.firstpairprice as returnusd,
			t.roi,
			c3.price AS btcprice
		FROM TRADES_MICRO t
		LEFT JOIN CURRENT_PRICE c3 ON(c3.symbol = 'BTC');`

	trades_rows, err := DbWebApp.Query(trades_sql, user.Id, isopen)
	if err != nil {
		panic(err.Error())
	}
	for trades_rows.Next() {
		trade := Trade{}
		if err = trades_rows.Scan(
			&trade.Id,
			&trade.Usertrade,
			&trade.Exchange,
			&trade.FirstPair,
			&trade.FirstPairPrice,
			&trade.SecondPair,
			&trade.SecondPairPrice,
			&trade.QtyBuys,
			&trade.QtySells,
			&trade.TotalBuys,
			&trade.TotalSells,
			&trade.QtyAvailable,
			&trade.CurrentPrice,
			&trade.ActualReturn,
			&trade.FutureReturn,
			&trade.TotalReturn,
			&trade.ReturnBtc,
			&trade.ReturnUsd,
			&trade.Roi,
			&trade.BtcPrice,
		); err != nil {
			panic(err)
		}

		subtrades_sql := `
			SELECT
				subtradeid,
				type,
				reason,
				TO_CHAR(tradetimestamp, 'YYYY-MM-DD"T"HH24:MI'),
				quantity,
				avgprice,
				total
			FROM subtrades
			WHERE tradeid = $1
			ORDER BY 1;`

		subtrades := []Subtrade{}
		subtrades_rows, err := DbWebApp.Query(subtrades_sql, trade.Usertrade)
		if err != nil {
			panic(err.Error())
		}

		for subtrades_rows.Next() {
			subtrade := Subtrade{}
			if err = subtrades_rows.Scan(
				&subtrade.SubtradeId,
				&subtrade.Type,
				&subtrade.Reason,
				&subtrade.Timestamp,
				&subtrade.Quantity,
				&subtrade.AvgPrice,
				&subtrade.Total); err != nil {
				return
			}
			subtrades = append(subtrades, subtrade)
		}

		trade.Subtrades = subtrades
		trades = append(trades, trade)
	}

	json.NewEncoder(w).Encode(trades)
}

func UpdateTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateTrade..."))

	_ = SelectSession(r)

	trade := struct {
		Id         int    `json:"Id"`
		Exchange   string `json:"Exchange"`
		FirstPair  string `json:"FirstPair"`
		SecondPair string `json:"SecondPair"`
		Subtrades  []struct {
			Timestamp string      `json:"Timestamp"`
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
		panic(err)
	}

	DbWebApp.QueryRow("DELETE FROM subtrades WHERE tradeid = $1;", trade.Id)

	for i, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (tradeid, subtradeid, tradetimestamp, type, reason, quantity, avgprice, total, createdat, updatedat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, current_timestamp, current_timestamp)`
		_, err = DbWebApp.Exec(
			subtrade_sql,
			trade.Id,
			i+1,
			subtrade.Timestamp,
			subtrade.Type,
			subtrade.Reason,
			subtrade.Quantity,
			subtrade.AvgPrice,
			subtrade.Total,
		)
	}
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode("OK")
}

func CloseTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting CloseTrade..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.QueryRow(`
		UPDATE trades
		SET isopen = False
		WHERE userid = $1
		AND usertrade = $2;`, user.Id, tradeid)

	json.NewEncoder(w).Encode("OK")
}

func OpenTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting OpenTrade..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.QueryRow(`
		UPDATE trades
		SET isopen = True
		WHERE userid = $1
		AND usertrade = $2;`, user.Id, tradeid)

	json.NewEncoder(w).Encode("OK")
}

func DeleteTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting DeleteTrade..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.QueryRow(`
		DELETE FROM subtrades
		WHERE tradeid IN (
			SELECT id
			FROM trades
			WHERE userid = $1
			AND usertrade = $2);`, user.Id, tradeid)

	DbWebApp.QueryRow(`
		DELETE FROM trades
		WHERE userid = $1
		AND usertrade = $2;`, user.Id, tradeid)

	json.NewEncoder(w).Encode("OK")
}
