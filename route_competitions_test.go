package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestInsertPrediction(t *testing.T) {
	// <setup code>
	params := []byte(`{
		"Competition": "first_competition",
		"Source": "test",
		"Prediction": 1999
	}`)

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
	Db.Exec(
		`INSERT INTO competitions (
			name, submissionendedat, submissionstartedat,
			competitionstartedat, competitionendedat)
		VALUES (
			'first_competition', current_timestamp,
			current_timestamp, current_timestamp, current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("POST", "/insert_prediction", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		InsertPrediction(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed insert prediction, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test origin not web"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("api")
		req := httptest.NewRequest("POST", "/insert_prediction", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		InsertPrediction(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed insert prediction, origin not web")
		}
	})

	t.Run(fmt.Sprintf("Test successfully insert prediction"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web")
		req := httptest.NewRequest("POST", "/insert_prediction", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		InsertPrediction(w, req)
		var prediction string
		_ = Db.QueryRow(`
			SELECT payload
			FROM submissions
			WHERE userwallet = $1;`,
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A").Scan(&prediction)

		if !strings.Contains(prediction, "1999") {
			t.Fatal("Failed successfully inserting prediction")
		}
	})

	t.Run(fmt.Sprintf("Test successfully update prediction"), func(t *testing.T) {
		new_params := []byte(`{
			"Competition": "first_competition",
			"Source": "test",
			"Prediction": 1111
		}`)
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web")
		req := httptest.NewRequest("POST", "/insert_prediction", bytes.NewBuffer(new_params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		InsertPrediction(w, req)
		var prediction string
		_ = Db.QueryRow(`
			SELECT payload
			FROM submissions
			WHERE userwallet = $1
			ORDER BY updatedat DESC
			LIMIT 1;`,
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A").Scan(&prediction)

		if !strings.Contains(prediction, "1111") {
			t.Fatal("Failed successfully updating prediction")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM competitions WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM submissions WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestSelectPrediction(t *testing.T) {
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
	Db.Exec(
		`INSERT INTO competitions (
			name, submissionendedat, submissionstartedat,
			competitionstartedat, competitionendedat)
		VALUES (
			'first_competition', current_timestamp,
			current_timestamp, current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO submissions (
			competitionname, userwallet, payload, updatedat)
		VALUES (
			'first_competition', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A',
			'{"Prediction": 144.99}', current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_prediction", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectPrediction(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed selecting prediction, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test origin not web"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("api")
		req := httptest.NewRequest("GET", "/select_prediction", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectPrediction(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed selecting prediction, origin not web")
		}
	})

	t.Run(fmt.Sprintf("Test successfully select prediction"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web")
		req := httptest.NewRequest("GET", "/select_prediction", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		SelectPrediction(w, req)
		prediction, _ := ioutil.ReadAll(w.Body)
		prediction_s := string(prediction)

		if !strings.Contains(prediction_s, "144.99") {
			t.Fatal("Failed successfully selecting prediction")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM competitions WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestDeletePrediction(t *testing.T) {
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
	Db.Exec(
		`INSERT INTO competitions (
			name, submissionendedat, submissionstartedat,
			competitionstartedat, competitionendedat)
		VALUES (
			'first_competition', current_timestamp,
			current_timestamp, current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO submissions (
			competitionname, userwallet, payload, updatedat)
		VALUES (
			'first_competition', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A',
			'{"prediction": 111.11}', current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delete_prediction", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		DeletePrediction(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed deleting prediction, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test origin not web"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("api")
		req := httptest.NewRequest("GET", "/deleting_prediction", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		DeletePrediction(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed deleting prediction, origin not web")
		}
	})

	t.Run(fmt.Sprintf("Test successfully delete prediction"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web")
		req := httptest.NewRequest("GET", "/delete_prediction", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		DeletePrediction(w, req)
		prediction, _ := ioutil.ReadAll(w.Body)
		prediction_s := string(prediction)

		if strings.Contains(prediction_s, "111.11") {
			t.Fatal("Failed successfully deleting prediction")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM competitions WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestGetCountPartecipants(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'userd', 'all', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userx', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (wallet, totalcounttrades, totalportfolio,
			totalreturn, totalroi, tradeqtyavailable, tradevalue, tradereturn,
			traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)
	Db.Exec(
		`INSERT INTO competitions (
			name, submissionendedat, submissionstartedat,
			competitionstartedat, competitionendedat)
		VALUES (
			'first_competition', current_timestamp,
			current_timestamp, current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO submissions (
			competitionname, userwallet, payload, updatedat)
		VALUES
			('first_competition', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A',
			'{"prediction": 144.99}', current_timestamp),
			('first_competition', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B',
			'{"prediction": 155.99}', current_timestamp);`)

	t.Run(fmt.Sprintf("Test successfully count predictions"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_count_partecipants", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		GetCountPartecipants(w, req)
		partecipants, _ := ioutil.ReadAll(w.Body)
		partecipants_s := string(partecipants)
		if !strings.Contains(partecipants_s, "2") {
			t.Fatal("Failed successfully counting partecipants")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM competitions WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestGetPartecipants(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (wallet, totalcounttrades, totalportfolio,
			totalreturn, totalroi, tradeqtyavailable, tradevalue, tradereturn,
			traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)
	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin');`)
	Db.Exec(`
		INSERT INTO lastprices (
			updatedat, coinid, price)
		VALUES
			(current_timestamp, 1, 45000);`)
	Db.Exec(
		`INSERT INTO competitions (
			name, submissionendedat, submissionstartedat,
			competitionstartedat, competitionendedat)
		VALUES (
			'first_competition', current_timestamp,
			current_timestamp, current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO submissions (
			competitionname, userwallet, payload, updatedat)
		VALUES
			('first_competition', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A',
			'{"Prediction": 46000}', current_timestamp),
			('first_competition', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B',
			'{"Prediction": 47000}', current_timestamp);`)

	t.Run(fmt.Sprintf("Test successfully get partipants"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_partecipants", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		vars := map[string]string{
			"competition": "first_competition",
		}
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		GetPartecipants(w, req)
		partecipants, _ := ioutil.ReadAll(w.Body)
		type Prediction struct {
			CreatedAt      string
			Username       string
			Wallet         string
			ProfilePicture string
			Prediction     string
			BtcPrice       string
			DeltaPerc      float32
			AbsDeltaPrice  string
		}

		predictions := []Prediction{}
		_ = json.Unmarshal(partecipants, &predictions)
		if predictions[0].Username != "usera" {
			t.Fatal("Failed successfully get partecipants")
		}
		if predictions[0].DeltaPerc != 2.22 {
			t.Fatal("Failed successfully get partecipants")
		}
		if predictions[1].DeltaPerc != 4.44 {
			t.Fatal("Failed successfully get partecipants")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM competitions WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
