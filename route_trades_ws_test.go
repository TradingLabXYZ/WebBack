package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestCheckPrivacy(t *testing.T) {

	// <setup code>
	Db.Exec(
		`INSERT INTO users (code, email, username, password, privacy, plan, createdat, updatedat) VALUES 
		('usera', 'usera@mail.com', 'usera', 'testpassword', 'all', 'basic', current_timestamp, current_timestamp), 
		('userb', 'userb@mail.com', 'userb', 'testpassword', 'private', 'basic', current_timestamp, current_timestamp),
		('userc', 'userc@mail.com', 'userc', 'testpassword', 'followers', 'basic', current_timestamp, current_timestamp), 
		('userd', 'userd@mail.com', 'userd', 'testpassword', 'subscribers', 'basic', current_timestamp, current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test user with privacy ALL is fully visibile"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "usera")
		user_b := User{Code: "userb"}
		session_b, _ := user_b.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_b.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Status != "OK" {
			t.Fatal("Failed test user with privacy ALL is fully visibile")
		}
	})

	t.Run(fmt.Sprintf("Test user not authenticated try to access not ALL users"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userb")
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer sessionId=NOTVALUD")
		status := CheckPrivacy(req, userToSee)
		if status.Status != "KO" {
			t.Fatal("Failed user not authenticated try to access not ALL users")
		}
	})

	t.Run(fmt.Sprintf("Test user PRIVATE always able to see its profile if authenticated"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userb")
		user_b := User{Code: "userb"}
		session_b, _ := user_b.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_b.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Status != "OK" {
			t.Fatal("Failed user PRIVATE always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user FOLLOWERS always able to see its profile if authenticated"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userc")
		user_c := User{Code: "userc"}
		session_c, _ := user_c.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_c.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Status != "OK" {
			t.Fatal("Failed user FOLLOWERS always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user SUBSCRIBERS always able to see its profile if authenticated"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userd")
		user_d := User{Code: "userd"}
		session_d, _ := user_d.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_d.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Status != "OK" {
			t.Fatal("Failed user SUBSCRIBERS always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when PRIVATE"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userb")
		user_a := User{Code: "usera"}
		session_a, _ := user_a.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_a.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Reason != "private" {
			t.Fatal("Failed user cannot access other user when PRIVATE")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when FOLLOWERS and not following"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userc")
		user_a := User{Code: "usera"}
		session_a, _ := user_a.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_a.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Reason != "user is not follower" {
			t.Fatal("Failed user cannot access other user when FOLLOWERS and not following")
		}
	})

	t.Run(fmt.Sprintf("Test user can access other user when FOLLOWERS and yes following"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userc")
		Db.Exec(`
			INSERT INTO followers (followfrom, followto, createdat)
			VALUES ('usera', 'userc', current_timestamp);`)
		user_a := User{Code: "usera"}
		session_a, _ := user_a.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_a.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Status != "OK" {
			t.Fatal("Failed user can access other user when FOLLOWERS and yes following")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when SUBSCRIBERS and not subscribers"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userd")
		user_a := User{Code: "usera"}
		session_a, _ := user_a.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_a.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Reason != "user is not subscriber" {
			t.Fatal("Failed user cannot access other user when SUBSCRIBERS and not subscribers")
		}
	})

	t.Run(fmt.Sprintf("Test user can access other user when SUBSCRIBERS and yes subscriber"), func(t *testing.T) {
		userToSee, _ := SelectUser("username", "userd")
		Db.Exec(`
			INSERT INTO subscribers (subscribefrom, subscribeto, createdat)
			VALUES ('usera', 'userd', current_timestamp);`)
		user_a := User{Code: "usera"}
		session_a, _ := user_a.CreateSession()
		req := httptest.NewRequest("GET", "/", nil)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session_a.Code, Expires: expiration}
		req.AddCookie(&cookie)
		status := CheckPrivacy(req, userToSee)
		if status.Status != "OK" {
			t.Fatal("Failed user can access other user when SUBSCRIBERS and yes subscriber")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestInstanciateTradeWs(t *testing.T) {

	// <test code>
	t.Run(fmt.Sprintf("Test instanciate ws from unvalid origin"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "http://notvalidorigin:9000")
		w := httptest.NewRecorder()
		ws, err := InstanciateTradeWs(w, req)
		if ws != nil && err.Error() != "CheckOrigin not accepted" {
			t.Fatal("Failed test instancate ws from unvalid origin")
		}
	})

	// <tear-down code>
}

func TestStartTradesWs(t *testing.T) {

	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			code, email, username, password, privacy,
			plan, createdat, updatedat)
		VALUES
			('usera', 'usera@r.r', 'usera', 'testpassword',
			'all', 'basic', current_timestamp, current_timestamp),
			('userb', 'userb@r.r', 'userb', 'testpassword',
			'private', 'basic', current_timestamp, current_timestamp);`)

	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(1000, 'USDC', 'USDC', 'usdc')`)

	Db.Exec(`
		INSERT INTO prices (
			createdat, coinid, price)
		VALUES
			(current_timestamp, 1, 80000),
			(current_timestamp, 1000, 1);`)

	Db.Exec(`
		INSERT INTO trades(
			code, usercode, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES
		('useratr', 'usera', current_timestamp, current_timestamp, 1000, 1, TRUE),
		('userbtr', 'userb', current_timestamp,current_timestamp, 1000, 1, TRUE);`)

	Db.Exec(`
		INSERT INTO subtrades(
			code, usercode, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES
		('userasub', 'usera', 'useratr', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART'),
		('userbsub', 'userb', 'userbtr', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test instanciate from invalid origin"), func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/usera/abc"}
		header := http.Header{}
		header.Set("Origin", "http://totallyinvalidorigin")
		_, _, err := websocket.DefaultDialer.Dial(u_new.String(), header)
		if err == nil {
			t.Fatal("Failed test instanciate from invalid origin")
		}

	})
	t.Run(fmt.Sprintf("Test wrong url path"), func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/wrong"}
		header := http.Header{}
		header.Set("Origin", "http://127.0.0.1")
		_, _, err := websocket.DefaultDialer.Dial(u_new.String(), header)
		if err == nil {
			t.Fatal("Failed test wrong url path")
		}
	})

	t.Run(fmt.Sprintf("Test invalid username"), func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/wrongusername/abc"}
		header := http.Header{}
		header.Set("Origin", "http://127.0.0.1")
		_, _, err := websocket.DefaultDialer.Dial(u_new.String(), header)
		if err == nil {
			t.Fatal("Failed test invalid username")
		}
	})

	t.Run(fmt.Sprintf("Test successfully receive initial snapshot"), func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/usera/sadkjfh"}
		header := http.Header{}
		header.Set("Origin", "http://127.0.0.1")
		ws, _, _ := websocket.DefaultDialer.Dial(u_new.String(), header)
		defer ws.Close()
		for {
			_, message, _ := ws.ReadMessage()
			trades_snapshot := TradesSnapshot{}
			json.Unmarshal([]byte(message), &trades_snapshot)
			if trades_snapshot.TotalReturnUsd != 15000 {
				t.Fatal("Failed successfully receive intial snapshot")
			}
			break
		}
	})

	t.Run(fmt.Sprintf("Test access PRIVATE user"), func(t *testing.T) {

		user_a := User{Code: "usera"}
		session_a, _ := user_a.CreateSession()

		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/userb/sadkjfh"}
		header := http.Header{}
		header.Set("Origin", "http://127.0.0.1")

		var cookies []*http.Cookie
		cookie := &http.Cookie{Name: "sessionId", Value: session_a.Code, MaxAge: 300}
		cookies = append(cookies, cookie)
		jar, _ := cookiejar.New(nil)
		urlObj, _ := url.Parse("http://127.0.0.1")
		jar.SetCookies(urlObj, cookies)

		dialer := websocket.DefaultDialer
		dialer.Jar = jar
		ws, _, _ := dialer.Dial(u_new.String(), header)
		defer ws.Close()
		for {
			_, message, _ := ws.ReadMessage()
			privacy_status := PrivacyStatus{}
			json.Unmarshal(message, &privacy_status)
			if privacy_status.Reason != "private" {
				t.Fatal("Failed test access PRIVATE user")
			}
			break
		}
	})

	t.Run(fmt.Sprintf("Test terminate ws"), func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/usera/sadkjfh"}
		header := http.Header{}
		header.Set("Origin", "http://127.0.0.1")
		ws, _, _ := websocket.DefaultDialer.Dial(u_new.String(), header)
		ws.Close()
		time.Sleep(time.Second)
		if len(trades_wss["usera"]) != 0 {
			t.Fatal("Failed test terminate ws")
		}
	})

	t.Run(fmt.Sprintf("Test receive snapshot after change"), func(t *testing.T) {
		go InstanciateActivityMonitor()
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/usera/sadkjfh"}
		header := http.Header{}
		header.Set("Origin", "http://127.0.0.1")
		ws, _, _ := websocket.DefaultDialer.Dial(u_new.String(), header)
		_, _, _ = ws.ReadMessage() // initial snapshot
		go func() {
			for {
				_, new_message, _ := ws.ReadMessage()
				trades_snapshot := TradesSnapshot{}
				json.Unmarshal([]byte(new_message), &trades_snapshot)
				if trades_snapshot.Trades[0].QtyBuys != 199 {
					t.Fatal("Failed test receive snapshot after update")
				}
				break
			}
		}()
		time.Sleep(time.Second)
		Db.Exec(`UPDATE subtrades SET quantity = 199 WHERE tradecode = 'useratr' RETURNING quantity;`)
		time.Sleep(time.Second)
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
}

func TestStartTradesWsIntegration(t *testing.T) {

	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			code, email, username, password, privacy,
			plan, createdat, updatedat)
		VALUES
			('usera', 'usera@r.r', 'usera', 'testpassword',
			'all', 'basic', current_timestamp, current_timestamp),
			('userb', 'userb@r.r', 'userb', 'testpassword',
			'all', 'basic', current_timestamp, current_timestamp),
			('userc', 'userc@r.r', 'userc', 'testpassword',
			'all', 'basic', current_timestamp, current_timestamp),
			('userd', 'userd@r.r', 'userd', 'testpassword',
			'all', 'basic', current_timestamp, current_timestamp);`)

	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(2, 'USDC', 'USDC', 'usdc'),
			(3, 'ETH', 'ETH', 'ethereum'),
			(4, 'DOT', 'DOT', 'polkadot'),
			(5, 'SOL', 'SOL', 'solana');`)

	Db.Exec(`
		INSERT INTO prices (
			createdat, coinid, price)
		VALUES
			(current_timestamp, 1, 65000),
			(current_timestamp, 2, 1),
			(current_timestamp, 3, 4500),
			(current_timestamp, 4, 50),
			(current_timestamp, 5, 200);`)

	Db.Exec(`
		INSERT INTO trades(
			code, usercode, createdat, updatedat,
			firstpair, secondpair, isopen)
		VALUES
		('usera', 'usera', current_timestamp, current_timestamp, 2, 1, TRUE),
		('userb', 'userb', current_timestamp, current_timestamp, 2, 4, TRUE),
		('userc', 'userc', current_timestamp, current_timestamp, 3, 1, TRUE),
		('userd', 'userd', current_timestamp, current_timestamp, 5, 2, TRUE);`)

	Db.Exec(`
		INSERT INTO subtrades(
			code, usercode, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES
		('usera', 'usera', 'usera', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART'),
		('userb', 'userb', 'userb', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART'),
		('userc', 'userc', 'userc', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART'),
		('userd', 'userd', 'userd', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test casual interaction without errors"), func(t *testing.T) {
		go InstanciateActivityMonitor()

		// usera
		server_a := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_a.Close()
		url_a := strings.TrimPrefix(server_a.URL, "http://")
		u_new_a := url.URL{Scheme: "ws", Host: url_a, Path: "get_trades/usera/requestusera"}
		header_a := http.Header{}
		header_a.Set("Origin", "http://127.0.0.1")
		ws_a, _, _ := websocket.DefaultDialer.Dial(u_new_a.String(), header_a)
		_, _, _ = ws_a.ReadMessage()

		// userb
		server_b := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_b.Close()
		url_b := strings.TrimPrefix(server_b.URL, "http://")
		u_new_b := url.URL{Scheme: "ws", Host: url_b, Path: "get_trades/usera/requestuserb"}
		header_b := http.Header{}
		header_b.Set("Origin", "http://127.0.0.1")
		ws_b, _, _ := websocket.DefaultDialer.Dial(u_new_b.String(), header_b)
		_, _, _ = ws_b.ReadMessage()

		// userc
		server_c := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_c.Close()
		url_c := strings.TrimPrefix(server_c.URL, "http://")
		u_new_c := url.URL{Scheme: "ws", Host: url_c, Path: "get_trades/usera/requestuserc"}
		header_c := http.Header{}
		header_c.Set("Origin", "http://127.0.0.1")
		ws_c, _, _ := websocket.DefaultDialer.Dial(u_new_c.String(), header_c)
		_, _, _ = ws_c.ReadMessage()

		// userc
		server_d := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_d.Close()
		url_d := strings.TrimPrefix(server_d.URL, "http://")
		u_new_d := url.URL{Scheme: "ws", Host: url_d, Path: "get_trades/usera/requestuserd"}
		header_d := http.Header{}
		header_d.Set("Origin", "http://127.0.0.1")
		ws_d, _, _ := websocket.DefaultDialer.Dial(u_new_d.String(), header_d)
		_, _, _ = ws_d.ReadMessage()

		fmt.Println(trades_wss["usera"])
		// WILL CONTINUE FROMHERE

	})
	// <tear-down code>

	fmt.Println("CIAO")
}
