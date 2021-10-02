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

	username := mux.Vars(r)["username"]

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, _ := upgrader.Upgrade(w, r, nil)

	closing_message := make(chan string)

	go func() {
		_, p, err := ws.ReadMessage()
		if string(p) != "" {
			closing_message <- "close"
		}
		if err != nil {
			log.Println(err)
			return
		}
	}()

	type TradePrice struct {
		TradeId string
		Price   float64
	}

	go func() {
		for {
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
				WHERE u.username = $1
				AND t.isopen = TRUE;`

			tradesprices := []TradePrice{}
			prices_rows, err := DbWebApp.Query(
				prices_sql,
				username)
			defer prices_rows.Close()
			if err != nil {
				log.Error(err)
			}
			for prices_rows.Next() {
				tradeprice := TradePrice{}
				if err = prices_rows.Scan(
					&tradeprice.TradeId,
					&tradeprice.Price,
				); err != nil {
					log.Error(err)
				}

				tradeprice.Price = tradeprice.Price + rand.Float64()*(1-0.01)/100000
				tradesprices = append(tradesprices, tradeprice)
			}
			err = ws.WriteJSON(tradesprices)
			if err != nil {
				return
			}
			select {
			case _ = <-closing_message:
				ws.Close()
				return
			default:
			}
			time.Sleep(time.Duration(rand.Intn(10)+2) * time.Second)
		}
	}()
}
