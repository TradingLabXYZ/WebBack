package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestListTrades(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'userd', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (wallet, totalcounttrades, totalportfolio,
			totalreturn, totalroi, tradeqtyavailable, tradevalue, tradereturn,
			traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)
	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(1000, 'USDC', 'USDC', 'usdc');`)
	Db.Exec(`
		INSERT INTO lastprices (
			updatedat, coinid, price)
		VALUES
			(current_timestamp, 1, 45000),
			(current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('useratr1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp, current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES ('userasub1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr1', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/list_trades", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		ListTrades(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed list trades, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test origin not api"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web")
		req := httptest.NewRequest("GET", "/list_trades", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		ListTrades(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed list trades, origin not api")
		}
	})

	t.Run(fmt.Sprintf("Test successfully receive trade codes"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("api")
		req := httptest.NewRequest("GET", "/list_trades", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		ListTrades(w, req)
		trades_codes, _ := ioutil.ReadAll(w.Body)
		trades_codes_s := string(trades_codes)
		if !strings.Contains(trades_codes_s, "useratr1") {
			t.Fatal("Failed list trades, successfully receive trade codes")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestListSubtrades(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'userd', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (wallet, totalcounttrades, totalportfolio,
			totalreturn, totalroi, tradeqtyavailable, tradevalue, tradereturn,
			traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)
	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(1000, 'USDC', 'USDC', 'usdc');`)
	Db.Exec(`
		INSERT INTO lastprices (
			updatedat, coinid, price)
		VALUES
			(current_timestamp, 1, 45000),
			(current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('useratr1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp, current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES ('userasub1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr1', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/list_subtrades", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		ListSubtrades(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed list subtrades, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test origin not api"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web")
		req := httptest.NewRequest("GET", "/list_subtrades", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		vars := map[string]string{
			"tradecode": "tempTradeCode",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		ListSubtrades(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed list subtrades, origin not api")
		}
	})

	t.Run(fmt.Sprintf("Test successfully receive subtrade codes"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("api")
		req := httptest.NewRequest("GET", "/list_subtrades", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		vars := map[string]string{
			"tradecode": "useratr1",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		ListSubtrades(w, req)
		subtrades_codes, _ := ioutil.ReadAll(w.Body)
		subtrades_codes_s := string(subtrades_codes)
		if !strings.Contains(subtrades_codes_s, "userasub1") {
			t.Fatal("Failed list subtrades, successfully receive subtrade codes")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestApiGetSnapshot(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'userd', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (wallet, totalcounttrades, totalportfolio,
			totalreturn, totalroi, tradeqtyavailable, tradevalue, tradereturn,
			traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)
	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(1000, 'USDC', 'USDC', 'usdc');`)
	Db.Exec(`
		INSERT INTO lastprices (
			updatedat, coinid, price)
		VALUES
			(current_timestamp, 1, 45000),
			(current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('useratr1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp, current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES ('userasub1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr1', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_snapshot", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		GetSnapshot(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed get snapshot, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test origin not api"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web")
		req := httptest.NewRequest("GET", "/get_snapshot", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		GetSnapshot(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed get snapshot, origin not api")
		}
	})

	t.Run(fmt.Sprintf("Test successfully get snapshot"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("api")
		req := httptest.NewRequest("GET", "/get_snapshot", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		GetSnapshot(w, req)
		message, _ := ioutil.ReadAll(w.Body)
		temp_snapshot := TradesSnapshot{}
		json.Unmarshal(message, &temp_snapshot)
		if temp_snapshot.TotalPortfolioUsd != "45,000" {
			t.Fatal("Failed get snapshot, successfully get snapshot")
		}
		if temp_snapshot.Trades[0].Code != "useratr1" {
			t.Fatal("Failed get snapshot, successfully get snapshot")
		}
		if temp_snapshot.Trades[0].Subtrades[0].Code != "userasub1" {
			t.Fatal("Failed get snapshot, successfully get snapshot")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
