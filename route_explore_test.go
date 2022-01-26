package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSelectExplore(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			wallet, username, privacy,
			plan, createdat, updatedat)
		VALUES
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera',
			'all', 'basic', current_timestamp, current_timestamp);`)

	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(1000, 'USDC', 'USDC', 'usdc');`)

	Db.Exec(`
		INSERT INTO lastprices (
			createdat, coinid, price)
		VALUES
			(current_timestamp, 1, 80000),
			(current_timestamp, 1000, 1);`)

	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('useratr1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp, current_timestamp, 1000, 1),
		('useratr2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp,current_timestamp, 1, 1000);`)

	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES
		('userasub1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr1', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART'),
		('userasub2', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr2', current_timestamp, current_timestamp, 'BUY', 1, 50000, 50000, 'TESTART'),
		('userasub3', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr2', current_timestamp, current_timestamp, 'SELL', 1, 60000, 60000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test empty offset"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_explore", nil)
		vars := map[string]string{
			"offset": "",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectExplore(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed empty offset")
		}
	})
	t.Run(fmt.Sprintf("Test offset not divisible by 10"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_explore", nil)
		vars := map[string]string{
			"offset": "13",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectExplore(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed offset not divisible by 10")
		}
	})
	t.Run(fmt.Sprintf("Test valid offset"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_explore", nil)
		vars := map[string]string{
			"offset": "0",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectExplore(w, req)
		res := w.Result()
		if res.StatusCode != 200 {
			t.Fatal("Failed valid offset")
		}
	})
	t.Run(fmt.Sprintf("Test valid explore response"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_explore", nil)
		vars := map[string]string{
			"offset": "0",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectExplore(w, req)
		bytes_body, _ := ioutil.ReadAll(w.Body)
		var slice []map[string]interface{}
		_ = json.Unmarshal([]byte(bytes_body), &slice)
		if len(slice) != 3 {
			t.Fatal("Failed valid explore response, explore length")
		}
		if slice[0]["firstpair"] != 1000.0 {
			t.Fatal("Failed valid explore response, first pair")
		}
		if slice[1]["firstpair"] != 1.0 {
			t.Fatal("Failed valid explore response, first pair")
		}
		if slice[0]["secondpairsymbol"] != "BTC" {
			t.Fatal("Failed valid explore response, second pair symold")
		}
		timeago := slice[0]["timeago"]
		if timeago != "0 seconds ago" && timeago != "1 seconds ago" {
			t.Fatal("Failed valid explore response, time ago")
		}
	})
	t.Run(fmt.Sprintf("Test empty explore when data mistake"), func(t *testing.T) {
		Db.Exec(`DELETE FROM coins WHERE coinid = 1`)
		req := httptest.NewRequest("GET", "/get_explore", nil)
		vars := map[string]string{
			"offset": "10",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectExplore(w, req)
		bytes_body, _ := ioutil.ReadAll(w.Body)
		var slice []map[string]interface{}
		_ = json.Unmarshal([]byte(bytes_body), &slice)
		if len(slice) != 0 {
			t.Fatal("Failed empty explore when data mistake")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
}
