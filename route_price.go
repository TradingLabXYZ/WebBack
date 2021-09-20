package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

func SelectPrice(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectPrice..."))

	first_pair := mux.Vars(r)["firstpair"]
	second_pair := mux.Vars(r)["secondpair"]

	_ = SelectSession(r)

	price_sql := `
		SELECT
			p2.price / p1.price AS price
		FROM (
				SELECT
					price
				FROM coinmarketcap
				WHERE symbol = $1
				AND createdat = (SELECT MAX(createdat) FROM coinmarketcap)) p1
		LEFT JOIN (
				SELECT
					price
				FROM coinmarketcap
				WHERE symbol = $2
				AND createdat = (SELECT MAX(createdat) FROM coinmarketcap)) p2
			ON(1 = 1);`

	var price float64
	err := DbWebApp.QueryRow(price_sql, first_pair, second_pair).Scan(&price)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(price)
}

func SelectPairs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectPairs..."))

	_ = SelectSession(r)

	pairs_sql := `
		SELECT
			ARRAY_AGG(DISTINCT symbol || ' - ' || name) AS pair
		FROM coinmarketcap
		ORDER BY 1;
		`
	var pairs []string
	err := DbWebApp.QueryRow(pairs_sql).Scan(pq.Array(&pairs))
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(pairs)
}
