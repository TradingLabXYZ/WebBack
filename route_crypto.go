package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func SelectPairs(w http.ResponseWriter, r *http.Request) {
	type PairInfo struct {
		CoinId int
		Name   string
		Slug   string
	}

	pairs := make(map[string]PairInfo)

	pairs_sql := `
		SELECT DISTINCT
			symbol,
			name,
			coinid,
			slug
		FROM coins
		ORDER BY 1;`
	pairs_rows, err := Db.Query(pairs_sql)
	defer pairs_rows.Close()
	if err != nil {
		log.Error(err)
	}
	for pairs_rows.Next() {
		var symbol string
		pair_info := PairInfo{}
		if err = pairs_rows.Scan(
			&symbol,
			&pair_info.Name,
			&pair_info.CoinId,
			&pair_info.Slug,
		); err != nil {
			log.Error(err)
		}
		pairs[symbol] = pair_info
	}

	json.NewEncoder(w).Encode(pairs)
}

func SelectPairRatio(w http.ResponseWriter, r *http.Request) {
	first_coin_id := mux.Vars(r)["firstPairCoinId"]
	second_coin_id := mux.Vars(r)["secondPairCoinId"]

	ratio_sql := `
		SELECT
			y.price / x.price
		FROM (
			SELECT
				ROUND(price, 6) AS price
			FROM prices
			WHERE coinid = $1
			ORDER BY createdat
			DESC LIMIT 1) x
		LEFT JOIN (
			SELECT
				ROUND(price, 6) AS price
			FROM prices
			WHERE coinid = $2
			ORDER BY createdat DESC
			LIMIT 1) y ON(1=1);`

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
