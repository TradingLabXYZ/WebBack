package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestInsertSubTrade(t *testing.T) {
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
			(1000, 'TestCoin', 'TC', 'testcoin'),
			(1001, 'TestCoin2', 'TC2', 'testcoin')`)

	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES (
			'MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
			current_timestamp, 1000, 1001, TRUE);`)

	// <test code>
	t.Run(fmt.Sprintf("Test successfully insert subtrade"), func(t *testing.T) {
		new_subtrades := []NewSubtrade{
			{
				CreatedAt: "2021-01-01",
				Type:      "BUY",
				Reason:    "Test",
				Quantity:  "1",
				AvgPrice:  "1",
				Total:     "1",
			},
		}

		new_trade := NewTrade{
			Exchange:     "Test",
			FirstPairId:  1000,
			SecondPairId: 1001,
			Subtrades:    new_subtrades,
			Code:         "MBMBMBM",
			UserWallet:   "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X",
		}

		err := new_trade.InsertSubTrades()
		if err != nil {
			t.Fatal("Failed successfully insert subtrade")
		}
	})

	t.Run(fmt.Sprintf("Test insert subtrade wrong user"), func(t *testing.T) {
		new_subtrades := []NewSubtrade{
			{
				CreatedAt: "2021-01-01",
				Type:      "BUY",
				Reason:    "Test",
				Quantity:  "1",
				AvgPrice:  "1",
				Total:     "1",
			},
		}

		new_trade := NewTrade{
			Exchange:     "Test",
			FirstPairId:  1000,
			SecondPairId: 1001,
			Subtrades:    new_subtrades,
			Code:         "MBMBMBM",
			UserWallet:   "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86T",
		}

		err := new_trade.InsertSubTrades()
		if err == nil {
			t.Fatal("Failed insert subtrade wrong user")
		}
	})

	t.Run(fmt.Sprintf("Test insert subtrade wrong createdat"), func(t *testing.T) {
		new_subtrades := []NewSubtrade{
			{
				CreatedAt: "2021-01-01TESTETEST",
				Type:      "BUY",
				Reason:    "Test",
				Quantity:  "1",
				AvgPrice:  "1",
				Total:     "1",
			},
		}

		new_trade := NewTrade{
			Exchange:     "Test",
			FirstPairId:  1000,
			SecondPairId: 1001,
			Subtrades:    new_subtrades,
			Code:         "MBMBMBM",
			UserWallet:   "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X",
		}

		err := new_trade.InsertSubTrades()
		if err == nil {
			t.Fatal("Failed insert subtrade wrong createdat")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestCreateSubTrade(t *testing.T) {
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
			(1000, 'TestCoin', 'TC', 'testcoin'),
			(1001, 'TestCoin2', 'TC2', 'testcoin')`)

	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES (
			'MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
			current_timestamp, 1000, 1001, TRUE);`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
	session, _ := user.InsertSession()

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/insert_subtrade", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		CreateSubtrade(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test insert subtrade, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test empty tradecode"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/insert_subtrade", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		CreateSubtrade(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test insert subtrade, empty tradecode")
		}
	})

	t.Run(fmt.Sprintf("Test successfully insert subtrade"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/insert_subtrade", nil)
		vars := map[string]string{
			"tradecode": "MBMBMBM",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		CreateSubtrade(w, req)
		CreateSubtrade(w, req)
		var count_subtrades int
		_ = Db.QueryRow(`
			SELECT COUNT(code)
			FROM subtrades
			WHERE tradecode = $1;`, "MBMBMBM").Scan(&count_subtrades)
		if count_subtrades != 2 {
			t.Fatal("Failed successfully insert subtrade")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestUpdateSubTrade(t *testing.T) {
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
			(1000, 'TestCoin', 'TC', 'testcoin'),
			(1001, 'TestCoin2', 'TC2', 'testcoin')`)

	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES (
			'MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
			current_timestamp, 1000, 1001, TRUE);`)

	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			quantity, avgprice, total, reason)
		VALUES (
			'SISISIS', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
			current_timestamp, 1, 1, 1, 'TESTART');`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
	session, _ := user.InsertSession()

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("POST", "/insert_subtrade", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		UpdateSubtrade(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test insert subtrade, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test empty payload"), func(t *testing.T) {
		req := httptest.NewRequest("POST", "/insert_subtrade", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateSubtrade(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test insert subtrade, empty payload")
		}
	})

	t.Run(fmt.Sprintf("Test wrong payload"), func(t *testing.T) {
		params := []byte(`{
			"Code": "SISISIS",
			"CreatedAt": "2021-10-10WROOOO"
			"Type": "BUY",
			"Reason": "Test",
			"Quantity": "10",
			"AvgPrice": "10",
			"Total": "10"
		}`)
		req := httptest.NewRequest("POST", "/insert_subtrade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateSubtrade(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test insert subtrade, empty payload")
		}
	})

	t.Run(fmt.Sprintf("Test successfully update subtrade"), func(t *testing.T) {
		params := []byte(`{
			"Code": "SISISIS",
			"CreatedAt": "2021-10-01T19:39",
			"Type": "BUY",
			"Reason": "This is it",
			"Quantity": "10",
			"AvgPrice": "10",
			"Total": "10"
		}`)
		req := httptest.NewRequest("POST", "/insert_subtrade", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateSubtrade(w, req)
		var expected_reason string
		_ = Db.QueryRow(`
			SELECT reason
			FROM subtrades
			WHERE code = $1;`, "SISISIS").Scan(&expected_reason)
		if expected_reason != "This is it" {
			t.Fatal("Failed successfully insert subtrade")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestDeleteSubTrade(t *testing.T) {
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
			(1000, 'TestCoin', 'TC', 'testcoin'),
			(1001, 'TestCoin2', 'TC2', 'testcoin')`)

	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES (
			'MBMBMBM', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', current_timestamp,
			current_timestamp, 1000, 1001, TRUE);`)

	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			quantity, avgprice, total, reason)
		VALUES (
			'SISISIS', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X', 'MBMBMBM', current_timestamp,
			current_timestamp, 1, 1, 1, 'TESTART');`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86X"}
	session, _ := user.InsertSession()

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delete_subtrade", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		DeleteSubtrade(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test insert subtrade, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test empty subtradecode"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delete_subtrade", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		DeleteSubtrade(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test deleting subtrade, empty subtradecode")
		}
	})

	t.Run(fmt.Sprintf("Test successfully delete subtrade"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delete_subtrade", nil)
		vars := map[string]string{
			"subtradecode": "SISISIS",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		DeleteSubtrade(w, req)
		var subtrade_code string
		_ = Db.QueryRow(`
			SELECT code
			FROM subtrades
			WHERE code = $1;`, "SISISIS").Scan(&subtrade_code)
		if subtrade_code != "" {
			t.Fatal("Failed successfully deleting subtrade")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
