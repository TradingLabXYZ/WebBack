package main

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestSelectActivity(t *testing.T) {
	// <setup code>
	session_id := "abcdef"
	fake_req := httptest.NewRequest("GET", "/", nil)
	fake_req.Header.Set("Origin", "http://testlocalhost:9000")
	fake_w := httptest.NewRecorder()
	ws, _ := InstanciateTradeWs(fake_w, fake_req)
	c := make(chan TradesSnapshot)

	wallet_1 := "wallet_1"
	observed_1 := User{Wallet: wallet_1}
	observer_1 := User{Wallet: "wallet_2"}
	ws_trade_1 := WsTrade{observer_1, observed_1, session_id, c, ws}
	trades_wss[wallet_1] = append(trades_wss[wallet_1], ws_trade_1)

	wallet_2 := "wallet_2"
	observed_2 := User{Wallet: wallet_2}
	observer_2 := User{Wallet: "wallet_3"}
	ws_trade_2 := WsTrade{observer_2, observed_2, session_id, c, ws}
	trades_wss[wallet_2] = append(trades_wss[wallet_2], ws_trade_2)

	ws_trade_3 := WsTrade{observer_2, observed_1, session_id, c, ws}
	trades_wss[wallet_1] = append(trades_wss[wallet_1], ws_trade_3)

	// <test code>
	t.Run(fmt.Sprintf("Test accessing admin with wrong token"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		vars := map[string]string{
			"token": "randomWrongToken",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectActivity(w, req)
		if w.Code != 400 {
			t.Fatal("Failed accessing admin with wrong token")
		}
	})

	t.Run(fmt.Sprintf("Test accessing admin with correct token"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		vars := map[string]string{
			"token": os.Getenv("ADMIN_TOKEN"),
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectActivity(w, req)
		if w.Code != 200 {
			t.Fatal("Failed acessing admin with correct token")
		}
	})

	t.Run(fmt.Sprintf("Test creating correct output"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		vars := map[string]string{
			"token": os.Getenv("ADMIN_TOKEN"),
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectActivity(w, req)
		bytes_body, _ := ioutil.ReadAll(w.Body)
		body := string(bytes_body)
		if strings.Count(body, "Observer: wallet_2") != 1 {
			t.Fatal("Failed creating correct output: observer wallet_2")
		}
		if strings.Count(body, "Observed: wallet_1") != 2 {
			t.Fatal("Failed creating correct output: observed wallet_1")
		}
		if !strings.Contains(body, "Online Users: 2") {
			t.Fatal("Failed creating correct output: online users")
		}
	})

	// <tear-down code>
	delete(trades_wss, "wallet_1")
	delete(trades_wss, "wallet_2")
}
