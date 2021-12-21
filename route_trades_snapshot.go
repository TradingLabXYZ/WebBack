package main

import (
	"math"

	log "github.com/sirupsen/logrus"
)

type UserDetails struct {
	Username       string
	Twitter        string
	Github         string
	Discord        string
	Followers      int
	Followings     int
	Subscribers    int
	ProfilePicture string
	JoinTime       string
}

type Subtrade struct {
	Code      string
	TradeCode string
	CreatedAt string
	Type      string
	Reason    string
	Quantity  float64
	AvgPrice  float64
	Total     float64
}

type Trade struct {
	Code              string
	Username          string
	Userwallet        string
	Exchange          string
	FirstPairId       int
	SecondPairId      int
	FirstPairName     string
	SecondPairName    string
	FirstPairSymbol   string
	SecondPairSymbol  string
	FirstPairPrice    float64
	SecondPairPrice   float64
	FirstPairUrlIcon  string
	SecondPairUrlIcon string
	CurrentPrice      float64
	QtyBuys           float64
	QtySells          float64
	QtyAvailable      float64
	TotalBuys         float64
	TotalBuysBtc      float64
	TotalBuysUsd      float64
	TotalSells        float64
	TotalSellsBtc     float64
	TotalSellsUsd     float64
	ActualReturn      float64
	FutureReturn      float64
	FutureReturnBtc   float64
	FutureReturnUsd   float64
	TotalReturn       float64
	TotalReturnBtc    float64
	TotalReturnUsd    float64
	Roi               float64
	BtcPrice          float64
	Subtrades         []Subtrade
}

type PrivacyStatus struct {
	Status string
	Reason string
}

type TradesSnapshot struct {
	UserDetails    UserDetails
	PrivacyStatus  PrivacyStatus
	IsFollower     bool
	IsSubscriber   bool
	Trades         []Trade
	CountTrades    int
	TotalReturnUsd float64
	TotalReturnBtc float64
	Roi            float64
}

func (user User) GetSnapshot() (snapshot TradesSnapshot) {
	snapshot.UserDetails = UserDetails{
		user.Username,
		user.Twitter,
		user.Github,
		user.Discord,
		user.Followers,
		user.Followings,
		user.Subscribers,
		user.ProfilePicture,
		user.JoinTime,
	}

	snapshot.Trades = user.SelectUserTrades()
	snapshot.CountTrades = len(snapshot.Trades)
	snapshot.CalculateTradesTotals()
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
					t.code,
					CASE WHEN u.username IS NULL THEN '' ELSE u.username END AS username,
					u.wallet AS userwallet,
					CASE WHEN t.exchange IS NULL THEN '' ELSE t.exchange END AS exchange,
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
				LEFT JOIN subtrades s ON(t.code  = s.tradecode)
				INNER JOIN users u ON(t.userwallet = u.wallet)
				WHERE u.wallet = $1
				GROUP BY 1, 2, 3, 4, 5, 6),
			TRADES_MICRO AS (
				SELECT
					t.code,
					t.username,
					t.userwallet,
					t.exchange,
					t.firstpair AS firstpairid,
					c1.name AS firstpairname,
					c1.symbol AS firstpairsymbol,
					c1.price AS firstpairprice,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.firstpair::TEXT || '.png' AS firstpairurlicon,
					t.secondpair AS secondpairid,
					c2.name AS secondpairname,
					c2.symbol AS secondpairsymbol,
					c2.price AS secondpairprice,
					'https://s2.coinmarketcap.com/static/img/coins/32x32/' || t.secondpair::TEXT || '.png' AS secondpairurlicon,
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
			t.code,
			t.username,
			t.userwallet,
			t.exchange,
			t.firstpairid,
			t.firstpairname,
			t.firstpairsymbol,
			t.firstpairprice,
			t.firstpairurlicon,
			t.secondpairid,
			t.secondpairname,
			t.secondpairsymbol,
			t.secondpairprice,
			t.secondpairurlicon,
			CASE WHEN t.currentprice > 1 THEN ROUND(t.currentprice, 2) ELSE ROUND(t.currentprice, 5) END as currentprice,
			t.qtybuys,
			t.qtysells,
			ROUND(t.qtyavailable, 2) as qtyavailable,
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
			ROUND(t.totalreturn, 2) as totalreturn,
			t.totalreturn * t.firstpairprice / c3.price AS returnbtc,
			t.totalreturn * t.firstpairprice AS returnusd,
			ROUND(t.roi, 1) AS roi,
			c3.price AS btcprice
		FROM TRADES_MICRO t
		LEFT JOIN CURRENT_PRICE c3 ON(c3.coinid = 1);`

	trades_rows, err := Db.Query(
		trades_sql,
		user.Wallet)
	defer trades_rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"wallet":     user.Wallet,
			"custom_msg": "Failed running trades_sql",
		}).Error(err)
	}
	for trades_rows.Next() {
		trade := Trade{}
		if err = trades_rows.Scan(
			&trade.Code,
			&trade.Username,
			&trade.Userwallet,
			&trade.Exchange,
			&trade.FirstPairId,
			&trade.FirstPairName,
			&trade.FirstPairSymbol,
			&trade.FirstPairPrice,
			&trade.FirstPairUrlIcon,
			&trade.SecondPairId,
			&trade.SecondPairName,
			&trade.SecondPairSymbol,
			&trade.SecondPairPrice,
			&trade.SecondPairUrlIcon,
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
				"wallet":     user.Wallet,
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
				code,
				tradecode,
				CASE WHEN type IS NULL THEN '' ELSE type END AS type,
				CASE WHEN reason IS NULL THEN '' ELSE reason END AS reason,
				TO_CHAR(createdat, 'YYYY-MM-DD"T"HH24:MI'),
				quantity,
				ROUND(avgprice, 6) AS avgprice,
				total
			FROM subtrades
			WHERE tradecode = $1
			ORDER BY 5;`

	subtrades_rows, err := Db.Query(
		subtrades_sql,
		trade.Code)
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
			&subtrade.Code,
			&subtrade.TradeCode,
			&subtrade.Type,
			&subtrade.Reason,
			&subtrade.CreatedAt,
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

func (snapshot *TradesSnapshot) CalculateTradesTotals() {
	var (
		totalReturnBtc  float64
		totalReturnUsd  float64
		totalBuysBtc    float64
		totalSellBtc    float64
		futureReturnBtc float64
	)
	for _, trade := range snapshot.Trades {
		totalReturnBtc = totalReturnBtc + trade.TotalReturnBtc
		totalReturnUsd = totalReturnUsd + trade.TotalReturnUsd
		totalBuysBtc = totalBuysBtc + trade.TotalBuysBtc
		totalSellBtc = totalSellBtc + trade.TotalSellsBtc
		futureReturnBtc = futureReturnBtc + trade.FutureReturnBtc
	}
	snapshot.TotalReturnBtc = math.Round(totalReturnBtc*100) / 100
	snapshot.TotalReturnUsd = math.Round(totalReturnUsd*100) / 100
	snapshot.Roi = math.Round(((futureReturnBtc+totalSellBtc)/totalBuysBtc-1)*100*100) / 100
	if math.IsNaN(snapshot.Roi) || math.IsInf(snapshot.Roi, 0) {
		snapshot.Roi = 0
	}
}
