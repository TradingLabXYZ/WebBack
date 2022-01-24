package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSelectSubscriptionMonthlyPrice(t *testing.T) {
	// <setup code>
	Db.Exec(`
		INSERT INTO users (
			wallet,
			username,
			privacy,
			plan,
			createdat,
			updatedat)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A',
			'r',
			'all',
			'basic',
			current_timestamp,
			current_timestamp);`)

	Db.Exec(`
			INSERT INTO smartcontractevents (
				createdat,
				transaction,
				contract,
				name,
				signature,
				sender,
				payload)
			VALUES(
				current_timestamp,
				'txABC',
				'contractCDE',
				'ChangePlan',
				'ABC',
				'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A',
				'{"Value": 12}');`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
	session, _ := user.InsertSession()

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/subscription", nil)
		vars := map[string]string{
			"wallet": "",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectSubscriptionMonthlyPrice(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test subscription, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test invalid user"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/subscription", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86Z",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		SelectSubscriptionMonthlyPrice(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test subscription, invalid user")
		}
	})

	t.Run(fmt.Sprintf("Test correct subscription plan"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/subscription", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		SelectSubscriptionMonthlyPrice(w, req)
		bytes_body, _ := ioutil.ReadAll(w.Body)
		var monthlyPlan int
		_ = json.Unmarshal([]byte(bytes_body), &monthlyPlan)
		if monthlyPlan != 12 {
			t.Fatal("Failed correct subscription plan")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM smartcontractevents WHERE 1 = 1;`)
}
