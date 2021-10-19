package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

type WsTrade struct {
	Channel   chan string
	Username  string
	RequestId string
}

type WsTradeOutput struct {
	UserDetails UserDetails
	Trades      []Trade
}

type UserDetails struct {
	Username string
	Twitter  string
}

type Subtrade struct {
	Id        int
	Timestamp string
	Type      string
	Reason    string
	Quantity  float64
	AvgPrice  float64
	Total     float64
}

type Trade struct {
	Id               string
	IsOpen           string
	Exchange         string
	FirstPairId      int
	SecondPairId     int
	FirstPairName    string
	SecondPairName   string
	FirstPairSymbol  string
	SecondPairSymbol string
	FirstPairPrice   float64
	SecondPairPrice  float64
	QtyBuys          float64
	QtySells         float64
	TotalBuys        float64
	TotalSells       float64
	QtyAvailable     float64
	CurrentPrice     float64
	ActualReturn     float64
	FutureReturn     float64
	TotalReturn      float64
	ReturnBtc        float64
	ReturnUsd        float64
	Roi              float64
	BtcPrice         float64
	Subtrades        []Subtrade
}

var tradesWss = make(map[string][]WsTrade)

func InstanciateTradesDispatcher() {
	for {
		var users []string
		user_sql := `
			SELECT DISTINCT
					u.username
			FROM subtrades s
			LEFT JOIN trades t on(s.tradeid=t.id)
			LEFT JOIN users u ON(t.userid = u.id)
			WHERE (
					s.updatedat > current_timestamp - interval '1 seconds' OR
					t.updatedat > current_timestamp - interval '1 seconds'
			);`
		user_rows, err := DbWebApp.Query(user_sql)
		defer user_rows.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"custom_msg": "Failed fetching user_sql",
			}).Error(err)
			return
		}
		for user_rows.Next() {
			var user string
			if err = user_rows.Scan(&user); err != nil {
				log.WithFields(log.Fields{
					"custom_msg": "Failed scanning user_sql",
				}).Error(err)
				return
			}
			users = append(users, user)
		}

		for _, x := range users {
			for _, q := range tradesWss[x] {
				q.Channel <- q.Username
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func GetTrades(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetTrades..."))

	username := mux.Vars(r)["username"]
	requestid := mux.Vars(r)["requestid"]

	c := make(chan string)
	listener := WsTrade{c, username, requestid}
	tradesWss[username] = append(tradesWss[username], listener)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		for _, origin := range Origins {
			if origin == r.Header["Origin"][0] {
				return true
			}
		}
		return false
	}
	ws, _ := upgrader.Upgrade(w, r, nil)

	wsTradeOutput := NewSelectUserTrades(username)
	err := ws.WriteJSON(wsTradeOutput)
	if err != nil {
		ws.Close()
		log.WithFields(log.Fields{
			"sessionid":  requestid,
			"username":   username,
			"custom_msg": "Failed running sending initial trades ws",
		}).Error(err)
		return
	}

	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				for i, v := range tradesWss[username] {
					if v.RequestId == requestid {
						tradesWss[username] = append(tradesWss[username][:i], tradesWss[username][i+1:]...)
					}
				}
				ws.Close()
				return
			} else {
				ws.Close()
				log.WithFields(log.Fields{
					"sessionid":  requestid,
					"username":   username,
					"custom_msg": "Failed running receiving trades ws",
				}).Error(err)
				return
			}
		}
	}()
	go func() {
		for {
			s1 := <-c
			if s1 == username {
				wsTradeOutput := NewSelectUserTrades(username)
				err := ws.WriteJSON(wsTradeOutput)
				if err != nil {
					ws.Close()
					log.WithFields(log.Fields{
						"sessionid":  requestid,
						"username":   username,
						"custom_msg": "Failed running sending trades ws",
					}).Error(err)
					return
				}
			}
		}
	}()
}

