package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestCreateTrade(t *testing.T) {

	// Preliminary data upload
	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'bitcoin'),
			(1027, 'Ethereum', 'ETH', 'ethereum')`)

	var user_id int
	_ = Db.QueryRow(
		`INSERT INTO users (
			code,
			email,
			username,
			password,
			privacy,
			plan,
			createdat,
			updatedat)
		VALUES (
			'ABABAB',
			'r@r.r',
			'r',
			'testpassword',
			'all',
			'basic',
			current_timestamp,
			current_timestamp)
		RETURNING id;`).Scan(&user_id)

	user := User{
		Id:    user_id,
		Email: "r@r.r",
	}
	session, _ := user.CreateSession()

	// Test insert trade with subtrade
	params := []byte(`{
		"Usercode": "ABABAB",
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
			}
		]
	}`)
	req := httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w := httptest.NewRecorder()
	CreateTrade(w, req)
	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatal("Failed test insert trade with subtrade")
	}

	// Test insert trade not valid cookie
	params = []byte(`{
		"Usercode": "ABABAB",
		"Exchange": "Binance",
		"FirstPair": 1,
		"SecondPair": 1027,
	}`)
	req = httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
	req.Header.Set("Authorization", "Bearer sessionId=testwrongcookie")
	w = httptest.NewRecorder()
	CreateTrade(w, req)
	res = w.Result()
	if res.StatusCode != 404 {
		t.Fatal("Failed test insert trade not valid cookie")
	}

	// Test insert trade not present cookie
	params = []byte(`{
		"Usercode": "ABABAB",
		"Exchange": "Binance",
		"FirstPair": 1,
		"SecondPair": 1027,
	}`)
	req = httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
	req.Header.Set("Authorization", "Bearer sessionId=")
	w = httptest.NewRecorder()
	CreateTrade(w, req)
	res = w.Result()
	if res.StatusCode != 404 {
		t.Fatal("Failed test insert trade not present cookie")
	}

	// Test insert trade bad payload
	params = []byte(`{
		"Usercode": "THIS USER DOES NOT EXISTS",
		"Exchange": "Binance",
		"FirstPair": 1,
		"SecondPair": 1027,
	}`)
	req = httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	CreateTrade(w, req)
	res = w.Result()
	if res.StatusCode != 400 {
		t.Fatal("Failed insert trade bad payload")
	}

	// Test insert trade bad payload 2
	params = []byte(`{
		"Usercode": "ABABAB",
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
				"Total": 30000
			}
		]
	}`)
	req = httptest.NewRequest("POST", "/insert_trade", bytes.NewBuffer(params))
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	CreateTrade(w, req)
	res = w.Result()
	if res.StatusCode != 400 {
		t.Fatal("Failed insert trade bad payload")
	}
}

func TestInsertTrade(t *testing.T) {

	// Preliminary data upload
	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(1000, 'TestCoin', 'TC', 'testcoin'),
			(1001, 'TestCoin2', 'TC2', 'testcoin')`)

	Db.Exec(
		`INSERT INTO users (
			code,
			email,
			username,
			password,
			privacy,
			plan,
			createdat,
			updatedat)
		VALUES (
			'JFJFJF',
			'jsjsjs@r.r',
			'jsjsjsj',
			'testpassword',
			'all',
			'basic',
			current_timestamp,
			current_timestamp)
		RETURNING id;`)

	// Test insert valid trade
	new_trade := NewTrade{
		Exchange:     "Binance",
		FirstPairId:  1000,
		SecondPairId: 1001,
		Usercode:     "JFJFJF",
	}
	err := new_trade.InsertTrade()
	if err != nil {
		t.Fatal("Failed insert valid trade")
	}

	// Test insert trade with not existing user
	new_trade = NewTrade{
		Exchange:     "Binance",
		FirstPairId:  1000,
		SecondPairId: 1001,
		Usercode:     "JFJFJFJSJSJSJ",
	}
	err = new_trade.InsertTrade()
	if err == nil {
		t.Fatal("Failed insert trade with not existing user")
	}

	// Test insert trade with not existing coinid
	new_trade = NewTrade{
		Exchange:     "Binance",
		FirstPairId:  1000,
		SecondPairId: 1001191,
		Usercode:     "JFJFJF",
	}
	err = new_trade.InsertTrade()
	if err == nil {
		t.Fatal("Failed insert trade with not existing coinid")
	}

	// Test insert trade with fully empty trade
	new_trade = NewTrade{}
	err = new_trade.InsertTrade()
	if err == nil {
		t.Fatal("Failed insert trade with fully empty trade")
	}
}

func TestCloseTrade(t *testing.T) {

	var user_id int
	_ = Db.QueryRow(
		`INSERT INTO users (
			code,
			email,
			username,
			password,
			privacy,
			plan,
			createdat,
			updatedat)
		VALUES (
			'QPQPQPQ',
			'pqpqpq@r.r',
			'pqpqpqp',
			'testpassword',
			'all',
			'basic',
			current_timestamp,
			current_timestamp)
		RETURNING id;`).Scan(&user_id)

	user := User{
		Id:    user_id,
		Email: "pqpqpq@r.r",
	}
	session, _ := user.CreateSession()

	var trade_code string // SONO ARRIVATO QUI
	_ = Db.QueryRow(`
		INSERT INTO trades	

	`).Scan(&trade_code)
	fmt.Println(trade_code)

	req := httptest.NewRequest("GET", "/close_trade"+trade_code, nil)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w := httptest.NewRecorder()
	CloseTrade(w, req)
	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatal("Failed test insert trade with subtrade")
	}
}

/* new_subtrades = []NewSubtrade{
	{
		CreatedAt: "2021dsd",
		Type:      "BUY",
		Reason:    "Volume",
		Quantity:  "1",
		AvgPrice:  "30000",
		Total:     "30000",
	},
}
new_trade = NewTrade{
	Exchange:     "Binance",
	FirstPairId:  1000,
	SecondPairId: 1001,
	Usercode:     "JFJFJF",
	Subtrades:    new_subtrades,
} */
