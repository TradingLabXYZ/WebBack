package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

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
	Username         string
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
	CurrentPrice     float64
	QtyBuys          float64
	QtySells         float64
	QtyAvailable     float64
	TotalBuys        float64
	TotalBuysBtc     float64
	TotalBuysUsd     float64
	TotalSells       float64
	TotalSellsBtc    float64
	TotalSellsUsd    float64
	ActualReturn     float64
	FutureReturn     float64
	FutureReturnBtc  float64
	FutureReturnUsd  float64
	TotalReturn      float64
	TotalReturnBtc   float64
	TotalReturnUsd   float64
	Roi              float64
	BtcPrice         float64
	Subtrades        []Subtrade
}

type TradesOutput struct {
	UserDetails    UserDetails
	Trades         []Trade
	CountTrades    float64
	TotalReturnUsd float64
	TotalReturnBtc float64
	Roi            float64
}

type WsTrade struct {
	Channel   chan TradesOutput
	RequestId string
}

var tradesWss = make(map[string][]WsTrade)

func InstanciateTradesDispatcher() {
	for {
		user_sql := `
			SELECT DISTINCT
					u.username
			FROM subtrades s
			LEFT JOIN trades t on(s.tradeid=t.id)
			LEFT JOIN users u ON(t.userid = u.id)
			WHERE (
					s.updatedat > current_timestamp - interval '1 seconds' OR
					t.updatedat > current_timestamp - interval '1 seconds' OR
					u.updatedat > current_timestamp - interval '1 seconds'
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
			var username string
			if err = user_rows.Scan(&username); err != nil {
				log.WithFields(log.Fields{
					"custom_msg": "Failed scanning user_sql",
				}).Error(err)
				return
			}
			go func() {
				user := SelectUser("username", username)
				userSnapshot := user.GetUserSnapshot()
				for _, q := range tradesWss[username] {
					q.Channel <- userSnapshot
				}
			}()
		}
		time.Sleep(1 * time.Second)
	}
}

func GetTrades(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetTrades..."))

	username := mux.Vars(r)["username"]
	userToSee := SelectUser("username", username)

	status := CheckPrivacy(r, userToSee)
	if status != "OK" {
		w.Write([]byte(status))
		return
	}

	c := make(chan TradesOutput)
	requestid := mux.Vars(r)["requestid"]
	listener := WsTrade{c, requestid}
	tradesWss[username] = append(tradesWss[username], listener)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		for _, origin := range Origins {
			if origin == r.Header["Origin"][0] {
				return true
			}
		}
		return false
	}
	ws, _ := upgrader.Upgrade(w, r, nil)

	wsTradeOutput := userToSee.GetUserSnapshot()
	err := ws.WriteJSON(wsTradeOutput)
	if err != nil {
		ws.Close()
		log.WithFields(log.Fields{
			"sessionid":  requestid,
			"username":   username,
			"custom_msg": "Failed running sending initial snapshot",
		}).Error(err)
		return
	}

	// RECEIVE MESSAGES
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

	// SEND MESSAGES
	go func() {
		for {
			s1 := <-c
			if s1.UserDetails.Username == username {
				err := ws.WriteJSON(s1)
				if err != nil {
					ws.Close()
					log.WithFields(log.Fields{
						"sessionid":  requestid,
						"username":   username,
						"custom_msg": "Failed running sending snapshot",
					}).Error(err)
					return
				}
			}
		}
	}()
}

func CheckPrivacy(request *http.Request, userToSee User) (status string) {
	fmt.Println(Gray(8-1, "Starting CheckUserPrivacy..."))

	if userToSee.Privacy == "all" {
		return "OK"
	}

	session, err := GetSession(request, "cookie")
	if err != nil {
		return "KO"
	}

	user := SelectUser("email", session.Email)
	if user.Id == userToSee.Id {
		return "OK"
	}

	switch userToSee.Privacy {
	case "private":
		return `{"Status": "denied", "Reason": "private"}`
	case "followers":
		var isfollower bool
		_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE followto = $1
					AND followfrom = $2;`, user.Id, userToSee.Id).Scan(
			&isfollower,
		)
		if isfollower {
			return "OK"
		} else {
			return `{"Status": "denied", "Reason": "follow"}`
		}
	case "subscribers":
		var issubscriber bool
		_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE subscribeto = $1
					AND subscribefrom = $2;`, user.Id, userToSee.Id).Scan(
			&issubscriber,
		)
		if issubscriber {
			return "OK"
		} else {
			return `{"Status": "denied", "Reason": "subscribe"}`
		}
	default:
		return `{"Status": "denied", "Reason": "unknown"}`
	}
}

func (user User) GetUserSnapshot() (tradesOutput TradesOutput) {

	tradesOutput.UserDetails = UserDetails{
		user.UserName,
		user.Twitter,
	}

	tradesOutput.Trades = user.SelectUserTrades()
	tradesOutput.CalculateTradesTotals()
	return
}

func (user User) SelectUserTrades() (trades []Trade) {

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
					u.username,
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
				GROUP BY 1, 2, 3, 4, 5, 6),
			TRADES_MICRO AS (
				SELECT
					t.id,
					t.username,
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
					(c2.price / c1.price) AS currentprice,
					t.qtybuys,
					t.qtysells,
					t.qtybuys - t.qtysells AS qtyavailable,
					t.totalbuys,
					t.totalsells,
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
			t.username,
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
			t.currentprice,
			t.qtybuys,
			t.qtysells,
			t.qtyavailable,
			t.totalbuys,
			t.totalbuys * t.firstpairprice / c3.price AS totalbuysbtc,
			t.totalbuys * t.firstpairprice AS totalbuysusd,
			t.totalsells,
			t.totalsells * t.firstpairprice / c3.price AS totalsellbtc,
			t.totalsells * t.firstpairprice AS totalsellusd,
			t.actualreturn,
			t.futurereturn,
			t.futurereturn * t.firstpairprice / c3.price AS futurereturnbtc,
			t.futurereturn * t.firstpairprice AS futurereturnusd,
			t.totalreturn,
			t.totalreturn * t.firstpairprice / c3.price AS returnbtc,
			t.totalreturn * t.firstpairprice AS returnusd,
			t.roi,
			c3.price AS btcprice
		FROM TRADES_MICRO t
		LEFT JOIN CURRENT_PRICE c3 ON(c3.coinid = 1);`

	trades_rows, err := DbWebApp.Query(
		trades_sql,
		user.UserName)
	defer trades_rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"username":   user.UserName,
			"custom_msg": "Failed running trades_sql",
		}).Error(err)
	}
	for trades_rows.Next() {
		trade := Trade{}
		if err = trades_rows.Scan(
			&trade.Id,
			&trade.Username,
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
			&trade.CurrentPrice,
			&trade.QtyBuys,
			&trade.QtySells,
			&trade.QtyAvailable,
			&trade.TotalBuys,
			&trade.TotalBuysBtc,
			&trade.TotalBuysUsd,
			&trade.TotalSells,
			&trade.TotalSellsBtc,
			&trade.TotalSellsUsd,
			&trade.ActualReturn,
			&trade.FutureReturn,
			&trade.FutureReturnBtc,
			&trade.FutureReturnUsd,
			&trade.TotalReturn,
			&trade.TotalReturnBtc,
			&trade.TotalReturnUsd,
			&trade.Roi,
			&trade.BtcPrice,
		); err != nil {
			log.WithFields(log.Fields{
				"username":   user.UserName,
				"custom_msg": "Failed parsing trades_sql",
			}).Error(err)
		}

		subtrades := trade.SelectTradeSubtrades()
		trade.Subtrades = subtrades

		trades = append(trades, trade)
	}
	return
}

