package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

func CheckUserPermissions(next http.Handler) http.Handler {
	fmt.Println(Gray(8-1, "Starting CheckUserPermissions..."))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]
		userA := UserByUsername(username)
		permission := userA.Permission
		if permission == "all" {
			next.ServeHTTP(w, r)
		} else {
			session := SelectSession(r)
			if session.Id == 0 {
				w.Write([]byte("Denied, need login"))
				return
			}
			userB := UserByEmail(session.Email)
			switch permission {
			case "private":
				w.Write([]byte("Denied, user private"))
				return
			case "followers":
				var isfollower bool
				_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE usera = $1
					AND userb = $2;`, userA.Id, userB.Id).Scan(
					&isfollower,
				)
				if isfollower {
					next.ServeHTTP(w, r)
				} else {
					w.Write([]byte("Denied, need follow"))
					return
				}
			case "subscribers":
				var issubscriber bool
				_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE usera = $1
					AND userb = $2;`, userA.Id, userB.Id).Scan(
					&issubscriber,
				)
				if issubscriber {
					next.ServeHTTP(w, r)
				} else {
					w.Write([]byte("Denied, need subscribe"))
					return
				}
			case "individuals":
				var isindividual bool
				_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM individuals
					WHERE usera = $1
					AND userb = $2;`, userA.Id, userB.Id).Scan(
					&isindividual,
				)
				if isindividual {
					next.ServeHTTP(w, r)
				} else {
					w.Write([]byte("Denied, need individual"))
					return
				}
			}
		}
	})
}

func SelectTrades(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectTrades..."))

	isopen := mux.Vars(r)["isopen"]
	username := mux.Vars(r)["username"]

	user := UserByUsername(username)

	type Subtrade struct {
		Id        int
		Timestamp string
		Type      string
		Reason    string
		Quantity  float64
		AvgPrice  float64
		Total     float64
	}

	type Trade struct {
		Id               string
		Exchange         string
		FirstPairId      int
		SecondPairId     int
		FirstPairName    string
		SecondPairName   string
		FirstPairSymbol  string
		SecondPairSymbol string
		FirstPairPrice   float64
		SecondPairPrice  float64
		QtyBuys          float64
		QtySells         float64
		TotalBuys        float64
		TotalSells       float64
		QtyAvailable     float64
		CurrentPrice     float64
		ActualReturn     float64
		FutureReturn     float64
		TotalReturn      float64
		ReturnBtc        float64
		ReturnUsd        float64
		Roi              float64
		Subtrades        []Subtrade
		BtcPrice         float64
	}

	trades := []Trade{}

	trades_sql := `
		WITH
			CURRENT_PRICE AS (
				SELECT
					p.coinid,
					c.name,
					c.symbol,
					p.price
				FROM prices p
				LEFT JOIN coins c ON(p.coinid = c.coinid)
				WHERE createdat = (SELECT MAX(createdat) FROM prices)),
			TRADES_MACRO AS (
				SELECT
					t.id,
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
				LEFT JOIN subtrades s ON(t.id  = s.tradeid)
				WHERE t.userid = $1
				AND t.isopen = $2
				GROUP BY 1, 2, 3, 4),
			TRADES_MICRO AS (
				SELECT
					t.id,
					t.exchange,
					t.firstpair AS firstpairid,
					c1.name AS firstpairname,
					c1.symbol AS firstpairsymbol,
					c1.price AS firstpairprice,
					t.secondpair AS secondpairid,
					c2.name AS secondpairname,
					c2.symbol AS secondpairsymbol,
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
				LEFT JOIN CURRENT_PRICE c1 ON(t.firstpair = c1.coinid)
				LEFT JOIN CURRENT_PRICE c2 ON(t.secondpair = c2.coinid))
		SELECT
			t.id,
			t.exchange,
			t.firstpairid,
			t.firstpairname,
			t.firstpairsymbol,
			t.firstpairprice,
			t.secondpairid,
			t.secondpairname,
			t.secondpairsymbol,
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
		LEFT JOIN CURRENT_PRICE c3 ON(c3.coinid = 1);`

	trades_rows, err := DbWebApp.Query(
		trades_sql,
		user.Id,
		isopen)
	defer trades_rows.Close()
	if err != nil {
		panic(err.Error())
	}
	for trades_rows.Next() {
		trade := Trade{}
		if err = trades_rows.Scan(
			&trade.Id,
			&trade.Exchange,
			&trade.FirstPairId,
			&trade.FirstPairName,
			&trade.FirstPairSymbol,
			&trade.FirstPairPrice,
			&trade.SecondPairId,
			&trade.SecondPairName,
			&trade.SecondPairSymbol,
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
				id,
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
		subtrades_rows, err := DbWebApp.Query(
			subtrades_sql,
			trade.Id)
		defer subtrades_rows.Close()
		if err != nil {
			panic(err.Error())
		}

		for subtrades_rows.Next() {
			subtrade := Subtrade{}
			if err = subtrades_rows.Scan(
				&subtrade.Id,
				&subtrade.Type,
				&subtrade.Reason,
				&subtrade.Timestamp,
				&subtrade.Quantity,
				&subtrade.AvgPrice,
				&subtrade.Total); err != nil {
				panic(err.Error())
			}

			subtrades = append(subtrades, subtrade)
		}

		trade.Subtrades = subtrades
		trades = append(trades, trade)
	}

	json.NewEncoder(w).Encode(trades)
}

func InsertTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting InsertTrade..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	trade := struct {
		Exchange     string `json:"Exchange"`
		FirstPairId  int    `json:"FirstPair"`
		SecondPairId int    `json:"SecondPair"`
		Subtrades    []struct {
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

	var trade_id string
	trade_sql := `
		INSERT INTO trades (id, userid, exchange, firstpair, secondpair, createdat, updatedat, isopen)
		VALUES (SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4, current_timestamp, current_timestamp, true)
		RETURNING id;`
	err = DbWebApp.QueryRow(
		trade_sql,
		user.Id,
		trade.Exchange,
		trade.FirstPairId,
		trade.SecondPairId,
	).Scan(&trade_id)
	if err != nil {
		panic(err.Error())
	}

	for _, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (tradeid, tradetimestamp, type, reason, quantity, avgprice, total, createdat, updatedat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp)`
		DbWebApp.Exec(
			subtrade_sql,
			trade_id,
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

func UpdateTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateTrade..."))

	_ = SelectSession(r)

	trade := struct {
		Id           string `json:"Id"`
		Exchange     string `json:"Exchange"`
		FirstPairId  int    `json:"FirstPairId"`
		SecondPairId int    `json:"SecondPairId"`
		Subtrades    []struct {
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

	DbWebApp.Exec(`
		DELETE FROM subtrades
		WHERE tradeid = $1;
	`, trade.Id)

	for _, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (tradeid, tradetimestamp, type, reason, quantity, avgprice, total, createdat, updatedat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp)`
		DbWebApp.Exec(
			subtrade_sql,
			trade.Id,
			subtrade.Timestamp,
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

	_ = SelectSession(r)

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.Exec(`
		UPDATE trades
		SET isopen = False
		WHERE id = $1;
		`, tradeid)

	json.NewEncoder(w).Encode("OK")
}

func OpenTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting OpenTrade..."))

	_ = SelectSession(r)

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.Exec(`
		UPDATE trades
		SET isopen = True
		WHERE id = $1;
		`, tradeid)

	json.NewEncoder(w).Encode("OK")
}

func DeleteTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting DeleteTrade..."))

	_ = SelectSession(r)

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.Exec(`
		DELETE FROM subtrades
		WHERE tradeid IN (
			SELECT id
			FROM trades
			WHERE id = $1);
		`, tradeid)

	DbWebApp.Exec(`
		DELETE FROM trades
		WHERE id = $1;
		`, tradeid)

	json.NewEncoder(w).Encode("OK")
}
