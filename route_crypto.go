package main

import (
	"encoding/json"
	"net/http"

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
