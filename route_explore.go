package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func SelectExplore(w http.ResponseWriter, r *http.Request) {
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
					s.createdat,
					CASE
						WHEN EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60 < 1
							THEN ROUND(EXTRACT(EPOCH FROM (NOW() - t.createdat)))::TEXT || ' seconds ago'
						WHEN (EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60 > 1) AND (EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60 < 60)
							THEN ROUND(EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60)::TEXT || ' minutes ago'
						WHEN (EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60 >= 60) AND (EXTRACT(EPOCH FROM(NOW() - t.createdat)) / 60 < 1440)  
							THEN ROUND(EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60 / 60)::TEXT || ' hours ago'
						ELSE ROUND(EXTRACT(EPOCH FROM (NOW() - t.createdat)) / 60 / 60 / 60)::TEXT || ' days ago'
					END AS timeago,
					s.userwallet,
					s.type,
					s.reason,
					s.quantity,
					CASE WHEN s.avgprice > 1 THEN ROUND(s.avgprice, 2) ELSE ROUND(s.avgprice, 5) END as avgprice,
					CASE WHEN s.total > 1 THEN ROUND(s.total, 2) ELSE ROUND(s.total, 5) END as total,
					u.profilepicture,
					t.firstpair,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.firstpair::TEXT || '.png' AS firstpairurlicon,
					c1.symbol AS firstpairsymbol,
					t.secondpair,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.secondpair::TEXT || '.png' AS secondpairurlicon,
					c2.symbol AS secondpairsymbol,
					CASE
						WHEN (cp2.price / cp1.price)  > 1 THEN ROUND((cp2.price / cp1.price) , 2)
						ELSE ROUND((cp2.price / cp1.price) , 5)
					END as currentprice,
					ROUND(((((cp2.price / cp1.price) / s.avgprice) - 1) * 100), 1) AS deltapriceperc
				FROM subtrades s
				LEFT JOIN trades t ON(s.tradecode = t.code)
				LEFT JOIN coins c1 ON(t.firstpair = c1.coinid)
				LEFT JOIN coins c2 ON(t.secondpair = c2.coinid)
				LEFT JOIN users u ON(s.userwallet = u.wallet)
				LEFT JOIN CURRENT_PRICE cp1 ON(t.firstpair = cp1.coinid)
				LEFT JOIN CURRENT_PRICE cp2 ON(t.secondpair = cp2.coinid)
				LIMIT 20),
			events AS (
				SELECT
					userwallet,
					createdat,
					ROW_TO_JSON(subtrades) AS row_json
				FROM SUBTRADES
				ORDER BY createdat DESC)
		SELECT
			json_agg(row_json)
		FROM events c;`

	var explore_json string
	err := Db.QueryRow(explore_sql).Scan(&explore_json)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed sending explore snapshot",
		}).Error(err.Error())
	}
	w.Write([]byte(explore_json))
}
