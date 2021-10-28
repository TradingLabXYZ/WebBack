package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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
	if res.StatusCode != 401 {
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
	if res.StatusCode != 401 {
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

func TestChangeTradeStatus(t *testing.T) {

	// Preliminary data upload
	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(9999, 'TestCoin', 'TC', 'testcoin'),
			(8888, 'TestCoin2', 'TC2', 'testcoin')`)

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

	// Test wrong header
	req := httptest.NewRequest("GET", "/change_trade", nil)
	req.Header.Set("Authorization", "Bearer sessionId=")
	w := httptest.NewRecorder()
	ChangeTradeStatus(w, req)
	res := w.Result()
	if res.StatusCode != 401 {
		t.Fatal("Failed test change trade status, wrong header")
	}

	// Test empty tradecode
	req = httptest.NewRequest("GET", "/change_trade", nil)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	ChangeTradeStatus(w, req)
	res = w.Result()
	if res.StatusCode != 400 {
		t.Fatal("Failed test change trade status, empty tradecode")
	}

	// Test not existing tradecode
	req = httptest.NewRequest("GET", "/change_trade", nil)
	vars := map[string]string{
		"tradecode": "TEGDGDHGKJEHS",
		"tostatus":  "true",
	}
	req = mux.SetURLVars(req, vars)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	ChangeTradeStatus(w, req)
	res = w.Result()
	if res.StatusCode != 400 {
		t.Fatal("Failed test change trade status, not existing tradecode")
	}

	// Test change trade status to false
	Db.Exec(`
		INSERT INTO trades(
			code,
			usercode,
			createdat,
			updatedat,
			firstpair,
			secondpair,
			isopen
		) VALUES (
			'PQPQP',
			'QPQPQPQ',
			current_timestamp,
			current_timestamp,
			9999,
			8888,
			TRUE);`)

	req = httptest.NewRequest("GET", "/change_trade", nil)
	vars = map[string]string{
		"tradecode": "PQPQP",
		"tostatus":  "false",
	}
	req = mux.SetURLVars(req, vars)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	ChangeTradeStatus(w, req)
	isopen := false
	_ = Db.QueryRow(`
		SELECT
			isopen
		FROM trades
		WHERE code = 'PQPQP'`).Scan(isopen)
	if isopen {
		t.Fatal("Failed test change trade status, true to false")
	}

	// Test change trade status to true
	Db.Exec(`
		INSERT INTO trades(
			code,
			usercode,
			createdat,
			updatedat,
			firstpair,
			secondpair,
			isopen
		) VALUES (
			'PQPQP',
			'QPQPQPQ',
			current_timestamp,
			current_timestamp,
			9999,
			8888,
			FALSE);`)

	req = httptest.NewRequest("GET", "/change_trade", nil)
	vars = map[string]string{
		"tradecode": "PQPQP",
		"tostatus":  "true",
	}
	req = mux.SetURLVars(req, vars)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	ChangeTradeStatus(w, req)
	isopen = true
	_ = Db.QueryRow(`SELECT isopen FROM trades WHERE code = 'PQPQP'`).Scan(isopen)
	if !isopen {
		t.Fatal("Failed test change trade status, false to true")
	}
}

func TestDeleteTrade(t *testing.T) {

	// Preliminary data upload
	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(367213, 'TestCoin', 'TC', 'testcoin'),
			(123123, 'TestCoin2', 'TC2', 'testcoin')`)

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
			'MBMBMBM',
			'MBMBMBM@r.r',
			'MBMBMBM',
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
	_ = session

	Db.Exec(`
		INSERT INTO trades(
			code,
			usercode,
			createdat,
			updatedat,
			firstpair,
			secondpair,
			isopen
		) VALUES (
			'MBMBMBM',
			'MBMBMBM',
			current_timestamp,
			current_timestamp,
			9999,
			8888,
			TRUE);`)

	// Test wrong header
	req := httptest.NewRequest("GET", "/delete_trade", nil)
	req.Header.Set("Authorization", "Bearer sessionId=")
	w := httptest.NewRecorder()
	DeleteTrade(w, req)
	res := w.Result()
	if res.StatusCode != 401 {
		t.Fatal("Failed test delete trade, wrong header")
	}

	// Test empty tradecode
	req = httptest.NewRequest("GET", "/delete_trade", nil)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	DeleteTrade(w, req)
	res = w.Result()
	if res.StatusCode != 400 {
		t.Fatal("Failed test delete, empty tradecode")
	}

	// Test not existing tradecode
	req = httptest.NewRequest("GET", "/delete_trade", nil)
	vars := map[string]string{
		"tradecode": "TEGDGDHGKJEHS",
	}
	req = mux.SetURLVars(req, vars)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	DeleteTrade(w, req)
	res = w.Result()
	if res.StatusCode != 400 {
		t.Fatal("Failed delete trade, not existing tradecode")
	}

	// Test successfully delete trade previously created
	req = httptest.NewRequest("GET", "/delete_trade", nil)
	vars = map[string]string{
		"tradecode": "MBMBMBM",
	}
	req = mux.SetURLVars(req, vars)
	req.Header.Set("Authorization", "Bearer sessionId="+session.Uuid)
	w = httptest.NewRecorder()
	DeleteTrade(w, req)
	var tradecode string
	_ = Db.QueryRow(`
		SELECT
			code
		FROM trades
		WHERE code = 'MBMBMBM'`).Scan(tradecode)
	if tradecode != "" {
		t.Fatal("Failed deleting trade")
	}
}
