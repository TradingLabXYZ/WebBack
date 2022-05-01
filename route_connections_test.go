package main

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestUpdateFollower(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (
			wallet, totalcounttrades, totalportfolio, totalreturn, totalroi, tradeqtyavailable, tradevalue,
			tradereturn, traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
	session, _ := user.InsertSession("web", "Europe|Berlin")
	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/follow", nil)
		vars := map[string]string{
			"wallet": "1234",
			"status": "false",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		UpdateFollower(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed wrong header")
		}
	})
	t.Run(fmt.Sprintf("Test empty status"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/follow", nil)
		vars := map[string]string{
			"wallet": "1234",
			"status": "",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateFollower(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed empty status")
		}
	})
	t.Run(fmt.Sprintf("Test correctly set followers TRUE"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/follow", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B",
			"status": "false",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateFollower(w, req)
		var is_follower bool
		_ = Db.QueryRow(`
				SELECT TRUE
				FROM followers
				WHERE followfrom = $1
				AND followto = $2;`,
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A",
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B").Scan(&is_follower)
		if !is_follower {
			t.Fatal("Failed correctly set followers TRUE")
		}
	})
	t.Run(fmt.Sprintf("Test correctly set followers FALSE"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/follow", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B",
			"status": "true",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateFollower(w, req)
		var is_follower bool
		_ = Db.QueryRow(`
				SELECT TRUE
				FROM followers
				WHERE followfrom = $1
				AND followto = $2;`,
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A",
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B").Scan(&is_follower)
		if is_follower {
			t.Fatal("Failed correctly set followers FALSE")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestUpdateSubscribers(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (
			wallet, totalcounttrades, totalportfolio, totalreturn, totalroi, tradeqtyavailable, tradevalue,
			tradereturn, traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)

	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
	session, _ := user.InsertSession("web", "Europe|Berlin")
	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/subscribe", nil)
		vars := map[string]string{
			"wallet": "1234",
			"status": "false",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		UpdateSubscriber(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed wrong header")
		}
	})
	t.Run(fmt.Sprintf("Test empty status"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/subscribe", nil)
		vars := map[string]string{
			"wallet": "1234",
			"status": "",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateSubscriber(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed empty status")
		}
	})
	t.Run(fmt.Sprintf("Test correctly set subscriber TRUE"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/follow", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B",
			"status": "false",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateSubscriber(w, req)
		var is_subscriber bool
		_ = Db.QueryRow(`
				SELECT TRUE
				FROM subscribers
				WHERE subscribefrom = $1
				AND subscribeto = $2;`,
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A",
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B").Scan(&is_subscriber)
		if !is_subscriber {
			t.Fatal("Failed correctly set subscriber TRUE")
		}
	})
	t.Run(fmt.Sprintf("Test correctly set subscriber FALSE"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/subscribe", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B",
			"status": "true",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateSubscriber(w, req)
		var is_subscriber bool
		_ = Db.QueryRow(`
				SELECT TRUE
				FROM subscribers
				WHERE subscribefrom = $1
				AND subscribeto = $2;`,
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A",
			"0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B").Scan(&is_subscriber)
		if is_subscriber {
			t.Fatal("Failed correctly set subscriber FALSE")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestSelectConnection(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'subscribers', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userc', 'all', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', 'userd', 'followers', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (
			wallet, totalcounttrades, totalportfolio, totalreturn, totalroi, tradeqtyavailable, tradevalue,
			tradereturn, traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)
	Db.Exec(
		`INSERT INTO followers (followfrom, followto, createdat)
		VALUES
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', current_timestamp);`)
	user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
	session, _ := user.InsertSession("web", "Europe|Berlin")
	_ = session

	// <test code>
	t.Run(fmt.Sprintf("Test invalid user"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_connections", nil)
		vars := map[string]string{
			"wallet": "1234notexists",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectConnections(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed wrong header")
		}
	})
	t.Run(fmt.Sprintf("Test correctely accessing from user not loggedin"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_connections", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		SelectConnections(w, req)
		res := w.Result()
		if res.StatusCode != 200 {
			t.Fatal("Failed correctely accessing from user not loggedin")
		}
	})
	t.Run(fmt.Sprintf("Test response correctly not returning data"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_connections", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		SelectConnections(w, req)
		bytes_body, _ := ioutil.ReadAll(w.Body)
		body := string(bytes_body)
		if !strings.Contains(body, `"Status":"KO"`) {
			t.Fatal("Failed response correctly not returning data")
		}
	})
	t.Run(fmt.Sprintf("Test response correctly returning data"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/get_connections", nil)
		vars := map[string]string{
			"wallet": "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		SelectConnections(w, req)
		bytes_body, _ := ioutil.ReadAll(w.Body)
		body := string(bytes_body)
		if strings.Contains(body, `"Status":"KO"`) {
			t.Fatal("Failed response correctly returning data")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
