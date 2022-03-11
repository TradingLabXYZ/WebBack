package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateTrade(t *testing.T) {
	// <setup code>
	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'bitcoin'),
			(1027, 'Ethereum', 'ETH', 'ethereum')`)

	Db.Exec(
		`INSERT INTO users (
			wallet, username, privacy,
			createdat, updatedat)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'r',
			'all', current_timestamp, current_timestamp);`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
	session, _ := user.InsertSession("web")

	// <test code>
	t.Run(fmt.Sprintf("Test insert trade with subtrade"), func(t *testing.T) {
		params := []byte(`{
			"Exchange": "Binance",
			"FirstPair": 1,
			"SecondPair": 1027,
			"Subtrades": [
				{
					"CreatedAt": "2021-10-01T19:39",
					"Type": "BUY",
					"Reason": "Volume",
					"Quantity": 1,
					"AvgPrice": 30000,
					"Total": 30000
				}]}`)
		req := httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		CreateTrade(w, req)
		res := w.Result()
		if res.StatusCode != 200 {
			t.Fatal("Failed test insert trade with subtrade")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade not valid cookie"), func(t *testing.T) {
		params := []byte(`{
			"Exchange": "Binance",
			"FirstPair": 1,
			"SecondPair": 1027,
		}`)
		req := httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId=testwrongcookie")
		w := httptest.NewRecorder()
		CreateTrade(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test insert trade not valid cookie")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade not present cookie"), func(t *testing.T) {
		params := []byte(`{
				"Usercode": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X",
				"Exchange": "Binance",
				"FirstPair": 1,
				"SecondPair": 1027,
			}`)
		req := httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		CreateTrade(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test insert trade not present cookie")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade missing subtrades"), func(t *testing.T) {
		params := []byte(`{
			"Exchange": "Binance",
			"FirstPair": 1,
			"SecondPair": 1027,
			"Subtrades": []
		}`)
		req := httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		CreateTrade(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed insert trade missing subtrades")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade invalid coinid"), func(t *testing.T) {
		params := []byte(`{
			"Exchange": "Binance",
			"FirstPair": "1",
			"SecondPair": 1027,
			"Subtrades": [
				{
					"CreatedAt": "2021-10-01T19:39",
					"Type": "BUY",
					"Reason": "Volume",
					"Quantity": 1,
					"AvgPrice": 30000,
					"Total": 30000,
				}]}`)
		req := httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		CreateTrade(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed insert trade invalid coinid")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade invalid subtrades"), func(t *testing.T) {
		Db.Exec(
			`INSERT INTO users (
			wallet, username, privacy,
			createdat, updatedat)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86T', 'ZTZTZT',
			'all', current_timestamp, current_timestamp);`)
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86T"}
		session, _ := user.InsertSession("web")

		params := []byte(`{
			"Exchange": "Binance",
			"FirstPair": 1,
			"SecondPair": 1027,
			"Subtrades": [
				{
					"CreatedAt": "2021-10-1000ABCABD",
					"Type": "BUY",
					"Reason": "Volume",
					"Quantity": 1,
					"AvgPrice": 30000,
					"Total": 30000
				}]}`)
		req := httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		CreateTrade(w, req)
		var count_trades int
		err := Db.QueryRow(`SELECT COUNT(code) FROM trades WHERE userwallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86T'`).Scan(&count_trades)
		if err != nil || count_trades > 0 {
			t.Fatal("Failed insert trade invalid subtrades")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestInsertTrade(t *testing.T) {
	// <setup code>
	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(1000, 'TestCoin', 'TC', 'testcoin'),
			(1001, 'TestCoin2', 'TC2', 'testcoin')`)

	Db.Exec(
		`INSERT INTO users (
			wallet, username, privacy,
			createdat, updatedat)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'jsjsjsj',
			'all', current_timestamp, current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test insert valid trade"), func(t *testing.T) {
		new_trade := NewTrade{
			Exchange:     "Binance",
			FirstPairId:  1000,
			SecondPairId: 1001,
			UserWallet:   "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X",
		}
		err := new_trade.InsertTrade()
		if err != nil {
			t.Fatal("Failed insert valid trade")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade with not existing user"), func(t *testing.T) {
		new_trade := NewTrade{
			Exchange:     "Binance",
			FirstPairId:  1000,
			SecondPairId: 1001,
			UserWallet:   "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86T",
		}
		err := new_trade.InsertTrade()
		if err == nil {
			t.Fatal("Failed insert trade with not existing user")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade with not existing coinid"), func(t *testing.T) {
		new_trade := NewTrade{
			Exchange:     "Binance",
			FirstPairId:  1000,
			SecondPairId: 1001191,
			UserWallet:   "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X",
		}
		err := new_trade.InsertTrade()
		if err == nil {
			t.Fatal("Failed insert trade with not existing coinid")
		}
	})

	t.Run(fmt.Sprintf("Test insert trade with fully empty trade"), func(t *testing.T) {
		new_trade := NewTrade{}
		err := new_trade.InsertTrade()
		if err == nil {
			t.Fatal("Failed insert trade with fully empty trade")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestDeleteTrade(t *testing.T) {
	// <setup code>
	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(367213, 'TestCoin', 'TC', 'testcoin'),
			(123123, 'TestCoin2', 'TC2', 'testcoin')`)

	Db.Exec(
		`INSERT INTO users (
			wallet, username, privacy, createdat, updatedat)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM',
			'all', current_timestamp, current_timestamp);`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
	session, _ := user.InsertSession("web")

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delete_trade", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		DeleteTrade(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test delete trade, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test empty tradecode"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delete_trade", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		DeleteTrade(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test delete, empty tradecode")
		}
	})

	t.Run(fmt.Sprintf("Test successfully delete trade"), func(t *testing.T) {
		Db.Exec(`
			INSERT INTO trades(
				code, userwallet, createdat, updatedat,
				firstpair, secondpair
			) VALUES (
				'IUIUIUIU', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
				current_timestamp, 9999, 8888);`)
		Db.Exec(`
			INSERT INTO subtrades(
				code, tradecode, userwallet, createdat, updatedat)
				VALUES ('fjhfdjsa', 'IUIUIUIU', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X',
				current_timestamp, current_timestamp);`)

		req := httptest.NewRequest("GET", "/delete_trade", nil)
		vars := map[string]string{
			"tradecode": "IUIUIUIU",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		DeleteTrade(w, req)
		var tradecode string
		_ = Db.QueryRow(`SELECT code FROM trades WHERE code = 'IUIUIUIU'`).Scan(&tradecode)
		if tradecode != "" {
			t.Fatal("Failed deleting trade")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
