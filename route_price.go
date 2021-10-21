package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetPrices(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetPrices..."))

	upgrader.CheckOrigin = func(r *http.Request) bool {
		for _, origin := range Origins {
			if origin == r.Header["Origin"][0] {
				return true
			}
		}
		return false
	}

	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	usercode := mux.Vars(r)["usercode"]
	user := SelectUser("code", usercode)
	if user.Id == 0 {
		log.WithFields(log.Fields{
			"usercode": usercode,
		}).Warn("Unknow usercode, closing wb")
		return
	}

	type TradePrice struct {
		TradeId string
		Price   float64
	}

	prices_sql := `
		WITH
			latest_price AS (
				SELECT
					coinid,
					price
				FROM prices
				WHERE createdat = (
					SELECT
						MAX(createdat)
					FROM prices))
		SELECT 
			t.id,
			l2.price / l1.price AS price
		FROM users u
		LEFT JOIN trades t ON(u.id = t.userid)
		LEFT JOIN latest_price l1 ON(t.firstpair = l1.coinid)
		LEFT JOIN latest_price l2 ON(t.secondpair = l2.coinid)
		WHERE u.code = $1
		AND t.isopen = TRUE;`

	for {
		tradesprices := []TradePrice{}
		prices_rows, err := DbWebApp.Query(
			prices_sql,
			usercode)
		defer prices_rows.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"usercode":   usercode,
				"custom_msg": "Failed running prices query",
			}).Error(err)
			return
		}
		for prices_rows.Next() {
			tradeprice := TradePrice{}
			if err = prices_rows.Scan(
				&tradeprice.TradeId,
				&tradeprice.Price,
			); err != nil {
				log.WithFields(log.Fields{
					"usercode":   usercode,
					"custom_msg": "Failed scanning price into struct",
				}).Error(err)
				return
			}

			tradeprice.Price = tradeprice.Price + rand.Float64()*(1-0.01)/10000
			tradesprices = append(tradesprices, tradeprice)
		}
		err = ws.WriteJSON(tradesprices)
		if err != nil {
			log.WithFields(log.Fields{
				"usercode":   usercode,
				"custom_msg": "Failed returning pricing to ws",
			}).Error(err)
			return
		}

		time.Sleep(time.Duration(rand.Intn(20)+8) * time.Second)
	}
}
