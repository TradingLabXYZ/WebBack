package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestSelectPairs(t *testing.T) {
	// <setup code>
	type TempPairInfo struct {
		Name   string
		Slug   string
		Symbol string
		CoinId int
	}
	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1000, 'TestCoinA', 'A', 'TestA'),
			(1001, 'TestCoinB', 'B', 'TestB'),
			(1002, 'TestCoinC', 'C', 'TestC'),
			(1003, 'TestCoinD', 'D', 'TestD')`)

	// <test code>
	t.Run(fmt.Sprintf("Test successfully extract pairs info"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_pairs", nil)
		w := httptest.NewRecorder()
		SelectPairs(w, req)
		readBuf, _ := ioutil.ReadAll(w.Body)
		var sec map[string]TempPairInfo
		_ = json.Unmarshal([]byte(*&readBuf), &sec)
		if sec["1000"].Name != "TestCoinA" {
			t.Error("Failed successfully extract pairs info A")
		}
		if sec["1001"].Name != "TestCoinB" {
			t.Error("Failed successfully extract pairs info B")
		}
		if sec["1002"].Slug != "TestC" {
			t.Error("Failed successfully extract pairs info C")
		}
		if sec["1003"].Slug != "TestD" {
			t.Error("Failed successfully extract pairs info D")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
}

func TestSelectPairRatio(t *testing.T) {
	// <setup code>
	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1000, 'TestCoinA', 'A', 'TestA'),
			(1001, 'TestCoinB', 'B', 'TestB')`)

	Db.Exec(`
		INSERT INTO lastprices (
			updatedat, coinid, price)
		VALUES
			(current_timestamp, 1000, 200),
			(current_timestamp, 1001, 100);`)

	// <test code>
	t.Run(fmt.Sprintf("Test passing empty value"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_pair_ratio", nil)
		vars := map[string]string{
			"firstPairCoinId":  "",
			"secondPairCoinId": "1001",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectPairRatio(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed passing empty values")
		}
	})
	t.Run(fmt.Sprintf("Test passing wrong coinid"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_pair_ratio", nil)
		vars := map[string]string{
			"firstPairCoinId":  "ABC",
			"secondPairCoinId": "1001",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectPairRatio(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed passing wrong coinid")
		}
	})
	t.Run(fmt.Sprintf("Test successfully extract pair ratio"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_pair_ratio", nil)
		vars := map[string]string{
			"firstPairCoinId":  "1000",
			"secondPairCoinId": "1001",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectPairRatio(w, req)
		ratio, _ := ioutil.ReadAll(w.Body)
		ratio_s := string(ratio)

		ratio_s = strings.Replace(ratio_s, "\n", "", -1)
		if s, err := strconv.ParseFloat(ratio_s, 32); err == nil {
			if s != 0.5 {
				t.Error("Failed successfully extract pair ratio")
			}
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
}
