package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func SelectExplore(w http.ResponseWriter, r *http.Request) {
	offset_string := mux.Vars(r)["offset"]
	offset, err := strconv.Atoi(offset_string)
	if offset%10 != 0 || err != nil {
		log.Warn("Attempted accessing Explore with invalid offset")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	explore_sql := `
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
			SUBTRADES AS (
				SELECT
					'subtrade' AS eventtype,
					s.tradecode,
					s.updatedat,
					CASE
						WHEN EXTRACT(EPOCH FROM (NOW() - s.updatedat)) / 60 < 1
							THEN ROUND(EXTRACT(EPOCH FROM (NOW() - s.updatedat)))::TEXT || ' seconds ago'
						WHEN (EXTRACT(EPOCH FROM (NOW() - s.updatedat)) / 60 > 1) AND (EXTRACT(EPOCH FROM (NOW() - s.updatedat)) / 60 < 60)
							THEN ROUND(EXTRACT(EPOCH FROM (NOW() - s.updatedat)) / 60)::TEXT || ' minutes ago'
						WHEN (EXTRACT(EPOCH FROM (NOW() - s.updatedat)) / 60 >= 60) AND (EXTRACT(EPOCH FROM(NOW() - s.updatedat)) / 60 < 1440)  
							THEN ROUND(EXTRACT(EPOCH FROM (NOW() - s.updatedat)) / 60 / 60)::TEXT || ' hours ago'
						ELSE CEIL(EXTRACT(EPOCH FROM (NOW() - s.updatedat)) / 60 / 60 / 60)::TEXT || ' days ago'
					END AS timeago,
					s.userwallet,
					s.type,
					s.reason,
					CASE
						WHEN s.quantity > 100 THEN TO_CHAR(s.quantity, '999,999,999')
						WHEN s.quantity > 1 THEN RTRIM(RTRIM(TO_CHAR(s.quantity, '999,999,999.00'), '0'), '.')
						ELSE RTRIM(RTRIM(TO_CHAR(s.quantity, '999,999,999.00000'), '0'), '.')
					END as quantity,
					CASE
						WHEN s.avgprice > 100 THEN TO_CHAR(s.avgprice, '999,999,999')
						WHEN s.avgprice > 1 THEN TO_CHAR(s.avgprice, '999,999,999.00')
						ELSE TO_CHAR(s.avgprice, '999,999,999.00000')
					END as avgprice,
					CASE
						WHEN s.total > 100 THEN TO_CHAR(s.total, '999,999,999')
						WHEN s.total > 1 THEN TO_CHAR(s.total, '999,999,999.00')
						ELSE TO_CHAR(s.total, '999,999,999.00000')
					END as total,
					u.profilepicture,
					t.firstpair,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.firstpair::TEXT || '.png' AS firstpairurlicon,
					c1.symbol AS firstpairsymbol,
					t.secondpair,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.secondpair::TEXT || '.png' AS secondpairurlicon,
					c2.symbol AS secondpairsymbol,
					CASE
						WHEN (cp2.price / cp1.price)  > 100 THEN TO_CHAR((cp2.price / cp1.price), '999,999,999') 
						WHEN (cp2.price / cp1.price)  > 1 THEN TO_CHAR((cp2.price / cp1.price), '999,999,999.00') 
						ELSE TO_CHAR((cp2.price / cp1.price), '999,999,999.00000')
					END as currentprice,
					ROUND(((((cp2.price / cp1.price) / s.avgprice) - 1) * 100), 1) AS deltapriceperc
				FROM subtrades s
				LEFT JOIN trades t ON(s.tradecode = t.code)
				LEFT JOIN coins c1 ON(t.firstpair = c1.coinid)
				LEFT JOIN coins c2 ON(t.secondpair = c2.coinid)
				LEFT JOIN users u ON(s.userwallet = u.wallet)
				LEFT JOIN CURRENT_PRICE cp1 ON(t.firstpair = cp1.coinid)
				LEFT JOIN CURRENT_PRICE cp2 ON(t.secondpair = cp2.coinid)
				ORDER BY s.updatedat DESC
				LIMIT 10
				OFFSET $1),
			events AS (
				SELECT
					userwallet,
					updatedat,
					ROW_TO_JSON(subtrades) AS row_json
				FROM SUBTRADES)
		SELECT
			json_agg(row_json)
		FROM events c;`

	var explore_json string
	err = Db.QueryRow(explore_sql, offset).Scan(&explore_json)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("{}"))
		return
	}
	w.Write([]byte(explore_json))
}