func (trade Trade) SelectTradeSubtrades() (subtrades []Subtrade) {
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

	subtrades_rows, err := DbWebApp.Query(
		subtrades_sql,
		trade.Id)
	defer subtrades_rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"username":   trade.Username,
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
				"username":   trade.Username,
				"custom_msg": "Failed parsing subtrades_sql",
			}).Error(err)
		}

		subtrades = append(subtrades, subtrade)
	}
	return
}

func (tradesOutput *TradesOutput) CalculateTradesTotals() {
	var totalReturnBtc float64
	var totalReturnUsd float64
	var totalBuysBtc float64
	var totalSellBtc float64
	var futureReturnBtc float64
	for _, trade := range tradesOutput.Trades {
		totalReturnBtc = totalReturnBtc + trade.TotalReturnBtc
		totalReturnUsd = totalReturnUsd + trade.TotalReturnUsd
		totalBuysBtc = totalBuysBtc + trade.TotalBuysBtc
		totalSellBtc = totalSellBtc + trade.TotalSellsBtc
		futureReturnBtc = futureReturnBtc + trade.FutureReturnBtc
	}
	tradesOutput.TotalReturnBtc = totalReturnBtc
	tradesOutput.TotalReturnUsd = totalReturnUsd
	tradesOutput.Roi = ((futureReturnBtc+totalSellBtc)/totalBuysBtc - 1) * 100
}
