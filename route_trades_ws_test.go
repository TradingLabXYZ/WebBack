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

type WebsocketServer struct {
	upgrader websocket.Upgrader
	addr     *string
	conn     *websocket.Conn
}

func TestCheckRelation(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, plan, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', 'basic', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'all', 'basic', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userd', 'all', 'basic', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO followers (followfrom, followto, createdat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', current_timestamp);`)
	Db.Exec(
		`INSERT INTO subscribers (subscribefrom, subscribeto, createdat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', current_timestamp);`)
	// <test code>
	t.Run(fmt.Sprintf("Test user not follow when userid is null"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: ""}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		if user_relation.IsFollower || user_relation.IsSubscriber {
			t.Fatal("Failed user not follow when userid is null")
		}
	})
	t.Run(fmt.Sprintf("Test usera follows userb"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		if !user_relation.IsFollower {
			t.Fatal("Failed test usera follows userb")
		}
	})
	t.Run(fmt.Sprintf("Test userc does not follows userb"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		if user_relation.IsFollower {
			t.Fatal("Failed test userc does not follows userb")
		}
	})
	t.Run(fmt.Sprintf("Test userb is subscribed to userc"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		if !user_relation.IsSubscriber {
			t.Fatal("Failed test userb is subscibed to userc")
		}
	})
	t.Run(fmt.Sprintf("Test usera is not subscribed to userc"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		if user_relation.IsSubscriber {
			t.Fatal("Failed test usera is not subscribed to userc")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestCheckPrivacy(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, plan, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', 'basic', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'private', 'basic', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userc', 'followers', 'basic', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', 'userd', 'subscribers', 'basic', current_timestamp, current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test user with privacy ALL is fully visibile"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		fmt.Println(user_relation.Privacy)
		if user_relation.Privacy.Status != "OK" {
			t.Fatal("Failed test user with privacy ALL is fully visibile")
		}
	})

	t.Run(fmt.Sprintf("Test user not authenticated try to access not ALL users"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: ""}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Status != "KO" {
			t.Fatal("Failed user not authenticated try to access not ALL users")
		}
	})

	t.Run(fmt.Sprintf("Test user PRIVATE always able to see its profile if authenticated"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Status != "OK" {
			t.Fatal("Failed user PRIVATE always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user FOLLOWERS always able to see its profile if authenticated"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Status != "OK" {
			t.Fatal("Failed user FOLLOWERS always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user SUBSCRIBERS always able to see its profile if authenticated"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Status != "OK" {
			t.Fatal("Failed user SUBSCRIBERS always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when PRIVATE"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Reason != "private" {
			t.Fatal("Failed user cannot access other user when PRIVATE")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when FOLLOWERS and not following"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Reason != "user is not follower" {
			t.Fatal("Failed user cannot access other user when FOLLOWERS and not following")
		}
	})

	t.Run(fmt.Sprintf("Test user can access other user when FOLLOWERS and yes following"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		Db.Exec(`
				INSERT INTO followers (followfrom, followto, createdat)
				VALUES ('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', current_timestamp);`)
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Status != "OK" {
			t.Fatal("Failed user can access other user when FOLLOWERS and yes following")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when SUBSCRIBERS and not subscribers"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Reason != "user is not subscriber" {
			t.Fatal("Failed user cannot access other user when SUBSCRIBERS and not subscribers")
		}
	})

	t.Run(fmt.Sprintf("Test user can access other user when SUBSCRIBERS and yes subscriber"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D")
		Db.Exec(`
				INSERT INTO subscribers (subscribefrom, subscribeto, createdat)
				VALUES ('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', current_timestamp);`)
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Status != "OK" {
			t.Fatal("Failed user can access other user when SUBSCRIBERS and yes subscriber")
		}
	})

	t.Run(fmt.Sprintf("Test user with privacy not legit"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		observed.Privacy = "random"
		user_relation := Relation{
			Observer: observer,
			Observed: observed,
		}
		user_relation.CheckRelation()
		user_relation.CheckPrivacy()
		if user_relation.Privacy.Reason != "unknown reason" {
			t.Fatal("Failed user with privacy not legit")
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
			wallet, username, privacy,
			plan, createdat, updatedat)
		VALUES
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera',
			'all', 'basic', current_timestamp, current_timestamp),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb',
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
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('useratr', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp, current_timestamp, 1000, 1),
		('userbtr', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', current_timestamp,current_timestamp, 1000, 1);`)

	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES
		('userasub', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART'),
		('userbsub', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userbtr', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test instanciate from invalid origin"), func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/abc"}
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
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/undefined"}
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

	t.Run(fmt.Sprintf("Test access PRIVATE user not logged in"), func(t *testing.T) {
		user_a := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session_a, _ := user_a.InsertSession()
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B/" + session_a.Code}
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
			snapshot := TradesSnapshot{}
			json.Unmarshal(message, &snapshot)
			if snapshot.PrivacyStatus.Reason != "private" {
				t.Fatal("Failed test access PRIVATE user")
			}
			break
		}
	})

	t.Run(fmt.Sprintf("Test terminate ws"), func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/undefined"}
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
		user_a := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		session_a, _ := user_a.InsertSession()
		go InstanciateActivityMonitor()
		s := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer s.Close()
		u := strings.TrimPrefix(s.URL, "http://")
		u_new := url.URL{Scheme: "ws", Host: u, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/" + session_a.Code}
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
			wallet, username, privacy,
			plan, createdat, updatedat)
		VALUES
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera',
			'all', 'basic', current_timestamp, current_timestamp),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb',
			'all', 'basic', current_timestamp, current_timestamp),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userc',
			'all', 'basic', current_timestamp, current_timestamp),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', 'userd',
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
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('usera', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp, current_timestamp, 2, 1),
		('userb', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', current_timestamp, current_timestamp, 2, 4),
		('userc', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', current_timestamp, current_timestamp, 3, 1),
		('userd', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', current_timestamp, current_timestamp, 5, 2);`)

	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES
		('usera', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART'),
		('userb', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART'),
		('userc', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userc', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART'),
		('userd', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', 'userd', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART');`)

	// <test code>
	t.Run(fmt.Sprintf("Test casual interaction without errors"), func(t *testing.T) {
		go InstanciateActivityMonitor()

		// test 1
		server_a := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_a.Close()
		url_a := strings.TrimPrefix(server_a.URL, "http://")
		u_new_a := url.URL{Scheme: "ws", Host: url_a, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/undefined"}
		header_a := http.Header{}
		header_a.Set("Origin", "http://127.0.0.1")
		ws_a, _, _ := websocket.DefaultDialer.Dial(u_new_a.String(), header_a)
		_, _, _ = ws_a.ReadMessage()

		// test 2
		server_b := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_b.Close()
		url_b := strings.TrimPrefix(server_b.URL, "http://")
		u_new_b := url.URL{Scheme: "ws", Host: url_b, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/undefined"}
		header_b := http.Header{}
		header_b.Set("Origin", "http://127.0.0.1")
		ws_b, _, _ := websocket.DefaultDialer.Dial(u_new_b.String(), header_b)
		_, _, _ = ws_b.ReadMessage()

		// test 3
		user_c := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C"}
		session_c, _ := user_c.InsertSession()
		server_c := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_c.Close()
		url_c := strings.TrimPrefix(server_c.URL, "http://")
		u_new_c := url.URL{Scheme: "ws", Host: url_c, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/" + session_c.Code}
		header_c := http.Header{}
		header_c.Set("Origin", "http://127.0.0.1")
		ws_c, _, _ := websocket.DefaultDialer.Dial(u_new_c.String(), header_c)
		_, _, _ = ws_c.ReadMessage()

		// test 4
		user_d := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D"}
		session_d, _ := user_d.InsertSession()
		server_d := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_d.Close()
		url_d := strings.TrimPrefix(server_d.URL, "http://")
		u_new_d := url.URL{Scheme: "ws", Host: url_d, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A/" + session_d.Code}
		header_d := http.Header{}
		header_d.Set("Origin", "http://127.0.0.1")
		ws_d, _, _ := websocket.DefaultDialer.Dial(u_new_d.String(), header_d)
		_, _, _ = ws_d.ReadMessage()

		// This test needs to be expanded way further
		if len(trades_wss) != 1 {
			t.Fatal("Test failed casual interaction without errors")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
}

// This test creates problem, need to investigate
/* func TestVisitPrivateUserNotCreareWebSocket(t *testing.T) {
	// <setup code>

	trades_wss = make(map[string][]WsTrade)

	Db.Exec(
		`INSERT INTO users (
			wallet, username, privacy,
			plan, createdat, updatedat)
		VALUES
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera',
			'all', 'basic', current_timestamp, current_timestamp),
			('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb',
			'private', 'basic', current_timestamp, current_timestamp);`)

	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(2, 'USDC', 'USDC', 'usdc');`)
	Db.Exec(`
		INSERT INTO prices (
			createdat, coinid, price)
		VALUES
			(current_timestamp, 1, 65000),
			(current_timestamp, 2, 1);`)

	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('userb', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', current_timestamp, current_timestamp, 2, 1);`)
	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES
		('userb', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', current_timestamp, current_timestamp, 'BUY', 1, 10000, 10000, 'TESTART')`)

	// <test code>
	t.Run(fmt.Sprintf("Test web socket not create when user is private"), func(t *testing.T) {
		go InstanciateActivityMonitor()

		server_a := httptest.NewServer(http.HandlerFunc(StartTradesWs))
		defer server_a.Close()
		url_a := strings.TrimPrefix(server_a.URL, "http://")
		u_new_a := url.URL{Scheme: "ws", Host: url_a, Path: "get_trades/0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B/undefined"}
		header_a := http.Header{}
		header_a.Set("Origin", "http://127.0.0.1")
		_, _, _ = websocket.DefaultDialer.Dial(u_new_a.String(), header_a)

		// LEAVE LIKE THIS OTHERWISE RACE CONDITION (not sure why)
		temp_trade_wss := trades_wss
		counter := 0
		for _ = range temp_trade_wss {
			counter++
		}

		if counter != 0 {
			t.Fatal("Test failed web socket not create when user is private")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
} */