func NewSelectUserTrades(username string) (wsTradeOutput WsTradeOutput) {

	user_details_sql := `
		SELECT
			username,
			CASE WHEN twitter IS NULL THEN '' ELSE twitter END AS twitter
		FROM users u
		WHERE username = $1;`

	userDetails := UserDetails{}
	err := DbWebApp.QueryRow(
		user_details_sql,
		username).Scan(&userDetails.Username, &userDetails.Twitter)
	if err != nil {
		log.WithFields(log.Fields{
			"username":   username,
			"custom_msg": "Failed running user_details",
		}).Error(err)
		return
	}

	trades := []Trade{}

	trades_sql := `
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
			TRADES_MACRO AS (
				SELECT
					t.id,
					t.isopen,
					t.exchange,
					t.firstpair,
					t.secondpair,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'BUY' THEN s.quantity END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'BUY' THEN s.quantity END)
					END AS qtybuys,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'SELL' THEN s.quantity END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'SELL' THEN s.quantity END)
					END AS qtysells,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'BUY' THEN s.total END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'BUY' THEN s.total END)
					END AS totalbuys,
					CASE
						WHEN SUM(CASE WHEN s."type" = 'SELL' THEN s.total END) IS NULL THEN 0
						ELSE SUM(CASE WHEN s."type" = 'SELL' THEN s.total END)
					END AS totalsells
				FROM trades t
				LEFT JOIN subtrades s ON(t.id  = s.tradeid)
				INNER JOIN users u ON(t.userid = u.id)
				WHERE u.username = $1
				GROUP BY 1, 2, 3, 4, 5),
			TRADES_MICRO AS (
				SELECT
					t.id,
					t.isopen,
					t.exchange,
					t.firstpair AS firstpairid,
					c1.name AS firstpairname,
					c1.symbol AS firstpairsymbol,
					c1.price AS firstpairprice,
					t.secondpair AS secondpairid,
					c2.name AS secondpairname,
					c2.symbol AS secondpairsymbol,
					c2.price AS secondpairprice,
					t.qtybuys,
					t.qtysells,
					t.totalbuys,
					t.totalsells,
					t.qtybuys - t.qtysells AS qtyavailable,
					(c2.price / c1.price) AS currentprice,
					t.totalsells - t.totalbuys AS actualreturn,
					(t.qtybuys - t.qtysells) * (c2.price / c1.price) AS futurereturn,
					t.totalsells - t.totalbuys + (t.qtybuys - t.qtysells) * (c2.price / c1.price) AS totalreturn,
					CASE
						WHEN t.totalbuys = 0 THEN 0
						ELSE (((t.qtybuys - t.qtysells) * (c2.price / c1.price) + t.totalsells) / t.totalbuys - 1) * 100
					END AS roi
				FROM TRADES_MACRO t
				LEFT JOIN CURRENT_PRICE c1 ON(t.firstpair = c1.coinid)
				LEFT JOIN CURRENT_PRICE c2 ON(t.secondpair = c2.coinid))
		SELECT
			t.id,
			t.isopen,
			t.exchange,
			t.firstpairid,
			t.firstpairname,
			t.firstpairsymbol,
			t.firstpairprice,
			t.secondpairid,
			t.secondpairname,
			t.secondpairsymbol,
			t.secondpairprice,
			t.qtybuys,
			t.qtysells,
			t.totalbuys,
			t.totalsells,
			t.qtyavailable,
			t.currentprice,
			t.actualreturn,
			t.futurereturn,
			t.totalreturn,
			t.totalreturn * t.firstpairprice / c3.price as returnbtc,
			t.totalreturn * t.firstpairprice as returnusd,
			t.roi,
			c3.price AS btcprice
		FROM TRADES_MICRO t
		LEFT JOIN CURRENT_PRICE c3 ON(c3.coinid = 1);`

	trades_rows, err := DbWebApp.Query(
		trades_sql,
		username)
	defer trades_rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"username":   username,
			"custom_msg": "Failed running trades_sql",
		}).Error(err)
	}
	for trades_rows.Next() {
		trade := Trade{}
		if err = trades_rows.Scan(
			&trade.Id,
			&trade.IsOpen,
			&trade.Exchange,
			&trade.FirstPairId,
			&trade.FirstPairName,
			&trade.FirstPairSymbol,
			&trade.FirstPairPrice,
			&trade.SecondPairId,
			&trade.SecondPairName,
			&trade.SecondPairSymbol,
			&trade.SecondPairPrice,
			&trade.QtyBuys,
			&trade.QtySells,
			&trade.TotalBuys,
			&trade.TotalSells,
			&trade.QtyAvailable,
			&trade.CurrentPrice,
			&trade.ActualReturn,
			&trade.FutureReturn,
			&trade.TotalReturn,
			&trade.ReturnBtc,
			&trade.ReturnUsd,
			&trade.Roi,
			&trade.BtcPrice,
		); err != nil {
			log.WithFields(log.Fields{
				"username":   username,
				"custom_msg": "Failed parsing trades_sql",
			}).Error(err)
		}

		subtrades_sql := `
			SELECT
				id,
				type,
				reason,
				TO_CHAR(tradetimestamp, 'YYYY-MM-DD"T"HH24:MI'),
				quantity,
				avgprice,
				total
			FROM subtrades
			WHERE tradeid = $1
			ORDER BY 1;`

		subtrades := []Subtrade{}
		subtrades_rows, err := DbWebApp.Query(
			subtrades_sql,
			trade.Id)
		defer subtrades_rows.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"username":   username,
				"custom_msg": "Failed running subtrades_sql",
			}).Error(err)

		}

		for subtrades_rows.Next() {
			subtrade := Subtrade{}
			if err = subtrades_rows.Scan(
				&subtrade.Id,
				&subtrade.Type,
				&subtrade.Reason,
				&subtrade.Timestamp,
				&subtrade.Quantity,
				&subtrade.AvgPrice,
				&subtrade.Total); err != nil {
				log.WithFields(log.Fields{
					"username":   username,
					"custom_msg": "Failed parsing subtrades_sql",
				}).Error(err)

			}

			subtrades = append(subtrades, subtrade)
		}

		trade.Subtrades = subtrades
		trades = append(trades, trade)
	}

	wsTradeOutput.Trades = trades
	wsTradeOutput.UserDetails = userDetails

	return
}
