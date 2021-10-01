package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

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
		FROM coins;`
	pairs_rows, err := DbWebApp.Query(pairs_sql)
	defer pairs_rows.Close()
	if err != nil {
		log.Error(err)
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
			log.Error(err)
		}
		pairs[name] = pair_info
	}

	json.NewEncoder(w).Encode(pairs)
}
