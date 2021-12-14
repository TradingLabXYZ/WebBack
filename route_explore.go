package main

import (
	"fmt"
	"net/http"
)

func SelectExplore(w http.ResponseWriter, r *http.Request) {
	explore_sql := `
		WITH
			trades AS (
				SELECT
					'trade' AS eventtype,
					t.code,
					t.createdat,
					ROUND(EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60) AS minuteago,
					t.userwallet,
					t.exchange,
					t.firstpair,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.firstpair::TEXT || '.png' AS firstpairurlicon,
					c1.symbol AS firstpairsymbol,
					t.secondpair,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.secondpair::TEXT || '.png' AS secondpairurlicon,
					c2.symbol AS secondpairsymbol,
					u.profilepicture
				FROM trades t
				LEFT JOIN users u ON(t.userwallet = u.wallet)
				LEFT JOIN coins c1 ON(t.firstpair = c1.coinid)
				LEFT JOIN coins c2 ON(t.secondpair = c2.coinid)),
			subtrades AS (
				SELECT
					'subtrade' AS eventtype,
					s.tradecode,
					s.createdat,
					ROUND(EXTRACT(EPOCH FROM (NOW() - s.createdat)) / 60) AS minuteago,
					s.userwallet,
					s.type,
					s.reason,
					s.quantity,
					s.avgprice,
					s.total,
					u.profilepicture,
					t.secondpair,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.secondpair::TEXT || '.png' AS secondpairurlicon,
					c.symbol AS secondpairsymbol
				FROM subtrades s
				LEFT JOIN trades t ON(s.tradecode = t.code)
				LEFT JOIN coins c ON(t.secondpair = c.coinid)
				LEFT JOIN users u ON(s.userwallet = u.wallet)),
			events AS (
				SELECT
					userwallet,
					createdat,
					ROW_TO_JSON(trades) AS row_json
				FROM trades
				UNION ALL
				SELECT
					userwallet,
					createdat,
					ROW_TO_JSON(subtrades) AS row_json
				FROM subtrades
				ORDER BY createdat DESC)
		SELECT
			json_agg(row_json)
		FROM events c
		LIMIT 1;`

	var explore_json string
	err := Db.QueryRow(explore_sql).Scan(&explore_json)
	if err != nil {
		// manage error
		fmt.Println(err)
	}
	w.Write([]byte(explore_json))
	// json.NewEncoder(w).Encode(explore_json)
}
