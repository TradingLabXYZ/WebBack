package main

import (
	"math"

	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

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
				SELECT DISTINCT ON(p.coinid)
					p.coinid,
					p.createdat,
					c.name,
					c.symbol,
					p.price
				FROM prices p
				LEFT JOIN coins c ON(p.coinid = c.coinid)
				ORDER BY 1, 2 DESC)
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
			CASE WHEN t.currentprice > 1 THEN TO_CHAR(ROUND(t.currentprice, 2), '999,999,999.00') ELSE TO_CHAR(ROUND(t.currentprice, 5), '999,999,999.00000') END as currentprice,
			t.qtybuys,
			t.qtysells,
			CASE
				WHEN RIGHT(TO_CHAR(ROUND(t.qtyavailable, 2), '999,999,999.00'), 2) = '00' THEN TO_CHAR(ROUND(t.qtyavailable, 2), '999,999,999')
				ELSE TO_CHAR(ROUND(t.qtyavailable, 2), '999,999,999.00')
			END as qtyavailable,
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
			c3.price AS btcprice,
			ROUND(t.qtyavailable * t.secondpairprice, 2) AS totalvalueusd
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
			&trade.TotalValueUsd,
		); err != nil {
			log.WithFields(log.Fields{
				"wallet":     user.Wallet,
				"custom_msg": "Failed parsing trades_sql",
			}).Error(err)
		}

		trade.TotalValueUsdS = humanize.Commaf(trade.TotalValueUsd)
		trade.TotalReturnS = humanize.Commaf(trade.TotalReturn)
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
		totalReturnBtc    float64
		totalReturnUsd    float64
		totalPortfolioUsd float64
		totalBuysBtc      float64
		totalSellBtc      float64
		futureReturnBtc   float64
	)
	for _, trade := range snapshot.Trades {
		totalReturnBtc = totalReturnBtc + trade.TotalReturnBtc
		totalReturnUsd = totalReturnUsd + trade.TotalReturnUsd
		totalBuysBtc = totalBuysBtc + trade.TotalBuysBtc
		totalSellBtc = totalSellBtc + trade.TotalSellsBtc
		futureReturnBtc = futureReturnBtc + trade.FutureReturnBtc
		totalPortfolioUsd = totalPortfolioUsd + trade.TotalValueUsd
	}

	snapshot.TotalReturnBtc = humanize.Commaf(math.Round(totalReturnBtc*100) / 100)
	snapshot.TotalReturnUsd = humanize.Commaf(math.Round(totalReturnUsd*100) / 100)
	snapshot.TotalPortfolioUsd = humanize.Commaf(math.Round(totalPortfolioUsd*100) / 100)
	snapshot.Roi = math.Round(((futureReturnBtc+totalSellBtc)/totalBuysBtc-1)*100*100) / 100

	if math.IsNaN(snapshot.Roi) || math.IsInf(snapshot.Roi, 0) {
		snapshot.Roi = 0
	}
}
