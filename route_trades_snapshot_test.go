package main

import (
	"fmt"
	"math"
	"testing"
)

func TestGetSnapshot(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			wallet, username, privacy,
			plan, createdat, updatedat)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'jsjsjsj',
			'all', 'basic', current_timestamp, current_timestamp);`)

	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'BTC', 'BTC', 'bitcoin'),
			(1001, 'USDC', 'USDC', 'usdollar')`)

	Db.Exec(`
		INSERT INTO prices (
			createdat, coinid, price)
		VALUES
			(current_timestamp, 1, 65000),
			(current_timestamp, 1001, 1);`)

	// <test code>
	t.Run(fmt.Sprintf("Test snapshot single buy"), func(t *testing.T) {
		Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES (
			'MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
			current_timestamp, 1001, 1, TRUE);`)
		Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES (
			'SISISIS', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
			current_timestamp, 'BUY', 1, 50000, 50000, 'TESTART');`)
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
		snapshot := TradesSnapshot{}
		snapshot.Trades = user.SelectUserTrades()
		snapshot.CountTrades = len(snapshot.Trades)
		snapshot.CalculateTradesTotals()
		if snapshot.Trades[0].QtyBuys != 1 {
			t.Fatal("Failed test snapshot single buy, trade[0].QtyBuys")
		}
		if snapshot.Trades[0].TotalBuysBtc != 1.0*50000/65000 {
			t.Fatal("Failed test snapshot single buy, trade[0].TotalBuysBtc")
		}
		if snapshot.Trades[0].QtySells != 0 {
			t.Fatal("Failed test snapshot single buy, trade[0].QtySells")
		}
		if snapshot.Trades[0].TotalSellsBtc != 0 {
			t.Fatal("Failed test snapshot single buy, trade[0].TotalSellsBtc")
		}
		if snapshot.Trades[0].QtyAvailable != 1 {
			t.Fatal("Failed test snapshot single buy, trade[0].QtyAvailable")
		}
		if snapshot.Trades[0].Roi != 100*(65000.0/50000-1) {
			t.Fatal("Failed test snapshot single buy, trade[0].Roi")
		}
		if snapshot.CountTrades != 1 {
			t.Fatal("Failed test snapshot single buy, QtyBuys")
		}
		if snapshot.TotalReturnUsd != 15000 {
			t.Fatal("Failed test snapshot single buy, TotalReturnUsd")
		}
		if snapshot.TotalReturnBtc != math.Round(1.0*15000/65000*100)/100 {
			t.Fatal("Failed test snapshot single buy, TotalReturnBtc")
		}
		if math.Round(snapshot.Roi) != 100*(65000.0/50000-1) {
			t.Fatal("Failed test snapshot single buy, ROI")
		}
		Db.Exec(`DELETE FROM trades WHERE 1 = 1;`)
	})

	t.Run(fmt.Sprintf("Test snapshot buy and sell"), func(t *testing.T) {
		Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES (
			'MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
			current_timestamp, 1001, 1, TRUE);`)
		Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES
			('SISISIS', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
			current_timestamp, 'BUY', 1, 50000, 50000, 'TESTART'),
			('SISISIS2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
			current_timestamp, 'SELL', 0.5, 80000, 40000, 'TESTART');`)
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
		snapshot := TradesSnapshot{}
		snapshot.Trades = user.SelectUserTrades()
		snapshot.CountTrades = len(snapshot.Trades)
		snapshot.CalculateTradesTotals()
		if snapshot.Trades[0].QtyBuys != 1 {
			t.Fatal("Failed test snapshot buy and sell, trade[0].QtyBuys")
		}
		if snapshot.Trades[0].TotalBuysBtc != 1.0*50000/65000 {
			t.Fatal("Failed test snapshot buy and sell, trade[0].TotalBuysBtc")
		}
		if snapshot.Trades[0].QtySells != 0.5 {
			t.Fatal("Failed test snapshot buy and sell, trade[0].QtySells")
		}
		if snapshot.Trades[0].TotalSellsBtc != 1.0*0.5*80000/65000 {
			t.Fatal("Failed test snapshot buy and sell, trade[0].TotalSellsBtc")
		}
		if snapshot.Trades[0].QtyAvailable != 0.5 {
			t.Fatal("Failed test snapshot buy and sell, trade[0].QtyAvailable")
		}
		if snapshot.Trades[0].Roi != 100*((65000.0*0.5+80000.0*0.5)/50000-1) {
			t.Fatal("Failed test snapshot buy and sell, trade[0].Roi")
		}
		if snapshot.CountTrades != 1 {
			t.Fatal("Failed test snapshot buy and sell, QtyBuys")
		}
		if snapshot.TotalReturnUsd != 22500 {
			t.Fatal("Failed test snapshot buy and sell, TotalReturnUsd")
		}
		if snapshot.TotalReturnBtc != math.Round((80000*0.5-1*50000+65000*0.5)/65000*100)/100 {
			t.Fatal("Failed test snapshot buy and sell, TotalReturnBtc")
		}
		if math.Round(snapshot.Roi) != 100*((80000*0.5+65000*0.5)/50000-1) {
			t.Fatal("Failed test snapshot buy and sell, ROI")
		}
		Db.Exec(`DELETE FROM trades WHERE 1 = 1;`)
	})
	t.Run(fmt.Sprintf("Test snapshot multiple buy and sell"), func(t *testing.T) {
		Db.Exec(`
			INSERT INTO trades(
				code, userwallet, createdat, updatedat,
				firstpair, secondpair, isopen)
			VALUES (
				'MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
				current_timestamp, 1001, 1, TRUE);`)
		Db.Exec(`
			INSERT INTO subtrades(
				code, userwallet, tradecode, createdat, updatedat,
				type, quantity, avgprice, total, reason)
			VALUES
				('SISISIS', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'BUY', 1, 50000, 50000, 'TESTART'),
				('SISISIS2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'BUY', 2, 70000, 140000, 'TESTART'),
				('SISISIS3', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'SELL', 1.5, 100000, 150000, 'TESTART'),
				('SISISIS4', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'SELL', 0.5, 80000, 40000, 'TESTART');`)
		snapshot := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}.GetSnapshot()
		if snapshot.Trades[0].QtyBuys != 3 {
			t.Fatal("Failed test snapshot multiple buy and sell, trade[0].QtyBuys")
		}
		if snapshot.Trades[0].TotalBuysBtc != (1.0*50000+2*70000)/65000 {
			t.Fatal("Failed test snapshot multiple buy and sell, trade[0].TotalBuysBtc")
		}
		if snapshot.Trades[0].QtySells != 2 {
			t.Fatal("Failed test snapshot multiple buy and sell, trade[0].QtySells")
		}
		if snapshot.Trades[0].TotalSellsBtc != (1.5*100000+0.5*80000)/65000 {
			t.Fatal("Failed test snapshot multiple buy and sell, trade[0].TotalSellsBtc")
		}
		if snapshot.Trades[0].QtyAvailable != 1 {
			t.Fatal("Failed test snapshot multiple buy and sell, trade[0].QtyAvailable")
		}
		if snapshot.Trades[0].Roi != math.Round((((1.0*65000+1.5*100000+0.5*80000)/(1.0*50000+2*70000)-1)*100)*10)/10 {
			t.Fatal("Failed test snapshot multiple buy and sell, trade[0].Roi")
		}
		if snapshot.CountTrades != 1 {
			t.Fatal("Failed test snapshot multiple buy and sell, QtyBuys")
		}
		if snapshot.TotalReturnUsd != (1.0*65000+1.5*100000+0.5*80000)-(1.0*50000+2*70000) {
			t.Fatal("Failed test snapshot multiple buy and sell, TotalReturnUsd")
		}
		if snapshot.TotalReturnBtc != ((1.0*65000+1.5*100000+0.5*80000)-(1.0*50000+2*70000))/65000 {
			t.Fatal("Failed test snapshot multiple buy and sell, TotalReturnBtc")
		}
		if math.Round(snapshot.Roi) != math.Round(100*((1.0*65000+1.5*100000+0.5*80000)/(1.0*50000+2*70000)-1)) {
			t.Fatal("Failed test snapshot multiple buy and sell, ROI")
		}
		Db.Exec(`DELETE FROM trades WHERE 1 = 1;`)
	})

	t.Run(fmt.Sprintf("Test snapshot multiple trades"), func(t *testing.T) {
		Db.Exec(`
			INSERT INTO trades(
				code, userwallet, createdat, updatedat,
				firstpair, secondpair, isopen)
			VALUES
				('MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
				current_timestamp, 1001, 1, TRUE),
				('MBMBMBM2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
				current_timestamp, 1001, 1, TRUE);`)
		Db.Exec(`
			INSERT INTO subtrades(
				code, userwallet, tradecode, createdat, updatedat,
				type, quantity, avgprice, total, reason)
			VALUES
				('SISISIS', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'BUY', 1, 50000, 50000, 'TESTART'),
				('SISISIS2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'SELL', 1, 80000, 80000, 'TESTART'),
				('SISISIS3', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM2', current_timestamp,
				current_timestamp, 'BUY', 1, 50000, 50000, 'TESTART'),
				('SISISIS4', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM2', current_timestamp,
				current_timestamp, 'SELL', 1, 80000, 80000, 'TESTART');`)
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
		snapshot := TradesSnapshot{}
		snapshot.Trades = user.SelectUserTrades()
		snapshot.CountTrades = len(snapshot.Trades)
		snapshot.CalculateTradesTotals()
		if snapshot.CountTrades != 2 {
			t.Fatal("Failed test snapshot multiple trades, QtyBuys")
		}
		if snapshot.TotalReturnUsd != 60000 {
			t.Fatal("Failed test snapshot multiple trades, TotalReturnUsd")
		}
		if snapshot.TotalReturnBtc != math.Round((1.0*60000/65000)*100)/100 {
			t.Fatal("Failed test snapshot multiple trades, TotalReturnBtc")
		}
		if math.Round(snapshot.Roi) != math.Round(100*(160000.0/100000-1)) {
			t.Fatal("Failed test snapshot multiple trades, ROI")
		}
		Db.Exec(`DELETE FROM trades WHERE 1 = 1;`)
	})

	t.Run(fmt.Sprintf("Test snapshot multiple trades multiple pairs"), func(t *testing.T) {
		Db.Exec(`
			INSERT INTO coins (
				coinid, name, symbol, slug)
			VALUES
				(2, 'DOT', 'DOT', 'POLKADOT'),
				(3, 'SOL', 'SOL', 'SOLANA'),
				(4, 'LUNA', 'LUNA', 'LUNA'),
				(5, 'XLM', 'XLM', 'STELLAR');`)
		Db.Exec(`
			INSERT INTO prices (
				createdat, coinid, price)
			VALUES
				(current_timestamp, 1, 100000),
				(current_timestamp, 2, 100),
				(current_timestamp, 3, 1000),
				(current_timestamp, 4, 200),
				(current_timestamp, 5, 400);`)
		Db.Exec(`
			INSERT INTO trades(
				code, userwallet, createdat, updatedat,
				firstpair, secondpair, isopen)
			VALUES
				('MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
				current_timestamp, 2, 3, TRUE),
				('MBMBMBM2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
				current_timestamp, 4, 5, TRUE);`)
		Db.Exec(`
			INSERT INTO subtrades(
				code, userwallet, tradecode, createdat, updatedat,
				type, quantity, avgprice, total, reason)
			VALUES
				('SISISIS', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'BUY', 1, 5, 5, 'TESTART'),
				('SISISIS2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
				current_timestamp, 'SELL', 1, 10, 10, 'TESTART'),
				('SISISIS3', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM2', current_timestamp,
				current_timestamp, 'BUY', 1, 1, 1, 'TESTART'),
				('SISISIS4', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM2', current_timestamp,
				current_timestamp, 'SELL', 1, 2, 2, 'TESTART');`)
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
		snapshot := TradesSnapshot{}
		snapshot.Trades = user.SelectUserTrades()
		snapshot.CountTrades = len(snapshot.Trades)
		snapshot.CalculateTradesTotals()
		if snapshot.CountTrades != 2 {
			t.Fatal("Failed test snapshot multiple trades multiple pairs, QtyBuys")
		}
		if snapshot.TotalReturnUsd != 700 {
			t.Fatal("Failed test snapshot multiple trades multiple pairs, TotalReturnUsd")
		}
		if snapshot.TotalReturnBtc != math.Round((1.0*700/100000)*100)/100 {
			t.Fatal("Failed test snapshot multiple trades multiple pairs, TotalReturnBtc")
		}
		if math.Round(snapshot.Roi) != math.Round(100*((10*100+2*200)/(5*100+1*200)-1)) {
			t.Fatal("Failed test snapshot multiple trades multiple pairs, ROI")
		}
		Db.Exec(`DELETE FROM trades WHERE 1 = 1;`)
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
