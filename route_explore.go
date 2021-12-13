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
					code,
					createdat,
					EXTRACT(EPOCH FROM (NOW() - createdat)) / 60 AS minuteago,
					userwallet,
					exchange,
					firstpair,		
					secondpair
				FROM trades),
			subtrades AS (
				SELECT
					'subtrade' AS eventtype,
					tradecode,
					createdat,
					EXTRACT(EPOCH FROM (NOW() - createdat)) / 60 AS minuteago,
					userwallet,
					type,
					reason,
					quantity,
					avgprice,
					total
				FROM subtrades),
			trade_subtrade AS (
				SELECT
					createdat,
					ROW_TO_JSON(trades) AS row_json
				FROM trades

				UNION ALL

				SELECT
					createdat,
					ROW_TO_JSON(subtrades) AS row_json
				FROM subtrades

				ORDER BY createdat DESC)
		SELECT
			json_agg(row_json)
		FROM trade_subtrade c
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
