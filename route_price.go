package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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

func GetPriceSocket(w http.ResponseWriter, r *http.Request) {

	str_trade_id := mux.Vars(r)["tradeid"]
	trade_id, err := strconv.Atoi(str_trade_id)

	str_first_pair_id := mux.Vars(r)["firstpairid"]
	first_pair_id, err := strconv.Atoi(str_first_pair_id)

	str_second_pair_id := mux.Vars(r)["secondpairid"]
	second_pair_id, err := strconv.Atoi(str_second_pair_id)

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

	go func() {
		for {
			price := GetPrice(first_pair_id, second_pair_id)
			price = price + rand.Float64()*(1-0.01)/10000
			response := struct {
				TradeId      int
				FirstPairId  int
				SelectPairId int
				Price        float64
			}{
				trade_id,
				first_pair_id,
				second_pair_id,
				price,
			}
			err = ws.WriteJSON(response)
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

func GetPrice(first_pair_id int, second_pair_id int) (price float64) {
	fmt.Println(Gray(8-1, "Get SelectPrice..."))

	price_sql := `
		SELECT
			p2.price / p1.price AS price
		FROM (
				SELECT
					price
				FROM prices
				WHERE coinid = $1
				AND createdat = (
					SELECT MAX(createdat)
					FROM prices)) p1
		LEFT JOIN (
				SELECT
					price
				FROM prices
				WHERE coinid = $2
				AND createdat = (
					SELECT MAX(createdat)
					FROM prices)) p2
			ON(1 = 1);`

	err := DbWebApp.QueryRow(
		price_sql,
		first_pair_id,
		second_pair_id,
	).Scan(&price)
	if err != nil {
		log.Error(err)
	}

	return
}
