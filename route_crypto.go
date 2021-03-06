package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func SelectPairs(w http.ResponseWriter, r *http.Request) {
	pairs := make(map[int]PairInfo)

	pairs_sql := `
		SELECT DISTINCT
			coinid,
			symbol,
			name,
			slug
		FROM coins
		ORDER BY coinid;`
	pairs_rows, err := Db.Query(pairs_sql)
	defer pairs_rows.Close()
	if err != nil {
		log.Error(err)
	}
	for pairs_rows.Next() {
		var coinid int
		pair_info := PairInfo{}
		if err = pairs_rows.Scan(
			&coinid,
			&pair_info.Symbol,
			&pair_info.Name,
			&pair_info.Slug,
		); err != nil {
			log.Error(err)
		}
		pairs[coinid] = pair_info
	}
	json.NewEncoder(w).Encode(pairs)
}

func SelectPairRatio(w http.ResponseWriter, r *http.Request) {
	first_coin_id := mux.Vars(r)["firstPairCoinId"]
	second_coin_id := mux.Vars(r)["secondPairCoinId"]

	if first_coin_id == "" || second_coin_id == "" {
		log.WithFields(log.Fields{
			"firstPairCoinId":  first_coin_id,
			"secondPairCoinId": second_coin_id,
		}).Error("Failed extracting pair ratio, empty value")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ratio_sql := `
		SELECT
			y.price / x.price
		FROM (
			SELECT
				ROUND(price, 6) AS price
			FROM lastprices
			WHERE coinid = $1) x
		LEFT JOIN (
			SELECT
				ROUND(price, 6) AS price
			FROM lastprices
			WHERE coinid = $2) y ON(1=1);`

	var pair_ratio float64
	err := Db.QueryRow(
		ratio_sql,
		first_coin_id,
		second_coin_id).Scan(&pair_ratio)
	if err != nil {
		log.WithFields(log.Fields{
			"firstPairCoinId":  first_coin_id,
			"secondPairCoinId": second_coin_id,
			"customMsg":        "Failed extracting pair ratio",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(pair_ratio)
}
