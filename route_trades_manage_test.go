package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestCreateTrade(t *testing.T) {
	fmt.Println("Start TestCreateTrade")

	Db.Exec(`
		INSERT INTO coins (coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'bitcoin'),
			(1027, 'Ethereum', 'ETH', 'ethereum')`)

	// Test insert trade with 1 subtrade
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
			'12e91acdae83c7225a0b16950a4268083b3eae7f',
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
		t.Fatal("Failed TestLogin empty body")
	}
}
