package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

func SelectPrice(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectPrice..."))

	str_first_pair_id := mux.Vars(r)["firstpairid"]
	first_pair_id, err := strconv.Atoi(str_first_pair_id)

	str_second_pair_id := mux.Vars(r)["secondpairid"]
	second_pair_id, err := strconv.Atoi(str_second_pair_id)

	_ = SelectSession(r)

	price_sql := `
		SELECT
			p2.price / p1.price AS price
		FROM (
				SELECT
					price
				FROM prices
				WHERE coinid = $1
				AND createdat = (SELECT MAX(createdat) FROM prices)) p1
		LEFT JOIN (
				SELECT
					price
				FROM prices
				WHERE coinid = $2
				AND createdat = (SELECT MAX(createdat) FROM prices)) p2
			ON(1 = 1);`

	var price float64
	err = DbWebApp.QueryRow(price_sql, first_pair_id, second_pair_id).Scan(&price)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(price)
}

func SelectPairs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectPairs..."))

	_ = SelectSession(r)

	type PairInfo struct {
		CoinId int
		Symbol string
		Slug   string
	}

	pairs := make(map[string]PairInfo)

	pairs_sql := `
		SELECT DISTINCT
			name,
			coinid,
			symbol,
			slug
		FROM coins;
		`
	pairs_rows, err := DbWebApp.Query(pairs_sql)
	if err != nil {
		panic(err.Error())
	}
	for pairs_rows.Next() {
		var name string
		pair_info := PairInfo{}
		if err = pairs_rows.Scan(
			&name,
			&pair_info.CoinId,
			&pair_info.Symbol,
			&pair_info.Slug,
		); err != nil {
			panic(err)
		}
		pairs[name] = pair_info
	}

	json.NewEncoder(w).Encode(pairs)
}
